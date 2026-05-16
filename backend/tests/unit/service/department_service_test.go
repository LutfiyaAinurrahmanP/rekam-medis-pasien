package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	departmentservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/department"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============= Test Cases: GetDepartmentByID =============

func TestGetDepartmentByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	departmentModel := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Heart Department", "Floor 2")

	mockRepo.On("FindById", uint(1)).Return(departmentModel, nil)
	service := departmentservice.NewDepartmentService(mockRepo, cfg)

	result, err := service.GetDepartmentByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Cardiology", result.Name)
	assert.Equal(t, "CARD001", result.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetDepartmentByID_DepartmentNotFound(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	mockRepo.On("FindById", uint(999)).Return(nil, errors.New("department not found"))
	service := departmentservice.NewDepartmentService(mockRepo, cfg)

	result, err := service.GetDepartmentByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "department not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: ListDepartments =============

func TestListDepartments_Success(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	departments := []models.Department{
		*mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Heart Department", "Floor 2"),
		*mocks.NewTestDepartmentWithData(2, "Neurology", "NEURO001", "Brain Department", "Floor 3"),
	}

	query := mocks.NewDepartmentPaginationQuery(1, 10)
	mockRepo.On("List", query).Return(departments, int64(2), nil)
	service := departmentservice.NewDepartmentService(mockRepo, cfg)

	result, err := service.ListDepartments(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, int64(2), result.Meta.TotalItems)
	assert.Equal(t, 1, result.Meta.TotalPages)
	mockRepo.AssertExpectations(t)
}

func TestListDepartments_WithSearch(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	departments := []models.Department{
		*mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Heart Department", "Floor 2"),
	}

	query := &dto.DepartmentPaginationQuery{
		Page:     1,
		PageSize: 10,
		Search:   "Cardiology",
		SortBy:   "created_at",
		SortDir:  "desc",
	}

	mockRepo.On("List", mock.MatchedBy(func(q *dto.DepartmentPaginationQuery) bool {
		return q.Search == "Cardiology"
	})).Return(departments, int64(1), nil)

	service := departmentservice.NewDepartmentService(mockRepo, cfg)
	result, err := service.ListDepartments(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Data))
	mockRepo.AssertExpectations(t)
}

func TestListDepartments_Empty(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	query := mocks.NewDepartmentPaginationQuery(1, 10)
	mockRepo.On("List", query).Return([]models.Department{}, int64(0), nil)
	service := departmentservice.NewDepartmentService(mockRepo, cfg)

	result, err := service.ListDepartments(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Data))
	assert.Equal(t, int64(0), result.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestListDepartments_DefaultPagination(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	query := &dto.DepartmentPaginationQuery{
		Page:     0,
		PageSize: 0,
	}

	departments := mocks.NewTestDepartmentList(10)
	mockRepo.On("List", mock.MatchedBy(func(q *dto.DepartmentPaginationQuery) bool {
		return q.Page == 1 && q.PageSize == 10
	})).Return(departments, int64(10), nil)

	service := departmentservice.NewDepartmentService(mockRepo, cfg)
	result, err := service.ListDepartments(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, len(result.Data))
}

// ============= Test Cases: CreateDepartment =============

func TestCreateDepartment_Success(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	req := mocks.NewCreateDepartmentRequest("Cardiology", "CARD001", "Heart Department", "Floor 2")

	mockRepo.On("IsCodeExists", "CARD001", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.MatchedBy(func(d *models.Department) bool {
		return d.Name == "Cardiology" && d.Code == "CARD001"
	})).Return(nil)

	service := departmentservice.NewDepartmentService(mockRepo, cfg)
	result, err := service.CreateDepartment(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Cardiology", result.Name)
	assert.Equal(t, "CARD001", result.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateDepartment_CodeAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	req := mocks.NewCreateDepartmentRequest("Cardiology", "CARD001", "Heart Department", "Floor 2")

	mockRepo.On("IsCodeExists", "CARD001", mock.Anything).Return(true, nil)
	service := departmentservice.NewDepartmentService(mockRepo, cfg)

	result, err := service.CreateDepartment(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "code already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateDepartment_CreateFailed(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	req := mocks.NewCreateDepartmentRequest("Cardiology", "CARD001", "Heart Department", "Floor 2")

	mockRepo.On("IsCodeExists", "CARD001", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.Anything).Return(errors.New("database error"))

	service := departmentservice.NewDepartmentService(mockRepo, cfg)
	result, err := service.CreateDepartment(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: UpdateDepartment =============

func TestUpdateDepartment_Success(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	existingDept := mocks.NewTestDepartmentWithData(1, "Old Cardiology", "CARD001", "Old Description", "Floor 2")
	req := mocks.NewUpdateDepartmentRequest("New Cardiology", "CARD001", "New Description", "Floor 3")

	mockRepo.On("FindById", uint(1)).Return(existingDept, nil)
	mockRepo.On("Update", mock.MatchedBy(func(d *models.Department) bool {
		return d.Name == "New Cardiology" && d.Description == "New Description"
	})).Return(nil)

	service := departmentservice.NewDepartmentService(mockRepo, cfg)
	result, err := service.UpdateDepartment(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Cardiology", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateDepartment_DepartmentNotFound(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	req := mocks.NewUpdateDepartmentRequest("New Cardiology", "CARD001", "New Description", "Floor 3")

	mockRepo.On("FindById", uint(999)).Return(nil, errors.New("department not found"))
	service := departmentservice.NewDepartmentService(mockRepo, cfg)

	result, err := service.UpdateDepartment(999, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "department not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestUpdateDepartment_ChangeCode(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	existingDept := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Description", "Floor 2")
	newCode := "CARD002"
	req := &dto.UpdateDepartmentRequest{
		Code: &newCode,
	}

	mockRepo.On("FindById", uint(1)).Return(existingDept, nil)
	mockRepo.On("IsCodeExists", "CARD002", []uint{uint(1)}).Return(false, nil)
	mockRepo.On("Update", mock.MatchedBy(func(d *models.Department) bool {
		return d.Code == "CARD002"
	})).Return(nil)

	service := departmentservice.NewDepartmentService(mockRepo, cfg)
	result, err := service.UpdateDepartment(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateDepartment_CodeAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	existingDept := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Description", "Floor 2")
	newCode := "CARD002"
	req := &dto.UpdateDepartmentRequest{
		Code: &newCode,
	}

	mockRepo.On("FindById", uint(1)).Return(existingDept, nil)
	mockRepo.On("IsCodeExists", "CARD002", []uint{uint(1)}).Return(true, nil)

	service := departmentservice.NewDepartmentService(mockRepo, cfg)
	result, err := service.UpdateDepartment(1, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "code already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: SoftDeleteDepartment =============

func TestSoftDeleteDepartment_Success(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	dept := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Description", "Floor 2")

	mockRepo.On("FindById", uint(1)).Return(dept, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	service := departmentservice.NewDepartmentService(mockRepo, cfg)
	err := service.SoftDeleteDepartment(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSoftDeleteDepartment_DepartmentNotFound(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	mockRepo.On("FindById", uint(999)).Return(nil, errors.New("department not found"))
	service := departmentservice.NewDepartmentService(mockRepo, cfg)

	err := service.SoftDeleteDepartment(999)

	assert.Error(t, err)
	assert.Equal(t, "department not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: RestoreDepartment =============

func TestRestoreDepartment_Success(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	mockRepo.On("Restore", uint(1)).Return(nil)
	service := departmentservice.NewDepartmentService(mockRepo, cfg)

	err := service.RestoreDepartment(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: HardDeleteDepartment =============

func TestHardDeleteDepartment_Success(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{}

	mockRepo.On("HardDelete", uint(1)).Return(nil)
	service := departmentservice.NewDepartmentService(mockRepo, cfg)

	err := service.HardDeleteDepartment(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: DeleteListDepartments =============

func TestDeleteListDepartments_Success(t *testing.T) {
	mockRepo := new(mocks.MockDepartmentRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	deletedDepartments := []models.Department{
		*mocks.NewTestDepartmentWithData(1, "Old Cardiology", "CARD001", "Description", "Floor 2"),
	}

	query := mocks.NewDepartmentPaginationQuery(1, 10)
	mockRepo.On("DeleteList", query).Return(deletedDepartments, int64(1), nil)
	service := departmentservice.NewDepartmentService(mockRepo, cfg)

	result, err := service.DeleteListDepartments(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Data))
	mockRepo.AssertExpectations(t)
}
