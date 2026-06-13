package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMedicalHistoryRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockMedicalHistoryRepository)
	query := mocks.NewMedicalHistoryPaginationQuery(1, 10)
	records := mocks.NewTestMedicalHistoryPatientList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestMedicalHistoryRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockMedicalHistoryRepository)
	expectedRecord := mocks.NewTestMedicalHistoryPatientWithData(1)

	mockRepo.On("FindByID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}
