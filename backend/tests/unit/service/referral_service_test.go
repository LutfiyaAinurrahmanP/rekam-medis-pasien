package service

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/referral"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestReferralService() (*mocks.MockReferralRepository, referral.ReferralService) {
	mockRepo := new(mocks.MockReferralRepository)
	service := referral.NewReferralService(mockRepo)
	return mockRepo, service
}

func TestReferralService_List_Success(t *testing.T) {
	mockRepo, service := setupTestReferralService()
	query := mocks.NewReferralPaginationQuery(1, 10)
	records := mocks.NewTestReferralList(2)
	meta := dto.ReferralPaginationMeta{TotalItems: 2}

	mockRepo.On("List", *query).Return(records, meta, nil)

	res, err := service.List(*query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestReferralService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestReferralService()
	record := mocks.NewTestReferralWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(record, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, record.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestReferralService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestReferralService()
	req := mocks.NewCreateReferralRequest(1)

	mockRepo.On("GenerateReferralNumber").Return("REF-2024-000001", nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.Referral")).Return(nil)
	mockRepo.On("FindByID", mock.Anything).Return(mocks.NewTestReferralWithData(1, 1), nil)

	res, err := service.Create(*req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.PatientID, res.PatientID)
	mockRepo.AssertExpectations(t)
}

func TestReferralService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestReferralService()
	record := mocks.NewTestReferralWithData(1, 1)
	req := mocks.NewUpdateReferralRequest()

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Referral")).Return(nil)

	res, err := service.Update(1, *req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestReferralService_Accept_Success(t *testing.T) {
	mockRepo, service := setupTestReferralService()
	record := mocks.NewTestReferralWithData(1, 1)
	req := dto.AcceptReferralRequest{Notes: "Ok"}

	mockRepo.On("FindByID", uint(1)).Return(record, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Referral")).Return(nil)

	res, err := service.Accept(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestReferralService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestReferralService()

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
