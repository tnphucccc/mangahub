package models

import "time"

// Manga represents a manga entry.
type Manga struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Genres        []string  `json:"genres"`
	Status        string    `json:"status"`
	TotalChapters int       `json:"total_chapters"`
	Description   string    `json:"description"`
	CoverImageURL string    `json:"cover_image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// MangaSearchQuery represents the query parameters for manga search.
type MangaSearchQuery struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Genre  string `json:"genre"`
	Status string `json:"status"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type Meta struct {
	Total      int  `json:"total"`
	Count      int  `json:"count"`
	Limit      int  `json:"limit"`
	Offset     int  `json:"offset"`
	HasMore    bool `json:"has_more"`
	Page       int  `json:"page"`
	TotalPages int  `json:"total_pages"`
}

// MangaListResponse represents the response for a list of manga.
type MangaListResponse struct {
	Items []Manga `json:"items"`
}

// MangaDetailResponse represents the response for a single manga.
type MangaDetailResponse struct {
	Manga Manga `json:"manga"`
}
