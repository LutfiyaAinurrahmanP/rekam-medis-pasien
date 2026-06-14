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
	billingservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/billing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupBillingRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	billingRepo := repository.NewBillingRepository(db)
	billingService := billingservice.NewBillingService(billingRepo)
	billingHandler := handler.NewBillingHandler(billingService)

	routeCfg := &routes.RouteConfig{
		Config:         cfg,
		BillingHandler: billingHandler,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupBillingRoutes(v1, routeCfg, billingHandler)

	return r, routeCfg, db
}

func setupBillingTestPatient(db *gorm.DB) *models.Patient {
	patient := &models.Patient{
		FullName:    "Billing Patient",
		PatientCode: "RM-999",
		DateOfBirth: "1990-01-01",
		Gender:      "male",
		BloodType:   "O",
	}
	db.Create(patient)
	return patient
}

// ==============================================================================
// CREATE Tests
// ==============================================================================

func TestIntegration_Billing_Create_Success(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	receptionistToken := GenerateTestToken(1, models.RoleReceptionist, cfg.Config)
	patient := setupBillingTestPatient(repo)

	reqBody := dto.CreateBillingRequest{
		PatientID:     patient.ID,
		InvoiceNumber: "INV-001",
		BillingDate:   time.Now().Format(time.RFC3339),
		DueDate:       time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
		TotalAmount:   100000,
		Status:        "pending",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/billing", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+receptionistToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Billing_Create_ForbiddenPatient(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	patientToken := GenerateTestToken(1, models.RolePatient, cfg.Config)
	patient := setupBillingTestPatient(repo)

	reqBody := dto.CreateBillingRequest{
		PatientID:     patient.ID,
		InvoiceNumber: "INV-002",
		TotalAmount:   100000,
		Status:        "pending",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/billing", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+patientToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

// Helper to create billing and get ID
func createBillingAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, inv string) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	patient := setupBillingTestPatient(db)

	reqBody := dto.CreateBillingRequest{
		PatientID:     patient.ID,
		InvoiceNumber: inv,
		BillingDate:   time.Now().Format(time.RFC3339),
		DueDate:       time.Now().AddDate(0, 0, 7).Format(time.RFC3339),
		TotalAmount:   500000,
		Status:        "pending",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/billing", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	return uint(res["data"].(map[string]interface{})["id"].(float64))
}

// ==============================================================================
// GET Tests
// ==============================================================================

func TestIntegration_Billing_FindByID(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-100")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config) // All auth users

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/billing/%d", billingID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Billing_FindByInvoiceNumber(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	createBillingAndGetID(r, cfg, repo, "INV-200")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/billing/invoice/INV-200", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Billing_FindByPatientID(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	createBillingAndGetID(r, cfg, repo, "INV-250")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/billing/patient/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Billing_List(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	createBillingAndGetID(r, cfg, repo, "INV-300")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config) // Staff

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/billing", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==============================================================================
// UPDATE & ACTION Tests
// ==============================================================================

func TestIntegration_Billing_Update(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-400")
	token := GenerateTestToken(1, models.RoleReceptionist, cfg.Config)

	reqBody := dto.UpdateBillingRequest{
		TotalAmount: 600000,
		Notes:       "Updated amount",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/billing/%d", billingID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Billing_RecordPayment(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-500")
	token := GenerateTestToken(1, models.RoleReceptionist, cfg.Config)

	reqBody := dto.RecordPaymentRequest{
		PaidAmount:    500000,
		PaymentMethod: "cash",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/billing/%d/pay", billingID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Billing_Cancel(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-600")
	token := GenerateTestToken(1, models.RoleReceptionist, cfg.Config)

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/billing/%d/cancel", billingID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==============================================================================
// ITEM Tests
// ==============================================================================

func createBillingItemAndGetID(r *gin.Engine, cfg *routes.RouteConfig, billingID uint) uint {
	token := GenerateTestToken(1, models.RoleReceptionist, cfg.Config)
	
	reqBody := dto.CreateBillingItemRequest{
		ItemType:    "consultation",
		Description: "General Checkup",
		Quantity:    1,
		UnitPrice:   250000,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/billing/%d/items", billingID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	return uint(res["data"].(map[string]interface{})["id"].(float64))
}

func TestIntegration_Billing_CreateItem(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-700")
	token := GenerateTestToken(1, models.RoleReceptionist, cfg.Config)

	reqBody := dto.CreateBillingItemRequest{
		ItemType:    "medicine",
		Description: "Paracetamol",
		Quantity:    2,
		UnitPrice:   15000,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/billing/%d/items", billingID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Billing_ListItems(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-800")
	createBillingItemAndGetID(r, cfg, billingID)
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/billing/%d/items", billingID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Billing_FindItemByID(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-850")
	itemID := createBillingItemAndGetID(r, cfg, billingID)
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/billing/%d/items/%d", billingID, itemID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Billing_UpdateItem(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-900")
	itemID := createBillingItemAndGetID(r, cfg, billingID)
	token := GenerateTestToken(1, models.RoleReceptionist, cfg.Config)

	reqBody := dto.UpdateBillingItemRequest{
		Quantity:  3,
		UnitPrice: 20000,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/billing/%d/items/%d", billingID, itemID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Billing_DeleteItem(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-1000")
	itemID := createBillingItemAndGetID(r, cfg, billingID)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/billing/%d/items/%d", billingID, itemID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ==============================================================================
// DELETE & RESTORE Tests
// ==============================================================================

func TestIntegration_Billing_SoftDelete(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-1100")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/billing/%d", billingID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Billing_Restore(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-1200")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete first
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/billing/%d", billingID), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/billing/%d/restore", billingID), nil)
	req2.Header.Set("Authorization", "Bearer "+token)

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Billing_DeletedList(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-1300")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete first
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/billing/%d", billingID), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// List deleted
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/billing/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Billing_HardDelete(t *testing.T) {
	r, cfg, repo := setupBillingRouter()
	billingID := createBillingAndGetID(r, cfg, repo, "INV-1400")
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/billing/%d/hard-delete", billingID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
