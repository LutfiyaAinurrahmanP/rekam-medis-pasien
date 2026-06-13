package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockMedicalRecordService struct {
	mock.Mock
}

func (m *MockMedicalRecordService) List(query *dto.MedicalRecordPaginationQuery) (*dto.MedicalRecordListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalRecordListResponse), args.Error(1)
}

func (m *MockMedicalRecordService) DeletedList(query *dto.MedicalRecordPaginationQuery) (*dto.MedicalRecordDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalRecordDeletedListResponse), args.Error(1)
}

func (m *MockMedicalRecordService) FindByID(id uint) (*dto.MedicalRecordResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalRecordResponse), args.Error(1)
}

func (m *MockMedicalRecordService) FindByIDUnscoped(id uint) (*dto.MedicalRecordResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalRecordResponse), args.Error(1)
}

func (m *MockMedicalRecordService) Create(req *dto.CreateMedicalRecordRequest) (*dto.MedicalRecordResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalRecordResponse), args.Error(1)
}

func (m *MockMedicalRecordService) Update(id uint, req *dto.UpdateMedicalRecordRequest) (*dto.MedicalRecordResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.MedicalRecordResponse), args.Error(1)
}

func (m *MockMedicalRecordService) Finalize(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicalRecordService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicalRecordService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMedicalRecordService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
