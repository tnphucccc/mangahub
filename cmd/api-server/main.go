package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tnphucccc/mangahub/internal/auth"
	"github.com/tnphucccc/mangahub/internal/manga"
	"github.com/tnphucccc/mangahub/internal/middleware"
	"github.com/tnphucccc/mangahub/internal/user"
	"github.com/tnphucccc/mangahub/pkg/config"
	"github.com/tnphucccc/mangahub/pkg/database"
)

func main() {
	// Load configuration
	configPath := getEnv("CONFIG_PATH", "./configs/dev.yaml")
	cfg, err := config.LoadFromEnv(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	dbConfig := database.Config{
		Path:            cfg.Database.Path,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 0,
	}
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	log.Println("Connected to database successfully")

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpiryDays)

	// Initialize repositories
	userRepo := user.NewRepository(db)
	mangaRepo := manga.NewRepository(db)

	// Initialize services
	userService := user.NewService(userRepo, jwtManager)
	mangaService := manga.NewService(mangaRepo)

	// Initialize handlers
	userHandler := user.NewHandler(userService)
	mangaHandler := manga.NewHandler(mangaService)

	// Initialize Gin router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		// Check database health
		if err := database.HealthCheck(db); err != nil {
			c.JSON(500, gin.H{
				"status":  "unhealthy",
				"service": "MangaHub HTTP API Server",
				"error":   "database connection failed",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "MangaHub HTTP API Server",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	api := router.Group("/api/v1")
	{
		// Public auth routes
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", userHandler.Register)
			authRoutes.POST("/login", userHandler.Login)
		}

		// Public manga routes
		mangaRoutes := api.Group("/manga")
		{
			mangaRoutes.GET("", mangaHandler.Search)          // Search manga
			mangaRoutes.GET("/all", mangaHandler.GetAll)      // Get all manga
			mangaRoutes.GET("/:id", mangaHandler.GetByID)     // Get manga by ID
		}

		// Protected user routes (require authentication)
		userRoutes := api.Group("/users")
		userRoutes.Use(middleware.AuthMiddleware(userService))
		{
			userRoutes.GET("/me", userHandler.GetProfile)                          // Get current user profile
			userRoutes.GET("/library", mangaHandler.GetLibrary)                    // Get user's library
			userRoutes.POST("/library", mangaHandler.AddToLibrary)                 // Add manga to library
			userRoutes.GET("/progress/:manga_id", mangaHandler.GetProgress)        // Get progress for manga
			userRoutes.PUT("/progress/:manga_id", mangaHandler.UpdateProgress)     // Update reading progress
		}
	}

	// Start server
	addr := cfg.GetHTTPAddress()
	log.Printf("HTTP API Server starting on %s", addr)
	log.Printf("Database: %s", cfg.Database.Path)
	log.Printf("Endpoints available:")
	log.Printf("  - Health check: GET /health")
	log.Printf("  - Register: POST /api/v1/auth/register")
	log.Printf("  - Login: POST /api/v1/auth/login")
	log.Printf("  - Search manga: GET /api/v1/manga?q=query")
	log.Printf("  - Get manga: GET /api/v1/manga/:id")
	log.Printf("  - User library: GET /api/v1/users/library (protected)")
	log.Printf("  - Add to library: POST /api/v1/users/library (protected)")
	log.Printf("  - Update progress: PUT /api/v1/users/progress/:manga_id (protected)")

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
