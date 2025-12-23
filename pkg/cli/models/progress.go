package models

// ProgressUpdateRequest represents the request body for updating reading progress.
type ProgressUpdateRequest struct {
	CurrentChapter int    `json:"current_chapter"`
	Status         string `json:"status,omitempty"`
	Rating         *int   `json:"rating,omitempty"`
}
