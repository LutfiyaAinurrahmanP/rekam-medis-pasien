package repository

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTypeTestRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)
	query := mocks.NewTypeTestPaginationQuery(1, 10)
	typeTests := mocks.NewTestTypeTestList(2)

	mockRepo.On("List", query).Return(typeTests, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)
	query := mocks.NewTypeTestPaginationQuery(1, 10)
	typeTests := mocks.NewTestTypeTestList(2)

	mockRepo.On("DeletedList", query).Return(typeTests, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_ActiveList(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)
	query := mocks.NewTypeTestPaginationQuery(1, 10)
	typeTests := mocks.NewTestTypeTestList(2)

	mockRepo.On("ActiveList", query).Return(typeTests, int64(2), nil)

	res, total, err := mockRepo.ActiveList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_InactiveList(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)
	query := mocks.NewTypeTestPaginationQuery(1, 10)
	typeTests := mocks.NewTestTypeTestList(2)

	mockRepo.On("InactiveList", query).Return(typeTests, int64(2), nil)

	res, total, err := mockRepo.InactiveList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)
	typeTest := mocks.NewTestTypeTestWithData(1, "Test1", "T001", 1, "Desc1", 50000, true)

	mockRepo.On("FindByID", uint(1)).Return(typeTest, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, typeTest.Name, res.Name)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_FindByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("type test not found"))

	res, err := mockRepo.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "type test not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)
	typeTest := mocks.NewTestTypeTestWithData(1, "Test1", "T001", 1, "Desc1", 50000, true)

	mockRepo.On("Create", typeTest).Return(nil)

	err := mockRepo.Create(typeTest)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)
	typeTest := mocks.NewTestTypeTestWithData(1, "Test1 Updated", "T001", 1, "Desc1", 50000, true)

	mockRepo.On("Update", typeTest).Return(nil)

	err := mockRepo.Update(typeTest)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_Activate_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)

	mockRepo.On("Activate", uint(1)).Return(nil)

	err := mockRepo.Activate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_Deactivate_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)

	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := mockRepo.Deactivate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_IsNameExists_True(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)

	mockRepo.On("IsNameExists", "Test1", mock.Anything).Return(true, nil)

	exists, err := mockRepo.IsNameExists("Test1")
	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestTypeTestRepository_IsCodeExists_False(t *testing.T) {
	mockRepo := new(mocks.MockTypeTestRepository)

	mockRepo.On("IsCodeExists", "T001", mock.Anything).Return(false, nil)

	exists, err := mockRepo.IsCodeExists("T001")
	assert.NoError(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}
