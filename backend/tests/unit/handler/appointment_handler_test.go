package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAppointmentHandlerTest() (*gin.Engine, *mocks.MockAppointmentService, *mocks.MockDoctorRepository, *mocks.MockPatientRepository) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockAppointmentService)
	mockDoctorRepo := new(mocks.MockDoctorRepository)
	mockPatientRepo := new(mocks.MockPatientRepository)

	h := handler.NewAppointmentHandler(mockService, mockDoctorRepo, mockPatientRepo)

	r := gin.Default()

	// Inject role patient manually for testing ownership
	r.Use(func(c *gin.Context) {
		c.Set("role", "patient")
		c.Set("user_id", uint(1))
		c.Next()
	})

	r.GET("/appointments/me", h.MyAppointments)
	r.GET("/appointments/my-schedule", h.MySchedule)
	r.GET("/appointments", h.List)
	r.GET("/appointments/today", h.TodayList)
	r.GET("/appointments/upcoming", h.UpcomingList)
	r.GET("/appointments/past", h.PastList)
	r.GET("/appointments/cancelled", h.CancelledList)
	r.GET("/appointments/deleted", h.DeletedList)
	r.GET("/appointments/:id", h.FindByID)
	r.POST("/appointments", h.Create)
	r.PUT("/appointments/:id", h.Update)
	r.PATCH("/appointments/:id/confirm", h.Confirm)
	r.PATCH("/appointments/:id/start", h.Start)
	r.PATCH("/appointments/:id/complete", h.Complete)
	r.PATCH("/appointments/:id/cancel", h.Cancel)
	r.PATCH("/appointments/:id/reschedule", h.Reschedule)
	r.PATCH("/appointments/:id/no-show", h.NoShow)
	r.DELETE("/appointments/:id", h.SoftDelete)
	r.PATCH("/appointments/:id/restore", h.Restore)
	r.DELETE("/appointments/:id/hard", h.HardDelete)

	return r, mockService, mockDoctorRepo, mockPatientRepo
}

func TestAppointmentHandler_List_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()

	expectedRes := &dto.AppointmentListResponse{
		Data: []dto.AppointmentResponse{
			*mocks.NewTestAppointmentResponse(mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false)),
		},
		Meta: dto.AppointmentPaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.AppointmentPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/appointments?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_FindByID_Success(t *testing.T) {
	r, mockService, _, mockPatientRepo := setupAppointmentHandlerTest()

	expectedRes := mocks.NewTestAppointmentResponse(mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false))

	// Mock for ownership check
	uid := uint(1)
	mockPatientRepo.On("FindByUserID", uint(1)).Return(&models.Patient{ID: 1, UserID: &uid}, nil)
	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/appointments/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
	mockPatientRepo.AssertExpectations(t)
}

func TestAppointmentHandler_Create_Success(t *testing.T) {
	r, mockService, _, mockPatientRepo := setupAppointmentHandlerTest()
	createReq := mocks.NewCreateAppointmentRequest(1, 1, "2023-12-01", "10:00", 30, "Reason", "Notes")
	expectedRes := mocks.NewTestAppointmentResponse(mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false))

	body, _ := json.Marshal(createReq)

	uid := uint(1)
	mockPatientRepo.On("FindByUserID", uint(1)).Return(&models.Patient{ID: 1, UserID: &uid}, nil)
	mockService.On("Create", mock.AnythingOfType("*dto.CreateAppointmentRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/appointments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_Update_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()
	updateReq := mocks.NewUpdateAppointmentRequest("2023-12-02", "11:00", 45, "Reason2", "Notes2")
	expectedRes := mocks.NewTestAppointmentResponse(mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-02", "11:00", "scheduled", "Reason2", "Notes2", 45, false))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateAppointmentRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/appointments/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_Confirm_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()

	mockService.On("Confirm", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/appointments/1/confirm", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_Start_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()

	mockService.On("Start", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/appointments/1/start", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_Complete_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()

	mockService.On("Complete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/appointments/1/complete", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_Cancel_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()
	cancelReq := &dto.CancelAppointmentRequest{Reason: "Patient cancelled"}

	body, _ := json.Marshal(cancelReq)

	mockService.On("Cancel", uint(1), mock.AnythingOfType("*dto.CancelAppointmentRequest")).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/appointments/1/cancel", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_Reschedule_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()
	rescheduleReq := &dto.RescheduleAppointmentRequest{AppointmentDate: "2023-12-05", AppointmentTime: "14:00", Reason: "Doctor unavailable"}

	body, _ := json.Marshal(rescheduleReq)

	mockService.On("Reschedule", uint(1), mock.AnythingOfType("*dto.RescheduleAppointmentRequest")).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/appointments/1/reschedule", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_NoShow_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()

	mockService.On("NoShow", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/appointments/1/no-show", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_SoftDelete_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()

	mockService.On("SoftDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/appointments/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_Restore_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()

	mockService.On("Restore", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/appointments/1/restore", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAppointmentHandler_HardDelete_Success(t *testing.T) {
	r, mockService, _, _ := setupAppointmentHandlerTest()

	mockService.On("HardDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/appointments/1/hard", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
