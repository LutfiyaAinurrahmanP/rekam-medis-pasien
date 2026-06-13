package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPrescriptionRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)
	query := mocks.NewPrescriptionPaginationQuery(1, 10)
	records := mocks.NewTestPrescriptionList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)
	query := mocks.NewPrescriptionPaginationQuery(1, 10)
	records := mocks.NewTestPrescriptionList(2)

	mockRepo.On("DeletedList", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)
	expectedRecord := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", false)

	mockRepo.On("FindByID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)
	expectedRecord := mocks.NewTestPrescriptionWithData(1, 1, 1, "pending", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByIDUnscoped(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.Prescription")).Return(nil)

	err := mockRepo.Create(&models.Prescription{MedicalRecordID: 1, DoctorID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)

	mockRepo.On("Update", mock.AnythingOfType("*models.Prescription")).Return(nil)

	err := mockRepo.Update(&models.Prescription{ID: 1, Status: "dispensed"})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_FindItemByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)
	expectedRecord := mocks.NewTestPrescriptionItemWithData(1, 1, 1)

	mockRepo.On("FindItemByID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindItemByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_CreateItem_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)

	mockRepo.On("CreateItem", mock.AnythingOfType("*models.PrescriptionItem")).Return(nil)

	err := mockRepo.CreateItem(&models.PrescriptionItem{PrescriptionID: 1, MedicineID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_UpdateItem_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)

	mockRepo.On("UpdateItem", mock.AnythingOfType("*models.PrescriptionItem")).Return(nil)

	err := mockRepo.UpdateItem(&models.PrescriptionItem{ID: 1, Dosage: "2 tabs"})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPrescriptionRepository_DeleteItem_Success(t *testing.T) {
	mockRepo := new(mocks.MockPrescriptionRepository)

	mockRepo.On("DeleteItem", uint(1)).Return(nil)

	err := mockRepo.DeleteItem(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
