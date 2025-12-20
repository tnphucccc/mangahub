package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/tnphucccc/mangahub/internal/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	client pb.MangaServiceClient
	ctx    context.Context
)

func main() {
	// Command line flags
	host := flag.String("host", "localhost", "gRPC server host")
	port := flag.String("port", "9092", "gRPC server port")
	flag.Parse()

	// Connect to gRPC server
	addr := fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("Connecting to gRPC server at %s...", addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client = pb.NewMangaServiceClient(conn)
	ctx = context.Background()

	log.Println("✅ Connected to gRPC server")

	// Interactive command loop
	fmt.Println("\n=== gRPC Client Connected ===")
	fmt.Println("Commands:")
	fmt.Println("  get <manga_id>                              - Get manga by ID")
	fmt.Println("  search title=<query>                        - Search by title")
	fmt.Println("  search author=<query>                       - Search by author")
	fmt.Println("  search status=<status>                      - Search by status")
	fmt.Println("  search title=<query> limit=<n> offset=<n>   - Search with pagination")
	fmt.Println("  update <user_id> <manga_id> <chapter>       - Update progress")
	fmt.Println("  update <user_id> <manga_id> <chapter> <rating> <status> - Full update")
	fmt.Println("  help                                         - Show commands")
	fmt.Println("  quit                                         - Exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse and execute command
		handleCommand(line)
	}
}

func handleCommand(line string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return
	}

	cmd := parts[0]

	switch cmd {
	case "quit", "exit":
		log.Println("Exiting...")
		os.Exit(0)

	case "help":
		showHelp()

	case "get":
		if len(parts) < 2 {
			fmt.Println("Usage: get <manga_id>")
			return
		}
		handleGetManga(parts[1])

	case "search":
		if len(parts) < 2 {
			fmt.Println("Usage: search title=<query> | author=<query> | status=<status>")
			return
		}
		handleSearch(parts[1:])

	case "update":
		if len(parts) < 4 {
			fmt.Println("Usage: update <user_id> <manga_id> <chapter> [rating] [status]")
			return
		}
		handleUpdateProgress(parts[1:])

	default:
		fmt.Println("Unknown command. Type 'help' for available commands.")
	}
}

func handleGetManga(mangaID string) {
	req := &pb.GetMangaRequest{
		MangaId: mangaID,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.GetManga(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("❌ Error: %s (Code: %s)\n", st.Message(), st.Code())
		return
	}

	fmt.Println("✅ Manga retrieved:")
	fmt.Printf("   ID: %s\n", resp.Id)
	fmt.Printf("   Title: %s\n", resp.Title)
	fmt.Printf("   Author: %s\n", resp.Author)
	fmt.Printf("   Genres: %v\n", resp.Genres)
	fmt.Printf("   Status: %s\n", resp.Status)
	fmt.Printf("   Total Chapters: %d\n", resp.TotalChapters)
	fmt.Printf("   Description: %s\n", truncateString(resp.Description, 100))
	if resp.CoverUrl != "" {
		fmt.Printf("   Cover URL: %s\n", resp.CoverUrl)
	}
}

func handleSearch(args []string) {
	req := &pb.SearchRequest{
		Limit:  10, // Default limit
		Offset: 0,  // Default offset
	}

	// Parse arguments (key=value format)
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid argument format: %s (expected key=value)\n", arg)
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "title":
			req.Title = value
		case "author":
			req.Author = value
		case "genre":
			req.Genre = value
		case "status":
			req.Status = value
		case "limit":
			if n, err := strconv.ParseInt(value, 10, 32); err == nil {
				req.Limit = int32(n)
			}
		case "offset":
			if n, err := strconv.ParseInt(value, 10, 32); err == nil {
				req.Offset = int32(n)
			}
		default:
			fmt.Printf("Unknown search parameter: %s\n", key)
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.SearchManga(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("❌ Error: %s (Code: %s)\n", st.Message(), st.Code())
		return
	}

	fmt.Printf("✅ Found %d manga:\n", len(resp.Manga))
	for i, manga := range resp.Manga {
		fmt.Printf("   %d. [ID:%s] %s by %s (%s) - %d chapters\n",
			i+1, manga.Id, manga.Title, manga.Author, manga.Status, manga.TotalChapters)
	}

	if len(resp.Manga) == 0 {
		fmt.Println("   No results found")
	}
}

func handleUpdateProgress(args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: update <user_id> <manga_id> <chapter> [rating] [status]")
		return
	}

	userID := args[0]
	mangaID := args[1]

	chapter, err := strconv.ParseInt(args[2], 10, 32)
	if err != nil {
		fmt.Printf("Invalid chapter number: %s\n", args[2])
		return
	}

	req := &pb.UpdateProgressRequest{
		UserId:  userID,
		MangaId: mangaID,
		Chapter: int32(chapter),
	}

	// Optional rating
	if len(args) >= 4 {
		if rating, err := strconv.ParseInt(args[3], 10, 32); err == nil {
			req.Rating = int32(rating)
		}
	}

	// Optional status
	if len(args) >= 5 {
		req.Status = args[4]
	} else {
		req.Status = "reading" // Default status
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.UpdateProgress(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("❌ Error: %s (Code: %s)\n", st.Message(), st.Code())
		return
	}

	fmt.Println("✅ Progress updated:")
	fmt.Printf("   User ID: %s\n", resp.Progress.UserId)
	fmt.Printf("   Manga ID: %s\n", resp.Progress.MangaId)
	fmt.Printf("   Current Chapter: %d\n", resp.Progress.CurrentChapter)
	fmt.Printf("   Status: %s\n", resp.Progress.Status)
	if resp.Progress.Rating > 0 {
		fmt.Printf("   Rating: %d/10\n", resp.Progress.Rating)
	}
	fmt.Printf("   Updated At: %s\n", resp.Progress.UpdatedAt)
}

func showHelp() {
	fmt.Println("\nAvailable commands:")
	fmt.Println()
	fmt.Println("Manga Retrieval:")
	fmt.Println("  get <manga_id>")
	fmt.Println("    Example: get manga-123")
	fmt.Println()
	fmt.Println("Search:")
	fmt.Println("  search title=<query>")
	fmt.Println("    Example: search title=One")
	fmt.Println("  search author=<query>")
	fmt.Println("    Example: search author=Oda")
	fmt.Println("  search status=<status>")
	fmt.Println("    Example: search status=ongoing")
	fmt.Println("  search title=<query> limit=<n> offset=<n>")
	fmt.Println("    Example: search title=One limit=5 offset=0")
	fmt.Println()
	fmt.Println("Progress Update:")
	fmt.Println("  update <user_id> <manga_id> <chapter>")
	fmt.Println("    Example: update user-123 manga-456 50")
	fmt.Println("  update <user_id> <manga_id> <chapter> <rating> <status>")
	fmt.Println("    Example: update user-123 manga-456 50 8 reading")
	fmt.Println()
	fmt.Println("Status values: reading, completed, plan_to_read, on_hold, dropped")
	fmt.Println("Rating: 1-10")
	fmt.Println()
	fmt.Println("Other:")
	fmt.Println("  help - Show this help message")
	fmt.Println("  quit - Exit the client")
	fmt.Println()
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
