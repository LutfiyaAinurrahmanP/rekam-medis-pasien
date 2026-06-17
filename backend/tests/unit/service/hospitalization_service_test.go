package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/hospitalization"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestHospitalizationService() (*mocks.MockHospitalizationRepository, hospitalization.HospitalizationService) {
	mockRepo := new(mocks.MockHospitalizationRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := hospitalization.NewHospitalizationService(mockRepo, cfg)
	return mockRepo, service
}

func TestHospitalizationService_List_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	query := mocks.NewHospitalizationPaginationQuery(1, 10)
	query.SortBy = "admission_date"
	query.SortDir = "desc"
	records := mocks.NewTestHospitalizationList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, 2, res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	query := mocks.NewHospitalizationPaginationQuery(1, 10)
	query.SortBy = "admission_date"
	query.SortDir = "desc"
	records := mocks.NewTestHospitalizationList(2)

	mockRepo.On("DeletedList", query).Return(records, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, 2, res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_FindByID_NotFound(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	res, err := service.FindByID(99)

	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)

	res, err := service.FindByIDUnscoped(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	req := mocks.NewCreateHospitalizationRequest(1, 1, 1, "2023-12-01T10:00:00Z", "Reason")

	mockRepo.On("Create", mock.AnythingOfType("*models.Hospitalization")).Return(nil)
	mockRepo.On("FindByID", mock.Anything).Return(mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false), nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.PatientID, res.PatientID)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_Create_InvalidDate(t *testing.T) {
	_, service := setupTestHospitalizationService()
	req := mocks.NewCreateHospitalizationRequest(1, 1, 1, "invalid-date", "Reason")

	res, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "invalid admission_date format")
}

func TestHospitalizationService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false)
	req := mocks.NewUpdateHospitalizationRequest(2, 2, "New Reason")

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", uint(1), mock.Anything).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_Update_Discharged(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "discharged", false)
	req := mocks.NewUpdateHospitalizationRequest(2, 2, "New Reason")

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.Update(1, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "cannot update discharged hospitalization", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_Discharge_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false)
	req := &dto.DischargeHospitalizationRequest{DischargeSummary: "Patient improved"}

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Discharge", uint(1), mock.Anything).Return(nil)

	res, err := service.Discharge(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_Discharge_AlreadyDischarged(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "discharged", false)
	req := &dto.DischargeHospitalizationRequest{DischargeSummary: "Patient improved"}

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.Discharge(1, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "patient has already been discharged", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_Transfer_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false)
	req := &dto.TransferHospitalizationRequest{Notes: "Transfer to ICU"}

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Transfer", uint(1), mock.Anything).Return(nil)

	res, err := service.Transfer(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_Transfer_Discharged(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "discharged", false)
	req := &dto.TransferHospitalizationRequest{Notes: "Transfer to ICU"}

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.Transfer(1, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "cannot transfer discharged patient", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_Activate_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "transferred", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Activate", uint(1), mock.Anything).Return(nil)

	res, err := service.Activate(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_Deactivate_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Deactivate", uint(1), mock.Anything).Return(nil)

	res, err := service.Deactivate(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "cancelled", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_SoftDelete_Admitted(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "admitted", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	err := service.SoftDelete(1)

	assert.Error(t, err)
	assert.Equal(t, "cannot delete admitted hospitalization", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "cancelled", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)
	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestHospitalizationService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestHospitalizationService()
	record := mocks.NewTestHospitalizationWithData(1, 1, 1, 1, "2023-12-01", "10:00:00", "cancelled", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)
	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
