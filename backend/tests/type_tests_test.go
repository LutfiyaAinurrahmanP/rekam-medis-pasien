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
	ttTestDB     *gorm.DB
	ttTestRouter *gin.Engine
	ttTestConfig *config.Config
)

func setupTypeTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	ttTestConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-key-for-typetest-testing",
			ExpiredTime: 24 * 60 * time.Minute,
		},
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	ttTestDB = database.InitTestDB()

	ttTestDB.AutoMigrate(&models.User{}, &models.TypeTest{})

	userRepo := repository.NewUserRepository(ttTestDB)
	typeTestRepo := repository.NewTypeTestRepository(ttTestDB)

	userService := service.NewUserService(userRepo, ttTestConfig)
	typeTestService := service.NewTypeTestService(typeTestRepo, ttTestConfig)

	userHandler := handler.NewUserHandler(userService)
	typeTestHandler := handler.NewTypeTestHandler(typeTestService)

	router := gin.New()
	api := router.Group("/api/v1")
	routeConfig := &routes.RouteConfig{
		Config:          ttTestConfig,
		UserHandler:     userHandler,
		TypeTestHandler: typeTestHandler,
	}
	routes.SetupAPIRouter(api, routeConfig)

	return router
}

func cleanupTypeTestDB() {
	ttTestDB.Exec("TRUNCATE TABLE type_tests RESTART IDENTITY CASCADE")
	ttTestDB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
}

func createTTTestUser(username, email, phone, password, role string, isActive bool) (*models.User, string, error) {
	hashedPassword, _ := utils.HashPassword(password)
	user := &models.User{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: hashedPassword,
		Role:     role,
		IsActive: isActive,
	}
	if err := ttTestDB.Create(user).Error; err != nil {
		return nil, "", err
	}
	token, _, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role,
		ttTestConfig.JWT.Secret, ttTestConfig.JWT.ExpiredTime)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

func createTTTestTypeTest(name, code, category, description string, price float64, isActive bool) (*models.TypeTest, error) {
	tt := &models.TypeTest{
		Name:        name,
		Code:        code,
		Category:    category,
		Description: description,
		Price:       price,
		IsActive:    isActive,
	}
	if err := ttTestDB.Create(tt).Error; err != nil {
		return nil, err
	}
	return tt, nil
}

func performTTRequest(method, path, token string, bodyJSON map[string]interface{}) *httptest.ResponseRecorder {
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
	ttTestRouter.ServeHTTP(w, req)
	return w
}

func init() {
	ttTestRouter = setupTypeTestRouter()
}

// ==================== POST /api/v1/test-types ====================

func Test_TypeTest_Create_Success(t *testing.T) {
	cleanupTypeTestDB()

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	body := map[string]interface{}{
		"name":        "Complete Blood Count",
		"code":        "LAB-HEM-001",
		"category":    "Hematologi",
		"description": "Pemeriksaan darah lengkap",
		"price":       150000.0,
		"is_active":   true,
	}

	w := performTTRequest("POST", "/api/v1/test-types", adminToken, body)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "Complete Blood Count", data["name"])
	assert.Equal(t, "LAB-HEM-001", data["code"])
	assert.Equal(t, "Hematologi", data["category"])
	assert.Equal(t, 150000.0, data["price"])

	t.Logf("[PASS] POST /api/v1/test-types - Success")
}

func Test_TypeTest_Create_SuperAdmin_Success(t *testing.T) {
	cleanupTypeTestDB()

	_, superAdminToken, _ := createTTTestUser("superadmin", "superadmin@tt.com", "082222222222", "password123", models.RoleSuperAdmin, true)

	body := map[string]interface{}{
		"name": "Hemoglobin",
		"code": "LAB-HEM-002",
	}

	w := performTTRequest("POST", "/api/v1/test-types", superAdminToken, body)
	assert.Equal(t, http.StatusCreated, w.Code)
	t.Logf("[PASS] POST /api/v1/test-types - SuperAdmin Success")
}

func Test_TypeTest_Create_Forbidden_Patient(t *testing.T) {
	cleanupTypeTestDB()

	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	body := map[string]interface{}{
		"name": "CBC",
		"code": "LAB-001",
	}

	w := performTTRequest("POST", "/api/v1/test-types", patientToken, body)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] POST /api/v1/test-types - Forbidden Patient")
}

func Test_TypeTest_Create_Forbidden_Doctor(t *testing.T) {
	cleanupTypeTestDB()

	_, doctorToken, _ := createTTTestUser("doctor", "doctor@tt.com", "084444444444", "password123", models.RoleDoctor, true)

	body := map[string]interface{}{
		"name": "CBC",
		"code": "LAB-001",
	}

	w := performTTRequest("POST", "/api/v1/test-types", doctorToken, body)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] POST /api/v1/test-types - Forbidden Doctor")
}

func Test_TypeTest_Create_Forbidden_Receptionist(t *testing.T) {
	cleanupTypeTestDB()

	_, recepToken, _ := createTTTestUser("recep", "recep@tt.com", "085555555555", "password123", models.RoleReceptionist, true)

	body := map[string]interface{}{
		"name": "CBC",
		"code": "LAB-001",
	}

	w := performTTRequest("POST", "/api/v1/test-types", recepToken, body)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] POST /api/v1/test-types - Forbidden Receptionist")
}

func Test_TypeTest_Create_Unauthorized(t *testing.T) {
	cleanupTypeTestDB()

	body := map[string]interface{}{
		"name": "CBC",
		"code": "LAB-001",
	}

	w := performTTRequest("POST", "/api/v1/test-types", "", body)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Logf("[PASS] POST /api/v1/test-types - Unauthorized")
}

func Test_TypeTest_Create_MissingName(t *testing.T) {
	cleanupTypeTestDB()

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	body := map[string]interface{}{
		"code": "LAB-001",
	}

	w := performTTRequest("POST", "/api/v1/test-types", adminToken, body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	t.Logf("[PASS] POST /api/v1/test-types - Missing Name")
}

func Test_TypeTest_Create_MissingCode(t *testing.T) {
	cleanupTypeTestDB()

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	body := map[string]interface{}{
		"name": "CBC",
	}

	w := performTTRequest("POST", "/api/v1/test-types", adminToken, body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	t.Logf("[PASS] POST /api/v1/test-types - Missing Code")
}

func Test_TypeTest_Create_DuplicateCode(t *testing.T) {
	cleanupTypeTestDB()

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)
	createTTTestTypeTest("CBC Existing", "LAB-HEM-001", "Hematologi", "", 100000, true)

	body := map[string]interface{}{
		"name": "CBC New",
		"code": "LAB-HEM-001",
	}

	w := performTTRequest("POST", "/api/v1/test-types", adminToken, body)
	assert.Equal(t, http.StatusConflict, w.Code)
	t.Logf("[PASS] POST /api/v1/test-types - Duplicate Code")
}

// ==================== GET /api/v1/test-types ====================

func Test_TypeTest_List_Success(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("CBC", "LAB-HEM-001", "Hematologi", "", 150000, true)
	createTTTestTypeTest("Hemoglobin", "LAB-HEM-002", "Hematologi", "", 50000, true)
	createTTTestTypeTest("Urinalysis", "LAB-URI-001", "Urinalisis", "", 75000, false)

	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("GET", "/api/v1/test-types?page=1&page_size=10", patientToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	data := resp["data"].(map[string]interface{})
	list := data["data"].([]interface{})
	assert.GreaterOrEqual(t, len(list), 3)

	t.Logf("[PASS] GET /api/v1/test-types - Success")
}

func Test_TypeTest_List_WithSearch(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("Complete Blood Count", "LAB-HEM-001", "Hematologi", "", 150000, true)
	createTTTestTypeTest("Urinalysis", "LAB-URI-001", "Urinalisis", "", 75000, true)

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("GET", "/api/v1/test-types?search=blood", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	list := data["data"].([]interface{})
	assert.GreaterOrEqual(t, len(list), 1)

	t.Logf("[PASS] GET /api/v1/test-types - With Search")
}

func Test_TypeTest_List_WithCategoryFilter(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("CBC", "LAB-HEM-001", "Hematologi", "", 150000, true)
	createTTTestTypeTest("Hemoglobin", "LAB-HEM-002", "Hematologi", "", 50000, true)
	createTTTestTypeTest("Urinalysis", "LAB-URI-001", "Urinalisis", "", 75000, true)

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("GET", "/api/v1/test-types?category=Hematologi", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	list := data["data"].([]interface{})
	assert.Equal(t, 2, len(list))

	t.Logf("[PASS] GET /api/v1/test-types - With Category Filter")
}

func Test_TypeTest_List_WithPagination(t *testing.T) {
	cleanupTypeTestDB()

	for i := 1; i <= 5; i++ {
		createTTTestTypeTest(
			fmt.Sprintf("Test Type %d", i),
			fmt.Sprintf("LAB-00%d", i),
			"Hematologi", "", float64(i*10000), true,
		)
	}

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("GET", "/api/v1/test-types?page=1&page_size=2", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	list := data["data"].([]interface{})
	meta := data["meta"].(map[string]interface{})

	assert.Equal(t, 2, len(list))
	assert.Equal(t, float64(1), meta["page"])
	assert.Equal(t, float64(2), meta["page_size"])

	t.Logf("[PASS] GET /api/v1/test-types - With Pagination")
}

func Test_TypeTest_List_Unauthorized(t *testing.T) {
	cleanupTypeTestDB()

	w := performTTRequest("GET", "/api/v1/test-types", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types - Unauthorized")
}

// ==================== GET /api/v1/test-types/active ====================

func Test_TypeTest_ListActive_Success(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("CBC", "LAB-HEM-001", "Hematologi", "", 150000, true)
	createTTTestTypeTest("Inactive Test", "LAB-OLD-001", "Obsolete", "", 0, false)

	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("GET", "/api/v1/test-types/active", patientToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	data := resp["data"].(map[string]interface{})
	list := data["data"].([]interface{})

	// All returned should be active
	for _, item := range list {
		tt := item.(map[string]interface{})
		assert.Equal(t, true, tt["is_active"])
	}

	t.Logf("[PASS] GET /api/v1/test-types/active - Success")
}

func Test_TypeTest_ListActive_Unauthorized(t *testing.T) {
	cleanupTypeTestDB()

	w := performTTRequest("GET", "/api/v1/test-types/active", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/active - Unauthorized")
}

// ==================== GET /api/v1/test-types/inactive ====================

func Test_TypeTest_ListInactive_Success_Doctor(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("Active Test", "LAB-ACT-001", "Hematologi", "", 100000, true)
	createTTTestTypeTest("Inactive Test", "LAB-OLD-001", "Obsolete", "", 0, false)

	_, doctorToken, _ := createTTTestUser("doctor", "doctor@tt.com", "084444444444", "password123", models.RoleDoctor, true)

	w := performTTRequest("GET", "/api/v1/test-types/inactive", doctorToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	list := data["data"].([]interface{})

	// All returned should be inactive
	for _, item := range list {
		tt := item.(map[string]interface{})
		assert.Equal(t, false, tt["is_active"])
	}

	t.Logf("[PASS] GET /api/v1/test-types/inactive - Success Doctor")
}

func Test_TypeTest_ListInactive_Success_Admin(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("Inactive Test", "LAB-OLD-001", "Obsolete", "", 0, false)

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("GET", "/api/v1/test-types/inactive", adminToken, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/inactive - Success Admin")
}

func Test_TypeTest_ListInactive_Forbidden_Patient(t *testing.T) {
	cleanupTypeTestDB()

	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("GET", "/api/v1/test-types/inactive", patientToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/inactive - Forbidden Patient")
}

// ==================== GET /api/v1/test-types/deleted ====================

func Test_TypeTest_ListDeleted_Success_Admin(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("To Delete", "LAB-DEL-001", "Hematologi", "", 100000, true)
	ttTestDB.Delete(&models.TypeTest{}, tt.ID)

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("GET", "/api/v1/test-types/deleted", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	data := resp["data"].(map[string]interface{})
	list := data["data"].([]interface{})
	assert.GreaterOrEqual(t, len(list), 1)

	t.Logf("[PASS] GET /api/v1/test-types/deleted - Success Admin")
}

func Test_TypeTest_ListDeleted_Forbidden_Patient(t *testing.T) {
	cleanupTypeTestDB()

	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("GET", "/api/v1/test-types/deleted", patientToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/deleted - Forbidden Patient")
}

func Test_TypeTest_ListDeleted_Forbidden_Doctor(t *testing.T) {
	cleanupTypeTestDB()

	_, doctorToken, _ := createTTTestUser("doctor", "doctor@tt.com", "084444444444", "password123", models.RoleDoctor, true)

	w := performTTRequest("GET", "/api/v1/test-types/deleted", doctorToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/deleted - Forbidden Doctor")
}

// ==================== GET /api/v1/test-types/:id ====================

func Test_TypeTest_GetByID_Success(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("CBC", "LAB-HEM-001", "Hematologi", "Darah lengkap", 150000, true)

	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("GET", fmt.Sprintf("/api/v1/test-types/%d", tt.ID), patientToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "CBC", data["name"])
	assert.Equal(t, "LAB-HEM-001", data["code"])

	t.Logf("[PASS] GET /api/v1/test-types/:id - Success")
}

func Test_TypeTest_GetByID_NotFound(t *testing.T) {
	cleanupTypeTestDB()

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("GET", "/api/v1/test-types/99999", adminToken, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/:id - Not Found")
}

func Test_TypeTest_GetByID_Unauthorized(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("CBC", "LAB-HEM-001", "Hematologi", "", 150000, true)
	w := performTTRequest("GET", fmt.Sprintf("/api/v1/test-types/%d", tt.ID), "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/:id - Unauthorized")
}

// ==================== GET /api/v1/test-types/code/:code ====================

func Test_TypeTest_GetByCode_Success(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("CBC", "LAB-HEM-001", "Hematologi", "", 150000, true)

	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("GET", "/api/v1/test-types/code/LAB-HEM-001", patientToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "LAB-HEM-001", data["code"])
	assert.Equal(t, "CBC", data["name"])

	t.Logf("[PASS] GET /api/v1/test-types/code/:code - Success")
}

func Test_TypeTest_GetByCode_NotFound(t *testing.T) {
	cleanupTypeTestDB()

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("GET", "/api/v1/test-types/code/NOTEXIST", adminToken, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/code/:code - Not Found")
}

func Test_TypeTest_GetByCode_Unauthorized(t *testing.T) {
	cleanupTypeTestDB()

	w := performTTRequest("GET", "/api/v1/test-types/code/LAB-001", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/code/:code - Unauthorized")
}

// ==================== GET /api/v1/test-types/category/:category ====================

func Test_TypeTest_GetByCategory_Success(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("CBC", "LAB-HEM-001", "Hematologi", "", 150000, true)
	createTTTestTypeTest("Hemoglobin", "LAB-HEM-002", "Hematologi", "", 50000, true)
	createTTTestTypeTest("Urinalysis", "LAB-URI-001", "Urinalisis", "", 75000, true)

	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("GET", "/api/v1/test-types/category/Hematologi", patientToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "Hematologi", data["category"])
	list := data["data"].([]interface{})
	assert.Equal(t, 2, len(list))

	t.Logf("[PASS] GET /api/v1/test-types/category/:category - Success")
}

func Test_TypeTest_GetByCategory_Unauthorized(t *testing.T) {
	cleanupTypeTestDB()

	w := performTTRequest("GET", "/api/v1/test-types/category/Hematologi", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/category/:category - Unauthorized")
}

// ==================== GET /api/v1/test-types/search ====================

func Test_TypeTest_Search_Success(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("Complete Blood Count", "LAB-HEM-001", "Hematologi", "Darah lengkap", 150000, true)
	createTTTestTypeTest("Golongan Darah", "LAB-HEM-004", "Hematologi", "Golongan darah ABO", 75000, true)
	createTTTestTypeTest("Urinalysis", "LAB-URI-001", "Urinalisis", "Urin lengkap", 50000, true)

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("GET", "/api/v1/test-types/search?keyword=darah&category=Hematologi", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	data := resp["data"].(map[string]interface{})
	_, hasCriteria := data["search_criteria"]
	assert.True(t, hasCriteria)
	list := data["data"].([]interface{})
	assert.GreaterOrEqual(t, len(list), 1)

	t.Logf("[PASS] GET /api/v1/test-types/search - Success")
}

func Test_TypeTest_Search_WithPriceFilter(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("Cheap Test", "LAB-CHE-001", "Hematologi", "", 50000, true)
	createTTTestTypeTest("Expensive Test", "LAB-EXP-001", "Hematologi", "", 500000, true)

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("GET", "/api/v1/test-types/search?max_price=100000", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	list := data["data"].([]interface{})

	// Only the cheap test should be returned
	assert.Equal(t, 1, len(list))

	t.Logf("[PASS] GET /api/v1/test-types/search - With Price Filter")
}

func Test_TypeTest_Search_Unauthorized(t *testing.T) {
	cleanupTypeTestDB()

	w := performTTRequest("GET", "/api/v1/test-types/search?keyword=test", "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Logf("[PASS] GET /api/v1/test-types/search - Unauthorized")
}

// ==================== PUT /api/v1/test-types/:id ====================

func Test_TypeTest_Update_Success(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("CBC", "LAB-HEM-001", "Hematologi", "", 150000, true)
	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	body := map[string]interface{}{
		"name":  "CBC Updated",
		"price": 175000.0,
	}

	w := performTTRequest("PUT", fmt.Sprintf("/api/v1/test-types/%d", tt.ID), adminToken, body)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "CBC Updated", data["name"])
	assert.Equal(t, 175000.0, data["price"])

	t.Logf("[PASS] PUT /api/v1/test-types/:id - Success")
}

func Test_TypeTest_Update_NotFound(t *testing.T) {
	cleanupTypeTestDB()

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	body := map[string]interface{}{
		"name": "Updated",
	}

	w := performTTRequest("PUT", "/api/v1/test-types/99999", adminToken, body)
	assert.Equal(t, http.StatusNotFound, w.Code)
	t.Logf("[PASS] PUT /api/v1/test-types/:id - Not Found")
}

func Test_TypeTest_Update_DuplicateCode(t *testing.T) {
	cleanupTypeTestDB()

	createTTTestTypeTest("CBC", "LAB-HEM-001", "Hematologi", "", 150000, true)
	tt2, _ := createTTTestTypeTest("Hemoglobin", "LAB-HEM-002", "Hematologi", "", 50000, true)

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	body := map[string]interface{}{
		"code": "LAB-HEM-001", // duplicate
	}

	w := performTTRequest("PUT", fmt.Sprintf("/api/v1/test-types/%d", tt2.ID), adminToken, body)
	assert.Equal(t, http.StatusConflict, w.Code)
	t.Logf("[PASS] PUT /api/v1/test-types/:id - Duplicate Code")
}

func Test_TypeTest_Update_Forbidden_Patient(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("CBC", "LAB-HEM-001", "Hematologi", "", 150000, true)
	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	body := map[string]interface{}{
		"name": "Updated",
	}

	w := performTTRequest("PUT", fmt.Sprintf("/api/v1/test-types/%d", tt.ID), patientToken, body)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] PUT /api/v1/test-types/:id - Forbidden Patient")
}

// ==================== PATCH /api/v1/test-types/:id/activate ====================

func Test_TypeTest_Activate_Success(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("Inactive Test", "LAB-INA-001", "Hematologi", "", 100000, false)
	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("PATCH", fmt.Sprintf("/api/v1/test-types/%d/activate", tt.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	t.Logf("[PASS] PATCH /api/v1/test-types/:id/activate - Success")
}

func Test_TypeTest_Activate_Forbidden_Patient(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("Test", "LAB-001", "Hematologi", "", 100000, false)
	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("PATCH", fmt.Sprintf("/api/v1/test-types/%d/activate", tt.ID), patientToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] PATCH /api/v1/test-types/:id/activate - Forbidden Patient")
}

// ==================== PATCH /api/v1/test-types/:id/deactivate ====================

func Test_TypeTest_Deactivate_Success(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("Active Test", "LAB-ACT-001", "Hematologi", "", 100000, true)
	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("PATCH", fmt.Sprintf("/api/v1/test-types/%d/deactivate", tt.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	t.Logf("[PASS] PATCH /api/v1/test-types/:id/deactivate - Success")
}

func Test_TypeTest_Deactivate_Forbidden_Doctor(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("Test", "LAB-001", "Hematologi", "", 100000, true)
	_, doctorToken, _ := createTTTestUser("doctor", "doctor@tt.com", "084444444444", "password123", models.RoleDoctor, true)

	w := performTTRequest("PATCH", fmt.Sprintf("/api/v1/test-types/%d/deactivate", tt.ID), doctorToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] PATCH /api/v1/test-types/:id/deactivate - Forbidden Doctor")
}

// ==================== DELETE /api/v1/test-types/:id ====================

func Test_TypeTest_SoftDelete_Success(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("To Delete", "LAB-DEL-001", "Hematologi", "", 100000, true)
	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("DELETE", fmt.Sprintf("/api/v1/test-types/%d", tt.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	// Verify it's no longer accessible
	w2 := performTTRequest("GET", fmt.Sprintf("/api/v1/test-types/%d", tt.ID), adminToken, nil)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	t.Logf("[PASS] DELETE /api/v1/test-types/:id - Success")
}

func Test_TypeTest_SoftDelete_Forbidden_Patient(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("Test", "LAB-001", "Hematologi", "", 100000, true)
	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("DELETE", fmt.Sprintf("/api/v1/test-types/%d", tt.ID), patientToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] DELETE /api/v1/test-types/:id - Forbidden Patient")
}

func Test_TypeTest_SoftDelete_Unauthorized(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("Test", "LAB-001", "Hematologi", "", 100000, true)
	w := performTTRequest("DELETE", fmt.Sprintf("/api/v1/test-types/%d", tt.ID), "", nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	t.Logf("[PASS] DELETE /api/v1/test-types/:id - Unauthorized")
}

// ==================== PATCH /api/v1/test-types/:id/restore ====================

func Test_TypeTest_Restore_Success(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("To Restore", "LAB-RES-001", "Hematologi", "", 100000, true)
	ttTestDB.Delete(&models.TypeTest{}, tt.ID)

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("PATCH", fmt.Sprintf("/api/v1/test-types/%d/restore", tt.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	t.Logf("[PASS] PATCH /api/v1/test-types/:id/restore - Success")
}

func Test_TypeTest_Restore_Forbidden_Patient(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("Test", "LAB-001", "Hematologi", "", 100000, true)
	ttTestDB.Delete(&models.TypeTest{}, tt.ID)

	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("PATCH", fmt.Sprintf("/api/v1/test-types/%d/restore", tt.ID), patientToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] PATCH /api/v1/test-types/:id/restore - Forbidden Patient")
}

// ==================== DELETE /api/v1/test-types/:id/hard-delete ====================

func Test_TypeTest_HardDelete_Success(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("To Hard Delete", "LAB-HD-001", "Hematologi", "", 100000, true)
	ttTestDB.Delete(&models.TypeTest{}, tt.ID) // soft delete first

	_, superAdminToken, _ := createTTTestUser("superadmin", "superadmin@tt.com", "082222222222", "password123", models.RoleSuperAdmin, true)

	w := performTTRequest("DELETE", fmt.Sprintf("/api/v1/test-types/%d/hard-delete", tt.ID), superAdminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, true, resp["success"])

	t.Logf("[PASS] DELETE /api/v1/test-types/:id/hard-delete - Success")
}

func Test_TypeTest_HardDelete_Forbidden_Admin(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("Test", "LAB-001", "Hematologi", "", 100000, true)
	ttTestDB.Delete(&models.TypeTest{}, tt.ID)

	_, adminToken, _ := createTTTestUser("admin", "admin@tt.com", "081111111111", "password123", models.RoleAdmin, true)

	w := performTTRequest("DELETE", fmt.Sprintf("/api/v1/test-types/%d/hard-delete", tt.ID), adminToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] DELETE /api/v1/test-types/:id/hard-delete - Forbidden Admin")
}

func Test_TypeTest_HardDelete_Forbidden_Patient(t *testing.T) {
	cleanupTypeTestDB()

	tt, _ := createTTTestTypeTest("Test", "LAB-001", "Hematologi", "", 100000, true)
	ttTestDB.Delete(&models.TypeTest{}, tt.ID)

	_, patientToken, _ := createTTTestUser("patient", "patient@tt.com", "083333333333", "password123", models.RolePatient, true)

	w := performTTRequest("DELETE", fmt.Sprintf("/api/v1/test-types/%d/hard-delete", tt.ID), patientToken, nil)
	assert.Equal(t, http.StatusForbidden, w.Code)
	t.Logf("[PASS] DELETE /api/v1/test-types/:id/hard-delete - Forbidden Patient")
}
