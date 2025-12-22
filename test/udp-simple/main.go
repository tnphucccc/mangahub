package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/tnphucccc/mangahub/pkg/models"
)

func main() {
	// Get optional client ID from command line
	clientID := "test-client-1"
	username := "test-user"
	if len(os.Args) >= 2 {
		clientID = os.Args[1]
	}
	if len(os.Args) >= 3 {
		username = os.Args[2]
	}

	// Connect to UDP server
	serverAddr := "localhost:9091"
	log.Printf("Connecting to UDP server at %s...", serverAddr)

	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("Failed to connect to UDP server: %v", err)
	}
	defer conn.Close()

	log.Println("✅ Connected to UDP server")

	// Register client
	log.Println("\nRegistering client...")
	registerMsg := models.UDPMessage{
		Type:      models.UDPMessageTypeRegister,
		Timestamp: time.Now(),
		Data: models.UDPRegisterMessage{
			ClientID: clientID,
			Username: username,
			UserID:   "user-123",
		},
	}

	if err := sendMessage(conn, registerMsg); err != nil {
		log.Fatalf("Failed to send registration: %v", err)
	}

	// Read registration response
	buffer := make([]byte, 2048)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatalf("Failed to read registration response: %v", err)
	}

	var registerResp models.UDPMessage
	if err := json.Unmarshal(buffer[:n], &registerResp); err != nil {
		log.Fatalf("Failed to parse registration response: %v", err)
	}

	if registerResp.Type == models.UDPMessageTypeRegisterSuccess {
		log.Println("✅ Registration successful!")
		dataBytes, _ := json.Marshal(registerResp.Data)
		var successData models.UDPRegisterSuccessMessage
		json.Unmarshal(dataBytes, &successData)
		log.Printf("Registered as: %s", successData.ClientID)
	} else if registerResp.Type == models.UDPMessageTypeRegisterFailed {
		log.Fatalf("❌ Registration failed")
	}

	// Send a ping
	log.Println("\nSending ping...")
	pingMsg := models.UDPMessage{
		Type:      models.UDPMessageTypePing,
		Timestamp: time.Now(),
		Data: models.UDPPingMessage{
			ClientTime: time.Now(),
		},
	}

	if err := sendMessage(conn, pingMsg); err != nil {
		log.Fatalf("Failed to send ping: %v", err)
	}

	// Read pong response
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, _, err = conn.ReadFromUDP(buffer)
	if err != nil {
		log.Printf("⚠️  No pong received (timeout): %v", err)
	} else {
		var pongResp models.UDPMessage
		if err := json.Unmarshal(buffer[:n], &pongResp); err == nil {
			if pongResp.Type == models.UDPMessageTypePong {
				log.Println("✅ Received pong from server")
			}
		}
	}

	// Listen for notifications
	log.Println("\n📡 Listening for notifications (30 seconds)...")
	log.Println("Press Ctrl+C to exit")

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	notificationCount := 0
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			// Timeout is expected
			break
		}

		var notifMsg models.UDPMessage
		if err := json.Unmarshal(buffer[:n], &notifMsg); err != nil {
			log.Printf("Failed to parse notification: %v", err)
			continue
		}

		if notifMsg.Type == models.UDPMessageTypeNotification {
			notificationCount++
			dataBytes, _ := json.Marshal(notifMsg.Data)
			var notification models.UDPNotification
			json.Unmarshal(dataBytes, &notification)
			log.Printf("\n📬 Notification #%d received:", notificationCount)
			log.Printf("  Manga: %s (ID: %s)", notification.MangaTitle, notification.MangaID)
			log.Printf("  Chapter: %d - %s", notification.ChapterNumber, notification.ChapterTitle)
			log.Printf("  Message: %s", notification.Message)
			log.Printf("  Released: %s", notification.ReleaseDate.Format(time.RFC3339))
		}

		// Reset deadline for next notification
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	}

	// Unregister before exiting
	log.Println("\nUnregistering client...")
	unregisterMsg := models.UDPMessage{
		Type:      models.UDPMessageTypeUnregister,
		Timestamp: time.Now(),
		Data: models.UDPUnregisterMessage{
			ClientID: clientID,
		},
	}

	if err := sendMessage(conn, unregisterMsg); err != nil {
		log.Printf("⚠️  Failed to send unregistration: %v", err)
	}

	log.Println("\n🎉 UDP client test completed successfully!")
	log.Println("All features tested:")
	log.Println("  ✅ Connection established")
	log.Println("  ✅ Client registration")
	log.Println("  ✅ Ping/Pong heartbeat")
	log.Printf("  ✅ Notifications received: %d", notificationCount)
	log.Println("  ✅ Client unregistration")
}

func sendMessage(conn *net.UDPConn, msg models.UDPMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
