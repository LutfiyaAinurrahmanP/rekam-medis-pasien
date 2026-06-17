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
	type_test "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/type-test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTypeTestRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewTypeTestRepository(db)
	service := type_test.NewTypeTestService(repo, cfg)
	handlerInst := handler.NewTypeTestHandler(service)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupTypeTestRouter(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createPrerequisiteTTC(db *gorm.DB) uint {
	ttc := &models.TypeTestCategory{
		Name: "Lab Tests",
		Code: "LAB",
	}
	db.Create(ttc)
	return ttc.ID
}

func createTypeTestAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, code string) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	ttcID := createPrerequisiteTTC(db)

	p := 150000.0
	isActive := true
	reqBody := dto.CreateTypeTestRequest{
		Name:               "Test " + code,
		Code:               code,
		TypeTestCategoryID: ttcID,
		Description:        "Desc for " + code,
		Price:              &p,
		IsActive:           &isActive,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/type-tests", bytes.NewBuffer(body))
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

func TestIntegration_TypeTest_Create(t *testing.T) {
	r, cfg, db := setupTypeTestRouter()
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	ttcID := createPrerequisiteTTC(db)

	p := 200000.0
	reqBody := dto.CreateTypeTestRequest{
		Name:               "Blood Sugar Test",
		Code:               "BS-01",
		TypeTestCategoryID: ttcID,
		Price:              &p,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/type-tests", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_TypeTest_FindByID(t *testing.T) {
	r, cfg, db := setupTypeTestRouter()
	id := createTypeTestAndGetID(r, cfg, db, "TT-100")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/type-tests/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TypeTest_List(t *testing.T) {
	r, cfg, db := setupTypeTestRouter()
	createTypeTestAndGetID(r, cfg, db, "TT-200")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/type-tests", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TypeTest_ActiveList(t *testing.T) {
	r, cfg, db := setupTypeTestRouter()
	createTypeTestAndGetID(r, cfg, db, "TT-300")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/type-tests/active", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TypeTest_Update(t *testing.T) {
	r, cfg, db := setupTypeTestRouter()
	id := createTypeTestAndGetID(r, cfg, db, "TT-400")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := map[string]string{
		"name": "Updated Test Name",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/type-tests/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TypeTest_DeactivateAndActivate(t *testing.T) {
	r, cfg, db := setupTypeTestRouter()
	id := createTypeTestAndGetID(r, cfg, db, "TT-500")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Deactivate
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/type-tests/%d/deactivate", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Activate
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/type-tests/%d/activate", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_TypeTest_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupTypeTestRouter()
	id := createTypeTestAndGetID(r, cfg, db, "TT-600")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/type-tests/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/type-tests/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_TypeTest_DeletedList(t *testing.T) {
	r, cfg, db := setupTypeTestRouter()
	id := createTypeTestAndGetID(r, cfg, db, "TT-700")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/type-tests/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// List
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/type-tests/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_TypeTest_HardDelete(t *testing.T) {
	r, cfg, db := setupTypeTestRouter()
	id := createTypeTestAndGetID(r, cfg, db, "TT-800")
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/type-tests/%d/hard-delete", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
