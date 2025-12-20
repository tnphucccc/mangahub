package unit

import (
	"testing"
	"time"

	"github.com/tnphucccc/mangahub/internal/auth"
	"github.com/tnphucccc/mangahub/internal/tcp"
	"github.com/tnphucccc/mangahub/pkg/models"
)

// ==========================================
// TCP Message Structure Tests
// ==========================================

func TestTCPMessage_Structure(t *testing.T) {
	msg := models.TCPMessage{
		Type:      models.TCPMessageTypeAuth,
		Timestamp: time.Now(),
		Data: models.TCPAuthMessage{
			Token: "test-token",
		},
	}

	if msg.Type != models.TCPMessageTypeAuth {
		t.Errorf("Expected type 'auth', got '%s'", msg.Type)
	}

	if msg.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}

	t.Logf("✓ TCP message structure correct")
}

func TestTCPMessageTypes_Constants(t *testing.T) {
	types := []models.TCPMessageType{
		models.TCPMessageTypeAuth,
		models.TCPMessageTypeAuthSuccess,
		models.TCPMessageTypeAuthFailed,
		models.TCPMessageTypePing,
		models.TCPMessageTypePong,
		models.TCPMessageTypeProgress,
		models.TCPMessageTypeBroadcast,
		models.TCPMessageTypeError,
	}

	expectedTypes := []string{
		"auth",
		"auth_success",
		"auth_failed",
		"ping",
		"pong",
		"progress",
		"broadcast",
		"error",
	}

	for i, msgType := range types {
		if string(msgType) != expectedTypes[i] {
			t.Errorf("Expected type '%s', got '%s'", expectedTypes[i], msgType)
		}
	}

	t.Logf("✓ TCP message type constants correct")
}

func TestTCPAuthMessage_Structure(t *testing.T) {
	authMsg := models.TCPAuthMessage{
		Token: "jwt-token-here",
	}

	if authMsg.Token != "jwt-token-here" {
		t.Errorf("Expected token 'jwt-token-here', got '%s'", authMsg.Token)
	}

	t.Logf("✓ TCP auth message structure correct")
}

func TestTCPAuthSuccessMessage_Structure(t *testing.T) {
	successMsg := models.TCPAuthSuccessMessage{
		UserID:   "user-123",
		Username: "testuser",
		Message:  "Authentication successful",
	}

	if successMsg.UserID != "user-123" {
		t.Errorf("Expected UserID 'user-123', got '%s'", successMsg.UserID)
	}

	if successMsg.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", successMsg.Username)
	}

	t.Logf("✓ TCP auth success message structure correct")
}

func TestTCPAuthFailedMessage_Structure(t *testing.T) {
	failedMsg := models.TCPAuthFailedMessage{
		Reason: "Invalid token",
	}

	if failedMsg.Reason != "Invalid token" {
		t.Errorf("Expected reason 'Invalid token', got '%s'", failedMsg.Reason)
	}

	t.Logf("✓ TCP auth failed message structure correct")
}

func TestTCPPingPongMessages_Structure(t *testing.T) {
	now := time.Now()

	pingMsg := models.TCPPingMessage{
		ClientTime: now,
	}

	pongMsg := models.TCPPongMessage{
		ServerTime: now,
	}

	if pingMsg.ClientTime.IsZero() {
		t.Error("Expected client time to be set")
	}

	if pongMsg.ServerTime.IsZero() {
		t.Error("Expected server time to be set")
	}

	t.Logf("✓ TCP ping/pong message structures correct")
}

func TestTCPProgressMessage_Structure(t *testing.T) {
	progressMsg := models.TCPProgressMessage{
		MangaID:        "manga-123",
		CurrentChapter: 50,
		Status:         models.ReadingStatusReading,
	}

	if progressMsg.MangaID != "manga-123" {
		t.Errorf("Expected manga ID 'manga-123', got '%s'", progressMsg.MangaID)
	}

	if progressMsg.CurrentChapter != 50 {
		t.Errorf("Expected chapter 50, got %d", progressMsg.CurrentChapter)
	}

	if progressMsg.Status != models.ReadingStatusReading {
		t.Errorf("Expected status 'reading', got '%s'", progressMsg.Status)
	}

	t.Logf("✓ TCP progress message structure correct")
}

func TestTCPProgressBroadcast_Structure(t *testing.T) {
	broadcast := models.TCPProgressBroadcast{
		UserID:         "user-123",
		Username:       "testuser",
		MangaID:        "manga-123",
		CurrentChapter: 50,
		Status:         models.ReadingStatusReading,
		Timestamp:      time.Now(),
	}

	if broadcast.UserID != "user-123" {
		t.Errorf("Expected UserID 'user-123', got '%s'", broadcast.UserID)
	}

	if broadcast.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", broadcast.Username)
	}

	if broadcast.MangaID != "manga-123" {
		t.Errorf("Expected manga ID 'manga-123', got '%s'", broadcast.MangaID)
	}

	if broadcast.CurrentChapter != 50 {
		t.Errorf("Expected chapter 50, got %d", broadcast.CurrentChapter)
	}

	t.Logf("✓ TCP progress broadcast structure correct")
}

func TestTCPErrorMessage_Structure(t *testing.T) {
	errorMsg := models.TCPErrorMessage{
		Code:    "INVALID_MESSAGE",
		Message: "Invalid message format",
	}

	if errorMsg.Code != "INVALID_MESSAGE" {
		t.Errorf("Expected code 'INVALID_MESSAGE', got '%s'", errorMsg.Code)
	}

	if errorMsg.Message != "Invalid message format" {
		t.Errorf("Expected message 'Invalid message format', got '%s'", errorMsg.Message)
	}

	t.Logf("✓ TCP error message structure correct")
}

// ==========================================
// TCP Server Tests
// ==========================================

func TestTCPServer_NewServer(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", 7)
	server := tcp.NewServer("9090", jwtManager)

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.Port != "9090" {
		t.Errorf("Expected port '9090', got '%s'", server.Port)
	}

	t.Logf("✓ TCP server created successfully")
}

func TestTCPServer_GetStats(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", 7)
	server := tcp.NewServer("9090", jwtManager)

	stats := server.GetStats()

	if stats == nil {
		t.Fatal("Expected stats to be returned, got nil")
	}

	totalClients, ok := stats["total_clients"].(int)
	if !ok {
		t.Fatal("Expected total_clients to be int")
	}

	if totalClients != 0 {
		t.Errorf("Expected 0 clients initially, got %d", totalClients)
	}

	totalUsers, ok := stats["total_users"].(int)
	if !ok {
		t.Fatal("Expected total_users to be int")
	}

	if totalUsers != 0 {
		t.Errorf("Expected 0 users initially, got %d", totalUsers)
	}

	broadcastQueue, ok := stats["broadcast_queue"].(int)
	if !ok {
		t.Fatal("Expected broadcast_queue to be int")
	}

	if broadcastQueue != 0 {
		t.Errorf("Expected 0 items in broadcast queue initially, got %d", broadcastQueue)
	}

	t.Logf("✓ TCP server stats correct: %d clients, %d users, %d queued", totalClients, totalUsers, broadcastQueue)
}

func TestTCPServer_BroadcastProgress(t *testing.T) {
	jwtManager := auth.NewJWTManager("test-secret", 7)
	server := tcp.NewServer("9090", jwtManager)

	broadcast := models.TCPProgressBroadcast{
		UserID:         "user-123",
		Username:       "testuser",
		MangaID:        "manga-123",
		CurrentChapter: 50,
		Status:         models.ReadingStatusReading,
		Timestamp:      time.Now(),
	}

	// This should not panic or block
	server.BroadcastProgress(broadcast)

	// Check stats to verify broadcast was queued
	stats := server.GetStats()
	broadcastQueue := stats["broadcast_queue"].(int)

	if broadcastQueue != 1 {
		t.Errorf("Expected 1 item in broadcast queue, got %d", broadcastQueue)
	}

	t.Logf("✓ TCP broadcast queued successfully")
}

// ==========================================
// TCP Protocol Flow Tests
// ==========================================

func TestTCPProtocolFlow_Authentication(t *testing.T) {
	// Test the expected message flow for authentication
	jwtManager := auth.NewJWTManager("test-secret", 7)

	// Create test user and generate token
	user := &models.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	token, err := jwtManager.GenerateToken(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Create auth message
	authMsg := models.TCPMessage{
		Type:      models.TCPMessageTypeAuth,
		Timestamp: time.Now(),
		Data: models.TCPAuthMessage{
			Token: token,
		},
	}

	if authMsg.Type != models.TCPMessageTypeAuth {
		t.Errorf("Expected auth message type, got '%s'", authMsg.Type)
	}

	// Validate token (simulating what server would do)
	claims, err := jwtManager.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("Expected UserID '%s', got '%s'", user.ID, claims.UserID)
	}

	// Create success response (what server would send)
	successMsg := models.TCPMessage{
		Type:      models.TCPMessageTypeAuthSuccess,
		Timestamp: time.Now(),
		Data: models.TCPAuthSuccessMessage{
			UserID:   claims.UserID,
			Username: claims.Username,
			Message:  "Authentication successful",
		},
	}

	if successMsg.Type != models.TCPMessageTypeAuthSuccess {
		t.Errorf("Expected auth_success message type, got '%s'", successMsg.Type)
	}

	t.Logf("✓ TCP authentication flow validated")
}

func TestTCPProtocolFlow_PingPong(t *testing.T) {
	clientTime := time.Now()

	// Client sends ping
	pingMsg := models.TCPMessage{
		Type:      models.TCPMessageTypePing,
		Timestamp: clientTime,
		Data: models.TCPPingMessage{
			ClientTime: clientTime,
		},
	}

	if pingMsg.Type != models.TCPMessageTypePing {
		t.Errorf("Expected ping message type, got '%s'", pingMsg.Type)
	}

	// Server responds with pong
	serverTime := time.Now()
	pongMsg := models.TCPMessage{
		Type:      models.TCPMessageTypePong,
		Timestamp: serverTime,
		Data: models.TCPPongMessage{
			ServerTime: serverTime,
		},
	}

	if pongMsg.Type != models.TCPMessageTypePong {
		t.Errorf("Expected pong message type, got '%s'", pongMsg.Type)
	}

	// Verify server time is after client time
	if serverTime.Before(clientTime) {
		t.Error("Server time should be after or equal to client time")
	}

	t.Logf("✓ TCP ping/pong flow validated")
}

func TestTCPProtocolFlow_ProgressUpdate(t *testing.T) {
	// Client sends progress update
	progressMsg := models.TCPMessage{
		Type:      models.TCPMessageTypeProgress,
		Timestamp: time.Now(),
		Data: models.TCPProgressMessage{
			MangaID:        "manga-123",
			CurrentChapter: 75,
			Status:         models.ReadingStatusReading,
		},
	}

	if progressMsg.Type != models.TCPMessageTypeProgress {
		t.Errorf("Expected progress message type, got '%s'", progressMsg.Type)
	}

	// Server broadcasts to all clients
	broadcastMsg := models.TCPMessage{
		Type:      models.TCPMessageTypeBroadcast,
		Timestamp: time.Now(),
		Data: models.TCPProgressBroadcast{
			UserID:         "user-123",
			Username:       "testuser",
			MangaID:        "manga-123",
			CurrentChapter: 75,
			Status:         models.ReadingStatusReading,
			Timestamp:      time.Now(),
		},
	}

	if broadcastMsg.Type != models.TCPMessageTypeBroadcast {
		t.Errorf("Expected broadcast message type, got '%s'", broadcastMsg.Type)
	}

	t.Logf("✓ TCP progress update flow validated")
}

func TestTCPProtocolFlow_ErrorHandling(t *testing.T) {
	// Test error message structure
	errorMsg := models.TCPMessage{
		Type:      models.TCPMessageTypeError,
		Timestamp: time.Now(),
		Data: models.TCPErrorMessage{
			Code:    "INVALID_MESSAGE",
			Message: "Invalid message format",
		},
	}

	if errorMsg.Type != models.TCPMessageTypeError {
		t.Errorf("Expected error message type, got '%s'", errorMsg.Type)
	}

	t.Logf("✓ TCP error handling flow validated")
}
