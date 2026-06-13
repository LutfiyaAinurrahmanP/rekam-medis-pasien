package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBillingRepository_List_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	query := *mocks.NewBillingPaginationQuery(1, 10)
	expectedList := mocks.NewTestBillingList(2)
	expectedMeta := dto.BillingPaginationMeta{TotalItems: 2, TotalPages: 1}

	mockRepo.On("List", query).Return(expectedList, expectedMeta, nil)

	list, meta, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, int64(2), meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_DeletedList_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	query := *mocks.NewBillingPaginationQuery(1, 10)
	expectedList := mocks.NewTestBillingList(1)
	expectedMeta := dto.BillingPaginationMeta{TotalItems: 1, TotalPages: 1}

	mockRepo.On("DeletedList", query).Return(expectedList, expectedMeta, nil)

	list, meta, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, int64(1), meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	expectedItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(expectedItem, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedItem.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	expectedItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(expectedItem, nil)

	res, err := mockRepo.FindByIDUnscoped(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedItem.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_FindByInvoiceNumber_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	expectedItem := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("FindByInvoiceNumber", "INV-2024-0001").Return(expectedItem, nil)

	res, err := mockRepo.FindByInvoiceNumber("INV-2024-0001")
	assert.NoError(t, err)
	assert.Equal(t, expectedItem.InvoiceNumber, res.InvoiceNumber)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_FindByPatientID_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	query := *mocks.NewBillingPaginationQuery(1, 10)
	expectedList := mocks.NewTestBillingList(2)
	expectedMeta := dto.BillingPaginationMeta{TotalItems: 2, TotalPages: 1}

	mockRepo.On("FindByPatientID", uint(1), query).Return(expectedList, expectedMeta, nil)

	list, meta, err := mockRepo.FindByPatientID(1, query)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, int64(2), meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	item := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("Create", item).Return(nil)

	err := mockRepo.Create(item)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	item := mocks.NewTestBillingWithData(1, 1)

	mockRepo.On("Update", item).Return(nil)

	err := mockRepo.Update(item)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := mockRepo.Delete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Items

func TestBillingRepository_ListItems_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	query := dto.BillingItemPaginationQuery{}
	expectedList := mocks.NewTestBillingItemList(2)

	mockRepo.On("ListItems", uint(1), query).Return(expectedList, nil)

	list, err := mockRepo.ListItems(1, query)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_FindItemByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	expectedItem := mocks.NewTestBillingItemWithData(1, 1)

	mockRepo.On("FindItemByID", uint(1)).Return(expectedItem, nil)

	res, err := mockRepo.FindItemByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedItem.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_CreateItem_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	item := mocks.NewTestBillingItemWithData(1, 1)

	mockRepo.On("CreateItem", item).Return(nil)

	err := mockRepo.CreateItem(item)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_UpdateItem_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)
	item := mocks.NewTestBillingItemWithData(1, 1)

	mockRepo.On("UpdateItem", item).Return(nil)

	err := mockRepo.UpdateItem(item)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestBillingRepository_DeleteItem_Success(t *testing.T) {
	mockRepo := new(mocks.MockBillingRepository)

	mockRepo.On("DeleteItem", uint(1)).Return(nil)

	err := mockRepo.DeleteItem(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
