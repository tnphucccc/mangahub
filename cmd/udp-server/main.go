package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/tnphucccc/mangahub/pkg/utils"
)

func main() {
	// Load configuration
	port := utils.GetEnv("UDP_PORT", "9091")

	// Start UDP listener
	addr := fmt.Sprintf(":%s", port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Failed to start UDP listener: %v", err)
	}
	defer conn.Close()

	log.Printf("UDP Notification Server listening on %s", addr)

	// Handle graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Println("Shutting down UDP server...")
		conn.Close()
		os.Exit(0)
	}()

	// Buffer for incoming messages
	buffer := make([]byte, 1024)

	// Listen for messages
	for {
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP message: %v", err)
			continue
		}

		// Handle message in goroutine
		go handleMessage(conn, clientAddr, buffer[:n])
	}
}

func handleMessage(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
	log.Printf("Received message from %s: %s", addr.String(), string(data))

	// TODO: Implement notification protocol
	// - Register client for notifications
	// - Broadcast chapter release notifications
	// - Handle client unregistration
}
