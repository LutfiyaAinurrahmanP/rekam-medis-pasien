package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHospitalizationRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)
	query := mocks.NewHospitalizationPaginationQuery(1, 10)
	records := mocks.NewTestHospitalizationList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)
	query := mocks.NewHospitalizationPaginationQuery(1, 10)
	records := mocks.NewTestHospitalizationList(2)

	mockRepo.On("DeletedList", query).Return(records, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)
	expectedRecord := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false)

	mockRepo.On("FindByID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)
	expectedRecord := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByIDUnscoped(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.Hospitalization")).Return(nil)

	err := mockRepo.Create(&models.Hospitalization{PatientID: 1, DoctorID: 1, RoomID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)
	updates := map[string]interface{}{"room_id": uint(2)}

	mockRepo.On("Update", uint(1), updates).Return(nil)

	err := mockRepo.Update(1, updates)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_Discharge_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)
	updates := map[string]interface{}{"status": "discharged"}

	mockRepo.On("Discharge", uint(1), updates).Return(nil)

	err := mockRepo.Discharge(1, updates)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_Transfer_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)
	updates := map[string]interface{}{"status": "transferred"}

	mockRepo.On("Transfer", uint(1), updates).Return(nil)

	err := mockRepo.Transfer(1, updates)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_Activate_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)
	updates := map[string]interface{}{"status": "admitted"}

	mockRepo.On("Activate", uint(1), updates).Return(nil)

	err := mockRepo.Activate(1, updates)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_Deactivate_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)
	updates := map[string]interface{}{"status": "discharged"}

	mockRepo.On("Deactivate", uint(1), updates).Return(nil)

	err := mockRepo.Deactivate(1, updates)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockHospitalizationRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
