package services

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

type RateLimiter struct {
	redis *RedisClient
	limit int
}

func NewRedisClient(redisURL string) *RedisClient {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Warning: Invalid Redis URL, using default: %v", err)
		opt = &redis.Options{Addr: "localhost:6379"}
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	}

	return &RedisClient{client: client}
}

func NewRateLimiter(redisClient *RedisClient, limit int) *RateLimiter {
	return &RateLimiter{
		redis: redisClient,
		limit: limit,
	}
}

func (rl *RateLimiter) CheckLimit(clientIP string) bool {
	if rl.redis == nil || rl.redis.client == nil {
		return true // Allow if Redis is not available
	}

	ctx := context.Background()
	key := "rate_limit:" + clientIP

	current, err := rl.redis.client.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("Redis error: %v", err)
		return true // Allow on error
	}

	if current == 1 {
		rl.redis.client.Expire(ctx, key, time.Minute)
	}

	return current <= int64(rl.limit)
}