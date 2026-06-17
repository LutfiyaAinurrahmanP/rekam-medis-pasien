package dto

import (
	"time"
)

type VitalSignMedicalRecord struct {
	ID             uint   `json:"id"`
	VisitDate      string `json:"visit_date"`
	ChiefComplaint string `json:"chief_complaint"`
}

type VitalSignResponse struct {
	ID                     uint                    `json:"id"`
	MedicalRecordID        uint                    `json:"medical_record_id"`
	BloodPressureSystolic  *int                    `json:"blood_pressure_systolic"`
	BloodPressureDiastolic *int                    `json:"blood_pressure_diastolic"`
	HeartRate              *int                    `json:"heart_rate"`
	Temperature            *float64                `json:"temperature"`
	RespiratoryRate        *int                    `json:"respiratory_rate"`
	OxygenSaturation       *float64                `json:"oxygen_saturation"`
	WeightKG               *float64                `json:"weight_kg"`
	HeightCM               *int                    `json:"height_cm"`
	BMI                    *float64                `json:"bmi"`
	RecordedAt             string                  `json:"recorded_at"`
	CreatedAt              time.Time               `json:"created_at"`
	UpdatedAt              time.Time               `json:"updated_at"`
	MedicalRecord          *VitalSignMedicalRecord `json:"medical_record,omitempty"`
}

type DeletedVitalSignResponse struct {
	VitalSignResponse
	DeletedAt *time.Time `json:"deleted_at"`
}

type VitalSignListResponse struct {
	Data []VitalSignResponse     `json:"data"`
	Meta VitalSignPaginationMeta `json:"meta"`
}

type VitalSignDeletedListResponse struct {
	Data []DeletedVitalSignResponse `json:"data"`
	Meta VitalSignPaginationMeta    `json:"meta"`
}

type CreateVitalSignRequest struct {
	MedicalRecordID        uint     `json:"medical_record_id" binding:"required"`
	RecordedAt             string   `json:"recorded_at" binding:"required,datetime=2006-01-02T15:04:05Z07:00"`
	BloodPressureSystolic  *int     `json:"blood_pressure_systolic" binding:"omitempty"`
	BloodPressureDiastolic *int     `json:"blood_pressure_diastolic" binding:"omitempty"`
	HeartRate              *int     `json:"heart_rate" binding:"omitempty"`
	Temperature            *float64 `json:"temperature" binding:"omitempty"`
	RespiratoryRate        *int     `json:"respiratory_rate" binding:"omitempty"`
	OxygenSaturation       *float64 `json:"oxygen_saturation" binding:"omitempty"`
	WeightKG               *float64 `json:"weight_kg" binding:"omitempty"`
	HeightCM               *int     `json:"height_cm" binding:"omitempty"`
}

type UpdateVitalSignRequest struct {
	BloodPressureSystolic  *int     `json:"blood_pressure_systolic" binding:"omitempty"`
	BloodPressureDiastolic *int     `json:"blood_pressure_diastolic" binding:"omitempty"`
	HeartRate              *int     `json:"heart_rate" binding:"omitempty"`
	Temperature            *float64 `json:"temperature" binding:"omitempty"`
	RespiratoryRate        *int     `json:"respiratory_rate" binding:"omitempty"`
	OxygenSaturation       *float64 `json:"oxygen_saturation" binding:"omitempty"`
	WeightKG               *float64 `json:"weight_kg" binding:"omitempty"`
	HeightCM               *int     `json:"height_cm" binding:"omitempty"`
}

type VitalSignPaginationQuery struct {
	Page            int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize        int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	MedicalRecordID *uint  `form:"medical_record_id" binding:"omitempty"`
	SortBy          string `form:"sort_by,default=recorded_at" binding:"omitempty,oneof=id recorded_at created_at"`
	SortDir         string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type VitalSignPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
