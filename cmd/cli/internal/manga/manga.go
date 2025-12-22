package manga

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/tnphucccc/mangahub/cmd/cli/internal/config"
	climodels "github.com/tnphucccc/mangahub/pkg/cli/models"
)

func HandleMangaCommand() {
	if len(os.Args) < 3 {
		printMangaUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "search":
		mangaSearch()
	case "get":
		mangaGet()
	case "all":
		mangaGetAll()
	default:
		fmt.Printf("Unknown manga subcommand: %s\n", subcommand)
		printMangaUsage()
		os.Exit(1)
	}
}

func printMangaUsage() {
	fmt.Println("Usage: mangahub manga <subcommand> [options]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  search               Search for manga by title, author, genre, or status")
	fmt.Println("  get <id>             Get details for a specific manga by ID")
	fmt.Println("  all                  Get all manga with pagination")
}

func mangaSearch() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	queryParams := url.Values{}
	for i := 3; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg[2:], "=", 2)
			if len(parts) == 2 {
				queryParams.Add(parts[0], parts[1])
			} else {
				fmt.Printf("Invalid argument format: %s\n", arg)
				os.Exit(1)
			}
		}
	}

	apiURL := fmt.Sprintf("http://%s:%d/api/v1/manga?%s", cliConfig.Server.Host, cliConfig.Server.HTTPPort, queryParams.Encode())
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Error connecting to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var apiResp struct {
			Success bool                        `json:"success"`
			Data    climodels.MangaListResponse `json:"data"`
			Meta    climodels.Meta              `json:"meta"`
			Error   climodels.ErrorResponse     `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding API response: %v\n", err)
			os.Exit(1)
		}

		if !apiResp.Success {
			fmt.Printf("❌ Manga search failed: %s\n", apiResp.Error.Error)
			os.Exit(1)
		}

		fmt.Println("Manga Search Results:")
		for _, m := range apiResp.Data.Items {
			fmt.Printf("  ID: %s\n", m.ID)
			fmt.Printf("  Title: %s\n", m.Title)
			fmt.Printf("  Author: %s\n", m.Author)
			fmt.Printf("  Genres: %s\n", strings.Join(m.Genres, ", "))
			fmt.Printf("  Status: %s\n", m.Status)
			fmt.Printf("  Total Chapters: %d\n", m.TotalChapters)
			fmt.Printf("  --------------------\n")
		}
	} else {
		var apiResp struct {
			Success bool                    `json:"success"`
			Data    interface{}             `json:"data"`
			Error   climodels.ErrorResponse `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding API error response: %v\n", err)
			os.Exit(1)
		}
		errMsg := apiResp.Error.Error
		if errMsg == "" {
			errMsg = fmt.Sprintf("API returned status %d", resp.StatusCode)
		}
		fmt.Printf("❌ Manga search failed: %s\n", errMsg)
		os.Exit(1)
	}
}

func mangaGet() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: mangahub manga get <id>")
		os.Exit(1)
	}
	mangaID := os.Args[3]

	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	apiURL := fmt.Sprintf("http://%s:%d/api/v1/manga/%s", cliConfig.Server.Host, cliConfig.Server.HTTPPort, mangaID)
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Error connecting to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var apiResp struct {
			Success bool                          `json:"success"`
			Data    climodels.MangaDetailResponse `json:"data"`
			Error   climodels.ErrorResponse       `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding API response: %v\n", err)
			os.Exit(1)
		}

		if !apiResp.Success {
			fmt.Printf("❌ Failed to get manga: %s\n", apiResp.Error.Error)
			os.Exit(1)
		}

		m := apiResp.Data.Manga
		fmt.Printf("Manga Details (ID: %s):\n", m.ID)
		fmt.Printf("  Title: %s\n", m.Title)
		fmt.Printf("  Author: %s\n", m.Author)
		fmt.Printf("  Genres: %s\n", strings.Join(m.Genres, ", "))
		fmt.Printf("  Status: %s\n", m.Status)
		fmt.Printf("  Total Chapters: %d\n", m.TotalChapters)
		fmt.Printf("  Description: %s\n", m.Description)
		fmt.Printf("  Cover Image URL: %s\n", m.CoverImageURL)
		fmt.Printf("  Created At: %s\n", m.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Updated At: %s\n", m.UpdatedAt.Format("2006-01-02 15:04:05"))
	} else {
		var apiResp struct {
			Success bool                    `json:"success"`
			Data    interface{}             `json:"data"`
			Error   climodels.ErrorResponse `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding API error response: %v\n", err)
			os.Exit(1)
		}
		errMsg := apiResp.Error.Error
		if errMsg == "" {
			errMsg = fmt.Sprintf("API returned status %d", resp.StatusCode)
		}
		fmt.Printf("❌ Failed to get manga: %s\n", errMsg)
		os.Exit(1)
	}
}

func mangaGetAll() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	queryParams := url.Values{}
	for i := 3; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg[2:], "=", 2)
			if len(parts) == 2 {
				queryParams.Add(parts[0], parts[1])
			} else {
				fmt.Printf("Invalid argument format: %s\n", arg)
				os.Exit(1)
			}
		}
	}

	apiURL := fmt.Sprintf("http://%s:%d/api/v1/manga/all?%s", cliConfig.Server.Host, cliConfig.Server.HTTPPort, queryParams.Encode())
	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Error connecting to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var apiResp struct {
			Success bool                        `json:"success"`
			Data    climodels.MangaListResponse `json:"data"`
			Meta    climodels.Meta              `json:"meta"`
			Error   climodels.ErrorResponse     `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding API response: %v\n", err)
			os.Exit(1)
		}

		if !apiResp.Success {
			fmt.Printf("❌ Failed to get all manga: %s\n", apiResp.Error.Error)
			os.Exit(1)
		}

		fmt.Println("All Manga:")
		for _, m := range apiResp.Data.Items {
			fmt.Printf("  ID: %s\n", m.ID)
			fmt.Printf("  Title: %s\n", m.Title)
			fmt.Printf("  Author: %s\n", m.Author)
			fmt.Printf("  Genres: %s\n", strings.Join(m.Genres, ", "))
			fmt.Printf("  Status: %s\n", m.Status)
			fmt.Printf("  Total Chapters: %d\n", m.TotalChapters)
			fmt.Printf("  --------------------\n")
		}
	} else {
		var apiResp struct {
			Success bool                    `json:"success"`
			Data    interface{}             `json:"data"`
			Error   climodels.ErrorResponse `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding API error response: %v\n", err)
			os.Exit(1)
		}
		errMsg := apiResp.Error.Error
		if errMsg == "" {
			errMsg = fmt.Sprintf("API returned status %d", resp.StatusCode)
		}
		fmt.Printf("❌ Failed to get all manga: %s\n", errMsg)
		os.Exit(1)
	}
}
