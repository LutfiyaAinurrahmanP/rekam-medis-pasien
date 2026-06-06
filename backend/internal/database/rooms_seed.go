package database

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedRooms(tx *gorm.DB, count int, departments []models.Department) ([]models.Room, error) {
	var roomTypes []models.RoomType
	if err := tx.Where("deleted_at IS NULL").Find(&roomTypes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch room types: %w", err)
	}
	if len(roomTypes) == 0 {
		return nil, fmt.Errorf("no room types available to seed rooms")
	}

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
			RoomTypeID:    &roomType.ID,
			DepartmentID:  departmentID,
			BedCapacity:   bedCapacity,
			AvailableBeds: availableBeds,
			PricePerDay:   roomPriceByType(roomType.Code),
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

func seedDeletedRooms(tx *gorm.DB) error {
	var roomTypes []models.RoomType
	if err := tx.Where("deleted_at IS NULL").Find(&roomTypes).Error; err != nil {
		return fmt.Errorf("failed to fetch room types: %w", err)
	}
	if len(roomTypes) == 0 {
		return fmt.Errorf("no room types available to seed deleted rooms")
	}

	// Get active departments for assignment
	var departments []models.Department
	if err := tx.Where("deleted_at IS NULL").Find(&departments).Error; err != nil {
		return fmt.Errorf("failed to fetch departments: %w", err)
	}

	deletedTime := time.Now().AddDate(0, -1, 0) // 1 bulan yang lalu
	deletedAt := gorm.DeletedAt{Time: deletedTime, Valid: true}

	rooms := make([]models.Room, 0, 12)
	for i := 4001; i <= 4012; i++ {
		roomType := roomTypes[(i-1)%len(roomTypes)]
		bedCapacity := (i % 4) + 1
		availableBeds := 0

		var departmentID *uint
		if len(departments) > 0 {
			selected := departments[(i-1)%len(departments)].ID
			departmentID = &selected
		}

		room := models.Room{
			RoomNumber:    fmt.Sprintf("RM-DEL-%03d", i-4000),
			RoomTypeID:    &roomType.ID,
			DepartmentID:  departmentID,
			BedCapacity:   bedCapacity,
			AvailableBeds: availableBeds,
			PricePerDay:   roomPriceByType(roomType.Code),
			IsActive:      false,
			DeletedAt:     deletedAt,
		}
		rooms = append(rooms, room)
	}

	if err := tx.Create(&rooms).Error; err != nil {
		return fmt.Errorf("failed to seed deleted rooms: %w", err)
	}

	return nil
}
