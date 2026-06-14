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
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/room"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupRoomRouter() (*gin.Engine, *routes.RouteConfig, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	db := SetupTestDatabase()
	cfg := SetupTestConfig()

	roomRepo := repository.NewRoomRepository(db)
	roomService := room.NewRoomService(roomRepo, cfg)
	roomHandler := handler.NewRoomHandler(roomService)

	routeCfg := &routes.RouteConfig{
		Config: cfg,
	}

	r := gin.Default()
	v1 := r.Group("/api/v1")
	routes.SetupRoomsRouter(v1, routeCfg, roomHandler)

	return r, routeCfg, db
}

func createPrerequisitesForRoom(db *gorm.DB) (uint, uint) {
	dept := &models.Department{
		Name: "General",
		Code: "GEN-ROOM",
	}
	db.Create(dept)

	rt := &models.RoomType{
		Name: "VIP",
		Code: "VIP-RM",
	}
	db.Create(rt)

	return dept.ID, rt.ID
}

func createRoomAndGetID(r *gin.Engine, cfg *routes.RouteConfig, db *gorm.DB, num string) uint {
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	deptID, rtID := createPrerequisitesForRoom(db)

	p := 100000.0
	reqBody := dto.CreateRoomRequest{
		RoomNumber:   num,
		RoomTypeID:   &rtID,
		DepartmentID: &deptID,
		BedCapacity:  2,
		PricePerDay:  &p,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewBuffer(body))
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

func TestIntegration_Room_Create(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)
	deptID, rtID := createPrerequisitesForRoom(db)

	reqBody := dto.CreateRoomRequest{
		RoomNumber:   "RM-001",
		RoomTypeID:   &rtID,
		DepartmentID: &deptID,
		BedCapacity:  4,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestIntegration_Room_FindByID(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	id := createRoomAndGetID(r, cfg, db, "RM-100")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/rooms/%d", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Room_List(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	createRoomAndGetID(r, cfg, db, "RM-200")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Room_AvailableList(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	createRoomAndGetID(r, cfg, db, "RM-300")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms/available", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Room_ActiveList(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	createRoomAndGetID(r, cfg, db, "RM-301")
	token := GenerateTestToken(1, models.RolePatient, cfg.Config)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms/active", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Room_Update(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	id := createRoomAndGetID(r, cfg, db, "RM-400")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	reqBody := map[string]string{
		"room_number": "RM-400-UPDATED",
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/rooms/%d", id), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_Room_OccupyAndRelease(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	id := createRoomAndGetID(r, cfg, db, "RM-500") // BedCapacity is 2
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Occupy 1 bed
	reqBodyOccupy := dto.OccupyRoomRequest{Beds: 1}
	bodyO, _ := json.Marshal(reqBodyOccupy)
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/rooms/%d/occupy", id), bytes.NewBuffer(bodyO))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Release 1 bed
	reqBodyRelease := dto.ReleaseRoomRequest{Beds: 1}
	bodyR, _ := json.Marshal(reqBodyRelease)
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/rooms/%d/release", id), bytes.NewBuffer(bodyR))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Room_DeactivateAndActivate(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	id := createRoomAndGetID(r, cfg, db, "RM-600")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Deactivate
	req1, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/rooms/%d/deactivate", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Activate
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/rooms/%d/activate", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Room_DeleteAndRestore(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	id := createRoomAndGetID(r, cfg, db, "RM-700")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/rooms/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Restore
	req2, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/rooms/%d/restore", id), nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Room_DeletedList(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	id := createRoomAndGetID(r, cfg, db, "RM-800")
	token := GenerateTestToken(1, models.RoleAdmin, cfg.Config)

	// Delete
	req1, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/rooms/%d", id), nil)
	req1.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(httptest.NewRecorder(), req1)

	// List Deleted
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms/deleted", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}

func TestIntegration_Room_HardDelete(t *testing.T) {
	r, cfg, db := setupRoomRouter()
	id := createRoomAndGetID(r, cfg, db, "RM-900")
	token := GenerateTestToken(1, models.RoleSuperAdmin, cfg.Config)

	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/rooms/%d/hard-delete", id), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
