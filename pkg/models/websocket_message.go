package models

import "time"

// WebSocketMessageType represents the type of WebSocket message
type WebSocketMessageType string

const (
	// WSJoinRoom is sent when a client wants to join a chat room
	WSJoinRoom WebSocketMessageType = "join_room"
	// WSLeaveRoom is sent when a client wants to leave a chat room
	WSLeaveRoom WebSocketMessageType = "leave_room"
	// WSChatMessage is a regular chat message
	WSChatMessage WebSocketMessageType = "chat"
	// WSSystemMessage is a system notification (user joined, left, etc.)
	WSSystemMessage WebSocketMessageType = "system"
	// WSError is an error message from the server
	WSError WebSocketMessageType = "error"
)

// WebSocketMessage is the base message structure for WebSocket communication
type WebSocketMessage struct {
	Type      WebSocketMessageType `json:"type"`
	Username  string               `json:"username,omitempty"`
	Room      string               `json:"room,omitempty"`
	Content   string               `json:"content"`
	Timestamp time.Time            `json:"timestamp"`
}

// NewChatMessage creates a new chat message
func NewChatMessage(username, room, content string) WebSocketMessage {
	return WebSocketMessage{
		Type:      WSChatMessage,
		Username:  username,
		Room:      room,
		Content:   content,
		Timestamp: time.Now(),
	}
}

// NewSystemMessage creates a new system message
func NewSystemMessage(room, content string) WebSocketMessage {
	return WebSocketMessage{
		Type:      WSSystemMessage,
		Room:      room,
		Content:   content,
		Timestamp: time.Now(),
	}
}

// NewJoinRoomMessage creates a join room request message
func NewJoinRoomMessage(username, room string) WebSocketMessage {
	return WebSocketMessage{
		Type:      WSJoinRoom,
		Username:  username,
		Room:      room,
		Timestamp: time.Now(),
	}
}

// NewErrorMessage creates an error message
func NewErrorMessage(content string) WebSocketMessage {
	return WebSocketMessage{
		Type:      WSError,
		Content:   content,
		Timestamp: time.Now(),
	}
}
