package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tnphucccc/mangahub/pkg/models"
)

const (
	DataFile    = "data/manga.json"
	ManualCount = 100
	ApiCount    = 100
	MangaDexURL = "https://api.mangadex.org/manga?limit=100&includes[]=author&includes[]=cover_art&contentRating[]=safe&contentRating[]=suggestive"
)

var (
	allManga   []models.Manga
	genresList = []string{"Action", "Adventure", "Comedy", "Drama", "Fantasy", "Horror", "Mystery", "Romance", "Sci-Fi", "Slice of Life", "Sports", "Supernatural"}
	statuses   = []models.MangaStatus{models.MangaStatusOngoing, models.MangaStatusCompleted, models.MangaStatusHiatus}

	adjectives = []string{"Silent", "Iron", "Blue", "Dark", "Eternal", "Broken", "Lost", "Infinite", "Crimson", "Last", "Hidden", "Divine", "Shattered", "Crystal", "Burning", "Frozen", "Golden", "Silver"}
	nouns      = []string{"Soul", "Alchemist", "Hunter", "Warrior", "Exorcist", "Titan", "Ghoul", "Note", "Piece", "Slayer", "Hero", "Academia", "Clover", "Tail", "Gate", "Crown", "Sword", "Shield"}
	authors    = []string{"Akira Toriyama", "Eiichiro Oda", "Masashi Kishimoto", "Hirohiko Araki", "Naoko Takeuchi", "Rumiko Takahashi", "Osamu Tezuka", "Kentaro Miura", "Hiromu Arakawa", "Yoshihiro Togashi", "Junji Ito"}
)

// MangaDex API types
type MDResponse struct {
	Data []MDManga `json:"data"`
}

type MDManga struct {
	ID         string `json:"id"`
	Attributes struct {
		Title       map[string]string `json:"title"`
		Description map[string]string `json:"description"`
		Status      string            `json:"status"`
		LastChapter string            `json:"lastChapter"`
		Tags        []struct {
			Attributes struct {
				Name  map[string]string `json:"name"`
				Group string            `json:"group"`
			} `json:"attributes"`
		} `json:"tags"`
	} `json:"attributes"`
	Relationships []struct {
		ID         string                 `json:"id"`
		Type       string                 `json:"type"`
		Attributes map[string]interface{} `json:"attributes"`
	} `json:"relationships"`
}

func main() {
	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	fmt.Println("Generating manga entries...")

	// 1. Web Scraping Practice (Requirement)
	scrapeQuotes()
	checkHttpBin()

	// 2. Generate "Manual" entries
	generateManualEntries()

	// 3. Fetch real data from MangaDex
	fmt.Println("Fetching data from MangaDex API...")
	fetchMangaDexEntries()

	// 4. Save to single file
	saveAllManga()

	fmt.Printf("✓ Successfully populated %s (%d entries)\n", DataFile, len(allManga))
}

var scrapedQuotes []string

func scrapeQuotes() {
	fmt.Println("Scraping educational quotes from quotes.toscrape.com...")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("http://quotes.toscrape.com")
	if err != nil {
		log.Printf("Warning: Failed to scrape quotes: %v", err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := os.ReadFile("scripts/generate_data/main.go") // Dummy read to avoid unused
	_ = bodyBytes

	// Simple string parsing for educational purposes (no external scraping lib required for basic task)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	// Extract a few quotes using simple string splitting
	parts := strings.Split(body, "<span class=\"text\" itemprop=\"text\">“")
	for i := 1; i < len(parts) && i <= 5; i++ {
		quote := strings.Split(parts[i], "”</span>")[0]
		scrapedQuotes = append(scrapedQuotes, quote)
	}
	fmt.Printf("✓ Scraped %d quotes for educational practice\n", len(scrapedQuotes))
}

func checkHttpBin() {
	fmt.Println("Verifying scraper connection with httpbin.org...")
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Printf("Warning: httpbin check failed: %v", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("✓ httpbin.org check successful")
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
			m.ID = slugify(m.Title)
		} else {
			m = generateRandomManga(i, "manual")
		}

		// Append a scraped quote to the description for educational practice demonstration
		if len(scrapedQuotes) > 0 {
			quote := scrapedQuotes[rand.Intn(len(scrapedQuotes))]
			m.Description = fmt.Sprintf("%s\n\nNote: %s", m.Description, quote)
		}

		m.CoverImageURL = fmt.Sprintf("https://via.placeholder.com/300x450?text=%s", urlEncode(m.Title))
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()

		saveManga(m)
	}
}

func fetchMangaDexEntries() {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(MangaDexURL)
	if err != nil {
		log.Printf("Error fetching from MangaDex: %v", err)
		return
	}
	defer resp.Body.Close()

	var mdResp MDResponse
	if err := json.NewDecoder(resp.Body).Decode(&mdResp); err != nil {
		log.Printf("Error decoding MangaDex response: %v", err)
		return
	}

	for _, md := range mdResp.Data {
		m := models.Manga{
			ID:          "mangadex-" + md.ID,
			Title:       getLocalized(md.Attributes.Title, "en"),
			Description: getLocalized(md.Attributes.Description, "en"),
			Status:      mapStatus(md.Attributes.Status),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Parse last chapter
		if md.Attributes.LastChapter != "" {
			if ch, err := strconv.ParseFloat(md.Attributes.LastChapter, 64); err == nil {
				m.TotalChapters = int(ch)
			}
		}

		// Extract genres from tags
		for _, tag := range md.Attributes.Tags {
			if tag.Attributes.Group == "genre" {
				m.Genres = append(m.Genres, getLocalized(tag.Attributes.Name, "en"))
			}
		}

		// Find author and cover from relationships
		var coverFile string
		for _, rel := range md.Relationships {
			if rel.Type == "author" && rel.Attributes != nil {
				if name, ok := rel.Attributes["name"].(string); ok {
					m.Author = name
				}
			}
			if rel.Type == "cover_art" && rel.Attributes != nil {
				if fileName, ok := rel.Attributes["fileName"].(string); ok {
					coverFile = fileName
				}
			}
		}

		if m.Author == "" {
			m.Author = "Unknown"
		}

		if coverFile != "" {
			m.CoverImageURL = fmt.Sprintf("https://uploads.mangadex.org/covers/%s/%s", md.ID, coverFile)
		} else {
			m.CoverImageURL = fmt.Sprintf("https://via.placeholder.com/300x450?text=%s", urlEncode(m.Title))
		}

		saveManga(m)
	}
}

func getLocalized(m map[string]string, lang string) string {
	if val, ok := m[lang]; ok {
		return val
	}
	// Fallback to first available
	for _, val := range m {
		return val
	}
	return ""
}

func mapStatus(mdStatus string) models.MangaStatus {
	switch mdStatus {
	case "ongoing":
		return models.MangaStatusOngoing
	case "completed":
		return models.MangaStatusCompleted
	case "hiatus":
		return models.MangaStatusHiatus
	case "cancelled":
		return models.MangaStatusCancelled
	default:
		return models.MangaStatusOngoing
	}
}

func slugify(s string) string {
	s = strings.ToLower(s)
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		} else if r == ' ' || r == '-' || r == '_' {
			result.WriteRune('-')
		}
	}
	res := result.String()
	for strings.Contains(res, "--") {
		res = strings.ReplaceAll(res, "--", "-")
	}
	return strings.Trim(res, "-")
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

func saveAllManga() {
	file, err := os.Create(DataFile)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", DataFile, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(allManga); err != nil {
		log.Fatalf("Failed to encode manga: %v", err)
	}
}

func saveManga(m models.Manga) {
	allManga = append(allManga, m)
}

func urlEncode(s string) string {
	return strings.ReplaceAll(s, " ", "+")
}
