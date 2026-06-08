package database

import (
	"errors"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	log.Println("🔄 Running database migrations...")

	if err := db.AutoMigrate(
		&models.User{},
		&models.Department{},
		&models.Patient{},
		&models.DoctorSpecialization{},
		&models.Room{},
		&models.TypeTestCategory{},
		&models.TypeTest{},
		&models.Medicine{},
		&models.Appointment{},
		&models.MedicineType{},
		&models.RoomType{},
	); err != nil {
		return err
	}

	if err := migrateDoctorSpecializationID(db); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Doctor{}); err != nil {
		return err
	}

	// Hapus kolom lama jika masih ada untuk menghindari error NOT NULL saat seeding
	if err := db.Exec(`ALTER TABLE doctors DROP COLUMN IF EXISTS specialization`).Error; err != nil {
		return err
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

func migrateDoctorSpecializationID(db *gorm.DB) error {
	if !db.Migrator().HasTable(&models.Doctor{}) {
		return nil
	}

	if db.Migrator().HasColumn(&models.Doctor{}, "SpecializationID") {
		return nil
	}

	if err := db.Exec(`ALTER TABLE doctors ADD COLUMN IF NOT EXISTS specialization_id bigint`).Error; err != nil {
		return err
	}

	defaultSpecializationID, err := ensureDefaultDoctorSpecialization(db)
	if err != nil {
		return err
	}

	if err := db.Exec(`UPDATE doctors SET specialization_id = ? WHERE specialization_id IS NULL`, defaultSpecializationID).Error; err != nil {
		return err
	}

	if err := db.Exec(`ALTER TABLE doctors ALTER COLUMN specialization_id SET NOT NULL`).Error; err != nil {
		return err
	}

	return nil
}

func ensureDefaultDoctorSpecialization(db *gorm.DB) (uint, error) {
	var spec models.DoctorSpecialization
	if err := db.Where("code = ?", "unspecified").First(&spec).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			spec = models.DoctorSpecialization{
				Name:        "Unspecified",
				Code:        "unspecified",
				Description: "Default specialization assigned to legacy doctor records.",
				IsActive:    true,
			}
			if err := db.Create(&spec).Error; err != nil {
				return 0, err
			}
			return spec.ID, nil
		}
		return 0, err
	}

	return spec.ID, nil
}
