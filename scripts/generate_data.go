package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tnphucccc/mangahub/pkg/models"
)

const (
	DataDir     = "data/manga"
	ManualCount = 100
	ApiCount    = 100
)

var (
	genresList = []string{"Action", "Adventure", "Comedy", "Drama", "Fantasy", "Horror", "Mystery", "Romance", "Sci-Fi", "Slice of Life", "Sports", "Supernatural"}
	statuses   = []models.MangaStatus{models.MangaStatusOngoing, models.MangaStatusCompleted, models.MangaStatusHiatus}

	adjectives = []string{"Silent", "Iron", "Blue", "Dark", "Eternal", "Broken", "Lost", "Infinite", "Crimson", "Last", "Hidden", "Divine", "Shattered", "Crystal", "Burning", "Frozen", "Golden", "Silver"}
	nouns      = []string{"Soul", "Alchemist", "Hunter", "Warrior", "Exorcist", "Titan", "Ghoul", "Note", "Piece", "Slayer", "Hero", "Academia", "Clover", "Tail", "Gate", "Crown", "Sword", "Shield"}
	authors    = []string{"Akira Toriyama", "Eiichiro Oda", "Masashi Kishimoto", "Hirohiko Araki", "Naoko Takeuchi", "Rumiko Takahashi", "Osamu Tezuka", "Kentaro Miura", "Hiromu Arakawa", "Yoshihiro Togashi", "Junji Ito"}
)

func main() {
	// Clear existing data
	os.RemoveAll(DataDir)

	if err := os.MkdirAll(DataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	fmt.Println("Generating 200 manga entries (100 manual + 100 simulated API)...")

	// 1. Generate "Manual" entries
	generateManualEntries()

	// 2. Generate "Simulated API" entries
	generateSimulatedAPIEntries()

	fmt.Printf("✓ Successfully generated %d manga JSON files in %s\n", ManualCount+ApiCount, DataDir)
}

func generateManualEntries() {
	// Real data provided by user
	realManga := []models.Manga{
		{
			Title:         "My Robot Has Been Acting Strange Lately",
			Author:        "GoHome_kun",
			Genres:        []string{"Sci-Fi", "Romance", "Web Comic"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "A story about a robot acting strange.",
		},
		{
			Title:         "She's Likely Aiming for My Older Brother",
			Author:        "Unknown",
			Genres:        []string{"Romance", "Comedy", "School Life", "Slice of Life"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "A high school guy who doesn't trust girls anymore has been tricked too many times by girls who just wanted to get close to his good-looking older brother.",
		},
		{
			Title:         "Versatile Mage",
			Author:        "Chaos",
			Genres:        []string{"Action", "Comedy", "Magic", "Harem", "Isekai", "Drama", "School Life", "Fantasy", "Supernatural"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 1181,
			Description:   "Our hero, Mo Fan, inherits a magical necklace—the next day, he wakes up to find that the world has changed.",
		},
		{
			Title:         "Pet",
			Author:        "Ranjou Miyake",
			Genres:        []string{"Psychological", "Drama", "Supernatural", "Mystery", "Tragedy"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 55,
			Description:   "The story revolves around people who possess the ability to infiltrate human minds and manipulate memories.",
		},
		{
			Title:         "Sakamoto Days (Official Colored)",
			Author:        "Yuto Suzuki",
			Genres:        []string{"Sci-Fi", "Action", "Comedy", "Martial Arts", "Mafia", "Delinquents", "Slice of Life", "Supernatural"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "Taro Sakamoto was the ultimate assassin, feared by villains and admired by hitmen. But one day…he fell in love!",
		},
		{
			Title:         "Cinderella's Pocket-Sized Protector",
			Author:        "Unknown",
			Genres:        []string{"Romance", "Comedy", "Magic", "Drama", "Fantasy"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "After marrying the prince, Cinderella’s “happily ever after” ends when he falls for another woman and proposes to her.",
		},
		{
			Title:         "I'm Sick and Tired of My Childhood Friend's Abuse",
			Author:        "Unknown",
			Genres:        []string{"Romance", "Drama", "School Life"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "My childhood friend is also my girlfriend. A dream come true, right? Hell no. She insults me on a daily basis.",
		},
		{
			Title:         "My Life With Amelia",
			Author:        "KlyptoKicks",
			Genres:        []string{"Romance", "Comedy", "School Life", "Web Comic", "Slice of Life"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "The story of a scarred girl named Amelia and how she slowly opens up.",
		},
		{
			Title:         "Anecdotes with Master After the Novel Transmigration",
			Author:        "Unknown",
			Genres:        []string{"Historical", "Romance", "Boys' Love", "Fantasy", "Web Comic"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "After transmigration, both Shizun and the original protagonist's personas went out of character!!",
		},
		{
			Title:         "WITCHRIV",
			Author:        "Unknown",
			Genres:        []string{"Action", "Survival", "Magical Girls", "Magic", "Gore", "Drama", "Fantasy", "Tragedy"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "A girl named Nona lives in plain sight among humans, hiding her true identity as a mage.",
		},
		{
			Title:         "Baby Steps",
			Author:        "Hikaru Katsuki",
			Genres:        []string{"Romance", "Comedy", "Sports", "Drama", "School Life", "Slice of Life"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 455,
			Description:   "Maruo Eiichirou (Ei-Chan), a first year honor student, one day decides he's unhappy with the way things are and lacks exercise.",
		},
		{
			Title:         "This Kind of Thing is Fine",
			Author:        "Kouji Oishi",
			Genres:        []string{"Comedy", "Drama", "Slice of Life"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "Do what we want, when we want, with the person we want. Such a relationship needs no label.",
		},
		{
			Title:         "Zilbagias the Demon Prince",
			Author:        "Unknown",
			Genres:        []string{"Reincarnation", "Action", "Demons", "Comedy", "Adventure", "Magic", "Harem", "Drama", "Fantasy"},
			Status:        models.MangaStatusCancelled,
			TotalChapters: 20,
			Description:   "The hero Alexander and his comrades unleash a daring raid on the Demon King's castle.",
		},
		{
			Title:         "Kowloon Generic Romance",
			Author:        "Jun Mayuzuki",
			Genres:        []string{"Sci-Fi", "Psychological", "Romance", "Drama", "Slice of Life", "Mystery"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "The greatest labyrinth of the 20th century, a drama for working men and women in a town called Kowloon Walled City.",
		},
		{
			Title:         "Imaginater!",
			Author:        "Unknown",
			Genres:        []string{"Reincarnation", "Monsters", "Adventure", "Magic", "Isekai", "Fantasy"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "Hinata, a Heisei gal, died from an incurable illness. When she was reincarnated, it was into a fantasy world where imagination literally becomes magic.",
		},
		{
			Title:         "The One-Eyed, One-Armed, One-Legged Sorcerer",
			Author:        "Unknown",
			Genres:        []string{"Action", "Demons", "Adventure", "Magic", "Fantasy", "Supernatural", "Mystery"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "The Autonomous City of Ains Territory prospers through highly advanced sorcery and magic.",
		},
		{
			Title:         "Nijigasaki of the Rebellion",
			Author:        "Choboraunyopomi",
			Genres:        []string{"Sci-Fi", "Action", "Comedy", "Adventure", "Web Comic", "Music"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "Hangyaku no Nijigasaki (Nijigasaki of the Rebellion) is an official Love Live! manga.",
		},
		{
			Title:         "The Reptiles Wish To Be Praised!",
			Author:        "Unknown",
			Genres:        []string{"Animals", "Monster Girls", "Slice of Life"},
			Status:        models.MangaStatusCompleted,
			TotalChapters: 1,
			Description:   "Two maids working hard(?) at their residence. They are reptiles that have taken on human form.",
		},
		{
			Title:         "Yakuza Reincarnation",
			Author:        "Takeshi Natsuhara",
			Genres:        []string{"Reincarnation", "Monsters", "Action", "Demons", "Comedy", "Martial Arts", "Mafia", "Adventure", "Magic", "Isekai", "Fantasy", "Supernatural"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "An old school yakuza, Nagamasa Ryumatsu, ended up losing his life due to certain circumstances.",
		},
		{
			Title:         "Destiny Unchain Online",
			Author:        "Unknown",
			Genres:        []string{"Genderswap", "Action", "Adventure", "Virtual Reality", "Video Games", "Magic", "Fantasy", "Vampires"},
			Status:        models.MangaStatusOngoing,
			TotalChapters: 0,
			Description:   "Mitsuki Kou is such a feared PVP gamer that he's earned the nickname \"Crim the Headhunter\".",
		},
	}

	for i := 0; i < ManualCount; i++ {
		var m models.Manga
		if i < len(realManga) {
			m = realManga[i]
			m.ID = slugify(m.Title) // Force slug ID
		} else {
			m = generateRandomManga(i, "manual")
			// generateRandomManga already slugifies ID somewhat, but let's ensure it's clean
		}

		m.CoverImageURL = fmt.Sprintf("https://via.placeholder.com/300x450?text=%s", urlEncode(m.Title))
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()

		saveManga(m)
	}
}

func slugify(s string) string {
	s = strings.ToLower(s)
	// Replace non-alphanumeric with hyphen
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}
	// Trim hyphens and duplicate hyphens
	res := result.String()
	for strings.Contains(res, "--") {
		res = strings.ReplaceAll(res, "--", "-")
	}
	return strings.Trim(res, "-")
}

func generateSimulatedAPIEntries() {
	for i := 0; i < ApiCount; i++ {
		m := generateRandomManga(i+ManualCount, "mangadex")
		saveManga(m)
	}
}

func generateRandomManga(index int, prefix string) models.Manga {
	rand.Seed(time.Now().UnixNano() + int64(index))

	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	title := fmt.Sprintf("%s %s", adj, noun)

	id := fmt.Sprintf("%s-%d", prefix, index)

	numGenres := rand.Intn(3) + 1
	myGenres := make([]string, 0)
	for k := 0; k < numGenres; k++ {
		myGenres = append(myGenres, genresList[rand.Intn(len(genresList))])
	}

	return models.Manga{
		ID:            id,
		Title:         title,
		Author:        authors[rand.Intn(len(authors))],
		Genres:        myGenres,
		Status:        statuses[rand.Intn(len(statuses))],
		TotalChapters: rand.Intn(200) + 10,
		Description:   fmt.Sprintf("An exciting %s manga about %s %s.", myGenres[0], strings.ToLower(adj), strings.ToLower(noun)),
		CoverImageURL: fmt.Sprintf("https://via.placeholder.com/300x450?text=%s", urlEncode(title)),
	}
}

func saveManga(m models.Manga) {
	filename := filepath.Join(DataDir, fmt.Sprintf("%s.json", m.ID))
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create file %s: %v", filename, err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(m); err != nil {
		log.Printf("Failed to encode manga %s: %v", m.Title, err)
	}
}

func urlEncode(s string) string {
	return strings.ReplaceAll(s, " ", "+")
}
