package main

import (
	"fmt"
	"os"
)

const version = "1.0.0-dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version":
		fmt.Printf("MangaHub CLI v%s\n", version)
	case "help":
		printUsage()
	case "init":
		initConfig()
	case "server":
		handleServerCommand()
	case "auth":
		handleAuthCommand()
	case "manga":
		handleMangaCommand()
	case "library":
		handleLibraryCommand()
	case "progress":
		handleProgressCommand()
	case "chat":
		handleChatCommand()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("MangaHub CLI - Manga Tracking System")
	fmt.Println("\nUsage:")
	fmt.Println("  mangahub <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Println("  version              Show version information")
	fmt.Println("  help                 Show this help message")
	fmt.Println("  init                 Initialize configuration")
	fmt.Println("  server               Manage servers (start, stop, status)")
	fmt.Println("  auth                 Authentication (register, login, logout)")
	fmt.Println("  manga                Manga operations (search, info, list)")
	fmt.Println("  library              Library management (add, remove, list)")
	fmt.Println("  progress             Progress tracking (update, history)")
	fmt.Println("  chat                 Chat system (join, send)")
	fmt.Println("\nFor more information on a command:")
	fmt.Println("  mangahub <command> help")
}

func initConfig() {
	fmt.Println("Initializing MangaHub configuration...")
	// TODO: Create config file at ~/.mangahub/config.yaml
	fmt.Println("Configuration initialized successfully!")
}

func handleServerCommand() {
	// TODO: Implement server management
	fmt.Println("Server command - TODO")
}

func handleAuthCommand() {
	// TODO: Implement authentication
	fmt.Println("Auth command - TODO")
}

func handleMangaCommand() {
	// TODO: Implement manga operations
	fmt.Println("Manga command - TODO")
}

func handleLibraryCommand() {
	// TODO: Implement library operations
	fmt.Println("Library command - TODO")
}

func handleProgressCommand() {
	// TODO: Implement progress tracking
	fmt.Println("Progress command - TODO")
}

func handleChatCommand() {
	// TODO: Implement chat
	fmt.Println("Chat command - TODO")
}
