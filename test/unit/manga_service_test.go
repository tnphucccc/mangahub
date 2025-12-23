package unit

import (
	"testing"

	"github.com/tnphucccc/mangahub/pkg/models"
)

// Test manga model and search query structures
func TestMangaSearchQuery_Defaults(t *testing.T) {
	query := models.MangaSearchQuery{}

	if query.Limit != 0 {
		t.Errorf("Expected default limit 0, got %d", query.Limit)
	}

	if query.Offset != 0 {
		t.Errorf("Expected default offset 0, got %d", query.Offset)
	}

	t.Logf("✓ MangaSearchQuery defaults correct")
}

func TestMangaSearchQuery_WithParams(t *testing.T) {
	query := models.MangaSearchQuery{
		Title:  "One Piece",
		Author: "Oda",
		Genre:  "Action",
		Status: models.MangaStatusOngoing,
		Limit:  20,
		Offset: 0,
	}

	if query.Title != "One Piece" {
		t.Errorf("Expected title 'One Piece', got '%s'", query.Title)
	}

	if query.Limit != 20 {
		t.Errorf("Expected limit 20, got %d", query.Limit)
	}

	t.Logf("✓ MangaSearchQuery parameters set correctly")
}

func TestMangaModel(t *testing.T) {
	manga := models.Manga{
		ID:            "test-1",
		Title:         "Test Manga",
		Author:        "Test Author",
		Genres:        []string{"Action", "Adventure"},
		Status:        models.MangaStatusOngoing,
		TotalChapters: 100,
		Description:   "Test description",
	}

	if manga.ID != "test-1" {
		t.Errorf("Expected ID 'test-1', got '%s'", manga.ID)
	}

	if len(manga.Genres) != 2 {
		t.Errorf("Expected 2 genres, got %d", len(manga.Genres))
	}

	if manga.Status != models.MangaStatusOngoing {
		t.Errorf("Expected status 'ongoing', got '%s'", manga.Status)
	}

	t.Logf("✓ Manga model structure correct")
}

func TestMangaStatus_Constants(t *testing.T) {
	statuses := []models.MangaStatus{
		models.MangaStatusOngoing,
		models.MangaStatusCompleted,
		models.MangaStatusHiatus,
		models.MangaStatusCancelled,
	}

	expectedStatuses := []string{"ongoing", "completed", "hiatus", "cancelled"}

	for i, status := range statuses {
		if string(status) != expectedStatuses[i] {
			t.Errorf("Expected status '%s', got '%s'", expectedStatuses[i], status)
		}
	}

	t.Logf("✓ Manga status constants correct")
}

func TestProgressUpdateRequest(t *testing.T) {
	status := models.ReadingStatusReading
	rating := 8
	chapter := 50

	req := models.ProgressUpdateRequest{
		CurrentChapter: &chapter,
		Status:         &status,
		Rating:         &rating,
	}

	if *req.CurrentChapter != 50 {
		t.Errorf("Expected chapter 50, got %d", *req.CurrentChapter)
	}

	if *req.Status != models.ReadingStatusReading {
		t.Errorf("Expected status 'reading', got '%s'", *req.Status)
	}

	if *req.Rating != 8 {
		t.Errorf("Expected rating 8, got %d", *req.Rating)
	}

	t.Logf("✓ ProgressUpdateRequest structure correct")
}

func TestReadingStatus_Constants(t *testing.T) {
	statuses := []models.ReadingStatus{
		models.ReadingStatusReading,
		models.ReadingStatusCompleted,
		models.ReadingStatusPlanToRead,
		models.ReadingStatusOnHold,
		models.ReadingStatusDropped,
	}

	expectedStatuses := []string{"reading", "completed", "plan_to_read", "on_hold", "dropped"}

	for i, status := range statuses {
		if string(status) != expectedStatuses[i] {
			t.Errorf("Expected status '%s', got '%s'", expectedStatuses[i], status)
		}
	}

	t.Logf("✓ Reading status constants correct")
}

func TestLibraryAddRequest(t *testing.T) {
	req := models.LibraryAddRequest{
		MangaID:        "manga-123",
		Status:         models.ReadingStatusPlanToRead,
		CurrentChapter: 0,
	}

	if req.MangaID != "manga-123" {
		t.Errorf("Expected manga ID 'manga-123', got '%s'", req.MangaID)
	}

	if req.Status != models.ReadingStatusPlanToRead {
		t.Errorf("Expected status 'plan_to_read', got '%s'", req.Status)
	}

	t.Logf("✓ LibraryAddRequest structure correct")
}

func TestUserProgress_SafeGetters(t *testing.T) {
	progress := models.UserProgress{
		UserID:         "user-1",
		MangaID:        "manga-1",
		CurrentChapter: 50,
		Status:         models.ReadingStatusReading,
	}

	// Test safe rating getter (when rating is not set, Valid is false)
	rating := progress.GetRatingValue()
	if rating != 0 {
		t.Errorf("Expected rating 0 for invalid, got %d", rating)
	}

	// Set rating using *int
	val := 8
	progress.Rating = &val
	rating = progress.GetRatingValue()
	if rating != 8 {
		t.Errorf("Expected rating 8, got %d", rating)
	}

	t.Logf("✓ UserProgress safe getters working")
}
