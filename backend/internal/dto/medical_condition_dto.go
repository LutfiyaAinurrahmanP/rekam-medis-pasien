package dto

import "time"

// ─── Medical Condition Response ────────────────────────────────────────────────

type MedicalConditionPatientResponse struct {
	ID          uint   `json:"id"`
	PatientCode string `json:"patient_code"`
	FullName    string `json:"full_name"`
}

type MedicalConditionResponse struct {
	ID            uint                             `json:"id"`
	PatientID     uint                             `json:"patient_id"`
	ConditionName string                           `json:"condition_name"`
	ICDCode       string                           `json:"icd_code"`
	DiagnosedDate string                           `json:"diagnosed_date"`
	Status        string                           `json:"status"`
	Notes         string                           `json:"notes"`
	CreatedAt     time.Time                        `json:"created_at"`
	UpdatedAt     time.Time                        `json:"updated_at"`
	Patient       *MedicalConditionPatientResponse `json:"patient,omitempty"`
}

type MedicalConditionListResponse struct {
	Data []MedicalConditionResponse     `json:"data"`
	Meta MedicalConditionPaginationMeta `json:"meta"`
}

// ─── Medical Condition Request ─────────────────────────────────────────────────

type CreateMedicalConditionRequest struct {
	PatientID     uint   `json:"patient_id" binding:"required"`
	ConditionName string `json:"condition_name" binding:"required"`
	ICDCode       string `json:"icd_code" binding:"omitempty"`
	DiagnosedDate string `json:"diagnosed_date" binding:"omitempty,datetime=2006-01-02"`
	Status        string `json:"status" binding:"omitempty,oneof=ongoing resolved managed"`
	Notes         string `json:"notes" binding:"omitempty"`
}

type UpdateMedicalConditionRequest struct {
	ConditionName string `json:"condition_name" binding:"omitempty"`
	ICDCode       string `json:"icd_code" binding:"omitempty"`
	DiagnosedDate string `json:"diagnosed_date" binding:"omitempty,datetime=2006-01-02"`
	Status        string `json:"status" binding:"omitempty,oneof=ongoing resolved managed"`
	Notes         string `json:"notes" binding:"omitempty"`
}

// ─── Medical Condition Pagination ──────────────────────────────────────────────

type MedicalConditionPaginationQuery struct {
	Page      int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	PatientID *uint  `form:"patient_id" binding:"omitempty"`
	Status    string `form:"status" binding:"omitempty,oneof=ongoing resolved managed"`
	SortBy    string `form:"sort_by,default=created_at" binding:"omitempty,oneof=id created_at condition_name"`
	SortDir   string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type MedicalConditionPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
