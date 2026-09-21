package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

const redis_url = "redis://default:H1OuWhD69wVLnfd1BfTTA2WU2xNMcOX2@redis-17276.c281.us-east-1-2.ec2.cloud.redislabs.com:17276"

func ConnectRedis() {
	ctx := context.Background()

	// Get address from environment variable, fallback to 127.0.0.1:6379
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379" // Use 127.0.0.1 to avoid IPv6 issues on Windows
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Ping the database to ensure connection is working
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v\n", err)
	}
	
	log.Println("Connected to Redis successfully!")
}