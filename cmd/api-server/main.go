package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	port := getEnv("HTTP_PORT", "8080")

	// Initialize Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "MangaHub HTTP API Server",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// TODO: Add route groups here
		// auth := api.Group("/auth")
		// manga := api.Group("/manga")
		// users := api.Group("/users")
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("HTTP API Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
