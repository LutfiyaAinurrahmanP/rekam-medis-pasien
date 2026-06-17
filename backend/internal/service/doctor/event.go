package doctor

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

// eventDoctorService membungkus DoctorService dengan Kafka event publishing.
type eventDoctorService struct {
	inner     DoctorService
	publisher kafka.EventPublisher
}

// NewEventDoctorService mengembalikan DoctorService dengan event publishing.
func NewEventDoctorService(inner DoctorService, publisher kafka.EventPublisher) DoctorService {
	if publisher == nil {
		return inner
	}
	return &eventDoctorService{inner: inner, publisher: publisher}
}

// helper: ekstrak payload dari DoctorResponse
func doctorIsActive(r *dto.DoctorResponse) bool {
	if r.IsActive == nil {
		return false
	}
	return *r.IsActive
}

// ─── Read operations ──────────────────────────────────────────────────────────

func (s *eventDoctorService) GetMyDoctorData(userID uint) (*dto.DoctorResponse, error) {
	return s.inner.GetMyDoctorData(userID)
}

func (s *eventDoctorService) ListDoctors(query *dto.DoctorPaginationQuery) (*dto.DoctorListResponse, error) {
	return s.inner.ListDoctors(query)
}

func (s *eventDoctorService) DeletedListDoctors(query *dto.DoctorPaginationQuery) (*dto.DoctorDeletedListResponse, error) {
	return s.inner.DeletedListDoctors(query)
}

func (s *eventDoctorService) GetDoctorByID(id uint) (*dto.DoctorResponse, error) {
	return s.inner.GetDoctorByID(id)
}

func (s *eventDoctorService) ActiveList(query *dto.DoctorPaginationQuery) (*dto.DoctorListResponse, error) {
	return s.inner.ActiveList(query)
}

func (s *eventDoctorService) InactiveList(query *dto.DoctorPaginationQuery) (*dto.DoctorListResponse, error) {
	return s.inner.InactiveList(query)
}

// ─── Write operations ─────────────────────────────────────────────────────────

func (s *eventDoctorService) CreateDoctor(req *dto.CreateDoctorRequest) (*dto.DoctorResponse, error) {
	result, err := s.inner.CreateDoctor(req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicDoctorCreated,
		events.NewDoctorCreatedEvent(
			result.ID, result.FullName, result.SpecializationID,
			result.Phone, result.Email, result.DepartmentID, doctorIsActive(result),
		))
	return result, nil
}

func (s *eventDoctorService) UpdateDoctor(id uint, req *dto.UpdateDoctorRequest) (*dto.DoctorResponse, error) {
	result, err := s.inner.UpdateDoctor(id, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicDoctorUpdated,
		events.NewDoctorUpdatedEvent(
			result.ID, result.FullName, result.SpecializationID,
			result.Phone, result.Email, result.DepartmentID, doctorIsActive(result), "admin_update",
		))
	return result, nil
}

func (s *eventDoctorService) UpdateMyDoctorData(userID uint, req *dto.UpdateDoctorRequest) (*dto.DoctorResponse, error) {
	result, err := s.inner.UpdateMyDoctorData(userID, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicDoctorUpdated,
		events.NewDoctorUpdatedEvent(
			result.ID, result.FullName, result.SpecializationID,
			result.Phone, result.Email, result.DepartmentID, doctorIsActive(result), "self_update",
		))
	return result, nil
}

func (s *eventDoctorService) ActivateDoctor(id uint) (*dto.DoctorResponse, error) {
	result, err := s.inner.ActivateDoctor(id)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicDoctorUpdated,
		events.NewDoctorUpdatedEvent(
			result.ID, result.FullName, result.SpecializationID,
			result.Phone, result.Email, result.DepartmentID, true, "activate",
		))
	return result, nil
}

func (s *eventDoctorService) DeactivateDoctor(id uint) (*dto.DoctorResponse, error) {
	result, err := s.inner.DeactivateDoctor(id)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicDoctorUpdated,
		events.NewDoctorUpdatedEvent(
			result.ID, result.FullName, result.SpecializationID,
			result.Phone, result.Email, result.DepartmentID, false, "deactivate",
		))
	return result, nil
}

func (s *eventDoctorService) SoftDeleteDoctor(id uint) error {
	doctor, _ := s.inner.GetDoctorByID(id)
	if err := s.inner.SoftDeleteDoctor(id); err != nil {
		return err
	}
	name := ""
	if doctor != nil {
		name = doctor.FullName
	}
	s.publisher.PublishAsync(kafka.TopicDoctorDeleted,
		events.NewDoctorDeletedEvent(id, name, "soft_delete"))
	return nil
}

func (s *eventDoctorService) HardDeleteDoctor(id uint) error {
	doctor, _ := s.inner.GetDoctorByID(id)
	if err := s.inner.HardDeleteDoctor(id); err != nil {
		return err
	}
	name := ""
	if doctor != nil {
		name = doctor.FullName
	}
	s.publisher.PublishAsync(kafka.TopicDoctorDeleted,
		events.NewDoctorDeletedEvent(id, name, "hard_delete"))
	return nil
}

func (s *eventDoctorService) RestoreDoctor(id uint) error {
	if err := s.inner.RestoreDoctor(id); err != nil {
		return err
	}
	if doctor, err := s.inner.GetDoctorByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicDoctorRestored,
			events.NewDoctorRestoredEvent(doctor.ID, doctor.FullName))
	}
	return nil
}
