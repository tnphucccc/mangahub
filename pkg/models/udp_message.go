package models

import "time"

// UDPMessageType represents the type of UDP message
type UDPMessageType string

const (
	// Client to Server messages
	UDPMessageTypeRegister   UDPMessageType = "register"   // Register for notifications
	UDPMessageTypeUnregister UDPMessageType = "unregister" // Unregister from notifications
	UDPMessageTypePing       UDPMessageType = "ping"       // Heartbeat ping

	// Server to Client messages
	UDPMessageTypeRegisterSuccess UDPMessageType = "register_success" // Registration successful
	UDPMessageTypeRegisterFailed  UDPMessageType = "register_failed"  // Registration failed
	UDPMessageTypePong            UDPMessageType = "pong"             // Heartbeat pong response
	UDPMessageTypeNotification    UDPMessageType = "notification"     // Chapter release notification
	UDPMessageTypeError           UDPMessageType = "error"            // Error message
)

// UDPMessage is the base message structure for UDP communication
type UDPMessage struct {
	Type      UDPMessageType `json:"type"`
	Timestamp time.Time      `json:"timestamp"`
	Data      interface{}    `json:"data,omitempty"`
}

// UDPRegisterMessage contains registration data
type UDPRegisterMessage struct {
	ClientID string `json:"client_id"` // Unique client identifier
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
}

// UDPRegisterSuccessMessage contains successful registration response
type UDPRegisterSuccessMessage struct {
	ClientID string `json:"client_id"`
	Message  string `json:"message"`
}

// UDPRegisterFailedMessage contains failed registration response
type UDPRegisterFailedMessage struct {
	Reason string `json:"reason"`
}

// UDPUnregisterMessage contains unregistration data
type UDPUnregisterMessage struct {
	ClientID string `json:"client_id"`
}

// UDPNotification contains chapter release notification data
type UDPNotification struct {
	MangaID       string    `json:"manga_id"`
	MangaTitle    string    `json:"manga_title"`
	ChapterNumber int       `json:"chapter_number"`
	ChapterTitle  string    `json:"chapter_title,omitempty"`
	ReleaseDate   time.Time `json:"release_date"`
	Message       string    `json:"message"`
}

// UDPErrorMessage contains error information
type UDPErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// UDPPingMessage is sent by client to check connectivity
type UDPPingMessage struct {
	ClientTime time.Time `json:"client_time"`
}

// UDPPongMessage is sent by server in response to ping
type UDPPongMessage struct {
	ServerTime time.Time `json:"server_time"`
	ClientTime time.Time `json:"client_time"`
}
