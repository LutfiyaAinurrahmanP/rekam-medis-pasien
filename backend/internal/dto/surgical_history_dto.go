package dto

import "time"

// ─── Surgical History Response ─────────────────────────────────────────────────

type SurgicalHistoryPatientResponse struct {
	ID          uint   `json:"id"`
	PatientCode string `json:"patient_code"`
	FullName    string `json:"full_name"`
}

type SurgicalHistoryResponse struct {
	ID            uint                            `json:"id"`
	PatientID     uint                            `json:"patient_id"`
	ProcedureName string                          `json:"procedure_name"`
	SurgeryDate   string                          `json:"surgery_date"`
	SurgeonName   string                          `json:"surgeon_name"`
	Hospital      string                          `json:"hospital"`
	Complication  string                          `json:"complication"`
	Notes         string                          `json:"notes"`
	CreatedAt     time.Time                       `json:"created_at"`
	UpdatedAt     time.Time                       `json:"updated_at"`
	Patient       *SurgicalHistoryPatientResponse `json:"patient,omitempty"`
}

type SurgicalHistoryListResponse struct {
	Data []SurgicalHistoryResponse     `json:"data"`
	Meta SurgicalHistoryPaginationMeta `json:"meta"`
}

// ─── Surgical History Request ──────────────────────────────────────────────────

type CreateSurgicalHistoryRequest struct {
	PatientID     uint   `json:"patient_id" binding:"required"`
	ProcedureName string `json:"procedure_name" binding:"required"`
	SurgeryDate   string `json:"surgery_date" binding:"required,datetime=2006-01-02"`
	SurgeonName   string `json:"surgeon_name" binding:"omitempty"`
	Hospital      string `json:"hospital" binding:"omitempty"`
	Complication  string `json:"complication" binding:"omitempty"`
	Notes         string `json:"notes" binding:"omitempty"`
}

type UpdateSurgicalHistoryRequest struct {
	ProcedureName string `json:"procedure_name" binding:"omitempty"`
	SurgeryDate   string `json:"surgery_date" binding:"omitempty,datetime=2006-01-02"`
	SurgeonName   string `json:"surgeon_name" binding:"omitempty"`
	Hospital      string `json:"hospital" binding:"omitempty"`
	Complication  string `json:"complication" binding:"omitempty"`
	Notes         string `json:"notes" binding:"omitempty"`
}

// ─── Surgical History Pagination ───────────────────────────────────────────────

type SurgicalHistoryPaginationQuery struct {
	Page      int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	PatientID *uint  `form:"patient_id" binding:"omitempty"`
	SortBy    string `form:"sort_by,default=surgery_date" binding:"omitempty,oneof=id created_at surgery_date"`
	SortDir   string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type SurgicalHistoryPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
