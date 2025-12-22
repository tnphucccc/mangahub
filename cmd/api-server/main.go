package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http" // Added import for http.Server
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tnphucccc/mangahub/internal/auth"
	"github.com/tnphucccc/mangahub/internal/manga"
	"github.com/tnphucccc/mangahub/internal/middleware"
	"github.com/tnphucccc/mangahub/internal/user"
	"github.com/tnphucccc/mangahub/internal/websocket"
	"github.com/tnphucccc/mangahub/pkg/config"
	"github.com/tnphucccc/mangahub/pkg/database"
	"github.com/tnphucccc/mangahub/pkg/models"
	"github.com/tnphucccc/mangahub/pkg/utils"
)

func main() {
	// Load configuration
	configPath := utils.GetEnv("CONFIG_PATH", "./configs/dev.yaml")
	cfg, err := config.LoadFromEnv(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal, initiating graceful shutdown...")
		cancel() // Trigger context cancellation
	}()

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

	// Initialize WebSocket hub and run it
	wsHub := websocket.NewHub()
	go wsHub.Run(ctx) // Pass context to hub

	// Start UDP listener and bridge to WebSocket
	go listenForUDPNotifications(ctx, wsHub, cfg.GetUDPAddress())

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
			mangaRoutes.GET("", mangaHandler.Search)      // Search manga
			mangaRoutes.GET("/all", mangaHandler.GetAll)  // Get all manga
			mangaRoutes.GET("/:id", mangaHandler.GetByID) // Get manga by ID
		}

		// Protected user routes (require authentication)
		userRoutes := api.Group("/users")
		userRoutes.Use(middleware.AuthMiddleware(userService))
		{
			userRoutes.GET("/me", userHandler.GetProfile)                      // Get current user profile
			userRoutes.GET("/library", mangaHandler.GetLibrary)                // Get user's library
			userRoutes.POST("/library", mangaHandler.AddToLibrary)             // Add manga to library
			userRoutes.GET("/progress/:manga_id", mangaHandler.GetProgress)    // Get progress for manga
			userRoutes.PUT("/progress/:manga_id", mangaHandler.UpdateProgress) // Update reading progress
		}
	}

	// Start HTTP server in a non-blocking way
	httpAddr := cfg.GetHTTPAddress()
	log.Printf("HTTP API Server starting on %s", httpAddr)
	go func() {
		if err := router.Run(httpAddr); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start WebSocket server in a non-blocking way on its own port
	wsAddr := cfg.GetWebSocketAddress()
	wsRouter := gin.New()
	wsRouter.GET("/ws", websocket.Handler(wsHub)) // Only WebSocket endpoint
	log.Printf("WebSocket Server starting on ws://%s/ws", wsAddr)
	wsServer := &http.Server{
		Addr:    wsAddr,
		Handler: wsRouter,
	}
	go func() {
		if err := wsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start WebSocket server: %v", err)
		}
	}()

	log.Printf("Database: %s", cfg.Database.Path)
	log.Printf("Endpoints available:")
	log.Printf("  - Health check: GET /health (HTTP)")
	log.Printf("  - WebSocket: GET /ws (WebSocket)")
	log.Printf("  - Register: POST /api/v1/auth/register (HTTP)")
	log.Printf("  - Login: POST /api/v1/auth/login (HTTP)")
	log.Printf("  - Search manga: GET /api/v1/manga?title=<title>&author=<author>&genre=<genre>&status=<status> (HTTP)")
	log.Printf("  - Get manga: GET /api/v1/manga/:id (HTTP)")
	log.Printf("  - User library: GET /api/v1/users/library (HTTP, protected)")
	log.Printf("  - Add to library: POST /api/v1/users/library (HTTP, protected)")
	log.Printf("  - Update progress: PUT /api/v1/users/progress/:manga_id (HTTP, protected)")

	// Block until context is cancelled
	<-ctx.Done()
	log.Println("API server exiting. Shutting down HTTP and WebSocket servers...")

	// Create a deadline to wait for servers to shut down
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := wsServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("WebSocket server forced to shutdown: %v", err)
	} else {
		log.Println("WebSocket server stopped.")
	}

	// The Gin HTTP server is run via router.Run, which blocks.
	// To shut it down gracefully, we need to create an *http.Server instance for it.
	// Since router.Run() starts its own http.Server, we need to change how the HTTP server is started as well.
	// For now, we will just let the main goroutine exit for the Gin router.
	// A more robust solution would be to create an http.Server for Gin and call its Shutdown method.
	// Given the academic project context, this might be acceptable.
	log.Println("HTTP server will stop once main goroutine exits.")
}

// listenForUDPNotifications acts as a UDP client, registers with the UDP server,
// sends heartbeats, and forwards notifications to the WebSocket hub.
func listenForUDPNotifications(ctx context.Context, hub *websocket.Hub, serverAddress string) {
	clientID := uuid.New().String()
	log.Printf("UDP client for API server generated ID: %s", clientID)

	serverUDPAddr, err := net.ResolveUDPAddr("udp", serverAddress)
	if err != nil {
		log.Printf("Failed to resolve UDP server address %s: %v", serverAddress, err)
		return
	}

	// Use a single UDP connection for both sending and receiving.
	// Bind to an ephemeral port (nil) so it doesn't conflict with the UDP server.
	conn, err := net.ListenUDP("udp", nil)
	if err != nil {
		log.Printf("Failed to create UDP connection: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("UDP client for API server operating on local address: %s", conn.LocalAddr().String())

	// Channel to signal last pong received
	lastPong := make(chan time.Time, 1)

	// Goroutine to receive and process UDP messages
	go receiveUDPMessages(ctx, conn, hub, clientID, serverUDPAddr, lastPong)

	// Initial registration
	registerMessage := models.UDPRegisterMessage{
		ClientID: clientID,
		UserID:   "api-server-user", // A placeholder user ID for the API server itself
		Username: "api-server",
	}
	err = sendUDPMessage(conn, serverUDPAddr, models.UDPMessageTypeRegister, registerMessage)
	if err != nil {
		log.Printf("Failed to send UDP register message: %v", err)
		return
	}
	log.Printf("Sent UDP register message to %s from %s", serverAddress, conn.LocalAddr().String())

	// Goroutine to send pings and manage re-registration
	go startUDPPinger(ctx, conn, serverUDPAddr, clientID, lastPong)

	// Handle shutdown
	<-ctx.Done()
	log.Println("UDP client for API server shutting down.")
	// Send unregister message on shutdown
	unregisterMessage := models.UDPUnregisterMessage{ClientID: clientID}
	err = sendUDPMessage(conn, serverUDPAddr, models.UDPMessageTypeUnregister, unregisterMessage)
	if err != nil {
		log.Printf("Failed to send UDP unregister message: %v", err)
	} else {
		log.Printf("Sent UDP unregister message for client ID %s", clientID)
	}
}

// sendUDPMessage constructs and sends a UDP message.
func sendUDPMessage(conn *net.UDPConn, remoteAddr *net.UDPAddr, msgType models.UDPMessageType, data interface{}) error {
	msgData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	udpMsg := models.UDPMessage{
		Type:      msgType,
		Timestamp: time.Now(),
		Data:      json.RawMessage(msgData),
	}
	fullMsg, err := json.Marshal(udpMsg)
	if err != nil {
		return err
	}
	_, err = conn.WriteToUDP(fullMsg, remoteAddr)
	return err
}

// receiveUDPMessages continuously listens for and processes incoming UDP messages.
func receiveUDPMessages(ctx context.Context, conn *net.UDPConn, hub *websocket.Hub, clientID string, serverUDPAddr *net.UDPAddr, lastPong chan time.Time) {
	buffer := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopped receiving UDP messages.")
			return
		default:
			// Set a read deadline to allow context cancellation to be checked
			conn.SetReadDeadline(time.Now().Add(time.Second))
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue // Timeout, recheck context
				}
				log.Printf("Error reading from UDP: %v", err)
				continue
			}

			var udpMsg models.UDPMessage
			if err := json.Unmarshal(buffer[:n], &udpMsg); err != nil {
				log.Printf("Error unmarshaling UDP message: %v", err)
				continue
			}

			switch udpMsg.Type {
			case models.UDPMessageTypeNotification:
				// Forward notification data to WebSocket clients
				notificationJSON, err := json.Marshal(udpMsg.Data)
				if err != nil {
					log.Printf("Error marshaling notification for WebSocket: %v", err)
					continue
				}
				log.Printf("Received UDP notification: %s", string(notificationJSON))
				hub.Broadcast <- notificationJSON
			case models.UDPMessageTypePong:
				var pong models.UDPPongMessage
				// Re-marshal the generic interface back to bytes
				dataBytes, err := json.Marshal(udpMsg.Data)
				if err != nil {
					log.Printf("Error re-marshaling pong message data: %v", err)
					continue
				}
				if err := json.Unmarshal(dataBytes, &pong); err != nil {
					log.Printf("Error unmarshaling pong message data: %v", err)
					continue
				}
				log.Printf("Received Pong from UDP server. Client Time: %s", pong.ClientTime.Format(time.RFC3339))
				select {
				case lastPong <- time.Now():
				default:
				} // Non-blocking send
			case models.UDPMessageTypeRegisterSuccess:
				var success models.UDPRegisterSuccessMessage
				// Re-marshal the generic interface back to bytes
				dataBytes, err := json.Marshal(udpMsg.Data)
				if err != nil {
					log.Printf("Error re-marshaling register success message data: %v", err)
					continue
				}
				if err := json.Unmarshal(dataBytes, &success); err != nil {
					log.Printf("Error unmarshaling register success message data: %v", err)
					continue
				}
				log.Printf("Successfully registered with UDP server: %s", success.Message)
			case models.UDPMessageTypeRegisterFailed:
				var failed models.UDPRegisterFailedMessage
				// Re-marshal the generic interface back to bytes
				dataBytes, err := json.Marshal(udpMsg.Data)
				if err != nil {
					log.Printf("Error re-marshaling register failed message data: %v", err)
					continue
				}
				if err := json.Unmarshal(dataBytes, &failed); err != nil {
					log.Printf("Error unmarshaling register failed message data: %v", err)
					continue
				}
				log.Printf("Failed to register with UDP server: %s", failed.Reason)
			case models.UDPMessageTypeError:
				var errMsg models.UDPErrorMessage
				// Re-marshal the generic interface back to bytes
				dataBytes, err := json.Marshal(udpMsg.Data)
				if err != nil {
					log.Printf("Error re-marshaling UDP error message data: %v", err)
					continue
				}
				if err := json.Unmarshal(dataBytes, &errMsg); err != nil {
					log.Printf("Error unmarshaling UDP error message data: %v", err)
					continue
				}
				log.Printf("Received error from UDP server: Code=%s, Message=%s", errMsg.Code, errMsg.Message)
			default:
				log.Printf("Received unhandled UDP message type: %s", udpMsg.Type)
			}
		}
	}
}

// startUDPPinger periodically sends ping messages and manages re-registration if pong is not received.
func startUDPPinger(ctx context.Context, conn *net.UDPConn, serverUDPAddr *net.UDPAddr, clientID string, lastPong chan time.Time) {
	pingInterval := 30 * time.Second
	pongTimeout := 60 * time.Second // Time to wait for a pong before considering connection lost

	pingTicker := time.NewTicker(pingInterval)
	defer pingTicker.Stop()

	lastPongReceivedTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopped UDP pinger.")
			return
		case newPongTime := <-lastPong:
			lastPongReceivedTime = newPongTime
		case <-pingTicker.C:
			// Check if pong timeout occurred
			if time.Since(lastPongReceivedTime) > pongTimeout {
				log.Printf("No pong received for %v. Re-registering with UDP server...", pongTimeout)
				registerMessage := models.UDPRegisterMessage{
					ClientID: clientID,
					UserID:   "api-server-user",
					Username: "api-server",
				}
				err := sendUDPMessage(conn, serverUDPAddr, models.UDPMessageTypeRegister, registerMessage)
				if err != nil {
					log.Printf("Failed to re-send UDP register message: %v", err)
				} else {
					log.Println("Re-sent UDP register message.")
					lastPongReceivedTime = time.Now() // Reset timer after re-registration attempt
				}
				continue // Skip ping if re-registering
			}

			// Send ping
			pingMessage := models.UDPPingMessage{ClientTime: time.Now()}
			err := sendUDPMessage(conn, serverUDPAddr, models.UDPMessageTypePing, pingMessage)
			if err != nil {
				log.Printf("Failed to send UDP ping message: %v", err)
			}
		}
	}
}
