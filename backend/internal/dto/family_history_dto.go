package dto

import "time"

// ─── Family History Response ───────────────────────────────────────────────────

type FamilyHistoryPatientResponse struct {
	ID          uint   `json:"id"`
	PatientCode string `json:"patient_code"`
	FullName    string `json:"full_name"`
}

type FamilyHistoryResponse struct {
	ID            uint                          `json:"id"`
	PatientID     uint                          `json:"patient_id"`
	FamilyMember  string                        `json:"family_member"`
	ConditionName string                        `json:"condition_name"`
	Relation      string                        `json:"relation"`
	Notes         string                        `json:"notes"`
	CreatedAt     time.Time                     `json:"created_at"`
	UpdatedAt     time.Time                     `json:"updated_at"`
	Patient       *FamilyHistoryPatientResponse `json:"patient,omitempty"`
}

type FamilyHistoryListResponse struct {
	Data []FamilyHistoryResponse     `json:"data"`
	Meta FamilyHistoryPaginationMeta `json:"meta"`
}

// ─── Family History Request ────────────────────────────────────────────────────

type CreateFamilyHistoryRequest struct {
	PatientID     uint   `json:"patient_id" binding:"required"`
	FamilyMember  string `json:"family_member" binding:"required"`
	ConditionName string `json:"condition_name" binding:"required"`
	Relation      string `json:"relation" binding:"omitempty"`
	Notes         string `json:"notes" binding:"omitempty"`
}

type UpdateFamilyHistoryRequest struct {
	FamilyMember  string `json:"family_member" binding:"omitempty"`
	ConditionName string `json:"condition_name" binding:"omitempty"`
	Relation      string `json:"relation" binding:"omitempty"`
	Notes         string `json:"notes" binding:"omitempty"`
}

// ─── Family History Pagination ─────────────────────────────────────────────────

type FamilyHistoryPaginationQuery struct {
	Page      int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	PatientID *uint  `form:"patient_id" binding:"omitempty"`
	SortBy    string `form:"sort_by,default=created_at" binding:"omitempty,oneof=id created_at family_member"`
	SortDir   string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type FamilyHistoryPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
