package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/tnphucccc/mangahub/pkg/database"
	"github.com/tnphucccc/mangahub/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

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

	// Seed manga
	if err := seedManga(db); err != nil {
		log.Fatalf("Failed to seed manga: %v", err)
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
		_, err = db.Exec(`
			INSERT OR IGNORE INTO users (id, username, email, password_hash)
			VALUES (?, ?, ?, ?)
		`, userID, u.username, u.email, string(hashedPassword))

		if err != nil {
			return fmt.Errorf("failed to insert user %s: %w", u.username, err)
		}

		fmt.Printf("  Created user: %s (%s)\n", u.username, u.email)
	}

	return nil
}

func seedManga(db *sql.DB) error {
	fmt.Println("Seeding manga...")

	// Sample manga data
	mangaList := []models.Manga{
		{
			ID:            "manga-001",
			Title:         "One Piece",
			Author:        "Eiichiro Oda",
			Genres:        []string{"Action", "Adventure", "Fantasy"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 1100,
			Description:   "The story follows Monkey D. Luffy, a young man whose body gained the properties of rubber after unintentionally eating a Devil Fruit.",
			CoverImageURL: "https://example.com/onepiece.jpg",
		},
		{
			ID:            "manga-002",
			Title:         "Naruto",
			Author:        "Masashi Kishimoto",
			Genres:        []string{"Action", "Adventure", "Martial Arts"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 700,
			Description:   "The story follows Naruto Uzumaki, a young ninja who seeks recognition from his peers and dreams of becoming the Hokage.",
			CoverImageURL: "https://example.com/naruto.jpg",
		},
		{
			ID:            "manga-003",
			Title:         "Attack on Titan",
			Author:        "Hajime Isayama",
			Genres:        []string{"Action", "Drama", "Fantasy", "Horror"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 139,
			Description:   "The story is set in a world where humanity lives inside cities surrounded by enormous walls protecting from Titans.",
			CoverImageURL: "https://example.com/aot.jpg",
		},
		{
			ID:            "manga-004",
			Title:         "My Hero Academia",
			Author:        "Kohei Horikoshi",
			Genres:        []string{"Action", "Superhero", "School"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 400,
			Description:   "A world where people with superpowers are the norm and our main protagonist was born without them.",
			CoverImageURL: "https://example.com/mha.jpg",
		},
		{
			ID:            "manga-005",
			Title:         "Demon Slayer",
			Author:        "Koyoharu Gotouge",
			Genres:        []string{"Action", "Supernatural", "Historical"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 205,
			Description:   "A family is attacked by demons and only two members survive - Tanjiro and his sister Nezuko.",
			CoverImageURL: "https://example.com/demonslayer.jpg",
		},
		{
			ID:            "manga-006",
			Title:         "Death Note",
			Author:        "Tsugumi Ohba",
			Genres:        []string{"Mystery", "Psychological", "Supernatural"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 108,
			Description:   "A high school student discovers a supernatural notebook that allows him to kill anyone by writing their name.",
			CoverImageURL: "https://example.com/deathnote.jpg",
		},
		{
			ID:            "manga-007",
			Title:         "Fullmetal Alchemist",
			Author:        "Hiromu Arakawa",
			Genres:        []string{"Action", "Adventure", "Fantasy", "Steampunk"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 116,
			Description:   "Two brothers search for the Philosopher's Stone to restore their bodies after a failed alchemical ritual.",
			CoverImageURL: "https://example.com/fma.jpg",
		},
		{
			ID:            "manga-008",
			Title:         "Tokyo Ghoul",
			Author:        "Sui Ishida",
			Genres:        []string{"Action", "Horror", "Supernatural"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 143,
			Description:   "A college student is turned into a half-ghoul and must navigate both human and ghoul societies.",
			CoverImageURL: "https://example.com/tokyoghoul.jpg",
		},
		{
			ID:            "manga-009",
			Title:         "Jujutsu Kaisen",
			Author:        "Gege Akutami",
			Genres:        []string{"Action", "Supernatural", "School"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 250,
			Description:   "A high school student joins a secret organization of Jujutsu Sorcerers to kill a powerful Curse.",
			CoverImageURL: "https://example.com/jjk.jpg",
		},
		{
			ID:            "manga-010",
			Title:         "Chainsaw Man",
			Author:        "Tatsuki Fujimoto",
			Genres:        []string{"Action", "Horror", "Supernatural"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 97,
			Description:   "A young man merges with his pet devil and becomes a devil hunter working for the government.",
			CoverImageURL: "https://example.com/chainsawman.jpg",
		},
	}

	for _, manga := range mangaList {
		// Marshal genres to JSON
		genresJSON, err := manga.MarshalGenres()
		if err != nil {
			return fmt.Errorf("failed to marshal genres for %s: %w", manga.Title, err)
		}

		// Insert manga
		_, err = db.Exec(`
			INSERT OR IGNORE INTO manga (id, title, author, genres, status, total_chapters, description, cover_image_url)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, manga.ID, manga.Title, manga.Author, genresJSON, manga.Status, manga.TotalChapters, manga.Description, manga.CoverImageURL)

		if err != nil {
			return fmt.Errorf("failed to insert manga %s: %w", manga.Title, err)
		}

		fmt.Printf("  Created manga: %s by %s\n", manga.Title, manga.Author)
	}

	// Seed some user progress
	fmt.Println("Seeding user progress...")
	userProgress := []struct {
		userID         string
		mangaID        string
		currentChapter int
		status         models.ReadingStatus
	}{
		{"user-testuser", "manga-001", 50, models.ReadingStatusReading},
		{"user-testuser", "manga-002", 700, models.ReadingStatusCompleted},
		{"user-alice", "manga-003", 100, models.ReadingStatusReading},
		{"user-alice", "manga-006", 0, models.ReadingStatusPlanToRead},
		{"user-bob", "manga-004", 200, models.ReadingStatusReading},
	}

	for _, up := range userProgress {
		_, err := db.Exec(`
			INSERT OR IGNORE INTO user_progress (user_id, manga_id, current_chapter, status)
			VALUES (?, ?, ?, ?)
		`, up.userID, up.mangaID, up.currentChapter, up.status)

		if err != nil {
			return fmt.Errorf("failed to insert user progress: %w", err)
		}
	}

	fmt.Println("  Created sample user progress entries")

	return nil
}
