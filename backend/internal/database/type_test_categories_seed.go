package database

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedTypeTestCategories(tx *gorm.DB) ([]models.TypeTestCategory, error) {
	categories := []models.TypeTestCategory{
		// 12 active
		{Name: "Hematologi", Code: "HEM", Description: "Pemeriksaan yang berkaitan dengan darah dan komponen darah", IsActive: true},
		{Name: "Kimia Darah", Code: "KIM", Description: "Pemeriksaan kimia komponen darah", IsActive: true},
		{Name: "Mikrobiologi", Code: "MIK", Description: "Pemeriksaan mikroorganisme", IsActive: true},
		{Name: "Imunologi", Code: "IMU", Description: "Pemeriksaan sistem kekebalan tubuh", IsActive: true},
		{Name: "Serologi", Code: "SER", Description: "Pemeriksaan cairan serum", IsActive: true},
		{Name: "Urinalisis", Code: "URI", Description: "Pemeriksaan urin", IsActive: true},
		{Name: "Parasitologi", Code: "PAR", Description: "Pemeriksaan parasit", IsActive: true},
		{Name: "Patologi Anatomi", Code: "PAT", Description: "Pemeriksaan jaringan tubuh", IsActive: true},
		{Name: "Biologi Molekuler", Code: "MOL", Description: "Pemeriksaan molekuler seperti PCR", IsActive: true},
		{Name: "Toksikologi", Code: "TOK", Description: "Pemeriksaan racun dan obat dalam tubuh", IsActive: true},
		{Name: "Bank Darah", Code: "BDR", Description: "Pemeriksaan golongan darah dan crossmatch", IsActive: true},
		{Name: "Genetika", Code: "GEN", Description: "Pemeriksaan DNA dan genetik", IsActive: true},

		// 12 inactive
		{Name: "Tes Narkoba Konvensional", Code: "NRK-OLD", Description: "Metode lama tes narkoba", IsActive: false},
		{Name: "Serologi Kuno", Code: "SER-OLD", Description: "Tes serologi metode lama", IsActive: false},
		{Name: "Pemeriksaan Feses Rutin Lama", Code: "FES-OLD", Description: "Pemeriksaan feses metode lama", IsActive: false},
		{Name: "Tes Kehamilan Kualitatif Lama", Code: "PRG-OLD", Description: "Tes kehamilan metode lama", IsActive: false},
		{Name: "Pewarnaan Gram Manual", Code: "GRM-OLD", Description: "Pewarnaan gram metode konvensional", IsActive: false},
		{Name: "Pemeriksaan Sperma Manual", Code: "SPR-OLD", Description: "Analisis sperma metode lama", IsActive: false},
		{Name: "Analisa Gas Darah Konvensional", Code: "AGD-OLD", Description: "Analisa gas darah alat lama", IsActive: false},
		{Name: "Elektrolit Manual", Code: "ELK-OLD", Description: "Pemeriksaan elektrolit metode manual", IsActive: false},
		{Name: "Tes Alergi Konvensional", Code: "ALR-OLD", Description: "Tes alergi versi lama", IsActive: false},
		{Name: "Kultur Darah Manual", Code: "KLD-OLD", Description: "Kultur darah metode konvensional", IsActive: false},
		{Name: "Tes Widal", Code: "WDL-OLD", Description: "Tes Widal metode lama untuk tifus", IsActive: false},
		{Name: "Tes Cepat COVID-19", Code: "CVD-OLD", Description: "Rapid tes antibodi COVID-19", IsActive: false},
	}

	if err := tx.Create(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to seed type test categories: %w", err)
	}

	return categories, nil
}

func seedDeletedTypeTestCategories(tx *gorm.DB) error {
	deletedTime := gorm.DeletedAt{Time: time.Now(), Valid: true}
	var categories []models.TypeTestCategory

	for i := 1; i <= 22; i++ {
		categories = append(categories, models.TypeTestCategory{
			Name:        fmt.Sprintf("Deleted Type Test Category %d", i),
			Code:        fmt.Sprintf("DEL-TTC-%d", i),
			Description: fmt.Sprintf("Description for soft-deleted type test category %d", i),
			IsActive:    false,
			DeletedAt:   deletedTime,
		})
	}

	if err := tx.Create(&categories).Error; err != nil {
		return fmt.Errorf("failed to seed deleted type test categories: %w", err)
	}

	return nil
}
