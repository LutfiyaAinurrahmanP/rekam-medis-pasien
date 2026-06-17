package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/routes"
	patientservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/patient"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupPatientRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	patientRepo := repository.NewPatientRepository(db)
	patientService := patientservice.NewPatientService(patientRepo, cfg)
	patientHandler := handler.NewPatientHandler(patientService)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupPatientsRouter(v1, routeCfg, patientHandler)

	return r, routeCfg, db
}

func createPatientUserAndGetID(db *gorm.DB) uint {
	user := &models.User{
		Username: "patient123",
		Email:    "patient123@example.com",
		Role:     models.RolePatient,
	}
	db.Create(user)
	return user.ID
}

func createPatientAndGetID(r *gin.Engine, cfg *routes.RouteConfig, code string, userID *uint) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := dto.CreatePatientRequest{
		PatientCode: code,
		FullName:    "Test Patient " + code,
		DateOfBirth: "1990-01-01",
		Gender:      "male",
		UserID:      userID,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/patients", bytes.NewBuffer(body))
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

// ==============================================================================
// CREATE Tests
// ==============================================================================

func TestIntegration_Patient_Create_Success(t *testing.T) {
	r, cfg, _ := setupPatientRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := dto.CreatePatientRequest{
		PatientCode: "RM-001",
		FullName:    "John Doe",
		DateOfBirth: "1985-05-15",
		Gender:      "male",
		BloodType:   "O",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Patient_Create_ForbiddenForDoctor(t *testing.T) {
	r, cfg, _ := setupPatientRouter()
	doctorToken := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreatePatientRequest{
		PatientCode: "RM-002",
		FullName:    "Jane Doe",
		DateOfBirth: "1992-08-20",
		Gender:      "female",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/patients", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+doctorToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ==============================================================================
// GET Tests
// ==============================================================================

func TestIntegration_Patient_FindByID(t *testing.T) {
	r, cfg, _ := setupPatientRouter()
	patientID := createPatientAndGetID(r, cfg, "RM-100", nil)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/patients/%d", patientID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Patient_List(t *testing.T) {
	r, cfg, _ := setupPatientRouter()
	createPatientAndGetID(r, cfg, "RM-200", nil)
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config) // Doctor can list

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/patients", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==============================================================================
// UPDATE Tests
// ==============================================================================

func TestIntegration_Patient_Update(t *testing.T) {
	r, cfg, _ := setupPatientRouter()
	patientID := createPatientAndGetID(r, cfg, "RM-300", nil)
	token := GenerateTestToken(1, models.RoleReceptionist, cfg.Config)

	reqBody := map[string]string{
		"full_name":  "Updated Patient Name",
		"blood_type": "A",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/patients/%d", patientID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==============================================================================
// ME Tests (Self Service)
// ==============================================================================

func TestIntegration_Patient_GetMyPatientData(t *testing.T) {
	r, cfg, db := setupPatientRouter()
	userID := createPatientUserAndGetID(db)
	createPatientAndGetID(r, cfg, "RM-400", &userID)

	token := GenerateTestToken(userID, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/patients/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Patient_UpdateMyPatientData(t *testing.T) {
	r, cfg, db := setupPatientRouter()
	userID := createPatientUserAndGetID(db)
	createPatientAndGetID(r, cfg, "RM-500", &userID)

	token := GenerateTestToken(userID, models.RolePatient, cfg.Config)

	reqBody := map[string]string{
		"full_name": "My New Name",
		"phone":     "081234567890",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, "/api/v1/patients/me", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==============================================================================
// DELETE & RESTORE Tests
// ==============================================================================

func TestIntegration_Patient_SoftDelete(t *testing.T) {
	r, cfg, _ := setupPatientRouter()
	patientID := createPatientAndGetID(r, cfg, "RM-600", nil)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/patients/%d", patientID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Patient_Restore(t *testing.T) {
	r, cfg, _ := setupPatientRouter()
	patientID := createPatientAndGetID(r, cfg, "RM-700", nil)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Soft Delete First
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/patients/%d", patientID), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/patients/%d/restore", patientID), nil)
	req2.Header.Set("Authorization", "Bearer "+token)

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Patient_DeletedList(t *testing.T) {
	r, cfg, _ := setupPatientRouter()
	patientID := createPatientAndGetID(r, cfg, "RM-800", nil)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Soft Delete First
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/patients/%d", patientID), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// List Deleted
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/patients/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Patient_HardDelete(t *testing.T) {
	r, cfg, _ := setupPatientRouter()
	patientID := createPatientAndGetID(r, cfg, "RM-900", nil)
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/patients/%d/hard-delete", patientID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
