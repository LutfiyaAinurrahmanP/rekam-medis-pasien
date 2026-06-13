package service

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	appointment "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/appointment"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestAppointmentService() (*mocks.MockAppointmentRepository, appointment.AppointmentService) {
	mockRepo := new(mocks.MockAppointmentRepository)
	cfg := &config.Config{
		Pagination: config.PaginationConfig{
			DefaultPageSize: 10,
			MaxPageSize:     100,
		},
	}
	service := appointment.NewAppointmentService(mockRepo, cfg)
	return mockRepo, service
}

func TestAppointmentService_List_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	query := mocks.NewAppointmentPaginationQuery(1, 10)
	query.SortBy = "appointment_date"
	query.SortDir = "asc"
	appointments := mocks.NewTestAppointmentList(2)

	mockRepo.On("List", query).Return(appointments, int64(2), nil)

	res, err := service.List(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_DeletedList_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	query := mocks.NewAppointmentPaginationQuery(1, 10)
	query.SortBy = "appointment_date"
	query.SortDir = "desc"
	appointments := mocks.NewTestAppointmentList(2)

	mockRepo.On("DeletedList", query).Return(appointments, int64(2), nil)

	res, err := service.DeletedList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_UpcomingList_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	query := mocks.NewAppointmentPaginationQuery(1, 10)
	query.SortBy = "appointment_date"
	query.SortDir = "asc"
	appointments := mocks.NewTestAppointmentList(2)

	mockRepo.On("UpcomingList", query).Return(appointments, int64(2), nil)

	res, err := service.UpcomingList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_PastList_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	query := mocks.NewAppointmentPaginationQuery(1, 10)
	query.SortBy = "appointment_date"
	query.SortDir = "desc"
	appointments := mocks.NewTestAppointmentList(2)

	mockRepo.On("PastList", query).Return(appointments, int64(2), nil)

	res, err := service.PastList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_TodayList_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	query := mocks.NewAppointmentPaginationQuery(1, 10)
	query.SortBy = "appointment_date"
	query.SortDir = "asc"
	appointments := mocks.NewTestAppointmentList(2)

	mockRepo.On("TodayList", query).Return(appointments, int64(2), nil)

	res, err := service.TodayList(query)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res.Data, 2)
	assert.Equal(t, int64(2), res.Meta.TotalItems)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_FindByID_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false)

	mockRepo.On("FindByID", uint(1)).Return(appt, nil)

	res, err := service.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, appt.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_FindByID_NotFound(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()

	mockRepo.On("FindByID", uint(99)).Return(nil, errors.New("appointment not found"))

	res, err := service.FindByID(99)

	assert.Error(t, err)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_FindByIDUnscoped_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(appt, nil)

	res, err := service.FindByIDUnscoped(1)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, appt.ID, res.ID)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_Create_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	req := mocks.NewCreateAppointmentRequest(1, 1, "2023-12-01", "10:00", 30, "Reason", "Notes")

	mockRepo.On("Create", mock.AnythingOfType("*models.Appointment")).Return(nil)
	mockRepo.On("FindByID", mock.Anything).Return(mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false), nil)

	res, err := service.Create(req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.PatientID, res.PatientID)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_Update_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false)
	req := mocks.NewUpdateAppointmentRequest("2023-12-02", "11:00", 45, "New Reason", "New Notes")

	mockRepo.On("FindByID", uint(1)).Return(appt, nil)
	mockRepo.On("Update", mock.AnythingOfType("*models.Appointment")).Return(nil)

	res, err := service.Update(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "2023-12-02", res.AppointmentDate)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_Confirm_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false)

	mockRepo.On("FindByID", uint(1)).Return(appt, nil)
	mockRepo.On("Confirm", uint(1)).Return(nil)

	err := service.Confirm(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_Start_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "confirmed", "Reason", "Notes", 30, false)

	mockRepo.On("FindByID", uint(1)).Return(appt, nil)
	mockRepo.On("Start", uint(1)).Return(nil)

	err := service.Start(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_Complete_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "in_progress", "Reason", "Notes", 30, false)

	mockRepo.On("FindByID", uint(1)).Return(appt, nil)
	mockRepo.On("Complete", uint(1)).Return(nil)

	err := service.Complete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_Cancel_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false)
	req := &dto.CancelAppointmentRequest{Reason: "Patient cancelled"}

	mockRepo.On("FindByID", uint(1)).Return(appt, nil)
	mockRepo.On("Cancel", uint(1), "Patient cancelled").Return(nil)

	err := service.Cancel(1, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_Reschedule_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false)
	req := &dto.RescheduleAppointmentRequest{AppointmentDate: "2023-12-05", AppointmentTime: "14:00", Reason: "Doctor unavailable"}

	mockRepo.On("FindByID", uint(1)).Return(appt, nil)
	mockRepo.On("Reschedule", uint(1), "2023-12-05", "14:00").Return(nil)

	err := service.Reschedule(1, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_NoShow_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false)

	mockRepo.On("FindByID", uint(1)).Return(appt, nil)
	mockRepo.On("NoShow", uint(1)).Return(nil)

	err := service.NoShow(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_SoftDelete_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, false)

	mockRepo.On("FindByID", uint(1)).Return(appt, nil)
	mockRepo.On("SoftDelete", uint(1)).Return(nil)

	err := service.SoftDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_Restore_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(appt, nil)
	mockRepo.On("Restore", uint(1)).Return(nil)

	err := service.Restore(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestAppointmentService_HardDelete_Success(t *testing.T) {
	mockRepo, service := setupTestAppointmentService()
	appt := mocks.NewTestAppointmentWithData(1, 1, 1, "2023-12-01", "10:00", "scheduled", "Reason", "Notes", 30, true)

	mockRepo.On("FindByIDUnscoped", uint(1)).Return(appt, nil)
	mockRepo.On("HardDelete", uint(1)).Return(nil)

	err := service.HardDelete(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
