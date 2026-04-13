package database

import (
	"fmt"
	"time"

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

func seedDeletedPatients(tx *gorm.DB) error {
	genders := []string{"male", "female", "other"}
	bloodTypes := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "O+", "O-"}
	insuranceProviders := []string{"BPJS", "Mandiri Health", "Prudential", "Allianz", "AXA"}

	deletedTime := time.Now().AddDate(0, -1, 0) // 1 bulan yang lalu
	deletedAt := gorm.DeletedAt{Time: deletedTime, Valid: true}

	patients := make([]models.Patient, 0, 12)
	for i := 2001; i <= 2012; i++ {
		patient := models.Patient{
			PatientCode:           fmt.Sprintf("P-DEL-%04d", i-2000),
			FullName:              fmt.Sprintf("Deleted Patient %02d", i-2000),
			DateOfBirth:           fmt.Sprintf("%04d-%02d-%02d", 1980+(i%30), ((i-1)%12)+1, ((i-1)%27)+1),
			Gender:                genders[(i-1)%len(genders)],
			BloodType:             bloodTypes[(i-1)%len(bloodTypes)],
			Phone:                 fmt.Sprintf("0895300%05d", i),
			Email:                 fmt.Sprintf("deleted-patient%02d@sirekam.local", i-2000),
			Address:               fmt.Sprintf("Jl. Lama No. %d, Kota Terhapus", i-2000),
			EmergencyContactName:  fmt.Sprintf("Kontak Darurat Lama %02d", i-2000),
			EmergencyContactPhone: fmt.Sprintf("0819900%05d", i),
			InsuranceNumber:       fmt.Sprintf("INS-DEL-%04d", i-2000),
			InsuranceProvider:     insuranceProviders[(i-1)%len(insuranceProviders)],
			Allergies:             []string{"", "Penicillin", "Seafood", "Dust"}[(i-1)%4],
			DeletedAt:             deletedAt,
		}
		patients = append(patients, patient)
	}

	if err := tx.Create(&patients).Error; err != nil {
		return fmt.Errorf("failed to seed deleted patients: %w", err)
	}

	return nil
}
