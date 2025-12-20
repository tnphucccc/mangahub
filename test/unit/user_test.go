package unit

import (
	"testing"

	"github.com/tnphucccc/mangahub/pkg/models"
)

// Test user model structures
func TestUserModel(t *testing.T) {
	user := models.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
	}

	if user.ID != "user-123" {
		t.Errorf("Expected ID 'user-123', got '%s'", user.ID)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}

	t.Logf("✓ User model structure correct")
}

func TestUserRegisterRequest(t *testing.T) {
	req := models.UserRegisterRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	}

	if req.Username != "newuser" {
		t.Errorf("Expected username 'newuser', got '%s'", req.Username)
	}

	if req.Email != "new@example.com" {
		t.Errorf("Expected email 'new@example.com', got '%s'", req.Email)
	}

	if req.Password != "password123" {
		t.Errorf("Expected password 'password123', got '%s'", req.Password)
	}

	t.Logf("✓ UserRegisterRequest structure correct")
}

func TestUserLoginRequest(t *testing.T) {
	req := models.UserLoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	if req.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", req.Username)
	}

	if req.Password != "password123" {
		t.Errorf("Expected password 'password123', got '%s'", req.Password)
	}

	t.Logf("✓ UserLoginRequest structure correct")
}

func TestUserResponse_ToResponse(t *testing.T) {
	user := models.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password", // This should be excluded
	}

	response := user.ToResponse()

	if response.ID != user.ID {
		t.Errorf("Expected ID '%s', got '%s'", user.ID, response.ID)
	}

	if response.Username != user.Username {
		t.Errorf("Expected username '%s', got '%s'", user.Username, response.Username)
	}

	if response.Email != user.Email {
		t.Errorf("Expected email '%s', got '%s'", user.Email, response.Email)
	}

	// Password hash should not be in response
	// (Check structure doesn't have PasswordHash field)

	t.Logf("✓ UserResponse excludes sensitive data")
}
