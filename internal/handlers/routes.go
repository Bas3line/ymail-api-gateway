package handlers

import (
	"net/http"
	"time"

	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/models"
	"api-gateway/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, proxyService *services.ProxyService) {
	// Health check
	r.GET("/health", HealthCheck)
	r.GET("/api/health", HealthCheck)

	// Public routes (no auth required)
	public := r.Group("/api/v1")
	{
		public.GET("/auth/google", ProxyHandler(proxyService))
		public.GET("/auth/google/callback", ProxyHandler(proxyService))
		public.GET("/auth/microsoft", ProxyHandler(proxyService))
		public.GET("/auth/microsoft/callback", ProxyHandler(proxyService))
	}

	// Protected routes (auth required)
	protected := r.Group("/api/v1")
	protected.Use(middleware.Auth(cfg))
	{
		protected.GET("/emails", ProxyHandler(proxyService))
		protected.GET("/emails/:id", ProxyHandler(proxyService))
		protected.POST("/emails/send", ProxyHandler(proxyService))
		protected.POST("/auth/refresh", ProxyHandler(proxyService))
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.HealthResponse{
		Status:    "healthy",
		Service:   "api-gateway",
		Timestamp: time.Now().Unix(),
	})
}

func ProxyHandler(proxyService *services.ProxyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxyService.ServeHTTP(c.Writer, c.Request)
	}
}