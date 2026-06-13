package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestTypeTestWithData(id uint, name, code string, categoryID uint, description string, price float64, isActive bool) *models.TypeTest {
	return &models.TypeTest{
		ID:                 id,
		Name:               name,
		Code:               code,
		TypeTestCategoryID: categoryID,
		Description:        description,
		Price:              price,
		IsActive:           isActive,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}

func NewTestTypeTestList(count int) []models.TypeTest {
	typeTests := make([]models.TypeTest, count)
	for i := 0; i < count; i++ {
		typeTests[i] = *NewTestTypeTestWithData(
			uint(i+1),
			"Test "+string(rune(i+1)),
			"T00"+string(rune(i+1)),
			uint(1),
			"Description "+string(rune(i+1)),
			50000.0,
			true,
		)
	}
	return typeTests
}

func NewTypeTestPaginationQuery(page, pageSize int) *dto.TypeTestPaginationQuery {
	return &dto.TypeTestPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "name",
		SortDir:  "asc",
	}
}

func NewCreateTypeTestRequest(name, code string, categoryID uint, description string, price float64, isActive bool) *dto.CreateTypeTestRequest {
	return &dto.CreateTypeTestRequest{
		Name:               name,
		Code:               code,
		TypeTestCategoryID: categoryID,
		Description:        description,
		Price:              PtrFloat64(price),
		IsActive:           PtrBool(isActive),
	}
}

func NewUpdateTypeTestRequest(name, code string, categoryID uint, description string, price float64, isActive bool) *dto.UpdateTypeTestRequest {
	return &dto.UpdateTypeTestRequest{
		Name:               PtrString(name),
		Code:               PtrString(code),
		TypeTestCategoryID: PtrUint(categoryID),
		Description:        PtrString(description),
		Price:              PtrFloat64(price),
		IsActive:           PtrBool(isActive),
	}
}

func NewTestTypeTestResponseWithData(id uint, name, code string, categoryID uint, description string, price float64, isActive bool) *dto.TypeTestResponse {
	return &dto.TypeTestResponse{
		ID:                 id,
		Name:               name,
		Code:               code,
		TypeTestCategoryID: categoryID,
		Description:        description,
		Price:              price,
		IsActive:           isActive,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
