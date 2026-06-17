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
	labtest "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/lab-test"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupLabTestRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewLabTestRepository(db)
	service := labtest.NewLabTestService(repo, cfg)
	docRepo := repository.NewDoctorRepository(db)
	handlerInst := handler.NewLabTestHandler(service, docRepo)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupLabTestRouter(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createPrereqForLabTest(db *gorm.DB) (uint, uint, uint) {
	patient := &models.Patient{FullName: "Test Lab Patient", PatientCode: "PT-LAB", DateOfBirth: "1990-01-01", Gender: "male", BloodType: "O"}
	db.Create(patient)

	department := &models.Department{Name: "General"}
	db.Create(department)

	doctor := &models.Doctor{FullName: "Test Lab Doctor", DepartmentID: &department.ID}
	db.Create(doctor)

	mr := &models.MedicalRecord{
		PatientID:      patient.ID,
		DoctorID:       doctor.ID,
		VisitDate:      time.Now().Format("2006-01-02"),
		ChiefComplaint: "Fever",
		Diagnosis:      "Unknown",
		TreatmentPlan:  "Lab test required",
	}
	db.Create(mr)

	ttc := &models.TypeTestCategory{Name: "Blood"}
	db.Create(ttc)

	tt := &models.TypeTest{Name: "Complete Blood Count", Code: "CBC", TypeTestCategoryID: ttc.ID}
	db.Create(tt)

	return mr.ID, doctor.ID, tt.ID
}

func createLabTestAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, mrID, dID, ttID uint) uint {
	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateLabTestRequest{
		MedicalRecordID:   mrID,
		TestTypeID:        ttID,
		OrderedByDoctorID: dID,
		OrderDate:         time.Now().Format("2006-01-02"),
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/lab-tests", bytes.NewBuffer(body))
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

func TestIntegration_LabTest_Create(t *testing.T) {
	r, cfg, db := setupLabTestRouter()
	mrID, dID, ttID := createPrereqForLabTest(db)
	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreateLabTestRequest{
		MedicalRecordID:   mrID,
		TestTypeID:        ttID,
		OrderedByDoctorID: dID,
		OrderDate:         time.Now().Format("2006-01-02"),
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/lab-tests", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_LabTest_FindByID(t *testing.T) {
	r, cfg, db := setupLabTestRouter()
	mrID, dID, ttID := createPrereqForLabTest(db)
	id := createLabTestAndGetID(r, cfg, db, mrID, dID, ttID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/lab-tests/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_LabTest_Update(t *testing.T) {
	r, cfg, db := setupLabTestRouter()
	mrID, dID, ttID := createPrereqForLabTest(db)
	id := createLabTestAndGetID(r, cfg, db, mrID, dID, ttID)

	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)
	notes := "Fasting required"
	reqBody := dto.UpdateLabTestRequest{
		Notes: &notes,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/lab-tests/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_LabTest_Lifecycle(t *testing.T) {
	r, cfg, db := setupLabTestRouter()
	mrID, dID, ttID := createPrereqForLabTest(db)
	id := createLabTestAndGetID(r, cfg, db, mrID, dID, ttID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Collect Sample
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/lab-tests/%d/collect-sample", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Start
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/lab-tests/%d/start", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Complete
	resVal := "Normal"
	reqBody := dto.CompleteLabTestRequest{ResultValue: &resVal}
	body, _ := json.Marshal(reqBody)
	req3, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/lab-tests/%d/complete", id), bytes.NewBuffer(body))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}

func TestIntegration_LabTest_Cancel(t *testing.T) {
	r, cfg, db := setupLabTestRouter()
	mrID, dID, ttID := createPrereqForLabTest(db)
	id := createLabTestAndGetID(r, cfg, db, mrID, dID, ttID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/lab-tests/%d/cancel", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_LabTest_Lists(t *testing.T) {
	r, cfg, db := setupLabTestRouter()
	mrID, dID, ttID := createPrereqForLabTest(db)
	createLabTestAndGetID(r, cfg, db, mrID, dID, ttID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// List all
	req1, _ := http.NewRequest(http.MethodGet, "/api/v1/lab-tests", nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// By MR ID
	req2, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/lab-tests/medical-record/%d", mrID), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_LabTest_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupLabTestRouter()
	mrID, dID, ttID := createPrereqForLabTest(db)
	id := createLabTestAndGetID(r, cfg, db, mrID, dID, ttID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Cancel before delete if needed, but ordered state might be deletable depending on logic.
	// Let's cancel it first just in case.
	reqC, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/lab-tests/%d/cancel", id), nil)
	reqC.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqC)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/lab-tests/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/lab-tests/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_LabTest_DeletedListAndHardDelete(t *testing.T) {
	r, cfg, db := setupLabTestRouter()
	mrID, dID, ttID := createPrereqForLabTest(db)
	id := createLabTestAndGetID(r, cfg, db, mrID, dID, ttID)

	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	reqC, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/lab-tests/%d/cancel", id), nil)
	reqC.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqC)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/lab-tests/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Deleted List
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/lab-tests/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Hard Delete
	req3, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/lab-tests/%d/hard-delete", id), nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}
