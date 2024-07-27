package env

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Set the MASTER_KEY and DATABASE_URL environment variables
var DB *sql.DB = db()
var MK string = mk()
var CDFE string = customDomainFrontend()
var CDBE string = customDomainBackend()

func DBType() string {
	loadDotenv()
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "postgres" // default to postgres
	}
	return dbType
}

func loadDotenv() {
	// Load the .env file
	env := os.Getenv("ENVIRONMENT")
	if env == "development" || env == "" {
		godotenv.Load()
	}
}

func mk() string {
	// Load the MASTER_KEY from the environment
	loadDotenv()
	masterKey := os.Getenv("MASTER_KEY")
	if masterKey == "" {
		log.Fatal("MASTER_KEY environment variable not set")
	}

	return masterKey
}

func db() *sql.DB {
	// Load the DATABASE_URL and DB_TYPE from the environment
	loadDotenv()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	var db *sql.DB
	var err error

	if DBType() == "sqlite" {
		db, err = sql.Open("sqlite3", databaseURL)
	} else {
		db, err = sql.Open("postgres", databaseURL)
	}

	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Custom domain for CORS
func customDomainFrontend() string {
	loadDotenv()
	customDomain := os.Getenv("CUSTOM_DOMAIN_FRONTEND")
	if customDomain == "" {
		customDomain = "http://hbd.lotiguere.com"
	}
	return customDomain
}

// Custom domain for CORS
func customDomainBackend() string {
	loadDotenv()
	customDomain := os.Getenv("CUSTOM_DOMAIN_BACKEND")
	if customDomain == "" {
		customDomain = "http://hbd-api.lotiguere.com"
	}
	return customDomain
}
