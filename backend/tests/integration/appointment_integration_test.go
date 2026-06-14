package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/routes"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/appointment"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupAppointmentRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewAppointmentRepository(db)
	service := appointment.NewAppointmentService(repo, cfg)
	docRepo := repository.NewDoctorRepository(db)
	patRepo := repository.NewPatientRepository(db)
	handlerInst := handler.NewAppointmentHandler(service, docRepo, patRepo)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupAppointmentRouter(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createPrereqForAppointment(db *gorm.DB) (uint, uint) {
	patient := &models.Patient{
		FullName:    "Test Patient",
		PatientCode: "PT-001",
		DateOfBirth: "1990-01-01",
		Gender:      "male",
		BloodType:   "O",
	}
	db.Create(patient)

	department := &models.Department{Name: "General"}
	db.Create(department)

	doctor := &models.Doctor{
		FullName:     "Test Doctor",
		DepartmentID: &department.ID,
	}
	db.Create(doctor)

	return patient.ID, doctor.ID
}

func createAppointmentAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, pID uint, dID uint) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := dto.CreateAppointmentRequest{
		PatientID:       pID,
		DoctorID:        dID,
		AppointmentDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		AppointmentTime: "10:00",
		DurationMinutes: 30,
		Reason:          "Checkup",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/appointments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	if res["data"] == nil {
		return 0
	}
	return uint(res["data"].(map[string]interface{})["id"].(float64))
}

func TestIntegration_Appointment_Create(t *testing.T) {
	r, cfg, db := setupAppointmentRouter()
	pID, dID := createPrereqForAppointment(db)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := dto.CreateAppointmentRequest{
		PatientID:       pID,
		DoctorID:        dID,
		AppointmentDate: time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
		AppointmentTime: "10:00",
		DurationMinutes: 30,
		Reason:          "Checkup",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/appointments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Appointment_FindByID(t *testing.T) {
	r, cfg, db := setupAppointmentRouter()
	pID, dID := createPrereqForAppointment(db)
	id := createAppointmentAndGetID(r, cfg, db, pID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/appointments/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Appointment_Update(t *testing.T) {
	r, cfg, db := setupAppointmentRouter()
	pID, dID := createPrereqForAppointment(db)
	id := createAppointmentAndGetID(r, cfg, db, pID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	reason := "Updated Reason"
	reqBody := dto.UpdateAppointmentRequest{
		Reason: &reason,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/appointments/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Appointment_Lists(t *testing.T) {
	r, cfg, db := setupAppointmentRouter()
	pID, dID := createPrereqForAppointment(db)
	createAppointmentAndGetID(r, cfg, db, pID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	endpoints := []string{
		"/api/v1/appointments",
		"/api/v1/appointments/upcoming",
		"/api/v1/appointments/past",
		"/api/v1/appointments/today",
	}

	for _, ep := range endpoints {
		req, _ := http.NewRequest(http.MethodGet, ep, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, ep)
	}
}

func TestIntegration_Appointment_RescheduleAndCancel(t *testing.T) {
	r, cfg, db := setupAppointmentRouter()
	pID, dID := createPrereqForAppointment(db)
	id := createAppointmentAndGetID(r, cfg, db, pID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Reschedule
	reqBodyR := dto.RescheduleAppointmentRequest{
		AppointmentDate: time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
		AppointmentTime: "11:00",
		Reason:          "Patient request",
	}
	bodyR, _ := json.Marshal(reqBodyR)
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/appointments/%d/reschedule", id), bytes.NewBuffer(bodyR))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Cancel
	reqBodyC := dto.CancelAppointmentRequest{Reason: "Cannot make it"}
	bodyC, _ := json.Marshal(reqBodyC)
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/appointments/%d/cancel", id), bytes.NewBuffer(bodyC))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Appointment_LifecycleStates(t *testing.T) {
	r, cfg, db := setupAppointmentRouter()
	pID, dID := createPrereqForAppointment(db)
	id := createAppointmentAndGetID(r, cfg, db, pID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Confirm
	reqC, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/appointments/%d/confirm", id), nil)
	reqC.Header.Set("Authorization", "Bearer "+token)
	wC := httptest.NewRecorder()
	r.ServeHTTP(wC, reqC)
	assert.Equal(t, http.StatusOK, wC.Code)

	// Delete and recreate another one for NoShow
	// Since state machine probably doesn't allow confirm -> no-show easily or it might. Let's just create a new one.
	id2 := createAppointmentAndGetID(r, cfg, db, pID, dID)
	reqN, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/appointments/%d/no-show", id2), nil)
	reqN.Header.Set("Authorization", "Bearer "+token)
	wN := httptest.NewRecorder()
	r.ServeHTTP(wN, reqN)
	assert.Equal(t, http.StatusOK, wN.Code)
}

func TestIntegration_Appointment_DoctorActions(t *testing.T) {
	r, cfg, db := setupAppointmentRouter()
	pID, dID := createPrereqForAppointment(db)
	id := createAppointmentAndGetID(r, cfg, db, pID, dID)

	// Set appointment status to confirmed via db directly or via API so we can start it
	db.Model(&models.Appointment{}).Where("id = ?", id).Update("status", "confirmed")

	doctorToken := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)

	// Start
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/appointments/%d/start", id), nil)
	req1.Header.Set("Authorization", "Bearer "+doctorToken)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Complete
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/appointments/%d/complete", id), nil)
	req2.Header.Set("Authorization", "Bearer "+doctorToken)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Appointment_PersonalLists(t *testing.T) {
	r, cfg, db := setupAppointmentRouter()
	pID, dID := createPrereqForAppointment(db)
	createAppointmentAndGetID(r, cfg, db, pID, dID)

	// Patient List
	pToken := GenerateTestToken(pID, models.RolePatient, cfg.Config)
	req1, _ := http.NewRequest(http.MethodGet, "/api/v1/appointments/my-appointments", nil)
	req1.Header.Set("Authorization", "Bearer "+pToken)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Doctor List
	dToken := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/appointments/my-schedule", nil)
	req2.Header.Set("Authorization", "Bearer "+dToken)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Appointment_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupAppointmentRouter()
	pID, dID := createPrereqForAppointment(db)
	id := createAppointmentAndGetID(r, cfg, db, pID, dID)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/appointments/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/appointments/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Appointment_DeletedListAndHardDelete(t *testing.T) {
	r, cfg, db := setupAppointmentRouter()
	pID, dID := createPrereqForAppointment(db)
	id := createAppointmentAndGetID(r, cfg, db, pID, dID)
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/appointments/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Deleted List
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/appointments/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Hard Delete
	req3, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/appointments/%d/hard-delete", id), nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}
