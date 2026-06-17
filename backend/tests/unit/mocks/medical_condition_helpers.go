package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestMedicalConditionWithData(id uint, patientID uint) *models.MedicalCondition {
	now := time.Now()
	diagnosedDate := now.AddDate(-1, 0, 0)
	return &models.MedicalCondition{
		ID:            id,
		PatientID:     patientID,
		ConditionName: "Diabetes Mellitus Tipe 2",
		ICDCode:       "E11",
		DiagnosedDate: &diagnosedDate,
		Status:        "ongoing",
		Notes:         "Terkontrol dengan Metformin",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func NewTestMedicalConditionList(count int) []models.MedicalCondition {
	var list []models.MedicalCondition
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestMedicalConditionWithData(uint(i), 1))
	}
	return list
}

func NewMedicalConditionPaginationQuery(page, pageSize int) *dto.MedicalConditionPaginationQuery {
	return &dto.MedicalConditionPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewTestMedicalConditionResponse(c *models.MedicalCondition) *dto.MedicalConditionResponse {
	diagnosedDateStr := ""
	if c.DiagnosedDate != nil {
		diagnosedDateStr = c.DiagnosedDate.Format("2006-01-02")
	}

	return &dto.MedicalConditionResponse{
		ID:            c.ID,
		PatientID:     c.PatientID,
		ConditionName: c.ConditionName,
		ICDCode:       c.ICDCode,
		DiagnosedDate: diagnosedDateStr,
		Status:        c.Status,
		Notes:         c.Notes,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

func NewCreateMedicalConditionRequest(patientID uint) *dto.CreateMedicalConditionRequest {
	return &dto.CreateMedicalConditionRequest{
		PatientID:     patientID,
		ConditionName: "Diabetes Mellitus Tipe 2",
		ICDCode:       "E11",
		DiagnosedDate: "2023-01-01",
		Status:        "ongoing",
		Notes:         "Terkontrol dengan Metformin",
	}
}

func NewUpdateMedicalConditionRequest() *dto.UpdateMedicalConditionRequest {
	return &dto.UpdateMedicalConditionRequest{
		Status: "managed",
	}
}
