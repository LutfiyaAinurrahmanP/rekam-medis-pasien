package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAppointmentRepository struct {
	mock.Mock
}

func (m *MockAppointmentRepository) List(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Appointment), args.Get(1).(int64), args.Error(2)
}

func (m *MockAppointmentRepository) DeletedList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Appointment), args.Get(1).(int64), args.Error(2)
}

func (m *MockAppointmentRepository) UpcomingList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Appointment), args.Get(1).(int64), args.Error(2)
}

func (m *MockAppointmentRepository) PastList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Appointment), args.Get(1).(int64), args.Error(2)
}

func (m *MockAppointmentRepository) TodayList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]models.Appointment), args.Get(1).(int64), args.Error(2)
}

func (m *MockAppointmentRepository) FindByID(id uint) (*models.Appointment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) FindByIDUnscoped(id uint) (*models.Appointment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Appointment), args.Error(1)
}

func (m *MockAppointmentRepository) Create(appointment *models.Appointment) error {
	args := m.Called(appointment)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Update(appointment *models.Appointment) error {
	args := m.Called(appointment)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Confirm(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Start(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Complete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Cancel(id uint, reason string) error {
	args := m.Called(id, reason)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Reschedule(id uint, newDate, newTime string) error {
	args := m.Called(id, newDate, newTime)
	return args.Error(0)
}

func (m *MockAppointmentRepository) NoShow(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentRepository) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentRepository) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentRepository) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
