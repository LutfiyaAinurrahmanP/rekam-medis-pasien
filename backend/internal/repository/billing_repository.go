package repository

import (
	"errors"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type BillingRepository interface {
	List(query dto.BillingPaginationQuery) ([]models.Billing, dto.BillingPaginationMeta, error)
	DeletedList(query dto.BillingPaginationQuery) ([]models.Billing, dto.BillingPaginationMeta, error)
	FindByID(id uint) (*models.Billing, error)
	FindByIDUnscoped(id uint) (*models.Billing, error)
	FindByInvoiceNumber(invoice string) (*models.Billing, error)
	FindByPatientID(patientID uint, query dto.BillingPaginationQuery) ([]models.Billing, dto.BillingPaginationMeta, error)
	Create(billing *models.Billing) error
	Update(billing *models.Billing) error
	Delete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	
	// Items
	ListItems(billingID uint, query dto.BillingItemPaginationQuery) ([]models.BillingItem, error)
	FindItemByID(itemID uint) (*models.BillingItem, error)
	CreateItem(item *models.BillingItem) error
	UpdateItem(item *models.BillingItem) error
	DeleteItem(itemID uint) error
}

type billingRepository struct {
	db *gorm.DB
}

func NewBillingRepository(db *gorm.DB) BillingRepository {
	return &billingRepository{db: db}
}

func (r *billingRepository) baseQuery() *gorm.DB {
	return r.db.Model(&models.Billing{}).
		Preload("Patient").
		Preload("MedicalRecord").
		Preload("Items")
}

func (r *billingRepository) List(query dto.BillingPaginationQuery) ([]models.Billing, dto.BillingPaginationMeta, error) {
	var billings []models.Billing
	var totalItems int64

	q := r.baseQuery()

	if query.PatientID != nil {
		q = q.Where("patient_id = ?", *query.PatientID)
	}
	if query.Status != "" {
		q = q.Where("status = ?", query.Status)
	}
	if query.PaymentMethod != "" {
		q = q.Where("payment_method = ?", query.PaymentMethod)
	}
	if query.Search != "" {
		search := "%" + query.Search + "%"
		q = q.Where("invoice_number LIKE ? OR notes LIKE ?", search, search)
	}

	q.Count(&totalItems)

	if query.SortBy != "" {
		sortDir := "ASC"
		if query.SortDir == "desc" {
			sortDir = "DESC"
		}
		q = q.Order(query.SortBy + " " + sortDir)
	} else {
		q = q.Order("created_at DESC")
	}

	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		q = q.Offset(offset).Limit(query.PageSize)
	}

	if err := q.Find(&billings).Error; err != nil {
		return nil, dto.BillingPaginationMeta{}, err
	}

	totalPages := 0
	if query.PageSize > 0 {
		totalPages = int((totalItems + int64(query.PageSize) - 1) / int64(query.PageSize))
	}

	meta := dto.BillingPaginationMeta{
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return billings, meta, nil
}

func (r *billingRepository) DeletedList(query dto.BillingPaginationQuery) ([]models.Billing, dto.BillingPaginationMeta, error) {
	var billings []models.Billing
	var totalItems int64

	q := r.baseQuery().Unscoped().Where("deleted_at IS NOT NULL")

	if query.PatientID != nil {
		q = q.Where("patient_id = ?", *query.PatientID)
	}
	if query.Status != "" {
		q = q.Where("status = ?", query.Status)
	}
	if query.PaymentMethod != "" {
		q = q.Where("payment_method = ?", query.PaymentMethod)
	}
	if query.Search != "" {
		search := "%" + query.Search + "%"
		q = q.Where("invoice_number LIKE ? OR notes LIKE ?", search, search)
	}

	q.Count(&totalItems)

	if query.SortBy != "" {
		sortDir := "ASC"
		if query.SortDir == "desc" {
			sortDir = "DESC"
		}
		q = q.Order(query.SortBy + " " + sortDir)
	} else {
		q = q.Order("deleted_at DESC")
	}

	if query.Page > 0 && query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		q = q.Offset(offset).Limit(query.PageSize)
	}

	if err := q.Find(&billings).Error; err != nil {
		return nil, dto.BillingPaginationMeta{}, err
	}

	totalPages := 0
	if query.PageSize > 0 {
		totalPages = int((totalItems + int64(query.PageSize) - 1) / int64(query.PageSize))
	}

	meta := dto.BillingPaginationMeta{
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return billings, meta, nil
}

func (r *billingRepository) FindByID(id uint) (*models.Billing, error) {
	var billing models.Billing
	if err := r.baseQuery().First(&billing, id).Error; err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *billingRepository) FindByIDUnscoped(id uint) (*models.Billing, error) {
	var billing models.Billing
	if err := r.baseQuery().Unscoped().First(&billing, id).Error; err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *billingRepository) FindByInvoiceNumber(invoice string) (*models.Billing, error) {
	var billing models.Billing
	if err := r.baseQuery().Where("invoice_number = ?", invoice).First(&billing).Error; err != nil {
		return nil, err
	}
	return &billing, nil
}

func (r *billingRepository) FindByPatientID(patientID uint, query dto.BillingPaginationQuery) ([]models.Billing, dto.BillingPaginationMeta, error) {
	query.PatientID = &patientID
	return r.List(query)
}

func (r *billingRepository) Create(billing *models.Billing) error {
	return r.db.Create(billing).Error
}

func (r *billingRepository) Update(billing *models.Billing) error {
	return r.db.Save(billing).Error
}

func (r *billingRepository) Delete(id uint) error {
	return r.db.Delete(&models.Billing{}, id).Error
}

func (r *billingRepository) Restore(id uint) error {
	return r.db.Unscoped().Model(&models.Billing{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *billingRepository) HardDelete(id uint) error {
	// Delete items first to avoid foreign key constraints issues
	if err := r.db.Unscoped().Where("billing_id = ?", id).Delete(&models.BillingItem{}).Error; err != nil {
		return err
	}
	return r.db.Unscoped().Delete(&models.Billing{}, id).Error
}

func (r *billingRepository) ListItems(billingID uint, query dto.BillingItemPaginationQuery) ([]models.BillingItem, error) {
	var items []models.BillingItem
	q := r.db.Where("billing_id = ?", billingID)
	
	if query.ItemType != "" {
		q = q.Where("item_type = ?", query.ItemType)
	}
	
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *billingRepository) FindItemByID(itemID uint) (*models.BillingItem, error) {
	var item models.BillingItem
	if err := r.db.First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("billing item not found")
		}
		return nil, err
	}
	return &item, nil
}

func (r *billingRepository) CreateItem(item *models.BillingItem) error {
	return r.db.Create(item).Error
}

func (r *billingRepository) UpdateItem(item *models.BillingItem) error {
	return r.db.Save(item).Error
}

func (r *billingRepository) DeleteItem(itemID uint) error {
	return r.db.Delete(&models.BillingItem{}, itemID).Error
}
