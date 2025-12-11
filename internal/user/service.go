package user

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tnphucccc/mangahub/internal/auth"
	"github.com/tnphucccc/mangahub/pkg/models"
)

// Service handles user business logic
type Service struct {
	repo       *Repository
	jwtManager *auth.JWTManager
}

// NewService creates a new user service
func NewService(repo *Repository, jwtManager *auth.JWTManager) *Service {
	return &Service{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// Register registers a new user
func (s *Service) Register(req models.UserRegisterRequest) (*models.User, string, error) {
	// Check if user already exists
	exists, err := s.repo.Exists(req.Username, req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, "", fmt.Errorf("username or email already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	// Save to database
	if err := s.repo.Create(user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// Login authenticates a user
func (s *Service) Login(req models.UserLoginRequest) (*models.User, string, error) {
	// Find user by username
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, "", fmt.Errorf("invalid username or password")
	}

	// Compare password
	if err := auth.ComparePassword(user.PasswordHash, req.Password); err != nil {
		return nil, "", fmt.Errorf("invalid username or password")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// GetByID retrieves a user by ID
func (s *Service) GetByID(id string) (*models.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *Service) ValidateToken(tokenString string) (*models.User, error) {
	// Validate token
	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user from database
	user, err := s.repo.FindByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
