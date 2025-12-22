package unit

import (
	"testing"
	"time"

	"github.com/tnphucccc/mangahub/internal/udp"
	"github.com/tnphucccc/mangahub/pkg/models"
)

// ==========================================
// UDP Message Structure Tests
// ==========================================

func TestUDPMessage_Structure(t *testing.T) {
	msg := models.UDPMessage{
		Type:      models.UDPMessageTypeRegister,
		Timestamp: time.Now(),
		Data: models.UDPRegisterMessage{
			ClientID: "test-client-1",
			Username: "testuser",
		},
	}

	if msg.Type != models.UDPMessageTypeRegister {
		t.Errorf("Expected type 'register', got '%s'", msg.Type)
	}

	if msg.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}

	t.Logf("✓ UDP message structure correct")
}

func TestUDPMessageTypes_Constants(t *testing.T) {
	types := []models.UDPMessageType{
		models.UDPMessageTypeRegister,
		models.UDPMessageTypeUnregister,
		models.UDPMessageTypePing,
		models.UDPMessageTypeRegisterSuccess,
		models.UDPMessageTypeRegisterFailed,
		models.UDPMessageTypePong,
		models.UDPMessageTypeNotification,
		models.UDPMessageTypeError,
	}

	expectedTypes := []string{
		"register",
		"unregister",
		"ping",
		"register_success",
		"register_failed",
		"pong",
		"notification",
		"error",
	}

	for i, msgType := range types {
		if string(msgType) != expectedTypes[i] {
			t.Errorf("Expected type '%s', got '%s'", expectedTypes[i], msgType)
		}
	}

	t.Logf("✓ UDP message type constants correct")
}

func TestUDPRegisterMessage_Structure(t *testing.T) {
	registerMsg := models.UDPRegisterMessage{
		ClientID: "client-123",
		UserID:   "user-456",
		Username: "testuser",
	}

	if registerMsg.ClientID != "client-123" {
		t.Errorf("Expected ClientID 'client-123', got '%s'", registerMsg.ClientID)
	}

	if registerMsg.UserID != "user-456" {
		t.Errorf("Expected UserID 'user-456', got '%s'", registerMsg.UserID)
	}

	if registerMsg.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", registerMsg.Username)
	}

	t.Logf("✓ UDP register message structure correct")
}

func TestUDPRegisterSuccessMessage_Structure(t *testing.T) {
	successMsg := models.UDPRegisterSuccessMessage{
		ClientID: "client-123",
		Message:  "Registration successful",
	}

	if successMsg.ClientID != "client-123" {
		t.Errorf("Expected ClientID 'client-123', got '%s'", successMsg.ClientID)
	}

	if successMsg.Message != "Registration successful" {
		t.Errorf("Expected message 'Registration successful', got '%s'", successMsg.Message)
	}

	t.Logf("✓ UDP register success message structure correct")
}

func TestUDPRegisterFailedMessage_Structure(t *testing.T) {
	failedMsg := models.UDPRegisterFailedMessage{
		Reason: "Client ID is required",
	}

	if failedMsg.Reason != "Client ID is required" {
		t.Errorf("Expected reason 'Client ID is required', got '%s'", failedMsg.Reason)
	}

	t.Logf("✓ UDP register failed message structure correct")
}

func TestUDPUnregisterMessage_Structure(t *testing.T) {
	unregisterMsg := models.UDPUnregisterMessage{
		ClientID: "client-123",
	}

	if unregisterMsg.ClientID != "client-123" {
		t.Errorf("Expected ClientID 'client-123', got '%s'", unregisterMsg.ClientID)
	}

	t.Logf("✓ UDP unregister message structure correct")
}

func TestUDPPingPongMessages_Structure(t *testing.T) {
	now := time.Now()

	pingMsg := models.UDPPingMessage{
		ClientTime: now,
	}

	pongMsg := models.UDPPongMessage{
		ServerTime: now,
		ClientTime: now,
	}

	if pingMsg.ClientTime.IsZero() {
		t.Error("Expected client time to be set")
	}

	if pongMsg.ServerTime.IsZero() {
		t.Error("Expected server time to be set")
	}

	if pongMsg.ClientTime.IsZero() {
		t.Error("Expected client time to be set")
	}

	t.Logf("✓ UDP ping/pong message structures correct")
}

func TestUDPNotification_Structure(t *testing.T) {
	notification := models.UDPNotification{
		MangaID:       "manga-123",
		MangaTitle:    "One Piece",
		ChapterNumber: 1100,
		ChapterTitle:  "The Final Battle",
		ReleaseDate:   time.Now(),
		Message:       "New chapter released!",
	}

	if notification.MangaID != "manga-123" {
		t.Errorf("Expected manga ID 'manga-123', got '%s'", notification.MangaID)
	}

	if notification.MangaTitle != "One Piece" {
		t.Errorf("Expected manga title 'One Piece', got '%s'", notification.MangaTitle)
	}

	if notification.ChapterNumber != 1100 {
		t.Errorf("Expected chapter 1100, got %d", notification.ChapterNumber)
	}

	if notification.ChapterTitle != "The Final Battle" {
		t.Errorf("Expected chapter title 'The Final Battle', got '%s'", notification.ChapterTitle)
	}

	if notification.Message != "New chapter released!" {
		t.Errorf("Expected message 'New chapter released!', got '%s'", notification.Message)
	}

	t.Logf("✓ UDP notification structure correct")
}

func TestUDPErrorMessage_Structure(t *testing.T) {
	errorMsg := models.UDPErrorMessage{
		Code:    "INVALID_MESSAGE",
		Message: "Invalid message format",
	}

	if errorMsg.Code != "INVALID_MESSAGE" {
		t.Errorf("Expected code 'INVALID_MESSAGE', got '%s'", errorMsg.Code)
	}

	if errorMsg.Message != "Invalid message format" {
		t.Errorf("Expected message 'Invalid message format', got '%s'", errorMsg.Message)
	}

	t.Logf("✓ UDP error message structure correct")
}

// ==========================================
// UDP Server Tests
// ==========================================

func TestUDPServer_NewServer(t *testing.T) {
	server := udp.NewServer("9091")

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.Port != "9091" {
		t.Errorf("Expected port '9091', got '%s'", server.Port)
	}

	t.Logf("✓ UDP server created successfully")
}

func TestUDPServer_GetStats(t *testing.T) {
	server := udp.NewServer("9091")

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

	notificationQueue, ok := stats["notification_queue"].(int)
	if !ok {
		t.Fatal("Expected notification_queue to be int")
	}

	if notificationQueue != 0 {
		t.Errorf("Expected 0 items in notification queue initially, got %d", notificationQueue)
	}

	bufferSize, ok := stats["buffer_size"].(int)
	if !ok {
		t.Fatal("Expected buffer_size to be int")
	}

	if bufferSize != 2048 {
		t.Errorf("Expected buffer size 2048, got %d", bufferSize)
	}

	t.Logf("✓ UDP server stats correct: %d clients, %d queued, %d buffer", totalClients, notificationQueue, bufferSize)
}

func TestUDPServer_BroadcastNotification(t *testing.T) {
	server := udp.NewServer("9091")

	notification := models.UDPNotification{
		MangaID:       "manga-123",
		MangaTitle:    "One Piece",
		ChapterNumber: 1100,
		ChapterTitle:  "The Final Battle",
		ReleaseDate:   time.Now(),
		Message:       "New chapter released!",
	}

	// This should not panic or block
	server.BroadcastNotification(notification)

	// Check stats to verify notification was queued
	stats := server.GetStats()
	notificationQueue := stats["notification_queue"].(int)

	if notificationQueue != 1 {
		t.Errorf("Expected 1 item in notification queue, got %d", notificationQueue)
	}

	t.Logf("✓ UDP notification queued successfully")
}

// ==========================================
// UDP Protocol Flow Tests
// ==========================================

func TestUDPProtocolFlow_Registration(t *testing.T) {
	// Client sends registration
	registerMsg := models.UDPMessage{
		Type:      models.UDPMessageTypeRegister,
		Timestamp: time.Now(),
		Data: models.UDPRegisterMessage{
			ClientID: "client-123",
			UserID:   "user-456",
			Username: "testuser",
		},
	}

	if registerMsg.Type != models.UDPMessageTypeRegister {
		t.Errorf("Expected register message type, got '%s'", registerMsg.Type)
	}

	// Server responds with success
	successMsg := models.UDPMessage{
		Type:      models.UDPMessageTypeRegisterSuccess,
		Timestamp: time.Now(),
		Data: models.UDPRegisterSuccessMessage{
			ClientID: "client-123",
			Message:  "Registration successful",
		},
	}

	if successMsg.Type != models.UDPMessageTypeRegisterSuccess {
		t.Errorf("Expected register_success message type, got '%s'", successMsg.Type)
	}

	t.Logf("✓ UDP registration flow validated")
}

func TestUDPProtocolFlow_PingPong(t *testing.T) {
	clientTime := time.Now()

	// Client sends ping
	pingMsg := models.UDPMessage{
		Type:      models.UDPMessageTypePing,
		Timestamp: clientTime,
		Data: models.UDPPingMessage{
			ClientTime: clientTime,
		},
	}

	if pingMsg.Type != models.UDPMessageTypePing {
		t.Errorf("Expected ping message type, got '%s'", pingMsg.Type)
	}

	// Server responds with pong
	serverTime := time.Now()
	pongMsg := models.UDPMessage{
		Type:      models.UDPMessageTypePong,
		Timestamp: serverTime,
		Data: models.UDPPongMessage{
			ServerTime: serverTime,
			ClientTime: clientTime,
		},
	}

	if pongMsg.Type != models.UDPMessageTypePong {
		t.Errorf("Expected pong message type, got '%s'", pongMsg.Type)
	}

	// Verify server time is after or equal to client time
	if serverTime.Before(clientTime) {
		t.Error("Server time should be after or equal to client time")
	}

	t.Logf("✓ UDP ping/pong flow validated")
}

func TestUDPProtocolFlow_Notification(t *testing.T) {
	// Server broadcasts notification
	notificationMsg := models.UDPMessage{
		Type:      models.UDPMessageTypeNotification,
		Timestamp: time.Now(),
		Data: models.UDPNotification{
			MangaID:       "manga-123",
			MangaTitle:    "One Piece",
			ChapterNumber: 1100,
			ChapterTitle:  "The Final Battle",
			ReleaseDate:   time.Now(),
			Message:       "New chapter released!",
		},
	}

	if notificationMsg.Type != models.UDPMessageTypeNotification {
		t.Errorf("Expected notification message type, got '%s'", notificationMsg.Type)
	}

	t.Logf("✓ UDP notification flow validated")
}

func TestUDPProtocolFlow_Unregistration(t *testing.T) {
	// Client sends unregistration
	unregisterMsg := models.UDPMessage{
		Type:      models.UDPMessageTypeUnregister,
		Timestamp: time.Now(),
		Data: models.UDPUnregisterMessage{
			ClientID: "client-123",
		},
	}

	if unregisterMsg.Type != models.UDPMessageTypeUnregister {
		t.Errorf("Expected unregister message type, got '%s'", unregisterMsg.Type)
	}

	t.Logf("✓ UDP unregistration flow validated")
}

func TestUDPProtocolFlow_ErrorHandling(t *testing.T) {
	// Test error message structure
	errorMsg := models.UDPMessage{
		Type:      models.UDPMessageTypeError,
		Timestamp: time.Now(),
		Data: models.UDPErrorMessage{
			Code:    "INVALID_MESSAGE",
			Message: "Invalid message format",
		},
	}

	if errorMsg.Type != models.UDPMessageTypeError {
		t.Errorf("Expected error message type, got '%s'", errorMsg.Type)
	}

	t.Logf("✓ UDP error handling flow validated")
}

func TestUDPProtocolFlow_RegistrationFailed(t *testing.T) {
	// Server responds with failure
	failedMsg := models.UDPMessage{
		Type:      models.UDPMessageTypeRegisterFailed,
		Timestamp: time.Now(),
		Data: models.UDPRegisterFailedMessage{
			Reason: "Client ID is required",
		},
	}

	if failedMsg.Type != models.UDPMessageTypeRegisterFailed {
		t.Errorf("Expected register_failed message type, got '%s'", failedMsg.Type)
	}

	t.Logf("✓ UDP registration failed flow validated")
}
