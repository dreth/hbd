// Package observability reports API errors, traces and logs to GlitchTip (Sentry-compatible).
// Everything here is a no-op unless SENTRY_DSN is set.
package observability

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/attribute"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

const (
	// DefaultTracesSampleRate applies when SENTRY_TRACES_SAMPLE_RATE is unset or invalid.
	DefaultTracesSampleRate = 0.2
	service                 = "hbd-api"
	flushTimeout            = 2 * time.Second
)

var enabled bool

var (
	// Telegram bot tokens appear in request URLs (and so in net/http errors).
	telegramTokenPattern = regexp.MustCompile(`\d{5,}:[A-Za-z0-9_-]{30,}`)
	urlUserInfoPattern   = regexp.MustCompile(`(https?://)[^/\s@"']+@`)
	urlQueryPattern      = regexp.MustCompile(`(https?://[^\s?#"']+)\?[^\s#"']*`)
	stdLogPrefixPattern  = regexp.MustCompile(`^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}(\.\d+)? `)
	errorLinePattern     = regexp.MustCompile(`(?i)^(error|failed|fatal|panic)\b`)
)

// Init configures reporting from SENTRY_DSN, SENTRY_ENVIRONMENT, SENTRY_RELEASE and
// SENTRY_TRACES_SAMPLE_RATE, and mirrors the standard logger to Sentry logs. It returns false and
// changes nothing when SENTRY_DSN is unset.
func Init() bool {
	return initWithTransport(nil)
}

// initWithTransport lets tests replace the HTTP transport; nil keeps the SDK default.
func initWithTransport(transport sentry.Transport) bool {
	dsn := strings.TrimSpace(os.Getenv("SENTRY_DSN"))
	if dsn == "" {
		return false
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:                   dsn,
		Environment:           firstNonEmpty(os.Getenv("SENTRY_ENVIRONMENT"), os.Getenv("ENVIRONMENT"), "production"),
		Release:               strings.TrimSpace(os.Getenv("SENTRY_RELEASE")),
		EnableTracing:         true,
		TracesSampleRate:      TracesSampleRate(os.Getenv("SENTRY_TRACES_SAMPLE_RATE")),
		IgnoreTransactions:    []string{`^GET /api/health$`},
		AttachStacktrace:      true,
		DataCollection:        PrivateDataCollection(),
		BeforeSend:            scrubEvent,
		BeforeSendTransaction: scrubEvent,
		BeforeSendLog:         scrubLog,
		Transport:             transport,
	})
	if err != nil {
		log.Println("GlitchTip reporting disabled: invalid Sentry configuration")
		return false
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) { scope.SetTag("service", service) })
	logger := sentry.NewLogger(context.Background())
	logger.SetAttributes(attribute.String("service.name", service))
	log.SetOutput(io.MultiWriter(log.Writer(), &logWriter{logger: logger}))
	enabled = true
	return true
}

// PrivateDataCollection is the equivalent of send_default_pii off: no user/IP data, cookies,
// bodies or query strings, and only a few harmless request headers.
func PrivateDataCollection() *sentry.DataCollection {
	return &sentry.DataCollection{
		UserInfo: sentry.Set(false),
		Cookies:  &sentry.KeyValueCollectionBehavior{Mode: sentry.CollectionOff},
		HTTPHeaders: &sentry.HeaderCollectionConfig{
			Request: &sentry.KeyValueCollectionBehavior{
				Mode:  sentry.CollectionAllowList,
				Terms: []string{"user-agent", "content-type", "accept", "host"},
			},
			Response: &sentry.KeyValueCollectionBehavior{Mode: sentry.CollectionOff},
		},
		HTTPBodies:  []sentry.BodyType{},
		QueryParams: &sentry.KeyValueCollectionBehavior{Mode: sentry.CollectionOff},
	}
}

// Gin traces each request and reports panics, then re-panics so gin's recovery still answers 500.
func Gin() gin.HandlerFunc {
	return sentrygin.New(sentrygin.Options{Repanic: true})
}

// CaptureRequestError reports a handled request error when it is answered with a 5xx status.
func CaptureRequestError(c *gin.Context, statusCode int, err error) {
	if !enabled || err == nil || statusCode < http.StatusInternalServerError {
		return
	}
	if hub := sentrygin.GetHubFromContext(c); hub != nil {
		hub.CaptureException(err)
	}
}

// CaptureJobError reports a failure inside a scheduled job.
func CaptureJobError(job string, err error) {
	if !enabled || err == nil {
		return
	}
	hub := sentry.CurrentHub().Clone()
	hub.ConfigureScope(func(scope *sentry.Scope) { scope.SetTag("cron.job", job) })
	hub.CaptureException(err)
}

// Job wraps a cron job so a panic is reported before it propagates exactly as before.
func Job(name string, job func()) func() {
	return func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				if enabled {
					hub := sentry.CurrentHub().Clone()
					hub.ConfigureScope(func(scope *sentry.Scope) { scope.SetTag("cron.job", name) })
					hub.Recover(recovered)
					hub.Flush(flushTimeout)
				}
				panic(recovered)
			}
		}()
		job()
	}
}

// TracesSampleRate parses SENTRY_TRACES_SAMPLE_RATE, falling back to the default when invalid.
func TracesSampleRate(value string) float64 {
	rate, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || rate < 0 || rate > 1 {
		return DefaultTracesSampleRate
	}
	return rate
}

// Scrub removes Telegram bot tokens, URL credentials and URL query strings from text.
func Scrub(text string) string {
	text = telegramTokenPattern.ReplaceAllString(text, "[Filtered]")
	text = urlUserInfoPattern.ReplaceAllString(text, "$1[Filtered]@")
	return urlQueryPattern.ReplaceAllString(text, "$1?[Filtered]")
}

func scrubEvent(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
	event.Message = Scrub(event.Message)
	for i := range event.Exception {
		event.Exception[i].Value = Scrub(event.Exception[i].Value)
	}
	if event.Request != nil {
		event.Request.URL = Scrub(event.Request.URL)
		event.Request.QueryString = ""
	}
	for _, breadcrumb := range event.Breadcrumbs {
		breadcrumb.Message = Scrub(breadcrumb.Message)
	}
	return event
}

func scrubLog(entry *sentry.Log) *sentry.Log {
	entry.Body = Scrub(entry.Body)
	return entry
}

// logWriter forwards standard-library log lines to Sentry logs, as errors when they look like one.
type logWriter struct {
	logger sentry.Logger
}

func (w *logWriter) Write(p []byte) (int, error) {
	message, isError := parseLogLine(string(p))
	switch {
	case message == "":
	case isError:
		w.logger.Error().Emit(message)
	default:
		w.logger.Info().Emit(message)
	}
	return len(p), nil
}

func parseLogLine(line string) (message string, isError bool) {
	message = strings.TrimSpace(stdLogPrefixPattern.ReplaceAllString(line, ""))
	return message, errorLinePattern.MatchString(message)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
