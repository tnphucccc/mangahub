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

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

const mangaSelectFields = `id, title, author, genres, status, total_chapters, description, cover_image_url, created_at, updated_at`

func scanManga(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Manga, error) {
	var manga models.Manga
	var genresJSON string

	err := scanner.Scan(
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
		return nil, err
	}

	// Unmarshal genres JSON
	if err := manga.UnmarshalGenres(genresJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal genres: %w", err)
	}

	return &manga, nil
}

func (r *Repository) FindByID(id string) (*models.Manga, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM manga
		WHERE id = ?
	`, mangaSelectFields)

	manga, err := scanManga(r.db.QueryRow(query, id))
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("manga not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find manga: %w", err)
	}

	return manga, nil
}

func (r *Repository) Search(query models.MangaSearchQuery) ([]models.Manga, int, error) {
	whereClauses := []string{"1=1"}
	args := []interface{}{}

	// Add title search (case-insensitive partial match)
	if query.Title != "" {
		whereClauses = append(whereClauses, "LOWER(title) LIKE LOWER(?)")
		args = append(args, "%"+query.Title+"%")
	}

	// Add author search (case-insensitive partial match)
	if query.Author != "" {
		whereClauses = append(whereClauses, "LOWER(author) LIKE LOWER(?)")
		args = append(args, "%"+query.Author+"%")
	}

	// Add genre filter (searches within JSON array)
	if query.Genre != "" {
		whereClauses = append(whereClauses, "LOWER(genres) LIKE LOWER(?)")
		args = append(args, "%"+query.Genre+"%")
	}

	// Add status filter (exact match)
	if query.Status != "" {
		whereClauses = append(whereClauses, "status = ?")
		args = append(args, query.Status)
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	// 1. Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM manga WHERE %s", whereSQL)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get count: %w", err)
	}

	// 2. Get data
	sqlQuery := fmt.Sprintf(`
		SELECT %s
		FROM manga
		WHERE %s
		ORDER BY title ASC
	`, mangaSelectFields, whereSQL)

	// Add pagination only if limit > 0
	if query.Limit > 0 {
		sqlQuery += " LIMIT ?"
		args = append(args, query.Limit)

		if query.Offset > 0 {
			sqlQuery += " OFFSET ?"
			args = append(args, query.Offset)
		}
	}

	mangaList, err := r.queryMangaList(sqlQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	return mangaList, total, nil
}

// queryMangaList executes a query and returns a list of manga (reduces duplication)
func (r *Repository) queryMangaList(query string, args ...interface{}) ([]models.Manga, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query manga: %w", err)
	}
	defer rows.Close()

	var mangaList []models.Manga
	for rows.Next() {
		manga, err := scanManga(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan manga: %w", err)
		}
		mangaList = append(mangaList, *manga)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating manga rows: %w", err)
	}

	return mangaList, nil
}

// FindAll retrieves all manga with optional pagination
func (r *Repository) FindAll(limit, offset int) ([]models.Manga, int, error) {
	// 1. Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM manga").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// 2. Get data
	query := fmt.Sprintf(`
		SELECT %s
		FROM manga
		ORDER BY title ASC
	`, mangaSelectFields)

	args := []interface{}{}
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)

		if offset > 0 {
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}

	mangaList, err := r.queryMangaList(query, args...)
	if err != nil {
		return nil, 0, err
	}

	return mangaList, total, nil
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
