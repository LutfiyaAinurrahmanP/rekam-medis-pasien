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
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/hospitalization"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupHospitalizationRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewHospitalizationRepository(db)
	service := hospitalization.NewHospitalizationService(repo, cfg)
	handlerInst := handler.NewHospitalizationHandler(service)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupHospitalizationRouter(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createPrereqForHospitalization(db *gorm.DB) (uint, uint, uint) {
	patient := &models.Patient{FullName: "Test Hosp Patient", PatientCode: "PT-HOSP", DateOfBirth: "1990-01-01", Gender: "male", BloodType: "O"}
	db.Create(patient)

	department := &models.Department{Name: "General"}
	db.Create(department)

	doctor := &models.Doctor{FullName: "Test Hosp Doctor", DepartmentID: &department.ID}
	db.Create(doctor)

	roomType := &models.RoomType{Name: "ICU"}
	db.Create(roomType)

	room := &models.Room{RoomNumber: "101", RoomTypeID: &roomType.ID, IsActive: true}
	db.Create(room)

	return patient.ID, doctor.ID, room.ID
}

func createHospitalizationAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, pID, dID, rID uint) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := dto.CreateHospitalizationRequest{
		PatientID:         pID,
		AttendingDoctorID: dID,
		RoomID:            rID,
		AdmissionDate:     time.Now().Format("2006-01-02T15:04:05Z"),
		AdmissionReason:   "Severe fever",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/hospitalizations", bytes.NewBuffer(body))
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

func TestIntegration_Hospitalization_Create(t *testing.T) {
	r, cfg, db := setupHospitalizationRouter()
	pID, dID, rID := createPrereqForHospitalization(db)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := dto.CreateHospitalizationRequest{
		PatientID:         pID,
		AttendingDoctorID: dID,
		RoomID:            rID,
		AdmissionDate:     time.Now().Format("2006-01-02T15:04:05Z"),
		AdmissionReason:   "Severe fever",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/hospitalizations", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Hospitalization_FindByID(t *testing.T) {
	r, cfg, db := setupHospitalizationRouter()
	pID, dID, rID := createPrereqForHospitalization(db)
	id := createHospitalizationAndGetID(r, cfg, db, pID, dID, rID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/hospitalizations/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Hospitalization_Update(t *testing.T) {
	r, cfg, db := setupHospitalizationRouter()
	pID, dID, rID := createPrereqForHospitalization(db)
	id := createHospitalizationAndGetID(r, cfg, db, pID, dID, rID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	reason := "Updated reason"
	reqBody := dto.UpdateHospitalizationRequest{
		AdmissionReason: &reason,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/hospitalizations/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Hospitalization_DischargeAndTransfer(t *testing.T) {
	r, cfg, db := setupHospitalizationRouter()
	pID, dID, rID := createPrereqForHospitalization(db)
	id := createHospitalizationAndGetID(r, cfg, db, pID, dID, rID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Transfer
	reqBodyT := dto.TransferHospitalizationRequest{Notes: "Transfer to new room"}
	bodyT, _ := json.Marshal(reqBodyT)
	reqT, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/hospitalizations/%d/transfer", id), bytes.NewBuffer(bodyT))
	reqT.Header.Set("Content-Type", "application/json")
	reqT.Header.Set("Authorization", "Bearer "+token)
	wT := httptest.NewRecorder()
	r.ServeHTTP(wT, reqT)
	assert.Equal(t, http.StatusOK, wT.Code)

	// Discharge
	reqBodyD := dto.DischargeHospitalizationRequest{DischargeSummary: "Patient recovered"}
	bodyD, _ := json.Marshal(reqBodyD)
	reqD, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/hospitalizations/%d/discharge", id), bytes.NewBuffer(bodyD))
	reqD.Header.Set("Content-Type", "application/json")
	reqD.Header.Set("Authorization", "Bearer "+token)
	wD := httptest.NewRecorder()
	r.ServeHTTP(wD, reqD)
	assert.Equal(t, http.StatusOK, wD.Code)
}

func TestIntegration_Hospitalization_Lists(t *testing.T) {
	r, cfg, db := setupHospitalizationRouter()
	pID, dID, rID := createPrereqForHospitalization(db)
	createHospitalizationAndGetID(r, cfg, db, pID, dID, rID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	endpoints := []string{
		"/api/v1/hospitalizations",
		"/api/v1/hospitalizations/active",
		"/api/v1/hospitalizations/inactive",
	}

	for _, ep := range endpoints {
		req, _ := http.NewRequest(http.MethodGet, ep, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, ep)
	}
}

func TestIntegration_Hospitalization_ActivateAndDeactivate(t *testing.T) {
	r, cfg, db := setupHospitalizationRouter()
	pID, dID, rID := createPrereqForHospitalization(db)
	id := createHospitalizationAndGetID(r, cfg, db, pID, dID, rID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Deactivate
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/hospitalizations/%d/deactivate", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Activate
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/hospitalizations/%d/activate", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Hospitalization_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupHospitalizationRouter()
	pID, dID, rID := createPrereqForHospitalization(db)
	id := createHospitalizationAndGetID(r, cfg, db, pID, dID, rID)
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Deactivate before delete
	reqD, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/hospitalizations/%d/deactivate", id), nil)
	reqD.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqD)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/hospitalizations/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/hospitalizations/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Hospitalization_DeletedListAndHardDelete(t *testing.T) {
	r, cfg, db := setupHospitalizationRouter()
	pID, dID, rID := createPrereqForHospitalization(db)
	id := createHospitalizationAndGetID(r, cfg, db, pID, dID, rID)
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	// Deactivate before delete
	reqD, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/hospitalizations/%d/deactivate", id), nil)
	reqD.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), reqD)

	// Soft Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/hospitalizations/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// Deleted List
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/hospitalizations/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Hard Delete
	req3, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/hospitalizations/%d/hard-delete", id), nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)
}
