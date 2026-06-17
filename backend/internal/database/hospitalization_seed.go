package database

import (
	"fmt"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func SeedHospitalizations(db *gorm.DB) error {
	var count int64
	db.Unscoped().Model(&models.Hospitalization{}).Count(&count)
	if count > 0 {
		return nil // Already seeded
	}

	log.Println("🌱 Seeding hospitalizations...")

	// Fetch dependencies to use valid IDs
	var patients []models.Patient
	var doctors []models.Doctor
	var rooms []models.Room

	db.Find(&patients)
	db.Find(&doctors)
	db.Find(&rooms)

	if len(patients) == 0 || len(doctors) == 0 || len(rooms) == 0 {
		return fmt.Errorf("cannot seed hospitalizations: missing patients, doctors, or rooms")
	}

	var hospitalizations []models.Hospitalization
	now := time.Now()

	// 22 active records (deleted_at is nil)
	for i := 1; i <= 22; i++ {
		patientIdx := (i - 1) % len(patients)
		doctorIdx := (i - 1) % len(doctors)
		roomIdx := (i - 1) % len(rooms)

		h := models.Hospitalization{
			PatientID:          patients[patientIdx].ID,
			DoctorID:           doctors[doctorIdx].ID,
			RoomID:             rooms[roomIdx].ID,
			AdmissionDate:      now.AddDate(0, 0, -i).Format("2006-01-02"),
			AdmissionTime:      "08:00:00",
			ReasonForAdmission: fmt.Sprintf("Alasan medis untuk pasien %d (Active %d)", patients[patientIdx].ID, i),
			Status:             "admitted",
			CreatedAt:          now.AddDate(0, 0, -i),
			UpdatedAt:          now.AddDate(0, 0, -i),
		}

		// Some discharged/transferred records
		if i%3 == 0 {
			h.Status = "discharged"
			dischargeDate := now.AddDate(0, 0, -i+2).Format("2006-01-02")
			dischargeTime := "10:00:00"
			h.DischargeDate = &dischargeDate
			h.DischargeTime = &dischargeTime
			h.Notes = "Pasien membaik dan diperbolehkan pulang"
		} else if i%4 == 0 {
			h.Status = "transferred"
			h.Notes = "Pasien ditransfer ke ruang operasi"
		}

		hospitalizations = append(hospitalizations, h)
	}

	// 22 deleted records (deleted_at is not nil)
	for i := 1; i <= 22; i++ {
		patientIdx := (i - 1) % len(patients)
		doctorIdx := (i - 1) % len(doctors)
		roomIdx := (i - 1) % len(rooms)

		h := models.Hospitalization{
			PatientID:          patients[patientIdx].ID,
			DoctorID:           doctors[doctorIdx].ID,
			RoomID:             rooms[roomIdx].ID,
			AdmissionDate:      now.AddDate(0, 0, -(i + 30)).Format("2006-01-02"),
			AdmissionTime:      "09:30:00",
			ReasonForAdmission: fmt.Sprintf("Alasan medis untuk pasien %d (Deleted %d)", patients[patientIdx].ID, i),
			Status:             "cancelled",
			Notes:              "Dibatalkan karena kesalahan input",
			CreatedAt:          now.AddDate(0, 0, -(i + 30)),
			UpdatedAt:          now.AddDate(0, 0, -(i + 30)),
			DeletedAt:          gorm.DeletedAt{Time: now.AddDate(0, 0, -(i + 29)), Valid: true},
		}

		hospitalizations = append(hospitalizations, h)
	}

	if err := db.Unscoped().Create(&hospitalizations).Error; err != nil {
		return fmt.Errorf("failed to insert hospitalizations: %w", err)
	}

	return nil
}
