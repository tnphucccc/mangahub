package stats

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tnphucccc/mangahub/pkg/models"
	"github.com/tnphucccc/mangahub/pkg/response"
)

// Handler handles statistics HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new stats handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetStats returns the authenticated user's reading statistics
func (h *Handler) GetStats(c *gin.Context) {
	// Get user from context (set by auth middleware)
	u, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	user, ok := u.(*models.User)
	if !ok {
		response.InternalError(c, "Invalid user context")
		return
	}

	stats, err := h.service.GetUserStats(user.ID)
	if err != nil {
		response.InternalError(c, "Failed to get statistics")
		return
	}

	response.Success(c, http.StatusOK, stats)
}
