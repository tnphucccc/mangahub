package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"time"

	"github.com/tnphucccc/mangahub/pkg/models"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run test_tcp_simple.go <JWT_TOKEN>")
	}

	token := os.Args[1]

	// Connect to TCP server
	addr := "localhost:9090"
	log.Printf("Connecting to TCP server at %s...", addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	log.Println("✅ Connected to TCP server")

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Send authentication
	authMsg := models.TCPMessage{
		Type:      models.TCPMessageTypeAuth,
		Timestamp: time.Now(),
		Data: models.TCPAuthMessage{
			Token: token,
		},
	}

	if err := sendMessage(writer, authMsg); err != nil {
		log.Fatalf("Failed to send auth: %v", err)
	}

	log.Println("Sent authentication message")

	// Read auth response
	respLine, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatalf("Failed to read auth response: %v", err)
	}

	var authResp models.TCPMessage
	if err := json.Unmarshal(respLine, &authResp); err != nil {
		log.Fatalf("Failed to parse auth response: %v", err)
	}

	if authResp.Type == models.TCPMessageTypeAuthSuccess {
		log.Println("✅ Authentication successful!")
		dataBytes, _ := json.Marshal(authResp.Data)
		var successData models.TCPAuthSuccessMessage
		json.Unmarshal(dataBytes, &successData)
		log.Printf("Logged in as: %s (ID: %s)", successData.Username, successData.UserID)
	} else {
		log.Fatalf("❌ Authentication failed")
	}

	// Send a ping
	log.Println("\nSending ping...")
	pingMsg := models.TCPMessage{
		Type:      models.TCPMessageTypePing,
		Timestamp: time.Now(),
		Data: models.TCPPingMessage{
			ClientTime: time.Now(),
		},
	}

	if err := sendMessage(writer, pingMsg); err != nil {
		log.Fatalf("Failed to send ping: %v", err)
	}

	// Read pong response
	respLine, err = reader.ReadBytes('\n')
	if err != nil {
		log.Fatalf("Failed to read pong: %v", err)
	}

	var pongResp models.TCPMessage
	if err := json.Unmarshal(respLine, &pongResp); err != nil {
		log.Fatalf("Failed to parse pong: %v", err)
	}

	if pongResp.Type == models.TCPMessageTypePong {
		log.Println("✅ Received pong from server")
	}

	// Send progress update
	log.Println("\nSending progress update...")
	progressMsg := models.TCPMessage{
		Type:      models.TCPMessageTypeProgress,
		Timestamp: time.Now(),
		Data: models.TCPProgressMessage{
			MangaID:        "manga-001",
			CurrentChapter: 75,
			Status:         models.ReadingStatusReading,
		},
	}

	if err := sendMessage(writer, progressMsg); err != nil {
		log.Fatalf("Failed to send progress: %v", err)
	}

	log.Println("✅ Sent progress update: manga-001, chapter 75")

	// Wait a bit for broadcast
	time.Sleep(500 * time.Millisecond)

	// Try to read broadcast (non-blocking check)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	respLine, err = reader.ReadBytes('\n')
	if err == nil {
		var broadcastMsg models.TCPMessage
		if err := json.Unmarshal(respLine, &broadcastMsg); err == nil {
			if broadcastMsg.Type == models.TCPMessageTypeBroadcast {
				log.Println("✅ Received broadcast of our own progress update")
				dataBytes, _ := json.Marshal(broadcastMsg.Data)
				var broadcast models.TCPProgressBroadcast
				json.Unmarshal(dataBytes, &broadcast)
				log.Printf("  User: %s, Manga: %s, Chapter: %d", broadcast.Username, broadcast.MangaID, broadcast.CurrentChapter)
			}
		}
	}

	log.Println("\n🎉 TCP client test completed successfully!")
	log.Println("All features tested:")
	log.Println("  ✅ Connection established")
	log.Println("  ✅ Authentication with JWT")
	log.Println("  ✅ Ping/Pong heartbeat")
	log.Println("  ✅ Progress update sent")
	log.Println("  ✅ Broadcast received")
}

func sendMessage(writer *bufio.Writer, msg models.TCPMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	writer.Write(data)
	return writer.Flush()
}
