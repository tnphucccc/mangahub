package manga

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tnphucccc/mangahub/pkg/models"
)

// Repository handles manga data access
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new manga repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// FindByID finds a manga by ID
func (r *Repository) FindByID(id string) (*models.Manga, error) {
	query := `
		SELECT id, title, author, genres, status, total_chapters, description, cover_image_url, created_at, updated_at
		FROM manga
		WHERE id = ?
	`
	var manga models.Manga
	var genresJSON string

	err := r.db.QueryRow(query, id).Scan(
		&manga.ID,
		&manga.Title,
		&manga.Author,
		&genresJSON,
		&manga.Status,
		&manga.TotalChapters,
		&manga.Description,
		&manga.CoverImageURL,
		&manga.CreatedAt,
		&manga.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("manga not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find manga: %w", err)
	}

	// Unmarshal genres JSON
	if err := manga.UnmarshalGenres(genresJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal genres: %w", err)
	}

	return &manga, nil
}

// Search searches for manga by query
func (r *Repository) Search(query models.MangaSearchQuery) ([]models.Manga, error) {
	// Build SQL query
	sqlQuery := `
		SELECT id, title, author, genres, status, total_chapters, description, cover_image_url, created_at, updated_at
		FROM manga
		WHERE 1=1
	`
	args := []interface{}{}

	// Add title search
	if query.Query != "" {
		sqlQuery += " AND (title LIKE ? OR author LIKE ?)"
		searchTerm := "%" + query.Query + "%"
		args = append(args, searchTerm, searchTerm)
	}

	// Add genre filter
	if query.Genre != "" {
		sqlQuery += " AND genres LIKE ?"
		args = append(args, "%"+query.Genre+"%")
	}

	// Add status filter
	if query.Status != "" {
		sqlQuery += " AND status = ?"
		args = append(args, query.Status)
	}

	// Add ordering
	sqlQuery += " ORDER BY title ASC"

	// Add pagination
	if query.Limit > 0 {
		sqlQuery += " LIMIT ?"
		args = append(args, query.Limit)

		if query.Offset > 0 {
			sqlQuery += " OFFSET ?"
			args = append(args, query.Offset)
		}
	}

	// Execute query
	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search manga: %w", err)
	}
	defer rows.Close()

	// Scan results
	var mangaList []models.Manga
	for rows.Next() {
		var manga models.Manga
		var genresJSON string

		err := rows.Scan(
			&manga.ID,
			&manga.Title,
			&manga.Author,
			&genresJSON,
			&manga.Status,
			&manga.TotalChapters,
			&manga.Description,
			&manga.CoverImageURL,
			&manga.CreatedAt,
			&manga.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan manga: %w", err)
		}

		// Unmarshal genres
		if err := manga.UnmarshalGenres(genresJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genres: %w", err)
		}

		mangaList = append(mangaList, manga)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating manga rows: %w", err)
	}

	return mangaList, nil
}

// FindAll retrieves all manga
func (r *Repository) FindAll(limit, offset int) ([]models.Manga, error) {
	query := `
		SELECT id, title, author, genres, status, total_chapters, description, cover_image_url, created_at, updated_at
		FROM manga
		ORDER BY title ASC
	`

	args := []interface{}{}
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)

		if offset > 0 {
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find all manga: %w", err)
	}
	defer rows.Close()

	var mangaList []models.Manga
	for rows.Next() {
		var manga models.Manga
		var genresJSON string

		err := rows.Scan(
			&manga.ID,
			&manga.Title,
			&manga.Author,
			&genresJSON,
			&manga.Status,
			&manga.TotalChapters,
			&manga.Description,
			&manga.CoverImageURL,
			&manga.CreatedAt,
			&manga.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan manga: %w", err)
		}

		if err := manga.UnmarshalGenres(genresJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genres: %w", err)
		}

		mangaList = append(mangaList, manga)
	}

	return mangaList, rows.Err()
}

// GetUserLibrary retrieves a user's manga library with progress
func (r *Repository) GetUserLibrary(userID string) ([]models.UserProgressWithManga, error) {
	query := `
		SELECT
			up.user_id, up.manga_id, up.current_chapter, up.status, up.rating,
			up.started_at, up.completed_at, up.updated_at,
			m.id, m.title, m.author, m.genres, m.status, m.total_chapters,
			m.description, m.cover_image_url, m.created_at, m.updated_at
		FROM user_progress up
		JOIN manga m ON up.manga_id = m.id
		WHERE up.user_id = ?
		ORDER BY up.updated_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user library: %w", err)
	}
	defer rows.Close()

	var library []models.UserProgressWithManga
	for rows.Next() {
		var item models.UserProgressWithManga
		var genresJSON string

		err := rows.Scan(
			&item.UserID,
			&item.MangaID,
			&item.CurrentChapter,
			&item.Status,
			&item.Rating,
			&item.StartedAt,
			&item.CompletedAt,
			&item.UpdatedAt,
			&item.Manga.ID,
			&item.Manga.Title,
			&item.Manga.Author,
			&genresJSON,
			&item.Manga.Status,
			&item.Manga.TotalChapters,
			&item.Manga.Description,
			&item.Manga.CoverImageURL,
			&item.Manga.CreatedAt,
			&item.Manga.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan library item: %w", err)
		}

		if err := item.Manga.UnmarshalGenres(genresJSON); err != nil {
			return nil, fmt.Errorf("failed to unmarshal genres: %w", err)
		}

		library = append(library, item)
	}

	return library, rows.Err()
}

// AddToLibrary adds a manga to user's library
func (r *Repository) AddToLibrary(userID, mangaID string, status models.ReadingStatus, currentChapter int) error {
	query := `
		INSERT INTO user_progress (user_id, manga_id, status, current_chapter, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id, manga_id) DO UPDATE SET
			status = excluded.status,
			current_chapter = excluded.current_chapter,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(query, userID, mangaID, status, currentChapter)
	if err != nil {
		return fmt.Errorf("failed to add manga to library: %w", err)
	}

	return nil
}

// UpdateProgress updates user's reading progress
func (r *Repository) UpdateProgress(userID, mangaID string, currentChapter int, status *models.ReadingStatus, rating *int) error {
	// Build dynamic update query
	updates := []string{"current_chapter = ?", "updated_at = CURRENT_TIMESTAMP"}
	args := []interface{}{currentChapter}

	if status != nil {
		updates = append(updates, "status = ?")
		args = append(args, *status)
	}

	if rating != nil {
		updates = append(updates, "rating = ?")
		args = append(args, *rating)
	}

	// Add WHERE clause parameters
	args = append(args, userID, mangaID)

	query := fmt.Sprintf(`
		UPDATE user_progress
		SET %s
		WHERE user_id = ? AND manga_id = ?
	`, strings.Join(updates, ", "))

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("manga not in user's library")
	}

	return nil
}

// GetProgress retrieves user's progress for a specific manga
func (r *Repository) GetProgress(userID, mangaID string) (*models.UserProgress, error) {
	query := `
		SELECT user_id, manga_id, current_chapter, status, rating, started_at, completed_at, updated_at
		FROM user_progress
		WHERE user_id = ? AND manga_id = ?
	`

	var progress models.UserProgress
	err := r.db.QueryRow(query, userID, mangaID).Scan(
		&progress.UserID,
		&progress.MangaID,
		&progress.CurrentChapter,
		&progress.Status,
		&progress.Rating,
		&progress.StartedAt,
		&progress.CompletedAt,
		&progress.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("progress not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get progress: %w", err)
	}

	return &progress, nil
}
