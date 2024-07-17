package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// RunMigrations runs the database migrations using the provided SQL database connection.
// It initializes the migration instance and applies any pending migrations.
// If an error occurs during the migration process, it will log the error and exit the program.
func RunMigrations(db *sql.DB) {
	dbType := os.Getenv("DB_TYPE")
	var driver database.Driver
	var err error

	migrationPath := fmt.Sprintf("file://migrations/%s", dbType)

	switch dbType {
	case "sqlite":
		driver, err = sqlite.WithInstance(db, &sqlite.Config{})
	case "postgres":
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	default:
		log.Fatalf("Unsupported database type: %v", dbType)
	}

	if err != nil {
		log.Fatalf("Could not create database driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		dbType, driver)
	if err != nil {
		log.Fatalf("Could not initialize migration: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could not apply migrations: %v", err)
	}
}
