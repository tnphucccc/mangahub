package manga

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tnphucccc/mangahub/pkg/models"
	"github.com/tnphucccc/mangahub/pkg/response"
)

// Handler handles manga HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new manga handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Search searches for manga
// GET /manga?q=query&genre=action&status=ongoing&limit=20&offset=0
func (h *Handler) Search(c *gin.Context) {
	var query models.MangaSearchQuery

	// Bind query parameters
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "Invalid query parameters")
		return
	}

	// Search manga
	mangaList, err := h.service.Search(query)
	if err != nil {
		response.InternalError(c, "Failed to search manga")
		return
	}

	// Use pagination response for consistency
	response.SuccessWithMeta(c, http.StatusOK, gin.H{"items": mangaList}, &response.Meta{
		Count:  len(mangaList),
		Limit:  query.Limit,
		Offset: query.Offset,
	})
}

// GetByID retrieves a manga by ID
// GET /manga/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	manga, err := h.service.GetByID(id)
	if err != nil {
		response.NotFound(c, "Manga not found")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"manga": manga})
}

// GetAll retrieves all manga with pagination
// GET /manga/all?limit=20&offset=0
func (h *Handler) GetAll(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	mangaList, err := h.service.GetAll(limit, offset)
	if err != nil {
		response.InternalError(c, "Failed to get manga")
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, gin.H{"items": mangaList}, &response.Meta{
		Count:  len(mangaList),
		Limit:  limit,
		Offset: offset,
	})
}

// GetLibrary retrieves user's manga library
// GET /users/library
func (h *Handler) GetLibrary(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	user := userInterface.(*models.User)

	library, err := h.service.GetUserLibrary(user.ID)
	if err != nil {
		response.InternalError(c, "Failed to get library")
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, gin.H{"items": library}, &response.Meta{
		Count: len(library),
	})
}

// AddToLibrary adds a manga to user's library
// POST /users/library
func (h *Handler) AddToLibrary(c *gin.Context) {
	// Get user ID from context
	userInterface, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	user := userInterface.(*models.User)

	var req models.LibraryAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.AddToLibrary(user.ID, req); err != nil {
		if err.Error() == "manga not found" {
			response.NotFound(c, "Manga not found")
			return
		}

		response.InternalError(c, "Failed to add manga to library")
		return
	}

	response.Success(c, http.StatusCreated, gin.H{"message": "Manga added to library"})
}

// UpdateProgress updates user's reading progress
// PUT /users/progress/:manga_id
func (h *Handler) UpdateProgress(c *gin.Context) {
	// Get user ID from context
	userInterface, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	user := userInterface.(*models.User)
	mangaID := c.Param("manga_id")

	var req models.ProgressUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.UpdateProgress(user.ID, mangaID, req); err != nil {
		if err.Error() == "manga not found" {
			response.NotFound(c, "Manga not found")
			return
		}

		if err.Error() == "progress not found" || err.Error() == "manga not in user's library" {
			response.NotFound(c, "Manga not in library")
			return
		}

		if err.Error() == "invalid chapter number" {
			response.BadRequest(c, "Invalid chapter number")
			return
		}

		response.InternalError(c, "Failed to update progress")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "Progress updated"})
}

// GetProgress retrieves user's progress for a specific manga
// GET /users/progress/:manga_id
func (h *Handler) GetProgress(c *gin.Context) {
	// Get user ID from context
	userInterface, exists := c.Get("user")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}

	user := userInterface.(*models.User)
	mangaID := c.Param("manga_id")

	progress, err := h.service.GetProgress(user.ID, mangaID)
	if err != nil {
		response.NotFound(c, "Progress not found")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"progress": progress})
}
