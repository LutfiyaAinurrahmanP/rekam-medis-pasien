package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMedicalRecordRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockMedicalRecordRepository)
	query := mocks.NewMedicalRecordPaginationQuery(1, 10)
	records := mocks.NewTestMedicalRecordList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockMedicalRecordRepository)
	query := mocks.NewMedicalRecordPaginationQuery(1, 10)
	records := mocks.NewTestMedicalRecordList(2)

	mockRepo.On("DeletedList", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalRecordRepository)
	expectedRecord := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", false)

	mockRepo.On("FindByID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordRepository_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalRecordRepository)
	expectedRecord := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByIDUnscoped(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalRecordRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.MedicalRecord")).Return(nil)

	err := mockRepo.Create(&models.MedicalRecord{PatientID: 1, DoctorID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalRecordRepository)

	mockRepo.On("Update", mock.AnythingOfType("*models.MedicalRecord")).Return(nil)

	err := mockRepo.Update(&models.MedicalRecord{ID: 1, PatientID: 1, DoctorID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordRepository_Finalize_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalRecordRepository)

	mockRepo.On("Finalize", uint(1)).Return(nil)

	err := mockRepo.Finalize(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalRecordRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalRecordRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalRecordRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
