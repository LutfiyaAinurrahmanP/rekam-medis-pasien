package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestRoomWithData(id uint, roomNumber string, roomTypeID, deptID uint, capacity, available int, price float64, isActive bool) *models.Room {
	return &models.Room{
		ID:            id,
		RoomNumber:    roomNumber,
		RoomTypeID:    PtrUint(roomTypeID),
		DepartmentID:  PtrUint(deptID),
		BedCapacity:   capacity,
		AvailableBeds: available,
		PricePerDay:   price,
		IsActive:      isActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func NewTestRoomList(count int) []models.Room {
	rooms := make([]models.Room, count)
	for i := 0; i < count; i++ {
		rooms[i] = *NewTestRoomWithData(
			uint(i+1),
			"R00"+string(rune(i+1)),
			uint(1),
			uint(1),
			4,
			4,
			500000.0,
			true,
		)
	}
	return rooms
}

func NewRoomPaginationQuery(page, pageSize int) *dto.RoomPaginationQuery {
	return &dto.RoomPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewCreateRoomRequest(roomNumber string, roomTypeID, deptID uint, capacity, available int, price float64, isActive bool) *dto.CreateRoomRequest {
	return &dto.CreateRoomRequest{
		RoomNumber:    roomNumber,
		RoomTypeID:    PtrUint(roomTypeID),
		DepartmentID:  PtrUint(deptID),
		BedCapacity:   capacity,
		AvailableBeds: PtrInt(available),
		PricePerDay:   PtrFloat64(price),
		IsActive:      PtrBool(isActive),
	}
}

func NewUpdateRoomRequest(roomNumber string, roomTypeID, deptID uint, capacity int, price float64, isActive bool) *dto.UpdateRoomRequest {
	return &dto.UpdateRoomRequest{
		RoomNumber:    PtrString(roomNumber),
		RoomTypeID:    PtrUint(roomTypeID),
		DepartmentID:  PtrUint(deptID),
		BedCapacity:   PtrInt(capacity),
		PricePerDay:   PtrFloat64(price),
		IsActive:      PtrBool(isActive),
	}
}

func NewTestRoomResponseWithData(id uint, roomNumber string, roomTypeID, deptID uint, capacity, available int, price float64, isActive bool) *dto.RoomResponse {
	return &dto.RoomResponse{
		ID:            id,
		RoomNumber:    roomNumber,
		RoomTypeID:    PtrUint(roomTypeID),
		DepartmentID:  PtrUint(deptID),
		BedCapacity:   capacity,
		AvailableBeds: available,
		PricePerDay:   price,
		IsActive:      isActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
