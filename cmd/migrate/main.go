package main

import (
	"flag"
	"log"

	"github.com/ndxbinh1922001/VNalo-be/internal/initialize"
)

func main() {
	var command string
	flag.StringVar(&command, "cmd", "up", "Migration command: up, down, status")
	flag.Parse()

	log.Println("🗄️  VNalo Migration Tool (Goose + PostgreSQL)")

	cfg := initialize.LoadConfig()

	switch command {
	case "up":
		log.Println("▶️  Running migrations UP...")
		if err := initialize.RunMigrations(cfg); err != nil {
			log.Fatalf("❌ Migration failed: %v", err)
		}
		log.Println("✅ Migrations applied successfully")

	case "down":
		log.Println("◀️  Rolling back last migration...")
		if err := initialize.RollbackMigration(cfg); err != nil {
			log.Fatalf("❌ Rollback failed: %v", err)
		}
		log.Println("✅ Rollback completed successfully")

	case "status":
		log.Println("📊 Checking migration status...")
		if err := initialize.GetMigrationStatus(cfg); err != nil {
			log.Fatalf("❌ Failed to get status: %v", err)
		}

	default:
		log.Fatalf("❌ Unknown command: %s. Use: up, down, or status", command)
	}
}

