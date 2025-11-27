package models

import (
	"encoding/json"
	"time"
)

// MangaStatus represents the publication status of a manga
type MangaStatus string

const (
	MangaStatusOngoing   MangaStatus = "ongoing"
	MangaStatusCompleted MangaStatus = "completed"
	MangaStatusHiatus    MangaStatus = "hiatus"
	MangaStatusCancelled MangaStatus = "cancelled"
)

// Manga represents a manga series in the catalog
type Manga struct {
	ID            string      `json:"id" db:"id"`
	Title         string      `json:"title" db:"title"`
	Author        string      `json:"author" db:"author"`
	Genres        []string    `json:"genres" db:"genres"` // Stored as JSON in DB
	Status        MangaStatus `json:"status" db:"status"`
	TotalChapters int         `json:"total_chapters" db:"total_chapters"`
	Description   string      `json:"description" db:"description"`
	CoverImageURL string      `json:"cover_image_url" db:"cover_image_url"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}

// MangaCreateRequest represents data for creating a new manga
type MangaCreateRequest struct {
	ID            string      `json:"id" binding:"required"`
	Title         string      `json:"title" binding:"required"`
	Author        string      `json:"author"`
	Genres        []string    `json:"genres"`
	Status        MangaStatus `json:"status" binding:"required,oneof=ongoing completed hiatus cancelled"`
	TotalChapters int         `json:"total_chapters"`
	Description   string      `json:"description"`
	CoverImageURL string      `json:"cover_image_url"`
}

// MangaUpdateRequest represents data for updating a manga
type MangaUpdateRequest struct {
	Title         *string      `json:"title"`
	Author        *string      `json:"author"`
	Genres        *[]string    `json:"genres"`
	Status        *MangaStatus `json:"status"`
	TotalChapters *int         `json:"total_chapters"`
	Description   *string      `json:"description"`
	CoverImageURL *string      `json:"cover_image_url"`
}

// MangaSearchQuery represents search parameters
type MangaSearchQuery struct {
	Query  string      `form:"q"`
	Genre  string      `form:"genre"`
	Status MangaStatus `form:"status"`
	Limit  int         `form:"limit"`
	Offset int         `form:"offset"`
}

// MarshalGenres converts genres slice to JSON string for database storage
func (m *Manga) MarshalGenres() (string, error) {
	if m.Genres == nil {
		return "[]", nil
	}
	data, err := json.Marshal(m.Genres)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalGenres converts JSON string from database to genres slice
func (m *Manga) UnmarshalGenres(data string) error {
	if data == "" {
		m.Genres = []string{}
		return nil
	}
	return json.Unmarshal([]byte(data), &m.Genres)
}
