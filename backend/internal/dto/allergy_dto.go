package dto

import "time"

// ─── Allergy Response ──────────────────────────────────────────────────────────

type AllergyPatientResponse struct {
	ID          uint   `json:"id"`
	PatientCode string `json:"patient_code"`
	FullName    string `json:"full_name"`
}

type AllergyResponse struct {
	ID           uint                    `json:"id"`
	PatientID    uint                    `json:"patient_id"`
	AllergenType string                  `json:"allergen_type"`
	AllergenName string                  `json:"allergen_name"`
	Reaction     string                  `json:"reaction"`
	Severity     string                  `json:"severity"`
	Notes        string                  `json:"notes"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
	Patient      *AllergyPatientResponse `json:"patient,omitempty"`
}

type AllergyListResponse struct {
	Data []AllergyResponse     `json:"data"`
	Meta AllergyPaginationMeta `json:"meta"`
}

// ─── Allergy Request ───────────────────────────────────────────────────────────

type CreateAllergyRequest struct {
	PatientID    uint   `json:"patient_id" binding:"required"`
	AllergenType string `json:"allergen_type" binding:"required,oneof=food drug environmental insect latex other"`
	AllergenName string `json:"allergen_name" binding:"required"`
	Reaction     string `json:"reaction" binding:"required"`
	Severity     string `json:"severity" binding:"required,oneof=mild moderate severe"`
	Notes        string `json:"notes" binding:"omitempty"`
}

type UpdateAllergyRequest struct {
	AllergenType string `json:"allergen_type" binding:"omitempty,oneof=food drug environmental insect latex other"`
	AllergenName string `json:"allergen_name" binding:"omitempty"`
	Reaction     string `json:"reaction" binding:"omitempty"`
	Severity     string `json:"severity" binding:"omitempty,oneof=mild moderate severe"`
	Notes        string `json:"notes" binding:"omitempty"`
}

// ─── Allergy Pagination ────────────────────────────────────────────────────────

type AllergyPaginationQuery struct {
	Page      int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	PatientID *uint  `form:"patient_id" binding:"omitempty"`
	SortBy    string `form:"sort_by,default=created_at" binding:"omitempty,oneof=id created_at allergen_name"`
	SortDir   string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type AllergyPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
