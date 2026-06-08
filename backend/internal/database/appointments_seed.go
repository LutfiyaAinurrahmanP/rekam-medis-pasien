package database

import (
	"fmt"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func SeedAppointments(db *gorm.DB) error {
	var count int64
	db.Unscoped().Model(&models.Appointment{}).Count(&count)
	if count > 0 {
		log.Println("Appointments table is not empty, skipping seed")
		return nil
	}

	log.Println("Seeding Appointments data...")

	statuses := []string{"scheduled", "confirmed", "in_progress", "completed", "cancelled", "no_show"}

	var appointments []models.Appointment
	now := time.Now()

	// 1. Generate 22 active appointments
	for i := 1; i <= 22; i++ {
		patientID := uint((i % 5) + 1)
		doctorID := uint((i % 5) + 1)
		status := statuses[i%len(statuses)]
		
		// Spread dates across a few days
		appDate := now.AddDate(0, 0, i%10).Format("2006-01-02")
		appTime := fmt.Sprintf("%02d:00:00", 8+(i%8))

		appointments = append(appointments, models.Appointment{
			PatientID:       patientID,
			DoctorID:        doctorID,
			AppointmentDate: appDate,
			AppointmentTime: appTime,
			DurationMinutes: 30,
			Status:          status,
			Reason:          fmt.Sprintf("Keluhan Pasien %d - Sesi %d", patientID, i),
			Notes:           "Catatan dokter/receptionist",
			CreatedAt:       now,
			UpdatedAt:       now,
		})
	}

	// 2. Generate 22 soft-deleted appointments
	for i := 23; i <= 44; i++ {
		patientID := uint((i % 5) + 1)
		doctorID := uint((i % 5) + 1)
		status := statuses[i%len(statuses)]
		
		appDate := now.AddDate(0, -1, i%10).Format("2006-01-02") // in the past
		appTime := fmt.Sprintf("%02d:00:00", 8+(i%8))

		appointments = append(appointments, models.Appointment{
			PatientID:       patientID,
			DoctorID:        doctorID,
			AppointmentDate: appDate,
			AppointmentTime: appTime,
			DurationMinutes: 30,
			Status:          status,
			Reason:          fmt.Sprintf("Keluhan Pasien %d - Sesi %d (Deleted)", patientID, i),
			Notes:           "Catatan dihapus",
			CreatedAt:       now.Add(-48 * time.Hour),
			UpdatedAt:       now.Add(-24 * time.Hour),
			DeletedAt:       gorm.DeletedAt{Time: now, Valid: true},
		})
	}

	if err := db.Unscoped().Create(&appointments).Error; err != nil {
		log.Printf("Failed to seed appointments: %v\n", err)
		return err
	}

	log.Println("Successfully seeded Appointments data")
	return nil
}
