package mocks

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func NewTestMedicineWithData(id uint, name, genericName, brandName string, medicineTypeID uint, strength, manufacturer, unit string, stock int, price float64, isDeleted bool) *models.Medicine {
	now := time.Now()
	m := &models.Medicine{
		ID:             id,
		Name:           name,
		GenericName:    genericName,
		BrandName:      brandName,
		MedicineTypeID: medicineTypeID,
		Strength:       strength,
		Manufacturer:   manufacturer,
		Unit:           unit,
		StockQuantity:  stock,
		Price:          price,
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if isDeleted {
		m.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	}

	return m
}

func NewTestMedicineList(count int) []models.Medicine {
	var list []models.Medicine
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestMedicineWithData(
			uint(i),
			fmt.Sprintf("Medicine %d", i),
			fmt.Sprintf("Generic %d", i),
			fmt.Sprintf("Brand %d", i),
			1,
			"500mg",
			"Manufacturer",
			"Tablet",
			100,
			15000,
			false,
		))
	}
	return list
}

func NewTestMedicineResponse(m *models.Medicine) *dto.MedicineResponse {
	return &dto.MedicineResponse{
		ID:             m.ID,
		Name:           m.Name,
		GenericName:    m.GenericName,
		BrandName:      m.BrandName,
		MedicineTypeID: m.MedicineTypeID,
		Strength:       m.Strength,
		Manufacturer:   m.Manufacturer,
		Unit:           m.Unit,
		StockQuantity:  m.StockQuantity,
		Price:          m.Price,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func NewTestDeletedMedicineResponse(m *models.Medicine) *dto.DeletedMedicineResponse {
	return &dto.DeletedMedicineResponse{
		ID:             m.ID,
		Name:           m.Name,
		GenericName:    m.GenericName,
		BrandName:      m.BrandName,
		MedicineTypeID: m.MedicineTypeID,
		Strength:       m.Strength,
		Manufacturer:   m.Manufacturer,
		Unit:           m.Unit,
		StockQuantity:  m.StockQuantity,
		Price:          m.Price,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      &m.DeletedAt.Time,
	}
}

func NewMedicinePaginationQuery(page, pageSize int) *dto.MedicinePaginationQuery {
	return &dto.MedicinePaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewCreateMedicineRequest(name, genericName, brandName string, medicineTypeID uint, strength, manufacturer, unit string, stock int, price float64, isActive bool) *dto.CreateMedicineRequest {
	active := isActive
	stockQty := stock
	pr := price
	return &dto.CreateMedicineRequest{
		Name:           name,
		GenericName:    genericName,
		BrandName:      brandName,
		MedicineTypeID: &medicineTypeID,
		Strength:       strength,
		Manufacturer:   manufacturer,
		Unit:           unit,
		StockQuantity:  &stockQty,
		Price:          &pr,
		IsActive:       &active,
	}
}

func NewUpdateMedicineRequest(name, genericName, brandName string, medicineTypeID uint, strength, manufacturer, unit string, price float64) *dto.UpdateMedicineRequest {
	n := name
	gn := genericName
	bn := brandName
	mtID := medicineTypeID
	s := strength
	m := manufacturer
	u := unit
	p := price
	return &dto.UpdateMedicineRequest{
		Name:           &n,
		GenericName:    &gn,
		BrandName:      &bn,
		MedicineTypeID: &mtID,
		Strength:       &s,
		Manufacturer:   &m,
		Unit:           &u,
		Price:          &p,
	}
}
