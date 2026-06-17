package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestTypeTestCategoryWithData(id uint, name, code, description string, isActive bool) *models.TypeTestCategory {
	return &models.TypeTestCategory{
		ID:          id,
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewTestTypeTestCategoryList(count int) []models.TypeTestCategory {
	categories := make([]models.TypeTestCategory, count)
	for i := 0; i < count; i++ {
		categories[i] = *NewTestTypeTestCategoryWithData(
			uint(i+1),
			"Category "+string(rune(i+1)),
			"CAT00"+string(rune(i+1)),
			"Description "+string(rune(i+1)),
			true,
		)
	}
	return categories
}

func NewTypeTestCategoryPaginationQuery(page, pageSize int) *dto.TypeTestCategoryPaginationQuery {
	return &dto.TypeTestCategoryPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewCreateTypeTestCategoryRequest(name, code, description string, isActive bool) *dto.CreateTypeTestCategoryRequest {
	return &dto.CreateTypeTestCategoryRequest{
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    PtrBool(isActive),
	}
}

func NewUpdateTypeTestCategoryRequest(name, code, description string, isActive bool) *dto.UpdateTypeTestCategoryRequest {
	return &dto.UpdateTypeTestCategoryRequest{
		Name:        PtrString(name),
		Code:        PtrString(code),
		Description: PtrString(description),
		IsActive:    PtrBool(isActive),
	}
}

func NewTestTypeTestCategoryResponseWithData(id uint, name, code, description string, isActive bool) *dto.TypeTestCategoryResponse {
	return &dto.TypeTestCategoryResponse{
		ID:          id,
		Name:        name,
		Code:        code,
		Description: description,
		IsActive:    isActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
