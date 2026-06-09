package database

import (
	"fmt"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func SeedPrescriptions(db *gorm.DB) error {
	var count int64
	db.Unscoped().Model(&models.Prescription{}).Count(&count)
	if count > 0 {
		return nil // Already seeded
	}

	log.Println("🌱 Seeding prescriptions and items...")

	var medicalRecords []models.MedicalRecord
	var doctors []models.Doctor
	var medicines []models.Medicine

	db.Find(&medicalRecords)
	db.Find(&doctors)
	db.Find(&medicines)

	if len(medicalRecords) == 0 || len(doctors) == 0 || len(medicines) == 0 {
		return fmt.Errorf("cannot seed prescriptions: missing medical records, doctors, or medicines")
	}

	var prescriptions []models.Prescription
	now := time.Now()

	// 22 active records (deleted_at is nil)
	for i := 1; i <= 22; i++ {
		recordIdx := (i - 1) % len(medicalRecords)
		doctorIdx := (i - 1) % len(doctors)
		medicineIdx := (i - 1) % len(medicines)

		status := "pending"
		if i%2 == 0 {
			status = "dispensed"
		}

		p := models.Prescription{
			MedicalRecordID:  medicalRecords[recordIdx].ID,
			DoctorID:         doctors[doctorIdx].ID,
			PrescriptionDate: now.AddDate(0, 0, -i).Format("2006-01-02"),
			Notes:            fmt.Sprintf("Catatan resep untuk record %d", i),
			Status:           status,
			CreatedAt:        now.AddDate(0, 0, -i),
			UpdatedAt:        now.AddDate(0, 0, -i),
			Items: []models.PrescriptionItem{
				{
					MedicineID:   medicines[medicineIdx].ID,
					Dosage:       "1 tablet",
					Frequency:    "3x sehari",
					DurationDays: 3 + (i % 5),
					Quantity:     (3 + (i % 5)) * 3,
					Instructions: "Sesudah makan",
					CreatedAt:    now.AddDate(0, 0, -i),
					UpdatedAt:    now.AddDate(0, 0, -i),
				},
			},
		}

		prescriptions = append(prescriptions, p)
	}

	// 22 deleted records (deleted_at is not nil)
	for i := 1; i <= 22; i++ {
		recordIdx := (i - 1) % len(medicalRecords)
		doctorIdx := (i - 1) % len(doctors)
		medicineIdx := (i - 1) % len(medicines)

		p := models.Prescription{
			MedicalRecordID:  medicalRecords[recordIdx].ID,
			DoctorID:         doctors[doctorIdx].ID,
			PrescriptionDate: now.AddDate(0, 0, -(i + 30)).Format("2006-01-02"),
			Notes:            fmt.Sprintf("Catatan resep dibatalkan untuk record %d", i),
			Status:           "cancelled",
			CreatedAt:        now.AddDate(0, 0, -(i + 30)),
			UpdatedAt:        now.AddDate(0, 0, -(i + 30)),
			DeletedAt:        gorm.DeletedAt{Time: now.AddDate(0, 0, -(i + 29)), Valid: true},
			Items: []models.PrescriptionItem{
				{
					MedicineID:   medicines[medicineIdx].ID,
					Dosage:       "1 kapsul",
					Frequency:    "2x sehari",
					DurationDays: 5,
					Quantity:     10,
					Instructions: "Sebelum makan",
					CreatedAt:    now.AddDate(0, 0, -(i + 30)),
					UpdatedAt:    now.AddDate(0, 0, -(i + 30)),
					DeletedAt:    gorm.DeletedAt{Time: now.AddDate(0, 0, -(i + 29)), Valid: true},
				},
			},
		}

		prescriptions = append(prescriptions, p)
	}

	if err := db.Unscoped().Create(&prescriptions).Error; err != nil {
		return fmt.Errorf("failed to insert prescriptions: %w", err)
	}

	return nil
}
