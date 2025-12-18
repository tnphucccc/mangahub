package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tnphucccc/mangahub/pkg/models"
	"github.com/tnphucccc/mangahub/pkg/response"
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
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Register user
	user, token, err := h.service.Register(req)
	if err != nil {
		// Check for duplicate user
		if err.Error() == "username or email already exists" {
			response.Conflict(c, err.Error())
			return
		}

		response.InternalError(c, "Failed to register user")
		return
	}

	// Return user and token
	response.Success(c, http.StatusCreated, gin.H{
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
		response.BadRequest(c, "Invalid request body")
		return
	}

	// Login user
	user, token, err := h.service.Login(req)
	if err != nil {
		response.Unauthorized(c, "Invalid username or password")
		return
	}

	// Return user and token
	response.Success(c, http.StatusOK, gin.H{
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
		response.Unauthorized(c, "Unauthorized")
		return
	}

	user := userInterface.(*models.User)
	response.Success(c, http.StatusOK, gin.H{
		"user": user.ToResponse(),
	})
}
