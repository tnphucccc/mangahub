package manga

import (
	"fmt"

	"github.com/tnphucccc/mangahub/pkg/models"
)

// Service handles manga business logic
type Service struct {
	repo *Repository
}

// NewService creates a new manga service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetByID retrieves a manga by ID
func (s *Service) GetByID(id string) (*models.Manga, error) {
	manga, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("manga not found")
	}
	return manga, nil
}

// Search searches for manga
func (s *Service) Search(query models.MangaSearchQuery) ([]models.Manga, int, error) {
	mangaList, total, err := s.repo.Search(query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search manga: %w", err)
	}

	return mangaList, total, nil
}

// GetAll retrieves all manga
func (s *Service) GetAll(limit, offset int) ([]models.Manga, int, error) {
	mangaList, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get manga: %w", err)
	}

	return mangaList, total, nil
}

// GetUserLibrary retrieves a user's manga library
func (s *Service) GetUserLibrary(userID string) ([]models.UserProgressWithManga, error) {
	library, err := s.repo.GetUserLibrary(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user library: %w", err)
	}

	return library, nil
}

// AddToLibrary adds a manga to user's library
func (s *Service) AddToLibrary(userID string, req models.LibraryAddRequest) error {
	// Verify manga exists
	_, err := s.repo.FindByID(req.MangaID)
	if err != nil {
		return fmt.Errorf("manga not found")
	}

	// Add to library
	if err := s.repo.AddToLibrary(userID, req.MangaID, req.Status, req.CurrentChapter); err != nil {
		return fmt.Errorf("failed to add manga to library: %w", err)
	}

	return nil
}

// UpdateProgress updates user's reading progress
func (s *Service) UpdateProgress(userID, mangaID string, req models.ProgressUpdateRequest) error {
	// Verify manga exists
	manga, err := s.repo.FindByID(mangaID)
	if err != nil {
		return fmt.Errorf("manga not found")
	}

	// Validate chapter number
	if req.CurrentChapter < 0 || req.CurrentChapter > manga.TotalChapters {
		return fmt.Errorf("invalid chapter number")
	}

	// Update progress
	if err := s.repo.UpdateProgress(userID, mangaID, req.CurrentChapter, req.Status, req.Rating); err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	return nil
}

// GetProgress retrieves user's progress for a specific manga
func (s *Service) GetProgress(userID, mangaID string) (*models.UserProgress, error) {
	progress, err := s.repo.GetProgress(userID, mangaID)
	if err != nil {
		return nil, fmt.Errorf("progress not found")
	}

	return progress, nil
}
