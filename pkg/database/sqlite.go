package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Config holds database configuration
type Config struct {
	Path            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DefaultConfig returns default database configuration
func DefaultConfig() Config {
	return Config{
		Path:            "./data/mangahub.db",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

// Connect establishes a connection to the SQLite database
func Connect(config Config) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("sqlite3", config.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys (SQLite specific)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return db, nil
}

// Close gracefully closes the database connection
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// HealthCheck verifies database connectivity
func HealthCheck(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	return db.Ping()
}
