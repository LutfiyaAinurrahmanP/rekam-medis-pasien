package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDashboardHandlerTest() (*gin.Engine, *mocks.MockDashboardService, *mocks.MockDoctorRepository, *mocks.MockPatientRepository) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDashboardService)
	mockDoctorRepo := new(mocks.MockDoctorRepository)
	mockPatientRepo := new(mocks.MockPatientRepository)

	h := handler.NewDashboardHandler(mockService, mockDoctorRepo, mockPatientRepo)

	r := gin.Default()

	r.GET("/dashboards/overview", h.Overview)
	r.GET("/dashboards/admin", h.AdminDashboard)

	// Context mocking middleware
	authGroup := r.Group("/")
	authGroup.Use(func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("role", "doctor")
		c.Next()
	})
	authGroup.GET("/dashboards/doctor", h.DoctorDashboard)
	authGroup.GET("/dashboards/patient", h.PatientDashboard)
	authGroup.GET("/dashboards/receptionist", h.ReceptionistDashboard)
	authGroup.GET("/reports/appointments", h.AppointmentReport)
	authGroup.GET("/reports/revenue", h.RevenueReport)
	authGroup.GET("/reports/patients", h.PatientReport)

	return r, mockService, mockDoctorRepo, mockPatientRepo
}

func TestDashboardHandler_Overview_Success(t *testing.T) {
	r, mockService, _, _ := setupDashboardHandlerTest()

	expectedRes := mocks.NewTestDashboardOverviewResponse()
	mockService.On("GetOverview", mock.AnythingOfType("*dto.DashboardOverviewQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/dashboards/overview?date=2024-01-01", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_AdminDashboard_Success(t *testing.T) {
	r, mockService, _, _ := setupDashboardHandlerTest()

	expectedRes := mocks.NewTestDashboardAdminResponse()
	mockService.On("GetAdminDashboard", mock.AnythingOfType("*dto.DashboardPeriodQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/dashboards/admin?period=today", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_DoctorDashboard_Success(t *testing.T) {
	r, mockService, mockDoctorRepo, _ := setupDashboardHandlerTest()

	expectedRes := mocks.NewTestDashboardDoctorResponse()
	mockDoctor := &models.Doctor{ID: 1}

	mockDoctorRepo.On("FindByUserID", uint(1)).Return(mockDoctor, nil)
	mockService.On("GetDoctorDashboard", uint(1), mock.AnythingOfType("*dto.DashboardDoctorQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/dashboards/doctor?date=2024-01-01", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_ReceptionistDashboard_Success(t *testing.T) {
	r, mockService, _, _ := setupDashboardHandlerTest()

	expectedRes := mocks.NewTestDashboardReceptionistResponse()
	mockService.On("GetReceptionistDashboard").Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/dashboards/receptionist", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_PatientDashboard_Success(t *testing.T) {
	r, mockService, _, mockPatientRepo := setupDashboardHandlerTest()

	expectedRes := mocks.NewTestDashboardPatientResponse()
	mockPatient := &models.Patient{ID: 1}

	mockPatientRepo.On("FindByUserID", uint(1)).Return(mockPatient, nil)
	mockService.On("GetPatientDashboard", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/dashboards/patient", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_AppointmentReport_Success(t *testing.T) {
	r, mockService, mockDoctorRepo, _ := setupDashboardHandlerTest()

	expectedRes := mocks.NewTestDashboardAppointmentReportResponse()
	mockDoctor := &models.Doctor{ID: 1}

	mockDoctorRepo.On("FindByUserID", uint(1)).Return(mockDoctor, nil)
	mockService.On("GetAppointmentReport", mock.AnythingOfType("*dto.DashboardAppointmentReportQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/reports/appointments?period=today", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_RevenueReport_Success(t *testing.T) {
	r, mockService, _, _ := setupDashboardHandlerTest()

	expectedRes := mocks.NewTestDashboardRevenueReportResponse()
	mockService.On("GetRevenueReport", mock.AnythingOfType("*dto.DashboardRevenueReportQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/reports/revenue?period=today", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDashboardHandler_PatientReport_Success(t *testing.T) {
	r, mockService, _, _ := setupDashboardHandlerTest()

	expectedRes := mocks.NewTestDashboardPatientReportResponse()
	mockService.On("GetPatientReport", mock.AnythingOfType("*dto.DashboardPatientReportQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/reports/patients?period=today", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
