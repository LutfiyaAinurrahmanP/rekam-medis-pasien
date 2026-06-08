package appointment

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventAppointmentService struct {
	inner     AppointmentService
	publisher kafka.EventPublisher
}

func NewEventAppointmentService(inner AppointmentService, publisher kafka.EventPublisher) AppointmentService {
	if publisher == nil {
		return inner
	}
	return &eventAppointmentService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventAppointmentService) List(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	return s.inner.List(query)
}

func (s *eventAppointmentService) DeletedList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	return s.inner.DeletedList(query)
}

func (s *eventAppointmentService) UpcomingList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	return s.inner.UpcomingList(query)
}

func (s *eventAppointmentService) PastList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	return s.inner.PastList(query)
}

func (s *eventAppointmentService) TodayList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	return s.inner.TodayList(query)
}

func (s *eventAppointmentService) FindByID(id uint) (*dto.AppointmentResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventAppointmentService) FindByIDUnscoped(id uint) (*dto.AppointmentResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *eventAppointmentService) Create(req *dto.CreateAppointmentRequest) (*dto.AppointmentResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicAppointmentCreated,
		events.NewAppointmentCreatedEvent(
			result.ID, result.PatientID, result.DoctorID, result.AppointmentDate, result.AppointmentTime, result.DurationMinutes, result.Status,
		),
	)
	return result, nil
}

func (s *eventAppointmentService) Update(id uint, req *dto.UpdateAppointmentRequest) (*dto.AppointmentResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicAppointmentUpdated,
		events.NewAppointmentUpdatedEvent(
			result.ID, result.PatientID, result.DoctorID, result.AppointmentDate, result.AppointmentTime, result.DurationMinutes, result.Status, "update",
		),
	)
	return result, nil
}

func (s *eventAppointmentService) Confirm(id uint) error {
	if err := s.inner.Confirm(id); err != nil {
		return err
	}
	s.publisher.PublishAsync(kafka.TopicAppointmentConfirmed, events.NewAppointmentConfirmedEvent(id))
	return nil
}

func (s *eventAppointmentService) Start(id uint) error {
	if err := s.inner.Start(id); err != nil {
		return err
	}
	s.publisher.PublishAsync(kafka.TopicAppointmentStarted, events.NewAppointmentStartedEvent(id))
	return nil
}

func (s *eventAppointmentService) Complete(id uint) error {
	if err := s.inner.Complete(id); err != nil {
		return err
	}
	s.publisher.PublishAsync(kafka.TopicAppointmentCompleted, events.NewAppointmentCompletedEvent(id))
	return nil
}

func (s *eventAppointmentService) Cancel(id uint, req *dto.CancelAppointmentRequest) error {
	if err := s.inner.Cancel(id, req); err != nil {
		return err
	}
	s.publisher.PublishAsync(kafka.TopicAppointmentCancelled, events.NewAppointmentCancelledEvent(id, req.Reason))
	return nil
}

func (s *eventAppointmentService) Reschedule(id uint, req *dto.RescheduleAppointmentRequest) error {
	if err := s.inner.Reschedule(id, req); err != nil {
		return err
	}
	s.publisher.PublishAsync(kafka.TopicAppointmentRescheduled, events.NewAppointmentRescheduledEvent(id, req.AppointmentDate, req.AppointmentTime))
	return nil
}

func (s *eventAppointmentService) NoShow(id uint) error {
	if err := s.inner.NoShow(id); err != nil {
		return err
	}
	s.publisher.PublishAsync(kafka.TopicAppointmentNoShow, events.NewAppointmentNoShowEvent(id))
	return nil
}

func (s *eventAppointmentService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.publisher.PublishAsync(
		kafka.TopicAppointmentDeleted,
		events.NewAppointmentDeletedEvent(id, "soft_delete"),
	)
	return nil
}

func (s *eventAppointmentService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.publisher.PublishAsync(
		kafka.TopicAppointmentRestored,
		events.NewAppointmentRestoredEvent(id),
	)
	return nil
}

func (s *eventAppointmentService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.publisher.PublishAsync(
		kafka.TopicAppointmentDeleted,
		events.NewAppointmentDeletedEvent(id, "hard_delete"),
	)
	return nil
}
