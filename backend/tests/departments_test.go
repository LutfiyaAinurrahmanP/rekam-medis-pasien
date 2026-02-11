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
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	deptTestDB     *gorm.DB
	deptTestRouter *gin.Engine
	deptTestConfig *config.Config
)

// Setup test router dengan database real
func setupDepartmentTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Initialize config
	deptTestConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-key-for-dept-testing",
			ExpiredTime: 24 * 60 * time.Minute,
		},
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	// Initialize database (SQLite in-memory)
	deptTestDB = database.InitTestDB()

	// Run migrations
	deptTestDB.AutoMigrate(&models.User{}, &models.Department{})

	// Initialize repositories
	userRepo := repository.NewUserRepository(deptTestDB)
	deptRepo := repository.NewDepartmentRepository(deptTestDB)

	// Initialize services
	userService := service.NewUserService(userRepo, deptTestConfig)
	deptService := service.NewDepartmentService(deptRepo, deptTestConfig)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	deptHandler := handler.NewDepartmentHandler(deptService)

	// Setup router
	router := gin.New()

	// Setup routes
	api := router.Group("/api/v1")
	routeConfig := &routes.RouteConfig{
		Config:            deptTestConfig,
		UserHandler:       userHandler,
		DepartmentHandler: deptHandler,
	}
	routes.SetupAPIRouter(api, routeConfig)

	return router
}

// Cleanup database sebelum test
func cleanupDeptTestDB() {
	deptTestDB.Exec("DELETE FROM departments")
	deptTestDB.Exec("DELETE FROM users")
}

// Helper function untuk membuat user langsung di database
func createDeptTestUser(username, email, phone, password, role string, isActive bool) (*models.User, string, error) {
	hashedPassword, _ := utils.HashPassword(password)
	user := &models.User{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: hashedPassword,
		Role:     role,
		IsActive: isActive,
	}

	result := deptTestDB.Create(user)
	if result.Error != nil {
		return nil, "", result.Error
	}

	token, _, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role,
		deptTestConfig.JWT.Secret, deptTestConfig.JWT.ExpiredTime)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Helper function untuk membuat department langsung di database
func createDeptTestDepartment(name, code, description, floorLocation string) (*models.Department, error) {
	dept := &models.Department{
		Name:          name,
		Code:          code,
		Description:   description,
		FloorLocation: floorLocation,
	}

	result := deptTestDB.Create(dept)
	if result.Error != nil {
		return nil, result.Error
	}

	return dept, nil
}

// Helper function untuk perform request
func performDeptRequest(method, path, token string, bodyJSON map[string]interface{}) *httptest.ResponseRecorder {
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
	deptTestRouter.ServeHTTP(w, req)
	return w
}

func init() {
	deptTestRouter = setupDepartmentTestRouter()
}

// ==================== POST /api/v1/departments Tests ====================

func Test_CreateDepartment_Success(t *testing.T) {
	cleanupDeptTestDB()

	// Login sebagai admin
	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Request body
	requestBody := map[string]interface{}{
		"name":           "Kardiologi",
		"code":           "KARDIO",
		"description":    "Departemen jantung dan pembuluh darah",
		"floor_location": "Lantai 3",
	}

	// Request POST /api/v1/departments
	w := performDeptRequest("POST", "/api/v1/departments", adminToken, requestBody)

	// Validasi response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Kardiologi", data["name"])
	assert.Equal(t, "KARDIO", data["code"])
	assert.Equal(t, "Departemen jantung dan pembuluh darah", data["description"])
	assert.Equal(t, "Lantai 3", data["floor_location"])

	t.Logf("[PASS] POST /api/v1/departments - Success")
}

func Test_CreateDepartment_Forbidden_Patient(t *testing.T) {
	cleanupDeptTestDB()

	// Login sebagai patient
	_, patientToken, _ := createDeptTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"name": "Kardiologi",
		"code": "KARDIO",
	}

	w := performDeptRequest("POST", "/api/v1/departments", patientToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - Forbidden Patient")
}

func Test_CreateDepartment_Forbidden_Doctor(t *testing.T) {
	cleanupDeptTestDB()

	// Login sebagai doctor
	_, doctorToken, _ := createDeptTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	requestBody := map[string]interface{}{
		"name": "Kardiologi",
		"code": "KARDIO",
	}

	w := performDeptRequest("POST", "/api/v1/departments", doctorToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - Forbidden Doctor")
}

func Test_CreateDepartment_Forbidden_Receptionist(t *testing.T) {
	cleanupDeptTestDB()

	// Login sebagai receptionist
	_, receptionistToken, _ := createDeptTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"name": "Kardiologi",
		"code": "KARDIO",
	}

	w := performDeptRequest("POST", "/api/v1/departments", receptionistToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - Forbidden Receptionist")
}

func Test_CreateDepartment_Success_SuperAdmin(t *testing.T) {
	cleanupDeptTestDB()

	// Login sebagai super admin
	_, superAdminToken, _ := createDeptTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	requestBody := map[string]interface{}{
		"name":           "Neurologi",
		"code":           "NEURO",
		"description":    "Departemen saraf",
		"floor_location": "Lantai 2",
	}

	w := performDeptRequest("POST", "/api/v1/departments", superAdminToken, requestBody)

	assert.Equal(t, http.StatusCreated, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - Success SuperAdmin")
}

func Test_CreateDepartment_MissingName(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"code": "KARDIO",
	}

	w := performDeptRequest("POST", "/api/v1/departments", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - Missing Name")
}

func Test_CreateDepartment_MissingCode(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"name": "Kardiologi",
	}

	w := performDeptRequest("POST", "/api/v1/departments", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - Missing Code")
}

func Test_CreateDepartment_NameTooLong(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Name lebih dari 100 karakter
	longName := "Kardiologi dan Pembuluh Darah dan Jantung dan Segala Macam Yang Berkaitan Dengan Jantung dan Pembuluh Darah Serta Organ Lainnya"

	requestBody := map[string]interface{}{
		"name": longName,
		"code": "KARDIO",
	}

	w := performDeptRequest("POST", "/api/v1/departments", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - Name Too Long")
}

func Test_CreateDepartment_CodeTooLong(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Code lebih dari 20 karakter
	requestBody := map[string]interface{}{
		"name": "Kardiologi",
		"code": "KARDIO-VERY-LONG-CODE-123456",
	}

	w := performDeptRequest("POST", "/api/v1/departments", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - Code Too Long")
}

func Test_CreateDepartment_FloorLocationTooLong(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// FloorLocation lebih dari 50 karakter
	requestBody := map[string]interface{}{
		"name":           "Kardiologi",
		"code":           "KARDIO",
		"floor_location": "Lantai 3 Gedung A Wing Kiri Koridor Utara Bagian Selatan Nomor 123456789",
	}

	w := performDeptRequest("POST", "/api/v1/departments", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - FloorLocation Too Long")
}

func Test_CreateDepartment_DuplicateCode(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Buat department existing
	createDeptTestDepartment("Kardiologi Existing", "KARDIO", "Dept existing", "Lantai 1")

	// Coba buat department dengan code yang sama
	requestBody := map[string]interface{}{
		"name": "Kardiologi Baru",
		"code": "KARDIO", // Duplicate
	}

	w := performDeptRequest("POST", "/api/v1/departments", adminToken, requestBody)

	assert.Equal(t, http.StatusConflict, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - Duplicate Code")
}

func Test_CreateDepartment_Unauthorized(t *testing.T) {
	cleanupDeptTestDB()

	requestBody := map[string]interface{}{
		"name": "Kardiologi",
		"code": "KARDIO",
	}

	w := performDeptRequest("POST", "/api/v1/departments", "", requestBody)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	t.Logf("[PASS] POST /api/v1/departments - Unauthorized")
}

// ==================== GET /api/v1/departments Tests ====================

func Test_ListDepartments_Success(t *testing.T) {
	cleanupDeptTestDB()

	// Buat beberapa departments
	createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")
	createDeptTestDepartment("Neurologi", "NEURO", "Dept saraf", "Lantai 2")
	createDeptTestDepartment("Pediatri", "PEDIA", "Dept anak", "Lantai 1")

	// Login sebagai patient (all authenticated users can access)
	_, patientToken, _ := createDeptTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performDeptRequest("GET", "/api/v1/departments?page=1&page_size=10", patientToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	// Response structure: { success, message, data: { data: [], meta: {} } }
	dataMap := response["data"].(map[string]interface{})
	departments := dataMap["data"].([]interface{})

	assert.GreaterOrEqual(t, len(departments), 3)

	t.Logf("[PASS] GET /api/v1/departments - Success")
}

func Test_ListDepartments_WithPagination(t *testing.T) {
	cleanupDeptTestDB()

	// Buat 5 departments
	for i := 1; i <= 5; i++ {
		createDeptTestDepartment(fmt.Sprintf("Dept %d", i), fmt.Sprintf("DEPT%d", i), "Description", "Lantai 1")
	}

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDeptRequest("GET", "/api/v1/departments?page=1&page_size=2", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	dataMap := response["data"].(map[string]interface{})
	departments := dataMap["data"].([]interface{})
	meta := dataMap["meta"].(map[string]interface{})

	assert.Equal(t, 2, len(departments))
	assert.Equal(t, float64(1), meta["page"])
	assert.Equal(t, float64(2), meta["page_size"])

	t.Logf("[PASS] GET /api/v1/departments - With Pagination")
}

func Test_ListDepartments_WithSearch(t *testing.T) {
	cleanupDeptTestDB()

	createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")
	createDeptTestDepartment("Neurologi", "NEURO", "Dept saraf", "Lantai 2")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDeptRequest("GET", "/api/v1/departments?search=kardio", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	dataMap := response["data"].(map[string]interface{})
	departments := dataMap["data"].([]interface{})

	// Should find at least 1 department with "kardio"
	assert.GreaterOrEqual(t, len(departments), 1)

	t.Logf("[PASS] GET /api/v1/departments - With Search")
}

func Test_ListDepartments_WithSorting(t *testing.T) {
	cleanupDeptTestDB()

	createDeptTestDepartment("Zebra Dept", "ZEBRA", "Dept Z", "Lantai 1")
	createDeptTestDepartment("Alpha Dept", "ALPHA", "Dept A", "Lantai 2")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDeptRequest("GET", "/api/v1/departments?sort_by=name&sort_dir=asc", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	dataMap := response["data"].(map[string]interface{})
	departments := dataMap["data"].([]interface{})

	if len(departments) >= 2 {
		first := departments[0].(map[string]interface{})
		// First should be "Alpha Dept" when sorted by name ascending
		assert.Equal(t, "Alpha Dept", first["name"])
	}

	t.Logf("[PASS] GET /api/v1/departments - With Sorting")
}

func Test_ListDepartments_Unauthorized(t *testing.T) {
	cleanupDeptTestDB()

	w := performDeptRequest("GET", "/api/v1/departments", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	t.Logf("[PASS] GET /api/v1/departments - Unauthorized")
}

// ==================== GET /api/v1/departments/:id Tests ====================

func Test_GetDepartmentByID_Success(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, patientToken, _ := createDeptTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performDeptRequest("GET", fmt.Sprintf("/api/v1/departments/%d", dept.ID), patientToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Kardiologi", data["name"])
	assert.Equal(t, "KARDIO", data["code"])

	t.Logf("[PASS] GET /api/v1/departments/:id - Success")
}

func Test_GetDepartmentByID_NotFound(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDeptRequest("GET", "/api/v1/departments/99999", adminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] GET /api/v1/departments/:id - Not Found")
}

func Test_GetDepartmentByID_Unauthorized(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	w := performDeptRequest("GET", fmt.Sprintf("/api/v1/departments/%d", dept.ID), "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	t.Logf("[PASS] GET /api/v1/departments/:id - Unauthorized")
}

// ==================== PUT /api/v1/departments/:id Tests ====================

func Test_UpdateDepartment_Success(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"name":           "Kardiologi Updated",
		"code":           "KARDIO-NEW",
		"floor_location": "Lantai 4",
	}

	w := performDeptRequest("PUT", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Kardiologi Updated", data["name"])
	assert.Equal(t, "KARDIO-NEW", data["code"])
	assert.Equal(t, "Lantai 4", data["floor_location"])

	t.Logf("[PASS] PUT /api/v1/departments/:id - Success")
}

func Test_UpdateDepartment_Success_SuperAdmin(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, superAdminToken, _ := createDeptTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	requestBody := map[string]interface{}{
		"name": "Kardiologi by SuperAdmin",
	}

	w := performDeptRequest("PUT", fmt.Sprintf("/api/v1/departments/%d", dept.ID), superAdminToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PUT /api/v1/departments/:id - Success SuperAdmin")
}

func Test_UpdateDepartment_NotFound(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"name": "Updated Name",
	}

	w := performDeptRequest("PUT", "/api/v1/departments/99999", adminToken, requestBody)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] PUT /api/v1/departments/:id - Not Found")
}

func Test_UpdateDepartment_Forbidden_Patient(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, patientToken, _ := createDeptTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"name": "Updated Name",
	}

	w := performDeptRequest("PUT", fmt.Sprintf("/api/v1/departments/%d", dept.ID), patientToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PUT /api/v1/departments/:id - Forbidden Patient")
}

func Test_UpdateDepartment_Forbidden_Doctor(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, doctorToken, _ := createDeptTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	requestBody := map[string]interface{}{
		"name": "Updated Name",
	}

	w := performDeptRequest("PUT", fmt.Sprintf("/api/v1/departments/%d", dept.ID), doctorToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PUT /api/v1/departments/:id - Forbidden Doctor")
}

func Test_UpdateDepartment_DuplicateCode(t *testing.T) {
	cleanupDeptTestDB()

	createDeptTestDepartment("Dept Existing", "EXISTING", "Existing dept", "Lantai 1")
	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"code": "EXISTING", // Duplicate code
	}

	w := performDeptRequest("PUT", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusConflict, w.Code)

	t.Logf("[PASS] PUT /api/v1/departments/:id - Duplicate Code")
}

func Test_UpdateDepartment_NameTooLong(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	longName := "Kardiologi dan Pembuluh Darah dan Jantung dan Segala Macam Yang Berkaitan Dengan Jantung dan Pembuluh Darah Serta Organ Lainnya"

	requestBody := map[string]interface{}{
		"name": longName,
	}

	w := performDeptRequest("PUT", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] PUT /api/v1/departments/:id - Name Too Long")
}

func Test_UpdateDepartment_CodeTooLong(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"code": "KARDIO-VERY-LONG-CODE-123456",
	}

	w := performDeptRequest("PUT", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] PUT /api/v1/departments/:id - Code Too Long")
}

// ==================== DELETE /api/v1/departments/:id (Soft Delete) Tests ====================

func Test_SoftDeleteDepartment_Success(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	t.Logf("[PASS] DELETE /api/v1/departments/:id - Success")
}

func Test_SoftDeleteDepartment_Success_SuperAdmin(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, superAdminToken, _ := createDeptTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	w := performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), superAdminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] DELETE /api/v1/departments/:id - Success SuperAdmin")
}

func Test_SoftDeleteDepartment_NotFound(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDeptRequest("DELETE", "/api/v1/departments/99999", adminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] DELETE /api/v1/departments/:id - Not Found")
}

func Test_SoftDeleteDepartment_Forbidden_Patient(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, patientToken, _ := createDeptTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/departments/:id - Forbidden Patient")
}

func Test_SoftDeleteDepartment_Forbidden_Doctor(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, doctorToken, _ := createDeptTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), doctorToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/departments/:id - Forbidden Doctor")
}

func Test_SoftDeleteDepartment_Forbidden_Receptionist(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, receptionistToken, _ := createDeptTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), receptionistToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/departments/:id - Forbidden Receptionist")
}

// ==================== GET /api/v1/departments/deleted Tests ====================

func Test_ListDeletedDepartments_Success(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Soft delete department first
	performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, nil)

	// Get deleted departments
	w := performDeptRequest("GET", "/api/v1/departments/deleted", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	t.Logf("[PASS] GET /api/v1/departments/deleted - Success")
}

func Test_ListDeletedDepartments_Success_SuperAdmin(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, superAdminToken, _ := createDeptTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	// Soft delete department first
	performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), superAdminToken, nil)

	// Get deleted departments
	w := performDeptRequest("GET", "/api/v1/departments/deleted", superAdminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/departments/deleted - Success SuperAdmin")
}

func Test_ListDeletedDepartments_Forbidden_Patient(t *testing.T) {
	cleanupDeptTestDB()

	_, patientToken, _ := createDeptTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performDeptRequest("GET", "/api/v1/departments/deleted", patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/departments/deleted - Forbidden Patient")
}

func Test_ListDeletedDepartments_Forbidden_Doctor(t *testing.T) {
	cleanupDeptTestDB()

	_, doctorToken, _ := createDeptTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performDeptRequest("GET", "/api/v1/departments/deleted", doctorToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/departments/deleted - Forbidden Doctor")
}

func Test_ListDeletedDepartments_Forbidden_Receptionist(t *testing.T) {
	cleanupDeptTestDB()

	_, receptionistToken, _ := createDeptTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performDeptRequest("GET", "/api/v1/departments/deleted", receptionistToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/departments/deleted - Forbidden Receptionist")
}

// ==================== PATCH /api/v1/departments/:id/restore Tests ====================

func Test_RestoreDepartment_Success(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Soft delete first
	performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, nil)

	// Restore
	w := performDeptRequest("PATCH", fmt.Sprintf("/api/v1/departments/%d/restore", dept.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	t.Logf("[PASS] PATCH /api/v1/departments/:id/restore - Success")
}

func Test_RestoreDepartment_Success_SuperAdmin(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, superAdminToken, _ := createDeptTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	// Soft delete first
	performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), superAdminToken, nil)

	// Restore
	w := performDeptRequest("PATCH", fmt.Sprintf("/api/v1/departments/%d/restore", dept.ID), superAdminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/departments/:id/restore - Success SuperAdmin")
}

func Test_RestoreDepartment_NotFound(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDeptRequest("PATCH", "/api/v1/departments/99999/restore", adminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] PATCH /api/v1/departments/:id/restore - Not Found")
}

func Test_RestoreDepartment_Forbidden_Patient(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)
	_, patientToken, _ := createDeptTestUser("patient", "patient@example.com", "081234567891", "password123", models.RolePatient, true)

	// Soft delete first
	performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, nil)

	// Try restore as patient
	w := performDeptRequest("PATCH", fmt.Sprintf("/api/v1/departments/%d/restore", dept.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/departments/:id/restore - Forbidden Patient")
}

func Test_RestoreDepartment_Forbidden_Doctor(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)
	_, doctorToken, _ := createDeptTestUser("doctor", "doctor@example.com", "081234567891", "password123", models.RoleDoctor, true)

	// Soft delete first
	performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, nil)

	// Try restore as doctor
	w := performDeptRequest("PATCH", fmt.Sprintf("/api/v1/departments/%d/restore", dept.ID), doctorToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/departments/:id/restore - Forbidden Doctor")
}

// ==================== DELETE /api/v1/departments/:id/hard-delete Tests ====================

func Test_HardDeleteDepartment_Success_SuperAdmin(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, superAdminToken, _ := createDeptTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	// Soft delete first
	performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), superAdminToken, nil)

	// Hard delete
	w := performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d/hard-delete", dept.ID), superAdminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	t.Logf("[PASS] DELETE /api/v1/departments/:id/hard-delete - Success SuperAdmin")
}

func Test_HardDeleteDepartment_Forbidden_Admin(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Soft delete first
	performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, nil)

	// Try hard delete as admin (should be forbidden)
	w := performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d/hard-delete", dept.ID), adminToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/departments/:id/hard-delete - Forbidden Admin")
}

func Test_HardDeleteDepartment_Forbidden_Patient(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Dept jantung", "Lantai 3")

	_, superAdminToken, _ := createDeptTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)
	_, patientToken, _ := createDeptTestUser("patient", "patient@example.com", "081234567891", "password123", models.RolePatient, true)

	// Soft delete first
	performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d", dept.ID), superAdminToken, nil)

	// Try hard delete as patient
	w := performDeptRequest("DELETE", fmt.Sprintf("/api/v1/departments/%d/hard-delete", dept.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/departments/:id/hard-delete - Forbidden Patient")
}

func Test_HardDeleteDepartment_NotFound(t *testing.T) {
	cleanupDeptTestDB()

	_, superAdminToken, _ := createDeptTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	w := performDeptRequest("DELETE", "/api/v1/departments/99999/hard-delete", superAdminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] DELETE /api/v1/departments/:id/hard-delete - Not Found")
}

// ==================== Additional Edge Cases ====================

func Test_CreateDepartment_OnlyRequiredFields(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Only name and code (description and floor_location are optional)
	requestBody := map[string]interface{}{
		"name": "Minimal Dept",
		"code": "MINIMAL",
	}

	w := performDeptRequest("POST", "/api/v1/departments", adminToken, requestBody)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Minimal Dept", data["name"])
	assert.Equal(t, "MINIMAL", data["code"])

	t.Logf("[PASS] POST /api/v1/departments - Only Required Fields")
}

func Test_UpdateDepartment_PartialUpdate(t *testing.T) {
	cleanupDeptTestDB()

	dept, _ := createDeptTestDepartment("Kardiologi", "KARDIO", "Original description", "Lantai 3")

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Update only floor_location
	requestBody := map[string]interface{}{
		"floor_location": "Lantai 5",
	}

	w := performDeptRequest("PUT", fmt.Sprintf("/api/v1/departments/%d", dept.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	// Name should remain the same
	assert.Equal(t, "Kardiologi", data["name"])
	// Floor location should be updated
	assert.Equal(t, "Lantai 5", data["floor_location"])

	t.Logf("[PASS] PUT /api/v1/departments/:id - Partial Update")
}

func Test_ListDepartments_EmptyResult(t *testing.T) {
	cleanupDeptTestDB()

	_, patientToken, _ := createDeptTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performDeptRequest("GET", "/api/v1/departments", patientToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	t.Logf("[PASS] GET /api/v1/departments - Empty Result")
}

func Test_GetDepartmentByID_InvalidID(t *testing.T) {
	cleanupDeptTestDB()

	_, adminToken, _ := createDeptTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performDeptRequest("GET", "/api/v1/departments/invalid-id", adminToken, nil)

	// Should return 400 Bad Request or 404 Not Found depending on implementation
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusNotFound)

	t.Logf("[PASS] GET /api/v1/departments/:id - Invalid ID")
}
