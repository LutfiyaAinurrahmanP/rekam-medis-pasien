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
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/prescription"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupPrescriptionRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewPrescriptionRepository(db)
	service := prescription.NewPrescriptionService(repo, cfg)
	docRepo := repository.NewDoctorRepository(db)
	patRepo := repository.NewPatientRepository(db)
	handlerInst := handler.NewPrescriptionHandler(service, docRepo, patRepo)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupPrescriptionRoutes(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createPrereqForPrescription(db *gorm.DB) (uint, uint, uint) {
	patient := &models.Patient{FullName: "Test P Patient", PatientCode: "PT-P", DateOfBirth: "1990-01-01", Gender: "male", BloodType: "O"}
	db.Create(patient)

	department := &models.Department{Name: "General"}
	db.Create(department)

	doctor := &models.Doctor{FullName: "Test P Doctor", DepartmentID: &department.ID}
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

	mt := &models.MedicineType{Name: "Tablet"}
	db.Create(mt)

	med := &models.Medicine{Name: "Paracetamol", MedicineTypeID: mt.ID, IsActive: true}
	db.Create(med)

	return mr.ID, doctor.ID, med.ID
}

func createPrescriptionAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, mrID, dID uint) uint {
	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreatePrescriptionRequest{
		MedicalRecordID:  mrID,
		DoctorID:         dID,
		PrescriptionDate: time.Now().Format("2006-01-02"),
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/prescriptions", bytes.NewBuffer(body))
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

func TestIntegration_Prescription_Create(t *testing.T) {
	r, cfg, db := setupPrescriptionRouter()
	mrID, dID, _ := createPrereqForPrescription(db)
	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)

	reqBody := dto.CreatePrescriptionRequest{
		MedicalRecordID:  mrID,
		DoctorID:         dID,
		PrescriptionDate: time.Now().Format("2006-01-02"),
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/prescriptions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Prescription_FindByID(t *testing.T) {
	r, cfg, db := setupPrescriptionRouter()
	mrID, dID, _ := createPrereqForPrescription(db)
	id := createPrescriptionAndGetID(r, cfg, db, mrID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/prescriptions/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Prescription_Update(t *testing.T) {
	r, cfg, db := setupPrescriptionRouter()
	mrID, dID, _ := createPrereqForPrescription(db)
	id := createPrescriptionAndGetID(r, cfg, db, mrID, dID)

	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)
	notes := "Take after meal"
	reqBody := dto.UpdatePrescriptionRequest{
		Notes: &notes,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/prescriptions/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Prescription_Items(t *testing.T) {
	r, cfg, db := setupPrescriptionRouter()
	mrID, dID, medID := createPrereqForPrescription(db)
	id := createPrescriptionAndGetID(r, cfg, db, mrID, dID)

	token := GenerateTestToken(dID, models.RoleDoctor, cfg.Config)

	// Create Item
	reqBody := dto.CreatePrescriptionItemRequest{
		MedicineID:   medID,
		Dosage:       "500mg",
		Frequency:    "3x1",
		DurationDays: 3,
		Quantity:     9,
	}
	body, _ := json.Marshal(reqBody)
	req1, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/prescriptions/%d/items", id), bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	var res map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &res)
	itemID := uint(res["data"].(map[string]interface{})["id"].(float64))

	// List Items
	req2, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/prescriptions/%d/items", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Get Item
	req3, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/prescriptions/%d/items/%d", id, itemID), nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	// Update Item
	dos := "600mg"
	reqBodyU := dto.UpdatePrescriptionItemRequest{Dosage: &dos}
	bodyU, _ := json.Marshal(reqBodyU)
	req4, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/prescriptions/%d/items/%d", id, itemID), bytes.NewBuffer(bodyU))
	req4.Header.Set("Content-Type", "application/json")
	req4.Header.Set("Authorization", "Bearer "+token)
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)

	// Delete Item
	req5, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/prescriptions/%d/items/%d", id, itemID), nil)
	req5.Header.Set("Authorization", "Bearer "+token)
	w5 := httptest.NewRecorder()
	r.ServeHTTP(w5, req5)
	assert.Equal(t, http.StatusOK, w5.Code)
}

func TestIntegration_Prescription_Lifecycle(t *testing.T) {
	r, cfg, db := setupPrescriptionRouter()
	mrID, dID, _ := createPrereqForPrescription(db)
	id := createPrescriptionAndGetID(r, cfg, db, mrID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Dispense
	reqD, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/prescriptions/%d/dispense", id), nil)
	reqD.Header.Set("Authorization", "Bearer "+token)
	wD := httptest.NewRecorder()
	r.ServeHTTP(wD, reqD)
	assert.Equal(t, http.StatusOK, wD.Code)

	id2 := createPrescriptionAndGetID(r, cfg, db, mrID, dID)

	// Cancel
	reqC, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/prescriptions/%d/cancel", id2), nil)
	reqC.Header.Set("Authorization", "Bearer "+token)
	wC := httptest.NewRecorder()
	r.ServeHTTP(wC, reqC)
	assert.Equal(t, http.StatusOK, wC.Code)
}

func TestIntegration_Prescription_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupPrescriptionRouter()
	mrID, dID, _ := createPrereqForPrescription(db)
	id := createPrescriptionAndGetID(r, cfg, db, mrID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Cancel before delete if needed, though status pending might be deletable.
	reqC, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/prescriptions/%d/cancel", id), nil)
	reqC.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqC)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/prescriptions/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/prescriptions/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Prescription_DeletedListAndHardDelete(t *testing.T) {
	r, cfg, db := setupPrescriptionRouter()
	mrID, dID, _ := createPrereqForPrescription(db)
	id := createPrescriptionAndGetID(r, cfg, db, mrID, dID)

	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	reqC, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/prescriptions/%d/cancel", id), nil)
	reqC.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqC)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/prescriptions/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Deleted List
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/prescriptions/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Hard Delete
	req3, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/prescriptions/%d/hard-delete", id), nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}
