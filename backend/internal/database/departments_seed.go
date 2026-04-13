package database

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedDepartments(tx *gorm.DB, count int) ([]models.Department, error) {
	seedDepartments := []models.Department{
		{Name: "Kardiologi", Code: "KARDIO", Description: "Layanan jantung dan pembuluh darah", FloorLocation: "Lantai 3"},
		{Name: "Neurologi", Code: "NEURO", Description: "Layanan saraf dan otak", FloorLocation: "Lantai 4"},
		{Name: "Pediatri", Code: "PEDIA", Description: "Layanan kesehatan anak", FloorLocation: "Lantai 2"},
		{Name: "Obstetri dan Ginekologi", Code: "OBGYN", Description: "Layanan kandungan dan kebidanan", FloorLocation: "Lantai 5"},
		{Name: "Orthopedi", Code: "ORTHO", Description: "Layanan tulang dan sendi", FloorLocation: "Lantai 4"},
		{Name: "Radiologi", Code: "RADIO", Description: "Pemeriksaan pencitraan medis", FloorLocation: "Lantai 1"},
		{Name: "Laboratorium Klinik", Code: "LABKLI", Description: "Analisis laboratorium diagnostik", FloorLocation: "Lantai 1"},
		{Name: "Instalasi Gawat Darurat", Code: "IGD", Description: "Layanan emergensi 24 jam", FloorLocation: "Lantai 1"},
		{Name: "Bedah Umum", Code: "BEDAH", Description: "Layanan tindakan bedah umum", FloorLocation: "Lantai 6"},
		{Name: "Penyakit Dalam", Code: "INTER", Description: "Layanan penyakit internal", FloorLocation: "Lantai 3"},
		{Name: "Mata", Code: "MATA", Description: "Layanan oftalmologi", FloorLocation: "Lantai 2"},
		{Name: "THT", Code: "THT", Description: "Layanan telinga hidung tenggorokan", FloorLocation: "Lantai 2"},
		{Name: "Dermatologi", Code: "DERMA", Description: "Layanan kulit dan kelamin", FloorLocation: "Lantai 3"},
		{Name: "Urologi", Code: "URO", Description: "Layanan saluran kemih", FloorLocation: "Lantai 5"},
		{Name: "Onkologi", Code: "ONKO", Description: "Layanan kanker dan kemoterapi", FloorLocation: "Lantai 6"},
		{Name: "Nefrologi", Code: "NEFRO", Description: "Layanan ginjal", FloorLocation: "Lantai 4"},
		{Name: "Pulmonologi", Code: "PULMO", Description: "Layanan paru dan respirasi", FloorLocation: "Lantai 3"},
		{Name: "Geriatri", Code: "GERIA", Description: "Layanan lansia", FloorLocation: "Lantai 2"},
		{Name: "Rehabilitasi Medik", Code: "REHAB", Description: "Layanan fisioterapi dan rehabilitasi", FloorLocation: "Lantai 1"},
		{Name: "Anestesi", Code: "ANEST", Description: "Layanan anestesi dan ICU support", FloorLocation: "Lantai 6"},
		{Name: "Patologi Anatomi", Code: "PATAN", Description: "Pemeriksaan jaringan dan sel", FloorLocation: "Lantai 1"},
		{Name: "Farmasi", Code: "FARM", Description: "Manajemen obat dan dispensing", FloorLocation: "Lantai 1"},
	}

	departments := make([]models.Department, 0, count)
	for i := 0; i < count; i++ {
		if i < len(seedDepartments) {
			departments = append(departments, seedDepartments[i])
			continue
		}

		departments = append(departments, models.Department{
			Name:          fmt.Sprintf("Department %d", i+1),
			Code:          fmt.Sprintf("DEPT%03d", i+1),
			Description:   "Generated seed department",
			FloorLocation: fmt.Sprintf("Lantai %d", (i%6)+1),
		})
	}

	if err := tx.Create(&departments).Error; err != nil {
		return nil, fmt.Errorf("failed to seed departments: %w", err)
	}

	return departments, nil
}

func seedDeletedDepartments(tx *gorm.DB) error {
	deletedTime := time.Now().AddDate(0, -1, 0) // 1 bulan yang lalu
	deletedAt := gorm.DeletedAt{Time: deletedTime, Valid: true}

	departments := make([]models.Department, 0, 12)
	for i := 1; i <= 12; i++ {
		department := models.Department{
			Name:          fmt.Sprintf("Deleted Department %02d", i),
			Code:          fmt.Sprintf("DEL%03d", i),
			Description:   fmt.Sprintf("Departemen yang telah dihapus - #%02d", i),
			FloorLocation: fmt.Sprintf("Lantai %d", (i%6)+1),
			DeletedAt:     deletedAt,
		}
		departments = append(departments, department)
	}

	if err := tx.Create(&departments).Error; err != nil {
		return fmt.Errorf("failed to seed deleted departments: %w", err)
	}

	return nil
}
