package repository

import (
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppointmentRepository_List(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)
	query := mocks.NewAppointmentPaginationQuery(1, 10)
	appointments := mocks.NewTestAppointmentList(2)

	mockRepo.On("List", query).Return(appointments, int64(2), nil)

	res, total, err := mockRepo.List(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_DeletedList(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)
	query := mocks.NewAppointmentPaginationQuery(1, 10)
	appointments := mocks.NewTestAppointmentList(2)

	mockRepo.On("DeletedList", query).Return(appointments, int64(2), nil)

	res, total, err := mockRepo.DeletedList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_UpcomingList(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)
	query := mocks.NewAppointmentPaginationQuery(1, 10)
	appointments := mocks.NewTestAppointmentList(2)

	mockRepo.On("UpcomingList", query).Return(appointments, int64(2), nil)

	res, total, err := mockRepo.UpcomingList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_PastList(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)
	query := mocks.NewAppointmentPaginationQuery(1, 10)
	appointments := mocks.NewTestAppointmentList(2)

	mockRepo.On("PastList", query).Return(appointments, int64(2), nil)

	res, total, err := mockRepo.PastList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_TodayList(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)
	query := mocks.NewAppointmentPaginationQuery(1, 10)
	appointments := mocks.NewTestAppointmentList(2)

	mockRepo.On("TodayList", query).Return(appointments, int64(2), nil)

	res, total, err := mockRepo.TodayList(query)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, int64(2), total)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_FindByID_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)
	expectedAppt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false)

	mockRepo.On("FindByID", uint(1)).Return(expectedAppt, nil)

	res, err := mockRepo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedAppt.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)
	expectedAppt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(expectedAppt, nil)

	res, err := mockRepo.FindByIDUnscoped(1)
	assert.NoError(t, err)
	assert.Equal(t, expectedAppt.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_Create_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("Create", mock.AnythingOfType("*models.Appointment")).Return(nil)

	err := mockRepo.Create(&models.Appointment{PatientID: 1, DoctorID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_Update_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("Update", mock.AnythingOfType("*models.Appointment")).Return(nil)

	err := mockRepo.Update(&models.Appointment{ID: 1, PatientID: 1, DoctorID: 1})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_Confirm_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("Confirm", uint(1)).Return(nil)

	err := mockRepo.Confirm(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_Start_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("Start", uint(1)).Return(nil)

	err := mockRepo.Start(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_Complete_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("Complete", uint(1)).Return(nil)

	err := mockRepo.Complete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_Cancel_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("Cancel", uint(1), "Patient requested").Return(nil)

	err := mockRepo.Cancel(1, "Patient requested")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_Reschedule_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("Reschedule", uint(1), "2023-12-02", "11:00").Return(nil)

	err := mockRepo.Reschedule(1, "2023-12-02", "11:00")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_NoShow_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("NoShow", uint(1)).Return(nil)

	err := mockRepo.NoShow(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_SoftDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := mockRepo.SoftDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_Restore_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("Restore", uint(1)).Return(nil)

	err := mockRepo.Restore(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentRepository_HardDelete_Success(t *testing.T) {
	mockRepo := new(mocks.MockAppointmentRepository)

	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := mockRepo.HardDelete(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
