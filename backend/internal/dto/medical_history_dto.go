package dto

import "time"

type MedicalHistoryOverviewResponse struct {
	ID                 uint      `json:"id"`
	PatientID          uint      `json:"patient_id"`
	PatientName        string    `json:"patient_name"`
	AllergiesCount     int       `json:"allergies_count"`
	ConditionsCount    int       `json:"conditions_count"`
	SurgeriesCount     int       `json:"surgeries_count"`
	FamilyHistoryCount int       `json:"family_history_count"`
	LastUpdated        time.Time `json:"last_updated"`
}

type MedicalHistoryDetailResponse struct {
	ID                uint                       `json:"id"`
	PatientID         uint                       `json:"patient_id"`
	Allergies         []AllergyResponse          `json:"allergies"`
	MedicalConditions []MedicalConditionResponse `json:"medical_conditions"`
	SurgicalHistories []SurgicalHistoryResponse  `json:"surgical_history"`
	FamilyHistories   []FamilyHistoryResponse    `json:"family_history"`
}

type MedicalHistoryPaginationQuery struct {
	Page      int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	PatientID *uint  `form:"patient_id" binding:"omitempty"`
	SortBy    string `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at updated_at patient_code full_name"`
	SortDir   string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type MedicalHistoryPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type MedicalHistoryListResponse struct {
	Data []MedicalHistoryOverviewResponse `json:"data"`
	Meta MedicalHistoryPaginationMeta     `json:"meta"`
}
