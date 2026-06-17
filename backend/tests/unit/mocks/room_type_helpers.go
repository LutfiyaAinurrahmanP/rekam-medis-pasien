package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestRoomTypeWithData(id uint, name, code, description string, isActive bool) *models.RoomType {
	return &models.RoomType{
		ID:          id,
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewTestRoomTypeList(count int) []models.RoomType {
	roomTypes := make([]models.RoomType, count)
	for i := 0; i < count; i++ {
		roomTypes[i] = *NewTestRoomTypeWithData(
			uint(i+1),
			"Room Type "+string(rune(i+1)),
			"CODE"+string(rune(i+1)),
			"Desc "+string(rune(i+1)),
			true,
		)
	}
	return roomTypes
}

func NewRoomTypePaginationQuery(page, pageSize int) *dto.RoomTypePaginationQuery {
	return &dto.RoomTypePaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewCreateRoomTypeRequest(name, code, description string, isActive bool) *dto.CreateRoomTypeRequest {
	return &dto.CreateRoomTypeRequest{
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    &isActive,
	}
}

func NewUpdateRoomTypeRequest(name, code, description string, isActive bool) *dto.UpdateRoomTypeRequest {
	return &dto.UpdateRoomTypeRequest{
		Name:        PtrString(name),
		Code:        PtrString(code),
		Description: PtrString(description),
		IsActive:    PtrBool(isActive),
	}
}

func NewTestRoomTypeResponseWithData(id uint, name, code, description string, isActive bool) *dto.RoomTypeResponse {
	return &dto.RoomTypeResponse{
		ID:          id,
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
