package dto

import "time"

// ============================================================
// Request DTOs
// ============================================================

type CreateTypeTestRequest struct {
	Name        string   `json:"name" binding:"required,min=1,max=200"`
	Code        string   `json:"code" binding:"required,min=1,max=50"`
	Category    *string  `json:"category" binding:"omitempty,max=100"`
	Description *string  `json:"description" binding:"omitempty"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	IsActive    *bool    `json:"is_active" binding:"omitempty"`
}

type UpdateTypeTestRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=1,max=200"`
	Code        *string  `json:"code" binding:"omitempty,min=1,max=50"`
	Category    *string  `json:"category" binding:"omitempty,max=100"`
	Description *string  `json:"description" binding:"omitempty"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	IsActive    *bool    `json:"is_active" binding:"omitempty"`
}

// ============================================================
// Query DTOs
// ============================================================

type TypeTestPaginationQuery struct {
	Page     int      `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int      `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	Search   string   `form:"search" binding:"omitempty"`
	Category string   `form:"category" binding:"omitempty"`
	IsActive *bool    `form:"is_active" binding:"omitempty"`
	MinPrice *float64 `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice *float64 `form:"max_price" binding:"omitempty,min=0"`
	SortBy   string   `form:"sort_by,default=name" binding:"omitempty,oneof=name code category price created_at"`
	SortDir  string   `form:"sort_dir,default=asc" binding:"omitempty,oneof=asc desc"`
}

type TypeTestSearchQuery struct {
	Keyword  string   `form:"keyword" binding:"omitempty"`
	Category string   `form:"category" binding:"omitempty"`
	MinPrice *float64 `form:"min_price" binding:"omitempty,min=0"`
	MaxPrice *float64 `form:"max_price" binding:"omitempty,min=0"`
	IsActive *bool    `form:"is_active" binding:"omitempty"`
	Page     int      `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int      `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
}

// ============================================================
// Response DTOs
// ============================================================

type TypeTestResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DeletedTypeTestResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

// ============================================================
// List Response DTOs
// ============================================================

type TypeTestPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type TypeTestListResponse struct {
	Data []TypeTestResponse     `json:"data"`
	Meta TypeTestPaginationMeta `json:"meta"`
}

type TypeTestDeletedListResponse struct {
	Data []DeletedTypeTestResponse `json:"data"`
	Meta TypeTestPaginationMeta    `json:"meta"`
}

type ActiveTypeTestCategorySummary struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type ActiveTypeTestListResponse struct {
	TotalActiveTests int64                           `json:"total_active_tests"`
	Categories       []ActiveTypeTestCategorySummary `json:"categories"`
	Data             []TypeTestResponse              `json:"data"`
	Meta             TypeTestPaginationMeta          `json:"meta"`
}

type TypeTestSearchResponse struct {
	SearchCriteria map[string]interface{} `json:"search_criteria"`
	ResultsFound   int64                  `json:"results_found"`
	Data           []TypeTestResponse     `json:"data"`
	Meta           TypeTestPaginationMeta `json:"meta"`
}
