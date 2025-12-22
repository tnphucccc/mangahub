package models

// LibraryAddRequest represents the request body for adding a manga to the library.
type LibraryAddRequest struct {
	MangaID        string `json:"manga_id"`
	Status         string `json:"status"`
	CurrentChapter int    `json:"current_chapter"`
}

// NullInt64 is a helper for unmarshalling sql.NullInt64 from JSON.
type NullInt64 struct {
	Int64 int64
	Valid bool
}

// UserProgress represents a user's progress on a manga.
type UserProgress struct {
	UserID         string    `json:"user_id"`
	MangaID        string    `json:"manga_id"`
	CurrentChapter int       `json:"current_chapter"`
	Status         string    `json:"status"`
	Rating         NullInt64 `json:"rating"`
	UpdatedAt      string    `json:"updated_at"` // Using string for simplicity, can be time.Time
}

// UserProgressWithManga combines UserProgress with Manga details.
type UserProgressWithManga struct {
	UserProgress
	Manga Manga `json:"manga"`
}

// LibraryListResponse represents the response for a list of library entries.
type LibraryListResponse struct {
	Items []UserProgressWithManga `json:"items"`
}
