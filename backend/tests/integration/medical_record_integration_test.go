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
	medical_record "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-record"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupMedicalRecordRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewMedicalRecordRepository(db)
	service := medical_record.NewMedicalRecordService(repo, cfg)
	docRepo := repository.NewDoctorRepository(db)
	patRepo := repository.NewPatientRepository(db)
	handlerInst := handler.NewMedicalRecordHandler(service, docRepo, patRepo)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupMedicalRecordRouter(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createPrereqForMedicalRecord(db *gorm.DB) (uint, uint) {
	patient := &models.Patient{
		FullName:    "Test MR Patient",
		PatientCode: "PT-MR",
		DateOfBirth: "1990-01-01",
		Gender:      "male",
		BloodType:   "O",
	}
	db.Create(patient)

	department := &models.Department{Name: "General"}
	db.Create(department)

	doctor := &models.Doctor{
		FullName:     "Test MR Doctor",
		DepartmentID: &department.ID,
	}
	db.Create(doctor)

	return patient.ID, doctor.ID
}

func createMedicalRecordAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, pID uint, dID uint) uint {
	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateMedicalRecordRequest{
		PatientID:      pID,
		VisitDate:      time.Now().Format("2006-01-02"),
		ChiefComplaint: "Headache",
		Diagnosis:      "Migraine",
		TreatmentPlan:  "Rest and medicine",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medical-records", bytes.NewBuffer(body))
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

func TestIntegration_MedicalRecord_Create(t *testing.T) {
	r, cfg, db := setupMedicalRecordRouter()
	pID, dID := createPrereqForMedicalRecord(db)
	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateMedicalRecordRequest{
		PatientID:      pID,
		VisitDate:      time.Now().Format("2006-01-02"),
		ChiefComplaint: "Headache",
		Diagnosis:      "Migraine",
		TreatmentPlan:  "Rest and medicine",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/medical-records", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_MedicalRecord_FindByID(t *testing.T) {
	r, cfg, db := setupMedicalRecordRouter()
	pID, dID := createPrereqForMedicalRecord(db)
	id := createMedicalRecordAndGetID(r, cfg, db, pID, dID)

	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/medical-records/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_MedicalRecord_Update(t *testing.T) {
	r, cfg, db := setupMedicalRecordRouter()
	pID, dID := createPrereqForMedicalRecord(db)
	id := createMedicalRecordAndGetID(r, cfg, db, pID, dID)

	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)
	diagnosis := "Severe Migraine"
	reqBody := dto.UpdateMedicalRecordRequest{
		Diagnosis: &diagnosis,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/medical-records/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_MedicalRecord_Finalize(t *testing.T) {
	r, cfg, db := setupMedicalRecordRouter()
	pID, dID := createPrereqForMedicalRecord(db)
	id := createMedicalRecordAndGetID(r, cfg, db, pID, dID)

	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medical-records/%d/finalize", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_MedicalRecord_Lists(t *testing.T) {
	r, cfg, db := setupMedicalRecordRouter()
	pID, dID := createPrereqForMedicalRecord(db)
	createMedicalRecordAndGetID(r, cfg, db, pID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// List all
	req1, _ := http.NewRequest(http.MethodGet, "/api/v1/medical-records", nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// By Patient ID
	req2, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/medical-records/patient/%d", pID), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_MedicalRecord_MyRecords(t *testing.T) {
	r, cfg, db := setupMedicalRecordRouter()
	pID, dID := createPrereqForMedicalRecord(db)
	createMedicalRecordAndGetID(r, cfg, db, pID, dID)

	token := GenerateTestToken(pID, models.RolePatient, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/medical-records/my-records", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_MedicalRecord_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupMedicalRecordRouter()
	pID, dID := createPrereqForMedicalRecord(db)
	id := createMedicalRecordAndGetID(r, cfg, db, pID, dID)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medical-records/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/medical-records/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_MedicalRecord_DeletedListAndHardDelete(t *testing.T) {
	r, cfg, db := setupMedicalRecordRouter()
	pID, dID := createPrereqForMedicalRecord(db)
	id := createMedicalRecordAndGetID(r, cfg, db, pID, dID)
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medical-records/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Deleted List
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/medical-records/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Hard Delete
	req3, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/medical-records/%d/hard-delete", id), nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}
