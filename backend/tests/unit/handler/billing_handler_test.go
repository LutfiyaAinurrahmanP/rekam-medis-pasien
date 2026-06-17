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

func setupBillingHandlerTest() (*gin.Engine, *mocks.MockBillingService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockBillingService)
	h := handler.NewBillingHandler(mockService)

	r := gin.Default()

	r.GET("/billings", h.List)
	r.GET("/billings/deleted", h.DeletedList)
	r.GET("/billings/:id", h.FindByID)
	r.GET("/billings/patient/:patient_id", h.FindByPatientID)
	r.GET("/billings/invoice/:invoice_number", h.FindByInvoiceNumber)
	r.POST("/billings", h.Create)
	r.PUT("/billings/:id", h.Update)
	r.POST("/billings/:id/payment", h.RecordPayment)
	r.PUT("/billings/:id/cancel", h.Cancel)
	r.DELETE("/billings/:id", h.Delete)
	r.PUT("/billings/:id/restore", h.Restore)
	r.DELETE("/billings/:id/hard", h.HardDelete)

	// Items
	r.GET("/billings/:id/items", h.ListItems)
	r.GET("/billings/:id/items/:item_id", h.FindItemByID)
	r.POST("/billings/:id/items", h.CreateItem)
	r.PUT("/billings/:id/items/:item_id", h.UpdateItem)
	r.DELETE("/billings/:id/items/:item_id", h.DeleteItem)

	return r, mockService
}

func TestBillingHandler_List_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	expectedRes := &dto.BillingListResponse{
		Status: "success",
		Data:   []dto.BillingResponse{},
	}
	mockService.On("List", mock.AnythingOfType("dto.BillingPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/billings", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_DeletedList_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	expectedRes := &dto.BillingDeletedListResponse{
		Status: "success",
		Data:   []dto.BillingResponse{},
	}
	mockService.On("DeletedList", mock.AnythingOfType("dto.BillingPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/billings/deleted", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	expectedRes := mocks.NewTestBillingResponse(mocks.NewTestBillingWithData(1, 1))
	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/billings/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_FindByPatientID_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	expectedRes := &dto.BillingListResponse{
		Status: "success",
		Data:   []dto.BillingResponse{},
	}
	mockService.On("FindByPatientID", uint(1), mock.AnythingOfType("dto.BillingPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/billings/patient/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_FindByInvoiceNumber_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	expectedRes := mocks.NewTestBillingResponse(mocks.NewTestBillingWithData(1, 1))
	mockService.On("FindByInvoiceNumber", "INV-1").Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/billings/invoice/INV-1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_Create_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	createReq := mocks.NewCreateBillingRequest(1)
	expectedRes := mocks.NewTestBillingResponse(mocks.NewTestBillingWithData(1, 1))

	mockService.On("Create", mock.AnythingOfType("dto.CreateBillingRequest")).Return(expectedRes, nil)

	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest(http.MethodPost, "/billings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_Update_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	updateReq := mocks.NewUpdateBillingRequest()
	expectedRes := mocks.NewTestBillingResponse(mocks.NewTestBillingWithData(1, 1))

	mockService.On("Update", uint(1), mock.AnythingOfType("dto.UpdateBillingRequest")).Return(expectedRes, nil)

	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest(http.MethodPut, "/billings/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_RecordPayment_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	paymentReq := mocks.NewRecordPaymentRequest()
	expectedRes := mocks.NewTestBillingResponse(mocks.NewTestBillingWithData(1, 1))

	mockService.On("RecordPayment", uint(1), mock.AnythingOfType("dto.RecordPaymentRequest")).Return(expectedRes, nil)

	body, _ := json.Marshal(paymentReq)
	req, _ := http.NewRequest(http.MethodPost, "/billings/1/payment", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_Cancel_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	expectedRes := mocks.NewTestBillingResponse(mocks.NewTestBillingWithData(1, 1))
	mockService.On("Cancel", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/billings/1/cancel", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_Delete_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/billings/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_Restore_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	expectedRes := mocks.NewTestBillingResponse(mocks.NewTestBillingWithData(1, 1))
	mockService.On("Restore", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/billings/1/restore", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_HardDelete_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	mockService.On("HardDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/billings/1/hard", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

// Items

func TestBillingHandler_ListItems_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	expectedRes := []dto.BillingItemResponse{}
	mockService.On("ListItems", uint(1), mock.AnythingOfType("dto.BillingItemPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/billings/1/items", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_FindItemByID_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	expectedRes := mocks.NewTestBillingItemResponse(mocks.NewTestBillingItemWithData(1, 1))
	mockService.On("FindItemByID", uint(1), uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/billings/1/items/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_CreateItem_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	createReq := mocks.NewCreateBillingItemRequest()
	expectedRes := mocks.NewTestBillingItemResponse(mocks.NewTestBillingItemWithData(1, 1))

	mockService.On("CreateItem", uint(1), mock.AnythingOfType("dto.CreateBillingItemRequest")).Return(expectedRes, nil)

	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest(http.MethodPost, "/billings/1/items", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_UpdateItem_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	updateReq := mocks.NewUpdateBillingItemRequest()
	expectedRes := mocks.NewTestBillingItemResponse(mocks.NewTestBillingItemWithData(1, 1))

	mockService.On("UpdateItem", uint(1), uint(1), mock.AnythingOfType("dto.UpdateBillingItemRequest")).Return(expectedRes, nil)

	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest(http.MethodPut, "/billings/1/items/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestBillingHandler_DeleteItem_Success(t *testing.T) {
	r, mockService := setupBillingHandlerTest()

	mockService.On("DeleteItem", uint(1), uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/billings/1/items/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
