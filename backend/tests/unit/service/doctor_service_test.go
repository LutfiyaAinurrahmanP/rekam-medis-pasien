package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	doctorservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/doctor"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestDoctorService() (*mocks.MockDoctorRepository, doctorservice.DoctorService) {
	mockRepo := new(mocks.MockDoctorRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := doctorservice.NewDoctorService(mockRepo, cfg)
	return mockRepo, service
}

func TestDoctorService_GetMyDoctorData_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	doctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	mockRepo.On("FindByUserID", uint(100)).Return(doctor, nil)

	res, err := service.GetMyDoctorData(100)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(100), *res.UserID)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_GetMyDoctorData_NotFound(t *testing.T) {
	mockRepo, service := setupTestDoctorService()

	mockRepo.On("FindByUserID", uint(100)).Return(nil, errors.New("doctor not found"))

	res, err := service.GetMyDoctorData(100)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "doctor not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_UpdateMyDoctorData_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	doctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	req := &dto.UpdateDoctorRequest{
		Phone: mocks.PtrString("0899"),
	}

	mockRepo.On("FindByUserID", uint(100)).Return(doctor, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Doctor")).Return(nil)

	res, err := service.UpdateMyDoctorData(100, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "0899", res.Phone)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_GetDoctorByID_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	doctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	mockRepo.On("FindByID", uint(1)).Return(doctor, nil)

	res, err := service.GetDoctorByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, doctor.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_GetDoctorByID_NotFound(t *testing.T) {
	mockRepo, service := setupTestDoctorService()

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("doctor not found"))

	res, err := service.GetDoctorByID(999)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "doctor not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_List_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	query := mocks.NewDoctorPaginationQuery(1, 10)
	doctors := mocks.NewTestDoctorList(2)

	mockRepo.On("List", query).Return(doctors, int64(2), nil)

	res, err := service.ListDoctors(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	query := mocks.NewDoctorPaginationQuery(1, 10)
	query.SortBy = "deleted_at"
	doctors := mocks.NewTestDoctorList(2)

	mockRepo.On("DeleteList", query).Return(doctors, int64(2), nil)

	res, err := service.DeletedListDoctors(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_ActiveList_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	query := mocks.NewDoctorPaginationQuery(1, 10)
	doctors := mocks.NewTestDoctorList(2)

	mockRepo.On("ActiveList", query).Return(doctors, int64(2), nil)

	res, err := service.ActiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_InactiveList_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	query := mocks.NewDoctorPaginationQuery(1, 10)
	doctors := mocks.NewTestDoctorList(2)

	mockRepo.On("InactiveList", query).Return(doctors, int64(2), nil)

	res, err := service.InactiveList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	req := mocks.NewCreateDoctorRequest(100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	mockRepo.On("IsEmployeeIDExists", "DOC001", mock.Anything).Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.Doctor")).Return(nil)

	res, err := service.CreateDoctor(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.EmployeeID, res.EmployeeID)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_Create_EmployeeIDAlreadyExists(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	req := mocks.NewCreateDoctorRequest(100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	mockRepo.On("IsEmployeeIDExists", "DOC001", mock.Anything).Return(true, nil)

	res, err := service.CreateDoctor(req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Equal(t, "employee already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	doctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	req := mocks.NewUpdateDoctorRequest("Dr. Updated", "0899", "up@hospital.com", 2, 2, false)

	mockRepo.On("FindByID", uint(1)).Return(doctor, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Doctor")).Return(nil)

	res, err := service.UpdateDoctor(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Dr. Updated", res.FullName)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	doctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	mockRepo.On("FindByID", uint(1)).Return(doctor, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDeleteDoctor(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.RestoreDoctor(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDeleteDoctor(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_Activate_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	doctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, false)

	mockRepo.On("FindByID", uint(1)).Return(doctor, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Doctor")).Return(nil)

	res, err := service.ActivateDoctor(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, *res.IsActive)
	mockRepo.AssertExpectations(t)
}

func TestDoctorService_Deactivate_Success(t *testing.T) {
	mockRepo, service := setupTestDoctorService()
	doctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)

	mockRepo.On("FindByID", uint(1)).Return(doctor, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Doctor")).Return(nil)

	res, err := service.DeactivateDoctor(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.False(t, *res.IsActive)
	mockRepo.AssertExpectations(t)
}
