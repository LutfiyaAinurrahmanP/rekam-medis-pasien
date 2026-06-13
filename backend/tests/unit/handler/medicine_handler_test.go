package handler

import (
	"bytes"
	"encoding/json"
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

func setupMedicineHandlerTest() (*gin.Engine, *mocks.MockMedicineService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockMedicineService)
	h := handler.NewMedicineHandler(mockService)

	r := gin.Default()

	r.GET("/medicines", h.ListMedicines)
	r.GET("/medicines/active", h.ActiveListMedicines)
	r.GET("/medicines/inactive", h.InactiveListMedicines)
	r.GET("/medicines/deleted", h.DeletedListMedicines)
	r.GET("/medicines/available", h.AvailableListMedicines)
	r.GET("/medicines/low-stock", h.LowStockListMedicines)
	r.GET("/medicines/out-of-stock", h.OutStockListMedicines)
	r.GET("/medicines/:id", h.FindByID)
	r.POST("/medicines", h.CreateMedicine)
	r.PUT("/medicines/:id", h.UpdateMedicine)
	r.POST("/medicines/:id/add-stock", h.AddStock)
	r.POST("/medicines/:id/reduce-stock", h.ReduceStock)
	r.DELETE("/medicines/:id", h.SoftDelete)
	r.PATCH("/medicines/:id/restore", h.Restore)
	r.DELETE("/medicines/:id/hard", h.HardDelete)
	r.PATCH("/medicines/:id/activate", h.Activate)
	r.PATCH("/medicines/:id/deactivate", h.Deactivate)

	return r, mockService
}

func TestMedicineHandler_List_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	expectedRes := &dto.MedicineListResponse{
		Data: []dto.MedicineResponse{
			*mocks.NewTestMedicineResponse(mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)),
		},
		Meta: dto.MedicinePaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.MedicinePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicines?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_AvailableList_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	expectedRes := &dto.MedicineListResponse{
		Data: []dto.MedicineResponse{
			*mocks.NewTestMedicineResponse(mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)),
		},
		Meta: dto.MedicinePaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("AvailableList", mock.AnythingOfType("*dto.MedicinePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicines/available?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_LowStockList_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	expectedRes := &dto.MedicineListResponse{
		Data: []dto.MedicineResponse{
			*mocks.NewTestMedicineResponse(mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 2, 15000, false)),
		},
		Meta: dto.MedicinePaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("LowStockList", mock.AnythingOfType("*dto.MedicinePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicines/low-stock?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_OutStockList_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	expectedRes := &dto.MedicineListResponse{
		Data: []dto.MedicineResponse{
			*mocks.NewTestMedicineResponse(mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 0, 15000, false)),
		},
		Meta: dto.MedicinePaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("OutStockList", mock.AnythingOfType("*dto.MedicinePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicines/out-of-stock?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_ActiveList_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	expectedRes := &dto.MedicineListResponse{
		Data: []dto.MedicineResponse{
			*mocks.NewTestMedicineResponse(mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)),
		},
		Meta: dto.MedicinePaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("ActiveList", mock.AnythingOfType("*dto.MedicinePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicines/active?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_InactiveList_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	expectedRes := &dto.MedicineListResponse{
		Data: []dto.MedicineResponse{
			*mocks.NewTestMedicineResponse(mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false)),
		},
		Meta: dto.MedicinePaginationMeta{
			TotalItems: 1,
			TotalPages: 1,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("InactiveList", mock.AnythingOfType("*dto.MedicinePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicines/inactive?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_DeletedList_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	expectedRes := &dto.MedicineDeletedListResponse{
		Data: []dto.DeletedMedicineResponse{},
		Meta: dto.MedicinePaginationMeta{
			TotalItems: 0,
			TotalPages: 0,
			Page:       1,
			PageSize:   10,
		},
	}

	mockService.On("DeletedList", mock.AnythingOfType("*dto.MedicinePaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicines/deleted?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	expectedRes := mocks.NewTestMedicineResponse(mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false))

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medicines/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_Create_Success(t *testing.T) {
	mockService := new(mocks.MockMedicineService)
	createReq := mocks.NewCreateMedicineRequest("Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, true)
	expectedRes := mocks.NewTestMedicineResponse(mocks.NewTestMedicineWithData(1, "Med1", "Gen1", "Brand1", 1, "500mg", "Manuf1", "Tablet", 100, 15000, false))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateMedicineRequest")).Return(expectedRes, nil)

	h := handler.NewMedicineHandler(mockService)
	r := gin.Default()
	r.POST("/medicines", h.CreateMedicine)

	req, _ := http.NewRequest(http.MethodPost, "/medicines", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_Update_Success(t *testing.T) {
	mockService := new(mocks.MockMedicineService)
	updateReq := mocks.NewUpdateMedicineRequest("Med2", "Gen2", "Brand2", 2, "250mg", "Manuf2", "Capsule", 20000)
	expectedRes := mocks.NewTestMedicineResponse(mocks.NewTestMedicineWithData(1, "Med2", "Gen2", "Brand2", 2, "250mg", "Manuf2", "Capsule", 100, 20000, false))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateMedicineRequest")).Return(expectedRes, nil)

	h := handler.NewMedicineHandler(mockService)
	r := gin.Default()
	r.PUT("/medicines/:id", h.UpdateMedicine)

	req, _ := http.NewRequest(http.MethodPut, "/medicines/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_AddStock_Success(t *testing.T) {
	mockService := new(mocks.MockMedicineService)
	addStockReq := &dto.AddStockRequest{Quantity: 50}

	body, _ := json.Marshal(addStockReq)

	mockService.On("AddStock", uint(1), mock.AnythingOfType("*dto.AddStockRequest")).Return(nil)

	h := handler.NewMedicineHandler(mockService)
	r := gin.Default()
	r.POST("/medicines/:id/add-stock", h.AddStock)

	req, _ := http.NewRequest(http.MethodPost, "/medicines/1/add-stock", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_ReduceStock_Success(t *testing.T) {
	mockService := new(mocks.MockMedicineService)
	reduceStockReq := &dto.ReduceStockRequest{Quantity: 20}

	body, _ := json.Marshal(reduceStockReq)

	mockService.On("ReduceStock", uint(1), mock.AnythingOfType("*dto.ReduceStockRequest")).Return(nil)

	h := handler.NewMedicineHandler(mockService)
	r := gin.Default()
	r.POST("/medicines/:id/reduce-stock", h.ReduceStock)

	req, _ := http.NewRequest(http.MethodPost, "/medicines/1/reduce-stock", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_SoftDelete_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	mockService.On("SoftDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/medicines/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_Restore_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	mockService.On("Restore", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/medicines/1/restore", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_HardDelete_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	mockService.On("HardDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/medicines/1/hard", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_Activate_Success(t *testing.T) {
	r, mockService := setupMedicineHandlerTest()

	mockService.On("Activate", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodPatch, "/medicines/1/activate", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicineHandler_Deactivate_Success(t *testing.T) {
	mockService := new(mocks.MockMedicineService)
	deactReq := &dto.DeactivateMedicineRequest{}
	body, _ := json.Marshal(deactReq)

	mockService.On("Deactivate", uint(1), mock.AnythingOfType("*dto.DeactivateMedicineRequest")).Return(nil)

	h := handler.NewMedicineHandler(mockService)
	r := gin.Default()
	r.PATCH("/medicines/:id/deactivate", h.Deactivate)

	req, _ := http.NewRequest(http.MethodPatch, "/medicines/1/deactivate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
