package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/routes"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/dashboard"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupDashboardRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardSvc := dashboard.NewDashboardService(dashboardRepo)
	
	doctorRepo := repository.NewDoctorRepository(db)
	patientRepo := repository.NewPatientRepository(db)

	handlerInst := handler.NewDashboardHandler(dashboardSvc, doctorRepo, patientRepo)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupDashboardRouter(v1, routeCfg, handlerInst)

	return r, routeCfg, db
}

func createPrereqForDashboard(db *gorm.DB) (uint, uint, uint, uint) {
	// Admin User
	adminUser := &models.User{Username: "admin_dash", Email: "admin@dash.com", Phone: "1231", Role: models.RoleAdmin}
	db.Create(adminUser)

	// Doctor User
	dept := &models.Department{Name: "General Dash"}
	db.Create(dept)
	spec := &models.DoctorSpecialization{Name: "General Practitioner Dash"}
	db.Create(spec)

	doctorUser := &models.User{Username: "doc_dash", Email: "doc@dash.com", Phone: "1232", Role: models.RoleDoctor}
	db.Create(doctorUser)
	doctor := &models.Doctor{
		FullName:         "Dr. Dashboard",
		UserID:           &doctorUser.ID,
		DepartmentID:     &dept.ID,
		SpecializationID: spec.ID,
		EmployeeID:       "EMP-DASH",
		LicenseNumber:    "LIC-DASH",
	}
	db.Create(doctor)

	// Patient User
	patientUser := &models.User{Username: "pat_dash", Email: "pat@dash.com", Phone: "1233", Role: models.RolePatient}
	db.Create(patientUser)
	patient := &models.Patient{FullName: "Patient Dashboard", UserID: &patientUser.ID}
	db.Create(patient)

	// Receptionist User
	receptionistUser := &models.User{Username: "recep_dash", Email: "recep@dash.com", Phone: "1234", Role: models.RoleReceptionist}
	db.Create(receptionistUser)

	return adminUser.ID, doctorUser.ID, patientUser.ID, receptionistUser.ID
}

// ─── Tests ───────────────────────────────────────────────────────────────

func TestIntegration_Dashboard_Overview(t *testing.T) {
	r, cfg, db := setupDashboardRouter()
	adminID, _, _, _ := createPrereqForDashboard(db)

	token := GenerateTestToken(adminID, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/overview", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Dashboard_Admin(t *testing.T) {
	r, cfg, db := setupDashboardRouter()
	adminID, _, _, _ := createPrereqForDashboard(db)

	token := GenerateTestToken(adminID, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Dashboard_Doctor(t *testing.T) {
	r, cfg, db := setupDashboardRouter()
	_, docID, _, _ := createPrereqForDashboard(db)

	token := GenerateTestToken(docID, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/doctor", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Dashboard_Receptionist(t *testing.T) {
	r, cfg, db := setupDashboardRouter()
	_, _, _, recID := createPrereqForDashboard(db)

	token := GenerateTestToken(recID, models.RoleReceptionist, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/receptionist", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Dashboard_Patient(t *testing.T) {
	r, cfg, db := setupDashboardRouter()
	_, _, patID, _ := createPrereqForDashboard(db)

	token := GenerateTestToken(patID, models.RolePatient, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/patient", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Dashboard_RevenueReport(t *testing.T) {
	r, cfg, db := setupDashboardRouter()
	adminID, _, _, _ := createPrereqForDashboard(db)

	token := GenerateTestToken(adminID, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/reports/revenue", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Dashboard_AppointmentReport(t *testing.T) {
	r, cfg, db := setupDashboardRouter()
	_, docID, _, _ := createPrereqForDashboard(db)

	token := GenerateTestToken(docID, models.RoleDoctor, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/reports/appointments", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Dashboard_PatientReport(t *testing.T) {
	r, cfg, db := setupDashboardRouter()
	adminID, _, _, _ := createPrereqForDashboard(db)

	token := GenerateTestToken(adminID, models.RoleAdmin, cfg.Config)
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/dashboard/reports/patients", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
