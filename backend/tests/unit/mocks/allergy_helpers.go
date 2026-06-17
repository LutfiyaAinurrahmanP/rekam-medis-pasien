package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestAllergyWithData(id uint, patientID uint) *models.Allergy {
	now := time.Now()
	return &models.Allergy{
		ID:           id,
		PatientID:    patientID,
		AllergenType: "drug",
		AllergenName: "Penisilin",
		Reaction:     "Ruam",
		Severity:     "severe",
		Notes:        "Hindari golongan Beta-laktam",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func NewTestAllergyList(count int) []models.Allergy {
	var list []models.Allergy
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestAllergyWithData(uint(i), 1))
	}
	return list
}

func NewAllergyPaginationQuery(page, pageSize int) *dto.AllergyPaginationQuery {
	return &dto.AllergyPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewTestAllergyResponse(a *models.Allergy) *dto.AllergyResponse {
	return &dto.AllergyResponse{
		ID:           a.ID,
		PatientID:    a.PatientID,
		AllergenType: a.AllergenType,
		AllergenName: a.AllergenName,
		Reaction:     a.Reaction,
		Severity:     a.Severity,
		Notes:        a.Notes,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

func NewCreateAllergyRequest(patientID uint) *dto.CreateAllergyRequest {
	return &dto.CreateAllergyRequest{
		PatientID:    patientID,
		AllergenName: "Penisilin",
		AllergenType: "drug",
		Reaction:     "Ruam",
		Severity:     "severe",
		Notes:        "Hindari golongan Beta-laktam",
	}
}

func NewUpdateAllergyRequest() *dto.UpdateAllergyRequest {
	return &dto.UpdateAllergyRequest{
		Severity: "moderate",
	}
}
