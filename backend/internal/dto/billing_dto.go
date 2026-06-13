package dto

type BillingResponse struct {
	ID                uint                 `json:"id"`
	PatientID         uint                 `json:"patient_id"`
	MedicalRecordID   *uint                `json:"medical_record_id"`
	HospitalizationID *uint                `json:"hospitalization_id"`
	InvoiceNumber     string               `json:"invoice_number"`
	BillingDate       string               `json:"billing_date"`
	DueDate           string               `json:"due_date"`
	TotalAmount       float64              `json:"total_amount"`
	PaidAmount        float64              `json:"paid_amount"`
	RemainingAmount   float64              `json:"remaining_amount"`
	Status            string               `json:"status"`
	PaymentMethod     *string              `json:"payment_method"`
	Notes             *string              `json:"notes"`
	CreatedAt         string               `json:"created_at"`
	UpdatedAt         string               `json:"updated_at"`
	Patient           *PatientRef          `json:"patient,omitempty"`
	MedicalRecord     *MedicalRecordRef    `json:"medical_record,omitempty"`
	Items             []BillingItemResponse `json:"items,omitempty"`
}

type BillingItemResponse struct {
	ID              uint    `json:"id"`
	BillingID       uint    `json:"billing_id"`
	ItemType        *string `json:"item_type"`
	Description     string  `json:"description"`
	Quantity        int     `json:"quantity"`
	UnitPrice       float64 `json:"unit_price"`
	TotalPrice      float64 `json:"total_price"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type PatientRef struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	MedicalRecordNumber string `json:"medical_record_number"`
}

type MedicalRecordRef struct {
	ID        uint   `json:"id"`
	VisitDate string `json:"visit_date"`
}

type BillingListResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []BillingResponse      `json:"data"`
	Meta    BillingPaginationMeta  `json:"meta"`
}

type BillingDeletedListResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []BillingResponse      `json:"data"`
	Meta    BillingPaginationMeta  `json:"meta"`
}

type BillingPaginationQuery struct {
	Page          int    `form:"page" binding:"omitempty,min=1"`
	PageSize      int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy        string `form:"sort_by" binding:"omitempty"`
	SortDir       string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
	Search        string `form:"search" binding:"omitempty"`
	Status        string `form:"status" binding:"omitempty,oneof=pending partial paid cancelled"`
	PaymentMethod string `form:"payment_method" binding:"omitempty"`
	PatientID     *uint  `form:"-"`
}

type BillingPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type CreateBillingRequest struct {
	PatientID         uint    `json:"patient_id" binding:"required"`
	MedicalRecordID   *uint   `json:"medical_record_id" binding:"omitempty"`
	HospitalizationID *uint   `json:"hospitalization_id" binding:"omitempty"`
	InvoiceNumber     string  `json:"invoice_number" binding:"required,max=50"`
	BillingDate       string  `json:"billing_date" binding:"required"`
	DueDate           string  `json:"due_date" binding:"required"`
	TotalAmount       float64 `json:"total_amount" binding:"required,min=0"`
	Status            string  `json:"status" binding:"required,oneof=pending partial paid cancelled"`
	Notes             string  `json:"notes" binding:"omitempty"`
}

type UpdateBillingRequest struct {
	TotalAmount    float64 `json:"total_amount" binding:"omitempty,min=0"`
	Notes          string  `json:"notes" binding:"omitempty"`
}

type RecordPaymentRequest struct {
	PaidAmount    float64 `json:"paid_amount" binding:"required,min=0"`
	PaymentMethod string  `json:"payment_method" binding:"required,oneof=cash debit_card credit_card insurance transfer other"`
}

type CreateBillingItemRequest struct {
	ItemType        string  `json:"item_type" binding:"required,oneof=consultation medicine lab_test procedure room other"`
	Description     string  `json:"description" binding:"required,max=200"`
	Quantity        int     `json:"quantity" binding:"required,min=1"`
	UnitPrice       float64 `json:"unit_price" binding:"required,min=0"`
}

type UpdateBillingItemRequest struct {
	Quantity  int     `json:"quantity" binding:"omitempty,min=1"`
	UnitPrice float64 `json:"unit_price" binding:"omitempty,min=0"`
}

type BillingItemPaginationQuery struct {
	ItemType string `form:"item_type" binding:"omitempty,oneof=consultation medicine lab_test procedure room other"`
}
