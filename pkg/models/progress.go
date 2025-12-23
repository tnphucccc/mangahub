package models

import (
	"database/sql"
	"time"
)

// ReadingStatus represents the user's reading status for a manga
type ReadingStatus string

const (
	ReadingStatusReading    ReadingStatus = "reading"
	ReadingStatusCompleted  ReadingStatus = "completed"
	ReadingStatusPlanToRead ReadingStatus = "plan_to_read"
	ReadingStatusOnHold     ReadingStatus = "on_hold"
	ReadingStatusDropped    ReadingStatus = "dropped"
)

// UserProgress represents a user's reading progress for a manga
type UserProgress struct {
	UserID         string        `json:"user_id" db:"user_id"`
	MangaID        string        `json:"manga_id" db:"manga_id"`
	CurrentChapter int           `json:"current_chapter" db:"current_chapter"`
	Status         ReadingStatus `json:"status" db:"status"`
	Rating         sql.NullInt64 `json:"rating" db:"rating"` // Can be NULL
	StartedAt      sql.NullTime  `json:"started_at" db:"started_at"`
	CompletedAt    sql.NullTime  `json:"completed_at" db:"completed_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
}

// UserProgressWithManga includes manga details with progress
type UserProgressWithManga struct {
	UserProgress
	Manga Manga `json:"manga"`
}

// ProgressUpdateRequest represents data for updating reading progress
type ProgressUpdateRequest struct {
	CurrentChapter *int           `json:"current_chapter" binding:"required,min=0"`
	Status         *ReadingStatus `json:"status"`
	Rating         *int           `json:"rating" binding:"omitempty,min=1,max=10"`
}

// LibraryAddRequest represents data for adding manga to library
type LibraryAddRequest struct {
	MangaID        string        `json:"manga_id" binding:"required"`
	Status         ReadingStatus `json:"status" binding:"required,oneof=reading completed plan_to_read on_hold dropped"`
	CurrentChapter int           `json:"current_chapter"`
}

// ProgressSyncMessage represents TCP sync protocol message
type ProgressSyncMessage struct {
	Type           string        `json:"type"` // "update", "request", "response"
	UserID         string        `json:"user_id,omitempty"`
	MangaID        string        `json:"manga_id"`
	CurrentChapter int           `json:"current_chapter"`
	Status         ReadingStatus `json:"status"`
	Timestamp      time.Time     `json:"timestamp"`
}

// GetRatingValue safely returns the rating value or 0 if NULL
func (p *UserProgress) GetRatingValue() int {
	if p.Rating.Valid {
		return int(p.Rating.Int64)
	}
	return 0
}

// GetStartedAtValue safely returns the started_at value or zero time if NULL
func (p *UserProgress) GetStartedAtValue() time.Time {
	if p.StartedAt.Valid {
		return p.StartedAt.Time
	}
	return time.Time{}
}

// GetCompletedAtValue safely returns the completed_at value or zero time if NULL
func (p *UserProgress) GetCompletedAtValue() time.Time {
	if p.CompletedAt.Valid {
		return p.CompletedAt.Time
	}
	return time.Time{}
}

// GetProgress retrieves user's progress for a specific manga
func (p *UserProgress) GetProgress(userID, mangaID string) (*UserProgress, error) {
	return &UserProgress{
		UserID:         userID,
		MangaID:        mangaID,
		CurrentChapter: p.CurrentChapter,
		Status:         p.Status,
		Rating:         p.Rating,
		StartedAt:      p.StartedAt,
		CompletedAt:    p.CompletedAt,
		UpdatedAt:      p.UpdatedAt,
	}, nil
}
