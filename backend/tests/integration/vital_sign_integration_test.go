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
	vital_sign "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/vital-sign"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupVitalSignRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewVitalSignRepository(db)
	mrRepo := repository.NewMedicalRecordRepository(db)
	service := vital_sign.NewVitalSignService(repo, cfg, mrRepo)
	handlerInst := handler.NewVitalSignHandler(service)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupVitalSignRoutes(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createPrereqForVitalSign(db *gorm.DB) uint {
	patient := &models.Patient{FullName: "Test VS Patient", PatientCode: "PT-VS", DateOfBirth: "1990-01-01", Gender: "male", BloodType: "O"}
	db.Create(patient)

	department := &models.Department{Name: "General"}
	db.Create(department)

	doctor := &models.Doctor{FullName: "Test VS Doctor", DepartmentID: &department.ID}
	db.Create(doctor)

	mr := &models.MedicalRecord{
		PatientID:      patient.ID,
		DoctorID:       doctor.ID,
		VisitDate:      time.Now().Format("2006-01-02"),
		ChiefComplaint: "Fever",
		Diagnosis:      "Flu",
		TreatmentPlan:  "Rest",
	}
	db.Create(mr)

	return mr.ID
}

func createVitalSignAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, mrID uint) uint {
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	sys := 120
	dia := 80
	reqBody := dto.CreateVitalSignRequest{
		MedicalRecordID:        mrID,
		RecordedAt:             time.Now().Format("2006-01-02T15:04:05Z07:00"),
		BloodPressureSystolic:  &sys,
		BloodPressureDiastolic: &dia,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/vital-signs", bytes.NewBuffer(body))
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

func TestIntegration_VitalSign_Create(t *testing.T) {
	r, cfg, db := setupVitalSignRouter()
	mrID := createPrereqForVitalSign(db)
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	sys := 120
	reqBody := dto.CreateVitalSignRequest{
		MedicalRecordID:       mrID,
		RecordedAt:            time.Now().Format("2006-01-02T15:04:05Z07:00"),
		BloodPressureSystolic: &sys,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/vital-signs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_VitalSign_FindByID(t *testing.T) {
	r, cfg, db := setupVitalSignRouter()
	mrID := createPrereqForVitalSign(db)
	id := createVitalSignAndGetID(r, cfg, db, mrID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/vital-signs/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_VitalSign_Update(t *testing.T) {
	r, cfg, db := setupVitalSignRouter()
	mrID := createPrereqForVitalSign(db)
	id := createVitalSignAndGetID(r, cfg, db, mrID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	sys := 125
	reqBody := dto.UpdateVitalSignRequest{
		BloodPressureSystolic: &sys,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/vital-signs/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_VitalSign_List(t *testing.T) {
	r, cfg, db := setupVitalSignRouter()
	mrID := createPrereqForVitalSign(db)
	createVitalSignAndGetID(r, cfg, db, mrID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/vital-signs", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_VitalSign_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupVitalSignRouter()
	mrID := createPrereqForVitalSign(db)
	id := createVitalSignAndGetID(r, cfg, db, mrID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/vital-signs/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/vital-signs/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_VitalSign_DeletedListAndHardDelete(t *testing.T) {
	r, cfg, db := setupVitalSignRouter()
	mrID := createPrereqForVitalSign(db)
	id := createVitalSignAndGetID(r, cfg, db, mrID)

	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/vital-signs/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Deleted List
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/vital-signs/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Hard Delete
	req3, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/vital-signs/%d/hard-delete", id), nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}
