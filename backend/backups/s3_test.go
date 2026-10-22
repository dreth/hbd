package backups

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type recordingUploader struct {
	input *s3.PutObjectInput
	body  []byte
	err   error
}

func (u *recordingUploader) PutObject(_ context.Context, input *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	u.input = input
	var err error
	u.body, err = io.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}
	return &s3.PutObjectOutput{}, u.err
}
func TestBackupUploadPreservesObjectAndBytes(t *testing.T) {
	body := []byte("SQLite format 3\x00synthetic backup verification")
	path := filepath.Join(t.TempDir(), "hbd.db")
	if err := os.WriteFile(path, body, 0600); err != nil {
		t.Fatal(err)
	}
	client := &recordingUploader{}
	if err := uploadDBWithClient(context.Background(), client, path, "synthetic-bucket", "hbd.db"); err != nil {
		t.Fatal(err)
	}
	if string(client.body) != string(body) || *client.input.Bucket != "synthetic-bucket" || *client.input.Key != "hbd.db" || *client.input.ContentLength != int64(len(body)) || *client.input.ContentType != "application/x-sqlite3" {
		t.Fatal("Backup content or object metadata changed")
	}
}
func TestBackupUploadPropagatesProviderFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "hbd.db")
	if err := os.WriteFile(path, []byte("synthetic"), 0600); err != nil {
		t.Fatal(err)
	}
	expected := errors.New("synthetic provider failure")
	err := uploadDBWithClient(context.Background(), &recordingUploader{err: expected}, path, "synthetic-bucket", "hbd.db")
	if !errors.Is(err, expected) {
		t.Fatal("Provider failure was hidden")
	}
}
func TestMissingDatabaseIsNotUploaded(t *testing.T) {
	client := &recordingUploader{}
	if err := uploadDBWithClient(context.Background(), client, filepath.Join(t.TempDir(), "missing.db"), "synthetic-bucket", "hbd.db"); err == nil || client.input != nil {
		t.Fatal("Missing database was uploaded")
	}
}
