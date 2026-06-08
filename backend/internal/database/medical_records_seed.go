package database

import (
	"fmt"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func SeedMedicalRecords(db *gorm.DB) error {
	var count int64
	db.Unscoped().Model(&models.MedicalRecord{}).Count(&count)
	if count > 0 {
		log.Println("MedicalRecords table is not empty, skipping seed")
		return nil
	}

	log.Println("Seeding MedicalRecords data...")

	var records []models.MedicalRecord
	now := time.Now()

	// 1. Generate active medical records
	for i := 1; i <= 22; i++ {
		patientID := uint((i % 5) + 1)
		doctorID := uint((i % 5) + 1)
		appID := uint(i)

		isFinalized := i%3 != 0 // 2/3 are finalized
		status := "draft"
		if isFinalized {
			status = "finalized"
		}

		records = append(records, models.MedicalRecord{
			PatientID:           patientID,
			DoctorID:            doctorID,
			AppointmentID:       &appID,
			VisitDate:           now.AddDate(0, 0, -(i % 10)).Format("2006-01-02"),
			ChiefComplaint:      fmt.Sprintf("Keluhan utama pasien %d - Kunjungan %d", patientID, i),
			HistoryOfIllness:    "Pasien mengeluh sakit sejak 3 hari yang lalu.",
			PhysicalExamination: "Dalam batas normal",
			Diagnosis:           fmt.Sprintf("Diagnosis Penyakit X untuk Pasien %d", patientID),
			TreatmentPlan:       "Edukasi dan modifikasi gaya hidup, obat diteruskan",
			Notes:               "Kontrol 2 minggu lagi",
			Status:              status,
			CreatedAt:           now,
			UpdatedAt:           now,
		})
	}

	// 2. Generate deleted medical records
	for i := 23; i <= 25; i++ {
		patientID := uint((i % 5) + 1)
		doctorID := uint((i % 5) + 1)

		records = append(records, models.MedicalRecord{
			PatientID:      patientID,
			DoctorID:       doctorID,
			VisitDate:      now.AddDate(0, 0, -(i % 10)).Format("2006-01-02"),
			ChiefComplaint: "Data Salah Input",
			Diagnosis:      "Salah Diagnosis",
			TreatmentPlan:  "Salah Obat",
			Status:         "draft",
			CreatedAt:      now,
			UpdatedAt:      now,
			DeletedAt:      gorm.DeletedAt{Time: now, Valid: true},
		})
	}

	if err := db.Create(&records).Error; err != nil {
		log.Printf("Error seeding MedicalRecords: %v\n", err)
		return err
	}

	log.Println("MedicalRecords seeded successfully")
	return nil
}
