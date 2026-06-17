package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockDashboardRepository struct {
	mock.Mock
}

func (m *MockDashboardRepository) GetOverview(date string) (*dto.DashboardOverviewResponse, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardOverviewResponse), args.Error(1)
}

func (m *MockDashboardRepository) GetAdminDashboard(query *dto.DashboardPeriodQuery, startDate, endDate string) (*dto.DashboardAdminResponse, error) {
	args := m.Called(query, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardAdminResponse), args.Error(1)
}

func (m *MockDashboardRepository) GetDoctorDashboard(doctorID uint, date string) (*dto.DashboardDoctorResponse, error) {
	args := m.Called(doctorID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardDoctorResponse), args.Error(1)
}

func (m *MockDashboardRepository) GetReceptionistDashboard(date string) (*dto.DashboardReceptionistResponse, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardReceptionistResponse), args.Error(1)
}

func (m *MockDashboardRepository) GetPatientDashboard(patientID uint) (*dto.DashboardPatientResponse, error) {
	args := m.Called(patientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardPatientResponse), args.Error(1)
}

func (m *MockDashboardRepository) GetAppointmentReport(query *dto.DashboardAppointmentReportQuery, startDate, endDate string) (*dto.DashboardAppointmentReportResponse, error) {
	args := m.Called(query, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardAppointmentReportResponse), args.Error(1)
}

func (m *MockDashboardRepository) GetRevenueReport(query *dto.DashboardRevenueReportQuery, startDate, endDate string) (*dto.DashboardRevenueReportResponse, error) {
	args := m.Called(query, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardRevenueReportResponse), args.Error(1)
}

func (m *MockDashboardRepository) GetPatientReport(query *dto.DashboardPatientReportQuery, startDate, endDate string) (*dto.DashboardPatientReportResponse, error) {
	args := m.Called(query, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardPatientReportResponse), args.Error(1)
}
