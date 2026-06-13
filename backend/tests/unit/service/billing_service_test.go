package service

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/billing"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestBillingService() (*mocks.MockBillingRepository, billing.BillingService) {
	mockRepo := new(mocks.MockBillingRepository)
	service := billing.NewBillingService(mockRepo)
	return mockRepo, service
}

func TestBillingService_List_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	query := *mocks.NewBillingPaginationQuery(1, 10)
	expectedList := mocks.NewTestBillingList(2)
	expectedMeta := dto.BillingPaginationMeta{TotalItems: 2, TotalPages: 1}

	mockRepo.On("List", query).Return(expectedList, expectedMeta, nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	query := *mocks.NewBillingPaginationQuery(1, 10)
	expectedList := mocks.NewTestBillingList(1)
	expectedMeta := dto.BillingPaginationMeta{TotalItems: 1, TotalPages: 1}

	mockRepo.On("DeletedList", query).Return(expectedList, expectedMeta, nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.Len(t, res.Data, 1)
	assert.Equal(t, int64(1), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	expectedItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(expectedItem, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedItem.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	expectedItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(expectedItem, nil)

	res, err := service.FindByIDUnscoped(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedItem.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_FindByInvoiceNumber_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	expectedItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("FindByInvoiceNumber", "INV-2024-0001").Return(expectedItem, nil)

	res, err := service.FindByInvoiceNumber("INV-2024-0001")

	assert.NoError(t, err)
	assert.Equal(t, expectedItem.InvoiceNumber, res.InvoiceNumber)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_FindByPatientID_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	query := *mocks.NewBillingPaginationQuery(1, 10)
	expectedList := mocks.NewTestBillingList(2)
	expectedMeta := dto.BillingPaginationMeta{TotalItems: 2, TotalPages: 1}

	mockRepo.On("List", mock.AnythingOfType("dto.BillingPaginationQuery")).Return(expectedList, expectedMeta, nil)

	res, err := service.FindByPatientID(1, query)

	assert.NoError(t, err)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	req := *mocks.NewCreateBillingRequest(1)
	createdItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("Create", mock.AnythingOfType("*models.Billing")).Return(nil)
	mockRepo.On("FindByID", uint(0)).Return(createdItem, nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	req := *mocks.NewUpdateBillingRequest()
	existingItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(existingItem, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Billing")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_RecordPayment_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	req := *mocks.NewRecordPaymentRequest()
	existingItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(existingItem, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Billing")).Return(nil)

	res, err := service.RecordPayment(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_Cancel_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	existingItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(existingItem, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Billing")).Return(nil)

	res, err := service.Cancel(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_Delete_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.Delete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	existingItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("Restore", uint(1)).Return(nil)
	mockRepo.On("FindByID", uint(1)).Return(existingItem, nil)

	res, err := service.Restore(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Items

func TestBillingService_ListItems_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	query := dto.BillingItemPaginationQuery{}
	expectedList := mocks.NewTestBillingItemList(2)

	mockRepo.On("ListItems", uint(1), query).Return(expectedList, nil)

	res, err := service.ListItems(1, query)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_FindItemByID_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	expectedItem := mocks.NewTestBillingItemWithData(1, 1)

	mockRepo.On("FindItemByID", uint(1)).Return(expectedItem, nil)

	res, err := service.FindItemByID(1, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedItem.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_CreateItem_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	req := *mocks.NewCreateBillingItemRequest()
	billing := mocks.NewTestBillingWithData(1, 1)
	createdItem := mocks.NewTestBillingItemWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(billing, nil)
	mockRepo.On("CreateItem", mock.AnythingOfType("*models.BillingItem")).Return(nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Billing")).Return(nil)
	mockRepo.On("FindItemByID", uint(0)).Return(createdItem, nil)

	res, err := service.CreateItem(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_UpdateItem_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	req := *mocks.NewUpdateBillingItemRequest()
	billing := mocks.NewTestBillingWithData(1, 1)
	existingItem := mocks.NewTestBillingItemWithData(1, 1)

	mockRepo.On("FindItemByID", uint(1)).Return(existingItem, nil)
	mockRepo.On("FindByID", uint(1)).Return(billing, nil)
	mockRepo.On("UpdateItem", mock.AnythingOfType("*models.BillingItem")).Return(nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Billing")).Return(nil)
	// FindItemByID is called again
	// The mock needs to return item since it expects 1
	
	res, err := service.UpdateItem(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestBillingService_DeleteItem_Success(t *testing.T) {
	mockRepo, service := setupTestBillingService()
	billing := mocks.NewTestBillingWithData(1, 1)
	existingItem := mocks.NewTestBillingItemWithData(1, 1)

	mockRepo.On("FindItemByID", uint(1)).Return(existingItem, nil)
	mockRepo.On("FindByID", uint(1)).Return(billing, nil)
	mockRepo.On("DeleteItem", uint(1)).Return(nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Billing")).Return(nil)

	err := service.DeleteItem(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
