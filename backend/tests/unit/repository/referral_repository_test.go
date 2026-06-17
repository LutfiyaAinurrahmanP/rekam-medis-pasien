package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReferralRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockReferralRepository)
	query := mocks.NewReferralPaginationQuery(1, 10)
	records := mocks.NewTestReferralList(2)
	meta := dto.ReferralPaginationMeta{TotalItems: 2}

	mockRepo.On("List", *query).Return(records, meta, nil)

	res, resMeta, err := mockRepo.List(*query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), resMeta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestReferralRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockReferralRepository)
	expectedRecord := mocks.NewTestReferralWithData(1, 1)

	mockRepo.On("FindByID", uint(1)).Return(expectedRecord, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedRecord.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestReferralRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockReferralRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.Referral")).Return(nil)

	err := mockRepo.Create(&models.Referral{PatientID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReferralRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockReferralRepository)

	mockRepo.On("Update", mock.AnythingOfType("*models.Referral")).Return(nil)

	err := mockRepo.Update(&models.Referral{ID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReferralRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockReferralRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReferralRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockReferralRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReferralRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockReferralRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReferralRepository_GenerateReferralNumber_Success(t *testing.T) {
	mockRepo := new(mocks.MockReferralRepository)

	mockRepo.On("GenerateReferralNumber").Return("REF-2024-000001", nil)

	res, err := mockRepo.GenerateReferralNumber()
	assert.NoError(t, err)
	assert.Equal(t, "REF-2024-000001", res)
	mockRepo.AssertExpectations(t)
}
