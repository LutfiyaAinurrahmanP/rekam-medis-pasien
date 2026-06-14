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
	medicine_type "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medicine-type"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupMedicineTypeRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewMedicineTypeRepository(db)
	service := medicine_type.NewMedicineTypeService(repo, cfg)
	handlerInst := handler.NewMedicineTypeHandler(service)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupMedicineTypeRouter(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createMedicineTypeAndGetID(r *gin.Engine, cfg *routes.RouteConfig, name string) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	isActive := true
	reqBody := dto.CreateMedicineTypeRequest{
		Name:        name,
		Code:        name + "-CODE",
		Description: "Desc for " + name,
		IsActive:    &isActive,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medicine-types", bytes.NewBuffer(body))
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

func TestIntegration_MedicineType_Create(t *testing.T) {
	r, cfg, _ := setupMedicineTypeRouter()
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	isActive := true
	reqBody := dto.CreateMedicineTypeRequest{
		Name:        "Tablet",
		Code:        "TAB-01",
		Description: "Tablet medicine",
		IsActive:    &isActive,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medicine-types", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_MedicineType_FindByID(t *testing.T) {
	r, cfg, _ := setupMedicineTypeRouter()
	id := createMedicineTypeAndGetID(r, cfg, "MT1")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/medicine-types/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_MedicineType_List(t *testing.T) {
	r, cfg, _ := setupMedicineTypeRouter()
	createMedicineTypeAndGetID(r, cfg, "MT2")
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/medicine-types", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_MedicineType_ActiveList(t *testing.T) {
	r, cfg, _ := setupMedicineTypeRouter()
	createMedicineTypeAndGetID(r, cfg, "MT3")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/medicine-types/active", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_MedicineType_Update(t *testing.T) {
	r, cfg, _ := setupMedicineTypeRouter()
	id := createMedicineTypeAndGetID(r, cfg, "MT4")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	newName := "Updated Type"
	reqBody := dto.UpdateMedicineTypeRequest{
		Name: &newName,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/medicine-types/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_MedicineType_DeactivateAndActivate(t *testing.T) {
	r, cfg, _ := setupMedicineTypeRouter()
	id := createMedicineTypeAndGetID(r, cfg, "MT5")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Deactivate
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medicine-types/%d/deactivate", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Activate
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medicine-types/%d/activate", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_MedicineType_DeleteAndRestore(t *testing.T) {
	r, cfg, _ := setupMedicineTypeRouter()
	id := createMedicineTypeAndGetID(r, cfg, "MT6")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medicine-types/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medicine-types/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_MedicineType_DeletedList(t *testing.T) {
	r, cfg, _ := setupMedicineTypeRouter()
	id := createMedicineTypeAndGetID(r, cfg, "MT7")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medicine-types/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// List Deleted
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/medicine-types/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_MedicineType_HardDelete(t *testing.T) {
	r, cfg, _ := setupMedicineTypeRouter()
	id := createMedicineTypeAndGetID(r, cfg, "MT8")
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medicine-types/%d/hard-delete", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
