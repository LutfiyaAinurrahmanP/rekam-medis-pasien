package dto

import (
	"time"
)

type MedicalRecordPatientResponse struct {
	ID          uint   `json:"id"`
	PatientCode string `json:"patient_code"`
	FullName    string `json:"full_name"`
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
	BloodType   string `json:"blood_type,omitempty"`
	Allergies   string `json:"allergies,omitempty"`
}

type MedicalRecordDoctorResponse struct {
	ID             uint   `json:"id"`
	FullName       string `json:"full_name"`
	Specialization string `json:"specialization,omitempty"`
	Department     string `json:"department,omitempty"`
}

type VitalSignsResponse struct {
	BloodPressure      string  `json:"blood_pressure,omitempty"`
	HeartRate          int     `json:"heart_rate,omitempty"`
	RespiratoryRate    int     `json:"respiratory_rate,omitempty"`
	Temperature        float64 `json:"temperature,omitempty"`
	OxygenSaturation   int     `json:"oxygen_saturation,omitempty"`
	Weight             float64 `json:"weight,omitempty"`
	Height             float64 `json:"height,omitempty"`
	BMI                float64 `json:"bmi,omitempty"`
}

type MedicalHistorySummaryResponse struct {
	Allergies         []any `json:"allergies"`
	MedicalConditions []any `json:"medical_conditions"`
	SurgicalHistory   []any `json:"surgical_history"`
	FamilyHistory     []any `json:"family_history"`
}

type MedicalRecordResponse struct {
	ID                      uint                          `json:"id"`
	PatientID               uint                          `json:"patient_id"`
	Patient                 *MedicalRecordPatientResponse `json:"patient,omitempty"`
	DoctorID                uint                          `json:"doctor_id"`
	Doctor                  *MedicalRecordDoctorResponse  `json:"doctor,omitempty"`
	AppointmentID       *uint                         `json:"appointment_id,omitempty"`
	VisitDate           string                        `json:"visit_date"`
	ChiefComplaint      string                         `json:"chief_complaint"`
	HistoryOfIllness    string                         `json:"history_of_illness,omitempty"`
	PhysicalExamination string                         `json:"physical_examination,omitempty"`
	VitalSigns          *VitalSignsResponse            `json:"vital_signs,omitempty"`
	MedicalHistory      *MedicalHistorySummaryResponse `json:"medical_history,omitempty"`
	Diagnosis           string                         `json:"diagnosis"`
	TreatmentPlan       string                        `json:"treatment_plan"`
	Notes               string                        `json:"notes,omitempty"`
	Status              string                        `json:"status"`
	CreatedAt           time.Time                     `json:"created_at"`
	UpdatedAt           time.Time                     `json:"updated_at"`
}

type DeletedMedicalRecordResponse struct {
	MedicalRecordResponse
	DeletedAt *time.Time `json:"deleted_at"`
}

type MedicalRecordListResponse struct {
	Data []MedicalRecordResponse     `json:"data"`
	Meta MedicalRecordPaginationMeta `json:"meta"`
}

type MedicalRecordDeletedListResponse struct {
	Data []DeletedMedicalRecordResponse `json:"data"`
	Meta MedicalRecordPaginationMeta    `json:"meta"`
}

type CreateMedicalRecordRequest struct {
	PatientID           uint   `json:"patient_id" binding:"required"`
	AppointmentID       *uint  `json:"appointment_id" binding:"omitempty"`
	VisitDate           string `json:"visit_date" binding:"required,datetime=2006-01-02"`
	ChiefComplaint      string `json:"chief_complaint" binding:"required"`
	HistoryOfIllness    string `json:"history_of_illness" binding:"omitempty"`
	PhysicalExamination string `json:"physical_examination" binding:"omitempty"`
	Diagnosis           string `json:"diagnosis" binding:"required"`
	TreatmentPlan       string `json:"treatment_plan" binding:"required"`
	Notes               string `json:"notes" binding:"omitempty"`
	DoctorID            uint   `json:"-"` // set from auth
}

type UpdateMedicalRecordRequest struct {
	ChiefComplaint      *string `json:"chief_complaint" binding:"omitempty"`
	HistoryOfIllness    *string `json:"history_of_illness" binding:"omitempty"`
	PhysicalExamination *string `json:"physical_examination" binding:"omitempty"`
	Diagnosis           *string `json:"diagnosis" binding:"omitempty"`
	TreatmentPlan       *string `json:"treatment_plan" binding:"omitempty"`
	Notes               *string `json:"notes" binding:"omitempty"`
}

type MedicalRecordPaginationQuery struct {
	Page         int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize     int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	PatientID    *uint  `form:"patient_id" binding:"omitempty"`
	DoctorID     *uint  `form:"doctor_id" binding:"omitempty"`
	DepartmentID *uint  `form:"department_id" binding:"omitempty"`
	Status       string `form:"status" binding:"omitempty,oneof=draft finalized amended"`
	DateFrom     string `form:"date_from" binding:"omitempty,datetime=2006-01-02"`
	DateTo       string `form:"date_to" binding:"omitempty,datetime=2006-01-02"`
	SortBy       string `form:"sort_by,default=visit_date" binding:"omitempty,oneof=created_at visit_date"`
	SortDir      string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type MedicalRecordPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
