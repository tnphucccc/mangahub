package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tnphucccc/mangahub/internal/udp"
	"github.com/tnphucccc/mangahub/pkg/config"
	"github.com/tnphucccc/mangahub/pkg/utils"
)

func main() {
	// Load configuration
	configPath := utils.GetEnv("CONFIG_PATH", "./configs/dev.yaml")
	cfg, err := config.LoadFromEnv(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Starting UDP Notification Server...")
	log.Printf("Configuration loaded from: %s", configPath)

	// Create UDP notification server
	server := udp.NewServer(cfg.Server.UDPPort)

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start UDP server: %v", err)
	}

	log.Printf("UDP Notification Server started successfully")
	log.Printf("Listening on port: %s", cfg.Server.UDPPort)
	log.Printf("Waiting for client registrations...")
	log.Printf("Ready to broadcast chapter release notifications")

	// Handle graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal
	<-shutdown
	log.Println("Shutting down UDP server...")

	// Stop server
	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}

	log.Println("UDP server stopped")
}
