package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockReferralService struct {
	mock.Mock
}

func (m *MockReferralService) List(query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralListResponse), args.Error(1)
}

func (m *MockReferralService) DeletedList(query dto.ReferralPaginationQuery) (*dto.ReferralDeletedListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralDeletedListResponse), args.Error(1)
}

func (m *MockReferralService) FindMyReferrals(patientID uint, status string) (*dto.ReferralMyListResponse, error) {
	args := m.Called(patientID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralMyListResponse), args.Error(1)
}

func (m *MockReferralService) FindByID(id uint) (*dto.ReferralResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralResponse), args.Error(1)
}

func (m *MockReferralService) FindByIDUnscoped(id uint) (*dto.ReferralResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralResponse), args.Error(1)
}

func (m *MockReferralService) FindByPatientID(patientID uint, query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	args := m.Called(patientID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralListResponse), args.Error(1)
}

func (m *MockReferralService) FindByDoctorID(doctorID uint, query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	args := m.Called(doctorID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralListResponse), args.Error(1)
}

func (m *MockReferralService) Create(req dto.CreateReferralRequest) (*dto.ReferralResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralResponse), args.Error(1)
}

func (m *MockReferralService) Update(id uint, req dto.UpdateReferralRequest) (*dto.ReferralResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralResponse), args.Error(1)
}

func (m *MockReferralService) Accept(id uint, req dto.AcceptReferralRequest) (*dto.ReferralResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralResponse), args.Error(1)
}

func (m *MockReferralService) Reject(id uint, req dto.RejectReferralRequest) (*dto.ReferralResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralResponse), args.Error(1)
}

func (m *MockReferralService) Complete(id uint, req dto.CompleteReferralRequest) (*dto.ReferralResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralResponse), args.Error(1)
}

func (m *MockReferralService) Cancel(id uint, req dto.CancelReferralRequest) (*dto.ReferralResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ReferralResponse), args.Error(1)
}

func (m *MockReferralService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockReferralService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockReferralService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
