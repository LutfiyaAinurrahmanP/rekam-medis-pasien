package vitalsign

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventedVitalSignService struct {
	inner     VitalSignService
	publisher kafka.EventPublisher
}

func NewEventedVitalSignService(inner VitalSignService, publisher kafka.EventPublisher) VitalSignService {
	if publisher == nil {
		return inner
	}
	return &eventedVitalSignService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventedVitalSignService) List(query *dto.VitalSignPaginationQuery) (*dto.VitalSignListResponse, error) {
	return s.inner.List(query)
}

func (s *eventedVitalSignService) DeletedList(query *dto.VitalSignPaginationQuery) (*dto.VitalSignDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}

func (s *eventedVitalSignService) FindByID(id uint) (*dto.VitalSignResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventedVitalSignService) FindByIDUnscoped(id uint) (*dto.VitalSignResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *eventedVitalSignService) Create(req *dto.CreateVitalSignRequest) (*dto.VitalSignResponse, error) {
	resp, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicVitalSignCreated,
		events.NewVitalSignCreatedEvent(resp.ID, resp.MedicalRecordID),
	)
	return resp, nil
}

func (s *eventedVitalSignService) Update(id uint, req *dto.UpdateVitalSignRequest) (*dto.VitalSignResponse, error) {
	resp, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicVitalSignUpdated,
		events.NewVitalSignUpdatedEvent(resp.ID, resp.MedicalRecordID),
	)
	return resp, nil
}

func (s *eventedVitalSignService) SoftDelete(id uint) error {
	vitalSign, err := s.inner.FindByID(id)
	if err == nil && vitalSign != nil {
		errDelete := s.inner.SoftDelete(id)
		if errDelete == nil {
			s.publisher.PublishAsync(
				kafka.TopicVitalSignDeleted,
				events.NewVitalSignDeletedEvent(id, vitalSign.MedicalRecordID, "soft_delete"),
			)
		}
		return errDelete
	}
	return s.inner.SoftDelete(id)
}

func (s *eventedVitalSignService) Restore(id uint) error {
	err := s.inner.Restore(id)
	if err == nil {
		s.publisher.PublishAsync(
			kafka.TopicVitalSignRestored,
			events.NewVitalSignRestoredEvent(id, 0),
		)
	}
	return err
}

func (s *eventedVitalSignService) HardDelete(id uint) error {
	return s.inner.HardDelete(id)
}
