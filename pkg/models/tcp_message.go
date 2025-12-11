package models

import "time"

// TCPMessageType represents the type of TCP message
type TCPMessageType string

const (
	// Client to Server messages
	TCPMessageTypeAuth     TCPMessageType = "auth"     // Authentication with JWT token
	TCPMessageTypePing     TCPMessageType = "ping"     // Heartbeat ping
	TCPMessageTypeProgress TCPMessageType = "progress" // Progress update request

	// Server to Client messages
	TCPMessageTypeAuthSuccess TCPMessageType = "auth_success" // Authentication successful
	TCPMessageTypeAuthFailed  TCPMessageType = "auth_failed"  // Authentication failed
	TCPMessageTypePong        TCPMessageType = "pong"         // Heartbeat pong response
	TCPMessageTypeBroadcast   TCPMessageType = "broadcast"    // Progress broadcast to all clients
	TCPMessageTypeError       TCPMessageType = "error"        // Error message
)

// TCPMessage is the base message structure for TCP communication
type TCPMessage struct {
	Type      TCPMessageType `json:"type"`
	Timestamp time.Time      `json:"timestamp"`
	Data      interface{}    `json:"data,omitempty"`
}

// TCPAuthMessage contains authentication data
type TCPAuthMessage struct {
	Token string `json:"token"`
}

// TCPAuthSuccessMessage contains successful authentication response
type TCPAuthSuccessMessage struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// TCPAuthFailedMessage contains failed authentication response
type TCPAuthFailedMessage struct {
	Reason string `json:"reason"`
}

// TCPProgressMessage contains progress update data
type TCPProgressMessage struct {
	MangaID        string        `json:"manga_id"`
	CurrentChapter int           `json:"current_chapter"`
	Status         ReadingStatus `json:"status,omitempty"`
	Rating         *int          `json:"rating,omitempty"`
}

// TCPProgressBroadcast contains progress update broadcast to all clients
type TCPProgressBroadcast struct {
	UserID         string        `json:"user_id"`
	Username       string        `json:"username"`
	MangaID        string        `json:"manga_id"`
	MangaTitle     string        `json:"manga_title"`
	CurrentChapter int           `json:"current_chapter"`
	Status         ReadingStatus `json:"status"`
	Timestamp      time.Time     `json:"timestamp"`
}

// TCPErrorMessage contains error information
type TCPErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// TCPPingMessage is sent by client to keep connection alive
type TCPPingMessage struct {
	ClientTime time.Time `json:"client_time"`
}

// TCPPongMessage is sent by server in response to ping
type TCPPongMessage struct {
	ServerTime time.Time `json:"server_time"`
	ClientTime time.Time `json:"client_time"`
}
