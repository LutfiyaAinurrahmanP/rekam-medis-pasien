package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVitalSignRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockVitalSignRepository)
	query := mocks.NewVitalSignPaginationQuery(1, 10)
	records := mocks.NewTestVitalSignList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockVitalSignRepository)
	query := mocks.NewVitalSignPaginationQuery(1, 10)
	records := mocks.NewTestVitalSignList(2)

	mockRepo.On("DeletedList", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockVitalSignRepository)
	expectedRecord := mocks.NewTestVitalSignWithData(1, 1, false)

	mockRepo.On("FindByID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignRepository_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo := new(mocks.MockVitalSignRepository)
	expectedRecord := mocks.NewTestVitalSignWithData(1, 1, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByIDUnscoped(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignRepository_FindByMedicalRecordID_Success(t *testing.T) {
	mockRepo := new(mocks.MockVitalSignRepository)
	expectedRecord := mocks.NewTestVitalSignWithData(1, 1, false)

	mockRepo.On("FindByMedicalRecordID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByMedicalRecordID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockVitalSignRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.VitalSign")).Return(nil)

	err := mockRepo.Create(&models.VitalSign{MedicalRecordID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockVitalSignRepository)

	mockRepo.On("Update", mock.AnythingOfType("*models.VitalSign")).Return(nil)

	err := mockRepo.Update(&models.VitalSign{ID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockVitalSignRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockVitalSignRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockVitalSignRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
