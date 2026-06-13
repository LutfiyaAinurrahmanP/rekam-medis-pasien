package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockBillingService struct {
	mock.Mock
}

func (m *MockBillingService) List(query dto.BillingPaginationQuery) (*dto.BillingListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingListResponse), args.Error(1)
}

func (m *MockBillingService) DeletedList(query dto.BillingPaginationQuery) (*dto.BillingDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingDeletedListResponse), args.Error(1)
}

func (m *MockBillingService) FindByID(id uint) (*dto.BillingResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingResponse), args.Error(1)
}

func (m *MockBillingService) FindByIDUnscoped(id uint) (*dto.BillingResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingResponse), args.Error(1)
}

func (m *MockBillingService) FindByInvoiceNumber(invoice string) (*dto.BillingResponse, error) {
	args := m.Called(invoice)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingResponse), args.Error(1)
}

func (m *MockBillingService) FindByPatientID(patientID uint, query dto.BillingPaginationQuery) (*dto.BillingListResponse, error) {
	args := m.Called(patientID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingListResponse), args.Error(1)
}

func (m *MockBillingService) Create(req dto.CreateBillingRequest) (*dto.BillingResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingResponse), args.Error(1)
}

func (m *MockBillingService) Update(id uint, req dto.UpdateBillingRequest) (*dto.BillingResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingResponse), args.Error(1)
}

func (m *MockBillingService) RecordPayment(id uint, req dto.RecordPaymentRequest) (*dto.BillingResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingResponse), args.Error(1)
}

func (m *MockBillingService) Cancel(id uint) (*dto.BillingResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingResponse), args.Error(1)
}

func (m *MockBillingService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBillingService) Restore(id uint) (*dto.BillingResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingResponse), args.Error(1)
}

func (m *MockBillingService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBillingService) ListItems(billingID uint, query dto.BillingItemPaginationQuery) ([]dto.BillingItemResponse, error) {
	args := m.Called(billingID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.BillingItemResponse), args.Error(1)
}

func (m *MockBillingService) FindItemByID(billingID, itemID uint) (*dto.BillingItemResponse, error) {
	args := m.Called(billingID, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingItemResponse), args.Error(1)
}

func (m *MockBillingService) CreateItem(billingID uint, req dto.CreateBillingItemRequest) (*dto.BillingItemResponse, error) {
	args := m.Called(billingID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingItemResponse), args.Error(1)
}

func (m *MockBillingService) UpdateItem(billingID, itemID uint, req dto.UpdateBillingItemRequest) (*dto.BillingItemResponse, error) {
	args := m.Called(billingID, itemID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.BillingItemResponse), args.Error(1)
}

func (m *MockBillingService) DeleteItem(billingID, itemID uint) error {
	args := m.Called(billingID, itemID)
	return args.Error(0)
}
