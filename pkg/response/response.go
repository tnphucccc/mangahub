package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standardized API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// APIError represents error details in the response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Meta contains pagination and other metadata
type Meta struct {
	Total      int  `json:"total,omitempty"`
	Count      int  `json:"count,omitempty"`
	Limit      int  `json:"limit,omitempty"`
	Offset     int  `json:"offset,omitempty"`
	HasMore    bool `json:"has_more,omitempty"`
	Page       int  `json:"page,omitempty"`
	TotalPages int  `json:"total_pages,omitempty"`
}

// PaginatedData wraps data with pagination info
type PaginatedData struct {
	Items interface{} `json:"items"`
}

// Success sends a successful response with data
func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}

// SuccessWithMeta sends a successful response with data and metadata
func SuccessWithMeta(c *gin.Context, statusCode int, data interface{}, meta *Meta) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

// Common error responses

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message)
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, "CONFLICT", message)
}

// InternalError sends a 500 Internal Server Error response
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}

// ValidationError sends a 422 Unprocessable Entity response
func ValidationError(c *gin.Context, message string) {
	Error(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", message)
}

// Paginated sends a successful response with paginated data
func Paginated(c *gin.Context, items interface{}, total, limit, offset int) {
	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}

	totalPages := 0
	if limit > 0 && total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	hasMore := offset+limit < total

	SuccessWithMeta(c, http.StatusOK, gin.H{"items": items}, &Meta{
		Total:      total,
		Count:      limit,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
		Page:       page,
		TotalPages: totalPages,
	})
}
