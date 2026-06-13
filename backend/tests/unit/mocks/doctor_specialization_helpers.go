package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

// ============= Doctor Specialization Model Builders =============

func NewTestDoctorSpecialization() *models.DoctorSpecialization {
	return &models.DoctorSpecialization{
		ID:          1,
		Name:        "Cardiology",
		Code:        "CARD",
		Description: "Heart Specialization",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewTestDoctorSpecializationWithData(id uint, name, code, description string, isActive bool) *models.DoctorSpecialization {
	return &models.DoctorSpecialization{
		ID:          id,
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewTestDoctorSpecializationList(count int) []models.DoctorSpecialization {
	specs := make([]models.DoctorSpecialization, count)
	for i := 0; i < count; i++ {
		specs[i] = *NewTestDoctorSpecializationWithData(
			uint(i+1),
			"Specialization "+string(rune(i+1)),
			"SPEC"+string(rune(i+1)),
			"Description "+string(rune(i+1)),
			true,
		)
	}
	return specs
}

// ============= Doctor Specialization DTO Builders =============

func NewTestDoctorSpecializationResponseWithData(id uint, name, code, description string, isActive bool) *dto.DoctorSpecializationResponse {
	return &dto.DoctorSpecializationResponse{
		ID:          id,
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewDoctorSpecializationPaginationQuery(page, pageSize int) *dto.DoctorSpecializationPaginationQuery {
	return &dto.DoctorSpecializationPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewCreateDoctorSpecializationRequest(name, code, description string, isActive bool) *dto.CreateDoctorSpecializationRequest {
	return &dto.CreateDoctorSpecializationRequest{
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    &isActive,
	}
}

func NewUpdateDoctorSpecializationRequest(name, code, description string, isActive bool) *dto.UpdateDoctorSpecializationRequest {
	return &dto.UpdateDoctorSpecializationRequest{
		Name:        PtrString(name),
		Code:        PtrString(code),
		Description: PtrString(description),
		IsActive:    PtrBool(isActive),
	}
}
