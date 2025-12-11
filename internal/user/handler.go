package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tnphucccc/mangahub/pkg/models"
)

// Handler handles user HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new user handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register handles user registration
// POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req models.UserRegisterRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Register user
	user, token, err := h.service.Register(req)
	if err != nil {
		// Check for duplicate user
		if err.Error() == "username or email already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register user",
		})
		return
	}

	// Return user and token
	c.JSON(http.StatusCreated, gin.H{
		"user":  user.ToResponse(),
		"token": token,
	})
}

// Login handles user login
// POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req models.UserLoginRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Login user
	user, token, err := h.service.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Return user and token
	c.JSON(http.StatusOK, gin.H{
		"user":  user.ToResponse(),
		"token": token,
	})
}

// GetProfile returns the current user's profile
// GET /users/me
func (h *Handler) GetProfile(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	user := userInterface.(*models.User)
	c.JSON(http.StatusOK, gin.H{
		"user": user.ToResponse(),
	})
}
