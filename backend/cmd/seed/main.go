package main

import (
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/database"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	if err := database.RunMigration(db); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	if err := database.SeedDatabase(db); err != nil {
		log.Fatalf("❌ Failed to seed database: %v", err)
	}

	log.Println("✅ Seed command finished successfully")
}
