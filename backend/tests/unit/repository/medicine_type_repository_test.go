package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMedicineTypeRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)
	query := mocks.NewMedicineTypePaginationQuery(1, 10)
	types := mocks.NewTestMedicineTypeList(2)

	mockRepo.On("List", query).Return(types, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)
	query := mocks.NewMedicineTypePaginationQuery(1, 10)
	types := mocks.NewTestMedicineTypeList(2)

	mockRepo.On("DeletedList", query).Return(types, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_ActiveList(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)
	query := mocks.NewMedicineTypePaginationQuery(1, 10)
	types := mocks.NewTestMedicineTypeList(2)

	mockRepo.On("ActiveList", query).Return(types, int64(2), nil)

	res, total, err := mockRepo.ActiveList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_InactiveList(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)
	query := mocks.NewMedicineTypePaginationQuery(1, 10)
	types := mocks.NewTestMedicineTypeList(2)

	mockRepo.On("InactiveList", query).Return(types, int64(2), nil)

	res, total, err := mockRepo.InactiveList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)
	expectedType := mocks.NewTestMedicineTypeWithData(1, "Type1", "T001", "Desc", false)

	mockRepo.On("FindByID", uint(1)).Return(expectedType, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedType.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_FindByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)

	mockRepo.On("FindByID", uint(99)).Return(nil, assert.AnError)

	res, err := mockRepo.FindByID(99)
	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.MedicineType")).Return(nil)

	err := mockRepo.Create(&models.MedicineType{Name: "Type1"})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)

	mockRepo.On("Update", mock.AnythingOfType("*models.MedicineType")).Return(nil)

	err := mockRepo.Update(&models.MedicineType{ID: 1, Name: "Type1"})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_Activate_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)

	mockRepo.On("Activate", uint(1)).Return(nil)

	err := mockRepo.Activate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_Deactivate_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)

	mockRepo.On("Deactivate", uint(1)).Return(nil)

	err := mockRepo.Deactivate(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_IsNameExists_True(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)

	mockRepo.On("IsNameExists", "Type1", mock.Anything).Return(true, nil)

	exists, err := mockRepo.IsNameExists("Type1")
	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestMedicineTypeRepository_IsCodeExists_False(t *testing.T) {
	mockRepo := new(mocks.MockMedicineTypeRepository)

	mockRepo.On("IsCodeExists", "T001", mock.Anything).Return(false, nil)

	exists, err := mockRepo.IsCodeExists("T001")
	assert.NoError(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}
