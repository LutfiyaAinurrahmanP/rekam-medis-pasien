package database

import (
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedTypeTests(tx *gorm.DB, count int) ([]models.TypeTest, error) {
	testCategories := []string{"Hematologi", "Kimia Darah", "Mikrobiologi", "Urinalisis", "Imunologi", "Serologi"}

	typeTests := make([]models.TypeTest, 0, count)
	for i := 1; i <= count; i++ {
		category := testCategories[(i-1)%len(testCategories)]
		typeTest := models.TypeTest{
			Name:        fmt.Sprintf("%s Test %02d", category, i),
			Code:        fmt.Sprintf("LAB-%s-%03d", codePrefix(category), i),
			Category:    category,
			Description: fmt.Sprintf("Pemeriksaan %s untuk parameter klinis %02d", category, i),
			Price:       float64(75000 + (i * 5000)),
			IsActive:    i%9 != 0,
		}
		typeTests = append(typeTests, typeTest)
	}

	if err := tx.Create(&typeTests).Error; err != nil {
		return nil, fmt.Errorf("failed to seed type tests: %w", err)
	}

	return typeTests, nil
}

func codePrefix(category string) string {
	switch category {
	case "Hematologi":
		return "HEM"
	case "Kimia Darah":
		return "KIM"
	case "Mikrobiologi":
		return "MIK"
	case "Urinalisis":
		return "URI"
	case "Imunologi":
		return "IMU"
	case "Serologi":
		return "SER"
	default:
		return "GEN"
	}
}
