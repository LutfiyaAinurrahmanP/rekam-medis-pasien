package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func NewTestPrescriptionWithData(id uint, medicalRecordID, doctorID uint, status string, isDeleted bool) *models.Prescription {
	now := time.Now()
	p := &models.Prescription{
		ID:               id,
		MedicalRecordID:  medicalRecordID,
		DoctorID:         doctorID,
		PrescriptionDate: "2023-12-01",
		Notes:            "Test Notes",
		Status:           status,
		CreatedAt:        now,
		UpdatedAt:        now,
		Items: []models.PrescriptionItem{
			*NewTestPrescriptionItemWithData(1, id, 1),
		},
	}

	if isDeleted {
		p.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	}

	return p
}

func NewTestPrescriptionItemWithData(id, prescriptionID, medicineID uint) *models.PrescriptionItem {
	now := time.Now()
	return &models.PrescriptionItem{
		ID:             id,
		PrescriptionID: prescriptionID,
		MedicineID:     medicineID,
		Dosage:         "1 tablet",
		Frequency:      "3x sehari",
		DurationDays:   3,
		Quantity:       9,
		Instructions:   "Sesudah makan",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func NewTestPrescriptionList(count int) []models.Prescription {
	var list []models.Prescription
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestPrescriptionWithData(uint(i), 1, 1, "pending", false))
	}
	return list
}

func NewTestPrescriptionResponse(p *models.Prescription) *dto.PrescriptionResponse {
	items := make([]dto.PrescriptionItemResponse, len(p.Items))
	for i, it := range p.Items {
		items[i] = *NewTestPrescriptionItemResponse(&it)
	}

	return &dto.PrescriptionResponse{
		ID:               p.ID,
		MedicalRecordID:  p.MedicalRecordID,
		DoctorID:         p.DoctorID,
		PrescriptionDate: p.PrescriptionDate,
		Notes:            p.Notes,
		Status:           p.Status,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
		Items:            items,
	}
}

func NewTestPrescriptionItemResponse(it *models.PrescriptionItem) *dto.PrescriptionItemResponse {
	return &dto.PrescriptionItemResponse{
		ID:             it.ID,
		PrescriptionID: it.PrescriptionID,
		MedicineID:     it.MedicineID,
		Dosage:         it.Dosage,
		Frequency:      it.Frequency,
		DurationDays:   it.DurationDays,
		Quantity:       it.Quantity,
		Instructions:   it.Instructions,
	}
}

func NewPrescriptionPaginationQuery(page, pageSize int) *dto.PrescriptionPaginationQuery {
	return &dto.PrescriptionPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewCreatePrescriptionRequest(medicalRecordID, doctorID uint) *dto.CreatePrescriptionRequest {
	return &dto.CreatePrescriptionRequest{
		MedicalRecordID:  medicalRecordID,
		DoctorID:         doctorID,
		PrescriptionDate: "2023-12-01",
		Notes:            "Test Notes",
		Status:           "pending",
	}
}

func NewUpdatePrescriptionRequest(notes string) *dto.UpdatePrescriptionRequest {
	return &dto.UpdatePrescriptionRequest{
		Notes: &notes,
	}
}

func NewCreatePrescriptionItemRequest(medicineID uint) *dto.CreatePrescriptionItemRequest {
	return &dto.CreatePrescriptionItemRequest{
		MedicineID:   medicineID,
		Dosage:       "1 tablet",
		Frequency:    "3x sehari",
		DurationDays: 3,
		Quantity:     9,
		Instructions: "Sesudah makan",
	}
}

func NewUpdatePrescriptionItemRequest(dosage string) *dto.UpdatePrescriptionItemRequest {
	return &dto.UpdatePrescriptionItemRequest{
		Dosage: &dosage,
	}
}
