package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestFamilyHistoryWithData(id uint, patientID uint) *models.FamilyHistory {
	now := time.Now()
	return &models.FamilyHistory{
		ID:            id,
		PatientID:     patientID,
		FamilyMember:  "father",
		ConditionName: "Hipertensi",
		Relation:      "Biological Father",
		Notes:         "Meninggal karena stroke di usia 68 tahun",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func NewTestFamilyHistoryList(count int) []models.FamilyHistory {
	var list []models.FamilyHistory
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestFamilyHistoryWithData(uint(i), 1))
	}
	return list
}

func NewFamilyHistoryPaginationQuery(page, pageSize int) *dto.FamilyHistoryPaginationQuery {
	return &dto.FamilyHistoryPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewTestFamilyHistoryResponse(f *models.FamilyHistory) *dto.FamilyHistoryResponse {
	return &dto.FamilyHistoryResponse{
		ID:            f.ID,
		PatientID:     f.PatientID,
		FamilyMember:  f.FamilyMember,
		ConditionName: f.ConditionName,
		Relation:      f.Relation,
		Notes:         f.Notes,
		CreatedAt:     f.CreatedAt,
		UpdatedAt:     f.UpdatedAt,
	}
}

func NewCreateFamilyHistoryRequest(patientID uint) *dto.CreateFamilyHistoryRequest {
	return &dto.CreateFamilyHistoryRequest{
		PatientID:     patientID,
		FamilyMember:  "father",
		ConditionName: "Hipertensi",
		Relation:      "Biological Father",
		Notes:         "Meninggal karena stroke di usia 68 tahun",
	}
}

func NewUpdateFamilyHistoryRequest() *dto.UpdateFamilyHistoryRequest {
	return &dto.UpdateFamilyHistoryRequest{
		Notes: "Updated Notes",
	}
}
