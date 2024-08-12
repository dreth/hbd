package backups

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3EnvVarspresent checks if the required environment variables for S3 are set.
func S3EnvVarspresent() bool {
	// Check if the required environment variables are set
	if os.Getenv("HBD_BUCKET_NAME") == "" {
		log.Println("HBD_BUCKET_NAME is not set")
		return false
	}

	if os.Getenv("HBD_BUCKET_REGION") == "" {
		log.Println("HBD_BUCKET_REGION is not set")
		return false
	}

	if os.Getenv("HBD_USER_ACCESS_KEY_ID") == "" {
		log.Println("HBD_USER_ACCESS_KEY_ID is not set")
		return false
	}

	if os.Getenv("HBD_USER_SECRET_ACCESS_KEY") == "" {
		log.Println("HBD_USER_SECRET_ACCESS_KEY is not set")
		return false
	}

	return true
}

// BackupDBToS3 uploads the SQLite database to S3.
func BackupDBToS3() {
	// Check if backups are enabled
	if os.Getenv("ENABLE_BACKUP") != "true" {
		return
	}

	if !S3EnvVarspresent() {
		log.Println("S3 environment variables not set")
		return
	}

	// Get the database URL from the environment
	databaseURL := os.Getenv("DATABASE_URL")
	bucketName := os.Getenv("HBD_BUCKET_NAME")
	objectKey := databaseURL[strings.LastIndex(databaseURL, "/")+1:]

	err := uploadDBToS3(databaseURL, bucketName, objectKey)
	if err != nil {
		log.Printf("Failed to upload database: %v", err)
	}
}

// uploadDBToS3 uploads the database file to the specified S3 bucket and key.
func uploadDBToS3(dbPath, bucketName, objectKey string) error {
	// Create a new S3 session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("HBD_BUCKET_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("HBD_USER_ACCESS_KEY_ID"),
			os.Getenv("HBD_USER_SECRET_ACCESS_KEY"),
			"",
		),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// Create an S3 client
	svc := s3.New(sess)

	// Open the database file
	file, err := os.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database file: %w", err)
	}
	defer file.Close()

	// Get the file size and content type (optional, could set a default)
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()

	// Upload the file to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          file,
		ContentLength: aws.Int64(fileSize),
		ContentType:   aws.String("application/x-sqlite3"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}
