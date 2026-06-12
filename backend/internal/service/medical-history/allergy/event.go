package allergy

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventedAllergyService struct {
	inner     AllergyService
	publisher kafka.EventPublisher
}

func NewEventedAllergyService(inner AllergyService, publisher kafka.EventPublisher) AllergyService {
	if publisher == nil {
		return inner
	}
	return &eventedAllergyService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventedAllergyService) List(query *dto.AllergyPaginationQuery) (*dto.AllergyListResponse, error) {
	return s.inner.List(query)
}

func (s *eventedAllergyService) FindByID(id uint) (*dto.AllergyResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventedAllergyService) Create(req *dto.CreateAllergyRequest) (*dto.AllergyResponse, error) {
	resp, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicAllergyCreated,
		events.NewAllergyCreatedEvent(resp.ID, resp.PatientID, resp.AllergenType, resp.AllergenName, resp.Severity),
	)
	return resp, nil
}

func (s *eventedAllergyService) Update(id uint, req *dto.UpdateAllergyRequest) (*dto.AllergyResponse, error) {
	resp, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicAllergyUpdated,
		events.NewAllergyUpdatedEvent(resp.ID, resp.PatientID),
	)
	return resp, nil
}

func (s *eventedAllergyService) Delete(id uint) error {
	if err := s.inner.Delete(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicAllergyDeleted,
		events.NewAllergyDeletedEvent(id),
	)
	return nil
}
