package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/tnphucccc/mangahub/internal/auth"
	"github.com/tnphucccc/mangahub/pkg/models"
)

// Client represents a connected TCP client
type Client struct {
	ID       string
	UserID   string
	Username string
	Conn     net.Conn
	Writer   *bufio.Writer
	mu       sync.Mutex
}

// Server represents the TCP progress sync server
type Server struct {
	Port       string
	listener   net.Listener
	clients    map[string]*Client // clientID -> Client
	userIndex  map[string][]*Client // userID -> []*Client (multiple devices per user)
	mu         sync.RWMutex
	broadcast  chan models.TCPProgressBroadcast
	jwtManager *auth.JWTManager
	shutdown   chan struct{}
}

// NewServer creates a new TCP server instance
func NewServer(port string, jwtManager *auth.JWTManager) *Server {
	return &Server{
		Port:       port,
		clients:    make(map[string]*Client),
		userIndex:  make(map[string][]*Client),
		broadcast:  make(chan models.TCPProgressBroadcast, 100),
		jwtManager: jwtManager,
		shutdown:   make(chan struct{}),
	}
}

// Start starts the TCP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}

	s.listener = listener
	log.Printf("TCP Progress Sync Server listening on %s", addr)

	// Start broadcast goroutine
	go s.broadcastLoop()

	// Accept connections
	go s.acceptConnections()

	return nil
}

// acceptConnections accepts incoming TCP connections
func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.shutdown:
				return
			default:
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
		}

		// Handle connection in goroutine
		go s.handleConnection(conn)
	}
}

// handleConnection handles a single client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientID := conn.RemoteAddr().String()
	log.Printf("New connection from %s", clientID)

	client := &Client{
		ID:     clientID,
		Conn:   conn,
		Writer: bufio.NewWriter(conn),
	}

	reader := bufio.NewReader(conn)

	// Wait for authentication (30 second timeout)
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	// Read authentication message
	line, err := reader.ReadBytes('\n')
	if err != nil {
		log.Printf("Failed to read auth message from %s: %v", clientID, err)
		s.sendError(client, "AUTH_TIMEOUT", "Authentication timeout")
		return
	}

	var msg models.TCPMessage
	if err := json.Unmarshal(line, &msg); err != nil {
		log.Printf("Failed to parse message from %s: %v", clientID, err)
		s.sendError(client, "INVALID_MESSAGE", "Invalid message format")
		return
	}

	// Verify it's an auth message
	if msg.Type != models.TCPMessageTypeAuth {
		s.sendAuthFailed(client, "First message must be authentication")
		return
	}

	// Parse auth data
	authDataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		s.sendAuthFailed(client, "Invalid authentication data")
		return
	}

	var authData models.TCPAuthMessage
	if err := json.Unmarshal(authDataBytes, &authData); err != nil {
		s.sendAuthFailed(client, "Invalid authentication data")
		return
	}

	// Validate JWT token
	claims, err := s.jwtManager.ValidateToken(authData.Token)
	if err != nil {
		s.sendAuthFailed(client, "Invalid or expired token")
		return
	}

	// Authentication successful
	client.UserID = claims.UserID
	client.Username = claims.Username

	// Register client
	s.registerClient(client)
	defer s.unregisterClient(client)

	// Send auth success
	s.sendAuthSuccess(client)

	// Remove read deadline for regular messages
	conn.SetReadDeadline(time.Time{})

	log.Printf("Client %s authenticated as user %s (%s)", clientID, claims.Username, claims.UserID)

	// Read messages loop
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Client %s disconnected: %v", clientID, err)
			return
		}

		var msg models.TCPMessage
		if err := json.Unmarshal(line, &msg); err != nil {
			log.Printf("Invalid message from %s: %v", clientID, err)
			s.sendError(client, "INVALID_MESSAGE", "Invalid message format")
			continue
		}

		// Handle message based on type
		s.handleMessage(client, msg)
	}
}

// handleMessage handles different message types
func (s *Server) handleMessage(client *Client, msg models.TCPMessage) {
	switch msg.Type {
	case models.TCPMessageTypePing:
		s.handlePing(client, msg)
	case models.TCPMessageTypeProgress:
		s.handleProgressUpdate(client, msg)
	default:
		s.sendError(client, "UNKNOWN_MESSAGE_TYPE", fmt.Sprintf("Unknown message type: %s", msg.Type))
	}
}

// handlePing responds to ping messages
func (s *Server) handlePing(client *Client, msg models.TCPMessage) {
	pong := models.TCPMessage{
		Type:      models.TCPMessageTypePong,
		Timestamp: time.Now(),
		Data: models.TCPPongMessage{
			ServerTime: time.Now(),
		},
	}

	s.sendMessage(client, pong)
}

// handleProgressUpdate handles progress update messages
func (s *Server) handleProgressUpdate(client *Client, msg models.TCPMessage) {
	// Parse progress data
	progressDataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		s.sendError(client, "INVALID_DATA", "Invalid progress data")
		return
	}

	var progressData models.TCPProgressMessage
	if err := json.Unmarshal(progressDataBytes, &progressData); err != nil {
		s.sendError(client, "INVALID_DATA", "Invalid progress data")
		return
	}

	// Create broadcast message
	broadcast := models.TCPProgressBroadcast{
		UserID:         client.UserID,
		Username:       client.Username,
		MangaID:        progressData.MangaID,
		CurrentChapter: progressData.CurrentChapter,
		Status:         progressData.Status,
		Timestamp:      time.Now(),
	}

	// Send to broadcast channel
	s.broadcast <- broadcast

	log.Printf("Progress update from %s: manga=%s, chapter=%d", client.Username, progressData.MangaID, progressData.CurrentChapter)
}

// broadcastLoop listens for broadcast messages and sends to all clients
func (s *Server) broadcastLoop() {
	for {
		select {
		case broadcast := <-s.broadcast:
			s.broadcastToAll(broadcast)
		case <-s.shutdown:
			return
		}
	}
}

// broadcastToAll sends a message to all connected clients
func (s *Server) broadcastToAll(broadcast models.TCPProgressBroadcast) {
	s.mu.RLock()
	clients := make([]*Client, 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, client)
	}
	s.mu.RUnlock()

	msg := models.TCPMessage{
		Type:      models.TCPMessageTypeBroadcast,
		Timestamp: time.Now(),
		Data:      broadcast,
	}

	// Send to all clients in parallel (non-blocking)
	for _, client := range clients {
		go s.sendMessage(client, msg)
	}

	log.Printf("Broadcasted progress update to %d clients", len(clients))
}

// BroadcastProgress sends a progress update to all connected clients (called from HTTP API)
func (s *Server) BroadcastProgress(broadcast models.TCPProgressBroadcast) {
	select {
	case s.broadcast <- broadcast:
		// Successfully queued
	default:
		log.Printf("Warning: Broadcast channel full, dropping message")
	}
}

// registerClient registers a new client
func (s *Server) registerClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[client.ID] = client

	// Add to user index
	s.userIndex[client.UserID] = append(s.userIndex[client.UserID], client)

	log.Printf("Client registered: %s (total: %d)", client.ID, len(s.clients))
}

// unregisterClient removes a client
func (s *Server) unregisterClient(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, client.ID)

	// Remove from user index
	if clients, ok := s.userIndex[client.UserID]; ok {
		for i, c := range clients {
			if c.ID == client.ID {
				s.userIndex[client.UserID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		if len(s.userIndex[client.UserID]) == 0 {
			delete(s.userIndex, client.UserID)
		}
	}

	log.Printf("Client unregistered: %s (total: %d)", client.ID, len(s.clients))
}

// sendMessage sends a message to a client
func (s *Server) sendMessage(client *Client, msg models.TCPMessage) error {
	client.mu.Lock()
	defer client.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	data = append(data, '\n')

	if _, err := client.Writer.Write(data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := client.Writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

// sendAuthSuccess sends authentication success message
func (s *Server) sendAuthSuccess(client *Client) {
	msg := models.TCPMessage{
		Type:      models.TCPMessageTypeAuthSuccess,
		Timestamp: time.Now(),
		Data: models.TCPAuthSuccessMessage{
			UserID:   client.UserID,
			Username: client.Username,
			Message:  "Authentication successful",
		},
	}
	s.sendMessage(client, msg)
}

// sendAuthFailed sends authentication failed message
func (s *Server) sendAuthFailed(client *Client, reason string) {
	msg := models.TCPMessage{
		Type:      models.TCPMessageTypeAuthFailed,
		Timestamp: time.Now(),
		Data: models.TCPAuthFailedMessage{
			Reason: reason,
		},
	}
	s.sendMessage(client, msg)
}

// sendError sends an error message to client
func (s *Server) sendError(client *Client, code, message string) {
	msg := models.TCPMessage{
		Type:      models.TCPMessageTypeError,
		Timestamp: time.Now(),
		Data: models.TCPErrorMessage{
			Code:    code,
			Message: message,
		},
	}
	s.sendMessage(client, msg)
}

// GetStats returns server statistics
func (s *Server) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"total_clients":      len(s.clients),
		"total_users":        len(s.userIndex),
		"broadcast_queue":    len(s.broadcast),
	}
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	close(s.shutdown)

	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}
