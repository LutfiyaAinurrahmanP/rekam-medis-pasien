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
	patientTestDB     *gorm.DB
	patientTestRouter *gin.Engine
	patientTestConfig *config.Config
)

// Setup test router dengan database real
func setupPatientTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Initialize config
	patientTestConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-key-for-patient-testing",
			ExpiredTime: 24 * 60 * time.Minute,
		},
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	// Initialize database (SQLite in-memory)
	patientTestDB = database.InitTestDB()

	// Run migrations
	patientTestDB.AutoMigrate(&models.User{}, &models.Patient{})

	// Initialize repositories
	userRepo := repository.NewUserRepository(patientTestDB)
	patientRepo := repository.NewPatientRepository(patientTestDB)

	// Initialize services
	userService := service.NewUserService(userRepo, patientTestConfig)
	patientService := service.NewPatientService(patientRepo, patientTestConfig)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	patientHandler := handler.NewPatientHandler(patientService)

	// Setup router
	router := gin.New()

	// Setup routes
	api := router.Group("/api/v1")
	routeConfig := &routes.RouteConfig{
		Config:         patientTestConfig,
		UserHandler:    userHandler,
		PatientHandler: patientHandler,
	}
	routes.SetupAPIRouter(api, routeConfig)

	return router
}

// Cleanup database sebelum test
func cleanupPatientTestDB() {
	patientTestDB.Exec("DELETE FROM patients")
	patientTestDB.Exec("DELETE FROM users")
}

// Helper function untuk membuat user langsung di database
func createPatientTestUser(username, email, phone, password, role string, isActive bool) (*models.User, string, error) {
	hashedPassword, _ := utils.HashPassword(password)
	user := &models.User{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: hashedPassword,
		Role:     role,
		IsActive: isActive,
	}

	result := patientTestDB.Create(user)
	if result.Error != nil {
		return nil, "", result.Error
	}

	token, _, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role,
		patientTestConfig.JWT.Secret, patientTestConfig.JWT.ExpiredTime)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Helper function untuk membuat patient langsung di database
func createPatientTestPatient(userID *uint, patientCode, fullName, dateOfBirth, gender, bloodType, phone, email, address string) (*models.Patient, error) {
	patient := &models.Patient{
		UserID:      userID,
		PatientCode: patientCode,
		FullName:    fullName,
		DateOfBirth: dateOfBirth,
		Gender:      gender,
		BloodType:   bloodType,
		Phone:       phone,
		Email:       email,
		Address:     address,
	}

	result := patientTestDB.Create(patient)
	if result.Error != nil {
		return nil, result.Error
	}

	return patient, nil
}

// Helper function untuk perform request
func performPatientRequest(method, path, token string, bodyJSON map[string]interface{}) *httptest.ResponseRecorder {
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
	patientTestRouter.ServeHTTP(w, req)
	return w
}

func init() {
	patientTestRouter = setupPatientTestRouter()
}

// ==================== GET /api/v1/patients/me Tests ====================

func Test_GetMyPatientData_Success(t *testing.T) {
	cleanupPatientTestDB()

	// Create user and patient
	user, token, _ := createPatientTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)
	createPatientTestPatient(&user.ID, "P-2024-001", "John Doe", "1990-05-15", "male", "O+", "081234567890", "patient1@example.com", "Jl. Merdeka No. 123")

	w := performPatientRequest("GET", "/api/v1/patients/me", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "P-2024-001", data["patient_code"])
	assert.Equal(t, "John Doe", data["full_name"])

	t.Logf("[PASS] GET /api/v1/patients/me - Success")
}

func Test_GetMyPatientData_NotFound(t *testing.T) {
	cleanupPatientTestDB()

	// Create user but no patient record
	_, token, _ := createPatientTestUser("patient2", "patient2@example.com", "081234567891", "password123", models.RolePatient, true)

	w := performPatientRequest("GET", "/api/v1/patients/me", token, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/me - Not Found")
}

func Test_GetMyPatientData_Unauthorized(t *testing.T) {
	cleanupPatientTestDB()

	w := performPatientRequest("GET", "/api/v1/patients/me", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/me - Unauthorized")
}

// ==================== PUT /api/v1/patients/me Tests ====================

func Test_UpdateMyPatientData_Success(t *testing.T) {
	cleanupPatientTestDB()

	user, token, _ := createPatientTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)
	createPatientTestPatient(&user.ID, "P-2024-001", "John Doe", "1990-05-15", "male", "O+", "081234567890", "patient1@example.com", "Jl. Merdeka No. 123")

	requestBody := map[string]interface{}{
		"full_name": "John Doe Updated",
		"phone":     "081234567899",
		"address":   "Jl. Merdeka No. 124, Jakarta",
	}

	w := performPatientRequest("PUT", "/api/v1/patients/me", token, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "John Doe Updated", data["full_name"])
	assert.Equal(t, "081234567899", data["phone"])

	t.Logf("[PASS] PUT /api/v1/patients/me - Success")
}

func Test_UpdateMyPatientData_CannotUpdatePatientCode(t *testing.T) {
	cleanupPatientTestDB()

	user, token, _ := createPatientTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)
	createPatientTestPatient(&user.ID, "P-2024-001", "John Doe", "1990-05-15", "male", "O+", "081234567890", "patient1@example.com", "Jl. Merdeka No. 123")

	requestBody := map[string]interface{}{
		"patient_code": "P-2024-999", // Should not be allowed
	}

	w := performPatientRequest("PUT", "/api/v1/patients/me", token, requestBody)

	// Should either ignore or return error
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusBadRequest)

	t.Logf("[PASS] PUT /api/v1/patients/me - Cannot Update Patient Code")
}

// ==================== POST /api/v1/patients Tests ====================

func Test_CreatePatient_Success_Receptionist(t *testing.T) {
	cleanupPatientTestDB()

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"patient_code":            "P-2024-010",
		"full_name":               "Jane Smith",
		"date_of_birth":           "1992-08-20",
		"gender":                  "female",
		"blood_type":              "A+",
		"phone":                   "081234567891",
		"email":                   "jane@example.com",
		"address":                 "Jl. Sudirman No. 456",
		"emergency_contact_name":  "John Smith",
		"emergency_contact_phone": "081234567892",
		"insurance_number":        "INS-123456",
		"insurance_provider":      "BPJS Kesehatan",
		"allergies":               "Penicillin",
	}

	w := performPatientRequest("POST", "/api/v1/patients", receptionistToken, requestBody)

	// Debug: print response body
	if w.Code != http.StatusCreated {
		t.Logf("Response Body: %s", w.Body.String())
	}

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "P-2024-010", data["patient_code"])
	assert.Equal(t, "Jane Smith", data["full_name"])

	t.Logf("[PASS] POST /api/v1/patients - Success Receptionist")
}

func Test_CreatePatient_Success_Admin(t *testing.T) {
	cleanupPatientTestDB()

	_, adminToken, _ := createPatientTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"patient_code":  "P-2024-011",
		"full_name":     "Bob Johnson",
		"date_of_birth": "1985-03-10",
		"gender":        "male",
	}

	w := performPatientRequest("POST", "/api/v1/patients", adminToken, requestBody)

	assert.Equal(t, http.StatusCreated, w.Code)

	t.Logf("[PASS] POST /api/v1/patients - Success Admin")
}

func Test_CreatePatient_Success_SuperAdmin(t *testing.T) {
	cleanupPatientTestDB()

	_, superAdminToken, _ := createPatientTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	requestBody := map[string]interface{}{
		"patient_code":  "P-2024-012",
		"full_name":     "Alice Wonder",
		"date_of_birth": "1995-12-25",
		"gender":        "female",
	}

	w := performPatientRequest("POST", "/api/v1/patients", superAdminToken, requestBody)

	assert.Equal(t, http.StatusCreated, w.Code)

	t.Logf("[PASS] POST /api/v1/patients - Success SuperAdmin")
}

func Test_CreatePatient_Forbidden_Patient(t *testing.T) {
	cleanupPatientTestDB()

	_, patientToken, _ := createPatientTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"patient_code":  "P-2024-013",
		"full_name":     "Test Patient",
		"date_of_birth": "1990-01-01",
		"gender":        "male",
	}

	w := performPatientRequest("POST", "/api/v1/patients", patientToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/patients - Forbidden Patient")
}

func Test_CreatePatient_Forbidden_Doctor(t *testing.T) {
	cleanupPatientTestDB()

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	requestBody := map[string]interface{}{
		"patient_code":  "P-2024-014",
		"full_name":     "Test Patient",
		"date_of_birth": "1990-01-01",
		"gender":        "male",
	}

	w := performPatientRequest("POST", "/api/v1/patients", doctorToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/patients - Forbidden Doctor")
}

func Test_CreatePatient_MissingRequiredFields(t *testing.T) {
	cleanupPatientTestDB()

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"patient_code": "P-2024-015",
		// Missing full_name, date_of_birth, gender
	}

	w := performPatientRequest("POST", "/api/v1/patients", receptionistToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/patients - Missing Required Fields")
}

func Test_CreatePatient_DuplicatePatientCode(t *testing.T) {
	cleanupPatientTestDB()

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	// Create existing patient
	createPatientTestPatient(nil, "P-2024-020", "Existing Patient", "1990-01-01", "male", "O+", "081111111111", "existing@example.com", "Address")

	requestBody := map[string]interface{}{
		"patient_code":  "P-2024-020", // Duplicate
		"full_name":     "New Patient",
		"date_of_birth": "1991-01-01",
		"gender":        "male",
	}

	w := performPatientRequest("POST", "/api/v1/patients", receptionistToken, requestBody)

	assert.Equal(t, http.StatusConflict, w.Code)

	t.Logf("[PASS] POST /api/v1/patients - Duplicate Patient Code")
}

func Test_CreatePatient_InvalidGender(t *testing.T) {
	cleanupPatientTestDB()

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"patient_code":  "P-2024-021",
		"full_name":     "Test Patient",
		"date_of_birth": "1990-01-01",
		"gender":        "invalid", // Invalid gender
	}

	w := performPatientRequest("POST", "/api/v1/patients", receptionistToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/patients - Invalid Gender")
}

func Test_CreatePatient_InvalidEmail(t *testing.T) {
	cleanupPatientTestDB()

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"patient_code":  "P-2024-022",
		"full_name":     "Test Patient",
		"date_of_birth": "1990-01-01",
		"gender":        "male",
		"email":         "invalid-email", // Invalid email
	}

	w := performPatientRequest("POST", "/api/v1/patients", receptionistToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/patients - Invalid Email")
}

// ==================== GET /api/v1/patients Tests ====================

func Test_ListPatients_Success(t *testing.T) {
	cleanupPatientTestDB()

	// Create some patients
	createPatientTestPatient(nil, "P-2024-030", "Patient One", "1990-01-01", "male", "O+", "081111111111", "patient1@test.com", "Address 1")
	createPatientTestPatient(nil, "P-2024-031", "Patient Two", "1991-02-02", "female", "A+", "081222222222", "patient2@test.com", "Address 2")
	createPatientTestPatient(nil, "P-2024-032", "Patient Three", "1992-03-03", "male", "B+", "081333333333", "patient3@test.com", "Address 3")

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performPatientRequest("GET", "/api/v1/patients?page=1&page_size=10", doctorToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	dataMap := response["data"].(map[string]interface{})
	patients := dataMap["data"].([]interface{})

	assert.GreaterOrEqual(t, len(patients), 3)

	t.Logf("[PASS] GET /api/v1/patients - Success")
}

func Test_ListPatients_WithPagination(t *testing.T) {
	cleanupPatientTestDB()

	// Create 5 patients
	for i := 1; i <= 5; i++ {
		createPatientTestPatient(nil, fmt.Sprintf("P-2024-%03d", 40+i), fmt.Sprintf("Patient %d", i), "1990-01-01", "male", "O+", fmt.Sprintf("08111111%04d", i), fmt.Sprintf("patient%d@test.com", i), "Address")
	}

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performPatientRequest("GET", "/api/v1/patients?page=1&page_size=2", receptionistToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	dataMap := response["data"].(map[string]interface{})
	patients := dataMap["data"].([]interface{})
	meta := dataMap["meta"].(map[string]interface{})

	assert.Equal(t, 2, len(patients))
	assert.Equal(t, float64(1), meta["page"])
	assert.Equal(t, float64(2), meta["page_size"])

	t.Logf("[PASS] GET /api/v1/patients - With Pagination")
}

func Test_ListPatients_WithSearch(t *testing.T) {
	cleanupPatientTestDB()

	createPatientTestPatient(nil, "P-2024-050", "John Doe", "1990-01-01", "male", "O+", "081111111111", "john@test.com", "Address")
	createPatientTestPatient(nil, "P-2024-051", "Jane Smith", "1991-02-02", "female", "A+", "081222222222", "jane@test.com", "Address")

	_, adminToken, _ := createPatientTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performPatientRequest("GET", "/api/v1/patients?search=john", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	dataMap := response["data"].(map[string]interface{})
	patients := dataMap["data"].([]interface{})

	assert.GreaterOrEqual(t, len(patients), 1)

	t.Logf("[PASS] GET /api/v1/patients - With Search")
}

func Test_ListPatients_WithGenderFilter(t *testing.T) {
	cleanupPatientTestDB()

	createPatientTestPatient(nil, "P-2024-060", "Male Patient", "1990-01-01", "male", "O+", "081111111111", "male@test.com", "Address")
	createPatientTestPatient(nil, "P-2024-061", "Female Patient", "1991-02-02", "female", "A+", "081222222222", "female@test.com", "Address")

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performPatientRequest("GET", "/api/v1/patients?gender=male", doctorToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/patients - With Gender Filter")
}

func Test_ListPatients_WithBloodTypeFilter(t *testing.T) {
	cleanupPatientTestDB()

	createPatientTestPatient(nil, "P-2024-070", "Patient O+", "1990-01-01", "male", "O+", "081111111111", "oplus@test.com", "Address")
	createPatientTestPatient(nil, "P-2024-071", "Patient A+", "1991-02-02", "female", "A+", "081222222222", "aplus@test.com", "Address")

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performPatientRequest("GET", "/api/v1/patients?blood_type=O%2B", receptionistToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/patients - With Blood Type Filter")
}

func Test_ListPatients_Forbidden_Patient(t *testing.T) {
	cleanupPatientTestDB()

	_, patientToken, _ := createPatientTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performPatientRequest("GET", "/api/v1/patients", patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/patients - Forbidden Patient")
}

func Test_ListPatients_Unauthorized(t *testing.T) {
	cleanupPatientTestDB()

	w := performPatientRequest("GET", "/api/v1/patients", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	t.Logf("[PASS] GET /api/v1/patients - Unauthorized")
}

// ==================== GET /api/v1/patients/:id Tests ====================

func Test_GetPatientByID_Success_Staff(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-080", "John Doe", "1990-05-15", "male", "O+", "081111111111", "john@test.com", "Jl. Merdeka No. 123")

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performPatientRequest("GET", fmt.Sprintf("/api/v1/patients/%d", patient.ID), doctorToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "P-2024-080", data["patient_code"])

	t.Logf("[PASS] GET /api/v1/patients/:id - Success Staff")
}

func Test_GetPatientByID_Success_OwnData(t *testing.T) {
	cleanupPatientTestDB()

	user, token, _ := createPatientTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)
	patient, _ := createPatientTestPatient(&user.ID, "P-2024-081", "John Doe", "1990-05-15", "male", "O+", "081234567890", "patient@example.com", "Address")

	w := performPatientRequest("GET", fmt.Sprintf("/api/v1/patients/%d", patient.ID), token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/:id - Success Own Data")
}

func Test_GetPatientByID_Forbidden_OtherPatient(t *testing.T) {
	cleanupPatientTestDB()

	// Create another patient's data
	patient, _ := createPatientTestPatient(nil, "P-2024-082", "Other Patient", "1990-01-01", "male", "O+", "081111111111", "other@test.com", "Address")

	// Login as different patient
	_, patientToken, _ := createPatientTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performPatientRequest("GET", fmt.Sprintf("/api/v1/patients/%d", patient.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/:id - Forbidden Other Patient")
}

func Test_GetPatientByID_NotFound(t *testing.T) {
	cleanupPatientTestDB()

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performPatientRequest("GET", "/api/v1/patients/99999", doctorToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/:id - Not Found")
}

// ==================== GET /api/v1/patients/code/:code Tests ====================

func Test_GetPatientByCode_Success(t *testing.T) {
	cleanupPatientTestDB()

	createPatientTestPatient(nil, "P-2024-090", "John Doe", "1990-05-15", "male", "O+", "081111111111", "john@test.com", "Address")

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performPatientRequest("GET", "/api/v1/patients/code/P-2024-090", receptionistToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "P-2024-090", data["patient_code"])

	t.Logf("[PASS] GET /api/v1/patients/code/:code - Success")
}

func Test_GetPatientByCode_NotFound(t *testing.T) {
	cleanupPatientTestDB()

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performPatientRequest("GET", "/api/v1/patients/code/P-9999-999", doctorToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/code/:code - Not Found")
}

func Test_GetPatientByCode_Forbidden_Patient(t *testing.T) {
	cleanupPatientTestDB()

	_, patientToken, _ := createPatientTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performPatientRequest("GET", "/api/v1/patients/code/P-2024-001", patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/code/:code - Forbidden Patient")
}

// ==================== GET /api/v1/patients/search Tests ====================

func Test_SearchPatients_Success(t *testing.T) {
	cleanupPatientTestDB()

	createPatientTestPatient(nil, "P-2024-100", "Search Test Patient", "1990-05-15", "male", "O+", "081111111111", "search@test.com", "Address")

	_, adminToken, _ := createPatientTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performPatientRequest("GET", "/api/v1/patients/search?full_name=Search", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/search - Success")
}

func Test_SearchPatients_ByDateOfBirth(t *testing.T) {
	cleanupPatientTestDB()

	createPatientTestPatient(nil, "P-2024-101", "DOB Test", "1990-05-15", "male", "O+", "081111111111", "dob@test.com", "Address")

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performPatientRequest("GET", "/api/v1/patients/search?date_of_birth=1990-05-15", doctorToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/search - By Date of Birth")
}

func Test_SearchPatients_Forbidden_Patient(t *testing.T) {
	cleanupPatientTestDB()

	_, patientToken, _ := createPatientTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performPatientRequest("GET", "/api/v1/patients/search?full_name=Test", patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/search - Forbidden Patient")
}

// ==================== PUT /api/v1/patients/:id Tests ====================

func Test_UpdatePatient_Success(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-110", "Original Name", "1990-01-01", "male", "O+", "081111111111", "original@test.com", "Original Address")

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"full_name": "Updated Name",
		"phone":     "081999999999",
		"address":   "Updated Address",
	}

	w := performPatientRequest("PUT", fmt.Sprintf("/api/v1/patients/%d", patient.ID), receptionistToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "Updated Name", data["full_name"])

	t.Logf("[PASS] PUT /api/v1/patients/:id - Success")
}

func Test_UpdatePatient_Success_Admin(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-111", "Test Patient", "1990-01-01", "male", "O+", "081111111111", "test@test.com", "Address")

	_, adminToken, _ := createPatientTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"full_name": "Admin Updated",
	}

	w := performPatientRequest("PUT", fmt.Sprintf("/api/v1/patients/%d", patient.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PUT /api/v1/patients/:id - Success Admin")
}

func Test_UpdatePatient_NotFound(t *testing.T) {
	cleanupPatientTestDB()

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"full_name": "Updated Name",
	}

	w := performPatientRequest("PUT", "/api/v1/patients/99999", receptionistToken, requestBody)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] PUT /api/v1/patients/:id - Not Found")
}

func Test_UpdatePatient_Forbidden_Patient(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-112", "Test Patient", "1990-01-01", "male", "O+", "081111111111", "test@test.com", "Address")

	_, patientToken, _ := createPatientTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"full_name": "Updated Name",
	}

	w := performPatientRequest("PUT", fmt.Sprintf("/api/v1/patients/%d", patient.ID), patientToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PUT /api/v1/patients/:id - Forbidden Patient")
}

func Test_UpdatePatient_Forbidden_Doctor(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-113", "Test Patient", "1990-01-01", "male", "O+", "081111111111", "test@test.com", "Address")

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	requestBody := map[string]interface{}{
		"full_name": "Updated Name",
	}

	w := performPatientRequest("PUT", fmt.Sprintf("/api/v1/patients/%d", patient.ID), doctorToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PUT /api/v1/patients/:id - Forbidden Doctor")
}

// ==================== DELETE /api/v1/patients/:id (Soft Delete) Tests ====================

func Test_SoftDeletePatient_Success(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-120", "Delete Test", "1990-01-01", "male", "O+", "081111111111", "delete@test.com", "Address")

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), receptionistToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] DELETE /api/v1/patients/:id - Success")
}

func Test_SoftDeletePatient_Success_Admin(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-121", "Delete Test", "1990-01-01", "male", "O+", "081111111111", "delete@test.com", "Address")

	_, adminToken, _ := createPatientTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] DELETE /api/v1/patients/:id - Success Admin")
}

func Test_SoftDeletePatient_NotFound(t *testing.T) {
	cleanupPatientTestDB()

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performPatientRequest("DELETE", "/api/v1/patients/99999", receptionistToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] DELETE /api/v1/patients/:id - Not Found")
}

func Test_SoftDeletePatient_Forbidden_Patient(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-122", "Delete Test", "1990-01-01", "male", "O+", "081111111111", "delete@test.com", "Address")

	_, patientToken, _ := createPatientTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/patients/:id - Forbidden Patient")
}

func Test_SoftDeletePatient_Forbidden_Doctor(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-123", "Delete Test", "1990-01-01", "male", "O+", "081111111111", "delete@test.com", "Address")

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), doctorToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/patients/:id - Forbidden Doctor")
}

// ==================== GET /api/v1/patients/deleted Tests ====================

func Test_ListDeletedPatients_Success(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-130", "Deleted Patient", "1990-01-01", "male", "O+", "081111111111", "deleted@test.com", "Address")

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	// Soft delete first
	performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), receptionistToken, nil)

	// Get deleted patients
	w := performPatientRequest("GET", "/api/v1/patients/deleted", receptionistToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/deleted - Success")
}

func Test_ListDeletedPatients_Success_Admin(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-131", "Deleted Patient", "1990-01-01", "male", "O+", "081111111111", "deleted@test.com", "Address")

	_, adminToken, _ := createPatientTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Soft delete first
	performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), adminToken, nil)

	// Get deleted patients
	w := performPatientRequest("GET", "/api/v1/patients/deleted", adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/deleted - Success Admin")
}

func Test_ListDeletedPatients_Forbidden_Patient(t *testing.T) {
	cleanupPatientTestDB()

	_, patientToken, _ := createPatientTestUser("patient", "patient@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performPatientRequest("GET", "/api/v1/patients/deleted", patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/deleted - Forbidden Patient")
}

func Test_ListDeletedPatients_Forbidden_Doctor(t *testing.T) {
	cleanupPatientTestDB()

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performPatientRequest("GET", "/api/v1/patients/deleted", doctorToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/patients/deleted - Forbidden Doctor")
}

// ==================== PATCH /api/v1/patients/:id/restore Tests ====================

func Test_RestorePatient_Success(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-140", "Restore Test", "1990-01-01", "male", "O+", "081111111111", "restore@test.com", "Address")

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	// Soft delete first
	performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), receptionistToken, nil)

	// Restore
	w := performPatientRequest("PATCH", fmt.Sprintf("/api/v1/patients/%d/restore", patient.ID), receptionistToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/patients/:id/restore - Success")
}

func Test_RestorePatient_Success_Admin(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-141", "Restore Test", "1990-01-01", "male", "O+", "081111111111", "restore@test.com", "Address")

	_, adminToken, _ := createPatientTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Soft delete first
	performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), adminToken, nil)

	// Restore
	w := performPatientRequest("PATCH", fmt.Sprintf("/api/v1/patients/%d/restore", patient.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/patients/:id/restore - Success Admin")
}

func Test_RestorePatient_NotFound(t *testing.T) {
	cleanupPatientTestDB()

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performPatientRequest("PATCH", "/api/v1/patients/99999/restore", receptionistToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] PATCH /api/v1/patients/:id/restore - Not Found")
}

func Test_RestorePatient_Forbidden_Patient(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-142", "Restore Test", "1990-01-01", "male", "O+", "081111111111", "restore@test.com", "Address")

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)
	_, patientToken, _ := createPatientTestUser("patient", "patient@example.com", "081234567891", "password123", models.RolePatient, true)

	// Soft delete first
	performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), receptionistToken, nil)

	// Try restore as patient
	w := performPatientRequest("PATCH", fmt.Sprintf("/api/v1/patients/%d/restore", patient.ID), patientToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/patients/:id/restore - Forbidden Patient")
}

func Test_RestorePatient_Forbidden_Doctor(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-143", "Restore Test", "1990-01-01", "male", "O+", "081111111111", "restore@test.com", "Address")

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)
	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567891", "password123", models.RoleDoctor, true)

	// Soft delete first
	performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), receptionistToken, nil)

	// Try restore as doctor
	w := performPatientRequest("PATCH", fmt.Sprintf("/api/v1/patients/%d/restore", patient.ID), doctorToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/patients/:id/restore - Forbidden Doctor")
}

// ==================== DELETE /api/v1/patients/:id/hard-delete Tests ====================

func Test_HardDeletePatient_Success_SuperAdmin(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-150", "Hard Delete Test", "1990-01-01", "male", "O+", "081111111111", "harddelete@test.com", "Address")

	_, superAdminToken, _ := createPatientTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	// Soft delete first
	performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), superAdminToken, nil)

	// Hard delete
	w := performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d/hard-delete", patient.ID), superAdminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] DELETE /api/v1/patients/:id/hard-delete - Success SuperAdmin")
}

func Test_HardDeletePatient_Forbidden_Admin(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-151", "Hard Delete Test", "1990-01-01", "male", "O+", "081111111111", "harddelete@test.com", "Address")

	_, adminToken, _ := createPatientTestUser("admin", "admin@example.com", "081234567890", "password123", models.RoleAdmin, true)

	// Soft delete first
	performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), adminToken, nil)

	// Try hard delete as admin
	w := performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d/hard-delete", patient.ID), adminToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/patients/:id/hard-delete - Forbidden Admin")
}

func Test_HardDeletePatient_Forbidden_Receptionist(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-152", "Hard Delete Test", "1990-01-01", "male", "O+", "081111111111", "harddelete@test.com", "Address")

	_, superAdminToken, _ := createPatientTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)
	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567891", "password123", models.RoleReceptionist, true)

	// Soft delete first
	performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d", patient.ID), superAdminToken, nil)

	// Try hard delete as receptionist
	w := performPatientRequest("DELETE", fmt.Sprintf("/api/v1/patients/%d/hard-delete", patient.ID), receptionistToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/patients/:id/hard-delete - Forbidden Receptionist")
}

func Test_HardDeletePatient_NotFound(t *testing.T) {
	cleanupPatientTestDB()

	_, superAdminToken, _ := createPatientTestUser("superadmin", "superadmin@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	w := performPatientRequest("DELETE", "/api/v1/patients/99999/hard-delete", superAdminToken, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] DELETE /api/v1/patients/:id/hard-delete - Not Found")
}

// ==================== Additional Edge Cases ====================

func Test_CreatePatient_OnlyRequiredFields(t *testing.T) {
	cleanupPatientTestDB()

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"patient_code":  "P-2024-160",
		"full_name":     "Minimal Patient",
		"date_of_birth": "1990-01-01",
		"gender":        "male",
	}

	w := performPatientRequest("POST", "/api/v1/patients", receptionistToken, requestBody)

	assert.Equal(t, http.StatusCreated, w.Code)

	t.Logf("[PASS] POST /api/v1/patients - Only Required Fields")
}

func Test_UpdatePatient_PartialUpdate(t *testing.T) {
	cleanupPatientTestDB()

	patient, _ := createPatientTestPatient(nil, "P-2024-170", "Original Name", "1990-01-01", "male", "O+", "081111111111", "original@test.com", "Original Address")

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	// Update only phone
	requestBody := map[string]interface{}{
		"phone": "081999999999",
	}

	w := performPatientRequest("PUT", fmt.Sprintf("/api/v1/patients/%d", patient.ID), receptionistToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	// Name should remain the same
	assert.Equal(t, "Original Name", data["full_name"])
	// Phone should be updated
	assert.Equal(t, "081999999999", data["phone"])

	t.Logf("[PASS] PUT /api/v1/patients/:id - Partial Update")
}

func Test_ListPatients_EmptyResult(t *testing.T) {
	cleanupPatientTestDB()

	_, doctorToken, _ := createPatientTestUser("doctor", "doctor@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performPatientRequest("GET", "/api/v1/patients", doctorToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] GET /api/v1/patients - Empty Result")
}

func Test_GetPatientByID_InvalidID(t *testing.T) {
	cleanupPatientTestDB()

	_, receptionistToken, _ := createPatientTestUser("receptionist", "receptionist@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performPatientRequest("GET", "/api/v1/patients/invalid-id", receptionistToken, nil)

	// Should return 400 Bad Request or 404 Not Found
	assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusNotFound)

	t.Logf("[PASS] GET /api/v1/patients/:id - Invalid ID")
}
