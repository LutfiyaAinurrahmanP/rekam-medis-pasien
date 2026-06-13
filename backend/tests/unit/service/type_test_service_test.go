package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	typetest "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/type-test"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestTypeTestService() (*mocks.MockTypeTestRepository, typetest.TypeTestService) {
	mockRepo := new(mocks.MockTypeTestRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := typetest.NewTypeTestService(mockRepo, cfg)
	return mockRepo, service
}

func TestTypeTestService_List_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()
	query := mocks.NewTypeTestPaginationQuery(1, 10)
	typeTests := mocks.NewTestTypeTestList(2)

	mockRepo.On("List", query).Return(typeTests, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()
	query := mocks.NewTypeTestPaginationQuery(1, 10)
	query.SortBy = "deleted_at"
	typeTests := mocks.NewTestTypeTestList(2)

	mockRepo.On("DeletedList", query).Return(typeTests, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_ActiveList_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()
	query := mocks.NewTypeTestPaginationQuery(1, 10)
	typeTests := mocks.NewTestTypeTestList(2)

	mockRepo.On("ActiveList", query).Return(typeTests, int64(2), nil)

	res, err := service.ActiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_InactiveList_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()
	query := mocks.NewTypeTestPaginationQuery(1, 10)
	typeTests := mocks.NewTestTypeTestList(2)

	mockRepo.On("InactiveList", query).Return(typeTests, int64(2), nil)

	res, err := service.InactiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()
	typeTest := mocks.NewTestTypeTestWithData(1, "Test1", "T001", 1, "Desc1", 50000, true)

	mockRepo.On("FindByID", uint(1)).Return(typeTest, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, typeTest.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_FindByID_NotFound(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("type test not found"))

	res, err := service.FindByID(999)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "type test not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()
	req := mocks.NewCreateTypeTestRequest("Test1", "T001", 1, "Desc1", 50000, true)

	mockRepo.On("IsCodeExists", "T001", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.TypeTest")).Return(nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_Create_CodeExists(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()
	req := mocks.NewCreateTypeTestRequest("Test1", "T001", 1, "Desc1", 50000, true)

	mockRepo.On("IsCodeExists", "T001", mock.Anything).Return(true, nil)

	res, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "code already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()
	typeTest := mocks.NewTestTypeTestWithData(1, "Test1", "T001", 1, "Desc1", 50000, true)
	req := mocks.NewUpdateTypeTestRequest("Test2", "T002", 2, "Desc2", 60000, false)

	mockRepo.On("FindByID", uint(1)).Return(typeTest, nil)
	mockRepo.On("IsCodeExists", "T002", []uint{1}).Return(false, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.TypeTest")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Test2", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()
	typeTest := mocks.NewTestTypeTestWithData(1, "Test1", "T001", 1, "Desc1", 50000, true)

	mockRepo.On("FindByID", uint(1)).Return(typeTest, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_Activate_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()

	mockRepo.On("Activate", uint(1)).Return(nil)

	err := service.Activate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestService_Deactivate_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestService()

	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := service.Deactivate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
