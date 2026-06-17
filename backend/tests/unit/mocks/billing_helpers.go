package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestBillingWithData(id uint, patientID uint) *models.Billing {
	now := time.Now()
	paymentMethod := "cash"

	return &models.Billing{
		ID:              id,
		PatientID:       patientID,
		InvoiceNumber:   "INV-2024-0001",
		BillingDate:     "2024-01-01",
		DueDate:         "2024-01-14",
		TotalAmount:     100000,
		PaidAmount:      50000,
		RemainingAmount: 50000,
		Status:          "partial",
		PaymentMethod:   &paymentMethod,
		CreatedAt:       now,
		UpdatedAt:       now,
		Patient: &models.Patient{
			ID:          patientID,
			FullName:    "Budi Santoso",
			PatientCode: "RM-001",
		},
	}
}

func NewTestBillingList(count int) []models.Billing {
	var list []models.Billing
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestBillingWithData(uint(i), 1))
	}
	return list
}

func NewTestBillingItemWithData(id uint, billingID uint) *models.BillingItem {
	now := time.Now()
	itemType := "consultation"

	return &models.BillingItem{
		ID:          id,
		BillingID:   billingID,
		ItemType:    &itemType,
		Description: "Konsultasi Dokter Umum",
		Quantity:    1,
		UnitPrice:   100000,
		TotalPrice:  100000,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewTestBillingItemList(count int) []models.BillingItem {
	var list []models.BillingItem
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestBillingItemWithData(uint(i), 1))
	}
	return list
}

func NewBillingPaginationQuery(page, pageSize int) *dto.BillingPaginationQuery {
	return &dto.BillingPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewTestBillingResponse(b *models.Billing) *dto.BillingResponse {
	return &dto.BillingResponse{
		ID:              b.ID,
		PatientID:       b.PatientID,
		InvoiceNumber:   b.InvoiceNumber,
		BillingDate:     b.BillingDate,
		DueDate:         b.DueDate,
		TotalAmount:     b.TotalAmount,
		PaidAmount:      b.PaidAmount,
		RemainingAmount: b.RemainingAmount,
		Status:          b.Status,
		PaymentMethod:   b.PaymentMethod,
		CreatedAt:       b.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       b.UpdatedAt.Format(time.RFC3339),
	}
}

func NewTestBillingItemResponse(item *models.BillingItem) *dto.BillingItemResponse {
	return &dto.BillingItemResponse{
		ID:          item.ID,
		BillingID:   item.BillingID,
		ItemType:    item.ItemType,
		Description: item.Description,
		Quantity:    item.Quantity,
		UnitPrice:   item.UnitPrice,
		TotalPrice:  item.TotalPrice,
		CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   item.UpdatedAt.Format(time.RFC3339),
	}
}

func NewCreateBillingRequest(patientID uint) *dto.CreateBillingRequest {
	return &dto.CreateBillingRequest{
		PatientID:     patientID,
		InvoiceNumber: "INV-2024-0001",
		BillingDate:   "2024-01-01",
		DueDate:       "2024-01-14",
		TotalAmount:   100000,
		Status:        "pending",
	}
}

func NewUpdateBillingRequest() *dto.UpdateBillingRequest {
	return &dto.UpdateBillingRequest{
		TotalAmount: 150000,
		Notes:       "Updated",
	}
}

func NewRecordPaymentRequest() *dto.RecordPaymentRequest {
	return &dto.RecordPaymentRequest{
		PaidAmount:    50000,
		PaymentMethod: "cash",
	}
}

func NewCreateBillingItemRequest() *dto.CreateBillingItemRequest {
	return &dto.CreateBillingItemRequest{
		ItemType:    "consultation",
		Description: "Konsultasi",
		Quantity:    1,
		UnitPrice:   100000,
	}
}

func NewUpdateBillingItemRequest() *dto.UpdateBillingItemRequest {
	return &dto.UpdateBillingItemRequest{
		Quantity:  2,
		UnitPrice: 100000,
	}
}
