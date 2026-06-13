package mocks

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func NewTestMedicineTypeWithData(id uint, name, code, desc string, isDeleted bool) *models.MedicineType {
	now := time.Now()
	mt := &models.MedicineType{
		ID:          id,
		Name:        name,
		Code:        code,
		Description: desc,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if isDeleted {
		mt.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	}

	return mt
}

func NewTestMedicineTypeList(count int) []models.MedicineType {
	var list []models.MedicineType
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestMedicineTypeWithData(
			uint(i),
			fmt.Sprintf("Medicine Type %d", i),
			fmt.Sprintf("MT%03d", i),
			fmt.Sprintf("Description for Medicine Type %d", i),
			false,
		))
	}
	return list
}

func NewTestMedicineTypeResponse(mt *models.MedicineType) *dto.MedicineTypeResponse {
	return &dto.MedicineTypeResponse{
		ID:          mt.ID,
		Name:        mt.Name,
		Code:        mt.Code,
		Description: mt.Description,
		IsActive:    mt.IsActive,
		CreatedAt:   mt.CreatedAt,
		UpdatedAt:   mt.UpdatedAt,
	}
}

func NewMedicineTypePaginationQuery(page, pageSize int) *dto.MedicineTypePaginationQuery {
	return &dto.MedicineTypePaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewCreateMedicineTypeRequest(name, code, desc string, isActive bool) *dto.CreateMedicineTypeRequest {
	active := isActive
	return &dto.CreateMedicineTypeRequest{
		Name:        name,
		Code:        code,
		Description: desc,
		IsActive:    &active,
	}
}

func NewUpdateMedicineTypeRequest(name, code, desc string, isActive bool) *dto.UpdateMedicineTypeRequest {
	n := name
	c := code
	d := desc
	a := isActive
	return &dto.UpdateMedicineTypeRequest{
		Name:        &n,
		Code:        &c,
		Description: &d,
		IsActive:    &a,
	}
}
