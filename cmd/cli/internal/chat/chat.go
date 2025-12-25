package chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/config"
	"github.com/tnphucccc/mangahub/pkg/models"
)

func HandleChatCommand() {
	if len(os.Args) < 3 {
		printChatUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "join":
		chatJoin()
	default:
		fmt.Printf("Unknown chat subcommand: %s\n", subcommand)
		printChatUsage()
		os.Exit(1)
	}
}

func printChatUsage() {
	fmt.Println("Usage: mangahub chat <subcommand> [flags]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  join                 Join the real-time chat")
	fmt.Println("\nFlags:")
	fmt.Println("  --manga-id <id>      Join manga-specific chat room (default: general)")
}

func chatJoin() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cliConfig.User.Token == "" {
		fmt.Println("Error: Not logged in. Please use 'mangahub auth login' first.")
		os.Exit(1)
	}

	// Parse flags for manga-id
	var mangaID string
	for i := 3; i < len(os.Args); i++ {
		if os.Args[i] == "--manga-id" && i+1 < len(os.Args) {
			mangaID = os.Args[i+1]
			break
		}
	}

	// Determine room name
	room := "general"
	if mangaID != "" {
		room = mangaID
	}

	// Build WebSocket URL with username and room as query parameters
	u := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", cliConfig.Server.Host, cliConfig.Server.WebSocketPort),
		Path:   "/ws",
	}
	q := u.Query()
	q.Set("username", cliConfig.User.Username)
	q.Set("room", room)
	u.RawQuery = q.Encode()

	fmt.Printf("Connecting to WebSocket chat server at %s...\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Goroutine to receive messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("Connection closed: %v", err)
				}
				return
			}

			// Parse the incoming message
			var wsMsg models.WebSocketMessage
			if err := json.Unmarshal(message, &wsMsg); err != nil {
				fmt.Printf("Error parsing message: %v\n", err)
				continue
			}

			// Display the message with proper formatting
			displayMessage(wsMsg)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	reader := bufio.NewReader(os.Stdin)

	// Display welcome message
	fmt.Printf("\n✓ Connected to Chat Room: #%s\n", room)
	fmt.Println("───────────────────────────────────────────────────────────")
	fmt.Println("\nYou are now in chat. Type your message and press Enter.")
	fmt.Println("Press Ctrl+C to leave.")
	fmt.Print("> ")

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			fmt.Println("\nLeaving chat...")

			// Cleanly close the connection
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error closing connection:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			fmt.Println("✓ Disconnected from chat server")
			return
		default:
			// Non-blocking read from stdin
			lineChan := make(chan string)
			errChan := make(chan error)
			go func() {
				line, err := reader.ReadString('\n')
				if err != nil {
					errChan <- err
					return
				}
				lineChan <- line
			}()

			select {
			case line := <-lineChan:
				line = strings.TrimSpace(line)
				if line != "" {
					// Create a chat message
					chatMsg := models.NewChatMessage(cliConfig.User.Username, room, line)
					msgBytes, err := json.Marshal(chatMsg)
					if err != nil {
						fmt.Printf("Error creating message: %v\n", err)
						continue
					}

					err = c.WriteMessage(websocket.TextMessage, msgBytes)
					if err != nil {
						log.Println("Error sending message:", err)
						return
					}
				}
				fmt.Print("> ")
			case err := <-errChan:
				log.Println("Error reading input:", err)
				return
			case <-done:
				return
			case <-interrupt:
				// Second interrupt, force exit
				fmt.Println("\nForcing exit...")
				return
			}
		}
	}
}

// displayMessage formats and displays a WebSocket message
func displayMessage(msg models.WebSocketMessage) {
	timestamp := msg.Timestamp.Format("15:04")

	switch msg.Type {
	case models.WSChatMessage:
		// Display chat message with [username] prefix
		fmt.Printf("\r[%s] %s: %s\n> ", timestamp, msg.Username, msg.Content)
	case models.WSSystemMessage:
		// Display system message
		fmt.Printf("\r[%s] * %s\n> ", timestamp, msg.Content)
	case models.WSError:
		// Display error message
		fmt.Printf("\r✗ Error: %s\n> ", msg.Content)
	default:
		// Display unknown message type
		fmt.Printf("\r[%s] %s\n> ", timestamp, msg.Content)
	}
}
