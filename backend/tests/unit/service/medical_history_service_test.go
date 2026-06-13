package service

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	medicalhistory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func setupTestMedicalHistoryService() (*mocks.MockMedicalHistoryRepository, medicalhistory.MedicalHistoryService) {
	mockRepo := new(mocks.MockMedicalHistoryRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := medicalhistory.NewMedicalHistoryService(mockRepo, cfg)
	return mockRepo, service
}

func TestMedicalHistoryService_List_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalHistoryService()
	query := mocks.NewMedicalHistoryPaginationQuery(1, 10)
	records := mocks.NewTestMedicalHistoryPatientList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicalHistoryService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalHistoryService()
	record := mocks.NewTestMedicalHistoryPatientWithData(1)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	assert.Len(t, res.Allergies, 1)
	assert.Len(t, res.MedicalConditions, 1)
	assert.Len(t, res.SurgicalHistories, 1)
	assert.Len(t, res.FamilyHistories, 1)
	mockRepo.AssertExpectations(t)
}
