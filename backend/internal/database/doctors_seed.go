package database

import (
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedDoctors(tx *gorm.DB, count int, users []models.User, departments []models.Department) ([]models.Doctor, error) {
	doctorUsers := make([]models.User, 0)
	for _, user := range users {
		if user.Role == models.RoleDoctor {
			doctorUsers = append(doctorUsers, user)
		}
	}

	specializations := []string{
		"Cardiology",
		"Neurology",
		"Pediatrics",
		"Obstetrics",
		"Orthopedics",
		"Radiology",
		"Internal Medicine",
		"General Surgery",
		"Pulmonology",
		"Dermatology",
	}

	doctors := make([]models.Doctor, 0, count)
	for i := 1; i <= count; i++ {
		var userID *uint
		if len(doctorUsers) > 0 {
			selected := doctorUsers[(i-1)%len(doctorUsers)].ID
			userID = &selected
		}

		var departmentID *uint
		if len(departments) > 0 {
			selected := departments[(i-1)%len(departments)].ID
			departmentID = &selected
		}

		doctor := models.Doctor{
			UserID:         userID,
			EmployeeID:     fmt.Sprintf("DOC-%04d", i),
			FullName:       fmt.Sprintf("dr. Sample Doctor %02d", i),
			Specialization: specializations[(i-1)%len(specializations)],
			LicenseNumber:  fmt.Sprintf("STR-2026-%05d", i),
			Phone:          fmt.Sprintf("0831100%05d", i),
			Email:          fmt.Sprintf("doctor%02d@sirekam.local", i),
			DepartmentID:   departmentID,
			IsActive:       i%6 != 0,
		}
		doctors = append(doctors, doctor)
	}

	if err := tx.Create(&doctors).Error; err != nil {
		return nil, fmt.Errorf("failed to seed doctors: %w", err)
	}

	return doctors, nil
}
