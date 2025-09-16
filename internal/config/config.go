package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port            string
	RustServiceURL  string
	RedisURL        string
	JWTSecret       string
	RateLimitRPM    int
	AllowedOrigins  []string
}

func Load() *Config {
	rpm, _ := strconv.Atoi(getEnv("RATE_LIMIT_RPM", "60"))
	origins := strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ",")

	return &Config{
		Port:            getEnv("PORT", "3000"),
		RustServiceURL:  getEnv("RUST_SERVICE_URL", "http://localhost:8080"),
		RedisURL:        getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:       getEnv("JWT_SECRET", "super-secret-key"),
		RateLimitRPM:    rpm,
		AllowedOrigins:  origins,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}