package library

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/tnphucccc/mangahub/cmd/cli/internal/config"
	climodels "github.com/tnphucccc/mangahub/pkg/cli/models"
)

func HandleLibraryCommand() {
	if len(os.Args) < 3 {
		printLibraryUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "add":
		libraryAdd()
	case "list":
		libraryList()
	default:
		fmt.Printf("Unknown library subcommand: %s\n", subcommand)
		printLibraryUsage()
		os.Exit(1)
	}
}

func printLibraryUsage() {
	fmt.Println("Usage: mangahub library <subcommand> [options]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  add <manga_id> --status=<status> [--chapter=<chapter>]")
	fmt.Println("  list                                                     List user's library")
}

func libraryAdd() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: mangahub library add <manga_id> --status=<status> [--chapter=<chapter>]")
		os.Exit(1)
	}

	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cliConfig.User.Token == "" {
		fmt.Println("Error: Not logged in. Please use 'mangahub auth login' first.")
		os.Exit(1)
	}

	mangaID := os.Args[3]
	status := ""
	chapter := 0

	for i := 4; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "--status=") {
			status = strings.TrimPrefix(arg, "--status=")
		} else if strings.HasPrefix(arg, "--chapter=") {
			chapter, err = strconv.Atoi(strings.TrimPrefix(arg, "--chapter="))
			if err != nil {
				fmt.Printf("Error: Invalid chapter number: %v\n", err)
				os.Exit(1)
			}
		}
	}

	if status == "" {
		fmt.Println("Error: --status is required (e.g., reading, completed, plan_to_read)")
		os.Exit(1)
	}

	reqBody := climodels.LibraryAddRequest{
		MangaID:        mangaID,
		Status:         status,
		CurrentChapter: chapter,
	}
	jsonReqBody, _ := json.Marshal(reqBody)

	apiURL := fmt.Sprintf("http://%s:%d/api/v1/users/library", cliConfig.Server.Host, cliConfig.Server.HTTPPort)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cliConfig.User.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error connecting to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("✅ Manga '%s' added to library successfully!\n", mangaID)
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
		fmt.Printf("❌ Failed to add manga to library: %s\n", errMsg)
		os.Exit(1)
	}
}

func libraryList() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cliConfig.User.Token == "" {
		fmt.Println("Error: Not logged in. Please use 'mangahub auth login' first.")
		os.Exit(1)
	}

	apiURL := fmt.Sprintf("http://%s:%d/api/v1/users/library", cliConfig.Server.Host, cliConfig.Server.HTTPPort)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer "+cliConfig.User.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error connecting to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var apiResp struct {
			Success bool                          `json:"success"`
			Data    climodels.LibraryListResponse `json:"data"`
			Meta    climodels.Meta                `json:"meta"`
			Error   climodels.ErrorResponse       `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding API response: %v\n", err)
			os.Exit(1)
		}

		if !apiResp.Success {
			fmt.Printf("❌ Failed to list library: %s\n", apiResp.Error.Error)
			os.Exit(1)
		}

		if len(apiResp.Data.Items) == 0 {
			fmt.Println("Your library is empty.")
			return
		}

		fmt.Println("Your Manga Library:")
		for _, item := range apiResp.Data.Items {
			fmt.Printf("  Title: %s\n", item.Manga.Title)
			fmt.Printf("  Status: %s\n", item.UserProgress.Status)
			fmt.Printf("  Current Chapter: %d\n", item.UserProgress.CurrentChapter)
			if item.UserProgress.Rating.Valid && item.UserProgress.Rating.Int64 > 0 {
				fmt.Printf("  Rating: %d/10\n", item.UserProgress.Rating.Int64)
			}
			fmt.Printf("  Updated At: %s\n", item.UserProgress.UpdatedAt)
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
		fmt.Printf("❌ Failed to list library: %s\n", errMsg)
		os.Exit(1)
	}
}
