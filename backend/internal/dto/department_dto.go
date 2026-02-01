package dto

import "time"

type DepartmentListResponse struct {
	Data []DepartmentResponse     `json:"data"`
	Meta DepartmentPaginationMeta `json:"meta"`
}

type DepartmentDeletedListResponse struct {
	Data []DeletedDepartmentResponse `json:"data"`
	Meta DepartmentPaginationMeta    `json:"meta"`
}

type CreateDepartmentRequest struct {
	Name          string `json:"name" binding:"required,min=3,max=100"`
	Code          string `json:"code" binding:"required,min=1,max=20"`
	Description   string `json:"description" binding:"omitempty"`
	FloorLocation string `json:"floor_location" binding:"omitempty,max=50"`
}

type UpdateDepartmentRequest struct {
	Name          *string `json:"name" binding:"omitempty,min=3,max=100"`
	Code          *string `json:"code" binding:"omitempty,min=1,max=20"`
	Description   *string `json:"description" binding:"omitempty"`
	FloorLocation *string `json:"floor_location" binding:"omitempty,max=50"`
}

type DepartmentResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	Description   string    `json:"description"`
	FloorLocation string    `json:"floor_location"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DeletedDepartmentResponse struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	Code          string     `json:"code"`
	Description   string     `json:"description"`
	FloorLocation string     `json:"floor_location"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type DepartmentPaginationQuery struct {
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search" binding:"omitempty"`
	SortBy   string `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at name code"`
	SortDir  string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type DepartmentPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
