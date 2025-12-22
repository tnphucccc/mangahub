package udp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/tnphucccc/mangahub/pkg/models"
)

// RegisteredClient represents a registered UDP client
type RegisteredClient struct {
	ClientID string
	UserID   string
	Username string
	Address  *net.UDPAddr
	LastSeen time.Time
}

// Server represents the UDP notification server
type Server struct {
	Port       string
	conn       *net.UDPConn
	clients    map[string]*RegisteredClient // clientID -> RegisteredClient
	mu         sync.RWMutex
	notify     chan models.UDPNotification
	shutdown   chan struct{}
	bufferSize int // Size of UDP receive buffer
}

// NewServer creates a new UDP notification server instance
func NewServer(port string) *Server {
	return &Server{
		Port:       port,
		clients:    make(map[string]*RegisteredClient),
		notify:     make(chan models.UDPNotification, 100),
		shutdown:   make(chan struct{}),
		bufferSize: 2048, // 2KB buffer for UDP packets
	}
}

// Start starts the UDP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("failed to start UDP listener: %w", err)
	}

	s.conn = conn
	log.Printf("UDP Notification Server listening on %s", addr)

	// Start notification broadcast goroutine
	go s.notificationLoop()

	// Start client cleanup goroutine (remove stale clients)
	go s.cleanupLoop()

	// Start listening for messages
	go s.listen()

	return nil
}

// listen listens for incoming UDP messages
func (s *Server) listen() {
	buffer := make([]byte, s.bufferSize)

	for {
		select {
		case <-s.shutdown:
			return
		default:
			n, clientAddr, err := s.conn.ReadFromUDP(buffer)
			if err != nil {
				select {
				case <-s.shutdown:
					return
				default:
					log.Printf("Error reading UDP message: %v", err)
					continue
				}
			}

			// Handle message in goroutine to avoid blocking
			go s.handleMessage(clientAddr, buffer[:n])
		}
	}
}

// handleMessage handles incoming UDP messages
func (s *Server) handleMessage(addr *net.UDPAddr, data []byte) {
	var msg models.UDPMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("Invalid message from %s: %v", addr.String(), err)
		s.sendError(addr, "INVALID_MESSAGE", "Invalid message format")
		return
	}

	// Handle message based on type
	switch msg.Type {
	case models.UDPMessageTypeRegister:
		s.handleRegister(addr, msg)
	case models.UDPMessageTypeUnregister:
		s.handleUnregister(addr, msg)
	case models.UDPMessageTypePing:
		s.handlePing(addr, msg)
	default:
		s.sendError(addr, "UNKNOWN_MESSAGE_TYPE", fmt.Sprintf("Unknown message type: %s", msg.Type))
	}
}

// handleRegister handles client registration
func (s *Server) handleRegister(addr *net.UDPAddr, msg models.UDPMessage) {
	// Parse registration data
	registerDataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		s.sendRegisterFailed(addr, "Invalid registration data")
		return
	}

	var registerData models.UDPRegisterMessage
	if err := json.Unmarshal(registerDataBytes, &registerData); err != nil {
		s.sendRegisterFailed(addr, "Invalid registration data")
		return
	}

	// Validate client ID
	if registerData.ClientID == "" {
		s.sendRegisterFailed(addr, "Client ID is required")
		return
	}

	// Register client
	client := &RegisteredClient{
		ClientID: registerData.ClientID,
		UserID:   registerData.UserID,
		Username: registerData.Username,
		Address:  addr,
		LastSeen: time.Now(),
	}

	s.registerClient(client)

	// Send success response
	s.sendRegisterSuccess(addr, registerData.ClientID)

	log.Printf("Client registered: %s (user: %s) from %s", registerData.ClientID, registerData.Username, addr.String())
}

// handleUnregister handles client unregistration
func (s *Server) handleUnregister(addr *net.UDPAddr, msg models.UDPMessage) {
	// Parse unregistration data
	unregisterDataBytes, err := json.Marshal(msg.Data)
	if err != nil {
		s.sendError(addr, "INVALID_DATA", "Invalid unregistration data")
		return
	}

	var unregisterData models.UDPUnregisterMessage
	if err := json.Unmarshal(unregisterDataBytes, &unregisterData); err != nil {
		s.sendError(addr, "INVALID_DATA", "Invalid unregistration data")
		return
	}

	// Unregister client
	s.unregisterClient(unregisterData.ClientID)

	log.Printf("Client unregistered: %s from %s", unregisterData.ClientID, addr.String())
}

// handlePing responds to ping messages
func (s *Server) handlePing(addr *net.UDPAddr, msg models.UDPMessage) {
	// Update last seen time for this client if registered
	s.mu.Lock()
	for _, client := range s.clients {
		if client.Address.String() == addr.String() {
			client.LastSeen = time.Now()
			break
		}
	}
	s.mu.Unlock()

	// Send pong response
	pong := models.UDPMessage{
		Type:      models.UDPMessageTypePong,
		Timestamp: time.Now(),
		Data: models.UDPPongMessage{
			ServerTime: time.Now(),
		},
	}

	s.sendMessage(addr, pong)
}

// notificationLoop listens for notifications and broadcasts them
func (s *Server) notificationLoop() {
	for {
		select {
		case notification := <-s.notify:
			s.broadcastNotification(notification)
		case <-s.shutdown:
			return
		}
	}
}

// cleanupLoop removes stale clients (no ping for 5 minutes)
func (s *Server) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.cleanupStaleClients()
		case <-s.shutdown:
			return
		}
	}
}

// cleanupStaleClients removes clients that haven't been seen for 5 minutes
func (s *Server) cleanupStaleClients() {
	s.mu.Lock()
	defer s.mu.Unlock()

	staleThreshold := time.Now().Add(-5 * time.Minute)
	staleCIDs := []string{}

	for clientID, client := range s.clients {
		if client.LastSeen.Before(staleThreshold) {
			staleCIDs = append(staleCIDs, clientID)
		}
	}

	for _, clientID := range staleCIDs {
		delete(s.clients, clientID)
		log.Printf("Removed stale client: %s", clientID)
	}

	if len(staleCIDs) > 0 {
		log.Printf("Cleanup: removed %d stale clients (total: %d)", len(staleCIDs), len(s.clients))
	}
}

// broadcastNotification sends a notification to all registered clients
func (s *Server) broadcastNotification(notification models.UDPNotification) {
	s.mu.RLock()
	clients := make([]*RegisteredClient, 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, client)
	}
	s.mu.RUnlock()

	msg := models.UDPMessage{
		Type:      models.UDPMessageTypeNotification,
		Timestamp: time.Now(),
		Data:      notification,
	}

	// Send to all clients in parallel (UDP is fire-and-forget)
	successCount := 0
	for _, client := range clients {
		if err := s.sendMessage(client.Address, msg); err == nil {
			successCount++
		}
	}

	log.Printf("Broadcasted notification to %d/%d clients: %s - Chapter %d",
		successCount, len(clients), notification.MangaTitle, notification.ChapterNumber)
}

// BroadcastNotification queues a notification for broadcast (called from HTTP API or other services)
func (s *Server) BroadcastNotification(notification models.UDPNotification) {
	select {
	case s.notify <- notification:
		// Successfully queued
	default:
		log.Printf("Warning: Notification channel full, dropping notification for %s", notification.MangaTitle)
	}
}

// registerClient registers a new client
func (s *Server) registerClient(client *RegisteredClient) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.clients[client.ClientID] = client

	log.Printf("Client registered: %s (total: %d)", client.ClientID, len(s.clients))
}

// unregisterClient removes a client
func (s *Server) unregisterClient(clientID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, clientID)

	log.Printf("Client unregistered: %s (total: %d)", clientID, len(s.clients))
}

// sendMessage sends a UDP message to a specific address
func (s *Server) sendMessage(addr *net.UDPAddr, msg models.UDPMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = s.conn.WriteToUDP(data, addr)
	if err != nil {
		return fmt.Errorf("failed to send UDP message: %w", err)
	}

	return nil
}

// sendRegisterSuccess sends registration success message
func (s *Server) sendRegisterSuccess(addr *net.UDPAddr, clientID string) {
	msg := models.UDPMessage{
		Type:      models.UDPMessageTypeRegisterSuccess,
		Timestamp: time.Now(),
		Data: models.UDPRegisterSuccessMessage{
			ClientID: clientID,
			Message:  "Registration successful",
		},
	}
	s.sendMessage(addr, msg)
}

// sendRegisterFailed sends registration failed message
func (s *Server) sendRegisterFailed(addr *net.UDPAddr, reason string) {
	msg := models.UDPMessage{
		Type:      models.UDPMessageTypeRegisterFailed,
		Timestamp: time.Now(),
		Data: models.UDPRegisterFailedMessage{
			Reason: reason,
		},
	}
	s.sendMessage(addr, msg)
}

// sendError sends an error message
func (s *Server) sendError(addr *net.UDPAddr, code, message string) {
	msg := models.UDPMessage{
		Type:      models.UDPMessageTypeError,
		Timestamp: time.Now(),
		Data: models.UDPErrorMessage{
			Code:    code,
			Message: message,
		},
	}
	s.sendMessage(addr, msg)
}

// GetStats returns server statistics
func (s *Server) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"total_clients":      len(s.clients),
		"notification_queue": len(s.notify),
		"buffer_size":        s.bufferSize,
	}
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	close(s.shutdown)

	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}
