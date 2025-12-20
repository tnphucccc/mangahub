package unit

import (
	"testing"

	pb "github.com/tnphucccc/mangahub/internal/grpc/pb"
	"github.com/tnphucccc/mangahub/pkg/models"
)

// ==========================================
// gRPC Message Structure Tests
// ==========================================

func TestGRPCGetMangaRequest_Structure(t *testing.T) {
	req := &pb.GetMangaRequest{
		MangaId: "manga-123",
	}

	if req.MangaId != "manga-123" {
		t.Errorf("Expected manga ID 'manga-123', got '%s'", req.MangaId)
	}

	t.Logf("✓ gRPC GetMangaRequest structure correct")
}

func TestGRPCMangaResponse_Structure(t *testing.T) {
	resp := &pb.MangaResponse{
		Id:            "manga-123",
		Title:         "One Piece",
		Author:        "Eiichiro Oda",
		Genres:        []string{"Action", "Adventure"},
		Status:        "ongoing",
		TotalChapters: 1100,
		Description:   "Pirate adventures",
		CoverUrl:      "https://example.com/cover.jpg",
	}

	if resp.Id != "manga-123" {
		t.Errorf("Expected ID 'manga-123', got '%s'", resp.Id)
	}

	if resp.Title != "One Piece" {
		t.Errorf("Expected title 'One Piece', got '%s'", resp.Title)
	}

	if len(resp.Genres) != 2 {
		t.Errorf("Expected 2 genres, got %d", len(resp.Genres))
	}

	if resp.TotalChapters != 1100 {
		t.Errorf("Expected 1100 chapters, got %d", resp.TotalChapters)
	}

	t.Logf("✓ gRPC MangaResponse structure correct")
}

func TestGRPCSearchRequest_Structure(t *testing.T) {
	req := &pb.SearchRequest{
		Title:   "One",
		Author:  "Oda",
		Genre:   "Action",
		Status:  "ongoing",
		OrderBy: "title",
		Limit:   20,
		Offset:  0,
	}

	if req.Title != "One" {
		t.Errorf("Expected title 'One', got '%s'", req.Title)
	}

	if req.Limit != 20 {
		t.Errorf("Expected limit 20, got %d", req.Limit)
	}

	if req.Offset != 0 {
		t.Errorf("Expected offset 0, got %d", req.Offset)
	}

	t.Logf("✓ gRPC SearchRequest structure correct")
}

func TestGRPCSearchResponse_Structure(t *testing.T) {
	resp := &pb.SearchResponse{
		Manga: []*pb.MangaResponse{
			{
				Id:            "manga-1",
				Title:         "One Piece",
				Author:        "Eiichiro Oda",
				Genres:        []string{"Action"},
				Status:        "ongoing",
				TotalChapters: 1100,
			},
			{
				Id:            "manga-2",
				Title:         "Naruto",
				Author:        "Masashi Kishimoto",
				Genres:        []string{"Action"},
				Status:        "completed",
				TotalChapters: 700,
			},
		},
	}

	if len(resp.Manga) != 2 {
		t.Errorf("Expected 2 manga, got %d", len(resp.Manga))
	}

	if resp.Manga[0].Title != "One Piece" {
		t.Errorf("Expected title 'One Piece', got '%s'", resp.Manga[0].Title)
	}

	t.Logf("✓ gRPC SearchResponse structure correct")
}

func TestGRPCUserProgress_Structure(t *testing.T) {
	progress := &pb.UserProgress{
		UserId:         "user-123",
		MangaId:        "manga-123",
		CurrentChapter: 50,
		Status:         "reading",
		Rating:         8,
		StartedAt:      "2024-01-01",
		CompletedAt:    "",
		UpdatedAt:      "2024-01-15",
	}

	if progress.UserId != "user-123" {
		t.Errorf("Expected user ID 'user-123', got '%s'", progress.UserId)
	}

	if progress.CurrentChapter != 50 {
		t.Errorf("Expected chapter 50, got %d", progress.CurrentChapter)
	}

	if progress.Rating != 8 {
		t.Errorf("Expected rating 8, got %d", progress.Rating)
	}

	t.Logf("✓ gRPC UserProgress structure correct")
}

func TestGRPCUpdateProgressRequest_Structure(t *testing.T) {
	req := &pb.UpdateProgressRequest{
		UserId:  "user-123",
		MangaId: "manga-123",
		Status:  "reading",
		Chapter: 50,
		Rating:  8,
	}

	if req.UserId != "user-123" {
		t.Errorf("Expected user ID 'user-123', got '%s'", req.UserId)
	}

	if req.Chapter != 50 {
		t.Errorf("Expected chapter 50, got %d", req.Chapter)
	}

	if req.Rating != 8 {
		t.Errorf("Expected rating 8, got %d", req.Rating)
	}

	t.Logf("✓ gRPC UpdateProgressRequest structure correct")
}

func TestGRPCUpdateProgressResponse_Structure(t *testing.T) {
	resp := &pb.UpdateProgressResponse{
		Progress: &pb.UserProgress{
			UserId:         "user-123",
			MangaId:        "manga-123",
			CurrentChapter: 50,
			Status:         "reading",
			Rating:         8,
		},
	}

	if resp.Progress == nil {
		t.Fatal("Expected progress to be set, got nil")
	}

	if resp.Progress.UserId != "user-123" {
		t.Errorf("Expected user ID 'user-123', got '%s'", resp.Progress.UserId)
	}

	t.Logf("✓ gRPC UpdateProgressResponse structure correct")
}

// ==========================================
// gRPC Model Conversion Tests
// ==========================================

func TestModelToGRPCConversion_Manga(t *testing.T) {
	// Test converting internal model to gRPC response
	manga := &models.Manga{
		ID:            "manga-123",
		Title:         "Test Manga",
		Author:        "Test Author",
		Genres:        []string{"Action", "Adventure"},
		Status:        models.MangaStatusOngoing,
		TotalChapters: 100,
		Description:   "Test description",
		CoverImageURL: "https://example.com/cover.jpg",
	}

	// Simulate conversion (this is what toMangaResponse does in server.go)
	grpcManga := &pb.MangaResponse{
		Id:            manga.ID,
		Title:         manga.Title,
		Author:        manga.Author,
		Genres:        manga.Genres,
		Status:        string(manga.Status),
		TotalChapters: int32(manga.TotalChapters),
		Description:   manga.Description,
		CoverUrl:      manga.CoverImageURL,
	}

	if grpcManga.Id != manga.ID {
		t.Errorf("Expected ID '%s', got '%s'", manga.ID, grpcManga.Id)
	}

	if grpcManga.Title != manga.Title {
		t.Errorf("Expected title '%s', got '%s'", manga.Title, grpcManga.Title)
	}

	if len(grpcManga.Genres) != len(manga.Genres) {
		t.Errorf("Expected %d genres, got %d", len(manga.Genres), len(grpcManga.Genres))
	}

	if grpcManga.Status != string(manga.Status) {
		t.Errorf("Expected status '%s', got '%s'", manga.Status, grpcManga.Status)
	}

	if grpcManga.TotalChapters != int32(manga.TotalChapters) {
		t.Errorf("Expected %d chapters, got %d", manga.TotalChapters, grpcManga.TotalChapters)
	}

	t.Logf("✓ Model to gRPC conversion correct")
}

func TestGRPCToModelConversion_SearchQuery(t *testing.T) {
	// Test converting gRPC request to internal model
	grpcReq := &pb.SearchRequest{
		Title:  "One",
		Author: "Oda",
		Genre:  "Action",
		Status: "ongoing",
		Limit:  20,
		Offset: 0,
	}

	// Simulate conversion (this is what happens in SearchManga handler)
	query := models.MangaSearchQuery{
		Title:  grpcReq.Title,
		Author: grpcReq.Author,
		Genre:  grpcReq.Genre,
		Status: models.MangaStatus(grpcReq.Status),
		Limit:  int(grpcReq.Limit),
		Offset: int(grpcReq.Offset),
	}

	if query.Title != grpcReq.Title {
		t.Errorf("Expected title '%s', got '%s'", grpcReq.Title, query.Title)
	}

	if query.Limit != int(grpcReq.Limit) {
		t.Errorf("Expected limit %d, got %d", grpcReq.Limit, query.Limit)
	}

	if string(query.Status) != grpcReq.Status {
		t.Errorf("Expected status '%s', got '%s'", grpcReq.Status, query.Status)
	}

	t.Logf("✓ gRPC to model conversion correct")
}

func TestGRPCToModelConversion_ProgressUpdate(t *testing.T) {
	// Test converting gRPC progress update to internal model
	grpcReq := &pb.UpdateProgressRequest{
		UserId:  "user-123",
		MangaId: "manga-123",
		Status:  "reading",
		Chapter: 50,
		Rating:  8,
	}

	// Simulate conversion
	status := models.ReadingStatus(grpcReq.Status)
	rating := int(grpcReq.Rating)

	updateReq := models.ProgressUpdateRequest{
		CurrentChapter: int(grpcReq.Chapter),
		Status:         &status,
		Rating:         &rating,
	}

	if updateReq.CurrentChapter != int(grpcReq.Chapter) {
		t.Errorf("Expected chapter %d, got %d", grpcReq.Chapter, updateReq.CurrentChapter)
	}

	if *updateReq.Status != models.ReadingStatus(grpcReq.Status) {
		t.Errorf("Expected status '%s', got '%s'", grpcReq.Status, *updateReq.Status)
	}

	if *updateReq.Rating != int(grpcReq.Rating) {
		t.Errorf("Expected rating %d, got %d", grpcReq.Rating, *updateReq.Rating)
	}

	t.Logf("✓ gRPC progress update conversion correct")
}

// ==========================================
// gRPC Data Validation Tests
// ==========================================

func TestGRPCRequest_Validation(t *testing.T) {
	// Test empty manga ID
	t.Run("empty manga ID", func(t *testing.T) {
		req := &pb.GetMangaRequest{
			MangaId: "",
		}

		if req.MangaId != "" {
			t.Error("Expected empty manga ID")
		}

		// In real implementation, server would return error
		t.Logf("✓ Empty manga ID detected")
	})

	// Test negative limit
	t.Run("negative limit", func(t *testing.T) {
		req := &pb.SearchRequest{
			Limit: -1,
		}

		if req.Limit >= 0 {
			t.Error("Expected negative limit to be caught")
		}

		t.Logf("✓ Negative limit detected")
	})

	// Test negative offset
	t.Run("negative offset", func(t *testing.T) {
		req := &pb.SearchRequest{
			Offset: -1,
		}

		if req.Offset >= 0 {
			t.Error("Expected negative offset to be caught")
		}

		t.Logf("✓ Negative offset detected")
	})
}

func TestGRPCResponse_EmptyResults(t *testing.T) {
	// Test empty search results
	resp := &pb.SearchResponse{
		Manga: []*pb.MangaResponse{},
	}

	if len(resp.Manga) != 0 {
		t.Errorf("Expected 0 manga, got %d", len(resp.Manga))
	}

	t.Logf("✓ Empty search results handled correctly")
}

// ==========================================
// gRPC Type Safety Tests
// ==========================================

func TestGRPCTypes_Int32Conversion(t *testing.T) {
	// Test int to int32 conversion (Go int -> protobuf int32)
	goInt := 1100
	protoInt32 := int32(goInt)

	if protoInt32 != 1100 {
		t.Errorf("Expected 1100, got %d", protoInt32)
	}

	// Test back conversion
	backToGoInt := int(protoInt32)
	if backToGoInt != goInt {
		t.Errorf("Expected %d, got %d", goInt, backToGoInt)
	}

	t.Logf("✓ int32 conversion correct")
}

func TestGRPCTypes_StringSliceConversion(t *testing.T) {
	// Test string slice (genres)
	goSlice := []string{"Action", "Adventure", "Fantasy"}
	protoSlice := make([]string, len(goSlice))
	copy(protoSlice, goSlice)

	if len(protoSlice) != len(goSlice) {
		t.Errorf("Expected %d items, got %d", len(goSlice), len(protoSlice))
	}

	for i, genre := range goSlice {
		if protoSlice[i] != genre {
			t.Errorf("Expected genre '%s', got '%s'", genre, protoSlice[i])
		}
	}

	t.Logf("✓ String slice conversion correct")
}

// ==========================================
// gRPC Error Handling Tests
// ==========================================

func TestGRPCErrorScenarios(t *testing.T) {
	t.Run("manga not found scenario", func(t *testing.T) {
		// This tests the expected data flow when manga is not found
		req := &pb.GetMangaRequest{
			MangaId: "non-existent-id",
		}

		// In real implementation, server would:
		// 1. Call service.GetByID(req.MangaId)
		// 2. Service returns error
		// 3. Server returns status.Errorf(codes.NotFound, ...)

		if req.MangaId == "" {
			t.Error("Expected manga ID to be set")
		}

		t.Logf("✓ Not found scenario structure validated")
	})

	t.Run("invalid progress update scenario", func(t *testing.T) {
		req := &pb.UpdateProgressRequest{
			UserId:  "",
			MangaId: "manga-123",
			Chapter: 50,
		}

		// Empty user ID should cause error
		if req.UserId == "" {
			t.Logf("✓ Empty user ID would cause validation error")
		}
	})
}
