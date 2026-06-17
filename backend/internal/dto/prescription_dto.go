package dto

import (
	"time"
)

type PrescriptionMedicalRecordResponse struct {
	ID        uint   `json:"id"`
	VisitDate string `json:"visit_date"`
}

type PrescriptionDoctorResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization,omitempty"`
}

type PrescriptionMedicineResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit,omitempty"`
}

type PrescriptionItemResponse struct {
	ID             uint                          `json:"id"`
	PrescriptionID uint                          `json:"prescription_id"`
	MedicineID     uint                          `json:"medicine_id"`
	Dosage         string                        `json:"dosage"`
	Frequency      string                        `json:"frequency"`
	DurationDays   int                           `json:"duration_days"`
	Quantity       int                           `json:"quantity"`
	Instructions   string                        `json:"instructions,omitempty"`
	Medicine       *PrescriptionMedicineResponse `json:"medicine,omitempty"`
}

type PrescriptionResponse struct {
	ID               uint                               `json:"id"`
	MedicalRecordID  uint                               `json:"medical_record_id"`
	MedicalRecord    *PrescriptionMedicalRecordResponse `json:"medical_record,omitempty"`
	DoctorID         uint                               `json:"doctor_id"`
	Doctor           *PrescriptionDoctorResponse        `json:"doctor,omitempty"`
	PrescriptionDate string                             `json:"prescription_date"`
	Notes            string                             `json:"notes,omitempty"`
	Status           string                             `json:"status"`
	CreatedAt        time.Time                          `json:"created_at"`
	UpdatedAt        time.Time                          `json:"updated_at"`
	Items            []PrescriptionItemResponse         `json:"items,omitempty"`
}

type DeletedPrescriptionResponse struct {
	PrescriptionResponse
	DeletedAt *time.Time `json:"deleted_at"`
}

type PrescriptionListResponse struct {
	Data []PrescriptionResponse     `json:"data"`
	Meta PrescriptionPaginationMeta `json:"meta"`
}

type PrescriptionDeletedListResponse struct {
	Data []DeletedPrescriptionResponse `json:"data"`
	Meta PrescriptionPaginationMeta    `json:"meta"`
}

type CreatePrescriptionRequest struct {
	MedicalRecordID  uint   `json:"medical_record_id" binding:"required"`
	DoctorID         uint   `json:"doctor_id" binding:"required"`
	PrescriptionDate string `json:"prescription_date" binding:"required,datetime=2006-01-02"`
	Notes            string `json:"notes" binding:"omitempty"`
	Status           string `json:"status" binding:"omitempty"`
}

type UpdatePrescriptionRequest struct {
	Notes *string `json:"notes" binding:"omitempty"`
}

type CreatePrescriptionItemRequest struct {
	MedicineID   uint   `json:"medicine_id" binding:"required"`
	Dosage       string `json:"dosage" binding:"required"`
	Frequency    string `json:"frequency" binding:"required"`
	DurationDays int    `json:"duration_days" binding:"required,min=1"`
	Quantity     int    `json:"quantity" binding:"required,min=1"`
	Instructions string `json:"instructions" binding:"omitempty"`
}

type UpdatePrescriptionItemRequest struct {
	Dosage       *string `json:"dosage" binding:"omitempty"`
	Frequency    *string `json:"frequency" binding:"omitempty"`
	DurationDays *int    `json:"duration_days" binding:"omitempty,min=1"`
	Quantity     *int    `json:"quantity" binding:"omitempty,min=1"`
	Instructions *string `json:"instructions" binding:"omitempty"`
}

type PrescriptionPaginationQuery struct {
	Page            int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize        int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	DoctorID        *uint  `form:"doctor_id" binding:"omitempty"`
	MedicalRecordID *uint  `form:"medical_record_id" binding:"omitempty"`
	Status          string `form:"status" binding:"omitempty,oneof=pending dispensed cancelled"`
	Search          string `form:"search" binding:"omitempty"`
	SortBy          string `form:"sort_by,default=created_at" binding:"omitempty,oneof=id prescription_date status created_at"`
	SortDir         string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type PrescriptionPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
