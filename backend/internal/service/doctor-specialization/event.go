package doctorspecialization

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventDoctorSpecializationService struct {
	inner     DoctorSpecializationService
	publisher kafka.EventPublisher
}

func NewEventDoctorSpecializationService(inner DoctorSpecializationService, publisher kafka.EventPublisher) DoctorSpecializationService {
	if publisher == nil {
		return inner
	}
	return &eventDoctorSpecializationService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventDoctorSpecializationService) List(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error) {
	return s.inner.List(query)
}
func (s *eventDoctorSpecializationService) ActiveList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error) {
	return s.inner.ActiveList(query)
}
func (s *eventDoctorSpecializationService) InactiveList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error) {
	return s.inner.InactiveList(query)
}
func (s *eventDoctorSpecializationService) DeletedList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}
func (s *eventDoctorSpecializationService) FindByID(id uint) (*dto.DoctorSpecializationResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventDoctorSpecializationService) Create(req *dto.CreateDoctorSpecializationRequest) (*dto.DoctorSpecializationResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicDoctorSpecializationCreated,
		events.NewDoctorSpecializationCreatedEvent(result.ID, result.Name, result.Code, result.Description, result.IsActive),
	)
	return result, nil
}
func (s *eventDoctorSpecializationService) Update(id uint, req *dto.UpdateDoctorSpecializationRequest) (*dto.DoctorSpecializationResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(
		kafka.TopicDoctorSpecializationUpdated,
		events.NewDoctorSpecializationUpdatedEvent(result.ID, result.Name, result.Code, result.Description, result.IsActive, "update"),
	)
	return result, nil
}
func (s *eventDoctorSpecializationService) SoftDelete(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	name := ""
	code := ""
	if ds != nil {
		name = ds.Name
		code = ds.Code
	}

	s.publisher.PublishAsync(
		kafka.TopicDoctorSpecializationDeleted,
		events.NewDoctorSpecializationDeletedEvent(id, name, code, "soft_delete"),
	)
	return nil
}
func (s *eventDoctorSpecializationService) Restore(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	name := ""
	code := ""
	if ds != nil {
		name = ds.Name
		code = ds.Code
	}

	s.publisher.PublishAsync(
		kafka.TopicDoctorSpecializationRestored,
		events.NewDoctorSpecializationRestoredEvent(id, name, code),
	)
	return nil
}
func (s *eventDoctorSpecializationService) HardDelete(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	name := ""
	code := ""
	if ds != nil {
		name = ds.Name
		code = ds.Code
	}

	s.publisher.PublishAsync(
		kafka.TopicDoctorSpecializationDeleted,
		events.NewDoctorSpecializationDeletedEvent(id, name, code, "hard_delete"),
	)
	return nil
}

func (s *eventDoctorSpecializationService) Activate(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.Activate(id); err != nil {
		return err
	}
	name := ""
	code := ""
	if ds != nil {
		name = ds.Name
		code = ds.Code
	}
	s.publisher.PublishAsync(
		kafka.TopicDoctorSpecializationUpdated,
		events.NewDoctorSpecializationUpdatedEvent(id, name, code, ds.Description, true, "activate"),
	)
	return nil
}

func (s *eventDoctorSpecializationService) Deactivate(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.Deactivate(id); err != nil {
		return err
	}
	name := ""
	code := ""
	if ds != nil {
		name = ds.Name
		code = ds.Code
	}
	s.publisher.PublishAsync(
		kafka.TopicDoctorSpecializationUpdated,
		events.NewDoctorSpecializationUpdatedEvent(id, name, code, ds.Description, false, "deactivate"),
	)
	return nil
}
