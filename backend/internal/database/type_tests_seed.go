package database

import (
	"fmt"
	"time"

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

func seedDeletedTypeTests(tx *gorm.DB) error {
	testCategories := []string{"Hematologi", "Kimia Darah", "Mikrobiologi", "Urinalisis", "Imunologi", "Serologi"}

	deletedTime := time.Now().AddDate(0, -1, 0) // 1 bulan yang lalu
	deletedAt := gorm.DeletedAt{Time: deletedTime, Valid: true}

	typeTests := make([]models.TypeTest, 0, 12)
	for i := 5001; i <= 5012; i++ {
		category := testCategories[(i-1)%len(testCategories)]
		typeTest := models.TypeTest{
			Name:        fmt.Sprintf("(Deleted) %s Test %02d", category, i-5000),
			Code:        fmt.Sprintf("LAB-DEL-%s-%03d", codePrefix(category), i-5000),
			Category:    category,
			Description: fmt.Sprintf("Pemeriksaan %s yang telah dihapus - #%02d", category, i-5000),
			Price:       float64(75000 + ((i - 5000) * 5000)),
			IsActive:    false,
			DeletedAt:   deletedAt,
		}
		typeTests = append(typeTests, typeTest)
	}

	if err := tx.Create(&typeTests).Error; err != nil {
		return fmt.Errorf("failed to seed deleted type tests: %w", err)
	}

	return nil
}
