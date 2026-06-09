package dto

import (
	"time"
)

type LabTestMedicalRecordResponse struct {
	ID             uint   `json:"id"`
	VisitDate      string `json:"visit_date"`
	ChiefComplaint string `json:"chief_complaint"`
}

type LabTestTypeResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	Category string `json:"category,omitempty"`
}

type LabTestDoctorResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization,omitempty"`
}

type LabTestResponse struct {
	ID                   uint                          `json:"id"`
	MedicalRecordID      uint                          `json:"medical_record_id"`
	MedicalRecord        *LabTestMedicalRecordResponse `json:"medical_record,omitempty"`
	TestTypeID           uint                          `json:"test_type_id"`
	TestType             *LabTestTypeResponse          `json:"test_type,omitempty"`
	OrderedByDoctorID    uint                          `json:"ordered_by_doctor_id"`
	OrderedByDoctor      *LabTestDoctorResponse        `json:"ordered_by_doctor,omitempty"`
	OrderDate            string                        `json:"order_date"`
	SampleCollectionDate *string                       `json:"sample_collection_date"`
	TestStartDate        *string                       `json:"test_start_date,omitempty"`
	ResultDate           *string                       `json:"result_date"`
	ResultValue          *string                       `json:"result_value"`
	ResultUnit           *string                       `json:"result_unit"`
	ReferenceRange       *string                       `json:"reference_range"`
	Status               string                        `json:"status"`
	Notes                string                        `json:"notes"`
	CreatedAt            time.Time                     `json:"created_at"`
	UpdatedAt            time.Time                     `json:"updated_at"`
}

type DeletedLabTestResponse struct {
	LabTestResponse
	DeletedAt *time.Time `json:"deleted_at"`
}

type LabTestListResponse struct {
	Data []LabTestResponse         `json:"data"`
	Meta LabTestPaginationMeta     `json:"meta"`
}

type LabTestDeletedListResponse struct {
	Data []DeletedLabTestResponse  `json:"data"`
	Meta LabTestPaginationMeta     `json:"meta"`
}

type LabTestPaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type CreateLabTestRequest struct {
	MedicalRecordID   uint   `json:"medical_record_id" binding:"required"`
	TestTypeID        uint   `json:"test_type_id" binding:"required"`
	OrderedByDoctorID uint   `json:"ordered_by_doctor_id" binding:"required"`
	OrderDate         string `json:"order_date" binding:"required,datetime=2006-01-02"`
	Status            string `json:"status" binding:"omitempty"`
	Notes             string `json:"notes" binding:"omitempty"`
}

type UpdateLabTestRequest struct {
	Notes          *string `json:"notes" binding:"omitempty"`
	ReferenceRange *string `json:"reference_range" binding:"omitempty"`
}

type CompleteLabTestRequest struct {
	ResultValue    *string `json:"result_value" binding:"omitempty"`
	ResultUnit     *string `json:"result_unit" binding:"omitempty,max=50"`
	ReferenceRange *string `json:"reference_range" binding:"omitempty,max=200"`
	Notes          *string `json:"notes" binding:"omitempty"`
}

type LabTestPaginationQuery struct {
	Page              int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize          int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	SortBy            string `form:"sort_by,default=created_at" binding:"omitempty,oneof=id order_date result_date status created_at"`
	SortDir           string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
	Search            string `form:"search" binding:"omitempty"`
	Status            string `form:"status" binding:"omitempty,oneof=ordered sample_collected in_progress completed cancelled"`
	NotStatus         string `form:"not_status" binding:"omitempty,oneof=ordered sample_collected in_progress completed cancelled"`
	MedicalRecordID   *uint  `form:"medical_record_id" binding:"omitempty"`
	TestTypeID        *uint  `form:"test_type_id" binding:"omitempty"`
	OrderedByDoctorID *uint  `form:"ordered_by_doctor_id" binding:"omitempty"`
}
