package env

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Set the MASTER_KEY and DATABASE_URL environment variables
var DB *sql.DB = db()
var MK string = mk()

func loadDotenv() {
	// Load the .env file
	env := os.Getenv("ENVIRONMENT")
	if env == "development" || env == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
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
	// Load the DATABASE_URL from the environment
	loadDotenv()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	var err error
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
