package observability

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type recordingTransport struct {
	mu     sync.Mutex
	events []*sentry.Event
}

func (t *recordingTransport) Configure(sentry.ClientOptions)        {}
func (t *recordingTransport) Flush(time.Duration) bool              { return true }
func (t *recordingTransport) FlushWithContext(context.Context) bool { return true }
func (t *recordingTransport) Close()                                {}
func (t *recordingTransport) SendEvent(event *sentry.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, event)
}

func (t *recordingTransport) byType(eventType string) []*sentry.Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	var matched []*sentry.Event
	for _, event := range t.events {
		if event.Type == eventType {
			matched = append(matched, event)
		}
	}
	return matched
}

func TestTracesSampleRate(t *testing.T) {
	cases := map[string]float64{"": 0.2, "abc": 0.2, "-1": 0.2, "1.5": 0.2, "0": 0, "1": 1, " 0.05 ": 0.05}
	for input, want := range cases {
		if got := TracesSampleRate(input); got != want {
			t.Errorf("TracesSampleRate(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestScrub(t *testing.T) {
	token := "123456789:AAH3x_example-token-value-0123456789"
	cases := map[string]string{
		`Post "https://api.telegram.org/bot` + token + `/sendMessage": dial tcp`: `Post "https://api.telegram.org/bot[Filtered]/sendMessage": dial tcp`,
		"GET https://example.com/a?token=abc&x=1 failed":                         "GET https://example.com/a?[Filtered] failed",
		"connect https://user:secret@db.example.com/x":                           "connect https://[Filtered]@db.example.com/x",
		"plain message": "plain message",
	}
	for input, want := range cases {
		if got := Scrub(input); got != want {
			t.Errorf("Scrub(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestParseLogLine(t *testing.T) {
	cases := []struct {
		line    string
		message string
		isError bool
	}{
		{"2026/11/09 15:04:05 Error: invalid token\n", "Error: invalid token", true},
		{"2026/11/09 15:04:05 Failed to upload database: boom\n", "Failed to upload database: boom", true},
		{"2026/11/09 15:04:05 Database backup complete\n", "Database backup complete", false},
		{"\n", "", false},
	}
	for _, tc := range cases {
		message, isError := parseLogLine(tc.line)
		if message != tc.message || isError != tc.isError {
			t.Errorf("parseLogLine(%q) = %q, %v; want %q, %v", tc.line, message, isError, tc.message, tc.isError)
		}
	}
}

func TestInitWithoutDSNChangesNothing(t *testing.T) {
	t.Setenv("SENTRY_DSN", "  ")
	before := log.Writer()
	if Init() {
		t.Fatal("Init() enabled reporting without a DSN")
	}
	if log.Writer() != before {
		t.Fatal("Init() replaced the standard logger output without a DSN")
	}
	CaptureRequestError(nil, http.StatusInternalServerError, errors.New("ignored"))
	CaptureJobError("job", errors.New("ignored"))
}

func TestReportsServerErrorsTracesAndLogs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("SENTRY_DSN", "https://public@glitchtip.example/1")
	t.Setenv("SENTRY_ENVIRONMENT", "test")
	t.Setenv("SENTRY_TRACES_SAMPLE_RATE", "1")
	originalOutput := log.Writer()
	var stderr bytes.Buffer
	log.SetOutput(&stderr)
	transport := &recordingTransport{}
	t.Cleanup(func() {
		log.SetOutput(originalOutput)
		enabled = false
		sentry.CurrentHub().BindClient(nil)
	})

	if !initWithTransport(transport) {
		t.Fatal("initWithTransport() = false with a DSN")
	}

	router := gin.New()
	router.Use(Gin())
	router.GET("/api/health", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.GET("/api/fail", func(c *gin.Context) {
		CaptureRequestError(c, http.StatusBadRequest, errors.New("client error"))
		CaptureRequestError(c, http.StatusInternalServerError,
			errors.New("send failed https://api.telegram.org/bot123456789:AAH3x_example-token-value-0123456789/sendMessage?chat=1"))
		c.Status(http.StatusInternalServerError)
	})
	for _, path := range []string{"/api/health", "/api/fail?secret=1"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		request.Header.Set("Authorization", "Bearer jwt")
		request.Header.Set("Cookie", "session=abc")
		router.ServeHTTP(httptest.NewRecorder(), request)
	}
	log.Println("Error: something broke")
	log.Println("Database backup complete")
	sentry.Flush(time.Second)

	errorsSent := transport.byType("")
	if len(errorsSent) != 1 {
		t.Fatalf("sent %d error events, want 1 (only the 5xx)", len(errorsSent))
	}
	event := errorsSent[0]
	if got := event.Exception[len(event.Exception)-1].Value; got != "send failed https://api.telegram.org/bot[Filtered]/sendMessage?[Filtered]" {
		t.Errorf("exception value not scrubbed: %q", got)
	}
	if event.Tags["service"] != "hbd-api" || event.Environment != "test" {
		t.Errorf("unexpected tags/environment: %v %q", event.Tags, event.Environment)
	}
	if event.Request == nil || event.Request.QueryString != "" || event.Request.Cookies != "" {
		t.Errorf("request data not private: %+v", event.Request)
	}
	for key, value := range event.Request.Headers {
		if strings.EqualFold(key, "Authorization") && value != "[Filtered]" {
			t.Errorf("authorization header leaked: %q", value)
		}
	}
	if event.User.IPAddress != "" {
		t.Errorf("user IP collected: %q", event.User.IPAddress)
	}

	var transactions []string
	for _, tx := range transport.byType("transaction") {
		transactions = append(transactions, tx.Transaction)
	}
	if len(transactions) != 1 || transactions[0] != "GET /api/fail" {
		t.Errorf("transactions = %v, want only GET /api/fail (health checks ignored)", transactions)
	}

	var logs []string
	for _, batch := range transport.byType("log") {
		for _, entry := range batch.Logs {
			logs = append(logs, string(entry.Level)+":"+entry.Body)
		}
	}
	if strings.Join(logs, "|") != "error:Error: something broke|info:Database backup complete" {
		t.Errorf("logs = %v", logs)
	}
	if !strings.Contains(stderr.String(), "Database backup complete") {
		t.Error("standard log output no longer reaches stderr")
	}
}
