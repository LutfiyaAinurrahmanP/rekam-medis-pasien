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
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/referral"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupReferralRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	repo := repository.NewReferralRepository(db)
	service := referral.NewReferralService(repo)
	handlerInst := handler.NewReferralHandler(service)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupReferralRoutes(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createPrereqForReferral(db *gorm.DB) (uint, uint, uint) {
	// Create Patient (ID 1)
	userID := uint(1)
	patient := &models.Patient{FullName: "Test Ref Patient", UserID: &userID}
	db.Create(patient)

	// Create Department
	dept := &models.Department{Name: "General", Code: "GEN"}
	db.Create(dept)

	// Create Specialization
	spec := &models.DoctorSpecialization{Name: "General Practitioner"}
	db.Create(spec)

	doctor := &models.Doctor{
		FullName:         "Dr. Referring",
		UserID:           &userID,
		DepartmentID:     &dept.ID,
		SpecializationID: spec.ID,
		EmployeeID:       "EMP-001",
		LicenseNumber:    "LIC-001",
	}
	db.Create(doctor)

	// Create MedicalRecord (ID 1)
	mr := &models.MedicalRecord{PatientID: patient.ID, DoctorID: doctor.ID, VisitDate: "2026-06-10", Status: "completed"}
	db.Create(mr)

	return patient.ID, mr.ID, doctor.ID
}

func createReferralAndGetID(r *gin.Engine, cfg *routes.RouteConfig, pID, mrID, dID uint) uint {
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)

	deptID := uint(1)
	reqBody := dto.CreateReferralRequest{
		PatientID:              pID,
		MedicalRecordID:        mrID,
		ReferringDoctorID:      dID,
		ReferralDate:           "2026-06-15",
		ReferralType:           "internal",
		Reason:                 "Need specialist consultation",
		Priority:               "routine",
		ReferredToDepartmentID: &deptID,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/referrals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code >= 400 {
		fmt.Println("Error creating referral:", w.Body.String())
	}

	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	if res["data"] == nil {
		return 0
	}
	return uint(res["data"].(map[string]interface{})["id"].(float64))
}

func TestIntegration_Referral_Create(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	reqBody := dto.CreateReferralRequest{
		PatientID:          pID,
		MedicalRecordID:    mrID,
		ReferringDoctorID:  dID,
		ReferralDate:       "2026-06-15",
		ReferralType:       "external",
		Reason:             "Need external MRI",
		Priority:           "urgent",
		ReferredToFacility: "RSCM",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/referrals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code >= 400 {
		fmt.Println("Error creating referral:", w.Body.String())
	}

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Referral_FindByID(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	id := createReferralAndGetID(r, cfg, pID, mrID, dID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/referrals/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Referral_Update(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	id := createReferralAndGetID(r, cfg, pID, mrID, dID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	notes := "Updated Notes"
	reqBody := dto.UpdateReferralRequest{
		Notes: &notes,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/referrals/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Referral_ListAll(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	createReferralAndGetID(r, cfg, pID, mrID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/referrals", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Referral_FindByPatientID(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	createReferralAndGetID(r, cfg, pID, mrID, dID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/referrals/patient/%d", pID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Referral_FindByDoctorID(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	createReferralAndGetID(r, cfg, pID, mrID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/referrals/doctor/%d", dID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Referral_FindMyReferrals(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	createReferralAndGetID(r, cfg, pID, mrID, dID)

	// User ID 1 is assigned to the patient in createPrereqForReferral
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/referrals/my-referrals", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Lifecycle Endpoints ──────────────────────────────────────────────────

func TestIntegration_Referral_Accept(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	id := createReferralAndGetID(r, cfg, pID, mrID, dID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	reqBody := dto.AcceptReferralRequest{Notes: "Accepted"}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/referrals/%d/accept", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Referral_Reject(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	id := createReferralAndGetID(r, cfg, pID, mrID, dID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	reqBody := dto.RejectReferralRequest{RejectionReason: "Not needed"}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/referrals/%d/reject", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Referral_Complete(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	id := createReferralAndGetID(r, cfg, pID, mrID, dID)

	// Accept it first
	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	reqBodyAcc := dto.AcceptReferralRequest{Notes: "Accepted"}
	bodyAcc, _ := json.Marshal(reqBodyAcc)
	reqAcc, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/referrals/%d/accept", id), bytes.NewBuffer(bodyAcc))
	reqAcc.Header.Set("Content-Type", "application/json")
	reqAcc.Header.Set("Authorization", "Bearer "+token)
	wAcc := httptest.NewRecorder()
	r.ServeHTTP(wAcc, reqAcc)

	// Then Complete
	reqBody := dto.CompleteReferralRequest{Notes: "Completed treatment"}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/referrals/%d/complete", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Referral_Cancel(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	id := createReferralAndGetID(r, cfg, pID, mrID, dID)

	token := GenerateTestToken(1, models.RoleDoctor, cfg.Config)
	reqBody := dto.CancelReferralRequest{CancellationReason: "Patient cancelled"}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/referrals/%d/cancel", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ─── Delete, Restore, Hard Delete ──────────────────────────────────────────

func TestIntegration_Referral_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	id := createReferralAndGetID(r, cfg, pID, mrID, dID)

	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	
	// Soft Delete
	reqDel, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/referrals/%d", id), nil)
	reqDel.Header.Set("Authorization", "Bearer "+token)
	wDel := httptest.NewRecorder()
	r.ServeHTTP(wDel, reqDel)
	assert.Equal(t, http.StatusOK, wDel.Code)

	// Restore
	reqRes, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/referrals/%d/restore", id), nil)
	reqRes.Header.Set("Authorization", "Bearer "+token)
	wRes := httptest.NewRecorder()
	r.ServeHTTP(wRes, reqRes)
	assert.Equal(t, http.StatusOK, wRes.Code)
}

func TestIntegration_Referral_DeletedListAndHardDelete(t *testing.T) {
	r, cfg, db := setupReferralRouter()
	pID, mrID, dID := createPrereqForReferral(db)
	id := createReferralAndGetID(r, cfg, pID, mrID, dID)

	tokenAdmin := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	tokenSuperAdmin := GenerateTestToken(2, models.RoleSuperAdmin, cfg.Config)

	// Soft Delete
	reqDel, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/referrals/%d", id), nil)
	reqDel.Header.Set("Authorization", "Bearer "+tokenAdmin)
	wDel := httptest.NewRecorder()
	r.ServeHTTP(wDel, reqDel)

	// Get Deleted List
	reqList, _ := http.NewRequest(http.MethodGet, "/api/v1/referrals/deleted", nil)
	reqList.Header.Set("Authorization", "Bearer "+tokenAdmin)
	wList := httptest.NewRecorder()
	r.ServeHTTP(wList, reqList)
	assert.Equal(t, http.StatusOK, wList.Code)

	// Hard Delete
	reqHard, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/referrals/%d/hard-delete", id), nil)
	reqHard.Header.Set("Authorization", "Bearer "+tokenSuperAdmin)
	wHard := httptest.NewRecorder()
	r.ServeHTTP(wHard, reqHard)
	assert.Equal(t, http.StatusOK, wHard.Code)
}
