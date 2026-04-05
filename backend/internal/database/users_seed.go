package database

import (
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"gorm.io/gorm"
)

func seedUsers(tx *gorm.DB, count int) ([]models.User, error) {
	hashedPassword, err := utils.HashPassword("Password123!")
	if err != nil {
		return nil, fmt.Errorf("failed to hash default password: %w", err)
	}

	roles := []string{
		models.RolePatient,
		models.RoleDoctor,
		models.RoleReceptionist,
		models.RoleAdmin,
		models.RoleSuperAdmin,
	}

	users := make([]models.User, 0, count)
	for i := 1; i <= count; i++ {
		user := models.User{
			Username: fmt.Sprintf("user_%02d", i),
			Email:    fmt.Sprintf("user%02d@sirekam.local", i),
			Phone:    fmt.Sprintf("0812300%05d", i),
			Password: hashedPassword,
			Role:     roles[(i-1)%len(roles)],
			IsActive: i%7 != 0,
		}
		users = append(users, user)
	}

	if err := tx.Create(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to seed users: %w", err)
	}

	return users, nil
}
