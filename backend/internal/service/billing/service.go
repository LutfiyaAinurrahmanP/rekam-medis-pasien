package billing

import (
	"errors"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"gorm.io/gorm"
)

type BillingService interface {
	List(query dto.BillingPaginationQuery) (*dto.BillingListResponse, error)
	DeletedList(query dto.BillingPaginationQuery) (*dto.BillingDeletedListResponse, error)
	FindByID(id uint) (*dto.BillingResponse, error)
	FindByIDUnscoped(id uint) (*dto.BillingResponse, error)
	FindByInvoiceNumber(invoice string) (*dto.BillingResponse, error)
	FindByPatientID(patientID uint, query dto.BillingPaginationQuery) (*dto.BillingListResponse, error)
	
	Create(req dto.CreateBillingRequest) (*dto.BillingResponse, error)
	Update(id uint, req dto.UpdateBillingRequest) (*dto.BillingResponse, error)
	RecordPayment(id uint, req dto.RecordPaymentRequest) (*dto.BillingResponse, error)
	Cancel(id uint) (*dto.BillingResponse, error)
	
	Delete(id uint) error
	Restore(id uint) (*dto.BillingResponse, error)
	HardDelete(id uint) error

	// Items
	ListItems(billingID uint, query dto.BillingItemPaginationQuery) ([]dto.BillingItemResponse, error)
	FindItemByID(billingID, itemID uint) (*dto.BillingItemResponse, error)
	CreateItem(billingID uint, req dto.CreateBillingItemRequest) (*dto.BillingItemResponse, error)
	UpdateItem(billingID, itemID uint, req dto.UpdateBillingItemRequest) (*dto.BillingItemResponse, error)
	DeleteItem(billingID, itemID uint) error
}

type billingService struct {
	repo repository.BillingRepository
}

func NewBillingService(repo repository.BillingRepository) BillingService {
	return &billingService{repo: repo}
}

func mapToResponse(billing *models.Billing) *dto.BillingResponse {
	billing = nullOrEmpty(billing)
	if billing == nil {
		return nil
	}
	resp := &dto.BillingResponse{
		ID:                billing.ID,
		PatientID:         billing.PatientID,
		MedicalRecordID:   billing.MedicalRecordID,
		HospitalizationID: billing.HospitalizationID,
		InvoiceNumber:     billing.InvoiceNumber,
		BillingDate:       billing.BillingDate,
		DueDate:           billing.DueDate,
		TotalAmount:       billing.TotalAmount,
		PaidAmount:        billing.PaidAmount,
		RemainingAmount:   billing.RemainingAmount,
		Status:            billing.Status,
		PaymentMethod:     billing.PaymentMethod,
		Notes:             billing.Notes,
		CreatedAt:         billing.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         billing.UpdatedAt.Format(time.RFC3339),
	}

	if billing.PaymentMethod != nil {
		resp.PaymentMethod = billing.PaymentMethod
	}

	if billing.Patient != nil {
		resp.Patient = &dto.PatientRef{
			ID:                  billing.Patient.ID,
			Name:                billing.Patient.FullName,
			MedicalRecordNumber: billing.Patient.PatientCode,
		}
	}

	if billing.MedicalRecord != nil {
		visitDateStr := billing.MedicalRecord.VisitDate
		resp.MedicalRecord = &dto.MedicalRecordRef{
			ID:        billing.MedicalRecord.ID,
			VisitDate: visitDateStr,
		}
	}

	if len(billing.Items) > 0 {
		var items []dto.BillingItemResponse
		for _, item := range billing.Items {
			items = append(items, *mapItemToResponse(&item))
		}
		resp.Items = items
	}

	return resp
}

func nullOrEmpty(b *models.Billing) *models.Billing {
	if b == nil || b.ID == 0 {
		return nil
	}
	return b
}

func mapItemToResponse(item *models.BillingItem) *dto.BillingItemResponse {
	if item == nil || item.ID == 0 {
		return nil
	}
	return &dto.BillingItemResponse{
		ID:              item.ID,
		BillingID:       item.BillingID,
		ItemType:        item.ItemType,
		Description:     item.Description,
		Quantity:        item.Quantity,
		UnitPrice:       item.UnitPrice,
		TotalPrice:      item.TotalPrice,
		CreatedAt:       item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       item.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *billingService) List(query dto.BillingPaginationQuery) (*dto.BillingListResponse, error) {
	billings, meta, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	var resps []dto.BillingResponse
	for _, b := range billings {
		resps = append(resps, *mapToResponse(&b))
	}

	if len(resps) == 0 {
		resps = []dto.BillingResponse{}
	}

	return &dto.BillingListResponse{
		Status:  "success",
		Message: "Billing records retrieved successfully",
		Data:    resps,
		Meta:    meta,
	}, nil
}

func (s *billingService) DeletedList(query dto.BillingPaginationQuery) (*dto.BillingDeletedListResponse, error) {
	billings, meta, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	var resps []dto.BillingResponse
	for _, b := range billings {
		resps = append(resps, *mapToResponse(&b))
	}

	if len(resps) == 0 {
		resps = []dto.BillingResponse{}
	}

	return &dto.BillingDeletedListResponse{
		Status:  "success",
		Message: "Deleted billing records retrieved successfully",
		Data:    resps,
		Meta:    meta,
	}, nil
}

func (s *billingService) FindByID(id uint) (*dto.BillingResponse, error) {
	billing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("billing not found")
		}
		return nil, err
	}
	return mapToResponse(billing), nil
}

func (s *billingService) FindByIDUnscoped(id uint) (*dto.BillingResponse, error) {
	billing, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("billing not found")
		}
		return nil, err
	}
	return mapToResponse(billing), nil
}

func (s *billingService) FindByInvoiceNumber(invoice string) (*dto.BillingResponse, error) {
	billing, err := s.repo.FindByInvoiceNumber(invoice)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("billing not found")
		}
		return nil, err
	}
	return mapToResponse(billing), nil
}

func (s *billingService) FindByPatientID(patientID uint, query dto.BillingPaginationQuery) (*dto.BillingListResponse, error) {
	query.PatientID = &patientID
	return s.List(query)
}

func (s *billingService) Create(req dto.CreateBillingRequest) (*dto.BillingResponse, error) {
	billing := &models.Billing{
		PatientID:         req.PatientID,
		MedicalRecordID:   req.MedicalRecordID,
		HospitalizationID: req.HospitalizationID,
		InvoiceNumber:     req.InvoiceNumber,
		BillingDate:       req.BillingDate,
		DueDate:           req.DueDate,
		TotalAmount:       req.TotalAmount,
		Status:            req.Status,
	}
	if req.Notes != "" {
		billing.Notes = &req.Notes
	}

	if err := s.repo.Create(billing); err != nil {
		return nil, err
	}

	return s.FindByID(billing.ID)
}

func (s *billingService) Update(id uint, req dto.UpdateBillingRequest) (*dto.BillingResponse, error) {
	billing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if billing.Status == "cancelled" {
		return nil, errors.New("cannot update a cancelled billing")
	}

	if req.TotalAmount >= 0 {
		billing.TotalAmount = req.TotalAmount
	}
	if req.Notes != "" {
		billing.Notes = &req.Notes
	}

	if billing.Status != "cancelled" {
		netAmount := billing.TotalAmount
		if billing.PaidAmount >= netAmount && netAmount > 0 {
			billing.Status = "paid"
		} else if billing.PaidAmount > 0 {
			billing.Status = "partial"
		} else {
			billing.Status = "unpaid"
		}
	}

	if err := s.repo.Update(billing); err != nil {
		return nil, err
	}

	return s.FindByID(id)
}

func (s *billingService) RecordPayment(id uint, req dto.RecordPaymentRequest) (*dto.BillingResponse, error) {
	billing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if billing.Status == "cancelled" {
		return nil, errors.New("cannot record payment for a cancelled billing")
	}

	billing.PaidAmount += req.PaidAmount
	billing.PaymentMethod = &req.PaymentMethod

	netAmount := billing.TotalAmount
	if billing.PaidAmount >= netAmount {
		billing.Status = "paid"
	} else if billing.PaidAmount > 0 {
		billing.Status = "partial"
	}

	if err := s.repo.Update(billing); err != nil {
		return nil, err
	}

	return s.FindByID(id)
}

func (s *billingService) Cancel(id uint) (*dto.BillingResponse, error) {
	billing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	billing.Status = "cancelled"

	if err := s.repo.Update(billing); err != nil {
		return nil, err
	}

	return s.FindByID(id)
}

func (s *billingService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *billingService) Restore(id uint) (*dto.BillingResponse, error) {
	if err := s.repo.Restore(id); err != nil {
		return nil, err
	}
	return s.FindByID(id)
}

func (s *billingService) HardDelete(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *billingService) ListItems(billingID uint, query dto.BillingItemPaginationQuery) ([]dto.BillingItemResponse, error) {
	items, err := s.repo.ListItems(billingID, query)
	if err != nil {
		return nil, err
	}

	var resps []dto.BillingItemResponse
	for _, item := range items {
		resps = append(resps, *mapItemToResponse(&item))
	}

	if len(resps) == 0 {
		resps = []dto.BillingItemResponse{}
	}

	return resps, nil
}

func (s *billingService) FindItemByID(billingID, itemID uint) (*dto.BillingItemResponse, error) {
	item, err := s.repo.FindItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if item.BillingID != billingID {
		return nil, errors.New("billing item does not belong to this billing")
	}
	return mapItemToResponse(item), nil
}

func (s *billingService) CreateItem(billingID uint, req dto.CreateBillingItemRequest) (*dto.BillingItemResponse, error) {
	billing, err := s.repo.FindByID(billingID)
	if err != nil {
		return nil, err
	}
	if billing.Status == "cancelled" {
		return nil, errors.New("Cannot add item to a cancelled billing")
	}

	totalPrice := float64(req.Quantity) * req.UnitPrice

	item := &models.BillingItem{
		BillingID:       billingID,
		Description: req.Description,
		Quantity:        req.Quantity,
		UnitPrice:       req.UnitPrice,
		TotalPrice:      totalPrice,
	}
	if req.ItemType != "" {
		item.ItemType = &req.ItemType
	}
	billing.RemainingAmount = billing.TotalAmount - billing.PaidAmount
	if err := s.repo.CreateItem(item); err != nil {
		return nil, err
	}

	// Update billing total amount
	billing.TotalAmount += totalPrice
	s.repo.Update(billing)

	return s.FindItemByID(billingID, item.ID)
}

func (s *billingService) UpdateItem(billingID, itemID uint, req dto.UpdateBillingItemRequest) (*dto.BillingItemResponse, error) {
	item, err := s.repo.FindItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if item.BillingID != billingID {
		return nil, errors.New("billing item does not belong to this billing")
	}

	billing, err := s.repo.FindByID(billingID)
	if err != nil {
		return nil, err
	}
	if billing.Status == "cancelled" {
		return nil, errors.New("cannot update item in a cancelled billing")
	}

	oldTotalPrice := item.TotalPrice

	if req.Quantity > 0 {
		item.Quantity = req.Quantity
	}
	if req.UnitPrice >= 0 {
		item.UnitPrice = req.UnitPrice
	}

	item.TotalPrice = float64(item.Quantity) * item.UnitPrice

	if err := s.repo.UpdateItem(item); err != nil {
		return nil, err
	}

	// Update billing total amount
	billing.TotalAmount = billing.TotalAmount - oldTotalPrice + item.TotalPrice
	s.repo.Update(billing)

	return s.FindItemByID(billingID, item.ID)
}

func (s *billingService) DeleteItem(billingID, itemID uint) error {
	item, err := s.repo.FindItemByID(itemID)
	if err != nil {
		return err
	}
	if item.BillingID != billingID {
		return errors.New("billing item does not belong to this billing")
	}

	billing, err := s.repo.FindByID(billingID)
	if err != nil {
		return err
	}
	if billing.Status == "cancelled" {
		return errors.New("cannot delete item from a cancelled billing")
	}

	if err := s.repo.DeleteItem(itemID); err != nil {
		return err
	}

	// Update billing total amount
	billing.TotalAmount -= item.TotalPrice
	if billing.TotalAmount < 0 {
		billing.TotalAmount = 0
	}
	s.repo.Update(billing)

	return nil
}
