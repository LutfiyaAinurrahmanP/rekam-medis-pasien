package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestSurgicalHistoryWithData(id uint, patientID uint) *models.SurgicalHistory {
	now := time.Now()
	surgeryDate := now.AddDate(-1, 0, 0)
	return &models.SurgicalHistory{
		ID:            id,
		PatientID:     patientID,
		ProcedureName: "Appendektomi",
		SurgeryDate:   surgeryDate,
		SurgeonName:   "dr. Bima",
		Hospital:      "RS Dr. Soetomo",
		Complication:  "None",
		Notes:         "Berjalan lancar",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func NewTestSurgicalHistoryList(count int) []models.SurgicalHistory {
	var list []models.SurgicalHistory
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestSurgicalHistoryWithData(uint(i), 1))
	}
	return list
}

func NewSurgicalHistoryPaginationQuery(page, pageSize int) *dto.SurgicalHistoryPaginationQuery {
	return &dto.SurgicalHistoryPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewTestSurgicalHistoryResponse(s *models.SurgicalHistory) *dto.SurgicalHistoryResponse {
	return &dto.SurgicalHistoryResponse{
		ID:            s.ID,
		PatientID:     s.PatientID,
		ProcedureName: s.ProcedureName,
		SurgeryDate:   s.SurgeryDate.Format("2006-01-02"),
		SurgeonName:   s.SurgeonName,
		Hospital:      s.Hospital,
		Complication:  s.Complication,
		Notes:         s.Notes,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func NewCreateSurgicalHistoryRequest(patientID uint) *dto.CreateSurgicalHistoryRequest {
	return &dto.CreateSurgicalHistoryRequest{
		PatientID:     patientID,
		ProcedureName: "Appendektomi",
		SurgeryDate:   "2023-01-01",
		SurgeonName:   "dr. Bima",
		Hospital:      "RS Dr. Soetomo",
		Complication:  "None",
		Notes:         "Berjalan lancar",
	}
}

func NewUpdateSurgicalHistoryRequest() *dto.UpdateSurgicalHistoryRequest {
	return &dto.UpdateSurgicalHistoryRequest{
		Notes: "Updated Notes",
	}
}
