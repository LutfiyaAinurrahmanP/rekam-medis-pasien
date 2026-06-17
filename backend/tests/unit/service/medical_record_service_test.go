package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	medicalrecord "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-record"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestMedicalRecordService() (*mocks.MockMedicalRecordRepository, medicalrecord.MedicalRecordService) {
	mockRepo := new(mocks.MockMedicalRecordRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := medicalrecord.NewMedicalRecordService(mockRepo, cfg)
	return mockRepo, service
}

func TestMedicalRecordService_List_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	query := mocks.NewMedicalRecordPaginationQuery(1, 10)
	query.SortBy = "visit_date"
	query.SortDir = "desc"
	records := mocks.NewTestMedicalRecordList(2)

	mockRepo.On("List", query).Return(records, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	query := mocks.NewMedicalRecordPaginationQuery(1, 10)
	query.SortBy = "visit_date"
	query.SortDir = "desc"
	records := mocks.NewTestMedicalRecordList(2)

	mockRepo.On("DeletedList", query).Return(records, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	record := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_FindByID_NotFound(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("record not found"))

	res, err := service.FindByID(99)

	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	record := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)

	res, err := service.FindByIDUnscoped(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	req := mocks.NewCreateMedicalRecordRequest(1, "2023-12-01", "Complaint", "Diagnosis", "Plan")

	mockRepo.On("Create", mock.AnythingOfType("*models.MedicalRecord")).Return(nil)
	mockRepo.On("FindByID", mock.Anything).Return(mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", false), nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.PatientID, res.PatientID)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	record := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", false)
	req := mocks.NewUpdateMedicalRecordRequest("New Complaint", "New Diagnosis", "New Plan")

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.MedicalRecord")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "New Complaint", res.ChiefComplaint)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_Update_NotDraft(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	record := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "finalized", false)
	req := mocks.NewUpdateMedicalRecordRequest("New Complaint", "New Diagnosis", "New Plan")

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.Update(1, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "medical record is not in draft status", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_Finalize_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	record := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Finalize", uint(1)).Return(nil)

	err := service.Finalize(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_Finalize_NotDraft(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	record := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "finalized", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	err := service.Finalize(1)

	assert.Error(t, err)
	assert.Equal(t, "medical record is not in draft status", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	record := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_SoftDelete_NotDraft(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	record := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "finalized", false)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	err := service.SoftDelete(1)

	assert.Error(t, err)
	assert.Equal(t, "medical record is not in draft status", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	record := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)
	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMedicalRecordService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestMedicalRecordService()
	record := mocks.NewTestMedicalRecordWithData(1, 1, 1, "2023-12-01", "Complaint", "Diagnosis", "Plan", "draft", true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(record, nil)
	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
