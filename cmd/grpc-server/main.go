package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	port := getEnv("GRPC_PORT", "9092")

	// Start gRPC listener
	addr := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to start gRPC listener: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// TODO: Register gRPC services
	// pb.RegisterMangaServiceServer(grpcServer, &mangaService{})
	// pb.RegisterProgressServiceServer(grpcServer, &progressService{})

	log.Printf("gRPC Internal Service listening on %s", addr)

	// Handle graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutdown
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	// Start serving
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
