package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockMedicineRepository struct {
	mock.Mock
}

func (m *MockMedicineRepository) List(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Medicine), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineRepository) DeletedList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Medicine), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineRepository) AvailableList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Medicine), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineRepository) LowStockList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Medicine), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineRepository) OutStockList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Medicine), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineRepository) ActiveList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Medicine), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineRepository) InactiveList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Medicine), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineRepository) FindByID(id uint) (*models.Medicine, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Medicine), args.Error(1)
}

func (m *MockMedicineRepository) FindByIDUnscoped(id uint) (*models.Medicine, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Medicine), args.Error(1)
}

func (m *MockMedicineRepository) Create(medicine *models.Medicine) error {
	args := m.Called(medicine)
	return args.Error(0)
}

func (m *MockMedicineRepository) Update(medicine *models.Medicine) error {
	args := m.Called(medicine)
	return args.Error(0)
}

func (m *MockMedicineRepository) AddStock(id uint, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

func (m *MockMedicineRepository) ReduceStock(id uint, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

func (m *MockMedicineRepository) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineRepository) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
