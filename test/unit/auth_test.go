package unit

import (
	"testing"
	"time"

	"github.com/tnphucccc/mangahub/internal/auth"
	"github.com/tnphucccc/mangahub/pkg/models"
)

// ==========================================
// JWT Tests
// ==========================================

func TestJWT_GenerateToken(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret-key", 7)

	testUser := &models.User{
		ID:       "test-user-123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	token, err := jwtManager.GenerateToken(testUser)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if token == "" {
		t.Fatal("Expected token to be non-empty")
	}

	t.Logf("✓ Generated token successfully")
}

func TestJWT_ValidateToken(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret-key", 7)

	testUser := &models.User{
		ID:       "test-user-123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	token, _ := jwtManager.GenerateToken(testUser)

	t.Run("valid token", func(t *testing.T) {
		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			t.Fatalf("Expected no error for valid token, got: %v", err)
		}

		if claims.UserID != testUser.ID {
			t.Errorf("Expected UserID %s, got %s", testUser.ID, claims.UserID)
		}

		if claims.Username != testUser.Username {
			t.Errorf("Expected Username %s, got %s", testUser.Username, claims.Username)
		}

		t.Logf("✓ Valid token verified")
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := jwtManager.ValidateToken("invalid.token.here")
		if err == nil {
			t.Fatal("Expected error for invalid token, got nil")
		}
		t.Logf("✓ Invalid token rejected")
	})

	t.Run("empty token", func(t *testing.T) {
		_, err := jwtManager.ValidateToken("")
		if err == nil {
			t.Fatal("Expected error for empty token, got nil")
		}
		t.Logf("✓ Empty token rejected")
	})
}

func TestJWT_TokenExpiration(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret-key", 0) // 0 days = immediate expiry

	testUser := &models.User{
		ID:       "test-user-123",
		Username: "testuser",
	}

	token, _ := jwtManager.GenerateToken(testUser)
	time.Sleep(1 * time.Second)

	_, err := jwtManager.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error for expired token, got nil")
	}

	t.Logf("✓ Expired token correctly rejected")
}

func TestJWT_RefreshToken(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret-key", 7)

	testUser := &models.User{
		ID:       "test-user-123",
		Username: "testuser",
	}

	oldToken, _ := jwtManager.GenerateToken(testUser)
	oldClaims, _ := jwtManager.ValidateToken(oldToken)

	// Wait a bit to ensure different timestamp
	time.Sleep(1 * time.Second)

	newToken, err := jwtManager.RefreshToken(oldToken)

	if err != nil {
		t.Fatalf("Failed to refresh token: %v", err)
	}

	if newToken == "" {
		t.Fatal("New token should not be empty")
	}

	newClaims, err := jwtManager.ValidateToken(newToken)
	if err != nil {
		t.Fatalf("Failed to validate refreshed token: %v", err)
	}

	if newClaims.UserID != oldClaims.UserID {
		t.Errorf("Expected UserID %s, got %s", oldClaims.UserID, newClaims.UserID)
	}

	t.Logf("✓ Token refreshed successfully")
}

// ==========================================
// Password Tests
// ==========================================

func TestPassword_Hash(t *testing.T) {
	password := "mysecurepassword123"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if hash == "" {
		t.Fatal("Expected hash to be non-empty")
	}

	if hash == password {
		t.Fatal("Hash should not be the same as plain password")
	}

	t.Logf("✓ Password hashed successfully")
}

func TestPassword_Compare(t *testing.T) {
	password := "mysecurepassword123"
	hash, _ := auth.HashPassword(password)

	t.Run("correct password", func(t *testing.T) {
		err := auth.ComparePassword(hash, password)
		if err != nil {
			t.Errorf("Expected no error for correct password, got: %v", err)
		}
		t.Logf("✓ Correct password verified")
	})

	t.Run("wrong password", func(t *testing.T) {
		err := auth.ComparePassword(hash, "wrongpassword")
		if err == nil {
			t.Fatal("Expected error for wrong password, got nil")
		}
		t.Logf("✓ Wrong password rejected")
	})

	t.Run("empty password", func(t *testing.T) {
		err := auth.ComparePassword(hash, "")
		if err == nil {
			t.Fatal("Expected error for empty password, got nil")
		}
		t.Logf("✓ Empty password rejected")
	})
}

func TestPassword_DifferentHashes(t *testing.T) {
	password := "mysecurepassword123"

	hash1, _ := auth.HashPassword(password)
	hash2, _ := auth.HashPassword(password)

	// Bcrypt generates different hashes due to random salt
	if hash1 == hash2 {
		t.Error("Expected different hashes for same password (due to salt)")
	}

	// But both should be valid
	if err := auth.ComparePassword(hash1, password); err != nil {
		t.Errorf("First hash should be valid: %v", err)
	}

	if err := auth.ComparePassword(hash2, password); err != nil {
		t.Errorf("Second hash should be valid: %v", err)
	}

	t.Logf("✓ Different hashes verified")
}

func TestPassword_SpecialCharacters(t *testing.T) {
	passwords := []string{
		"p@ssw0rd!#$%",
		"pass word with spaces",
		"pass\nword\t123",
	}

	for _, password := range passwords {
		hash, err := auth.HashPassword(password)
		if err != nil {
			t.Fatalf("Failed to hash password: %v", err)
		}

		if err := auth.ComparePassword(hash, password); err != nil {
			t.Errorf("Failed to verify password with special characters: %v", err)
		}
	}

	t.Logf("✓ Special character passwords handled correctly")
}
