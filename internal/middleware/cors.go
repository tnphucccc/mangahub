package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tnphucccc/mangahub/pkg/utils"
)

// CORSConfig holds CORS configuration options
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
	}
}

// CORS middleware for handling Cross-Origin Resource Sharing
// Uses environment variable CORS_ALLOWED_ORIGINS to configure allowed origins
// Default: "*" (all origins allowed - for development)
// Production example: CORS_ALLOWED_ORIGINS="https://mangahub.com,https://app.mangahub.com"
func CORS() gin.HandlerFunc {
	config := DefaultCORSConfig()

	// Override allowed origins from environment
	if origins := utils.GetEnv("CORS_ALLOWED_ORIGINS", ""); origins != "" {
		config.AllowOrigins = strings.Split(origins, ",")
	}

	return CORSWithConfig(config)
}

// CORSWithConfig creates a CORS middleware with custom configuration
func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	allowOrigins := strings.Join(config.AllowOrigins, ", ")
	allowMethods := strings.Join(config.AllowMethods, ", ")
	allowHeaders := strings.Join(config.AllowHeaders, ", ")

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowedOrigin := "*"
		if len(config.AllowOrigins) > 0 && config.AllowOrigins[0] != "*" {
			for _, o := range config.AllowOrigins {
				if strings.TrimSpace(o) == origin {
					allowedOrigin = origin
					break
				}
			}
			// If origin not in allowed list, use first allowed origin
			if allowedOrigin == "*" {
				allowedOrigin = strings.TrimSpace(config.AllowOrigins[0])
			}
		} else {
			allowedOrigin = allowOrigins
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", allowMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		if config.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
