package database

import (
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedRooms(tx *gorm.DB, count int, departments []models.Department) ([]models.Room, error) {
	roomTypes := []string{"vip", "class_1", "class_2", "class_3", "icu", "emergency"}

	rooms := make([]models.Room, 0, count)
	for i := 1; i <= count; i++ {
		roomType := roomTypes[(i-1)%len(roomTypes)]
		bedCapacity := (i % 4) + 1
		availableBeds := bedCapacity - (i % bedCapacity)

		var departmentID *uint
		if len(departments) > 0 {
			selected := departments[(i-1)%len(departments)].ID
			departmentID = &selected
		}

		room := models.Room{
			RoomNumber:    fmt.Sprintf("RM-%03d", i),
			RoomType:      roomType,
			DepartmentID:  departmentID,
			BedCapacity:   bedCapacity,
			AvailableBeds: availableBeds,
			PricePerDay:   roomPriceByType(roomType),
			IsActive:      i%8 != 0,
		}
		rooms = append(rooms, room)
	}

	if err := tx.Create(&rooms).Error; err != nil {
		return nil, fmt.Errorf("failed to seed rooms: %w", err)
	}

	return rooms, nil
}

func roomPriceByType(roomType string) float64 {
	switch roomType {
	case "vip":
		return 1500000
	case "class_1":
		return 900000
	case "class_2":
		return 600000
	case "class_3":
		return 350000
	case "icu":
		return 2200000
	case "emergency":
		return 800000
	default:
		return 500000
	}
}
