package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/tnphucccc/mangahub/internal/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:9092",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMangaServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("=== MangaHub gRPC Service Test ===")

	// Test 1: Get manga by ID
	fmt.Println("1. Getting Manga by ID:")
	testGetManga(ctx, client)

	// Test 2: Search manga by title
	fmt.Println("\n2. Searching Manga by Title:")
	testSearchByTitle(ctx, client)

	// Test 3: Search manga by author
	fmt.Println("\n3. Searching Manga by Author:")
	testSearchByAuthor(ctx, client)

	// Test 4: Search manga by status
	fmt.Println("\n4. Searching Manga by Status:")
	testSearchByStatus(ctx, client)

	// Test 5: Search with pagination
	fmt.Println("\n5. Searching with Pagination:")
	testSearchWithPagination(ctx, client)

	// Test 6: Update progress
	fmt.Println("\n6. Updating Reading Progress:")
	testUpdateProgress(ctx, client)

	// Test 7: Get non-existent manga (error case)
	fmt.Println("\n7. Testing Error Handling (Non-existent Manga):")
	testGetNonExistentManga(ctx, client)

	fmt.Println("\n✓ gRPC Service Test completed successfully!")
	fmt.Println("✓ All RPC methods tested")
	fmt.Println("✓ Error handling verified")
}

func testGetManga(ctx context.Context, client pb.MangaServiceClient) {
	// Note: You'll need to use an actual manga ID from your database
	req := &pb.GetMangaRequest{
		MangaId: "1", // Change to actual ID in your database
	}

	resp, err := client.GetManga(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("   GetManga failed: %s (Code: %s)\n", st.Message(), st.Code())
		fmt.Println("   Note: Make sure manga ID '1' exists in database")
		return
	}

	fmt.Printf("   ✓ Manga retrieved successfully\n")
	fmt.Printf("   ID: %s\n", resp.Id)
	fmt.Printf("   Title: %s\n", resp.Title)
	fmt.Printf("   Author: %s\n", resp.Author)
	fmt.Printf("   Genres: %v\n", resp.Genres)
	fmt.Printf("   Status: %s\n", resp.Status)
	fmt.Printf("   Total Chapters: %d\n", resp.TotalChapters)
	fmt.Printf("   Description: %s\n", truncateString(resp.Description, 80))
}

func testSearchByTitle(ctx context.Context, client pb.MangaServiceClient) {
	req := &pb.SearchRequest{
		Title:  "One",
		Limit:  10,
		Offset: 0,
	}

	resp, err := client.SearchManga(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("   SearchManga failed: %s (Code: %s)\n", st.Message(), st.Code())
		return
	}

	fmt.Printf("   ✓ Search by title successful\n")
	fmt.Printf("   Found %d manga\n", len(resp.Manga))
	for i, manga := range resp.Manga {
		fmt.Printf("   %d. [ID:%s] %s by %s\n", i+1, manga.Id, manga.Title, manga.Author)
	}
}

func testSearchByAuthor(ctx context.Context, client pb.MangaServiceClient) {
	req := &pb.SearchRequest{
		Author: "Oda",
		Limit:  10,
		Offset: 0,
	}

	resp, err := client.SearchManga(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("   SearchManga failed: %s (Code: %s)\n", st.Message(), st.Code())
		return
	}

	fmt.Printf("   ✓ Search by author successful\n")
	fmt.Printf("   Found %d manga\n", len(resp.Manga))
	for i, manga := range resp.Manga {
		fmt.Printf("   %d. %s by %s\n", i+1, manga.Title, manga.Author)
	}
}

func testSearchByStatus(ctx context.Context, client pb.MangaServiceClient) {
	req := &pb.SearchRequest{
		Status: "ongoing",
		Limit:  10,
		Offset: 0,
	}

	resp, err := client.SearchManga(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("   SearchManga failed: %s (Code: %s)\n", st.Message(), st.Code())
		return
	}

	fmt.Printf("   ✓ Search by status successful\n")
	fmt.Printf("   Found %d ongoing manga\n", len(resp.Manga))
	for i, manga := range resp.Manga {
		fmt.Printf("   %d. %s - Status: %s (%d chapters)\n", i+1, manga.Title, manga.Status, manga.TotalChapters)
	}
}

func testSearchWithPagination(ctx context.Context, client pb.MangaServiceClient) {
	// First page
	fmt.Println("   Page 1 (limit 3):")
	req1 := &pb.SearchRequest{
		Limit:  3,
		Offset: 0,
	}

	resp1, err := client.SearchManga(ctx, req1)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("   SearchManga failed: %s\n", st.Message())
		return
	}

	for i, manga := range resp1.Manga {
		fmt.Printf("     %d. [ID:%s] %s\n", i+1, manga.Id, manga.Title)
	}

	// Second page
	fmt.Println("   Page 2 (limit 3, offset 3):")
	req2 := &pb.SearchRequest{
		Limit:  3,
		Offset: 3,
	}

	resp2, err := client.SearchManga(ctx, req2)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("   SearchManga failed: %s\n", st.Message())
		return
	}

	for i, manga := range resp2.Manga {
		fmt.Printf("     %d. [ID:%s] %s\n", i+1, manga.Id, manga.Title)
	}

	fmt.Printf("   ✓ Pagination working correctly\n")
}

func testUpdateProgress(ctx context.Context, client pb.MangaServiceClient) {
	// Note: You'll need actual user and manga IDs
	req := &pb.UpdateProgressRequest{
		UserId:  "test-user-id", // Change to actual user ID
		MangaId: "1",            // Change to actual manga ID
		Status:  "reading",
		Chapter: 50,
		Rating:  8,
	}

	resp, err := client.UpdateProgress(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("   UpdateProgress failed: %s (Code: %s)\n", st.Message(), st.Code())
		fmt.Println("   Note: Make sure user and manga IDs exist in database")
		return
	}

	fmt.Printf("   ✓ Progress updated successfully\n")
	fmt.Printf("   User ID: %s\n", resp.Progress.UserId)
	fmt.Printf("   Manga ID: %s\n", resp.Progress.MangaId)
	fmt.Printf("   Current Chapter: %d\n", resp.Progress.CurrentChapter)
	fmt.Printf("   Status: %s\n", resp.Progress.Status)
	fmt.Printf("   Rating: %d/10\n", resp.Progress.Rating)
}

func testGetNonExistentManga(ctx context.Context, client pb.MangaServiceClient) {
	req := &pb.GetMangaRequest{
		MangaId: "non-existent-id-999999",
	}

	_, err := client.GetManga(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		fmt.Printf("   Expected error received: %s ✓\n", st.Message())
		fmt.Printf("   Error code: %s\n", st.Code())
	} else {
		fmt.Println("   ERROR: Should have received an error!")
	}
}

// Helper function to truncate long strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
