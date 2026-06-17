package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockMedicineTypeService struct {
	mock.Mock
}

func (m *MockMedicineTypeService) List(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineTypeListResponse), args.Error(1)
}

func (m *MockMedicineTypeService) DeletedList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineTypeDeletedListResponse), args.Error(1)
}

func (m *MockMedicineTypeService) ActiveList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineTypeListResponse), args.Error(1)
}

func (m *MockMedicineTypeService) InactiveList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineTypeListResponse), args.Error(1)
}

func (m *MockMedicineTypeService) FindByID(id uint) (*dto.MedicineTypeResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineTypeResponse), args.Error(1)
}

func (m *MockMedicineTypeService) Create(req *dto.CreateMedicineTypeRequest) (*dto.MedicineTypeResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineTypeResponse), args.Error(1)
}

func (m *MockMedicineTypeService) Update(id uint, req *dto.UpdateMedicineTypeRequest) (*dto.MedicineTypeResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineTypeResponse), args.Error(1)
}

func (m *MockMedicineTypeService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineTypeService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineTypeService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineTypeService) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineTypeService) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
