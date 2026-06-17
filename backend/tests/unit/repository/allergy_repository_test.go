package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAllergyRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockAllergyRepository)
	query := mocks.NewAllergyPaginationQuery(1, 10)
	records := mocks.NewTestAllergyList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestAllergyRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockAllergyRepository)
	expectedRecord := mocks.NewTestAllergyWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestAllergyRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockAllergyRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.Allergy")).Return(nil)

	err := mockRepo.Create(&models.Allergy{PatientID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAllergyRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockAllergyRepository)

	mockRepo.On("Update", mock.AnythingOfType("*models.Allergy")).Return(nil)

	err := mockRepo.Update(&models.Allergy{ID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAllergyRepository_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.MockAllergyRepository)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := mockRepo.Delete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
