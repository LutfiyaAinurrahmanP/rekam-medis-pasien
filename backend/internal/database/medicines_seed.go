package database

import (
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedMedicines(tx *gorm.DB, count int) ([]models.Medicine, error) {
	medicineTypes := []string{"tablet", "capsule", "syrup", "injection", "ointment", "other"}
	manufacturers := []string{"Kimia Farma", "Kalbe Farma", "Dexa Medica", "Sanbe", "Novell", "Tempo Scan"}
	units := []string{"strip", "bottle", "vial", "tube", "box"}

	medicines := make([]models.Medicine, 0, count)
	for i := 1; i <= count; i++ {
		medicineType := medicineTypes[(i-1)%len(medicineTypes)]
		medicine := models.Medicine{
			Name:          fmt.Sprintf("Medicine Sample %02d", i),
			GenericName:   fmt.Sprintf("Generic Compound %02d", i),
			BrandName:     fmt.Sprintf("Brand %02d", i),
			Type:          medicineType,
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
