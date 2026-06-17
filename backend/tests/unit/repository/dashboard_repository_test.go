package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDashboardRepository_GetOverview_Success(t *testing.T) {
	mockRepo := new(mocks.MockDashboardRepository)
	expectedRes := mocks.NewTestDashboardOverviewResponse()

	mockRepo.On("GetOverview", "2024-01-01").Return(expectedRes, nil)

	res, err := mockRepo.GetOverview("2024-01-01")
	assert.NoError(t, err)
	assert.Equal(t, expectedRes.SummaryDate, res.SummaryDate)
	mockRepo.AssertExpectations(t)
}

func TestDashboardRepository_GetAdminDashboard_Success(t *testing.T) {
	mockRepo := new(mocks.MockDashboardRepository)
	query := &dto.DashboardPeriodQuery{Period: "today"}
	expectedRes := mocks.NewTestDashboardAdminResponse()

	mockRepo.On("GetAdminDashboard", query, "2024-01-01", "2024-01-01").Return(expectedRes, nil)

	res, err := mockRepo.GetAdminDashboard(query, "2024-01-01", "2024-01-01")
	assert.NoError(t, err)
	assert.Equal(t, expectedRes.Period, res.Period)
	mockRepo.AssertExpectations(t)
}

func TestDashboardRepository_GetDoctorDashboard_Success(t *testing.T) {
	mockRepo := new(mocks.MockDashboardRepository)
	expectedRes := mocks.NewTestDashboardDoctorResponse()

	mockRepo.On("GetDoctorDashboard", uint(1), "2024-01-01").Return(expectedRes, nil)

	res, err := mockRepo.GetDoctorDashboard(1, "2024-01-01")
	assert.NoError(t, err)
	assert.Equal(t, expectedRes.Doctor.ID, res.Doctor.ID)
	mockRepo.AssertExpectations(t)
}

func TestDashboardRepository_GetReceptionistDashboard_Success(t *testing.T) {
	mockRepo := new(mocks.MockDashboardRepository)
	expectedRes := mocks.NewTestDashboardReceptionistResponse()

	mockRepo.On("GetReceptionistDashboard", "2024-01-01").Return(expectedRes, nil)

	res, err := mockRepo.GetReceptionistDashboard("2024-01-01")
	assert.NoError(t, err)
	assert.Equal(t, expectedRes.Date, res.Date)
	mockRepo.AssertExpectations(t)
}

func TestDashboardRepository_GetPatientDashboard_Success(t *testing.T) {
	mockRepo := new(mocks.MockDashboardRepository)
	expectedRes := mocks.NewTestDashboardPatientResponse()

	mockRepo.On("GetPatientDashboard", uint(1)).Return(expectedRes, nil)

	res, err := mockRepo.GetPatientDashboard(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRes.Patient.ID, res.Patient.ID)
	mockRepo.AssertExpectations(t)
}

func TestDashboardRepository_GetAppointmentReport_Success(t *testing.T) {
	mockRepo := new(mocks.MockDashboardRepository)
	query := &dto.DashboardAppointmentReportQuery{}
	expectedRes := mocks.NewTestDashboardAppointmentReportResponse()

	mockRepo.On("GetAppointmentReport", query, "2024-01-01", "2024-01-31").Return(expectedRes, nil)

	res, err := mockRepo.GetAppointmentReport(query, "2024-01-01", "2024-01-31")
	assert.NoError(t, err)
	assert.Equal(t, expectedRes.Totals.Total, res.Totals.Total)
	mockRepo.AssertExpectations(t)
}

func TestDashboardRepository_GetRevenueReport_Success(t *testing.T) {
	mockRepo := new(mocks.MockDashboardRepository)
	query := &dto.DashboardRevenueReportQuery{}
	expectedRes := mocks.NewTestDashboardRevenueReportResponse()

	mockRepo.On("GetRevenueReport", query, "2024-01-01", "2024-01-31").Return(expectedRes, nil)

	res, err := mockRepo.GetRevenueReport(query, "2024-01-01", "2024-01-31")
	assert.NoError(t, err)
	assert.Equal(t, expectedRes.Revenue.TotalBilled, res.Revenue.TotalBilled)
	mockRepo.AssertExpectations(t)
}

func TestDashboardRepository_GetPatientReport_Success(t *testing.T) {
	mockRepo := new(mocks.MockDashboardRepository)
	query := &dto.DashboardPatientReportQuery{}
	expectedRes := mocks.NewTestDashboardPatientReportResponse()

	mockRepo.On("GetPatientReport", query, "2024-01-01", "2024-01-31").Return(expectedRes, nil)

	res, err := mockRepo.GetPatientReport(query, "2024-01-01", "2024-01-31")
	assert.NoError(t, err)
	assert.Equal(t, expectedRes.Registrations.NewPatients, res.Registrations.NewPatients)
	mockRepo.AssertExpectations(t)
}
