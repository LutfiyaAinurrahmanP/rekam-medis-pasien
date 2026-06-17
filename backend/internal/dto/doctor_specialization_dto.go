package dto

import "time"

type DoctorSpecializationResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeletedDoctorSpecializationResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Description string     `json:"description"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type CreateDoctorSpecializationRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Code        string `json:"code" binding:"omitempty,max=20"`
	Description string `json:"description" binding:"omitempty"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

type UpdateDoctorSpecializationRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=100"`
	Code        *string `json:"code" binding:"omitempty,max=20"`
	Description *string `json:"description" binding:"omitempty"`
	IsActive    *bool   `json:"is_active" binding:"omitempty"`
}

type DoctorSpecializationListResponse struct {
	Data []DoctorSpecializationResponse     `json:"data"`
	Meta DoctorSpecializationPaginationMeta `json:"meta"`
}

type DoctorSpecializationDeletedListResponse struct {
	Data []DeletedDoctorSpecializationResponse `json:"data"`
	Meta DoctorSpecializationPaginationMeta    `json:"meta"`
}

type DoctorSpecializationPaginationQuery struct {
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search" binding:"omitempty"`
	SortBy   string `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at name code"`
	SortDir  string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type DoctorSpecializationPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
