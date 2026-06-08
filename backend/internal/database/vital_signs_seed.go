package database

import (
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func SeedVitalSigns(db *gorm.DB) error {
	var medicalRecords []models.MedicalRecord
	if err := db.Find(&medicalRecords).Error; err != nil {
		return err
	}

	if len(medicalRecords) == 0 {
		log.Println("⚠️  No medical records found, skipping vital signs seeding")
		return nil
	}

	log.Println("🌱 Seeding vital signs...")

	vitalSigns := make([]models.VitalSign, 0)
	now := time.Now()

	for i, record := range medicalRecords {
		// Create vital signs for every other medical record
		if i%2 == 0 {
			sys := 110 + (i % 30)
			dias := 70 + (i % 20)
			hr := 60 + (i % 40)
			temp := 36.5 + float64(i%10)/10.0
			rr := 16 + (i % 8)
			o2 := 95.0 + float64(i%6)
			weight := 60.0 + float64(i%30)
			height := 160 + (i % 30)
			bmi := weight / ((float64(height) / 100) * (float64(height) / 100))

			visitDate, err := time.Parse("2006-01-02", record.VisitDate)
			if err != nil {
				visitDate = now
			}

			vitalSigns = append(vitalSigns, models.VitalSign{
				MedicalRecordID:  record.ID,
				MeasurementDate:  visitDate,
				MeasurementTime:  "08:30:00",
				SystolicBP:       &sys,
				DiastolicBP:      &dias,
				HeartRate:        &hr,
				BodyTemperature:  &temp,
				RespiratoryRate:  &rr,
				OxygenSaturation: &o2,
				WeightKg:         &weight,
				HeightCm:         &height,
				BMI:              &bmi,
				Notes:            "Normal measurement",
				CreatedAt:        now,
				UpdatedAt:        now,
			})
		}
	}

	if len(vitalSigns) > 0 {
		if err := db.Create(&vitalSigns).Error; err != nil {
			log.Printf("❌ Failed to seed vital signs: %v", err)
			return err
		}
	}

	return nil
}
