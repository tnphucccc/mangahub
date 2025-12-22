package stats

import (
	"database/sql"
	"fmt"
)

// Repository handles statistics data access
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new stats repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// UserStats represents reading statistics for a user
type UserStats struct {
	TotalManga        int     `json:"total_manga"`
	TotalChaptersRead int     `json:"total_chapters_read"`
	CompletedManga    int     `json:"completed_manga"`
	ReadingManga      int     `json:"reading_manga"`
	PlanToReadManga   int     `json:"plan_to_read_manga"`
	AverageRating     float64 `json:"average_rating"`
}

// GetUserStats calculates statistics for a specific user
func (r *Repository) GetUserStats(userID string) (*UserStats, error) {
	stats := &UserStats{}

	// 1. Count by status
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN status = 'reading' THEN 1 ELSE 0 END) as reading,
			SUM(CASE WHEN status = 'plan_to_read' THEN 1 ELSE 0 END) as plan_to_read,
			COALESCE(SUM(current_chapter), 0) as total_chapters,
			AVG(CASE WHEN rating > 0 THEN rating ELSE NULL END) as avg_rating
		FROM user_progress
		WHERE user_id = ?
	`

	var avgRating sql.NullFloat64
	err := r.db.QueryRow(query, userID).Scan(
		&stats.TotalManga,
		&stats.CompletedManga,
		&stats.ReadingManga,
		&stats.PlanToReadManga,
		&stats.TotalChaptersRead,
		&avgRating,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	if avgRating.Valid {
		stats.AverageRating = avgRating.Float64
	}

	return stats, nil
}
