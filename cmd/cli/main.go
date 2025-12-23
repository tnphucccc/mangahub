package main

import (
	"fmt"
	"os"

	"github.com/tnphucccc/mangahub/cmd/cli/internal/auth"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/chat"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/config"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/grpc_client"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/library"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/manga"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/notify"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/progress"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/server"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/stats"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/sync"
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
		config.InitConfig()
	case "server":
		server.HandleServerCommand()
	case "auth":
		auth.HandleAuthCommand()
	case "manga":
		manga.HandleMangaCommand()
	case "library":
		library.HandleLibraryCommand()
	case "progress":
		progress.HandleProgressCommand()
	case "chat":
		chat.HandleChatCommand()
	case "stats":
		stats.HandleStatsCommand()
	case "sync":
		sync.HandleSyncCommand()
	case "notify":
		notify.HandleNotifyCommand()
	case "grpc":
		grpc_client.HandleGRPCCommand()
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
	fmt.Println("  auth                 Authentication (register, login, logout, status)")
	fmt.Println("  manga                Manga operations (search, info, list)")
	fmt.Println("  library              Library management (add, remove, list)")
	fmt.Println("  progress             Progress tracking (update, history)")
	fmt.Println("  chat                 Chat system (join, send)")
	fmt.Println("  stats                User statistics (view)")
	fmt.Println("  sync                 Real-time synchronization (TCP monitor)")
	fmt.Println("  notify               Notifications (UDP listener)")
	fmt.Println("  grpc                 Internal service queries (gRPC)")
	fmt.Println("\nFor more information on a command:")
	fmt.Println("  mangahub <command> help")
}
