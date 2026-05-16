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

// ============= Test Cases: ListDepartments =============

func TestListDepartments_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	expectedDepartments := mocks.NewTestDepartmentResponseList(2)
	expectedResponse := &dto.DepartmentListResponse{
		Data: expectedDepartments,
		Meta: dto.DepartmentPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 2,
			TotalPages: 1,
		},
	}

	mockService.On("ListDepartments", mock.Anything).Return(expectedResponse, nil)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments?page=1&page_size=10", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ListDepartments(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Departments retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])
}

func TestListDepartments_InvalidPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments?page=invalid", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ListDepartments(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListDepartments_WithSearch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	expectedDepartments := mocks.NewTestDepartmentResponseList(1)
	expectedResponse := &dto.DepartmentListResponse{
		Data: expectedDepartments,
		Meta: dto.DepartmentPaginationMeta{
			Page:       1,
			PageSize:   10,
			TotalItems: 1,
			TotalPages: 1,
		},
	}

	mockService.On("ListDepartments", mock.MatchedBy(func(q *dto.DepartmentPaginationQuery) bool {
		return q.Search == "Cardiology"
	})).Return(expectedResponse, nil)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments?page=1&search=Cardiology", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ListDepartments(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestListDepartments_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	mockService.On("ListDepartments", mock.Anything).Return(nil, errors.New("database error"))

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments?page=1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.ListDepartments(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: GetDepartmentByID =============

func TestGetDepartmentByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	expectedDept := mocks.NewTestDepartmentResponseWithData(1, "Cardiology", "CARD001", "Heart Department", "Floor 2")

	mockService.On("GetDepartmentByID", uint(1)).Return(expectedDept, nil)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetDepartmentByID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Department retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])
}

func TestGetDepartmentByID_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments/invalid", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.GetDepartmentByID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Invalid department ID", response["message"])
}

func TestGetDepartmentByID_DepartmentNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	mockService.On("GetDepartmentByID", uint(999)).Return(nil, errors.New("department not found"))

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/departments/999", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.GetDepartmentByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Department not found", response["message"])
	assert.Equal(t, "department not found", response["error"])
}

// ============= Test Cases: CreateDepartment =============

func TestCreateDepartment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	createReq := mocks.NewCreateDepartmentRequest("Cardiology", "CARD001", "Heart Department", "Floor 2")
	expectedResp := mocks.NewTestDepartmentResponseWithData(1, "Cardiology", "CARD001", "Heart Department", "Floor 2")

	mockService.On("CreateDepartment", mock.MatchedBy(func(req *dto.CreateDepartmentRequest) bool {
		return req.Name == "Cardiology" && req.Code == "CARD001"
	})).Return(expectedResp, nil)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/departments", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateDepartment(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Department created successfully", response["message"])
	assert.NotNil(t, response["data"])
}

func TestCreateDepartment_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	invalidBody := `{"name": "ab", "code": "C"}` // Too short
	req, _ := http.NewRequest("POST", "/departments", strings.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateDepartment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateDepartment_CodeAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	createReq := mocks.NewCreateDepartmentRequest("Cardiology", "CARD001", "Heart Department", "Floor 2")

	mockService.On("CreateDepartment", mock.Anything).Return(nil, errors.New("code already exists"))

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", "/departments", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.CreateDepartment(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Duplicate data", response["message"])
	assert.Equal(t, "code already exists", response["error"])
}

// ============= Test Cases: UpdateDepartment =============

func TestUpdateDepartment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	updateReq := mocks.NewUpdateDepartmentRequest("Updated Cardiology", "CARD001", "Updated Description", "Floor 3")
	expectedResp := mocks.NewTestDepartmentResponseWithData(1, "Updated Cardiology", "CARD001", "Updated Description", "Floor 3")

	mockService.On("UpdateDepartment", uint(1), mock.Anything).Return(expectedResp, nil)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/departments/1", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.UpdateDepartment(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Department updated successfully", response["message"])
}

func TestUpdateDepartment_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/departments/invalid", nil)
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.UpdateDepartment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateDepartment_DepartmentNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	updateReq := mocks.NewUpdateDepartmentRequest("Updated", "CARD001", "Updated Description", "Floor 3")

	mockService.On("UpdateDepartment", uint(999), mock.Anything).Return(nil, errors.New("department not found"))

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/departments/999", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.UpdateDepartment(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateDepartment_CodeAlreadyExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	updateReq := mocks.NewUpdateDepartmentRequest("Cardiology", "CARD002", "Description", "Floor 3")

	mockService.On("UpdateDepartment", uint(1), mock.Anything).Return(nil, errors.New("code already exists"))

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(updateReq)
	req, _ := http.NewRequest("PUT", "/departments/1", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.UpdateDepartment(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Duplicate data", response["message"])
}

// ============= Test Cases: SoftDeleteDepartment =============

func TestSoftDeleteDepartment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	mockService.On("SoftDeleteDepartment", uint(1)).Return(nil)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/departments/1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.SoftDeleteDepartment(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Department deleted successfully", response["message"])
}

func TestSoftDeleteDepartment_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/departments/invalid", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.SoftDeleteDepartment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSoftDeleteDepartment_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	mockService.On("SoftDeleteDepartment", uint(999)).Return(errors.New("department not found"))

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/departments/999", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.SoftDeleteDepartment(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: RestoreDepartment =============

func TestRestoreDepartment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	mockService.On("RestoreDepartment", uint(1)).Return(nil)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/departments/1/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.RestoreDepartment(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Department restored successfully", response["message"])
}

func TestRestoreDepartment_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/departments/invalid/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.RestoreDepartment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRestoreDepartment_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	mockService.On("RestoreDepartment", uint(999)).Return(errors.New("department not found"))

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/departments/999/restore", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.RestoreDepartment(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

// ============= Test Cases: HardDeleteDepartment =============

func TestHardDeleteDepartment_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	mockService.On("HardDeleteDepartment", uint(1)).Return(nil)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/departments/1/hard-delete", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.HardDeleteDepartment(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Department permanently deleted", response["message"])
}

func TestHardDeleteDepartment_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/departments/invalid/hard-delete", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "invalid"}}

	h.HardDeleteDepartment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHardDeleteDepartment_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.MockDepartmentService)

	mockService.On("HardDeleteDepartment", uint(999)).Return(errors.New("department not found"))

	h := handler.NewDepartmentHandler(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/departments/999/hard-delete", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	h.HardDeleteDepartment(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}
