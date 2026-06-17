package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockReferralRepository struct {
	mock.Mock
}

func (m *MockReferralRepository) List(query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.ReferralPaginationMeta), args.Error(2)
	}
	return args.Get(0).([]models.Referral), args.Get(1).(dto.ReferralPaginationMeta), args.Error(2)
}

func (m *MockReferralRepository) DeletedList(query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.ReferralPaginationMeta), args.Error(2)
	}
	return args.Get(0).([]models.Referral), args.Get(1).(dto.ReferralPaginationMeta), args.Error(2)
}

func (m *MockReferralRepository) FindMyReferrals(patientID uint, status string) ([]models.Referral, error) {
	args := m.Called(patientID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Referral), args.Error(1)
}

func (m *MockReferralRepository) FindByID(id uint) (*models.Referral, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Referral), args.Error(1)
}

func (m *MockReferralRepository) FindByIDUnscoped(id uint) (*models.Referral, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Referral), args.Error(1)
}

func (m *MockReferralRepository) FindByPatientID(patientID uint, query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error) {
	args := m.Called(patientID, query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.ReferralPaginationMeta), args.Error(2)
	}
	return args.Get(0).([]models.Referral), args.Get(1).(dto.ReferralPaginationMeta), args.Error(2)
}

func (m *MockReferralRepository) FindByDoctorID(doctorID uint, query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error) {
	args := m.Called(doctorID, query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(dto.ReferralPaginationMeta), args.Error(2)
	}
	return args.Get(0).([]models.Referral), args.Get(1).(dto.ReferralPaginationMeta), args.Error(2)
}

func (m *MockReferralRepository) Create(referral *models.Referral) error {
	args := m.Called(referral)
	return args.Error(0)
}

func (m *MockReferralRepository) Update(referral *models.Referral) error {
	args := m.Called(referral)
	return args.Error(0)
}

func (m *MockReferralRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockReferralRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockReferralRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockReferralRepository) GenerateReferralNumber() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
