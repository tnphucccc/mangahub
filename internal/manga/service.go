package manga

import (
	"fmt"
	"time"

	"github.com/tnphucccc/mangahub/internal/user"
	"github.com/tnphucccc/mangahub/pkg/models"
)

// Service handles manga business logic
type Service struct {
	repo                *Repository
	userRepo            *user.Repository
	TCPBroadcastChan    chan models.TCPProgressBroadcast
	UDPNotificationChan chan models.UDPNotification
}

// NewService creates a new manga service
func NewService(repo *Repository, userRepo *user.Repository) *Service {
	return &Service{
		repo:                repo,
		userRepo:            userRepo,
		TCPBroadcastChan:    make(chan models.TCPProgressBroadcast, 100),
		UDPNotificationChan: make(chan models.UDPNotification, 100),
	}
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
	manga, err := s.repo.FindByID(req.MangaID)
	if err != nil {
		return fmt.Errorf("manga not found")
	}

	// Check if already in library
	if _, err := s.repo.GetProgress(userID, req.MangaID); err == nil {
		return fmt.Errorf("manga already in library")
	}

	// Validate chapter number
	// If TotalChapters is 0, it means "n/a" (unknown), so we allow any chapter number.
	if req.CurrentChapter < 0 {
		return fmt.Errorf("invalid chapter number: cannot be negative")
	}
	if manga.TotalChapters > 0 && req.CurrentChapter > manga.TotalChapters {
		return fmt.Errorf("invalid chapter number: exceeds total chapters (%d)", manga.TotalChapters)
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

	// Get existing progress (needed for partial updates and validation)
	existingProgress, err := s.repo.GetProgress(userID, mangaID)
	if err != nil {
		return fmt.Errorf("manga not in user's library")
	}

	// Determine new chapter
	currentChapter := existingProgress.CurrentChapter
	if req.CurrentChapter != nil {
		currentChapter = *req.CurrentChapter

		// Validate chapter number
		if manga.TotalChapters <= 0 {
			if currentChapter > 0 {
				return fmt.Errorf("invalid chapter number: cannot update chapter progress when total chapters is N/A")
			}
		} else {
			if currentChapter < 0 {
				return fmt.Errorf("invalid chapter number: cannot be negative")
			}
			if currentChapter > manga.TotalChapters {
				return fmt.Errorf("invalid chapter number: exceeds total chapters (%d)", manga.TotalChapters)
			}
		}
	}

	// Update progress
	if err := s.repo.UpdateProgress(userID, mangaID, currentChapter, req.Status, req.Rating); err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	// Fetch username for broadcast
	username := "User"
	if s.userRepo != nil {
		if u, err := s.userRepo.FindByID(userID); err == nil {
			username = u.Username
		}
	}

	// Trigger TCP broadcast for real-time sync
	s.TCPBroadcastChan <- models.TCPProgressBroadcast{
		UserID:         userID,
		Username:       username,
		MangaID:        mangaID,
		MangaTitle:     manga.Title,
		CurrentChapter: currentChapter,
		Status:         derefStatus(req.Status),
		Timestamp:      time.Now(),
	}

	return nil
}

func derefStatus(s *models.ReadingStatus) models.ReadingStatus {
	if s == nil {
		return ""
	}
	return *s
}

// NotifyNotification sends a UDP notification (triggered by admin)
func (s *Service) NotifyNotification(notification models.UDPNotification) {
	s.UDPNotificationChan <- notification
}

// GetProgress retrieves user's progress for a specific manga
func (s *Service) GetProgress(userID, mangaID string) (*models.UserProgress, error) {
	progress, err := s.repo.GetProgress(userID, mangaID)
	if err != nil {
		return nil, fmt.Errorf("progress not found")
	}

	return progress, nil
}
