package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	specializationservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/doctor-specialization"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============= Test Cases: FindByID =============

func TestDoctorSpecializationService_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{}

	model := mocks.NewTestDoctorSpecializationWithData(1, "Cardiology", "CARD", "Heart Specialization", true)

	mockRepo.On("FindByID", uint(1)).Return(model, nil)
	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)

	result, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "Cardiology", result.Name)
	assert.Equal(t, "CARD", result.Code)
	mockRepo.AssertExpectations(t)
}

func TestDoctorSpecializationService_FindByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{}

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("doctor specialization not found"))
	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)

	result, err := service.FindByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "doctor specialization not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: List =============

func TestDoctorSpecializationService_List_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	specializations := []models.DoctorSpecialization{
		*mocks.NewTestDoctorSpecializationWithData(1, "Cardiology", "CARD", "Heart Specialization", true),
		*mocks.NewTestDoctorSpecializationWithData(2, "Neurology", "NEUR", "Brain Specialization", true),
	}

	query := mocks.NewDoctorSpecializationPaginationQuery(1, 10)
	mockRepo.On("List", query).Return(specializations, int64(2), nil)
	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)

	result, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, int64(2), result.Meta.TotalItems)
	assert.Equal(t, 1, result.Meta.TotalPages)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: Create =============

func TestDoctorSpecializationService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{}

	req := mocks.NewCreateDoctorSpecializationRequest("Cardiology", "CARD", "Heart Specialization", true)

	mockRepo.On("IsNameExists", "Cardiology", mock.Anything).Return(false, nil)
	mockRepo.On("IsCodeExists", "CARD", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.MatchedBy(func(d *models.DoctorSpecialization) bool {
		return d.Name == "Cardiology" && d.Code == "CARD"
	})).Return(nil)

	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)
	result, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Cardiology", result.Name)
	assert.Equal(t, "CARD", result.Code)
	mockRepo.AssertExpectations(t)
}

func TestDoctorSpecializationService_Create_NameAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{}

	req := mocks.NewCreateDoctorSpecializationRequest("Cardiology", "CARD", "Heart Specialization", true)

	mockRepo.On("IsNameExists", "Cardiology", mock.Anything).Return(true, nil)
	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)

	result, err := service.Create(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "name already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: Update =============

func TestDoctorSpecializationService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{}

	existingSpec := mocks.NewTestDoctorSpecializationWithData(1, "Old Cardiology", "CARD", "Old Description", true)
	req := mocks.NewUpdateDoctorSpecializationRequest("New Cardiology", "CARD", "New Description", true)

	mockRepo.On("FindByID", uint(1)).Return(existingSpec, nil)
	mockRepo.On("IsNameExists", "New Cardiology", []uint{uint(1)}).Return(false, nil)
	mockRepo.On("Update", mock.MatchedBy(func(d *models.DoctorSpecialization) bool {
		return d.Name == "New Cardiology" && d.Description == "New Description"
	})).Return(nil)

	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)
	result, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Cardiology", result.Name)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: SoftDelete =============

func TestDoctorSpecializationService_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{}

	spec := mocks.NewTestDoctorSpecializationWithData(1, "Cardiology", "CARD", "Description", true)

	mockRepo.On("FindByID", uint(1)).Return(spec, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)
	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: Restore =============

func TestDoctorSpecializationService_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{}

	mockRepo.On("Restore", uint(1)).Return(nil)
	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: HardDelete =============

func TestDoctorSpecializationService_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{}

	mockRepo.On("HardDelete", uint(1)).Return(nil)
	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: Active / Inactive List =============

func TestDoctorSpecializationService_ActiveList_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	query := mocks.NewDoctorSpecializationPaginationQuery(1, 10)
	mockRepo.On("ActiveList", query).Return([]models.DoctorSpecialization{}, int64(0), nil)
	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)

	result, err := service.ActiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Data))
	mockRepo.AssertExpectations(t)
}

func TestDoctorSpecializationService_InactiveList_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	query := mocks.NewDoctorSpecializationPaginationQuery(1, 10)
	mockRepo.On("InactiveList", query).Return([]models.DoctorSpecialization{}, int64(0), nil)
	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)

	result, err := service.InactiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Data))
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: Activate / Deactivate =============

func TestDoctorSpecializationService_Activate_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{}

	mockRepo.On("Activate", uint(1)).Return(nil)
	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)

	err := service.Activate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDoctorSpecializationService_Deactivate_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorSpecializationRepository)
	cfg := &config.Config{}

	mockRepo.On("Deactivate", uint(1)).Return(nil)
	service := specializationservice.NewDoctorSpecializationService(mockRepo, cfg)

	err := service.Deactivate(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
