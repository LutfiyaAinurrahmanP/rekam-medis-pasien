package database

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedDoctors(tx *gorm.DB, count int, users []models.User, departments []models.Department, specializations []models.DoctorSpecialization) ([]models.Doctor, error) {
	doctorUsers := make([]models.User, 0)
	for _, user := range users {
		if user.Role == models.RoleDoctor {
			doctorUsers = append(doctorUsers, user)
		}
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

		var specializationID uint
		if len(specializations) > 0 {
			specializationID = specializations[(i-1)%len(specializations)].ID
		}

		doctor := models.Doctor{
			UserID:           userID,
			EmployeeID:       fmt.Sprintf("DOC-%04d", i),
			FullName:         fmt.Sprintf("dr. Sample Doctor %02d", i),
			SpecializationID: specializationID,
			LicenseNumber:    fmt.Sprintf("STR-2026-%05d", i),
			Phone:            fmt.Sprintf("0831100%05d", i),
			Email:            fmt.Sprintf("doctor%02d@sirekam.local", i),
			DepartmentID:     departmentID,
			IsActive:         i%6 != 0,
		}
		doctors = append(doctors, doctor)
	}

	if err := tx.Create(&doctors).Error; err != nil {
		return nil, fmt.Errorf("failed to seed doctors: %w", err)
	}

	return doctors, nil
}

func seedDeletedDoctors(tx *gorm.DB) error {
	var specializations []models.DoctorSpecialization
	if err := tx.Where("deleted_at IS NULL").Find(&specializations).Error; err != nil {
		return fmt.Errorf("failed to fetch specializations: %w", err)
	}
	if len(specializations) == 0 {
		return fmt.Errorf("no active doctor specializations found for doctor assignment")
	}

	// Get active departments for assignment
	var departments []models.Department
	if err := tx.Where("deleted_at IS NULL").Find(&departments).Error; err != nil {
		return fmt.Errorf("failed to fetch departments: %w", err)
	}

	if len(departments) == 0 {
		return fmt.Errorf("no active departments found for doctor assignment")
	}

	deletedTime := time.Now().AddDate(0, -1, 0) // 1 bulan yang lalu
	deletedAt := gorm.DeletedAt{Time: deletedTime, Valid: true}

	doctors := make([]models.Doctor, 0, 12)
	for i := 3001; i <= 3012; i++ {
		selected := departments[(i-1)%len(departments)].ID
		selectedSpec := specializations[(i-1)%len(specializations)].ID

		doctor := models.Doctor{
			EmployeeID:       fmt.Sprintf("DOC-DEL-%04d", i-3000),
			FullName:         fmt.Sprintf("dr. Deleted Doctor %02d", i-3000),
			SpecializationID: selectedSpec,
			LicenseNumber:    fmt.Sprintf("STR-DEL-%05d", i),
			Phone:            fmt.Sprintf("0887100%05d", i),
			Email:            fmt.Sprintf("deleted-doctor%02d@sirekam.local", i-3000),
			DepartmentID:     &selected,
			IsActive:         false,
			DeletedAt:        deletedAt,
		}
		doctors = append(doctors, doctor)
	}

	if err := tx.Create(&doctors).Error; err != nil {
		return fmt.Errorf("failed to seed deleted doctors: %w", err)
	}

	return nil
}
