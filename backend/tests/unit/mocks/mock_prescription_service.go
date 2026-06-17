package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockPrescriptionService struct {
	mock.Mock
}

func (m *MockPrescriptionService) List(query *dto.PrescriptionPaginationQuery) (*dto.PrescriptionListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionListResponse), args.Error(1)
}

func (m *MockPrescriptionService) DeletedList(query *dto.PrescriptionPaginationQuery) (*dto.PrescriptionDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionDeletedListResponse), args.Error(1)
}

func (m *MockPrescriptionService) FindByID(id uint) (*dto.PrescriptionResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionResponse), args.Error(1)
}

func (m *MockPrescriptionService) FindByIDUnscoped(id uint) (*dto.PrescriptionResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionResponse), args.Error(1)
}

func (m *MockPrescriptionService) Create(req *dto.CreatePrescriptionRequest) (*dto.PrescriptionResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionResponse), args.Error(1)
}

func (m *MockPrescriptionService) Update(id uint, req *dto.UpdatePrescriptionRequest) (*dto.PrescriptionResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionResponse), args.Error(1)
}

func (m *MockPrescriptionService) Dispense(id uint) (*dto.PrescriptionResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionResponse), args.Error(1)
}

func (m *MockPrescriptionService) Cancel(id uint) (*dto.PrescriptionResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionResponse), args.Error(1)
}

func (m *MockPrescriptionService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPrescriptionService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPrescriptionService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPrescriptionService) ListItems(prescriptionID uint) ([]dto.PrescriptionItemResponse, error) {
	args := m.Called(prescriptionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.PrescriptionItemResponse), args.Error(1)
}

func (m *MockPrescriptionService) FindItemByID(prescriptionID, itemID uint) (*dto.PrescriptionItemResponse, error) {
	args := m.Called(prescriptionID, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionItemResponse), args.Error(1)
}

func (m *MockPrescriptionService) CreateItem(prescriptionID uint, req *dto.CreatePrescriptionItemRequest) (*dto.PrescriptionItemResponse, error) {
	args := m.Called(prescriptionID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionItemResponse), args.Error(1)
}

func (m *MockPrescriptionService) UpdateItem(prescriptionID, itemID uint, req *dto.UpdatePrescriptionItemRequest) (*dto.PrescriptionItemResponse, error) {
	args := m.Called(prescriptionID, itemID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PrescriptionItemResponse), args.Error(1)
}

func (m *MockPrescriptionService) DeleteItem(prescriptionID, itemID uint) error {
	args := m.Called(prescriptionID, itemID)
	return args.Error(0)
}
