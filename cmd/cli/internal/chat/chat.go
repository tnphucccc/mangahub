package chat

import (
	"bufio"
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
	fmt.Println("Usage: mangahub chat <subcommand>")
	fmt.Println("\nSubcommands:")
	fmt.Println("  join                 Join the real-time chat")
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

	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", cliConfig.Server.Host, cliConfig.Server.WebSocketPort), Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Connected to chat. Type your message and press Enter to send. Press Ctrl+C to exit.")

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
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
					err = c.WriteMessage(websocket.TextMessage, []byte(line))
					if err != nil {
						log.Println("write:", err)
						return
					}
				}
				fmt.Print("> ")
			case err := <-errChan:
				log.Println("read string:", err)
				return
			case <-done:
				return
			case <-interrupt:
				// Second interrupt, force exit
				log.Println("second interrupt, forcing exit")
				return
			}
		}
	}
}
