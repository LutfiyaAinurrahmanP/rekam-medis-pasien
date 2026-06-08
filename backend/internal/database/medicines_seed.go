package database

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedMedicines(tx *gorm.DB, count int, medicineTypes []models.MedicineType) ([]models.Medicine, error) {
	manufacturers := []string{"Kimia Farma", "Kalbe Farma", "Dexa Medica", "Sanbe", "Novell", "Tempo Scan"}
	units := []string{"strip", "bottle", "vial", "tube", "box"}

	medicines := make([]models.Medicine, 0, count)
	for i := 1; i <= count; i++ {
		medicineType := medicineTypes[(i-1)%len(medicineTypes)]
		medicine := models.Medicine{
			Name:          fmt.Sprintf("Medicine Sample %02d", i),
			GenericName:   fmt.Sprintf("Generic Compound %02d", i),
			BrandName:     fmt.Sprintf("Brand %02d", i),
			MedicineTypeID: medicineType.ID,
			Strength:      fmt.Sprintf("%d mg", 50+((i-1)%10)*50),
			Manufacturer:  manufacturers[(i-1)%len(manufacturers)],
			Unit:          units[(i-1)%len(units)],
			StockQuantity: 5 + (i * 2),
			Price:         float64(5000 + (i * 750)),
			IsActive:      i%10 != 0,
		}
		medicines = append(medicines, medicine)
	}

	if err := tx.Create(&medicines).Error; err != nil {
		return nil, fmt.Errorf("failed to seed medicines: %w", err)
	}

	return medicines, nil
}

func seedDeletedMedicines(tx *gorm.DB, medicineTypes []models.MedicineType) error {
	manufacturers := []string{"Kimia Farma", "Kalbe Farma", "Dexa Medica", "Sanbe", "Novell", "Tempo Scan"}
	units := []string{"strip", "bottle", "vial", "tube", "box"}

	deletedTime := time.Now().AddDate(0, -1, 0) // 1 bulan yang lalu
	deletedAt := gorm.DeletedAt{Time: deletedTime, Valid: true}

	medicines := make([]models.Medicine, 0, 12)
	for i := 6001; i <= 6012; i++ {
		medicineType := medicineTypes[(i-1)%len(medicineTypes)]
		medicine := models.Medicine{
			Name:          fmt.Sprintf("(Deleted) Medicine Sample %02d", i-6000),
			GenericName:   fmt.Sprintf("Generic Compound DEL-%02d", i-6000),
			BrandName:     fmt.Sprintf("Brand DEL-%02d", i-6000),
			MedicineTypeID: medicineType.ID,
			Strength:      fmt.Sprintf("%d mg", 50+((i-1)%10)*50),
			Manufacturer:  manufacturers[(i-1)%len(manufacturers)],
			Unit:          units[(i-1)%len(units)],
			StockQuantity: 0,
			Price:         float64(5000 + ((i - 6000) * 750)),
			IsActive:      false,
			DeletedAt:     deletedAt,
		}
		medicines = append(medicines, medicine)
	}

	if err := tx.Create(&medicines).Error; err != nil {
		return fmt.Errorf("failed to seed deleted medicines: %w", err)
	}

	return nil
}
