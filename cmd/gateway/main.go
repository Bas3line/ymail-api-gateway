package main

import (
	"log"
	"os"

	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/handlers"
	"api-gateway/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Setup Redis client
	redisClient := services.NewRedisClient(cfg.RedisURL)

	// Create services
	rateLimiter := services.NewRateLimiter(redisClient, cfg.RateLimitRPM)
	proxyService := services.NewProxyService(cfg.RustServiceURL)

	// Setup Gin
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Setup middleware
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(middleware.RequestFilter())
	r.Use(middleware.RateLimit(rateLimiter))

	// Setup routes
	handlers.SetupRoutes(r, cfg, proxyService)

	log.Printf("API Gateway starting on port %s", cfg.Port)
	log.Printf("Proxying to Rust service: %s", cfg.RustServiceURL)
	log.Printf("Rate limit: %d requests per minute", cfg.RateLimitRPM)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}