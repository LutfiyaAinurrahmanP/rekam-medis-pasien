package mocks

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func NewTestMedicalRecordWithData(id uint, patientID uint, doctorID uint, visitDate, complaint, diagnosis, plan, status string, isDeleted bool) *models.MedicalRecord {
	now := time.Now()
	m := &models.MedicalRecord{
		ID:                  id,
		PatientID:           patientID,
		DoctorID:            doctorID,
		VisitDate:           visitDate,
		ChiefComplaint:      complaint,
		HistoryOfIllness:    "History",
		PhysicalExamination: "Physical exam",
		Diagnosis:           diagnosis,
		TreatmentPlan:       plan,
		Notes:               "Notes",
		Status:              status,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if isDeleted {
		m.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	}

	return m
}

func NewTestMedicalRecordList(count int) []models.MedicalRecord {
	var list []models.MedicalRecord
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestMedicalRecordWithData(
			uint(i),
			uint(1),
			uint(1),
			"2023-12-01",
			fmt.Sprintf("Complaint %d", i),
			fmt.Sprintf("Diagnosis %d", i),
			fmt.Sprintf("Plan %d", i),
			"draft",
			false,
		))
	}
	return list
}

func NewTestMedicalRecordResponse(m *models.MedicalRecord) *dto.MedicalRecordResponse {
	return &dto.MedicalRecordResponse{
		ID:                  m.ID,
		PatientID:           m.PatientID,
		DoctorID:            m.DoctorID,
		AppointmentID:       m.AppointmentID,
		VisitDate:           m.VisitDate,
		ChiefComplaint:      m.ChiefComplaint,
		HistoryOfIllness:    m.HistoryOfIllness,
		PhysicalExamination: m.PhysicalExamination,
		Diagnosis:           m.Diagnosis,
		TreatmentPlan:       m.TreatmentPlan,
		Notes:               m.Notes,
		Status:              m.Status,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func NewMedicalRecordPaginationQuery(page, pageSize int) *dto.MedicalRecordPaginationQuery {
	return &dto.MedicalRecordPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "visit_date",
		SortDir:  "desc",
	}
}

func NewCreateMedicalRecordRequest(patientID uint, visitDate, complaint, diagnosis, plan string) *dto.CreateMedicalRecordRequest {
	return &dto.CreateMedicalRecordRequest{
		PatientID:      patientID,
		VisitDate:      visitDate,
		ChiefComplaint: complaint,
		Diagnosis:      diagnosis,
		TreatmentPlan:  plan,
	}
}

func NewUpdateMedicalRecordRequest(complaint, diagnosis, plan string) *dto.UpdateMedicalRecordRequest {
	return &dto.UpdateMedicalRecordRequest{
		ChiefComplaint: &complaint,
		Diagnosis:      &diagnosis,
		TreatmentPlan:  &plan,
	}
}
