package database

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedRoomTypes(tx *gorm.DB, count int) ([]models.RoomType, error) {
	activeTypes := []string{
		"VIP", "VVIP", "Class 1", "Class 2", "Class 3",
		"ICU", "NICU", "PICU", "Operation Theater", "Emergency",
		"Isolation", "Maternity",
	}
	inactiveTypes := []string{
		"Old Wing Class 1", "Old Wing Class 2", "Old Wing Class 3", "Temporary Ward A", "Temporary Ward B",
		"Renovation Ward 1", "Renovation Ward 2", "Closed ICU", "Closed NICU", "Closed Operation Theater",
		"Archive Room", "Storage Room",
	}

	roomTypes := make([]models.RoomType, 0, len(activeTypes)+len(inactiveTypes))

	for i, name := range activeTypes {
		code := fmt.Sprintf("RT-ACT-%03d", i+1)
		rt := models.RoomType{
			Name:        name,
			Code:        code,
			Description: fmt.Sprintf("Active Room Type for %s", name),
			IsActive:    true,
		}
		roomTypes = append(roomTypes, rt)
	}

	for i, name := range inactiveTypes {
		code := fmt.Sprintf("RT-INA-%03d", i+1)
		rt := models.RoomType{
			Name:        name,
			Code:        code,
			Description: fmt.Sprintf("Inactive Room Type for %s", name),
			IsActive:    false,
		}
		roomTypes = append(roomTypes, rt)
	}

	if err := tx.Create(&roomTypes).Error; err != nil {
		return nil, fmt.Errorf("failed to seed room types: %w", err)
	}

	return roomTypes, nil
}

func seedDeletedRoomTypes(tx *gorm.DB) error {
	deletedTime := time.Now().AddDate(0, -1, 0) // 1 month ago
	deletedAt := gorm.DeletedAt{Time: deletedTime, Valid: true}

	deletedTypes := []string{
		"Legacy VIP", "Legacy VVIP", "Legacy Class 1", "Legacy Class 2", "Legacy Class 3",
		"Legacy ICU", "Legacy NICU", "Legacy PICU", "Legacy Operation Theater", "Legacy Emergency",
		"Legacy Isolation", "Legacy Maternity", "Demolished Ward A", "Demolished Ward B", "Demolished Ward C",
		"Old Storage 1", "Old Storage 2", "Old Archive 1", "Old Archive 2", "Old Canteen",
		"Old Admin", "Old Lobby",
	}

	roomTypes := make([]models.RoomType, 0, len(deletedTypes))

	for i, name := range deletedTypes {
		code := fmt.Sprintf("RT-DEL-%03d", i+1)
		rt := models.RoomType{
			Name:        name,
			Code:        code,
			Description: fmt.Sprintf("Deleted Room Type for %s", name),
			IsActive:    false,
			DeletedAt:   deletedAt,
		}
		roomTypes = append(roomTypes, rt)
	}

	if err := tx.Create(&roomTypes).Error; err != nil {
		return fmt.Errorf("failed to seed deleted room types: %w", err)
	}

	return nil
}
