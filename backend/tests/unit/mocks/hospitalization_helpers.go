package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func NewTestHospitalizationWithData(id uint, patientID, doctorID, roomID uint, admissionDate, admissionTime, status string, isDeleted bool) *models.Hospitalization {
	now := time.Now()
	h := &models.Hospitalization{
		ID:                 id,
		PatientID:          patientID,
		DoctorID:           doctorID,
		RoomID:             roomID,
		AdmissionDate:      admissionDate,
		AdmissionTime:      admissionTime,
		ReasonForAdmission: "Reason",
		Status:             status,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if isDeleted {
		h.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	}

	return h
}

func NewTestHospitalizationList(count int) []models.Hospitalization {
	var list []models.Hospitalization
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestHospitalizationWithData(
			uint(i),
			uint(1),
			uint(1),
			uint(1),
			"2023-12-01",
			"10:00:00",
			"admitted",
			false,
		))
	}
	return list
}

func NewTestHospitalizationResponse(h *models.Hospitalization) *dto.HospitalizationResponse {
	return &dto.HospitalizationResponse{
		ID:              h.ID,
		PatientID:       h.PatientID,
		DoctorID:        h.DoctorID,
		RoomID:          h.RoomID,
		AdmissionDate:   h.AdmissionDate,
		AdmissionTime:   h.AdmissionTime,
		DischargeDate:   h.DischargeDate,
		DischargeTime:   h.DischargeTime,
		AdmissionReason: h.ReasonForAdmission,
		Status:          h.Status,
		Notes:           h.Notes,
		CreatedAt:       h.CreatedAt,
		UpdatedAt:       h.UpdatedAt,
	}
}

func NewHospitalizationPaginationQuery(page, pageSize int) *dto.HospitalizationPaginationQuery {
	return &dto.HospitalizationPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "admission_date",
		SortDir:  "desc",
	}
}

func NewCreateHospitalizationRequest(patientID, doctorID, roomID uint, admissionDate, reason string) *dto.CreateHospitalizationRequest {
	return &dto.CreateHospitalizationRequest{
		PatientID:         patientID,
		AttendingDoctorID: doctorID,
		RoomID:            roomID,
		AdmissionDate:     admissionDate,
		AdmissionReason:   reason,
	}
}

func NewUpdateHospitalizationRequest(roomID, doctorID uint, reason string) *dto.UpdateHospitalizationRequest {
	return &dto.UpdateHospitalizationRequest{
		RoomID:            &roomID,
		AttendingDoctorID: &doctorID,
		AdmissionReason:   &reason,
	}
}
