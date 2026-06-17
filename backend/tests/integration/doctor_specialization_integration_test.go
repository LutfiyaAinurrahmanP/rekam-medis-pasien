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
	doctorspecialization "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/doctor-specialization"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDoctorSpecializationRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	dsRepo := repository.NewDoctorSpecializationRepository(db)
	dsService := doctorspecialization.NewDoctorSpecializationService(dsRepo, cfg)
	dsHandler := handler.NewDoctorSpecializationHandler(dsService)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupDoctorSpecializationRouter(v1, routeCfg, dsHandler)

	return r, routeCfg, db
}

func createSpecializationAndGetID(r *gin.Engine, cfg *routes.RouteConfig, name string) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	isActive := true
	reqBody := dto.CreateDoctorSpecializationRequest{
		Name:        name,
		Code:        name + "-CODE",
		Description: "Desc for " + name,
		IsActive:    &isActive,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/doctor-specializations", bytes.NewBuffer(body))
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

func TestIntegration_DocSpec_Create(t *testing.T) {
	r, cfg, _ := setupDoctorSpecializationRouter()
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	isActive := true
	reqBody := dto.CreateDoctorSpecializationRequest{
		Name:        "Pediatrics",
		Code:        "PED-01",
		Description: "Children specialist",
		IsActive:    &isActive,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/doctor-specializations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_DocSpec_FindByID(t *testing.T) {
	r, cfg, _ := setupDoctorSpecializationRouter()
	id := createSpecializationAndGetID(r, cfg, "Spec1")
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/doctor-specializations/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_DocSpec_List(t *testing.T) {
	r, cfg, _ := setupDoctorSpecializationRouter()
	createSpecializationAndGetID(r, cfg, "Spec2")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/doctor-specializations", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_DocSpec_ActiveList(t *testing.T) {
	r, cfg, _ := setupDoctorSpecializationRouter()
	createSpecializationAndGetID(r, cfg, "Spec3")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/doctor-specializations/active", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_DocSpec_Update(t *testing.T) {
	r, cfg, _ := setupDoctorSpecializationRouter()
	id := createSpecializationAndGetID(r, cfg, "Spec4")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := map[string]string{
		"name": "Updated Spec",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/doctor-specializations/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_DocSpec_DeactivateAndActivate(t *testing.T) {
	r, cfg, _ := setupDoctorSpecializationRouter()
	id := createSpecializationAndGetID(r, cfg, "Spec5")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Deactivate
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/doctor-specializations/%d/deactivate", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Activate
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/doctor-specializations/%d/activate", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_DocSpec_DeleteAndRestore(t *testing.T) {
	r, cfg, _ := setupDoctorSpecializationRouter()
	id := createSpecializationAndGetID(r, cfg, "Spec6")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/doctor-specializations/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/doctor-specializations/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_DocSpec_DeletedList(t *testing.T) {
	r, cfg, _ := setupDoctorSpecializationRouter()
	id := createSpecializationAndGetID(r, cfg, "Spec7")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/doctor-specializations/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// List
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/doctor-specializations/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_DocSpec_HardDelete(t *testing.T) {
	r, cfg, _ := setupDoctorSpecializationRouter()
	id := createSpecializationAndGetID(r, cfg, "Spec8")
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/doctor-specializations/%d/hard-delete", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
