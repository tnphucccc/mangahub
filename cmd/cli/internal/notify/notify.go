package notify

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/tnphucccc/mangahub/cmd/cli/internal/config"
	"github.com/tnphucccc/mangahub/pkg/models"
)

func HandleNotifyCommand() {
	if len(os.Args) < 3 {
		printNotifyUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "listen":
		notifyListen()
	default:
		fmt.Printf("Unknown notify subcommand: %s\n", subcommand)
		printNotifyUsage()
		os.Exit(1)
	}
}

func printNotifyUsage() {
	fmt.Println("Usage: mangahub notify <subcommand>")
	fmt.Println("\nSubcommands:")
	fmt.Println("  listen               Listen for chapter release notifications (UDP)")
}

func notifyListen() {
	cliConfig, err := config.LoadCLIConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	clientID := uuid.New().String()
	addr := fmt.Sprintf("%s:%d", cliConfig.Server.Host, cliConfig.Server.UDPPort)
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)

	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Register with UDP server
	regMsg := models.UDPMessage{
		Type:      models.UDPMessageTypeRegister,
		Timestamp: time.Now(),
		Data: models.UDPRegisterMessage{
			ClientID: clientID,
			Username: cliConfig.User.Username,
		},
	}
	data, _ := json.Marshal(regMsg)
	conn.WriteToUDP(data, udpAddr)

	fmt.Printf("Listening for UDP notifications from %s...\n", addr)
	fmt.Println("(Press Ctrl+C to stop)")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		buf := make([]byte, 2048)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				return
			}

			var msg models.UDPMessage
			if err := json.Unmarshal(buf[:n], &msg); err != nil {
				continue
			}

			if msg.Type == models.UDPMessageTypeNotification {
				var n models.UDPNotification
				dataBytes, _ := json.Marshal(msg.Data)
				json.Unmarshal(dataBytes, &n)

				fmt.Printf("[%s] 🔔 NEW CHAPTER: %s - Chapter %d released!\n", 
					time.Now().Format("15:04:05"), n.MangaTitle, n.ChapterNumber)
				fmt.Printf("      Message: %s\n", n.Message)
			}
		}
	}()

	<-	interrupt
	// Unregister
	unregMsg := models.UDPMessage{
		Type:      models.UDPMessageTypeUnregister,
		Timestamp: time.Now(),
		Data: models.UDPUnregisterMessage{
			ClientID: clientID,
		},
	}
	data, _ = json.Marshal(unregMsg)
	conn.WriteToUDP(data, udpAddr)
	fmt.Println("\nStopping listener...")
}
