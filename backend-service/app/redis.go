package app

import (
	"context"
	"jirbthagoras/raksana-backend/helpers"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	cnf := helpers.NewConfig()
	addr := cnf.GetString("REDIS_ADDR")
	username := cnf.GetString("REDIS_USERNAME")
	password := cnf.GetString("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       0,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	slog.Info("Established connection to redis")

	return client
}
