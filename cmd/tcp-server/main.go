package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tnphucccc/mangahub/internal/auth"
	"github.com/tnphucccc/mangahub/internal/tcp"
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

	log.Println("Starting TCP Progress Sync Server...")
	log.Printf("Configuration loaded from: %s", configPath)

	// Initialize JWT manager for authentication
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiryDays)

	// Create TCP server
	server := tcp.NewServer(cfg.Server.TCPPort, jwtManager)

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}

	log.Printf("TCP Progress Sync Server started successfully")
	log.Printf("Listening on port: %s", cfg.Server.TCPPort)
	log.Printf("Waiting for client connections...")
	log.Printf("Clients must authenticate with JWT token")

	// Handle graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	log.Println("Shutting down TCP server...")

	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}

	log.Println("TCP server stopped")
}
