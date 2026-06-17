package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockDashboardService struct {
	mock.Mock
}

func (m *MockDashboardService) GetOverview(query *dto.DashboardOverviewQuery) (*dto.DashboardOverviewResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardOverviewResponse), args.Error(1)
}

func (m *MockDashboardService) GetAdminDashboard(query *dto.DashboardPeriodQuery) (*dto.DashboardAdminResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardAdminResponse), args.Error(1)
}

func (m *MockDashboardService) GetDoctorDashboard(doctorID uint, query *dto.DashboardDoctorQuery) (*dto.DashboardDoctorResponse, error) {
	args := m.Called(doctorID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardDoctorResponse), args.Error(1)
}

func (m *MockDashboardService) GetReceptionistDashboard() (*dto.DashboardReceptionistResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardReceptionistResponse), args.Error(1)
}

func (m *MockDashboardService) GetPatientDashboard(patientID uint) (*dto.DashboardPatientResponse, error) {
	args := m.Called(patientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardPatientResponse), args.Error(1)
}

func (m *MockDashboardService) GetAppointmentReport(query *dto.DashboardAppointmentReportQuery) (*dto.DashboardAppointmentReportResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardAppointmentReportResponse), args.Error(1)
}

func (m *MockDashboardService) GetRevenueReport(query *dto.DashboardRevenueReportQuery) (*dto.DashboardRevenueReportResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardRevenueReportResponse), args.Error(1)
}

func (m *MockDashboardService) GetPatientReport(query *dto.DashboardPatientReportQuery) (*dto.DashboardPatientReportResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DashboardPatientReportResponse), args.Error(1)
}
