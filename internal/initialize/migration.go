package initialize

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver for Goose
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

// MigrationModule provides migration functionality
var MigrationModule = fx.Options(
	fx.Invoke(runMigrations),
)

// runMigrations runs database migrations on startup
func runMigrations(cfg *Config) error {
	log.Println("🔄 Running database migrations with Goose...")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("⚠️  Migration failed: %v", err)
		return nil // Don't fail app startup
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("⚠️  Database ping failed: %v", err)
		return nil
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Printf("⚠️  Failed to set goose dialect: %v", err)
		return nil
	}

	migrationsDir := "./migrations/postgres"
	if err := goose.Up(db, migrationsDir); err != nil {
		log.Printf("⚠️  Migration failed: %v", err)
		return nil
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// Helper functions for CLI tool
func RunMigrations(cfg *Config) error {
	return runMigrations(cfg)
}

func RollbackMigration(cfg *Config) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	migrationsDir := "./migrations/postgres"
	if err := goose.Down(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("✅ Rollback completed successfully")
	return nil
}

func GetMigrationStatus(cfg *Config) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
		cfg.Postgres.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	migrationsDir := "./migrations/postgres"
	if err := goose.Status(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}

