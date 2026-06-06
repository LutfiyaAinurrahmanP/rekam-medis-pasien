package database

import (
	"fmt"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

const defaultSeedCount = 22

func SeedDatabase(db *gorm.DB) error {
	return SeedDatabaseWithCount(db, defaultSeedCount)
}

func SeedDatabaseWithCount(db *gorm.DB, count int) error {
	if count < defaultSeedCount {
		count = defaultSeedCount
	}

	log.Printf("🌱 Starting database seeding with %d records per domain", count)

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin seed transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := resetSeedData(tx); err != nil {
		tx.Rollback()
		return err
	}

	users, err := seedUsers(tx, count)
	if err != nil {
		tx.Rollback()
		return err
	}

	departments, err := seedDepartments(tx, count)
	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err := seedPatients(tx, count, users); err != nil {
		tx.Rollback()
		return err
	}

	doctorSpecs, err := seedDoctorSpecializations(tx, count)
	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err := seedDoctors(tx, count, users, departments, doctorSpecs); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := seedRoomTypes(tx, 12); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := seedRooms(tx, count, departments); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := seedTypeTests(tx, count); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := seedMedicines(tx, count); err != nil {
		tx.Rollback()
		return err
	}

	// Seed deleted records (soft-deleted data with 12 records each)
	if err := seedDeletedUsers(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := seedDeletedDepartments(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := seedDeletedPatients(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := seedDeletedDoctorSpecializations(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := seedDeletedDoctors(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := seedDeletedRoomTypes(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := seedDeletedRooms(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := seedDeletedTypeTests(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := seedDeletedMedicines(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit seed transaction: %w", err)
	}

	log.Println("✅ Database seeding completed")
	return nil
}

func resetSeedData(tx *gorm.DB) error {
	session := tx.Session(&gorm.Session{AllowGlobalUpdate: true})

	if err := session.Unscoped().Delete(&models.Room{}).Error; err != nil {
		return fmt.Errorf("failed to clear rooms: %w", err)
	}
	if err := session.Unscoped().Delete(&models.RoomType{}).Error; err != nil {
		return fmt.Errorf("failed to clear room types: %w", err)
	}
	if err := session.Unscoped().Delete(&models.Doctor{}).Error; err != nil {
		return fmt.Errorf("failed to clear doctors: %w", err)
	}
	if err := session.Unscoped().Delete(&models.DoctorSpecialization{}).Error; err != nil {
		return fmt.Errorf("failed to clear doctor specializations: %w", err)
	}
	if err := session.Unscoped().Delete(&models.Patient{}).Error; err != nil {
		return fmt.Errorf("failed to clear patients: %w", err)
	}
	if err := session.Unscoped().Delete(&models.TypeTest{}).Error; err != nil {
		return fmt.Errorf("failed to clear type tests: %w", err)
	}
	if err := session.Unscoped().Delete(&models.Medicine{}).Error; err != nil {
		return fmt.Errorf("failed to clear medicines: %w", err)
	}
	if err := session.Unscoped().Delete(&models.Department{}).Error; err != nil {
		return fmt.Errorf("failed to clear departments: %w", err)
	}
	if err := session.Unscoped().Delete(&models.User{}).Error; err != nil {
		return fmt.Errorf("failed to clear users: %w", err)
	}

	return nil
}
