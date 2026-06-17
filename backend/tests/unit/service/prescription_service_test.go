package service

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/prescription"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestPrescriptionService() (*mocks.MockPrescriptionRepository, prescription.PrescriptionService) {
	mockRepo := new(mocks.MockPrescriptionRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := prescription.NewPrescriptionService(mockRepo, cfg)
	return mockRepo, service
}

func TestPrescriptionService_List_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	query := mocks.NewPrescriptionPaginationQuery(1, 10)
	records := mocks.NewTestPrescriptionList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	query := mocks.NewPrescriptionPaginationQuery(1, 10)
	records := mocks.NewTestPrescriptionList(2)

	mockRepo.On("DeletedList", query).Return(records, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)

	res, err := service.FindByIDUnscoped(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	req := mocks.NewCreatePrescriptionRequest(1, 1)

	mockRepo.On("Create", mock.AnythingOfType("*models.Prescription")).Return(nil)
	mockRepo.On("FindByID", mock.Anything).Return(mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false), nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.MedicalRecordID, res.MedicalRecordID)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false)
	req := mocks.NewUpdatePrescriptionRequest("New Notes")

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Prescription")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_Dispense_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Prescription")).Return(nil)

	res, err := service.Dispense(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_Dispense_Cancelled(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "cancelled", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.Dispense(1)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "cannot dispense a cancelled prescription", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_Cancel_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Prescription")).Return(nil)

	res, err := service.Cancel(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_Cancel_Dispensed(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "dispensed", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.Cancel(1)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "cannot cancel a dispensed prescription", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)
	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)
	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_ListItems_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.ListItems(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 1)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_FindItemByID_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	item := mocks.NewTestPrescriptionItemWithData(1, 1, 1)

	mockRepo.On("FindItemByID", uint(1)).Return(item, nil)

	res, err := service.FindItemByID(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, item.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_CreateItem_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	record := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false)
	req := mocks.NewCreatePrescriptionItemRequest(1)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("CreateItem", mock.AnythingOfType("*models.PrescriptionItem")).Return(nil)
	mockRepo.On("FindItemByID", mock.Anything).Return(mocks.NewTestPrescriptionItemWithData(1, 1, 1), nil)

	res, err := service.CreateItem(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_UpdateItem_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	item := mocks.NewTestPrescriptionItemWithData(1, 1, 1)
	req := mocks.NewUpdatePrescriptionItemRequest("2 tablets")

	mockRepo.On("FindItemByID", uint(1)).Return(item, nil)
	mockRepo.On("UpdateItem", mock.AnythingOfType("*models.PrescriptionItem")).Return(nil)

	res, err := service.UpdateItem(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_UpdateItem_MismatchPrescription(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	item := mocks.NewTestPrescriptionItemWithData(1, 2, 1)
	req := mocks.NewUpdatePrescriptionItemRequest("2 tablets")

	mockRepo.On("FindItemByID", uint(1)).Return(item, nil)

	res, err := service.UpdateItem(1, 1, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "prescription item not found in this prescription", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionService_DeleteItem_Success(t *testing.T) {
	mockRepo, service := setupTestPrescriptionService()
	item := mocks.NewTestPrescriptionItemWithData(1, 1, 1)

	mockRepo.On("FindItemByID", uint(1)).Return(item, nil)
	mockRepo.On("DeleteItem", uint(1)).Return(nil)

	err := service.DeleteItem(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
