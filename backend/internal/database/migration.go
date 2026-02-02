package database

import (
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Department{},
	)

	if err != nil {
		return err
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}
