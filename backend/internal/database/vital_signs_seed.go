package database

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func SeedVitalSigns(db *gorm.DB) error {
	var medicalRecords []models.MedicalRecord
	if err := db.Unscoped().Find(&medicalRecords).Error; err != nil {
		return fmt.Errorf("failed to fetch medical records for vital signs seed: %w", err)
	}

	if len(medicalRecords) == 0 {
		log.Println("⚠️  No medical records found, skipping vital signs seed")
		return nil
	}

	var activeMRs []models.MedicalRecord
	var deletedMRs []models.MedicalRecord

	for _, mr := range medicalRecords {
		if mr.DeletedAt.Valid {
			deletedMRs = append(deletedMRs, mr)
		} else {
			activeMRs = append(activeMRs, mr)
		}
	}

	var activeVitalSigns []models.VitalSign
	var deletedVitalSigns []models.VitalSign

	rng := rand.New(rand.NewSource(42))
	now := time.Now()

	// target = len - 10
	targetActive := len(activeMRs) - 10
	if targetActive < 0 {
		targetActive = 0
	}
	targetDeleted := len(deletedMRs) - 10
	if targetDeleted < 0 {
		targetDeleted = 0
	}

	// Create Active Vital Signs
	for i, mr := range activeMRs {
		if i >= targetActive {
			break
		}
		activeVitalSigns = append(activeVitalSigns, generateVitalSign(rng, mr, now, false))
	}

	// Create Deleted Vital Signs
	for i, mr := range deletedMRs {
		if i >= targetDeleted {
			break
		}
		deletedVitalSigns = append(deletedVitalSigns, generateVitalSign(rng, mr, now, true))
	}

	if len(activeVitalSigns) > 0 {
		if err := db.Create(&activeVitalSigns).Error; err != nil {
			return fmt.Errorf("failed to seed active vital signs: %w", err)
		}
		log.Printf("✅ Seeded %d active vital signs", len(activeVitalSigns))
	}

	if len(deletedVitalSigns) > 0 {
		if err := db.Create(&deletedVitalSigns).Error; err != nil {
			return fmt.Errorf("failed to seed deleted vital signs: %w", err)
		}
		log.Printf("✅ Seeded %d deleted vital signs", len(deletedVitalSigns))
	}

	return nil
}

func generateVitalSign(rng *rand.Rand, mr models.MedicalRecord, now time.Time, isDeleted bool) models.VitalSign {
	weight := 45.0 + rng.Float64()*50.0     // 45kg - 95kg
	height := 150 + rng.Intn(40)            // 150cm - 190cm
	systolic := 90 + rng.Intn(50)           // 90 - 140
	diastolic := 60 + rng.Intn(30)          // 60 - 90
	heartRate := 60 + rng.Intn(40)          // 60 - 100
	temperature := 36.0 + rng.Float64()*3.0 // 36.0 - 39.0

	heightM := float64(height) / 100.0
	bmi := weight / (heightM * heightM)

	visitDate, err := time.Parse("2006-01-02", mr.VisitDate)
	if err != nil {
		visitDate = now
	}

	vitalSign := models.VitalSign{
		MedicalRecordID: mr.ID,
		MeasurementDate: visitDate,
		MeasurementTime: "08:30:00",
		SystolicBP:      &systolic,
		DiastolicBP:     &diastolic,
		HeartRate:       &heartRate,
		BodyTemperature: &temperature,
		WeightKg:        &weight,
		HeightCm:        &height,
		BMI:             &bmi,
		Notes:           fmt.Sprintf("Vital signs note for medical record %d", mr.ID),
	}

	if isDeleted {
		deletedTime := now.Add(-24 * time.Hour)
		vitalSign.DeletedAt = gorm.DeletedAt{Time: deletedTime, Valid: true}
	}

	return vitalSign
}
