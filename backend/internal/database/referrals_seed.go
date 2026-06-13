package database

import (
	"fmt"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func SeedReferrals(db *gorm.DB) error {
	var count int64
	db.Unscoped().Model(&models.Referral{}).Count(&count)
	if count > 0 {
		log.Println("Referrals table is not empty, skipping seed")
		return nil
	}

	log.Println("Seeding Referrals data...")

	var referrals []models.Referral
	now := time.Now()

	// 1. Generate 22 active referrals
	for i := 1; i <= 22; i++ {
		patientID := uint((i % 5) + 1)
		medicalRecordID := uint((i % 5) + 1)
		refDoctorID := uint((i % 5) + 1)

		referralType := "internal"
		var deptID *uint
		var toDocID *uint
		var facility string

		if i%3 == 0 {
			referralType = "external"
			facility = fmt.Sprintf("RSUD Rujukan Eksternal %d", i)
		} else {
			dID := uint((i % 4) + 1)
			docID := uint((i % 3) + 1)
			deptID = &dID
			toDocID = &docID
		}

		status := "pending"
		if i%4 == 0 {
			status = "accepted"
		} else if i%5 == 0 {
			status = "completed"
		}

		priority := "routine"
		if i%2 == 0 {
			priority = "urgent"
		}

		referral := models.Referral{
			ReferralNumber:         fmt.Sprintf("REF-2024-%06d", i),
			PatientID:              patientID,
			MedicalRecordID:        medicalRecordID,
			ReferringDoctorID:      refDoctorID,
			ReferralDate:           now.AddDate(0, 0, -(i % 10)).Format("2006-01-02"),
			ReferralType:           referralType,
			ReferredToDepartmentID: deptID,
			ReferredToDoctorID:     toDocID,
			ReferredToFacility:     facility,
			Reason:                 fmt.Sprintf("Alasan rujukan untuk pasien %d - Kasus %d", patientID, i),
			Diagnosis:              fmt.Sprintf("Diagnosis kerja %d", i),
			Priority:               priority,
			Status:                 status,
			Notes:                  "Harap ditindaklanjuti segera",
			CreatedAt:              now,
			UpdatedAt:              now,
		}

		if status == "accepted" || status == "completed" {
			acceptedTime := now.Add(-time.Hour * 24)
			referral.AcceptedAt = &acceptedTime
		}
		if status == "completed" {
			completedTime := now.Add(-time.Hour * 12)
			referral.CompletedAt = &completedTime
		}

		referrals = append(referrals, referral)
	}

	// 2. Generate 22 deleted referrals
	for i := 23; i <= 44; i++ {
		patientID := uint((i % 5) + 1)
		medicalRecordID := uint((i % 5) + 1)
		refDoctorID := uint((i % 5) + 1)
		dID := uint((i % 4) + 1)

		referrals = append(referrals, models.Referral{
			ReferralNumber:         fmt.Sprintf("REF-DEL-%06d", i),
			PatientID:              patientID,
			MedicalRecordID:        medicalRecordID,
			ReferringDoctorID:      refDoctorID,
			ReferralDate:           now.AddDate(0, 0, -(i % 10)).Format("2006-01-02"),
			ReferralType:           "internal",
			ReferredToDepartmentID: &dID,
			Reason:                 "Salah input pasien rujukan",
			Priority:               "routine",
			Status:                 "cancelled",
			CancellationReason:     "Salah input rujukan",
			CreatedAt:              now,
			UpdatedAt:              now,
			DeletedAt:              gorm.DeletedAt{Time: now, Valid: true},
		})
	}

	if err := db.Create(&referrals).Error; err != nil {
		log.Printf("Error seeding Referrals: %v\n", err)
		return err
	}

	log.Println("Referrals seeded successfully")
	return nil
}
