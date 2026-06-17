package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockDoctorService struct {
	mock.Mock
}

func (m *MockDoctorService) GetMyDoctorData(userID uint) (*dto.DoctorResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorResponse), args.Error(1)
}

func (m *MockDoctorService) UpdateMyDoctorData(userID uint, req *dto.UpdateDoctorRequest) (*dto.DoctorResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorResponse), args.Error(1)
}

func (m *MockDoctorService) ListDoctors(query *dto.DoctorPaginationQuery) (*dto.DoctorListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorListResponse), args.Error(1)
}

func (m *MockDoctorService) DeletedListDoctors(query *dto.DoctorPaginationQuery) (*dto.DoctorDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorDeletedListResponse), args.Error(1)
}

func (m *MockDoctorService) GetDoctorByID(id uint) (*dto.DoctorResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorResponse), args.Error(1)
}

func (m *MockDoctorService) ActiveList(query *dto.DoctorPaginationQuery) (*dto.DoctorListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorListResponse), args.Error(1)
}

func (m *MockDoctorService) InactiveList(query *dto.DoctorPaginationQuery) (*dto.DoctorListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorListResponse), args.Error(1)
}

func (m *MockDoctorService) CreateDoctor(req *dto.CreateDoctorRequest) (*dto.DoctorResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorResponse), args.Error(1)
}

func (m *MockDoctorService) UpdateDoctor(id uint, req *dto.UpdateDoctorRequest) (*dto.DoctorResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorResponse), args.Error(1)
}

func (m *MockDoctorService) ActivateDoctor(id uint) (*dto.DoctorResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorResponse), args.Error(1)
}

func (m *MockDoctorService) DeactivateDoctor(id uint) (*dto.DoctorResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DoctorResponse), args.Error(1)
}

func (m *MockDoctorService) SoftDeleteDoctor(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorService) RestoreDoctor(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDoctorService) HardDeleteDoctor(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
