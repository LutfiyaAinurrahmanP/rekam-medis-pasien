package database

import (
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedPatients(tx *gorm.DB, count int, users []models.User) ([]models.Patient, error) {
	patientUsers := make([]models.User, 0)
	for _, user := range users {
		if user.Role == models.RolePatient {
			patientUsers = append(patientUsers, user)
		}
	}

	genders := []string{"male", "female", "other"}
	bloodTypes := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
	insuranceProviders := []string{"BPJS", "Mandiri Health", "Prudential", "Allianz", "AXA"}

	patients := make([]models.Patient, 0, count)
	for i := 1; i <= count; i++ {
		var userID *uint
		if len(patientUsers) > 0 {
			selected := patientUsers[(i-1)%len(patientUsers)].ID
			userID = &selected
		}

		patient := models.Patient{
			UserID:                userID,
			PatientCode:           fmt.Sprintf("P-2026-%04d", i),
			FullName:              fmt.Sprintf("Patient Sample %02d", i),
			DateOfBirth:           fmt.Sprintf("%04d-%02d-%02d", 1980+(i%30), ((i-1)%12)+1, ((i-1)%27)+1),
			Gender:                genders[(i-1)%len(genders)],
			BloodType:             bloodTypes[(i-1)%len(bloodTypes)],
			Phone:                 fmt.Sprintf("0821100%05d", i),
			Email:                 fmt.Sprintf("patient%02d@sirekam.local", i),
			Address:               fmt.Sprintf("Jl. Sehat No. %d, Kota Medis", i),
			EmergencyContactName:  fmt.Sprintf("Kontak Darurat %02d", i),
			EmergencyContactPhone: fmt.Sprintf("0819900%05d", i),
			InsuranceNumber:       fmt.Sprintf("INS-%06d", i),
			InsuranceProvider:     insuranceProviders[(i-1)%len(insuranceProviders)],
			Allergies:             []string{"", "Penicillin", "Seafood", "Dust"}[(i-1)%4],
		}
		patients = append(patients, patient)
	}

	if err := tx.Create(&patients).Error; err != nil {
		return nil, fmt.Errorf("failed to seed patients: %w", err)
	}

	return patients, nil
}
