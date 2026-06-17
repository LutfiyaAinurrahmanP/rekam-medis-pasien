package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockBillingRepository struct {
	mock.Mock
}

func (m *MockBillingRepository) List(query dto.BillingPaginationQuery) ([]models.Billing, dto.BillingPaginationMeta, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.BillingPaginationMeta), args.Error(2)
	}
	return args.Get(0).([]models.Billing), args.Get(1).(dto.BillingPaginationMeta), args.Error(2)
}

func (m *MockBillingRepository) DeletedList(query dto.BillingPaginationQuery) ([]models.Billing, dto.BillingPaginationMeta, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.BillingPaginationMeta), args.Error(2)
	}
	return args.Get(0).([]models.Billing), args.Get(1).(dto.BillingPaginationMeta), args.Error(2)
}

func (m *MockBillingRepository) FindByID(id uint) (*models.Billing, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Billing), args.Error(1)
}

func (m *MockBillingRepository) FindByIDUnscoped(id uint) (*models.Billing, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Billing), args.Error(1)
}

func (m *MockBillingRepository) FindByInvoiceNumber(invoice string) (*models.Billing, error) {
	args := m.Called(invoice)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Billing), args.Error(1)
}

func (m *MockBillingRepository) FindByPatientID(patientID uint, query dto.BillingPaginationQuery) ([]models.Billing, dto.BillingPaginationMeta, error) {
	args := m.Called(patientID, query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.BillingPaginationMeta), args.Error(2)
	}
	return args.Get(0).([]models.Billing), args.Get(1).(dto.BillingPaginationMeta), args.Error(2)
}

func (m *MockBillingRepository) Create(billing *models.Billing) error {
	args := m.Called(billing)
	return args.Error(0)
}

func (m *MockBillingRepository) Update(billing *models.Billing) error {
	args := m.Called(billing)
	return args.Error(0)
}

func (m *MockBillingRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBillingRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBillingRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Items
func (m *MockBillingRepository) ListItems(billingID uint, query dto.BillingItemPaginationQuery) ([]models.BillingItem, error) {
	args := m.Called(billingID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.BillingItem), args.Error(1)
}

func (m *MockBillingRepository) FindItemByID(itemID uint) (*models.BillingItem, error) {
	args := m.Called(itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BillingItem), args.Error(1)
}

func (m *MockBillingRepository) CreateItem(item *models.BillingItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockBillingRepository) UpdateItem(item *models.BillingItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockBillingRepository) DeleteItem(itemID uint) error {
	args := m.Called(itemID)
	return args.Error(0)
}
