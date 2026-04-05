package dto

import "time"

type MedicineListResponse struct {
	Data []MedicineResponse     `json:"data"`
	Meta MedicinePaginationMeta `json:"meta"`
}

type MedicineDeletedListResponse struct {
	Data []DeletedMedicineResponse `json:"data"`
	Meta MedicinePaginationMeta    `json:"meta"`
}

type MedicineAvailableResponse struct {
	TotalAvailable  int64                  `json:"total_available"`
	TotalStockValue float64                `json:"total_stock_value"`
	MedicineTypes   []MedicineTypeCount    `json:"medicine_types"`
	Data            []MedicineResponse     `json:"data"`
	Meta            MedicinePaginationMeta `json:"meta"`
}

type MedicineTypeCount struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type MedicineLowStockResponse struct {
	Threshold     int                    `json:"threshold"`
	TotalLowStock int64                  `json:"total_low_stock"`
	CriticalCount int64                  `json:"critical_count"`
	Data          []MedicineLowStockItem `json:"data"`
	Meta          MedicinePaginationMeta `json:"meta"`
}

type MedicineLowStockItem struct {
	ID                uint    `json:"id"`
	Name              string  `json:"name"`
	GenericName       string  `json:"generic_name"`
	Type              string  `json:"type"`
	StockQuantity     int     `json:"stock_quantity"`
	MinimumStock      int     `json:"minimum_stock"`
	Price             float64 `json:"price"`
	StockStatus       string  `json:"stock_status"`
	DaysUntilStockout int     `json:"days_until_stockout"`
	IsActive          bool    `json:"is_active"`
}

type MedicineOutOfStockResponse struct {
	TotalOutOfStock int64                    `json:"total_out_of_stock"`
	Data            []MedicineOutOfStockItem `json:"data"`
	Meta            MedicinePaginationMeta   `json:"meta"`
}

type MedicineOutOfStockItem struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	GenericName    string     `json:"generic_name"`
	Type           string     `json:"type"`
	StockQuantity  int        `json:"stock_quantity"`
	Price          float64    `json:"price"`
	LastStockDate  *time.Time `json:"last_stock_date"`
	DaysOutOfStock int        `json:"days_out_of_stock"`
	IsActive       bool       `json:"is_active"`
}

type MedicineInactiveResponse struct {
	Data []MedicineInactiveItem `json:"data"`
	Meta MedicinePaginationMeta `json:"meta"`
}

type MedicineInactiveItem struct {
	ID                uint      `json:"id"`
	Name              string    `json:"name"`
	GenericName       string    `json:"generic_name"`
	Type              string    `json:"type"`
	StockQuantity     int       `json:"stock_quantity"`
	IsActive          bool      `json:"is_active"`
	DeactivatedReason string    `json:"deactivated_reason,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type MedicineByTypeResponse struct {
	Type               string                 `json:"type"`
	TotalMedicines     int64                  `json:"total_medicines"`
	AvailableMedicines int64                  `json:"available_medicines"`
	TotalStock         int64                  `json:"total_stock"`
	Data               []MedicineResponse     `json:"data"`
	Meta               MedicinePaginationMeta `json:"meta"`
}

type MedicineSearchResponse struct {
	SearchCriteria MedicineSearchCriteria `json:"search_criteria"`
	ResultsFound   int64                  `json:"results_found"`
	Data           []MedicineSearchItem   `json:"data"`
	Meta           MedicinePaginationMeta `json:"meta"`
}

type MedicineSearchCriteria struct {
	Keyword      string   `json:"keyword,omitempty"`
	Type         string   `json:"type,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	MinPrice     *float64 `json:"min_price,omitempty"`
	MaxPrice     *float64 `json:"max_price,omitempty"`
	HasStock     *bool    `json:"has_stock,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
}

type MedicineSearchItem struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name"`
	GenericName    string  `json:"generic_name"`
	BrandName      string  `json:"brand_name"`
	Type           string  `json:"type"`
	Strength       string  `json:"strength"`
	StockQuantity  int     `json:"stock_quantity"`
	Price          float64 `json:"price"`
	IsActive       bool    `json:"is_active"`
	RelevanceScore float64 `json:"relevance_score"`
}

type CreateMedicineRequest struct {
	Name          string  `json:"name" binding:"required,min=1,max=200"`
	GenericName   string  `json:"generic_name" binding:"omitempty,max=200"`
	BrandName     string  `json:"brand_name" binding:"omitempty,max=200"`
	Type          string  `json:"type" binding:"required,oneof=tablet capsule syrup injection ointment other"`
	Strength      string  `json:"strength" binding:"omitempty,max=50"`
	Manufacturer  string  `json:"manufacturer" binding:"omitempty,max=100"`
	Unit          string  `json:"unit" binding:"omitempty,max=20"`
	StockQuantity int     `json:"stock_quantity" binding:"omitempty,min=0"`
	Price         float64 `json:"price" binding:"omitempty,min=0"`
	IsActive      *bool   `json:"is_active" binding:"omitempty"`
}

type UpdateMedicineRequest struct {
	Name         *string  `json:"name" binding:"omitempty,min=1,max=200"`
	GenericName  *string  `json:"generic_name" binding:"omitempty,max=200"`
	BrandName    *string  `json:"brand_name" binding:"omitempty,max=200"`
	Type         *string  `json:"type" binding:"omitempty,oneof=tablet capsule syrup injection ointment other"`
	Strength     *string  `json:"strength" binding:"omitempty,max=50"`
	Manufacturer *string  `json:"manufacturer" binding:"omitempty,max=100"`
	Unit         *string  `json:"unit" binding:"omitempty,max=20"`
	Price        *float64 `json:"price" binding:"omitempty,min=0"`
	IsActive     *bool    `json:"is_active" binding:"omitempty"`
}

type MedicineResponse struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	GenericName   string    `json:"generic_name"`
	BrandName     string    `json:"brand_name"`
	Type          string    `json:"type"`
	Strength      string    `json:"strength"`
	Manufacturer  string    `json:"manufacturer"`
	Unit          string    `json:"unit"`
	StockQuantity int       `json:"stock_quantity"`
	Price         float64   `json:"price"`
	IsActive      bool      `json:"is_active"`
	IsLowStock    bool      `json:"is_low_stock"`
	IsOutOfStock  bool      `json:"is_out_of_stock"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type MedicineDetailResponse struct {
	ID            uint               `json:"id"`
	Name          string             `json:"name"`
	GenericName   string             `json:"generic_name"`
	BrandName     string             `json:"brand_name"`
	Type          string             `json:"type"`
	Strength      string             `json:"strength"`
	Manufacturer  string             `json:"manufacturer"`
	Unit          string             `json:"unit"`
	StockQuantity int                `json:"stock_quantity"`
	Price         float64            `json:"price"`
	IsActive      bool               `json:"is_active"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	StockInfo     MedicineStockInfo  `json:"stock_info"`
	Statistics    MedicineStatistics `json:"statistics"`
	Details       MedicineDetails    `json:"details"`
}

type MedicineStockInfo struct {
	IsLowStock        bool       `json:"is_low_stock"`
	IsOutOfStock      bool       `json:"is_out_of_stock"`
	MinimumStockLevel int        `json:"minimum_stock_level"`
	ReorderLevel      int        `json:"reorder_level"`
	LastRestockDate   *time.Time `json:"last_restock_date"`
	AverageDailyUsage float64    `json:"average_daily_usage"`
}

type MedicineStatistics struct {
	TotalPrescribed     int64   `json:"total_prescribed"`
	PrescribedThisMonth int64   `json:"prescribed_this_month"`
	TotalDispensed      int64   `json:"total_dispensed"`
	DispensedThisMonth  int64   `json:"dispensed_this_month"`
	RevenueThisMonth    float64 `json:"revenue_this_month"`
}

type MedicineDetails struct {
	Indication     string `json:"indication"`
	DosageForm     string `json:"dosage_form"`
	Route          string `json:"route"`
	Storage        string `json:"storage"`
	ExpiryTracking bool   `json:"expiry_tracking"`
}

type DeletedMedicineResponse struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	GenericName   string     `json:"generic_name"`
	BrandName     string     `json:"brand_name"`
	Type          string     `json:"type"`
	Strength      string     `json:"strength"`
	Manufacturer  string     `json:"manufacturer"`
	Unit          string     `json:"unit"`
	StockQuantity int        `json:"stock_quantity"`
	Price         float64    `json:"price"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type AddStockRequest struct {
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Notes    string `json:"notes" binding:"omitempty,max=500"`
}

type ReduceStockRequest struct {
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Notes    string `json:"notes" binding:"omitempty,max=500"`
}

type StockOperationResponse struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	PreviousStock   int       `json:"previous_stock"`
	ChangedQuantity int       `json:"changed_quantity"`
	CurrentStock    int       `json:"current_stock"`
	StockValue      float64   `json:"stock_value"`
	IsLowStock      bool      `json:"is_low_stock"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type MedicinePaginationQuery struct {
	Page         int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize     int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	Search       string `form:"search" binding:"omitempty"`
	Type         string `form:"type" binding:"omitempty,oneof=tablet capsule syrup injection ointment other"`
	Manufacturer string `form:"manufacturer" binding:"omitempty"`
	IsActive     *bool  `form:"is_active" binding:"omitempty"`
	HasStock     *bool  `form:"has_stock" binding:"omitempty"`
	MinStock     int    `form:"min_stock" binding:"omitempty,min=0"`
	MaxStock     int    `form:"max_stock" binding:"omitempty,min=0"`
	SortBy       string `form:"sort_by,default=name" binding:"omitempty,oneof=name generic_name stock_quantity price created_at"`
	SortDir      string `form:"sort_dir,default=asc" binding:"omitempty,oneof=asc desc"`
}

type MedicinePaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
