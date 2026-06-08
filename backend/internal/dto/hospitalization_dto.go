package dto

import (
	"time"
)

type HospitalizationPatientResponse struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	MedicalRecordNumber string `json:"medical_record_number"`
}

type HospitalizationDoctorResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization,omitempty"`
}

type HospitalizationRoomResponse struct {
	ID         uint   `json:"id"`
	RoomNumber string `json:"room_number"`
	RoomType   string `json:"room_type"`
}

type HospitalizationResponse struct {
	ID              uint                            `json:"id"`
	PatientID       uint                            `json:"patient_id"`
	Patient         *HospitalizationPatientResponse `json:"patient,omitempty"`
	DoctorID        uint                            `json:"attending_doctor_id"` // matching JSON attending_doctor_id
	Doctor          *HospitalizationDoctorResponse  `json:"attending_doctor,omitempty"`
	RoomID          uint                            `json:"room_id"`
	Room            *HospitalizationRoomResponse    `json:"room,omitempty"`
	AdmissionDate   string                          `json:"admission_date"`
	AdmissionTime   string                          `json:"admission_time"`
	DischargeDate   *string                         `json:"discharge_date,omitempty"`
	DischargeTime   *string                         `json:"discharge_time,omitempty"`
	AdmissionReason string                          `json:"admission_reason"` // mapping ReasonForAdmission
	Status          string                          `json:"status"`
	Notes           string                          `json:"notes,omitempty"`
	CreatedAt       time.Time                       `json:"created_at"`
	UpdatedAt       time.Time                       `json:"updated_at"`
}

type DeletedHospitalizationResponse struct {
	HospitalizationResponse
	DeletedAt *time.Time `json:"deleted_at"`
}

type HospitalizationListResponse struct {
	Data []HospitalizationResponse     `json:"data"`
	Meta HospitalizationPaginationMeta `json:"meta"`
}

type HospitalizationDeletedListResponse struct {
	Data []DeletedHospitalizationResponse `json:"data"`
	Meta HospitalizationPaginationMeta    `json:"meta"`
}

type HospitalizationPaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type CreateHospitalizationRequest struct {
	PatientID         uint   `json:"patient_id" binding:"required"`
	RoomID            uint   `json:"room_id" binding:"required"`
	AttendingDoctorID uint   `json:"attending_doctor_id" binding:"required"`
	AdmissionDate     string `json:"admission_date" binding:"required,datetime=2006-01-02T15:04:05Z"`
	AdmissionReason   string `json:"admission_reason" binding:"required"`
	Status            string `json:"status" binding:"omitempty"`
}

type UpdateHospitalizationRequest struct {
	RoomID            *uint   `json:"room_id" binding:"omitempty"`
	AttendingDoctorID *uint   `json:"attending_doctor_id" binding:"omitempty"`
	AdmissionReason   *string `json:"admission_reason" binding:"omitempty"`
}

type DischargeHospitalizationRequest struct {
	DischargeSummary string `json:"discharge_summary" binding:"omitempty"`
}

type TransferHospitalizationRequest struct {
	Notes string `json:"notes" binding:"omitempty"`
}

type HospitalizationPaginationQuery struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=id admission_date discharge_date status created_at"`
	SortDir   string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
	Search    string `form:"search" binding:"omitempty"`
	Status    string `form:"status" binding:"omitempty,oneof=admitted discharged transferred cancelled"`
	NotStatus string `form:"not_status" binding:"omitempty,oneof=admitted discharged transferred cancelled"`
	PatientID *uint  `form:"patient_id" binding:"omitempty"`
	RoomID    *uint  `form:"room_id" binding:"omitempty"`
	DoctorID  *uint  `form:"doctor_id" binding:"omitempty"`
}
