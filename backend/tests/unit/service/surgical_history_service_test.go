package service

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	surgicalHistory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/surgicalHistory"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestSurgicalHistoryService() (*mocks.MockSurgicalHistoryRepository, surgicalHistory.SurgicalHistoryService) {
	mockRepo := new(mocks.MockSurgicalHistoryRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := surgicalHistory.NewSurgicalHistoryService(mockRepo, cfg)
	return mockRepo, service
}

func TestSurgicalHistoryService_List_Success(t *testing.T) {
	mockRepo, service := setupTestSurgicalHistoryService()
	query := mocks.NewSurgicalHistoryPaginationQuery(1, 10)
	records := mocks.NewTestSurgicalHistoryList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestSurgicalHistoryService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestSurgicalHistoryService()
	record := mocks.NewTestSurgicalHistoryWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestSurgicalHistoryService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestSurgicalHistoryService()
	req := mocks.NewCreateSurgicalHistoryRequest(1)

	mockRepo.On("Create", mock.AnythingOfType("*models.SurgicalHistory")).Return(nil)
	mockRepo.On("FindByID", mock.Anything).Return(mocks.NewTestSurgicalHistoryWithData(1, 1), nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.PatientID, res.PatientID)
	mockRepo.AssertExpectations(t)
}

func TestSurgicalHistoryService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestSurgicalHistoryService()
	record := mocks.NewTestSurgicalHistoryWithData(1, 1)
	req := mocks.NewUpdateSurgicalHistoryRequest()

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.SurgicalHistory")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestSurgicalHistoryService_Delete_Success(t *testing.T) {
	mockRepo, service := setupTestSurgicalHistoryService()
	record := mocks.NewTestSurgicalHistoryWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
