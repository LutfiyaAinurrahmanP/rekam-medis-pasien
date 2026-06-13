package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMedicalConditionRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockMedicalConditionRepository)
	query := mocks.NewMedicalConditionPaginationQuery(1, 10)
	records := mocks.NewTestMedicalConditionList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicalConditionRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalConditionRepository)
	expectedRecord := mocks.NewTestMedicalConditionWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestMedicalConditionRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalConditionRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.MedicalCondition")).Return(nil)

	err := mockRepo.Create(&models.MedicalCondition{PatientID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicalConditionRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalConditionRepository)

	mockRepo.On("Update", mock.AnythingOfType("*models.MedicalCondition")).Return(nil)

	err := mockRepo.Update(&models.MedicalCondition{ID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicalConditionRepository_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalConditionRepository)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := mockRepo.Delete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
