package service

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	medicalCondition "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/medicalCondition"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestMedicalConditionService() (*mocks.MockMedicalConditionRepository, medicalCondition.MedicalConditionService) {
	mockRepo := new(mocks.MockMedicalConditionRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := medicalCondition.NewMedicalConditionService(mockRepo, cfg)
	return mockRepo, service
}

func TestMedicalConditionService_List_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalConditionService()
	query := mocks.NewMedicalConditionPaginationQuery(1, 10)
	records := mocks.NewTestMedicalConditionList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicalConditionService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalConditionService()
	record := mocks.NewTestMedicalConditionWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestMedicalConditionService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalConditionService()
	req := mocks.NewCreateMedicalConditionRequest(1)

	mockRepo.On("Create", mock.AnythingOfType("*models.MedicalCondition")).Return(nil)
	mockRepo.On("FindByID", mock.Anything).Return(mocks.NewTestMedicalConditionWithData(1, 1), nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.PatientID, res.PatientID)
	mockRepo.AssertExpectations(t)
}

func TestMedicalConditionService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalConditionService()
	record := mocks.NewTestMedicalConditionWithData(1, 1)
	req := mocks.NewUpdateMedicalConditionRequest()

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.MedicalCondition")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestMedicalConditionService_Delete_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalConditionService()
	record := mocks.NewTestMedicalConditionWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
