package database

import (
	"fmt"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func SeedLabTests(db *gorm.DB) error {
	var count int64
	db.Unscoped().Model(&models.LabTest{}).Count(&count)
	if count > 0 {
		return nil // Already seeded
	}

	log.Println("🌱 Seeding lab tests...")

	// Fetch dependencies to use valid IDs
	var medicalRecords []models.MedicalRecord
	var testTypes []models.TypeTest
	var doctors []models.Doctor

	db.Find(&medicalRecords)
	db.Find(&testTypes)
	db.Find(&doctors)

	if len(medicalRecords) == 0 || len(testTypes) == 0 || len(doctors) == 0 {
		return fmt.Errorf("cannot seed lab tests: missing medical records, test types, or doctors")
	}

	var labTests []models.LabTest
	now := time.Now()

	// 22 active records (deleted_at is nil)
	for i := 1; i <= 22; i++ {
		recordIdx := (i - 1) % len(medicalRecords)
		typeIdx := (i - 1) % len(testTypes)
		doctorIdx := (i - 1) % len(doctors)

		lt := models.LabTest{
			MedicalRecordID:   medicalRecords[recordIdx].ID,
			TestTypeID:        testTypes[typeIdx].ID,
			OrderedByDoctorID: doctors[doctorIdx].ID,
			OrderDate:         now.AddDate(0, 0, -i).Format("2006-01-02"),
			Status:            "ordered",
			CreatedAt:         now.AddDate(0, 0, -i),
			UpdatedAt:         now.AddDate(0, 0, -i),
		}

		// Some have transitioned to other statuses
		if i%3 == 0 {
			lt.Status = "completed"
			sampleDate := now.AddDate(0, 0, -i+1).Format("2006-01-02")
			startDate := now.AddDate(0, 0, -i+1).Format("2006-01-02")
			resultDate := now.AddDate(0, 0, -i+2).Format("2006-01-02")
			lt.SampleCollectionDate = &sampleDate
			lt.TestStartDate = &startDate
			lt.ResultDate = &resultDate
			notes := "Hasil normal, tidak ada indikasi penyakit kritis"
			resultVal := "120"
			lt.Notes = notes
			lt.ResultValue = &resultVal
		} else if i%4 == 0 {
			lt.Status = "in_progress"
			sampleDate := now.AddDate(0, 0, -i+1).Format("2006-01-02")
			startDate := now.AddDate(0, 0, -i+1).Format("2006-01-02")
			lt.SampleCollectionDate = &sampleDate
			lt.TestStartDate = &startDate
		} else if i%5 == 0 {
			lt.Status = "sample_collected"
			sampleDate := now.AddDate(0, 0, -i+1).Format("2006-01-02")
			lt.SampleCollectionDate = &sampleDate
		}

		labTests = append(labTests, lt)
	}

	// 22 deleted records (deleted_at is not nil)
	for i := 1; i <= 22; i++ {
		recordIdx := (i - 1) % len(medicalRecords)
		typeIdx := (i - 1) % len(testTypes)
		doctorIdx := (i - 1) % len(doctors)

		lt := models.LabTest{
			MedicalRecordID:   medicalRecords[recordIdx].ID,
			TestTypeID:        testTypes[typeIdx].ID,
			OrderedByDoctorID: doctors[doctorIdx].ID,
			OrderDate:         now.AddDate(0, 0, -(i + 30)).Format("2006-01-02"),
			Status:            "cancelled",
			CreatedAt:         now.AddDate(0, 0, -(i + 30)),
			UpdatedAt:         now.AddDate(0, 0, -(i + 30)),
			DeletedAt:         gorm.DeletedAt{Time: now.AddDate(0, 0, -(i + 29)), Valid: true},
		}

		labTests = append(labTests, lt)
	}

	if err := db.Unscoped().Create(&labTests).Error; err != nil {
		return fmt.Errorf("failed to insert lab tests: %w", err)
	}

	return nil
}
