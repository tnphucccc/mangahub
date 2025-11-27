package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a single database migration
type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

// Migrator handles database migrations
type Migrator struct {
	db             *sql.DB
	migrationsPath string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB, migrationsPath string) *Migrator {
	return &Migrator{
		db:             db,
		migrationsPath: migrationsPath,
	}
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	// Create migrations table if not exists
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Load migrations
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedVersions, err := m.getAppliedVersions()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if _, applied := appliedVersions[migration.Version]; !applied {
			if err := m.applyMigration(migration); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
			}
			fmt.Printf("Applied migration: %03d_%s\n", migration.Version, migration.Name)
		}
	}

	fmt.Println("All migrations applied successfully")
	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	// Get last applied migration
	var version int
	var name string
	err := m.db.QueryRow(`
		SELECT version, name FROM schema_migrations
		ORDER BY version DESC LIMIT 1
	`).Scan(&version, &name)

	if err == sql.ErrNoRows {
		fmt.Println("No migrations to roll back")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Load migration
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Find and rollback migration
	for _, migration := range migrations {
		if migration.Version == version {
			if err := m.rollbackMigration(migration); err != nil {
				return fmt.Errorf("failed to rollback migration %d: %w", version, err)
			}
			fmt.Printf("Rolled back migration: %03d_%s\n", version, name)
			return nil
		}
	}

	return fmt.Errorf("migration %d not found", version)
}

// createMigrationsTable creates the schema_migrations tracking table
func (m *Migrator) createMigrationsTable() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

// loadMigrations loads all migration files from the migrations directory
func (m *Migrator) loadMigrations() ([]Migration, error) {
	var migrations []Migration

	// Read migration files
	files, err := os.ReadDir(m.migrationsPath)
	if err != nil {
		return nil, err
	}

	// Group files by version
	migrationMap := make(map[int]*Migration)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		var version int
		var name string
		var direction string

		// Parse filename: 001_create_users_table.up.sql
		parts := strings.SplitN(filename, "_", 2)
		if len(parts) != 2 {
			continue
		}

		fmt.Sscanf(parts[0], "%d", &version)

		nameParts := strings.Split(parts[1], ".")
		if len(nameParts) < 3 {
			continue
		}

		name = nameParts[0]
		direction = nameParts[1]

		// Read SQL content
		content, err := os.ReadFile(filepath.Join(m.migrationsPath, filename))
		if err != nil {
			return nil, err
		}

		// Get or create migration
		if _, exists := migrationMap[version]; !exists {
			migrationMap[version] = &Migration{
				Version: version,
				Name:    name,
			}
		}

		// Set up or down SQL
		switch direction {
		case "up":
			migrationMap[version].UpSQL = string(content)
		case "down":
			migrationMap[version].DownSQL = string(content)
		}
	}

	// Convert map to slice
	for _, migration := range migrationMap {
		migrations = append(migrations, *migration)
	}

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedVersions returns a map of applied migration versions
func (m *Migrator) getAppliedVersions() (map[int]bool, error) {
	rows, err := m.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions[version] = true
	}

	return versions, rows.Err()
}

// applyMigration applies a single migration
func (m *Migrator) applyMigration(migration Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return err
	}

	// Record migration
	_, err = tx.Exec(`
		INSERT INTO schema_migrations (version, name) VALUES (?, ?)
	`, migration.Version, migration.Name)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// rollbackMigration rolls back a single migration
func (m *Migrator) rollbackMigration(migration Migration) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute rollback SQL
	if _, err := tx.Exec(migration.DownSQL); err != nil {
		return err
	}

	// Remove migration record
	_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = ?", migration.Version)
	if err != nil {
		return err
	}

	return tx.Commit()
}
