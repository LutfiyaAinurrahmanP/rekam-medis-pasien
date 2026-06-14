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
	type_test_category "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/type-test-category"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTypeTestCategoryRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewTypeTestCategoryRepository(db)
	service := type_test_category.NewTypeTestCategoryService(repo, cfg)
	handlerInst := handler.NewTypeTestCategoryHandler(service)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupTypeTestCategoryRouter(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createTypeTestCategoryAndGetID(r *gin.Engine, cfg *routes.RouteConfig, name string) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	isActive := true
	reqBody := dto.CreateTypeTestCategoryRequest{
		Name:        name,
		Code:        name + "-C",
		Description: "Category " + name,
		IsActive:    &isActive,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/type-test-categories", bytes.NewBuffer(body))
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

func TestIntegration_TTC_Create(t *testing.T) {
	r, cfg, _ := setupTypeTestCategoryRouter()
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	isActive := true
	reqBody := dto.CreateTypeTestCategoryRequest{
		Name:        "Hematology",
		Code:        "HEM",
		Description: "Blood tests",
		IsActive:    &isActive,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/type-test-categories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_TTC_FindByID(t *testing.T) {
	r, cfg, _ := setupTypeTestCategoryRouter()
	id := createTypeTestCategoryAndGetID(r, cfg, "Cat1")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/type-test-categories/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TTC_List(t *testing.T) {
	r, cfg, _ := setupTypeTestCategoryRouter()
	createTypeTestCategoryAndGetID(r, cfg, "Cat2")
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/type-test-categories", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TTC_ActiveList(t *testing.T) {
	r, cfg, _ := setupTypeTestCategoryRouter()
	createTypeTestCategoryAndGetID(r, cfg, "Cat3")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/type-test-categories/active", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TTC_Update(t *testing.T) {
	r, cfg, _ := setupTypeTestCategoryRouter()
	id := createTypeTestCategoryAndGetID(r, cfg, "Cat4")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := map[string]string{
		"name": "Updated Cat",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/type-test-categories/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_TTC_DeactivateAndActivate(t *testing.T) {
	r, cfg, _ := setupTypeTestCategoryRouter()
	id := createTypeTestCategoryAndGetID(r, cfg, "Cat5")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Deactivate
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/type-test-categories/%d/deactivate", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Activate
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/type-test-categories/%d/activate", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_TTC_DeleteAndRestore(t *testing.T) {
	r, cfg, _ := setupTypeTestCategoryRouter()
	id := createTypeTestCategoryAndGetID(r, cfg, "Cat6")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/type-test-categories/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/type-test-categories/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_TTC_DeletedList(t *testing.T) {
	r, cfg, _ := setupTypeTestCategoryRouter()
	id := createTypeTestCategoryAndGetID(r, cfg, "Cat7")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/type-test-categories/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// List Deleted
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/type-test-categories/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_TTC_HardDelete(t *testing.T) {
	r, cfg, _ := setupTypeTestCategoryRouter()
	id := createTypeTestCategoryAndGetID(r, cfg, "Cat8")
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/type-test-categories/%d/hard-delete", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
