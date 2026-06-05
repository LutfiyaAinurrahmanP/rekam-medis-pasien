package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

// MockPatientService is a mock implementation of PatientService interface
type MockPatientService struct {
	mock.Mock
}

func (m *MockPatientService) ListPatients(query *dto.PatientPaginationQuery) (*dto.PatientListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PatientListResponse), args.Error(1)
}

func (m *MockPatientService) DeleteListPatients(query *dto.PatientPaginationQuery) (*dto.PatientDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PatientDeletedListResponse), args.Error(1)
}

func (m *MockPatientService) GetPatientByID(id uint) (*dto.PatientResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PatientResponse), args.Error(1)
}


func (m *MockPatientService) GetMyPatientData(userID uint) (*dto.PatientResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PatientResponse), args.Error(1)
}

func (m *MockPatientService) UpdateMyPatientData(userID uint, req *dto.UpdatePatientRequest) (*dto.PatientResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PatientResponse), args.Error(1)
}

func (m *MockPatientService) CreatePatient(req *dto.CreatePatientRequest) (*dto.PatientResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PatientResponse), args.Error(1)
}

func (m *MockPatientService) UpdatePatient(id uint, req *dto.UpdatePatientRequest) (*dto.PatientResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PatientResponse), args.Error(1)
}

func (m *MockPatientService) SoftDeletePatient(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPatientService) RestorePatient(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPatientService) HardDeletePatient(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
