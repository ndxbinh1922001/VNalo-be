package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

var Pool *pgxpool.Pool

// NewPool creates a new PostgreSQL connection pool
func NewPool(cfg Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Connection pool settings
	config.MaxConns = 50
	config.MinConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	Pool = pool
	log.Println("PostgreSQL connection pool created successfully")
	return pool, nil
}

// GetPool returns the global pool instance
func GetPool() *pgxpool.Pool {
	if Pool == nil {
		log.Fatal("PostgreSQL pool not initialized. Call NewPool() first.")
	}
	return Pool
}

// Close closes the connection pool
func Close() {
	if Pool != nil {
		Pool.Close()
		log.Println("PostgreSQL connection pool closed")
	}
}

