package grpc_client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tnphucccc/mangahub/cmd/cli/internal/config"
	"github.com/tnphucccc/mangahub/internal/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func HandleGRPCCommand() {
	if len(os.Args) < 3 {
		printGRPCUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "get":
		grpcGet()
	default:
		fmt.Printf("Unknown grpc subcommand: %s\n", subcommand)
		printGRPCUsage()
		os.Exit(1)
	}
}

func printGRPCUsage() {
	fmt.Println("Usage: mangahub grpc <subcommand> [options]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  get <id>             Get manga details via gRPC")
}

func grpcGet() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: mangahub grpc get <id>")
		os.Exit(1)
	}
	mangaID := os.Args[3]

	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", cliConfig.Server.Host, cliConfig.Server.GRPCPort)
	fmt.Printf("Connecting to gRPC Server at %s...\n", addr)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewMangaServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetManga(ctx, &pb.GetMangaRequest{MangaId: mangaID})
	if err != nil {
		fmt.Printf("❌ gRPC Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Received Manga Details via gRPC:\n")
	fmt.Printf("  ID: %s\n", resp.Id)
	fmt.Printf("  Title: %s\n", resp.Title)
	fmt.Printf("  Author: %s\n", resp.Author)
	fmt.Printf("  Status: %s\n", resp.Status)
	fmt.Printf("  Total Chapters: %d\n", resp.TotalChapters)
}
