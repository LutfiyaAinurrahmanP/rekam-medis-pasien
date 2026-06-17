package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============= Test Cases: List =============

func TestListDoctorSpecializations_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	expectedResponse := &dto.DoctorSpecializationListResponse{
		Data: []dto.DoctorSpecializationResponse{},
		Meta: dto.DoctorSpecializationPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("List", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/specializations?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestListDoctorSpecializations_InvalidQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)
	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/specializations?page=invalid", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.List(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListDoctorSpecializations_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	mockService.On("List", mock.Anything).Return(nil, errors.New("database error"))

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/specializations?page=1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.List(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ============= Test Cases: FindByID =============

func TestFindDoctorSpecializationByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	expectedResp := mocks.NewTestDoctorSpecializationResponseWithData(1, "Cardiology", "CARD", "Heart Specialization", true)

	mockService.On("FindByID", uint(1)).Return(expectedResp, nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/specializations/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.FindByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestFindDoctorSpecializationByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	mockService.On("FindByID", uint(999)).Return(nil, errors.New("doctor specialization not found"))

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/specializations/999", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.FindByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// ============= Test Cases: Create =============

func TestCreateDoctorSpecialization_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	createReq := mocks.NewCreateDoctorSpecializationRequest("Cardiology", "CARD", "Heart Specialization", true)
	expectedResp := mocks.NewTestDoctorSpecializationResponseWithData(1, "Cardiology", "CARD", "Heart Specialization", true)

	mockService.On("Create", mock.Anything).Return(expectedResp, nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/specializations", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// ============= Test Cases: Update =============

func TestUpdateDoctorSpecialization_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	updateReq := mocks.NewUpdateDoctorSpecializationRequest("Cardiology Updated", "CARD", "Heart Specialization", true)
	expectedResp := mocks.NewTestDoctorSpecializationResponseWithData(1, "Cardiology Updated", "CARD", "Heart Specialization", true)

	mockService.On("Update", uint(1), mock.Anything).Return(expectedResp, nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/specializations/1", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateDoctorSpecialization_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	updateReq := mocks.NewUpdateDoctorSpecializationRequest("Cardiology Updated", "CARD", "Heart Specialization", true)

	mockService.On("Update", uint(1), mock.Anything).Return(nil, errors.New("name already exists"))

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/specializations/1", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Update(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

// ============= Test Cases: SoftDelete =============

func TestSoftDeleteDoctorSpecialization_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	mockService.On("SoftDelete", uint(1)).Return(nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/specializations/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SoftDelete(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============= Test Cases: Restore =============

func TestRestoreDoctorSpecialization_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	mockService.On("Restore", uint(1)).Return(nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/specializations/1/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Restore(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============= Test Cases: HardDelete =============

func TestHardDeleteDoctorSpecialization_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	mockService.On("HardDelete", uint(1)).Return(nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/specializations/1/hard-delete", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.HardDelete(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============= Test Cases: DeletedList =============

func TestDeletedListDoctorSpecializations_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	expectedResponse := &dto.DoctorSpecializationDeletedListResponse{
		Data: []dto.DeletedDoctorSpecializationResponse{},
		Meta: dto.DoctorSpecializationPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("DeletedList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/specializations/deleted?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.DeletedList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============= Test Cases: ActiveList =============

func TestActiveListDoctorSpecializations_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	expectedResponse := &dto.DoctorSpecializationListResponse{
		Data: []dto.DoctorSpecializationResponse{},
		Meta: dto.DoctorSpecializationPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("ActiveList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/specializations/active?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ActiveList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============= Test Cases: InactiveList =============

func TestInactiveListDoctorSpecializations_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	expectedResponse := &dto.DoctorSpecializationListResponse{
		Data: []dto.DoctorSpecializationResponse{},
		Meta: dto.DoctorSpecializationPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("InactiveList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/specializations/inactive?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.InactiveList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============= Test Cases: Activate =============

func TestActivateDoctorSpecialization_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	mockService.On("Activate", uint(1)).Return(nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/specializations/1/activate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Activate(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// ============= Test Cases: Deactivate =============

func TestDeactivateDoctorSpecialization_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorSpecializationService)

	mockService.On("Deactivate", uint(1)).Return(nil)

	h := handler.NewDoctorSpecializationHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/specializations/1/deactivate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.Deactivate(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

