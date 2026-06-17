package repository

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDoctorRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	doctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)
	mockRepo.On("Create", doctor).Return(nil)

	err := mockRepo.Create(doctor)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	expectedDoctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)
	mockRepo.On("FindByID", uint(1)).Return(expectedDoctor, nil)

	doctor, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, doctor)
	assert.Equal(t, expectedDoctor.ID, doctor.ID)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_FindByID_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("doctor not found"))

	doctor, err := mockRepo.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, doctor)
	assert.Equal(t, "doctor not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_FindByUserID_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	expectedDoctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)
	mockRepo.On("FindByUserID", uint(100)).Return(expectedDoctor, nil)

	doctor, err := mockRepo.FindByUserID(100)
	assert.NoError(t, err)
	assert.NotNil(t, doctor)
	assert.Equal(t, expectedDoctor.UserID, doctor.UserID)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_FindByDepartmentID_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	expectedDoctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test", "LIC001", "0812", "test@hospital.com", 1, 1, true)
	mockRepo.On("FindByDepartmentID", uint(1)).Return(expectedDoctor, nil)

	doctor, err := mockRepo.FindByDepartmentID(1)
	assert.NoError(t, err)
	assert.NotNil(t, doctor)
	assert.Equal(t, expectedDoctor.DepartmentID, doctor.DepartmentID)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_IsEmployeeIDExists_True(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	mockRepo.On("IsEmployeeIDExists", "DOC001", mock.Anything).Return(true, nil)

	exists, err := mockRepo.IsEmployeeIDExists("DOC001")
	assert.NoError(t, err)
	assert.True(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_IsEmployeeIDExists_False(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	mockRepo.On("IsEmployeeIDExists", "DOC001", mock.Anything).Return(false, nil)

	exists, err := mockRepo.IsEmployeeIDExists("DOC001")
	assert.NoError(t, err)
	assert.False(t, exists)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	doctor := mocks.NewTestDoctorWithData(1, 100, "DOC001", "Dr. Test Updated", "LIC001", "0812", "test@hospital.com", 1, 1, true)
	mockRepo.On("Update", doctor).Return(nil)

	err := mockRepo.Update(doctor)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_List_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	query := mocks.NewDoctorPaginationQuery(1, 10)
	doctors := mocks.NewTestDoctorList(2)
	mockRepo.On("List", query).Return(doctors, int64(2), nil)

	result, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_DeleteList_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	query := mocks.NewDoctorPaginationQuery(1, 10)
	doctors := mocks.NewTestDoctorList(2)
	mockRepo.On("DeleteList", query).Return(doctors, int64(2), nil)

	result, total, err := mockRepo.DeleteList(query)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_ActiveList_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	query := mocks.NewDoctorPaginationQuery(1, 10)
	doctors := mocks.NewTestDoctorList(2)
	mockRepo.On("ActiveList", query).Return(doctors, int64(2), nil)

	result, total, err := mockRepo.ActiveList(query)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_InactiveList_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	query := mocks.NewDoctorPaginationQuery(1, 10)
	doctors := mocks.NewTestDoctorList(2)
	mockRepo.On("InactiveList", query).Return(doctors, int64(2), nil)

	result, total, err := mockRepo.InactiveList(query)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDoctorRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockDoctorRepository)
	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
