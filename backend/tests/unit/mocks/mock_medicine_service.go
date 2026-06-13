package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockMedicineService struct {
	mock.Mock
}

func (m *MockMedicineService) List(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineListResponse), args.Error(1)
}

func (m *MockMedicineService) DeletedList(query *dto.MedicinePaginationQuery) (*dto.MedicineDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineDeletedListResponse), args.Error(1)
}

func (m *MockMedicineService) AvailableList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineListResponse), args.Error(1)
}

func (m *MockMedicineService) LowStockList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineListResponse), args.Error(1)
}

func (m *MockMedicineService) OutStockList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineListResponse), args.Error(1)
}

func (m *MockMedicineService) ActiveList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineListResponse), args.Error(1)
}

func (m *MockMedicineService) InactiveList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineListResponse), args.Error(1)
}

func (m *MockMedicineService) FindByID(id uint) (*dto.MedicineResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineResponse), args.Error(1)
}

func (m *MockMedicineService) FindByIDUnscoped(id uint) (*dto.DeletedMedicineResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DeletedMedicineResponse), args.Error(1)
}

func (m *MockMedicineService) Create(req *dto.CreateMedicineRequest) (*dto.MedicineResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineResponse), args.Error(1)
}

func (m *MockMedicineService) Update(id uint, req *dto.UpdateMedicineRequest) (*dto.MedicineResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicineResponse), args.Error(1)
}

func (m *MockMedicineService) AddStock(id uint, req *dto.AddStockRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockMedicineService) ReduceStock(id uint, req *dto.ReduceStockRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockMedicineService) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineService) Deactivate(id uint, req *dto.DeactivateMedicineRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockMedicineService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
