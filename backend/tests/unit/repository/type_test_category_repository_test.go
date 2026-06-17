package repository

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTypeTestCategoryRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)
	query := mocks.NewTypeTestCategoryPaginationQuery(1, 10)
	categories := mocks.NewTestTypeTestCategoryList(2)

	mockRepo.On("List", query).Return(categories, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)
	query := mocks.NewTypeTestCategoryPaginationQuery(1, 10)
	categories := mocks.NewTestTypeTestCategoryList(2)

	mockRepo.On("DeletedList", query).Return(categories, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_ActiveList(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)
	query := mocks.NewTypeTestCategoryPaginationQuery(1, 10)
	categories := mocks.NewTestTypeTestCategoryList(2)

	mockRepo.On("ActiveList", query).Return(categories, int64(2), nil)

	res, total, err := mockRepo.ActiveList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_InactiveList(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)
	query := mocks.NewTypeTestCategoryPaginationQuery(1, 10)
	categories := mocks.NewTestTypeTestCategoryList(2)

	mockRepo.On("InactiveList", query).Return(categories, int64(2), nil)

	res, total, err := mockRepo.InactiveList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)
	category := mocks.NewTestTypeTestCategoryWithData(1, "Cat1", "C001", "Desc1", true)

	mockRepo.On("FindByID", uint(1)).Return(category, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, category.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_FindByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("type test category not found"))

	res, err := mockRepo.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "type test category not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)
	category := mocks.NewTestTypeTestCategoryWithData(1, "Cat1", "C001", "Desc1", true)

	mockRepo.On("Create", category).Return(nil)

	err := mockRepo.Create(category)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)
	category := mocks.NewTestTypeTestCategoryWithData(1, "Cat1 Updated", "C001", "Desc1", true)

	mockRepo.On("Update", category).Return(nil)

	err := mockRepo.Update(category)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_Activate_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)

	mockRepo.On("Activate", uint(1)).Return(nil)

	err := mockRepo.Activate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_Deactivate_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)

	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := mockRepo.Deactivate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_IsNameExists_True(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)

	mockRepo.On("IsNameExists", "Cat1", mock.Anything).Return(true, nil)

	exists, err := mockRepo.IsNameExists("Cat1")
	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestCategoryRepository_IsCodeExists_False(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestCategoryRepository)

	mockRepo.On("IsCodeExists", "C001", mock.Anything).Return(false, nil)

	exists, err := mockRepo.IsCodeExists("C001")
	assert.NoError(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}
