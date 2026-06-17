package database

import (
	"fmt"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedMedicineTypes(tx *gorm.DB) ([]models.MedicineType, error) {
	log.Println("🌱 Seeding medicine types...")

	var created []models.MedicineType

	// 1. 12 data dengan is_active true, deleted_at null
	for i := 1; i <= 12; i++ {
		mt := models.MedicineType{
			Name:        fmt.Sprintf("Active Medicine Type %02d", i),
			Code:        fmt.Sprintf("MT-ACT-%02d", i),
			Description: fmt.Sprintf("Active medicine type description %02d", i),
			IsActive:    true,
		}
		if err := tx.Where("code = ?", mt.Code).FirstOrCreate(&mt).Error; err != nil {
			return nil, err
		}
		created = append(created, mt)
	}

	// 2. 12 data dengan is_active false, deleted_at null
	for i := 1; i <= 12; i++ {
		mt := models.MedicineType{
			Name:        fmt.Sprintf("Inactive Medicine Type %02d", i),
			Code:        fmt.Sprintf("MT-INA-%02d", i),
			Description: fmt.Sprintf("Inactive medicine type description %02d", i),
			IsActive:    false,
		}
		if err := tx.Where("code = ?", mt.Code).FirstOrCreate(&mt).Error; err != nil {
			return nil, err
		}
		created = append(created, mt)
	}

	return created, nil
}

func seedDeletedMedicineTypes(tx *gorm.DB) error {
	log.Println("🌱 Seeding deleted medicine types...")
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)

	// 3. 22 data dengan deleted_at terisi
	var deletedTypes []models.MedicineType
	for i := 1; i <= 22; i++ {
		deletedTypes = append(deletedTypes, models.MedicineType{
			Name:        fmt.Sprintf("Deleted Medicine Type %02d", i),
			Code:        fmt.Sprintf("MT-DEL-%02d", i),
			Description: fmt.Sprintf("Deleted medicine type description %02d", i),
			IsActive:    false,
			DeletedAt:   gorm.DeletedAt{Time: lastMonth, Valid: true},
		})
	}

	if err := tx.Create(&deletedTypes).Error; err != nil {
		return err
	}
	return nil
}
