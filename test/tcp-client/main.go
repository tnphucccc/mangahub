package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/tnphucccc/mangahub/pkg/models"
)

func main() {
	// Command line flags
	host := flag.String("host", "localhost", "TCP server host")
	port := flag.String("port", "9090", "TCP server port")
	token := flag.String("token", "", "JWT authentication token (required)")
	flag.Parse()

	if *token == "" {
		log.Fatal("JWT token is required. Use -token flag")
	}

	// Connect to TCP server
	addr := fmt.Sprintf("%s:%s", *host, *port)
	log.Printf("Connecting to TCP server at %s...", addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to TCP server")

	// Create reader and writer
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Send authentication message
	authMsg := models.TCPMessage{
		Type:      models.TCPMessageTypeAuth,
		Timestamp: time.Now(),
		Data: models.TCPAuthMessage{
			Token: *token,
		},
	}

	if err := sendMessage(writer, authMsg); err != nil {
		log.Fatalf("Failed to send auth message: %v", err)
	}

	log.Println("Sent authentication message")

	// Wait for auth response
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
	} else if authResp.Type == models.TCPMessageTypeAuthFailed {
		log.Fatalf("❌ Authentication failed")
	}

	// Start listening for broadcasts in background
	go listenForMessages(reader)

	// Interactive command loop
	fmt.Println("\n=== TCP Client Connected ===")
	fmt.Println("Commands:")
	fmt.Println("  progress <manga_id> <chapter> - Send progress update")
	fmt.Println("  ping                            - Send ping")
	fmt.Println("  quit                            - Exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse command
		var cmd, arg1, arg2 string
		fmt.Sscanf(line, "%s %s %s", &cmd, &arg1, &arg2)

		switch cmd {
		case "quit", "exit":
			log.Println("Exiting...")
			return

		case "ping":
			pingMsg := models.TCPMessage{
				Type:      models.TCPMessageTypePing,
				Timestamp: time.Now(),
				Data: models.TCPPingMessage{
					ClientTime: time.Now(),
				},
			}
			if err := sendMessage(writer, pingMsg); err != nil {
				log.Printf("Failed to send ping: %v", err)
			} else {
				log.Println("Sent ping")
			}

		case "progress":
			if arg1 == "" || arg2 == "" {
				fmt.Println("Usage: progress <manga_id> <chapter>")
				continue
			}

			var chapter int
			fmt.Sscanf(arg2, "%d", &chapter)

			progressMsg := models.TCPMessage{
				Type:      models.TCPMessageTypeProgress,
				Timestamp: time.Now(),
				Data: models.TCPProgressMessage{
					MangaID:        arg1,
					CurrentChapter: chapter,
					Status:         models.ReadingStatusReading,
				},
			}

			if err := sendMessage(writer, progressMsg); err != nil {
				log.Printf("Failed to send progress: %v", err)
			} else {
				log.Printf("Sent progress update: manga=%s, chapter=%d", arg1, chapter)
			}

		default:
			fmt.Println("Unknown command. Type 'quit' to exit.")
		}
	}
}

func sendMessage(writer *bufio.Writer, msg models.TCPMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	data = append(data, '\n')

	if _, err := writer.Write(data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func listenForMessages(reader *bufio.Reader) {
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			log.Printf("Connection closed: %v", err)
			os.Exit(0)
		}

		var msg models.TCPMessage
		if err := json.Unmarshal(line, &msg); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		switch msg.Type {
		case models.TCPMessageTypePong:
			fmt.Println("\n🏓 Received pong from server")
			fmt.Print("> ")

		case models.TCPMessageTypeBroadcast:
			dataBytes, _ := json.Marshal(msg.Data)
			var broadcast models.TCPProgressBroadcast
			json.Unmarshal(dataBytes, &broadcast)
			fmt.Printf("\n📢 Progress Update: %s read %s chapter %d (%s) at %s\n",
				broadcast.Username,
				broadcast.MangaTitle,
				broadcast.CurrentChapter,
				broadcast.Status,
				broadcast.Timestamp.Format("15:04:05"))
			fmt.Print("> ")

		case models.TCPMessageTypeError:
			dataBytes, _ := json.Marshal(msg.Data)
			var errMsg models.TCPErrorMessage
			json.Unmarshal(dataBytes, &errMsg)
			fmt.Printf("\n❌ Error: %s - %s\n", errMsg.Code, errMsg.Message)
			fmt.Print("> ")

		default:
			fmt.Printf("\nReceived message: type=%s\n", msg.Type)
			fmt.Print("> ")
		}
	}
}
