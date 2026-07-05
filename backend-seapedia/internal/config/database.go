package config

import (
	"backend-seapedia/internal/model"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildDSN(cfg *model.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
}

func ConnectDB(cfg *model.Config) (*pgxpool.Pool, error) {
	dsn := BuildDSN(cfg)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("gagal parse dsn database: %v", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("gagal konek ke database: %v", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("database tidak merespon: %v", err)
	}
	return pool, nil
}
