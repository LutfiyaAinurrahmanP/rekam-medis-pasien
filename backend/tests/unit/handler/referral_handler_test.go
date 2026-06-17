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

func setupReferralHandlerTest() (*gin.Engine, *mocks.MockReferralService) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockReferralService)

	h := handler.NewReferralHandler(mockService)

	r := gin.Default()

	r.GET("/referrals", h.List)
	r.GET("/referrals/:id", h.FindByID)
	r.POST("/referrals", h.Create)
	r.PUT("/referrals/:id", h.Update)
	r.POST("/referrals/:id/accept", h.Accept)
	r.POST("/referrals/:id/reject", h.Reject)
	r.DELETE("/referrals/:id", h.Delete)

	return r, mockService
}

func TestReferralHandler_List_Success(t *testing.T) {
	r, mockService := setupReferralHandlerTest()

	expectedRes := &dto.ReferralListResponse{
		Data: []dto.ReferralResponse{
			*mocks.NewTestReferralResponse(mocks.NewTestReferralWithData(1, 1)),
		},
	}

	mockService.On("List", mock.AnythingOfType("dto.ReferralPaginationQuery")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/referrals?page=1&page_size=10", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestReferralHandler_FindByID_Success(t *testing.T) {
	r, mockService := setupReferralHandlerTest()

	expectedRes := mocks.NewTestReferralResponse(mocks.NewTestReferralWithData(1, 1))

	mockService.On("FindByID", uint(1)).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodGet, "/referrals/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestReferralHandler_Create_Success(t *testing.T) {
	r, mockService := setupReferralHandlerTest()
	createReq := mocks.NewCreateReferralRequest(1)
	expectedRes := mocks.NewTestReferralResponse(mocks.NewTestReferralWithData(1, 1))

	body, _ := json.Marshal(createReq)

	mockService.On("Create", mock.AnythingOfType("dto.CreateReferralRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPost, "/referrals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestReferralHandler_Update_Success(t *testing.T) {
	r, mockService := setupReferralHandlerTest()
	updateReq := mocks.NewUpdateReferralRequest()
	expectedRes := mocks.NewTestReferralResponse(mocks.NewTestReferralWithData(1, 1))

	body, _ := json.Marshal(updateReq)

	mockService.On("Update", uint(1), mock.AnythingOfType("dto.UpdateReferralRequest")).Return(expectedRes, nil)

	req, _ := http.NewRequest(http.MethodPut, "/referrals/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestReferralHandler_Delete_Success(t *testing.T) {
	r, mockService := setupReferralHandlerTest()

	mockService.On("SoftDelete", uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/referrals/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
