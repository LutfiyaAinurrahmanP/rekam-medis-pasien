package dto

import "time"

type PatientListResponse struct {
	Data []PatientResponse     `json:"data"`
	Meta PatientPaginationMeta `json:"meta"`
}

type PatientDeletedListResponse struct {
	Data []DeletedPatientResponse     `json:"data"`
	Meta PatientPaginationMeta `json:"meta"`
}

type CreatePatientRequest struct {
	UserID                 uint   `json:"user_id" binding:"omitempty"`
	PatientCode            string  `json:"patient_code" binding:"required,min=1,max=20"`
	FullName               string  `json:"full_name" binding:"required,min=1,max=100"`
	DateOfBirth            string  `json:"date_of_birth" binding:"required"`
	Gender                 string  `json:"gender" binding:"required,oneof=male female other"`
	BloodType              string `json:"blood_type" binding:"omitempty,max=5"`
	Phone                  string `json:"phone" binding:"omitempty,max=15"`
	Email                  string `json:"email" binding:"omitempty,email"`
	Address                string `json:"address" binding:"omitempty"`
	EmergencyContactName   string `json:"emergency_contact_name" binding:"omitempty,max=100"`
	EmergencyContactPhone  string `json:"emergency_contact_phone" binding:"omitempty,max=15"`
	InsuranceNumber        string `json:"insurance_number" binding:"omitempty,max=50"`
	InsuranceProvider      string `json:"insurance_provider" binding:"omitempty,max=100"`
	Allergies              string `json:"allergies" binding:"omitempty"`
}

type UpdatePatientRequest struct {
	FullName               *string `json:"full_name" binding:"omitempty,min=1,max=100"`
	DateOfBirth            *string `json:"date_of_birth" binding:"omitempty"`
	Gender                 *string `json:"gender" binding:"omitempty,oneof=male female other"`
	BloodType              *string `json:"blood_type" binding:"omitempty,max=5"`
	Phone                  *string `json:"phone" binding:"omitempty,max=15"`
	Email                  *string `json:"email" binding:"omitempty,email"`
	Address                *string `json:"address" binding:"omitempty"`
	EmergencyContactName   *string `json:"emergency_contact_name" binding:"omitempty,max=100"`
	EmergencyContactPhone  *string `json:"emergency_contact_phone" binding:"omitempty,max=15"`
	InsuranceNumber        *string `json:"insurance_number" binding:"omitempty,max=50"`
	InsuranceProvider      *string `json:"insurance_provider" binding:"omitempty,max=100"`
	Allergies              *string `json:"allergies" binding:"omitempty"`
}

type PatientResponse struct {
	ID                     uint    `json:"id"`
	UserID                 uint   `json:"user_id"`
	PatientCode            string  `json:"patient_code"`
	FullName               string  `json:"full_name"`
	DateOfBirth            string  `json:"date_of_birth"`
	Age                    int    `json:"age"`
	Gender                 string  `json:"gender"`
	BloodType              string `json:"blood_type"`
	Phone                  string `json:"phone"`
	Email                  string `json:"email"`
	Address                string `json:"address"`
	EmergencyContactName   string `json:"emergency_contact_name"`
	EmergencyContactPhone  string `json:"emergency_contact_phone"`
	InsuranceNumber        string `json:"insurance_number"`
	InsuranceProvider      string `json:"insurance_provider"`
	Allergies              string `json:"allergies"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	// MedicalRecordsCount    int    `json:"medical_records_count"`
	// AppointmentsCount      int    `json:"appointments_count"`
	// LastVisit              time.Time `json:"last_visit"`
}

type DeletedPatientResponse struct {
	ID                     uint    `json:"id"`
	UserID                 uint   `json:"user_id"`
	PatientCode            string  `json:"patient_code"`
	FullName               string  `json:"full_name"`
	DateOfBirth            string  `json:"date_of_birth"`
	Age                    int    `json:"age"`
	Gender                 string  `json:"gender"`
	BloodType              string `json:"blood_type"`
	Phone                  string `json:"phone"`
	Email                  string `json:"email"`
	Address                string `json:"address"`
	EmergencyContactName   string `json:"emergency_contact_name"`
	EmergencyContactPhone  string `json:"emergency_contact_phone"`
	InsuranceNumber        string `json:"insurance_number"`
	InsuranceProvider      string `json:"insurance_provider"`
	Allergies              string `json:"allergies"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at"`
}

type PatientPaginationQuery struct {
	Page               int     `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize           int     `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	Search             string `form:"search" binding:"omitempty"`
	Gender             string `form:"gender" binding:"omitempty"`
	BloodType          string `form:"blood_type" binding:"omitempty"`
	InsuranceProvider  string `form:"insurance_provider" binding:"omitempty"`
	MinAge             int    `form:"min_age" binding:"omitempty"`
	MaxAge             int    `form:"max_age" binding:"omitempty"`
	SortBy             string `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at full_name patient_code date_of_birth"`
	SortDir            string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type PatientPaginationMeta struct {
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}