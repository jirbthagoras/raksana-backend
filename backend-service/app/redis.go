package app

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	slog.Info("Established connection to redis")

	return client
}
