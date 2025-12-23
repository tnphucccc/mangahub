package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpchandler "github.com/tnphucccc/mangahub/internal/grpc"
	"github.com/tnphucccc/mangahub/internal/grpc/pb"
	"github.com/tnphucccc/mangahub/internal/manga"
	"github.com/tnphucccc/mangahub/internal/user"
	"github.com/tnphucccc/mangahub/pkg/config"
	"github.com/tnphucccc/mangahub/pkg/database"
	"github.com/tnphucccc/mangahub/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	configPath := utils.GetEnv("CONFIG_PATH", "./configs/dev.yaml")
	cfg, err := config.LoadFromEnv(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	dbConfig := database.DefaultConfig()
	dbConfig.Path = cfg.Database.Path
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Dependency Injection
	mangaRepo := manga.NewRepository(db)
	userRepo := user.NewRepository(db)
	mangaService := manga.NewService(mangaRepo, userRepo)

	// Initialize gRPC Server
	grpcService := grpchandler.NewServer(mangaService)

	// Start gRPC listener
	addr := fmt.Sprintf(":%s", cfg.Server.GRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to start gRPC listener on %s: %v", addr, err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register gRPC services
	pb.RegisterMangaServiceServer(grpcServer, grpcService)

	log.Printf("gRPC Internal Service listening on %s", addr)

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
		cancel()
	}()

	// Start serving
	if err := grpcServer.Serve(listener); err != nil {
		// GracefulStop() will cause Serve() to return an error, so we don't fatal log it
		// if the context has been canceled.
		if ctx.Err() == nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}
	log.Println("gRPC server stopped.")
}
