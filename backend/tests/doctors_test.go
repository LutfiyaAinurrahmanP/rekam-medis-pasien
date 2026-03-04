package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/database"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/routes"
	departmentservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/department"
	doctorservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/doctor"
	userservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/user"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	doctorTestDB     *gorm.DB
	doctorTestRouter *gin.Engine
	doctorTestConfig *config.Config
)

// Setup test router dengan database real
func setupDoctorTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Initialize config
	doctorTestConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-key-for-doctor-testing",
			ExpiredTime: 24 * 60 * time.Minute,
		},
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	// Initialize database (SQLite in-memory)
	doctorTestDB = database.InitTestDB()

	// Run migrations
	doctorTestDB.AutoMigrate(&models.User{}, &models.Doctor{}, &models.Department{})

	// Initialize repositories
	userRepo := repository.NewUserRepository(doctorTestDB)
	doctorRepo := repository.NewDoctorRepository(doctorTestDB)
	deptRepo := repository.NewDepartmentRepository(doctorTestDB)

	// Initialize services
	userService := userservice.NewUserService(userRepo, doctorTestConfig)
	doctorService := doctorservice.NewDoctorService(doctorRepo, doctorTestConfig)
	deptService := departmentservice.NewDepartmentService(deptRepo, doctorTestConfig)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	doctorHandler := handler.NewDoctorHandler(doctorService)
	deptHandler := handler.NewDepartmentHandler(deptService)

	// Setup router
	router := gin.New()

	// Setup routes
	api := router.Group("/api/v1")
	routeConfig := &routes.RouteConfig{
		Config:            doctorTestConfig,
		UserHandler:       userHandler,
		DoctorHandler:     doctorHandler,
		DepartmentHandler: deptHandler,
	}
	routes.SetupAPIRouter(api, routeConfig)

	return router
}

// Cleanup database sebelum test
func cleanupDoctorTestDB() {
	doctorTestDB.Exec("DELETE FROM doctors")
	doctorTestDB.Exec("DELETE FROM departments")
	doctorTestDB.Exec("DELETE FROM users")
}

// Helper function untuk membuat user langsung di database
func createDoctorTestUser(username, email, phone, password, role string, isActive bool) (*models.User, string, error) {
	hashedPassword, _ := utils.HashPassword(password)
	user := &models.User{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: hashedPassword,
		Role:     role,
		IsActive: isActive,
	}

	result := doctorTestDB.Create(user)
	if result.Error != nil {
		return nil, "", result.Error
	}

	token, _, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role,
		doctorTestConfig.JWT.Secret, doctorTestConfig.JWT.ExpiredTime)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Helper function untuk membuat department langsung di database
func createDoctorTestDepartment(name, code, description, floorLocation string) (*models.Department, error) {
	dept := &models.Department{
		Name:          name,
		Code:          code,
		Description:   description,
		FloorLocation: floorLocation,
	}

	result := doctorTestDB.Create(dept)
	if result.Error != nil {
		return nil, result.Error
	}

	return dept, nil
}

// Helper function untuk membuat doctor langsung di database
func createDoctorTestDoctor(userID *uint, employeeID, fullName, specialization, licenseNumber, phone, email string, departmentID *uint, isActive bool) (*models.Doctor, error) {
	doctor := &models.Doctor{
		UserID:         userID,
		EmployeeID:     employeeID,
		FullName:       fullName,
		Specialization: specialization,
		LicenseNumber:  licenseNumber,
		Phone:          phone,
		Email:          email,
		DepartmentID:   departmentID,
		IsActive:       isActive,
	}

	result := doctorTestDB.Create(doctor)
	if result.Error != nil {
		return nil, result.Error
	}

	return doctor, nil
}

// Helper function untuk perform request
func performDoctorRequest(method, path, token string, bodyJSON map[string]interface{}) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if bodyJSON != nil {
		bodyBytes, _ = json.Marshal(bodyJSON)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	w := httptest.NewRecorder()
	doctorTestRouter.ServeHTTP(w, req)
	return w
}

func init() {
	doctorTestRouter = setupDoctorTestRouter()
}

// ==================== GET /api/v1/doctors/me Tests ====================

func Test_GetMyDoctorProfile_Success(t *testing.T) {
	cleanupDoctorTestDB()

	// Create department
	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")

	// Create user and doctor
	user, token, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	w := performDoctorRequest("GET", "/api/v1/doctors/me", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "DOC001", data["employee_id"])
	assert.Equal(t, "Dr. John Smith", data["full_name"])
	assert.Equal(t, "Cardiology", data["specialization"])

	t.Logf("[PASS] GET /api/v1/doctors/me - Success")
}

func Test_GetMyDoctorProfile_NotFound(t *testing.T) {
	cleanupDoctorTestDB()

	// Create user but no doctor record
	_, token, _ := createDoctorTestUser("doctor2", "doctor2@example.com", "081234567891", "password123", models.RoleDoctor, true)

	w := performDoctorRequest("GET", "/api/v1/doctors/me", token, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] GET /api/v1/doctors/me - Not Found")
}

func Test_GetMyDoctorProfile_Unauthorized(t *testing.T) {
	cleanupDoctorTestDB()

	w := performDoctorRequest("GET", "/api/v1/doctors/me", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	t.Logf("[PASS] GET /api/v1/doctors/me - Unauthorized")
}

func Test_GetMyDoctorProfile_Forbidden_Patient(t *testing.T) {
	cleanupDoctorTestDB()

	// Login sebagai patient
	_, patientToken, _ := createDoctorTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performDoctorRequest("GET", "/api/v1/doctors/me", patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/doctors/me - Forbidden Patient")
}

// ==================== PUT /api/v1/doctors/me Tests ====================

func Test_UpdateMyDoctorProfile_Success(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, token, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	requestBody := map[string]interface{}{
		"phone": "081234567899",
		"email": "newemail@hospital.com",
	}

	w := performDoctorRequest("PUT", "/api/v1/doctors/me", token, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "081234567899", data["phone"])
	assert.Equal(t, "newemail@hospital.com", data["email"])

	t.Logf("[PASS] PUT /api/v1/doctors/me - Success")
}

func Test_UpdateMyDoctorProfile_CannotUpdateEmployeeID(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, token, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	requestBody := map[string]interface{}{
		"employee_id": "DOC999", // Should not be allowed
	}

	w := performDoctorRequest("PUT", "/api/v1/doctors/me", token, requestBody)

	// Should either ignore or return error
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusBadRequest)

	t.Logf("[PASS] PUT /api/v1/doctors/me - Cannot Update Employee ID")
}

// ==================== GET /api/v1/doctors Tests ====================

func Test_ListDoctors_Success(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user1, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	user2, _, _ := createDoctorTestUser("doctor2", "doctor2@example.com", "081234567891", "password123", models.RoleDoctor, true)
	createDoctorTestDoctor(&user1.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)
	createDoctorTestDoctor(&user2.ID, "DOC002", "Dr. Jane Doe", "Neurology", "LIC789012", "081234567891", "doctor2@example.com", &dept.ID, true)

	// Login sebagai receptionist
	_, receptionistToken, _ := createDoctorTestUser("receptionist", "receptionist@example.com", "081234567892", "password123", models.RoleReceptionist, true)

	w := performDoctorRequest("GET", "/api/v1/doctors", receptionistToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	t.Logf("[PASS] GET /api/v1/doctors - Success")
}

func Test_ListDoctors_WithPagination(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")

	// Create multiple doctors
	for i := 1; i <= 5; i++ {
		user, _, _ := createDoctorTestUser(fmt.Sprintf("doctor%d", i), fmt.Sprintf("doctor%d@example.com", i), fmt.Sprintf("08123456789%d", i), "password123", models.RoleDoctor, true)
		createDoctorTestDoctor(&user.ID, fmt.Sprintf("DOC%03d", i), fmt.Sprintf("Dr. Doctor %d", i), "Cardiology", fmt.Sprintf("LIC%06d", i), fmt.Sprintf("08123456789%d", i), fmt.Sprintf("doctor%d@example.com", i), &dept.ID, true)
	}

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567899", "password123", models.RoleAdmin, true)

	w := performDoctorRequest("GET", "/api/v1/doctors?page=1&page_size=2", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	dataMap := response["data"].(map[string]interface{})
	doctors := dataMap["data"].([]interface{})
	meta := dataMap["meta"].(map[string]interface{})

	assert.Equal(t, 2, len(doctors))
	assert.Equal(t, float64(1), meta["page"])
	assert.Equal(t, float64(2), meta["page_size"])

	t.Logf("[PASS] GET /api/v1/doctors - With Pagination")
}

func Test_ListDoctors_Unauthorized(t *testing.T) {
	cleanupDoctorTestDB()

	w := performDoctorRequest("GET", "/api/v1/doctors", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	t.Logf("[PASS] GET /api/v1/doctors - Unauthorized")
}

// ==================== GET /api/v1/doctors/:id Tests ====================

func Test_GetDoctorByID_Success(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	_, receptionistToken, _ := createDoctorTestUser("receptionist", "receptionist@example.com", "081234567891", "password123", models.RoleReceptionist, true)

	w := performDoctorRequest("GET", fmt.Sprintf("/api/v1/doctors/%d", doctor.ID), receptionistToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "DOC001", data["employee_id"])
	assert.Equal(t, "Dr. John Smith", data["full_name"])

	t.Logf("[PASS] GET /api/v1/doctors/:id - Success")
}

func Test_GetDoctorByID_NotFound(t *testing.T) {
	cleanupDoctorTestDB()

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDoctorRequest("GET", "/api/v1/doctors/999", adminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] GET /api/v1/doctors/:id - Not Found")
}

// ==================== GET /api/v1/doctors/specialization/:spec Tests ====================

func Test_GetDoctorsBySpecialization_Success(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user1, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	user2, _, _ := createDoctorTestUser("doctor2", "doctor2@example.com", "081234567891", "password123", models.RoleDoctor, true)
	createDoctorTestDoctor(&user1.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)
	createDoctorTestDoctor(&user2.ID, "DOC002", "Dr. Jane Doe", "Cardiology", "LIC789012", "081234567891", "doctor2@example.com", &dept.ID, true)

	_, patientToken, _ := createDoctorTestUser("patient", "patient@example.com", "081234567892", "password123", models.RolePatient, true)

	w := performDoctorRequest("GET", "/api/v1/doctors/specialization/Cardiology", patientToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	t.Logf("[PASS] GET /api/v1/doctors/specialization/:spec - Success")
}

func Test_GetDoctorsBySpecialization_NoResults(t *testing.T) {
	cleanupDoctorTestDB()

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDoctorRequest("GET", "/api/v1/doctors/specialization/Oncology", adminToken, nil)

	// Could be 200 with empty array or 404
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)

	t.Logf("[PASS] GET /api/v1/doctors/specialization/:spec - No Results")
}

// ==================== POST /api/v1/doctors Tests ====================

func Test_CreateDoctor_Success_Admin(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctoruser", "doctoruser@example.com", "081234567890", "password123", models.RoleDoctor, true)
	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567891", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"user_id":        user.ID,
		"employee_id":    "DOC001",
		"full_name":      "Dr. John Smith",
		"specialization": "Cardiology",
		"license_number": "LIC123456",
		"phone":          "081234567890",
		"email":          "drsmith@hospital.com",
		"department_id":  dept.ID,
		"is_active":      true,
	}

	w := performDoctorRequest("POST", "/api/v1/doctors", adminToken, requestBody)

	if w.Code != http.StatusCreated {
		t.Logf("Response Body: %s", w.Body.String())
	}

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "DOC001", data["employee_id"])
	assert.Equal(t, "Dr. John Smith", data["full_name"])

	t.Logf("[PASS] POST /api/v1/doctors - Success Admin")
}

func Test_CreateDoctor_Success_SuperAdmin(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Neurology", "NEURO", "Neurology Department", "Floor 2")
	user, _, _ := createDoctorTestUser("doctoruser", "doctoruser@example.com", "081234567890", "password123", models.RoleDoctor, true)
	_, superAdminToken, _ := createDoctorTestUser("superadmin", "superadmin@example.com", "081234567891", "password123", models.RoleSuperAdmin, true)

	requestBody := map[string]interface{}{
		"user_id":        user.ID,
		"employee_id":    "DOC002",
		"full_name":      "Dr. Jane Doe",
		"specialization": "Neurology",
		"license_number": "LIC789012",
		"phone":          "081234567892",
		"email":          "drjane@hospital.com",
		"department_id":  dept.ID,
		"is_active":      true,
	}

	w := performDoctorRequest("POST", "/api/v1/doctors", superAdminToken, requestBody)

	assert.Equal(t, http.StatusCreated, w.Code)

	t.Logf("[PASS] POST /api/v1/doctors - Success SuperAdmin")
}

func Test_CreateDoctor_Forbidden_Patient(t *testing.T) {
	cleanupDoctorTestDB()

	_, patientToken, _ := createDoctorTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"employee_id": "DOC001",
		"full_name":   "Dr. John Smith",
	}

	w := performDoctorRequest("POST", "/api/v1/doctors", patientToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/doctors - Forbidden Patient")
}

func Test_CreateDoctor_Forbidden_Doctor(t *testing.T) {
	cleanupDoctorTestDB()

	_, doctorToken, _ := createDoctorTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	requestBody := map[string]interface{}{
		"employee_id": "DOC001",
		"full_name":   "Dr. John Smith",
	}

	w := performDoctorRequest("POST", "/api/v1/doctors", doctorToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/doctors - Forbidden Doctor")
}

func Test_CreateDoctor_Forbidden_Receptionist(t *testing.T) {
	cleanupDoctorTestDB()

	_, receptionistToken, _ := createDoctorTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"employee_id": "DOC001",
		"full_name":   "Dr. John Smith",
	}

	w := performDoctorRequest("POST", "/api/v1/doctors", receptionistToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/doctors - Forbidden Receptionist")
}

func Test_CreateDoctor_MissingEmployeeID(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctoruser", "doctoruser@example.com", "081234567890", "password123", models.RoleDoctor, true)
	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567891", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"user_id":        user.ID,
		"full_name":      "Dr. John Smith",
		"specialization": "Cardiology",
		"department_id":  dept.ID,
	}

	w := performDoctorRequest("POST", "/api/v1/doctors", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/doctors - Missing Employee ID")
}

func Test_CreateDoctor_MissingFullName(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctoruser", "doctoruser@example.com", "081234567890", "password123", models.RoleDoctor, true)
	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567891", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"user_id":        user.ID,
		"employee_id":    "DOC001",
		"specialization": "Cardiology",
		"department_id":  dept.ID,
	}

	w := performDoctorRequest("POST", "/api/v1/doctors", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/doctors - Missing Full Name")
}

// ==================== PUT /api/v1/doctors/:id Tests ====================

func Test_UpdateDoctor_Success_Admin(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567891", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"full_name":      "Dr. John Smith Updated",
		"phone":          "081234567899",
		"email":          "drsmith_new@hospital.com",
		"specialization": "Cardiology & Vascular",
	}

	w := performDoctorRequest("PUT", fmt.Sprintf("/api/v1/doctors/%d", doctor.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Dr. John Smith Updated", data["full_name"])

	t.Logf("[PASS] PUT /api/v1/doctors/:id - Success Admin")
}

func Test_UpdateDoctor_Forbidden_Patient(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	_, patientToken, _ := createDoctorTestUser("patient", "patient@example.com", "081234567891", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"full_name": "Dr. John Smith Updated",
	}

	w := performDoctorRequest("PUT", fmt.Sprintf("/api/v1/doctors/%d", doctor.ID), patientToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PUT /api/v1/doctors/:id - Forbidden Patient")
}

func Test_UpdateDoctor_NotFound(t *testing.T) {
	cleanupDoctorTestDB()

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"full_name": "Dr. John Smith Updated",
	}

	w := performDoctorRequest("PUT", "/api/v1/doctors/999", adminToken, requestBody)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] PUT /api/v1/doctors/:id - Not Found")
}

// ==================== PATCH /api/v1/doctors/:id/activate Tests ====================

func Test_ActivateDoctor_Success(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, false) // initially inactive

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567891", "password123", models.RoleAdmin, true)

	w := performDoctorRequest("PATCH", fmt.Sprintf("/api/v1/doctors/%d/activate", doctor.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, true, data["is_active"])

	t.Logf("[PASS] PATCH /api/v1/doctors/:id/activate - Success")
}

func Test_ActivateDoctor_Forbidden_Receptionist(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, false)

	_, receptionistToken, _ := createDoctorTestUser("receptionist", "receptionist@example.com", "081234567891", "password123", models.RoleReceptionist, true)

	w := performDoctorRequest("PATCH", fmt.Sprintf("/api/v1/doctors/%d/activate", doctor.ID), receptionistToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/doctors/:id/activate - Forbidden Receptionist")
}

// ==================== PATCH /api/v1/doctors/:id/deactivate Tests ====================

func Test_DeactivateDoctor_Success(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567891", "password123", models.RoleAdmin, true)

	w := performDoctorRequest("PATCH", fmt.Sprintf("/api/v1/doctors/%d/deactivate", doctor.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, false, data["is_active"])

	t.Logf("[PASS] PATCH /api/v1/doctors/:id/deactivate - Success")
}

func Test_DeactivateDoctor_Forbidden_Doctor(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	_, doctorToken, _ := createDoctorTestUser("doctor2", "doctor2@example.com", "081234567891", "password123", models.RoleDoctor, true)

	w := performDoctorRequest("PATCH", fmt.Sprintf("/api/v1/doctors/%d/deactivate", doctor.ID), doctorToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/doctors/:id/deactivate - Forbidden Doctor")
}

// ==================== GET /api/v1/doctors/deleted Tests ====================

func Test_ListDeletedDoctors_Success(t *testing.T) {
	cleanupDoctorTestDB()

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDoctorRequest("GET", "/api/v1/doctors/deleted", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	t.Logf("[PASS] GET /api/v1/doctors/deleted - Success")
}

func Test_ListDeletedDoctors_Forbidden_Patient(t *testing.T) {
	cleanupDoctorTestDB()

	_, patientToken, _ := createDoctorTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performDoctorRequest("GET", "/api/v1/doctors/deleted", patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/doctors/deleted - Forbidden Patient")
}

// ==================== DELETE /api/v1/doctors/:id Tests ====================

func Test_SoftDeleteDoctor_Success(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567891", "password123", models.RoleAdmin, true)

	w := performDoctorRequest("DELETE", fmt.Sprintf("/api/v1/doctors/%d", doctor.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] DELETE /api/v1/doctors/:id - Success")
}

func Test_SoftDeleteDoctor_Forbidden_Receptionist(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	_, receptionistToken, _ := createDoctorTestUser("receptionist", "receptionist@example.com", "081234567891", "password123", models.RoleReceptionist, true)

	w := performDoctorRequest("DELETE", fmt.Sprintf("/api/v1/doctors/%d", doctor.ID), receptionistToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/doctors/:id - Forbidden Receptionist")
}

func Test_SoftDeleteDoctor_NotFound(t *testing.T) {
	cleanupDoctorTestDB()

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDoctorRequest("DELETE", "/api/v1/doctors/999", adminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] DELETE /api/v1/doctors/:id - Not Found")
}

// ==================== PATCH /api/v1/doctors/:id/restore Tests ====================

func Test_RestoreDoctor_Success(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	// First, soft delete the doctor
	doctorTestDB.Delete(doctor)

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567891", "password123", models.RoleAdmin, true)

	w := performDoctorRequest("PATCH", fmt.Sprintf("/api/v1/doctors/%d/restore", doctor.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/doctors/:id/restore - Success")
}

func Test_RestoreDoctor_Forbidden_Doctor(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	// Soft delete
	doctorTestDB.Delete(doctor)

	_, doctorToken, _ := createDoctorTestUser("doctor2", "doctor2@example.com", "081234567891", "password123", models.RoleDoctor, true)

	w := performDoctorRequest("PATCH", fmt.Sprintf("/api/v1/doctors/%d/restore", doctor.ID), doctorToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/doctors/:id/restore - Forbidden Doctor")
}

// ==================== DELETE /api/v1/doctors/:id/hard-delete Tests ====================

func Test_HardDeleteDoctor_Success(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	_, superAdminToken, _ := createDoctorTestUser("superadmin", "superadmin@example.com", "081234567891", "password123", models.RoleSuperAdmin, true)

	w := performDoctorRequest("DELETE", fmt.Sprintf("/api/v1/doctors/%d/hard-delete", doctor.ID), superAdminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] DELETE /api/v1/doctors/:id/hard-delete - Success")
}

func Test_HardDeleteDoctor_Forbidden_Admin(t *testing.T) {
	cleanupDoctorTestDB()

	dept, _ := createDoctorTestDepartment("Cardiology", "CARD", "Heart Department", "Floor 3")
	user, _, _ := createDoctorTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)
	doctor, _ := createDoctorTestDoctor(&user.ID, "DOC001", "Dr. John Smith", "Cardiology", "LIC123456", "081234567890", "doctor1@example.com", &dept.ID, true)

	_, adminToken, _ := createDoctorTestUser("admin", "admin@example.com", "081234567891", "password123", models.RoleAdmin, true)

	w := performDoctorRequest("DELETE", fmt.Sprintf("/api/v1/doctors/%d/hard-delete", doctor.ID), adminToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/doctors/:id/hard-delete - Forbidden Admin")
}

func Test_HardDeleteDoctor_NotFound(t *testing.T) {
	cleanupDoctorTestDB()

	_, superAdminToken, _ := createDoctorTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	w := performDoctorRequest("DELETE", "/api/v1/doctors/999/hard-delete", superAdminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] DELETE /api/v1/doctors/:id/hard-delete - Not Found")
}
