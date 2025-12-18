package user

import (
	"database/sql"
	"fmt"

	"github.com/tnphucccc/mangahub/pkg/models"
)

// Repository handles user data access
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new user repository
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create creates a new user
func (r *Repository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// scanUser scans a row into a User model (reduces code duplication)
func (r *Repository) scanUser(row *sql.Row) (*models.User, error) {
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}
	return &user, nil
}

// findByField is a generic finder that reduces duplication for FindByID, FindByUsername, FindByEmail
func (r *Repository) findByField(field string, value interface{}) (*models.User, error) {
	query := fmt.Sprintf(`
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE %s = ?
	`, field)
	return r.scanUser(r.db.QueryRow(query, value))
}

// FindByID finds a user by ID
func (r *Repository) FindByID(id string) (*models.User, error) {
	return r.findByField("id", id)
}

// FindByUsername finds a user by username
func (r *Repository) FindByUsername(username string) (*models.User, error) {
	return r.findByField("username", username)
}

// FindByEmail finds a user by email
func (r *Repository) FindByEmail(email string) (*models.User, error) {
	return r.findByField("email", email)
}

// Exists checks if a user with given username or email already exists
func (r *Repository) Exists(username, email string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM users WHERE username = ? OR email = ?
	`
	var count int
	err := r.db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}
