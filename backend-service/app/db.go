package app

import (
	"context"
	"jirbthagoras/raksana-backend/helpers"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetConnection() *pgxpool.Pool {
	cnf := helpers.NewConfig()
	dbUrl := cnf.GetString("DATABASE_URL")

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		slog.Error("Failed to establish connection to Database", "err", err.Error())
		os.Exit(1)
	}
	config.MaxConns = 20
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		slog.Error("Failed to establish connection to Database", "err", err.Error())
		os.Exit(1)
	}

	slog.Debug("Established connection to db")
	return pool
}
