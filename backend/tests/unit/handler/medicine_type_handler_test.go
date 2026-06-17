package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMedicineTypeHandlerTest() (*gin.Engine, *mocks.MockMedicineTypeService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockMedicineTypeService)
	h := handler.NewMedicineTypeHandler(mockService)

	r := gin.Default()

	r.GET("/medicine-types", h.List)
	r.GET("/medicine-types/active", h.ActiveList)
	r.GET("/medicine-types/inactive", h.InactiveList)
	r.GET("/medicine-types/deleted", h.DeletedList)
	r.GET("/medicine-types/:id", h.FindByID)
	r.POST("/medicine-types", h.Create)
	r.PUT("/medicine-types/:id", h.Update)
	r.DELETE("/medicine-types/:id", h.SoftDelete)
	r.PATCH("/medicine-types/:id/restore", h.Restore)
	r.DELETE("/medicine-types/:id/hard", h.HardDelete)
	r.PATCH("/medicine-types/:id/activate", h.Activate)
	r.PATCH("/medicine-types/:id/deactivate", h.Deactivate)

	return r, mockService
}

func TestMedicineTypeHandler_List_Success(t *testing.T) {
	r, mockService := setupMedicineTypeHandlerTest()

	expectedRes := &dto.MedicineTypeListResponse{
		Data: []dto.MedicineTypeResponse{
			*mocks.NewTestMedicineTypeResponse(mocks.NewTestMedicineTypeWithData(1, "Type1", "T001", "Desc1", false)),
		},
		Meta: dto.MedicineTypePaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.MedicineTypePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicine-types?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_ActiveList_Success(t *testing.T) {
	r, mockService := setupMedicineTypeHandlerTest()

	expectedRes := &dto.MedicineTypeListResponse{
		Data: []dto.MedicineTypeResponse{
			*mocks.NewTestMedicineTypeResponse(mocks.NewTestMedicineTypeWithData(1, "Type1", "T001", "Desc1", false)),
		},
		Meta: dto.MedicineTypePaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("ActiveList", mock.AnythingOfType("*dto.MedicineTypePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicine-types/active?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_InactiveList_Success(t *testing.T) {
	r, mockService := setupMedicineTypeHandlerTest()

	expectedRes := &dto.MedicineTypeListResponse{
		Data: []dto.MedicineTypeResponse{
			*mocks.NewTestMedicineTypeResponse(mocks.NewTestMedicineTypeWithData(1, "Type1", "T001", "Desc1", false)),
		},
		Meta: dto.MedicineTypePaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("InactiveList", mock.AnythingOfType("*dto.MedicineTypePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicine-types/inactive?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_DeletedList_Success(t *testing.T) {
	r, mockService := setupMedicineTypeHandlerTest()

	expectedRes := &dto.MedicineTypeDeletedListResponse{
		Data: []dto.DeletedMedicineTypeResponse{},
		Meta: dto.MedicineTypePaginationMeta{
			TotalItems: 0,
			TotalPages: 0,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("DeletedList", mock.AnythingOfType("*dto.MedicineTypePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicine-types/deleted?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupMedicineTypeHandlerTest()

	expectedRes := mocks.NewTestMedicineTypeResponse(mocks.NewTestMedicineTypeWithData(1, "Type1", "T001", "Desc1", false))

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicine-types/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_Create_Success(t *testing.T) {
	mockService := new(mocks.MockMedicineTypeService)
	createReq := mocks.NewCreateMedicineTypeRequest("Type1", "T001", "Desc1", true)
	expectedRes := mocks.NewTestMedicineTypeResponse(mocks.NewTestMedicineTypeWithData(1, "Type1", "T001", "Desc1", false))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateMedicineTypeRequest")).Return(expectedRes, nil)

	h := handler.NewMedicineTypeHandler(mockService)
	r := gin.Default()
	r.POST("/medicine-types", h.Create)

	req, _ := http.NewRequest(http.MethodPost, "/medicine-types", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_Create_Conflict(t *testing.T) {
	mockService := new(mocks.MockMedicineTypeService)
	createReq := mocks.NewCreateMedicineTypeRequest("Type1", "T001", "Desc1", true)
	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateMedicineTypeRequest")).Return(nil, errors.New("name already exists"))

	h := handler.NewMedicineTypeHandler(mockService)
	r := gin.Default()
	r.POST("/medicine-types", h.Create)

	req, _ := http.NewRequest(http.MethodPost, "/medicine-types", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_Update_Success(t *testing.T) {
	mockService := new(mocks.MockMedicineTypeService)
	updateReq := mocks.NewUpdateMedicineTypeRequest("Type1", "T001", "Desc1", true)
	expectedRes := mocks.NewTestMedicineTypeResponse(mocks.NewTestMedicineTypeWithData(1, "Type1", "T001", "Desc1", false))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateMedicineTypeRequest")).Return(expectedRes, nil)

	h := handler.NewMedicineTypeHandler(mockService)
	r := gin.Default()
	r.PUT("/medicine-types/:id", h.Update)

	req, _ := http.NewRequest(http.MethodPut, "/medicine-types/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_SoftDelete_Success(t *testing.T) {
	r, mockService := setupMedicineTypeHandlerTest()

	mockService.On("SoftDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/medicine-types/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_Restore_Success(t *testing.T) {
	r, mockService := setupMedicineTypeHandlerTest()

	mockService.On("Restore", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/medicine-types/1/restore", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_HardDelete_Success(t *testing.T) {
	r, mockService := setupMedicineTypeHandlerTest()

	mockService.On("HardDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/medicine-types/1/hard", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_Activate_Success(t *testing.T) {
	r, mockService := setupMedicineTypeHandlerTest()

	mockService.On("Activate", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/medicine-types/1/activate", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineTypeHandler_Deactivate_Success(t *testing.T) {
	r, mockService := setupMedicineTypeHandlerTest()

	mockService.On("Deactivate", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/medicine-types/1/deactivate", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
