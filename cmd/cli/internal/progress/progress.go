package progress

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

func HandleProgressCommand() {
	if len(os.Args) < 3 {
		printProgressUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "update":
		progressUpdate()
	default:
		fmt.Printf("Unknown progress subcommand: %s\n", subcommand)
		printProgressUsage()
		os.Exit(1)
	}
}

func printProgressUsage() {
	fmt.Println("Usage: mangahub progress <subcommand> [options]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  update <manga_id> [--chapter=<chapter>] [--status=<status>] [--rating=<rating>]")
}

func progressUpdate() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: mangahub progress update <manga_id> [--chapter=<chapter>] [--status=<status>] [--rating=<rating>]")
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
	reqBody := climodels.ProgressUpdateRequest{}
	var rating *int

	for i := 4; i < len(os.Args); i++ {
		arg := os.Args[i]
		if strings.HasPrefix(arg, "--chapter=") {
			chapter, err := strconv.Atoi(strings.TrimPrefix(arg, "--chapter="))
			if err != nil {
				fmt.Printf("Error: Invalid chapter number: %v\n", err)
				os.Exit(1)
			}
			reqBody.CurrentChapter = chapter
		} else if strings.HasPrefix(arg, "--status=") {
			reqBody.Status = strings.TrimPrefix(arg, "--status=")
		} else if strings.HasPrefix(arg, "--rating=") {
			r, err := strconv.Atoi(strings.TrimPrefix(arg, "--rating="))
			if err != nil {
				fmt.Printf("Error: Invalid rating: %v\n", err)
				os.Exit(1)
			}
			rating = &r
		}
	}
	reqBody.Rating = rating

	jsonReqBody, _ := json.Marshal(reqBody)

	apiURL := fmt.Sprintf("http://%s:%d/api/v1/users/progress/%s", cliConfig.Server.Host, cliConfig.Server.HTTPPort, mangaID)
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonReqBody))
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

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Progress for manga '%s' updated successfully!\n", mangaID)
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
		fmt.Printf("❌ Failed to update progress: %s\n", errMsg)
		os.Exit(1)
	}
}
