package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockAppointmentService struct {
	mock.Mock
}

func (m *MockAppointmentService) List(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AppointmentListResponse), args.Error(1)
}

func (m *MockAppointmentService) DeletedList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AppointmentListResponse), args.Error(1)
}

func (m *MockAppointmentService) UpcomingList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AppointmentListResponse), args.Error(1)
}

func (m *MockAppointmentService) PastList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AppointmentListResponse), args.Error(1)
}

func (m *MockAppointmentService) TodayList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AppointmentListResponse), args.Error(1)
}

func (m *MockAppointmentService) FindByID(id uint) (*dto.AppointmentResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AppointmentResponse), args.Error(1)
}

func (m *MockAppointmentService) FindByIDUnscoped(id uint) (*dto.AppointmentResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AppointmentResponse), args.Error(1)
}

func (m *MockAppointmentService) Create(req *dto.CreateAppointmentRequest) (*dto.AppointmentResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AppointmentResponse), args.Error(1)
}

func (m *MockAppointmentService) Update(id uint, req *dto.UpdateAppointmentRequest) (*dto.AppointmentResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AppointmentResponse), args.Error(1)
}

func (m *MockAppointmentService) Confirm(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentService) Start(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentService) Complete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentService) Cancel(id uint, req *dto.CancelAppointmentRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockAppointmentService) Reschedule(id uint, req *dto.RescheduleAppointmentRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockAppointmentService) NoShow(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentService) SoftDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentService) Restore(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAppointmentService) HardDelete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
