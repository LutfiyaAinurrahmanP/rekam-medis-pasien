package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	vitalsign "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/vital-sign"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestVitalSignService() (*mocks.MockVitalSignRepository, *mocks.MockMedicalRecordRepository, vitalsign.VitalSignService) {
	mockRepo := new(mocks.MockVitalSignRepository)
	mockMedicalRecordRepo := new(mocks.MockMedicalRecordRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := vitalsign.NewVitalSignService(mockRepo, cfg, mockMedicalRecordRepo)
	return mockRepo, mockMedicalRecordRepo, service
}

func TestVitalSignService_List_Success(t *testing.T) {
	mockRepo, _, service := setupTestVitalSignService()
	query := mocks.NewVitalSignPaginationQuery(1, 10)
	records := mocks.NewTestVitalSignList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignService_DeletedList_Success(t *testing.T) {
	mockRepo, _, service := setupTestVitalSignService()
	query := mocks.NewVitalSignPaginationQuery(1, 10)
	records := mocks.NewTestVitalSignList(2)

	mockRepo.On("DeletedList", query).Return(records, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignService_FindByID_Success(t *testing.T) {
	mockRepo, _, service := setupTestVitalSignService()
	record := mocks.NewTestVitalSignWithData(1, 1, false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignService_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo, _, service := setupTestVitalSignService()
	record := mocks.NewTestVitalSignWithData(1, 1, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)

	res, err := service.FindByIDUnscoped(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignService_Create_Success(t *testing.T) {
	mockRepo, mockMedicalRecordRepo, service := setupTestVitalSignService()
	req := mocks.NewCreateVitalSignRequest(1)

	mockMedicalRecordRepo.On("FindByID", uint(1)).Return(&models.MedicalRecord{ID: 1}, nil)
	mockRepo.On("List", mock.AnythingOfType("*dto.VitalSignPaginationQuery")).Return([]models.VitalSign{}, int64(0), nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.VitalSign")).Return(nil)
	mockRepo.On("FindByID", mock.Anything).Return(mocks.NewTestVitalSignWithData(1, 1, false), nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.MedicalRecordID, res.MedicalRecordID)
	mockRepo.AssertExpectations(t)
	mockMedicalRecordRepo.AssertExpectations(t)
}

func TestVitalSignService_Create_AlreadyExists(t *testing.T) {
	mockRepo, mockMedicalRecordRepo, service := setupTestVitalSignService()
	req := mocks.NewCreateVitalSignRequest(1)

	mockMedicalRecordRepo.On("FindByID", uint(1)).Return(&models.MedicalRecord{ID: 1}, nil)
	mockRepo.On("List", mock.AnythingOfType("*dto.VitalSignPaginationQuery")).Return(mocks.NewTestVitalSignList(1), int64(1), nil)

	res, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "vital signs already recorded for this medical record", err.Error())
	mockRepo.AssertExpectations(t)
	mockMedicalRecordRepo.AssertExpectations(t)
}

func TestVitalSignService_Create_MedicalRecordNotFound(t *testing.T) {
	mockRepo, mockMedicalRecordRepo, service := setupTestVitalSignService()
	req := mocks.NewCreateVitalSignRequest(1)

	mockMedicalRecordRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	res, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "medical record not found", err.Error())
	mockRepo.AssertExpectations(t)
	mockMedicalRecordRepo.AssertExpectations(t)
}

func TestVitalSignService_Update_Success(t *testing.T) {
	mockRepo, _, service := setupTestVitalSignService()
	record := mocks.NewTestVitalSignWithData(1, 1, false)
	req := mocks.NewUpdateVitalSignRequest()

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.VitalSign")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignService_SoftDelete_Success(t *testing.T) {
	mockRepo, _, service := setupTestVitalSignService()
	record := mocks.NewTestVitalSignWithData(1, 1, false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignService_Restore_Success(t *testing.T) {
	mockRepo, _, service := setupTestVitalSignService()
	record := mocks.NewTestVitalSignWithData(1, 1, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)
	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVitalSignService_HardDelete_Success(t *testing.T) {
	mockRepo, _, service := setupTestVitalSignService()
	record := mocks.NewTestVitalSignWithData(1, 1, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)
	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
