package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLabTestRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockLabTestRepository)
	query := mocks.NewLabTestPaginationQuery(1, 10)
	records := mocks.NewTestLabTestList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestLabTestRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockLabTestRepository)
	query := mocks.NewLabTestPaginationQuery(1, 10)
	records := mocks.NewTestLabTestList(2)

	mockRepo.On("DeletedList", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestLabTestRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockLabTestRepository)
	expectedRecord := mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false)

	mockRepo.On("FindByID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestLabTestRepository_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo := new(mocks.MockLabTestRepository)
	expectedRecord := mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByIDUnscoped(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestLabTestRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockLabTestRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.LabTest")).Return(nil)

	err := mockRepo.Create(&models.LabTest{MedicalRecordID: 1, TestTypeID: 1, OrderedByDoctorID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLabTestRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockLabTestRepository)

	mockRepo.On("Update", mock.AnythingOfType("*models.LabTest")).Return(nil)

	err := mockRepo.Update(&models.LabTest{ID: 1, Status: "completed"})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLabTestRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockLabTestRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLabTestRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockLabTestRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLabTestRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockLabTestRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
