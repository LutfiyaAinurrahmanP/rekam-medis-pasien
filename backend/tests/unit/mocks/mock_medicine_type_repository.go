package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockMedicineTypeRepository struct {
	mock.Mock
}

func (m *MockMedicineTypeRepository) List(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.MedicineType), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineTypeRepository) DeletedList(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.MedicineType), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineTypeRepository) ActiveList(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.MedicineType), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineTypeRepository) InactiveList(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.MedicineType), args.Get(1).(int64), args.Error(2)
}

func (m *MockMedicineTypeRepository) FindByID(id uint) (*models.MedicineType, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MedicineType), args.Error(1)
}

func (m *MockMedicineTypeRepository) Create(medicineType *models.MedicineType) error {
	args := m.Called(medicineType)
	return args.Error(0)
}

func (m *MockMedicineTypeRepository) Update(medicineType *models.MedicineType) error {
	args := m.Called(medicineType)
	return args.Error(0)
}

func (m *MockMedicineTypeRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineTypeRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineTypeRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineTypeRepository) Activate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineTypeRepository) Deactivate(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicineTypeRepository) IsNameExists(name string, excludeID ...uint) (bool, error) {
	args := m.Called(name, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockMedicineTypeRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	args := m.Called(code, excludeID)
	return args.Bool(0), args.Error(1)
}
