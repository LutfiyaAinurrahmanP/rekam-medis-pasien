package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	typetestcategory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/type-test-category"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestTypeTestCategoryService() (*mocks.MockTypeTestCategoryRepository, typetestcategory.TypeTestCategoryService) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := typetestcategory.NewTypeTestCategoryService(mockRepo, cfg)
	return mockRepo, service
}

func TestTypeTestCategoryService_List_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()
	query := mocks.NewTypeTestCategoryPaginationQuery(1, 10)
	categories := mocks.NewTestTypeTestCategoryList(2)

	mockRepo.On("List", query).Return(categories, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()
	query := mocks.NewTypeTestCategoryPaginationQuery(1, 10)
	query.SortBy = "deleted_at"
	categories := mocks.NewTestTypeTestCategoryList(2)

	mockRepo.On("DeletedList", query).Return(categories, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_ActiveList_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()
	query := mocks.NewTypeTestCategoryPaginationQuery(1, 10)
	categories := mocks.NewTestTypeTestCategoryList(2)

	mockRepo.On("ActiveList", query).Return(categories, int64(2), nil)

	res, err := service.ActiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_InactiveList_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()
	query := mocks.NewTypeTestCategoryPaginationQuery(1, 10)
	categories := mocks.NewTestTypeTestCategoryList(2)

	mockRepo.On("InactiveList", query).Return(categories, int64(2), nil)

	res, err := service.InactiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()
	category := mocks.NewTestTypeTestCategoryWithData(1, "Cat1", "C001", "Desc1", true)

	mockRepo.On("FindByID", uint(1)).Return(category, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, category.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_FindByID_NotFound(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("type test category not found"))

	res, err := service.FindByID(999)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "type test category not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()
	req := mocks.NewCreateTypeTestCategoryRequest("Cat1", "C001", "Desc1", true)

	mockRepo.On("IsNameExists", "Cat1", mock.Anything).Return(false, nil)
	mockRepo.On("IsCodeExists", "C001", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.TypeTestCategory")).Return(nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_Create_NameExists(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()
	req := mocks.NewCreateTypeTestCategoryRequest("Cat1", "C001", "Desc1", true)

	mockRepo.On("IsNameExists", "Cat1", mock.Anything).Return(true, nil)

	res, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "name already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()
	category := mocks.NewTestTypeTestCategoryWithData(1, "Cat1", "C001", "Desc1", true)
	req := mocks.NewUpdateTypeTestCategoryRequest("Cat2", "C002", "Desc2", false)

	mockRepo.On("FindByID", uint(1)).Return(category, nil)
	mockRepo.On("IsNameExists", "Cat2", []uint{1}).Return(false, nil)
	mockRepo.On("IsCodeExists", "C002", []uint{1}).Return(false, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.TypeTestCategory")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Cat2", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()
	category := mocks.NewTestTypeTestCategoryWithData(1, "Cat1", "C001", "Desc1", true)

	mockRepo.On("FindByID", uint(1)).Return(category, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_Activate_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()

	mockRepo.On("Activate", uint(1)).Return(nil)

	err := service.Activate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryService_Deactivate_Success(t *testing.T) {
	mockRepo, service := setupTestTypeTestCategoryService()

	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := service.Deactivate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
