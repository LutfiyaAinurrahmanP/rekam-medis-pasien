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

func TestDoctorHandler_GetMyDoctorData_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	expectedResponse := mocks.NewTestDoctorResponseWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	mockService.On("GetMyDoctorData", uint(100)).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/doctors/me", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(100))

	h.GetMyDoctorData(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDoctorHandler_UpdateMyDoctorData_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	updateReq := &dto.UpdateDoctorRequest{
		Phone: mocks.PtrString("0899"),
	}
	body, _ := json.Marshal(updateReq)

	expectedResponse := mocks.NewTestDoctorResponseWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0899", "test@hospital.com", 1, 1, true)

	mockService.On("UpdateMyDoctorData", uint(100), mock.AnythingOfType("*dto.UpdateDoctorRequest")).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/doctors/me", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(100))

	h.UpdateMyDoctorData(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDoctorHandler_ListDoctors_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	expectedResponse := &dto.DoctorListResponse{
		Data: []dto.DoctorResponse{},
		Meta: dto.DoctorPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("ListDoctors", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/doctors?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ListDoctors(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoctorHandler_DeletedListDoctors_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	expectedResponse := &dto.DoctorDeletedListResponse{
		Data: []dto.DeletedDoctorResponse{},
		Meta: dto.DoctorPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("DeletedListDoctors", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/doctors/deleted?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.DeletedListDoctors(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoctorHandler_ActiveList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	expectedResponse := &dto.DoctorListResponse{
		Data: []dto.DoctorResponse{},
		Meta: dto.DoctorPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("ActiveList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/doctors/active?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ActiveList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoctorHandler_InactiveList_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	expectedResponse := &dto.DoctorListResponse{
		Data: []dto.DoctorResponse{},
		Meta: dto.DoctorPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 0,
			TotalPages: 0,
		},
	}

	mockService.On("InactiveList", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/doctors/inactive?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.InactiveList(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoctorHandler_GetDoctorByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	expectedResponse := mocks.NewTestDoctorResponseWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	mockService.On("GetDoctorByID", uint(1)).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/doctors/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetDoctorByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoctorHandler_CreateDoctor_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	createReq := mocks.NewCreateDoctorRequest(100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)
	body, _ := json.Marshal(createReq)

	expectedResponse := mocks.NewTestDoctorResponseWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	mockService.On("CreateDoctor", mock.AnythingOfType("*dto.CreateDoctorRequest")).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/doctors", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateDoctor(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestDoctorHandler_CreateDoctor_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	createReq := mocks.NewCreateDoctorRequest(100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)
	body, _ := json.Marshal(createReq)

	mockService.On("CreateDoctor", mock.AnythingOfType("*dto.CreateDoctorRequest")).Return(nil, errors.New("employee already exists"))

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/doctors", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateDoctor(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDoctorHandler_UpdateDoctor_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	updateReq := mocks.NewUpdateDoctorRequest("Dr. Updated", "0899", "up@hospital.com", 2, 2, false)
	body, _ := json.Marshal(updateReq)

	expectedResponse := mocks.NewTestDoctorResponseWithData(1, 100, "DOC001", "Dr. Updated", "LIC001", "0899", "up@hospital.com", 2, 2, false)

	mockService.On("UpdateDoctor", uint(1), mock.AnythingOfType("*dto.UpdateDoctorRequest")).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/doctors/1", bytes.NewBuffer(body))
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.UpdateDoctor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoctorHandler_ActivateDoctor_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	expectedResponse := mocks.NewTestDoctorResponseWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)
	mockService.On("ActivateDoctor", uint(1)).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/doctors/1/activate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.ActivateDoctor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoctorHandler_DeactivateDoctor_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	expectedResponse := mocks.NewTestDoctorResponseWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, false)
	mockService.On("DeactivateDoctor", uint(1)).Return(expectedResponse, nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/doctors/1/deactivate", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.DeactivateDoctor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoctorHandler_SoftDeleteDoctor_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	mockService.On("SoftDeleteDoctor", uint(1)).Return(nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/doctors/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SoftDeleteDoctor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoctorHandler_RestoreDoctor_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	mockService.On("RestoreDoctor", uint(1)).Return(nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/doctors/1/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.RestoreDoctor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDoctorHandler_HardDeleteDoctor_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDoctorService)

	mockService.On("HardDeleteDoctor", uint(1)).Return(nil)

	h := handler.NewDoctorHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/doctors/1/hard-delete", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.HardDeleteDoctor(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
