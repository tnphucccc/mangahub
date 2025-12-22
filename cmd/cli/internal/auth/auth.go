package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/tnphucccc/mangahub/cmd/cli/internal/config"
	climodels "github.com/tnphucccc/mangahub/pkg/cli/models"
	"golang.org/x/term"
)

func HandleAuthCommand() {
	if len(os.Args) < 3 {
		printAuthUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "register":
		authRegister()
	case "login":
		authLogin()
	case "logout":
		authLogout()
	case "status":
		authStatus()
	default:
		fmt.Printf("Unknown auth subcommand: %s\n", subcommand)
		printAuthUsage()
		os.Exit(1)
	}
}

func printAuthUsage() {
	fmt.Println("Usage: mangahub auth <subcommand> [options]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  register             Register a new user")
	fmt.Println("  login                Login an existing user")
	fmt.Println("  logout               Logout current user")
	fmt.Println("  status               Show current authentication status")
}

func authRegister() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	var username, email, password, confirmPassword string

	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	fmt.Print("Enter email: ")
	fmt.Scanln(&email)

	fmt.Print("Enter password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("\nError reading password: %v\n", err)
		os.Exit(1)
	}
	password = string(bytePassword)
	fmt.Println()

	fmt.Print("Confirm password: ")
	byteConfirmPassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("\nError confirming password: %v\n", err)
		os.Exit(1)
	}
	confirmPassword = string(byteConfirmPassword)
	fmt.Println()

	if password != confirmPassword {
		fmt.Println("Error: Passwords do not match.")
		os.Exit(1)
	}

	reqBody := climodels.UserRegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	jsonReqBody, _ := json.Marshal(reqBody)

	apiURL := fmt.Sprintf("http://%s:%d/api/v1/auth/register", cliConfig.Server.Host, cliConfig.Server.HTTPPort)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonReqBody))
	if err != nil {
		fmt.Printf("Error connecting to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		var apiResp struct {
			Success bool                    `json:"success"`
			Data    climodels.AuthResponse  `json:"data"`
			Error   climodels.ErrorResponse `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding API response: %v\n", err)
			os.Exit(1)
		}

		if !apiResp.Success {
			fmt.Printf("❌ Registration failed: %s\n", apiResp.Error.Message)
			os.Exit(1)
		}

		cliConfig.User.Username = apiResp.Data.User.Username
		cliConfig.User.Token = apiResp.Data.Token
		if err := config.SaveCLIConfig(cliConfig); err != nil {
			fmt.Printf("Warning: Could not save config: %v\n", err)
		}
		fmt.Printf("✅ User '%s' registered and logged in successfully!\n", apiResp.Data.User.Username)
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
		errMsg := apiResp.Error.Message
		if errMsg == "" {
			errMsg = fmt.Sprintf("API returned status %d", resp.StatusCode)
		}
		fmt.Printf("❌ Registration failed: %s\n", errMsg)
		os.Exit(1)
	}
}

func authLogin() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	var username, password string

	fmt.Print("Enter username: ")
	fmt.Scanln(&username)

	fmt.Print("Enter password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("\nError reading password: %v\n", err)
		os.Exit(1)
	}
	password = string(bytePassword)
	fmt.Println()

	reqBody := climodels.UserLoginRequest{
		Username: username,
		Password: password,
	}
	jsonReqBody, _ := json.Marshal(reqBody)

	apiURL := fmt.Sprintf("http://%s:%d/api/v1/auth/login", cliConfig.Server.Host, cliConfig.Server.HTTPPort)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonReqBody))
	if err != nil {
		fmt.Printf("Error connecting to API: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var apiResp struct {
			Success bool                    `json:"success"`
			Data    climodels.AuthResponse  `json:"data"`
			Error   climodels.ErrorResponse `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding API response: %v\n", err)
			os.Exit(1)
		}

		if !apiResp.Success {
			fmt.Printf("❌ Login failed: %s\n", apiResp.Error.Message)
			os.Exit(1)
		}

		cliConfig.User.Username = apiResp.Data.User.Username
		cliConfig.User.Token = apiResp.Data.Token
		if err := config.SaveCLIConfig(cliConfig); err != nil {
			fmt.Printf("Warning: Could not save config: %v\n", err)
		}
		fmt.Printf("✅ Logged in successfully as '%s'!\n", apiResp.Data.User.Username)
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
		errMsg := apiResp.Error.Message
		if errMsg == "" {
			errMsg = fmt.Sprintf("API returned status %d", resp.StatusCode)
		}
		fmt.Printf("❌ Login failed: %s\n", errMsg)
		os.Exit(1)
	}
}

func authLogout() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cliConfig.User.Token == "" {
		fmt.Println("You are not currently logged in.")
		return
	}

	cliConfig.User.Username = ""
	cliConfig.User.Token = ""
	if err := config.SaveCLIConfig(cliConfig); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Logged out successfully.")
}

func authStatus() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("--- Authentication Status ---")
	if cliConfig.User.Token != "" {
		fmt.Printf("Status: Logged In\n")
		fmt.Printf("Username: %s\n", cliConfig.User.Username)
		fmt.Println("Token: (present)")
	} else {
		fmt.Println("Status: Logged Out")
	}
	fmt.Println("---------------------------")
}
