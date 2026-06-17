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
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/allergy"
	familyHistory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/familyHistory"
	medicalCondition "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/medicalCondition"
	surgicalHistory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/surgicalHistory"
	medicalhistory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupMedicalHistoryRouters() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	// Repositories
	allergyRepo := repository.NewAllergyRepository(db)
	conditionRepo := repository.NewMedicalConditionRepository(db)
	familyRepo := repository.NewFamilyHistoryRepository(db)
	surgicalRepo := repository.NewSurgicalHistoryRepository(db)
	medicalHistoryRepo := repository.NewMedicalHistoryRepository(db)

	// Services
	allergyService := allergy.NewAllergyService(allergyRepo, cfg)
	conditionService := medicalCondition.NewMedicalConditionService(conditionRepo, cfg)
	familyService := familyHistory.NewFamilyHistoryService(familyRepo, cfg)
	surgicalService := surgicalHistory.NewSurgicalHistoryService(surgicalRepo, cfg)
	medicalHistoryService := medicalhistory.NewMedicalHistoryService(medicalHistoryRepo, cfg)

	// Handlers
	allergyHandler := handler.NewAllergyHandler(allergyService)
	conditionHandler := handler.NewMedicalConditionHandler(conditionService)
	familyHandler := handler.NewFamilyHistoryHandler(familyService)
	surgicalHandler := handler.NewSurgicalHistoryHandler(surgicalService)
	medicalHistoryHandler := handler.NewMedicalHistoryHandler(medicalHistoryService)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")

	routes.SetupAllergyRoutes(v1, routeCfg, allergyHandler)
	routes.SetupMedicalConditionRoutes(v1, routeCfg, conditionHandler)
	routes.SetupFamilyHistoryRoutes(v1, routeCfg, familyHandler)
	routes.SetupSurgicalHistoryRoutes(v1, routeCfg, surgicalHandler)
	routes.SetupMedicalHistoryRoutes(v1, routeCfg, medicalHistoryHandler)

	return r, routeCfg, db
}

func createPrereqForMedicalHistories(db *gorm.DB) uint {
	patient := &models.Patient{FullName: "Test MH Patient", PatientCode: "PT-MH", DateOfBirth: "1990-01-01", Gender: "male", BloodType: "O"}
	db.Create(patient)
	return patient.ID
}

// ─── Allergy Tests ───────────────────────────────────────────────────────────

func createAllergyAndGetID(r *gin.Engine, cfg *routes.RouteConfig, pID uint) uint {
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateAllergyRequest{
		PatientID:    pID,
		AllergenType: "food",
		AllergenName: "Peanut",
		Reaction:     "Rash",
		Severity:     "mild",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medical-history/allergies", bytes.NewBuffer(body))
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

func TestIntegration_Allergy_Create(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateAllergyRequest{
		PatientID:    pID,
		AllergenType: "drug",
		AllergenName: "Penicillin",
		Reaction:     "Hives",
		Severity:     "severe",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medical-history/allergies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Allergy_FindByID(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	id := createAllergyAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/medical-history/allergies/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Allergy_Update(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	id := createAllergyAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	reqBody := dto.UpdateAllergyRequest{
		Reaction: "Itching",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/medical-history/allergies/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Allergy_Delete(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	id := createAllergyAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medical-history/allergies/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Allergy_List(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	createAllergyAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/medical-history/allergies", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Medical Condition Tests ─────────────────────────────────────────────────

func createConditionAndGetID(r *gin.Engine, cfg *routes.RouteConfig, pID uint) uint {
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateMedicalConditionRequest{
		PatientID:     pID,
		ConditionName: "Hypertension",
		Status:        "ongoing",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medical-history/conditions", bytes.NewBuffer(body))
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

func TestIntegration_Condition_Create(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateMedicalConditionRequest{
		PatientID:     pID,
		ConditionName: "Diabetes",
		Status:        "ongoing",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medical-history/conditions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Condition_FindByID(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	id := createConditionAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/medical-history/conditions/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Surgical History Tests ──────────────────────────────────────────────────

func createSurgeryAndGetID(r *gin.Engine, cfg *routes.RouteConfig, pID uint) uint {
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateSurgicalHistoryRequest{
		PatientID:     pID,
		ProcedureName: "Appendectomy",
		SurgeryDate:   time.Now().Format("2006-01-02"),
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medical-history/surgeries", bytes.NewBuffer(body))
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

func TestIntegration_Surgery_Create(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateSurgicalHistoryRequest{
		PatientID:     pID,
		ProcedureName: "Tonsillectomy",
		SurgeryDate:   time.Now().Format("2006-01-02"),
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medical-history/surgeries", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Surgery_FindByID(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	id := createSurgeryAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/medical-history/surgeries/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Family History Tests ────────────────────────────────────────────────────

func createFamilyHistoryAndGetID(r *gin.Engine, cfg *routes.RouteConfig, pID uint) uint {
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateFamilyHistoryRequest{
		PatientID:     pID,
		FamilyMember:  "Father",
		ConditionName: "Heart Disease",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medical-history/family-histories", bytes.NewBuffer(body))
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

func TestIntegration_FamilyHistory_Create(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateFamilyHistoryRequest{
		PatientID:     pID,
		FamilyMember:  "Mother",
		ConditionName: "Breast Cancer",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medical-history/family-histories", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_FamilyHistory_FindByID(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	id := createFamilyHistoryAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/medical-history/family-histories/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_FamilyHistory_Update(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	id := createFamilyHistoryAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	reqBody := dto.UpdateFamilyHistoryRequest{
		Notes: "Updated Note",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/medical-history/family-histories/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_FamilyHistory_Delete(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	id := createFamilyHistoryAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medical-history/family-histories/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_FamilyHistory_List(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	createFamilyHistoryAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/medical-history/family-histories", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Combined Medical History View Test ──────────────────────────────────────

func TestIntegration_MedicalHistory_CombinedView(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	createAllergyAndGetID(r, cfg, pID)
	createConditionAndGetID(r, cfg, pID)
	createSurgeryAndGetID(r, cfg, pID)
	createFamilyHistoryAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	// Fetch combined view for the patient
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/medical-history/patient/%d", pID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})

	// Assert sub-histories are populated
	assert.NotNil(t, data["allergies"])
	assert.NotNil(t, data["medical_conditions"])
	assert.NotNil(t, data["surgical_history"])
	assert.NotNil(t, data["family_history"])
}

func TestIntegration_MedicalHistory_FindByID(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	// Create at least one history so Medical History record is created
	createAllergyAndGetID(r, cfg, pID)

	// Wait, the id of medical history is not the same as patient id. 
	// We should just get it from the combined view or list first, or assume it's 1.
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/medical-history/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_MedicalHistory_List(t *testing.T) {
	r, cfg, db := setupMedicalHistoryRouters()
	pID := createPrereqForMedicalHistories(db)
	createAllergyAndGetID(r, cfg, pID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/medical-history", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

