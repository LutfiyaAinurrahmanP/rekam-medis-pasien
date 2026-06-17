package database

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedTypeTests(tx *gorm.DB, count int) ([]models.TypeTest, error) {
	var categories []models.TypeTestCategory
	if err := tx.Where("is_active = ?", true).Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch type test categories: %w", err)
	}

	if len(categories) == 0 {
		return nil, fmt.Errorf("no active type test categories found")
	}

	typeTests := make([]models.TypeTest, 0, count)
	for i := 1; i <= count; i++ {
		category := categories[(i-1)%len(categories)]
		typeTest := models.TypeTest{
			Name:               fmt.Sprintf("%s Test %02d", category.Name, i),
			Code:               fmt.Sprintf("LAB-%s-%03d", category.Code, i),
			TypeTestCategoryID: category.ID,
			Description:        fmt.Sprintf("Pemeriksaan %s untuk parameter klinis %02d", category.Name, i),
			Price:              float64(75000 + (i * 5000)),
			IsActive:           i%9 != 0,
		}
		typeTests = append(typeTests, typeTest)
	}

	if err := tx.Create(&typeTests).Error; err != nil {
		return nil, fmt.Errorf("failed to seed type tests: %w", err)
	}

	return typeTests, nil
}

func seedDeletedTypeTests(tx *gorm.DB) error {
	var categories []models.TypeTestCategory
	if err := tx.Unscoped().Find(&categories).Error; err != nil {
		return fmt.Errorf("failed to fetch type test categories: %w", err)
	}

	if len(categories) == 0 {
		return fmt.Errorf("no type test categories found")
	}

	deletedTime := time.Now().AddDate(0, -1, 0) // 1 bulan yang lalu
	deletedAt := gorm.DeletedAt{Time: deletedTime, Valid: true}

	typeTests := make([]models.TypeTest, 0, 12)
	for i := 5001; i <= 5012; i++ {
		category := categories[(i-1)%len(categories)]
		typeTest := models.TypeTest{
			Name:               fmt.Sprintf("(Deleted) %s Test %02d", category.Name, i-5000),
			Code:               fmt.Sprintf("LAB-DEL-%s-%03d", category.Code, i-5000),
			TypeTestCategoryID: category.ID,
			Description:        fmt.Sprintf("Pemeriksaan %s yang telah dihapus - #%02d", category.Name, i-5000),
			Price:              float64(75000 + ((i - 5000) * 5000)),
			IsActive:           false,
			DeletedAt:          deletedAt,
		}
		typeTests = append(typeTests, typeTest)
	}

	if err := tx.Create(&typeTests).Error; err != nil {
		return fmt.Errorf("failed to seed deleted type tests: %w", err)
	}

	return nil
}
