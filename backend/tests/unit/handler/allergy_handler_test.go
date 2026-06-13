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

func setupAllergyHandlerTest() (*gin.Engine, *mocks.MockAllergyService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockAllergyService)

	h := handler.NewAllergyHandler(mockService)

	r := gin.Default()

	r.GET("/medical-history/allergies", h.List)
	r.GET("/medical-history/allergies/:id", h.FindByID)
	r.POST("/medical-history/allergies", h.Create)
	r.PUT("/medical-history/allergies/:id", h.Update)
	r.DELETE("/medical-history/allergies/:id", h.Delete)

	return r, mockService
}

func TestAllergyHandler_List_Success(t *testing.T) {
	r, mockService := setupAllergyHandlerTest()

	expectedRes := &dto.AllergyListResponse{
		Data: []dto.AllergyResponse{
			*mocks.NewTestAllergyResponse(mocks.NewTestAllergyWithData(1, 1)),
		},
	}

	mockService.On("List", mock.AnythingOfType("*dto.AllergyPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history/allergies?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAllergyHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupAllergyHandlerTest()

	expectedRes := mocks.NewTestAllergyResponse(mocks.NewTestAllergyWithData(1, 1))

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/medical-history/allergies/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAllergyHandler_Create_Success(t *testing.T) {
	r, mockService := setupAllergyHandlerTest()
	createReq := mocks.NewCreateAllergyRequest(1)
	expectedRes := mocks.NewTestAllergyResponse(mocks.NewTestAllergyWithData(1, 1))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("*dto.CreateAllergyRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/medical-history/allergies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestAllergyHandler_Update_Success(t *testing.T) {
	r, mockService := setupAllergyHandlerTest()
	updateReq := mocks.NewUpdateAllergyRequest()
	expectedRes := mocks.NewTestAllergyResponse(mocks.NewTestAllergyWithData(1, 1))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("*dto.UpdateAllergyRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/medical-history/allergies/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestAllergyHandler_Delete_Success(t *testing.T) {
	r, mockService := setupAllergyHandlerTest()

	mockService.On("Delete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/medical-history/allergies/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
