package service

import (
	"errors"
	"testing"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	patientservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/patient"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func expectedAgeFromDOB(t *testing.T, dateOfBirth string) int {
	t.Helper()
	parsedBirthDate, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		t.Fatalf("failed to parse test birth date %q: %v", dateOfBirth, err)
	}

	now := time.Now()
	age := now.Year() - parsedBirthDate.Year()
	if now.Month() < parsedBirthDate.Month() || (now.Month() == parsedBirthDate.Month() && now.Day() < parsedBirthDate.Day()) {
		age--
	}

	return age
}

// ============= Test Cases: GetPatientByID =============

func TestGetPatientByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	patientModel := mocks.NewTestPatientWithData(1, "PAT001", "John Doe", "1990-01-01", "male", "O")

	mockRepo.On("FindById", uint(1)).Return(patientModel, nil)
	service := patientservice.NewPatientService(mockRepo, cfg)

	result, err := service.GetPatientByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "PAT001", result.PatientCode)
	assert.Equal(t, "John Doe", result.FullName)
	assert.Equal(t, expectedAgeFromDOB(t, "1990-01-01"), result.Age)
	mockRepo.AssertExpectations(t)
}

func TestGetPatientByID_PatientNotFound(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	mockRepo.On("FindById", uint(999)).Return(nil, errors.New("patient not found"))
	service := patientservice.NewPatientService(mockRepo, cfg)

	result, err := service.GetPatientByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "patient not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: GetMyPatientData =============

func TestGetMyPatientData_Success(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	patientModel := mocks.NewTestPatientWithData(1, "PAT001", "John Doe", "1990-01-01", "male", "O")

	mockRepo.On("FindByUserID", uint(1)).Return(patientModel, nil)
	service := patientservice.NewPatientService(mockRepo, cfg)

	result, err := service.GetMyPatientData(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "PAT001", result.PatientCode)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: ListPatients =============

func TestListPatients_Success(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	patients := []models.Patient{
		*mocks.NewTestPatientWithData(1, "PAT001", "John Doe", "1990-01-01", "male", "O"),
		*mocks.NewTestPatientWithData(2, "PAT002", "Jane Doe", "1992-02-02", "female", "A"),
	}

	query := mocks.NewPatientPaginationQuery(1, 10)
	mockRepo.On("List", query).Return(patients, int64(2), nil)
	service := patientservice.NewPatientService(mockRepo, cfg)

	result, err := service.ListPatients(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result.Data))
	assert.Equal(t, int64(2), result.Meta.TotalItems)
	assert.Equal(t, 1, result.Meta.TotalPages)
	mockRepo.AssertExpectations(t)
}

func TestListPatients_Empty(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	query := mocks.NewPatientPaginationQuery(1, 10)
	mockRepo.On("List", query).Return([]models.Patient{}, int64(0), nil)
	service := patientservice.NewPatientService(mockRepo, cfg)

	result, err := service.ListPatients(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Data))
	assert.Equal(t, int64(0), result.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestListPatients_DefaultPagination(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	query := &dto.PatientPaginationQuery{Page: 0, PageSize: 0}
	patients := mocks.NewTestPatientList(10)
	mockRepo.On("List", mock.MatchedBy(func(q *dto.PatientPaginationQuery) bool {
		return q.Page == 1 && q.PageSize == 10 && q.SortBy == "created_at" && q.SortDir == "desc"
	})).Return(patients, int64(10), nil)

	service := patientservice.NewPatientService(mockRepo, cfg)
	result, err := service.ListPatients(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 10, len(result.Data))
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: DeleteListPatients =============

func TestDeleteListPatients_Success(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}

	deletedPatients := []models.Patient{
		*mocks.NewTestPatientWithData(1, "PATDEL001", "Deleted Patient", "1990-01-01", "female", "A"),
	}

	query := mocks.NewPatientPaginationQuery(1, 10)
	query.SortBy = "deleted_at"
	mockRepo.On("DeleteList", query).Return(deletedPatients, int64(1), nil)
	service := patientservice.NewPatientService(mockRepo, cfg)

	result, err := service.DeleteListPatients(query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Data))
	assert.Equal(t, expectedAgeFromDOB(t, "1990-01-01"), result.Data[0].Age)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: CreatePatient =============

func TestCreatePatient_Success(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	userID := uint(1)
	req := &dto.CreatePatientRequest{
		UserID:      &userID,
		PatientCode: "PAT001",
		FullName:    "John Doe",
		DateOfBirth: "1990-01-01",
		Gender:      "male",
		BloodType:   "O",
	}

	mockRepo.On("IsCodeExists", "PAT001", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.Anything).Return(nil)

	service := patientservice.NewPatientService(mockRepo, cfg)
	result, err := service.CreatePatient(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "PAT001", result.PatientCode)
	mockRepo.AssertExpectations(t)
}

func TestCreatePatient_CodeAlreadyExists(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	req := &dto.CreatePatientRequest{
		PatientCode: "PAT001",
		FullName:    "John Doe",
		DateOfBirth: "1990-01-01",
		Gender:      "male",
		BloodType:   "O",
	}

	mockRepo.On("IsCodeExists", "PAT001", mock.Anything).Return(true, nil)
	service := patientservice.NewPatientService(mockRepo, cfg)

	result, err := service.CreatePatient(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "patient code already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: UpdatePatient =============

func TestUpdatePatient_Success(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	existingPatient := mocks.NewTestPatientWithData(1, "PAT001", "John Doe", "1990-01-01", "male", "O")
	newFullName := "John Updated"
	req := &dto.UpdatePatientRequest{FullName: &newFullName}

	mockRepo.On("FindById", uint(1)).Return(existingPatient, nil)
	mockRepo.On("Update", mock.MatchedBy(func(p *models.Patient) bool {
		return p.FullName == "John Updated"
	})).Return(nil)

	service := patientservice.NewPatientService(mockRepo, cfg)
	result, err := service.UpdatePatient(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John Updated", result.FullName)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePatient_PatientNotFound(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	newFullName := "John Updated"
	req := &dto.UpdatePatientRequest{FullName: &newFullName}

	mockRepo.On("FindById", uint(999)).Return(nil, errors.New("patient not found"))
	service := patientservice.NewPatientService(mockRepo, cfg)

	result, err := service.UpdatePatient(999, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "patient not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: UpdateMyPatientData =============

func TestUpdateMyPatientData_Success(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	existingPatient := mocks.NewTestPatientWithData(1, "PAT001", "John Doe", "1990-01-01", "male", "O")
	newPhone := "08111111111"
	newEmail := "updated@example.com"
	req := &dto.UpdatePatientRequest{Phone: &newPhone, Email: &newEmail}

	mockRepo.On("FindByUserID", uint(1)).Return(existingPatient, nil)
	mockRepo.On("Update", mock.MatchedBy(func(p *models.Patient) bool {
		return p.Phone == "08111111111" && p.Email == "updated@example.com"
	})).Return(nil)

	service := patientservice.NewPatientService(mockRepo, cfg)
	result, err := service.UpdateMyPatientData(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "08111111111", result.Phone)
	assert.Equal(t, "updated@example.com", result.Email)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: SoftDeletePatient =============

func TestSoftDeletePatient_Success(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	mockRepo.On("FindById", uint(1)).Return(&models.Patient{ID: 1}, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	service := patientservice.NewPatientService(mockRepo, cfg)
	err := service.SoftDeletePatient(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSoftDeletePatient_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	mockRepo.On("FindById", uint(999)).Return(nil, errors.New("patient not found"))
	service := patientservice.NewPatientService(mockRepo, cfg)

	err := service.SoftDeletePatient(999)

	assert.Error(t, err)
	assert.Equal(t, "patient not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: RestorePatient =============

func TestRestorePatient_Success(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	mockRepo.On("Restore", uint(1)).Return(nil)
	service := patientservice.NewPatientService(mockRepo, cfg)

	err := service.RestorePatient(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ============= Test Cases: HardDeletePatient =============

func TestHardDeletePatient_Success(t *testing.T) {
	mockRepo := new(mocks.MockPatientRepository)
	cfg := &config.Config{}

	mockRepo.On("HardDelete", uint(1)).Return(nil)
	service := patientservice.NewPatientService(mockRepo, cfg)

	err := service.HardDeletePatient(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
