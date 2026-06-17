package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestReferralWithData(id uint, patientID uint) *models.Referral {
	now := time.Now()
	refDeptID := uint(2)
	refDocID := uint(3)

	return &models.Referral{
		ID:                     id,
		ReferralNumber:         "REF-2024-000001",
		PatientID:              patientID,
		MedicalRecordID:        1,
		ReferringDoctorID:      2,
		ReferralDate:           now.Format("2006-01-02"),
		ReferralType:           "internal",
		ReferredToDepartmentID: &refDeptID,
		ReferredToDoctorID:     &refDocID,
		Reason:                 "Konsultasi Spesialis Jantung",
		Diagnosis:              "Suspect CHD",
		Priority:               "routine",
		Status:                 "pending",
		CreatedAt:              now,
		UpdatedAt:              now,
		Patient: &models.Patient{
			ID:          patientID,
			FullName:    "Budi Santoso",
			PatientCode: "RM-001",
		},
		ReferringDoctor: &models.Doctor{
			ID:       2,
			FullName: "Dr. Andi",
		},
	}
}

func NewTestReferralList(count int) []models.Referral {
	var list []models.Referral
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestReferralWithData(uint(i), 1))
	}
	return list
}

func NewReferralPaginationQuery(page, pageSize int) *dto.ReferralPaginationQuery {
	return &dto.ReferralPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewTestReferralResponse(r *models.Referral) *dto.ReferralResponse {
	return &dto.ReferralResponse{
		ID:                r.ID,
		ReferralNumber:    r.ReferralNumber,
		PatientID:         r.PatientID,
		MedicalRecordID:   r.MedicalRecordID,
		ReferringDoctorID: r.ReferringDoctorID,
		ReferralDate:      r.ReferralDate,
		ReferralType:      r.ReferralType,
		Reason:            r.Reason,
		Diagnosis:         r.Diagnosis,
		Priority:          r.Priority,
		Status:            r.Status,
		CreatedAt:         r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         r.UpdatedAt.Format(time.RFC3339),
	}
}

func NewCreateReferralRequest(patientID uint) *dto.CreateReferralRequest {
	refDeptID := uint(2)
	return &dto.CreateReferralRequest{
		PatientID:              patientID,
		MedicalRecordID:        1,
		ReferringDoctorID:      2,
		ReferralDate:           "2024-01-01",
		ReferralType:           "internal",
		ReferredToDepartmentID: &refDeptID,
		Reason:                 "Konsultasi Spesialis",
		Priority:               "routine",
	}
}

func NewUpdateReferralRequest() *dto.UpdateReferralRequest {
	reason := "Updated Reason"
	return &dto.UpdateReferralRequest{
		Reason: &reason,
	}
}
