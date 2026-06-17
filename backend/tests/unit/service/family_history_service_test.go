package service

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	familyHistory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/familyHistory"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestFamilyHistoryService() (*mocks.MockFamilyHistoryRepository, familyHistory.FamilyHistoryService) {
	mockRepo := new(mocks.MockFamilyHistoryRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := familyHistory.NewFamilyHistoryService(mockRepo, cfg)
	return mockRepo, service
}

func TestFamilyHistoryService_List_Success(t *testing.T) {
	mockRepo, service := setupTestFamilyHistoryService()
	query := mocks.NewFamilyHistoryPaginationQuery(1, 10)
	records := mocks.NewTestFamilyHistoryList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestFamilyHistoryService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestFamilyHistoryService()
	record := mocks.NewTestFamilyHistoryWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestFamilyHistoryService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestFamilyHistoryService()
	req := mocks.NewCreateFamilyHistoryRequest(1)

	mockRepo.On("Create", mock.AnythingOfType("*models.FamilyHistory")).Return(nil)
	mockRepo.On("FindByID", mock.Anything).Return(mocks.NewTestFamilyHistoryWithData(1, 1), nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.PatientID, res.PatientID)
	mockRepo.AssertExpectations(t)
}

func TestFamilyHistoryService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestFamilyHistoryService()
	record := mocks.NewTestFamilyHistoryWithData(1, 1)
	req := mocks.NewUpdateFamilyHistoryRequest()

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.FamilyHistory")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestFamilyHistoryService_Delete_Success(t *testing.T) {
	mockRepo, service := setupTestFamilyHistoryService()
	record := mocks.NewTestFamilyHistoryWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
