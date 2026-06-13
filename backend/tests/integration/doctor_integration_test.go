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
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/doctor"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDoctorRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	doctorRepo := repository.NewDoctorRepository(db)
	doctorService := doctor.NewDoctorService(doctorRepo, cfg)
	doctorHandler := handler.NewDoctorHandler(doctorService)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupDoctorRouter(v1, routeCfg, doctorHandler)

	return r, routeCfg, db
}

func createPrerequisitesForDoctor(db *gorm.DB) (uint, uint) {
	dept := &models.Department{
		Name: "General",
		Code: "GEN",
	}
	db.Create(dept)

	spec := &models.DoctorSpecialization{
		Name: "General Practitioner",
		Code: "GP",
	}
	db.Create(spec)

	return dept.ID, spec.ID
}

func createDoctorAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, empID string, userID *uint) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	deptID, specID := createPrerequisitesForDoctor(db)

	reqBody := dto.CreateDoctorRequest{
		EmployeeID:       empID,
		FullName:         "Doctor " + empID,
		SpecializationID: &specID,
		LicenseNumber:    "LIC-" + empID,
		DepartmentID:     &deptID,
		UserID:           userID,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/doctors", bytes.NewBuffer(body))
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

func TestIntegration_Doctor_Create(t *testing.T) {
	r, cfg, db := setupDoctorRouter()
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	deptID, specID := createPrerequisitesForDoctor(db)

	reqBody := dto.CreateDoctorRequest{
		EmployeeID:       "DR-001",
		FullName:         "Dr. Strange",
		SpecializationID: &specID,
		LicenseNumber:    "LIC-001",
		DepartmentID:     &deptID,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/doctors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Doctor_FindByID(t *testing.T) {
	r, cfg, db := setupDoctorRouter()
	id := createDoctorAndGetID(r, cfg, db, "EMP-100", nil)
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/doctors/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Doctor_List(t *testing.T) {
	r, cfg, db := setupDoctorRouter()
	createDoctorAndGetID(r, cfg, db, "EMP-200", nil)
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/doctors", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Doctor_ActiveList(t *testing.T) {
	r, cfg, db := setupDoctorRouter()
	createDoctorAndGetID(r, cfg, db, "EMP-300", nil)
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/doctors/active", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Doctor_Update(t *testing.T) {
	r, cfg, db := setupDoctorRouter()
	id := createDoctorAndGetID(r, cfg, db, "EMP-400", nil)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := map[string]string{
		"full_name": "Updated Doctor Name",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/doctors/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Doctor_Me_GetAndUpdate(t *testing.T) {
	r, cfg, db := setupDoctorRouter()
	
	user := &models.User{
		Username: "doctor123",
		Email:    "doctor123@example.com",
		Role:     models.RoleDoctor,
	}
	db.Create(user)
	
	createDoctorAndGetID(r, cfg, db, "EMP-500", &user.ID)
	token := GenerateTestToken(user.ID, models.RoleDoctor, cfg.Config)

	// Get Me
	req1, _ := http.NewRequest(http.MethodGet, "/api/v1/doctors/me", nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Update Me
	reqBody := map[string]string{
		"phone": "08123456789",
	}
	body, _ := json.Marshal(reqBody)

	req2, _ := http.NewRequest(http.MethodPut, "/api/v1/doctors/me", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Doctor_DeactivateAndActivate(t *testing.T) {
	r, cfg, db := setupDoctorRouter()
	id := createDoctorAndGetID(r, cfg, db, "EMP-600", nil)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Deactivate
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/doctors/%d/deactivate", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Activate
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/doctors/%d/activate", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Doctor_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupDoctorRouter()
	id := createDoctorAndGetID(r, cfg, db, "EMP-700", nil)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/doctors/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/doctors/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Doctor_DeletedList(t *testing.T) {
	r, cfg, db := setupDoctorRouter()
	id := createDoctorAndGetID(r, cfg, db, "EMP-800", nil)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/doctors/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// List Deleted
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/doctors/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Doctor_HardDelete(t *testing.T) {
	r, cfg, db := setupDoctorRouter()
	id := createDoctorAndGetID(r, cfg, db, "EMP-900", nil)
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/doctors/%d/hard-delete", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
