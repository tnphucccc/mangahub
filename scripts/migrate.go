package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tnphucccc/mangahub/pkg/database"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/migrate.go [up|down]")
		fmt.Println("  up   - Apply all pending migrations")
		fmt.Println("  down - Rollback the last migration")
		os.Exit(1)
	}

	command := os.Args[1]

	// Get absolute path to migrations directory
	migrationsPath := "./migrations"
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// Connect to database
	dbConfig := database.DefaultConfig()

	// Ensure data directory exists
	dataDir := filepath.Dir(dbConfig.Path)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Create migrator
	migrator := database.NewMigrator(db, absPath)

	// Execute command
	switch command {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "down":
		if err := migrator.Down(); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Usage: go run scripts/migrate.go [up|down]")
		os.Exit(1)
	}
}
