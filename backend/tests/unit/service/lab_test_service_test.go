package service

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	labtest "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/lab-test"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestLabTestService() (*mocks.MockLabTestRepository, labtest.LabTestService) {
	mockRepo := new(mocks.MockLabTestRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := labtest.NewLabTestService(mockRepo, cfg)
	return mockRepo, service
}

func TestLabTestService_List_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	query := mocks.NewLabTestPaginationQuery(1, 10)
	records := mocks.NewTestLabTestList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, 2, res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	query := mocks.NewLabTestPaginationQuery(1, 10)
	records := mocks.NewTestLabTestList(2)

	mockRepo.On("DeletedList", query).Return(records, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, 2, res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)

	res, err := service.FindByIDUnscoped(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	req := mocks.NewCreateLabTestRequest(1, 1, 1)

	mockRepo.On("Create", mock.AnythingOfType("*models.LabTest")).Return(nil)
	mockRepo.On("FindByID", mock.Anything).Return(mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false), nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.MedicalRecordID, res.MedicalRecordID)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false)
	req := mocks.NewUpdateLabTestRequest("New Notes", "10-30")

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.LabTest")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_CollectSample_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.LabTest")).Return(nil)

	res, err := service.CollectSample(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_CollectSample_Cancelled(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "cancelled", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.CollectSample(1)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "cannot collect sample for a cancelled lab test", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_Start_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "sample_collected", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.LabTest")).Return(nil)

	res, err := service.Start(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_Start_Cancelled(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "cancelled", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.Start(1)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "cannot start a cancelled lab test", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_Complete_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "in_progress", false)
	val := "Normal"
	req := &dto.CompleteLabTestRequest{ResultValue: &val}

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.LabTest")).Return(nil)

	res, err := service.Complete(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_Complete_Cancelled(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "cancelled", false)
	val := "Normal"
	req := &dto.CompleteLabTestRequest{ResultValue: &val}

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.Complete(1, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "cannot complete a cancelled lab test", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_Cancel_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.LabTest")).Return(nil)

	res, err := service.Cancel(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_Cancel_Completed(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "completed", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.Cancel(1)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "cannot cancel a completed lab test", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)
	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestLabTestService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestLabTestService()
	record := mocks.NewTestLabTestWithData(1, 1, 1, 1, "ordered", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)
	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
