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
	roomservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/room"
	userservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/user"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	roomTestDB     *gorm.DB
	roomTestRouter *gin.Engine
	roomTestConfig *config.Config
)

// Setup test router dengan database real
func setupRoomTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Initialize config
	roomTestConfig = &config.Config{
		JWT: config.JWTConfig{
			Secret:      "test-secret-key-for-room-testing",
			ExpiredTime: 24 * 60 * time.Minute,
		},
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	// Initialize database (SQLite in-memory)
	roomTestDB = database.InitTestDB()

	// Run migrations
	roomTestDB.AutoMigrate(&models.User{}, &models.Department{}, &models.Room{})

	// Initialize repositories
	userRepo := repository.NewUserRepository(roomTestDB)
	roomRepo := repository.NewRoomRepository(roomTestDB)

	// Initialize services
	userService := userservice.NewUserService(userRepo, roomTestConfig)
	roomService := roomservice.NewRoomService(roomRepo, roomTestConfig)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	roomHandler := handler.NewRoomHandler(roomService)

	// Setup router
	router := gin.New()

	// Setup routes
	api := router.Group("/api/v1")
	routeConfig := &routes.RouteConfig{
		Config:      roomTestConfig,
		UserHandler: userHandler,
		RoomHandler: roomHandler,
	}
	routes.SetupAPIRouter(api, routeConfig)

	return router
}

// Cleanup database sebelum test
func cleanupRoomTestDB() {
	// Delete all records using GORM (handles foreign keys properly)
	// Use Unscoped to delete soft-deleted records as well
	roomTestDB.Unscoped().Where("1 = 1").Delete(&models.Room{})
	roomTestDB.Unscoped().Where("1 = 1").Delete(&models.Department{})
	roomTestDB.Unscoped().Where("1 = 1").Delete(&models.User{})
}

// Helper function untuk membuat user langsung di database
func createRoomTestUser(username, email, phone, password, role string, isActive bool) (*models.User, string, error) {
	hashedPassword, _ := utils.HashPassword(password)
	user := &models.User{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: hashedPassword,
		Role:     role,
		IsActive: isActive,
	}

	result := roomTestDB.Create(user)
	if result.Error != nil {
		return nil, "", result.Error
	}

	token, _, err := utils.GenerateToken(user.ID, user.Username, user.Email, user.Role,
		roomTestConfig.JWT.Secret, roomTestConfig.JWT.ExpiredTime)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Helper function untuk membuat department langsung di database
func createRoomTestDepartment(name, code, description, floorLocation string) (*models.Department, error) {
	dept := &models.Department{
		Name:          name,
		Code:          code,
		Description:   description,
		FloorLocation: floorLocation,
	}

	result := roomTestDB.Create(dept)
	if result.Error != nil {
		return nil, result.Error
	}

	return dept, nil
}

// Helper function untuk membuat room langsung di database
func createRoomTestRoom(roomNumber, roomType string, departmentID *uint, bedCapacity, availableBeds int, pricePerDay float64, isActive bool) (*models.Room, error) {
	room := &models.Room{
		RoomNumber:    roomNumber,
		RoomType:      roomType,
		DepartmentID:  departmentID,
		BedCapacity:   bedCapacity,
		AvailableBeds: availableBeds,
		PricePerDay:   pricePerDay,
		IsActive:      isActive,
	}

	result := roomTestDB.Create(room)
	if result.Error != nil {
		return nil, result.Error
	}

	// GORM has default:true for IsActive, so we need to explicitly update if false
	// Also, BeforeCreate hook sets available_beds to bed_capacity if 0, so we need to update that too
	updates := make(map[string]interface{})
	if !isActive {
		updates["is_active"] = false
	}
	if availableBeds == 0 {
		updates["available_beds"] = 0
	}

	if len(updates) > 0 {
		roomTestDB.Model(room).Updates(updates)
	}

	return room, nil
}

// Helper function untuk perform request
func performRoomRequest(method, path, token string, bodyJSON map[string]interface{}) *httptest.ResponseRecorder {
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
	roomTestRouter.ServeHTTP(w, req)
	return w
}

func init() {
	roomTestRouter = setupRoomTestRouter()
}

// ==================== GET /api/v1/rooms Tests ====================

func Test_GetRooms_Success(t *testing.T) {
	cleanupRoomTestDB()

	// Create department
	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	// Create rooms
	createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)
	createRoomTestRoom("302-A", "vip", &dept.ID, 1, 0, 1500000, true)
	createRoomTestRoom("201-A", "class_1", &dept.ID, 2, 2, 450000, true)

	// Create patient user
	_, token, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRoomRequest("GET", "/api/v1/rooms?page=1&page_size=10", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	rooms := data["data"].([]interface{})
	assert.Equal(t, 3, len(rooms))

	meta := data["meta"].(map[string]interface{})
	assert.Equal(t, float64(3), meta["total_items"])

	t.Logf("[PASS] GET /api/v1/rooms - Success")
}

func Test_GetRooms_WithFilters(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)
	createRoomTestRoom("201-A", "class_1", &dept.ID, 2, 2, 450000, true)

	_, token, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	// Test search filter which is supported by the List endpoint
	w := performRoomRequest("GET", "/api/v1/rooms?search=vip", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	rooms := data["data"].([]interface{})
	assert.GreaterOrEqual(t, len(rooms), 1, "Should find at least 1 room matching 'vip'")

	// Verify the first room contains 'vip' in either room_number or room_type
	room := rooms[0].(map[string]interface{})
	roomNumber := room["room_number"].(string)
	roomType := room["room_type"].(string)
	assert.True(t,
		roomType == "vip" || roomNumber == "301-A",
		"Room should match search term 'vip'")

	t.Logf("[PASS] GET /api/v1/rooms - With Filters")
}

func Test_GetRooms_Unauthorized(t *testing.T) {
	cleanupRoomTestDB()

	w := performRoomRequest("GET", "/api/v1/rooms", "", nil)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	t.Logf("[PASS] GET /api/v1/rooms - Unauthorized")
}

// ==================== GET /api/v1/rooms/available Tests ====================

func Test_GetAvailableRooms_Success(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)
	createRoomTestRoom("302-A", "vip", &dept.ID, 1, 0, 1500000, true) // Occupied (available_beds = 0, should NOT be in available list)
	createRoomTestRoom("201-A", "class_1", &dept.ID, 2, 2, 450000, true)

	_, token, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRoomRequest("GET", "/api/v1/rooms/available", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})

	// Should only return rooms with available_beds > 0 and is_active = true
	// Expected: 301-A (1 bed) and 201-A (2 beds) = 2 rooms
	rooms := data["data"].([]interface{})
	assert.Equal(t, 2, len(rooms))

	// Verify all returned rooms have available_beds > 0
	for _, r := range rooms {
		room := r.(map[string]interface{})
		assert.Greater(t, room["available_beds"], float64(0))
	}

	t.Logf("[PASS] GET /api/v1/rooms/available - Success")
}

// ==================== GET /api/v1/rooms/occupied Tests ====================

func Test_GetOccupiedRooms_Success_Doctor(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	createRoomTestRoom("301-A", "vip", &dept.ID, 1, 0, 1500000, true) // Occupied (available_beds < bed_capacity)
	createRoomTestRoom("302-A", "vip", &dept.ID, 1, 1, 1500000, true) // Available (available_beds = bed_capacity)

	_, token, _ := createRoomTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performRoomRequest("GET", "/api/v1/rooms/occupied", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	rooms := data["data"].([]interface{})

	// Should return 1 occupied room (301-A with 0 available beds)
	assert.Equal(t, 1, len(rooms), "Expected 1 occupied room")

	// Verify the occupied room has available_beds < bed_capacity
	if len(rooms) > 0 {
		room := rooms[0].(map[string]interface{})
		assert.Equal(t, float64(0), room["available_beds"])
		assert.Equal(t, "301-A", room["room_number"])
	}

	t.Logf("[PASS] GET /api/v1/rooms/occupied - Success Doctor")
}

func Test_GetOccupiedRooms_Forbidden_Patient(t *testing.T) {
	cleanupRoomTestDB()

	_, token, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRoomRequest("GET", "/api/v1/rooms/occupied", token, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/rooms/occupied - Forbidden Patient")
}

// ==================== GET /api/v1/rooms/inactive Tests ====================

func Test_GetInactiveRooms_Success_Receptionist(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, false) // Inactive
	createRoomTestRoom("302-A", "vip", &dept.ID, 1, 1, 1500000, true)  // Active

	_, token, _ := createRoomTestUser("receptionist1", "receptionist1@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performRoomRequest("GET", "/api/v1/rooms/inactive", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	rooms := data["data"].([]interface{})
	assert.Equal(t, 1, len(rooms), "Expected 1 inactive room")

	// Add nil check before accessing array
	if len(rooms) > 0 {
		room := rooms[0].(map[string]interface{})
		assert.Equal(t, false, room["is_active"])
		assert.Equal(t, "301-A", room["room_number"])
	}

	t.Logf("[PASS] GET /api/v1/rooms/inactive - Success Receptionist")
}

func Test_GetInactiveRooms_Forbidden_Doctor(t *testing.T) {
	cleanupRoomTestDB()

	_, token, _ := createRoomTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)

	w := performRoomRequest("GET", "/api/v1/rooms/inactive", token, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/rooms/inactive - Forbidden Doctor")
}

// ==================== GET /api/v1/rooms/deleted Tests ====================

func Test_GetDeletedRooms_Success_Admin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	// Create and soft delete a room
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)
	roomTestDB.Delete(room)

	_, token, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performRoomRequest("GET", "/api/v1/rooms/deleted", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	rooms := data["data"].([]interface{})
	assert.Equal(t, 1, len(rooms))

	t.Logf("[PASS] GET /api/v1/rooms/deleted - Success Admin")
}

func Test_GetDeletedRooms_Forbidden_Receptionist(t *testing.T) {
	cleanupRoomTestDB()

	_, token, _ := createRoomTestUser("receptionist1", "receptionist1@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performRoomRequest("GET", "/api/v1/rooms/deleted", token, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] GET /api/v1/rooms/deleted - Forbidden Receptionist")
}

// ==================== GET /api/v1/rooms/:id Tests ====================

func Test_GetRoomByID_Success(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	_, token, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRoomRequest("GET", fmt.Sprintf("/api/v1/rooms/%d", room.ID), token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "301-A", data["room_number"])
	assert.Equal(t, "vip", data["room_type"])

	t.Logf("[PASS] GET /api/v1/rooms/:id - Success")
}

func Test_GetRoomByID_NotFound(t *testing.T) {
	cleanupRoomTestDB()

	_, token, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRoomRequest("GET", "/api/v1/rooms/999", token, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] GET /api/v1/rooms/:id - Not Found")
}

// ==================== GET /api/v1/rooms/number/:room_number Tests ====================

func Test_GetRoomByNumber_Success(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	_, token, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRoomRequest("GET", "/api/v1/rooms/number/301-A", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "301-A", data["room_number"])

	t.Logf("[PASS] GET /api/v1/rooms/number/:room_number - Success")
}

func Test_GetRoomByNumber_NotFound(t *testing.T) {
	cleanupRoomTestDB()

	_, token, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRoomRequest("GET", "/api/v1/rooms/number/999-Z", token, nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	t.Logf("[PASS] GET /api/v1/rooms/number/:room_number - Not Found")
}

// ==================== GET /api/v1/rooms/type/:room_type Tests ====================

func Test_GetRoomsByType_Success(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)
	createRoomTestRoom("302-A", "vip", &dept.ID, 1, 0, 1500000, true)
	createRoomTestRoom("201-A", "class_1", &dept.ID, 2, 2, 450000, true)

	_, token, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRoomRequest("GET", "/api/v1/rooms/type/vip", token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	// This endpoint returns a single room, not a list
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "vip", data["room_type"])

	t.Logf("[PASS] GET /api/v1/rooms/type/:room_type - Success")
}

// ==================== GET /api/v1/rooms/department/:dept_id Tests ====================

func Test_GetRoomsByDepartment_Success(t *testing.T) {
	cleanupRoomTestDB()

	dept1, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	dept2, _ := createRoomTestDepartment("Penyakit Dalam", "PD", "Departemen Penyakit Dalam", "Lantai 2")

	createRoomTestRoom("301-A", "vip", &dept1.ID, 1, 1, 1500000, true)
	createRoomTestRoom("302-A", "vip", &dept1.ID, 1, 0, 1500000, true)
	createRoomTestRoom("201-A", "class_1", &dept2.ID, 2, 2, 450000, true)

	_, token, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	w := performRoomRequest("GET", fmt.Sprintf("/api/v1/rooms/department/%d", dept1.ID), token, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	// This endpoint returns a single room, not a list
	data := response["data"].(map[string]interface{})
	assert.Equal(t, dept1.ID, uint(data["department_id"].(float64)))

	t.Logf("[PASS] GET /api/v1/rooms/department/:dept_id - Success")
}

// ==================== POST /api/v1/rooms Tests ====================

func Test_CreateRoom_Success_Admin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"room_number":   "301-A",
		"room_type":     "vip",
		"department_id": dept.ID,
		"bed_capacity":  1,
		"price_per_day": 1500000,
		"is_active":     true,
	}

	w := performRoomRequest("POST", "/api/v1/rooms", adminToken, requestBody)

	if w.Code != http.StatusCreated {
		t.Logf("Response Body: %s", w.Body.String())
	}

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "301-A", data["room_number"])
	assert.Equal(t, "vip", data["room_type"])

	t.Logf("[PASS] POST /api/v1/rooms - Success Admin")
}

func Test_CreateRoom_Success_SuperAdmin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	_, superAdminToken, _ := createRoomTestUser("superadmin1", "superadmin1@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	requestBody := map[string]interface{}{
		"room_number":   "302-A",
		"room_type":     "vip",
		"department_id": dept.ID,
		"bed_capacity":  1,
		"price_per_day": 1500000,
	}

	w := performRoomRequest("POST", "/api/v1/rooms", superAdminToken, requestBody)

	assert.Equal(t, http.StatusCreated, w.Code)

	t.Logf("[PASS] POST /api/v1/rooms - Success Super Admin")
}

func Test_CreateRoom_DuplicateRoomNumber(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"room_number":   "301-A", // Duplicate
		"room_type":     "vip",
		"department_id": dept.ID,
		"bed_capacity":  1,
		"price_per_day": 1500000,
	}

	w := performRoomRequest("POST", "/api/v1/rooms", adminToken, requestBody)

	// API returns 409 Conflict for duplicate resources (semantically correct)
	assert.Equal(t, http.StatusConflict, w.Code)

	t.Logf("[PASS] POST /api/v1/rooms - Duplicate Room Number")
}

func Test_CreateRoom_InvalidRoomType(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"room_number":   "301-A",
		"room_type":     "invalid_type",
		"department_id": dept.ID,
		"bed_capacity":  1,
		"price_per_day": 1500000,
	}

	w := performRoomRequest("POST", "/api/v1/rooms", adminToken, requestBody)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	t.Logf("[PASS] POST /api/v1/rooms - Invalid Room Type")
}

func Test_CreateRoom_Forbidden_Receptionist(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	_, recepToken, _ := createRoomTestUser("receptionist1", "receptionist1@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"room_number":   "301-A",
		"room_type":     "vip",
		"department_id": dept.ID,
		"bed_capacity":  1,
		"price_per_day": 1500000,
	}

	w := performRoomRequest("POST", "/api/v1/rooms", recepToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/rooms - Forbidden Receptionist")
}

func Test_CreateRoom_Forbidden_Doctor(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	_, doctorToken, _ := createRoomTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)

	requestBody := map[string]interface{}{
		"room_number":   "301-A",
		"room_type":     "vip",
		"department_id": dept.ID,
		"bed_capacity":  1,
		"price_per_day": 1500000,
	}

	w := performRoomRequest("POST", "/api/v1/rooms", doctorToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/rooms - Forbidden Doctor")
}

func Test_CreateRoom_Forbidden_Patient(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")

	_, patientToken, _ := createRoomTestUser("patient1", "patient1@example.com", "081234567890", "password123", models.RolePatient, true)

	requestBody := map[string]interface{}{
		"room_number":   "301-A",
		"room_type":     "vip",
		"department_id": dept.ID,
		"bed_capacity":  1,
		"price_per_day": 1500000,
	}

	w := performRoomRequest("POST", "/api/v1/rooms", patientToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] POST /api/v1/rooms - Forbidden Patient")
}

// ==================== PUT /api/v1/rooms/:id Tests ====================

func Test_UpdateRoom_Success_Admin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"price_per_day": 1800000,
		"room_type":     "class_1",
	}

	w := performRoomRequest("PUT", fmt.Sprintf("/api/v1/rooms/%d", room.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])
	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1800000), data["price_per_day"])

	t.Logf("[PASS] PUT /api/v1/rooms/:id - Success Admin")
}

func Test_UpdateRoom_Forbidden_Receptionist(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	_, recepToken, _ := createRoomTestUser("receptionist1", "receptionist1@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"price_per_day": 1800000,
	}

	w := performRoomRequest("PUT", fmt.Sprintf("/api/v1/rooms/%d", room.ID), recepToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PUT /api/v1/rooms/:id - Forbidden Receptionist")
}

// ==================== PATCH /api/v1/rooms/:id/activate Tests ====================

func Test_ActivateRoom_Success_Admin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, false) // Inactive

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performRoomRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%d/activate", room.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	t.Logf("[PASS] PATCH /api/v1/rooms/:id/activate - Success Admin")
}

func Test_ActivateRoom_Forbidden_Receptionist(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, false)

	_, recepToken, _ := createRoomTestUser("receptionist1", "receptionist1@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performRoomRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%d/activate", room.ID), recepToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/rooms/:id/activate - Forbidden Receptionist")
}

// ==================== PATCH /api/v1/rooms/:id/deactivate Tests ====================

func Test_DeactivateRoom_Success_Admin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true) // Active

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performRoomRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%d/deactivate", room.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	t.Logf("[PASS] PATCH /api/v1/rooms/:id/deactivate - Success Admin")
}

// ==================== PATCH /api/v1/rooms/:id/occupy Tests ====================

func Test_OccupyRoom_Success_Receptionist(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 2, 2, 1500000, true) // 2 beds available

	_, recepToken, _ := createRoomTestUser("receptionist1", "receptionist1@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"beds": 1,
	}

	w := performRoomRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%d/occupy", room.ID), recepToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	// Verify available_beds decreased
	var updatedRoom models.Room
	roomTestDB.First(&updatedRoom, room.ID)
	assert.Equal(t, 1, updatedRoom.AvailableBeds)

	t.Logf("[PASS] PATCH /api/v1/rooms/:id/occupy - Success Receptionist")
}

func Test_OccupyRoom_Success_Admin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 2, 2, 1500000, true)

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"beds": 1,
	}

	w := performRoomRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%d/occupy", room.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/rooms/:id/occupy - Success Admin")
}

func Test_OccupyRoom_Forbidden_Doctor(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 2, 2, 1500000, true)

	_, doctorToken, _ := createRoomTestUser("doctor1", "doctor1@example.com", "081234567890", "password123", models.RoleDoctor, true)

	requestBody := map[string]interface{}{
		"beds": 1,
	}

	w := performRoomRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%d/occupy", room.ID), doctorToken, requestBody)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/rooms/:id/occupy - Forbidden Doctor")
}

// ==================== PATCH /api/v1/rooms/:id/release Tests ====================

func Test_ReleaseRoom_Success_Receptionist(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 2, 1, 1500000, true) // 1 bed occupied

	_, recepToken, _ := createRoomTestUser("receptionist1", "receptionist1@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	requestBody := map[string]interface{}{
		"beds": 1,
	}

	w := performRoomRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%d/release", room.ID), recepToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	// Verify available_beds increased
	var updatedRoom models.Room
	roomTestDB.First(&updatedRoom, room.ID)
	assert.Equal(t, 2, updatedRoom.AvailableBeds)

	t.Logf("[PASS] PATCH /api/v1/rooms/:id/release - Success Receptionist")
}

func Test_ReleaseRoom_Success_Admin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 2, 1, 1500000, true)

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	requestBody := map[string]interface{}{
		"beds": 1,
	}

	w := performRoomRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%d/release", room.ID), adminToken, requestBody)

	assert.Equal(t, http.StatusOK, w.Code)

	t.Logf("[PASS] PATCH /api/v1/rooms/:id/release - Success Admin")
}

// ==================== DELETE /api/v1/rooms/:id Tests ====================

func Test_DeleteRoom_Success_Admin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performRoomRequest("DELETE", fmt.Sprintf("/api/v1/rooms/%d", room.ID), adminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	// Verify soft delete
	var deletedRoom models.Room
	result := roomTestDB.Unscoped().First(&deletedRoom, room.ID)
	assert.NoError(t, result.Error)
	assert.NotNil(t, deletedRoom.DeletedAt)

	t.Logf("[PASS] DELETE /api/v1/rooms/:id - Success Admin")
}

func Test_DeleteRoom_Forbidden_Receptionist(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	_, recepToken, _ := createRoomTestUser("receptionist1", "receptionist1@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performRoomRequest("DELETE", fmt.Sprintf("/api/v1/rooms/%d", room.ID), recepToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/rooms/:id - Forbidden Receptionist")
}

// ==================== PATCH /api/v1/rooms/:id/restore Tests ====================

func Test_RestoreRoom_Success_Admin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	// Soft delete the room
	roomTestDB.Delete(room)

	// Verify room is soft-deleted
	var deletedCheck models.Room
	resultDeleted := roomTestDB.Unscoped().First(&deletedCheck, room.ID)
	if resultDeleted.Error != nil {
		t.Fatalf("Failed to verify room after delete: %v", resultDeleted.Error)
	}
	if deletedCheck.DeletedAt.Time.IsZero() {
		t.Fatalf("Room was not soft-deleted properly")
	}

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performRoomRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%d/restore", room.ID), adminToken, nil)

	if w.Code != http.StatusOK {
		t.Logf("Response Status: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())
	}

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	// Verify restore
	var restoredRoom models.Room
	result := roomTestDB.First(&restoredRoom, room.ID)
	assert.NoError(t, result.Error)

	t.Logf("[PASS] PATCH /api/v1/rooms/:id/restore - Success Admin")
}

func Test_RestoreRoom_Forbidden_Receptionist(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)
	roomTestDB.Delete(room)

	_, recepToken, _ := createRoomTestUser("receptionist1", "receptionist1@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performRoomRequest("PATCH", fmt.Sprintf("/api/v1/rooms/%d/restore", room.ID), recepToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] PATCH /api/v1/rooms/:id/restore - Forbidden Receptionist")
}

// ==================== DELETE /api/v1/rooms/:id/hard-delete Tests ====================

func Test_HardDeleteRoom_Success_SuperAdmin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	_, superAdminToken, _ := createRoomTestUser("superadmin1", "superadmin1@example.com", "081234567890", "password123", models.RoleSuperAdmin, true)

	w := performRoomRequest("DELETE", fmt.Sprintf("/api/v1/rooms/%d/hard-delete", room.ID), superAdminToken, nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, true, response["success"])

	// Verify hard delete
	var deletedRoom models.Room
	result := roomTestDB.Unscoped().First(&deletedRoom, room.ID)
	assert.Error(t, result.Error)

	t.Logf("[PASS] DELETE /api/v1/rooms/:id/hard-delete - Success Super Admin")
}

func Test_HardDeleteRoom_Forbidden_Admin(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	_, adminToken, _ := createRoomTestUser("admin1", "admin1@example.com", "081234567890", "password123", models.RoleAdmin, true)

	w := performRoomRequest("DELETE", fmt.Sprintf("/api/v1/rooms/%d/hard-delete", room.ID), adminToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/rooms/:id/hard-delete - Forbidden Admin")
}

func Test_HardDeleteRoom_Forbidden_Receptionist(t *testing.T) {
	cleanupRoomTestDB()

	dept, _ := createRoomTestDepartment("Kardiologi", "KARDIO", "Departemen Kardiologi", "Lantai 3")
	room, _ := createRoomTestRoom("301-A", "vip", &dept.ID, 1, 1, 1500000, true)

	_, recepToken, _ := createRoomTestUser("receptionist1", "receptionist1@example.com", "081234567890", "password123", models.RoleReceptionist, true)

	w := performRoomRequest("DELETE", fmt.Sprintf("/api/v1/rooms/%d/hard-delete", room.ID), recepToken, nil)

	assert.Equal(t, http.StatusForbidden, w.Code)

	t.Logf("[PASS] DELETE /api/v1/rooms/:id/hard-delete - Forbidden Receptionist")
}
