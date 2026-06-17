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

func setupMedicalConditionHandlerTest() (*gin.Engine, *mocks.MockMedicalConditionService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockMedicalConditionService)

	h := handler.NewMedicalConditionHandler(mockService)

	r := gin.Default()

	r.GET("/medical-history/conditions", h.List)
	r.GET("/medical-history/conditions/:id", h.FindByID)
	r.POST("/medical-history/conditions", h.Create)
	r.PUT("/medical-history/conditions/:id", h.Update)
	r.DELETE("/medical-history/conditions/:id", h.Delete)

	return r, mockService
}

func TestMedicalConditionHandler_List_Success(t *testing.T) {
	r, mockService := setupMedicalConditionHandlerTest()

	expectedRes := &dto.MedicalConditionListResponse{
		Data: []dto.MedicalConditionResponse{
			*mocks.NewTestMedicalConditionResponse(mocks.NewTestMedicalConditionWithData(1, 1)),
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.MedicalConditionPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history/conditions?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalConditionHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupMedicalConditionHandlerTest()

	expectedRes := mocks.NewTestMedicalConditionResponse(mocks.NewTestMedicalConditionWithData(1, 1))

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history/conditions/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalConditionHandler_Create_Success(t *testing.T) {
	r, mockService := setupMedicalConditionHandlerTest()
	createReq := mocks.NewCreateMedicalConditionRequest(1)
	expectedRes := mocks.NewTestMedicalConditionResponse(mocks.NewTestMedicalConditionWithData(1, 1))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateMedicalConditionRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/medical-history/conditions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalConditionHandler_Update_Success(t *testing.T) {
	r, mockService := setupMedicalConditionHandlerTest()
	updateReq := mocks.NewUpdateMedicalConditionRequest()
	expectedRes := mocks.NewTestMedicalConditionResponse(mocks.NewTestMedicalConditionWithData(1, 1))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateMedicalConditionRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/medical-history/conditions/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestMedicalConditionHandler_Delete_Success(t *testing.T) {
	r, mockService := setupMedicalConditionHandlerTest()

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/medical-history/conditions/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
