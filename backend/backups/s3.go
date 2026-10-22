package backups

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3EnvVarspresent() bool {
	for _, name := range []string{"HBD_BUCKET_NAME", "HBD_BUCKET_REGION", "HBD_USER_ACCESS_KEY_ID", "HBD_USER_SECRET_ACCESS_KEY"} {
		if os.Getenv(name) == "" {
			log.Printf("%s is not set", name)
			return false
		}
	}
	return true
}

func BackupDBToS3() {
	if os.Getenv("HBD_ENABLE_BACKUP") != "true" {
		return
	}
	if !S3EnvVarspresent() {
		log.Println("Backups are enabled, but S3 environment variables not set")
		return
	}
	databaseURL := os.Getenv("DATABASE_URL")
	objectKey := databaseURL[strings.LastIndex(databaseURL, "/")+1:]
	if err := uploadDBToS3(databaseURL, os.Getenv("HBD_BUCKET_NAME"), objectKey); err != nil {
		log.Printf("Failed to upload database: %v", err)
		return
	}
	log.Println("Database backup complete")
}

type objectUploader interface {
	PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func uploadDBToS3(dbPath, bucketName, objectKey string) error {
	config := aws.Config{
		Region: os.Getenv("HBD_BUCKET_REGION"),
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("HBD_USER_ACCESS_KEY_ID"), os.Getenv("HBD_USER_SECRET_ACCESS_KEY"), ""),
		HTTPClient: &http.Client{Timeout: 60 * time.Second},
	}
	return uploadDBWithClient(context.Background(), s3.NewFromConfig(config), dbPath, bucketName, objectKey)
}

func uploadDBWithClient(ctx context.Context, client objectUploader, dbPath, bucketName, objectKey string) error {
	file, err := os.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database file: %w", err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName), Key: aws.String(objectKey), Body: file,
		ContentLength: aws.Int64(info.Size()), ContentType: aws.String("application/x-sqlite3"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}
	return nil
}
