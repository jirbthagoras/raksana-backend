package app

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

func NewRedisClient() *redis.Client {
	// Create new Redis Client
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	// Pings the Redis Client
	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	slog.Info("Connected to Redis")

	// Returning the created Client
	return client
}
