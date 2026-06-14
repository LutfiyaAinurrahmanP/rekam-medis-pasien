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
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medicine"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupMedicineRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewMedicineRepository(db)
	service := medicine.NewMedicineService(repo, cfg)
	handlerInst := handler.NewMedicineHandler(service)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupMedicineRouter(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createPrerequisiteMedicineType(db *gorm.DB) uint {
	mt := &models.MedicineType{
		Name: "Syrup",
		Code: "SYR",
	}
	db.Create(mt)
	return mt.ID
}

func createMedicineAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, name string) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	mtID := createPrerequisiteMedicineType(db)

	p := 15000.0
	sq := 100
	isActive := true
	reqBody := dto.CreateMedicineRequest{
		Name:           name,
		GenericName:    "Generic " + name,
		MedicineTypeID: &mtID,
		StockQuantity:  &sq,
		Price:          &p,
		IsActive:       &isActive,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medicines", bytes.NewBuffer(body))
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

func TestIntegration_Medicine_Create(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	mtID := createPrerequisiteMedicineType(db)

	p := 25000.0
	sq := 50
	reqBody := dto.CreateMedicineRequest{
		Name:           "Paracetamol",
		MedicineTypeID: &mtID,
		StockQuantity:  &sq,
		Price:          &p,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medicines", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Medicine_FindByID(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	id := createMedicineAndGetID(r, cfg, db, "Med1")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/medicines/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Medicine_List(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	createMedicineAndGetID(r, cfg, db, "Med2")
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/medicines", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Medicine_ActiveList(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	createMedicineAndGetID(r, cfg, db, "Med3")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/medicines/active", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Medicine_AvailableList(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	createMedicineAndGetID(r, cfg, db, "Med4")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/medicines/available", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Medicine_LowStockAndOutStockList(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	id := createMedicineAndGetID(r, cfg, db, "Med5")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Set stock to 0 to trigger out of stock / low stock
	reqBodyStock := dto.ReduceStockRequest{Quantity: 100}
	bodyS, _ := json.Marshal(reqBodyStock)
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medicines/%d/reduce-stock", id), bytes.NewBuffer(bodyS))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Low Stock
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/medicines/low-stock", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Out of Stock
	req3, _ := http.NewRequest(http.MethodGet, "/api/v1/medicines/out-of-stock", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestIntegration_Medicine_Update(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	id := createMedicineAndGetID(r, cfg, db, "Med6")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	newName := "Updated Medicine"
	reqBody := dto.UpdateMedicineRequest{
		Name: &newName,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/medicines/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Medicine_AddAndReduceStock(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	id := createMedicineAndGetID(r, cfg, db, "Med7")
	token := GenerateTestToken(1, models.RoleReceptionist, cfg.Config)

	// Add Stock
	reqBodyAdd := dto.AddStockRequest{Quantity: 50}
	bodyA, _ := json.Marshal(reqBodyAdd)
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medicines/%d/add-stock", id), bytes.NewBuffer(bodyA))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Reduce Stock
	reqBodyReduce := dto.ReduceStockRequest{Quantity: 30}
	bodyR, _ := json.Marshal(reqBodyReduce)
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medicines/%d/reduce-stock", id), bytes.NewBuffer(bodyR))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Medicine_DeactivateAndActivate(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	id := createMedicineAndGetID(r, cfg, db, "Med8")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Deactivate
	reqBody := map[string]interface{}{}
	body, _ := json.Marshal(reqBody)
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medicines/%d/deactivate", id), bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Activate
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medicines/%d/activate", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Medicine_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	id := createMedicineAndGetID(r, cfg, db, "Med9")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medicines/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medicines/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Medicine_DeletedList(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	id := createMedicineAndGetID(r, cfg, db, "Med10")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medicines/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// List Deleted
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/medicines/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Medicine_HardDelete(t *testing.T) {
	r, cfg, db := setupMedicineRouter()
	id := createMedicineAndGetID(r, cfg, db, "Med11")
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	// Soft delete first
	reqSD, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medicines/%d", id), nil)
	reqSD.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqSD)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medicines/%d/hard-delete", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
