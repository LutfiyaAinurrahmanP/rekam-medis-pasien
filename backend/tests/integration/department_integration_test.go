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
	departmentservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/department"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDepartmentRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	deptRepo := repository.NewDepartmentRepository(db)
	deptService := departmentservice.NewDepartmentService(deptRepo, cfg)
	deptHandler := handler.NewDepartmentHandler(deptService)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupDepartmentsRouter(v1, routeCfg, deptHandler)

	return r, routeCfg, db
}

// Helper to create a department and return its ID
func createDepartmentAndGetID(r *gin.Engine, cfg *routes.RouteConfig, code string) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := dto.CreateDepartmentRequest{
		Name:          "Test Department " + code,
		Code:          code,
		Description:   "Description for " + code,
		FloorLocation: "1st Floor",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/departments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	
	// Check if data is null, wait what is the structure?
	if res["data"] == nil {
		return 0
	}
	return uint(res["data"].(map[string]interface{})["id"].(float64))
}

// ==============================================================================
// CREATE Tests
// ==============================================================================

func TestIntegration_Department_Create_Success(t *testing.T) {
	r, cfg, _ := setupDepartmentRouter()
	adminToken := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := dto.CreateDepartmentRequest{
		Name:          "Cardiology",
		Code:          "CARD-01",
		Description:   "Heart specialists",
		FloorLocation: "2nd Floor",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/departments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Department_Create_ForbiddenForDoctor(t *testing.T) {
	r, cfg, _ := setupDepartmentRouter()
	doctorToken := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateDepartmentRequest{
		Name:          "Neurology",
		Code:          "NEUR-01",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/departments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+doctorToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// ==============================================================================
// GET Tests
// ==============================================================================

func TestIntegration_Department_FindByID(t *testing.T) {
	r, cfg, _ := setupDepartmentRouter()
	deptID := createDepartmentAndGetID(r, cfg, "DEP-100")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/departments/%d", deptID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Department_List(t *testing.T) {
	r, cfg, _ := setupDepartmentRouter()
	createDepartmentAndGetID(r, cfg, "DEP-200")
	token := GenerateTestToken(1, models.RoleReceptionist, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/departments", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==============================================================================
// UPDATE Tests
// ==============================================================================

func TestIntegration_Department_Update(t *testing.T) {
	r, cfg, _ := setupDepartmentRouter()
	deptID := createDepartmentAndGetID(r, cfg, "DEP-300")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := map[string]string{
		"name":           "Updated Department",
		"code":           "DEP-300",
		"description":    "Updated description",
		"floor_location": "3rd Floor",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/departments/%d", deptID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==============================================================================
// DELETE & RESTORE Tests
// ==============================================================================

func TestIntegration_Department_SoftDelete(t *testing.T) {
	r, cfg, _ := setupDepartmentRouter()
	deptID := createDepartmentAndGetID(r, cfg, "DEP-400")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/departments/%d", deptID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Department_Restore(t *testing.T) {
	r, cfg, _ := setupDepartmentRouter()
	deptID := createDepartmentAndGetID(r, cfg, "DEP-500")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete first
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/departments/%d", deptID), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/departments/%d/restore", deptID), nil)
	req2.Header.Set("Authorization", "Bearer "+token)

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Department_DeletedList(t *testing.T) {
	r, cfg, _ := setupDepartmentRouter()
	deptID := createDepartmentAndGetID(r, cfg, "DEP-600")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete first
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/departments/%d", deptID), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// List deleted
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/departments/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Department_HardDelete(t *testing.T) {
	r, cfg, _ := setupDepartmentRouter()
	deptID := createDepartmentAndGetID(r, cfg, "DEP-700")
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/departments/%d/hard-delete", deptID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
