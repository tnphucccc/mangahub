package sync

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tnphucccc/mangahub/cmd/cli/internal/config"
	"github.com/tnphucccc/mangahub/pkg/models"
)

func HandleSyncCommand() {
	if len(os.Args) < 3 {
		printSyncUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "monitor":
		syncMonitor()
	default:
		fmt.Printf("Unknown sync subcommand: %s\n", subcommand)
		printSyncUsage()
		os.Exit(1)
	}
}

func printSyncUsage() {
	fmt.Println("Usage: mangahub sync <subcommand>")
	fmt.Println("\nSubcommands:")
	fmt.Println("  monitor              Monitor real-time progress updates (TCP)")
}

func syncMonitor() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if cliConfig.User.Token == "" {
		fmt.Println("Error: Not logged in. Please use 'mangahub auth login' first.")
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", cliConfig.Server.Host, cliConfig.Server.TCPPort)
	fmt.Printf("Connecting to TCP Sync Server at %s...\n", addr)

	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		fmt.Printf("Error: Could not connect to TCP server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 1. Authenticate
	authMsg := models.TCPMessage{
		Type:      models.TCPMessageTypeAuth,
		Timestamp: time.Now(),
		Data: models.TCPAuthMessage{
			Token: cliConfig.User.Token,
		},
	}

	authData, _ := json.Marshal(authMsg)
	authData = append(authData, '\n')
	conn.Write(authData)

	reader := bufio.NewReader(conn)
	
	// Wait for auth response
	line, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Printf("Error reading from server: %v\n", err)
		return
	}

	var resp models.TCPMessage
	json.Unmarshal(line, &resp)

	if resp.Type != models.TCPMessageTypeAuthSuccess {
		fmt.Printf("❌ Authentication failed: %v\n", resp.Data)
		return
	}

	fmt.Println("✅ Connected and authenticated. Monitoring progress updates...")
	fmt.Println("(Press Ctrl+C to stop)")

	// Handle interrupt for graceful exit
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Listen for updates in background
	go func() {
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Printf("\nDisconnected from server: %v\n", err)
				os.Exit(0)
			}

			var msg models.TCPMessage
			if err := json.Unmarshal(line, &msg); err != nil {
				continue
			}

			if msg.Type == models.TCPMessageTypeBroadcast {
				var b models.TCPProgressBroadcast
				dataBytes, _ := json.Marshal(msg.Data)
				json.Unmarshal(dataBytes, &b)

				fmt.Printf("[%s] 📢 %s updated %s to Chapter %d (%s)\n", 
					b.Timestamp.Format("15:04:05"), 
					b.Username, b.MangaTitle, b.CurrentChapter, b.Status)
			}
		}
	}()

	<-interrupt
	fmt.Println("\nStopping monitor...")
}
