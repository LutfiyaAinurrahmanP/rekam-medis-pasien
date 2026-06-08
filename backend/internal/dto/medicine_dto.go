package dto

import "time"

type MedicineResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	GenericName    string    `json:"generic_name"`
	BrandName      string    `json:"brand_name"`
	MedicineTypeID uint      `json:"medicine_type_id"`
	Strength       string    `json:"strength"`
	Manufacturer   string    `json:"manufacturer"`
	Unit           string    `json:"unit"`
	StockQuantity  int       `json:"stock_quantity"`
	Price          float64   `json:"price"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type DeletedMedicineResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	GenericName    string    `json:"generic_name"`
	BrandName      string    `json:"brand_name"`
	MedicineTypeID uint      `json:"medicine_type_id"`
	Strength       string    `json:"strength"`
	Manufacturer   string    `json:"manufacturer"`
	Unit           string    `json:"unit"`
	StockQuantity  int       `json:"stock_quantity"`
	Price          float64   `json:"price"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type MedicineListResponse struct {
	Data []MedicineResponse     `json:"data"`
	Meta MedicinePaginationMeta `json:"meta"`
}

type MedicineDeletedListResponse struct {
	Data []DeletedMedicineResponse `json:"data"`
	Meta MedicinePaginationMeta   `json:"meta"`
}

type CreateMedicineRequest struct {
	Name          string  `json:"name" binding:"required,min=1,max=200"`
	GenericName   string  `json:"generic_name" binding:"omitempty,max=200"`
	BrandName     string  `json:"brand_name" binding:"omitempty,max=200"`
	MedicineTypeID *uint  `json:"medicine_type_id" binding:"required"`
	Strength      string  `json:"strength" binding:"omitempty,max=50"`
	Manufacturer  string  `json:"manufacturer" binding:"omitempty,max=100"`
	Unit          string  `json:"unit" binding:"omitempty,max=20"`
	StockQuantity *int     `json:"stock_quantity" binding:"required,min=1"`
	Price         *float64 `json:"price" binding:"omitempty,min=0"`
	IsActive      *bool   `json:"is_active" binding:"omitempty"`
}

type UpdateMedicineRequest struct {
	Name         *string  `json:"name" binding:"omitempty,min=1,max=200"`
	GenericName  *string  `json:"generic_name" binding:"omitempty,max=200"`
	BrandName    *string  `json:"brand_name" binding:"omitempty,max=200"`
	MedicineTypeID *uint   `json:"medicine_type_id" binding:"omitempty"`
	Strength     *string  `json:"strength" binding:"omitempty,max=50"`
	Manufacturer *string  `json:"manufacturer" binding:"omitempty,max=100"`
	Unit         *string  `json:"unit" binding:"omitempty,max=20"`
	Price        *float64 `json:"price" binding:"omitempty,min=0"`
}

type ActivateMedicineRequest struct {
	IsActive *bool `json:"is_active" binding:"omitempty"`
}

type DeactivateMedicineRequest struct {
	IsActive *bool `json:"is_active" binding:"omitempty"`
}

type AddStockRequest struct {
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type ReduceStockRequest struct {
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

type MedicinePaginationQuery struct {
	Page     int `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize int `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	Search string `form:"search" binding:"omitempty"`
	IsActive       *bool  `form:"is_active" binding:"omitempty"`
	MedicineTypeID *uint  `form:"medicine_type_id" binding:"omitempty"`
	StockStatus    string `form:"stock_status" binding:"omitempty,oneof=available out_of_stock low_stock"`
	SortBy  string `form:"sort_by,default=name" binding:"omitempty,oneof=name generic_name brand_name manufacturer stock_quantity price created_at updated_at"`
	SortDir string `form:"sort_dir,default=asc" binding:"omitempty,oneof=asc desc"`
}

type MedicinePaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
