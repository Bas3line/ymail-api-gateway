package middleware

import (
	"net/http"
	"strings"

	"api-gateway/internal/models"
	"api-gateway/internal/services"

	"github.com/gin-gonic/gin"
)

func RateLimit(rateLimiter *services.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rateLimiter.CheckLimit(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
				Error: "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequestFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// Block suspicious patterns
		suspiciousPatterns := []string{
			"/admin", "/wp-admin", "/phpmyadmin", "/.env",
			"../", "script>", "javascript:", "eval(",
		}

		for _, pattern := range suspiciousPatterns {
			if strings.Contains(path, pattern) {
				c.JSON(http.StatusForbidden, models.ErrorResponse{
					Error: "Forbidden request",
				})
				c.Abort()
				return
			}
		}

		// Validate content type for POST/PUT requests
		if method == "POST" || method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") &&
			   !strings.Contains(contentType, "multipart/form-data") &&
			   !strings.Contains(contentType, "application/x-www-form-urlencoded") {
				c.JSON(http.StatusUnsupportedMediaType, models.ErrorResponse{
					Error: "Unsupported media type",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}