package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tnphucccc/mangahub/pkg/database"
	"github.com/tnphucccc/mangahub/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

const DataFile = "data/manga.json"

func main() {
	// Connect to database
	dbConfig := database.DefaultConfig()
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	fmt.Println("Seeding database...")

	// Seed users
	if err := seedUsers(db); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	// Seed manga from JSON file
	if err := seedMangaFromJSON(db); err != nil {
		log.Fatalf("Failed to seed manga: %v", err)
	}

	// Seed user progress
	if err := seedUserProgress(db); err != nil {
		log.Fatalf("Failed to seed user progress: %v", err)
	}

	fmt.Println("Database seeding completed successfully!")
}

func seedUsers(db *sql.DB) error {
	fmt.Println("Seeding users...")

	users := []struct {
		username string
		email    string
		password string
	}{
		{"testuser", "testuser@example.com", "password123"},
		{"alice", "alice@example.com", "alice123"},
		{"bob", "bob@example.com", "bob123"},
	}

	for _, u := range users {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Generate UUID (simple version for now)
		userID := fmt.Sprintf("user-%s", u.username)

		// Insert user
		_, err = db.Exec(
			`
			INSERT OR IGNORE INTO users (id, username, email, password_hash)
			VALUES (?, ?, ?, ?)
		`,
			userID, u.username, u.email, string(hashedPassword))

		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", u.username, err)
		}

		fmt.Printf("  Created user: %s (%s)\n", u.username, u.email)
	}

	return nil
}

func seedMangaFromJSON(db *sql.DB) error {
	fmt.Printf("Seeding manga from %s...\n", DataFile)

	content, err := os.ReadFile(DataFile)
	if err != nil {
		return fmt.Errorf("failed to read data file: %w", err)
	}

	var mangaList []models.Manga
	if err := json.Unmarshal(content, &mangaList); err != nil {
		return fmt.Errorf("failed to unmarshal manga list: %w", err)
	}

	count := 0
	for _, manga := range mangaList {
		// Marshal genres to JSON for DB
		genresJSON, err := manga.MarshalGenres()
		if err != nil {
			log.Printf("Failed to marshal genres for %s: %v", manga.Title, err)
			continue
		}

		// Insert manga
		_, err = db.Exec(`
			INSERT OR IGNORE INTO manga (id, title, author, genres, status, total_chapters, description, cover_image_url)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
			manga.ID, manga.Title, manga.Author, genresJSON, manga.Status, manga.TotalChapters, manga.Description, manga.CoverImageURL)

		if err != nil {
			log.Printf("Failed to insert manga %s: %v", manga.Title, err)
			continue
		}
		count++
	}

	fmt.Printf("  Imported %d manga from JSON file\n", count)
	return nil
}

func seedUserProgress(db *sql.DB) error {
	fmt.Println("Seeding user progress...")

	// IDs must match what's in the generated JSON file
	userProgress := []struct {
		userID         string
		mangaID        string
		currentChapter int
		status         models.ReadingStatus
	}{
		{"user-testuser", "my-robot-has-been-acting-strange-lately", 10, models.ReadingStatusReading},
		{"user-testuser", "versatile-mage", 500, models.ReadingStatusCompleted},
		{"user-alice", "pet", 55, models.ReadingStatusCompleted},
		{"user-alice", "baby-steps", 0, models.ReadingStatusPlanToRead},
		{"user-bob", "yakuza-reincarnation", 20, models.ReadingStatusReading},
	}

	for _, up := range userProgress {
		// Verify manga exists first to avoid FK errors (optional but good)
		var exists int
		err := db.QueryRow("SELECT 1 FROM manga WHERE id = ?", up.mangaID).Scan(&exists)
		if err != nil {
			log.Printf("  Skipping progress for %s: manga not found", up.mangaID)
			continue
		}

		_, err = db.Exec(
			`
			INSERT OR IGNORE INTO user_progress (user_id, manga_id, current_chapter, status)
			VALUES (?, ?, ?, ?)
		`,
			up.userID, up.mangaID, up.currentChapter, up.status)

		if err != nil {
			return fmt.Errorf("failed to insert user progress: %w", err)
		}
	}

	fmt.Println("  Created sample user progress entries")
	return nil
}