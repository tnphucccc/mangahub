package stats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/tnphucccc/mangahub/cmd/cli/internal/config"
	climodels "github.com/tnphucccc/mangahub/pkg/cli/models"
)

func HandleStatsCommand() {
	if len(os.Args) < 3 {
		printStatsUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "view":
		viewStats()
	default:
		fmt.Printf("Unknown stats subcommand: %s\n", subcommand)
		printStatsUsage()
		os.Exit(1)
	}
}

func printStatsUsage() {
	fmt.Println("Usage: mangahub stats <subcommand>")
	fmt.Println("\nSubcommands:")
	fmt.Println("  view                 View your reading statistics")
}

func viewStats() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cliConfig.User.Token == "" {
		fmt.Println("Error: Not logged in. Please use 'mangahub auth login' first.")
		os.Exit(1)
	}

	apiURL := fmt.Sprintf("http://%s:%d/api/v1/users/stats", cliConfig.Server.Host, cliConfig.Server.HTTPPort)
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
			Success bool                    `json:"success"`
			Data    map[string]interface{}  `json:"data"`
			Error   climodels.ErrorResponse `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding API response: %v\n", err)
			os.Exit(1)
		}

		if !apiResp.Success {

			fmt.Printf("❌ Failed to get statistics: %s\n", apiResp.Error.Message)

			os.Exit(1)

		}

		fmt.Println("--- Your Reading Statistics ---")

		fmt.Printf("Total Manga in Library:  %v\n", apiResp.Data["total_manga"])

		fmt.Printf("Total Chapters Read:     %v\n", apiResp.Data["total_chapters_read"])

		fmt.Printf("Manga Reading:           %v\n", apiResp.Data["reading_manga"])

		fmt.Printf("Manga Completed:         %v\n", apiResp.Data["completed_manga"])

		fmt.Printf("Manga Plan to Read:      %v\n", apiResp.Data["plan_to_read_manga"])

		fmt.Printf("Average Rating:          %.2f/10\n", apiResp.Data["average_rating"])

		fmt.Println("-------------------------------")

	} else {

		var apiResp struct {
			Success bool `json:"success"`

			Error climodels.ErrorResponse `json:"error"`
		}

		json.NewDecoder(resp.Body).Decode(&apiResp)

		fmt.Printf("❌ Failed to get statistics: %s\n", apiResp.Error.Message)

		os.Exit(1)

	}
}
