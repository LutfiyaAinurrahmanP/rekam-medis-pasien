package medicalcondition

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventedMedicalConditionService struct {
	inner     MedicalConditionService
	publisher kafka.EventPublisher
}

func NewEventedMedicalConditionService(inner MedicalConditionService, publisher kafka.EventPublisher) MedicalConditionService {
	if publisher == nil {
		return inner
	}
	return &eventedMedicalConditionService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventedMedicalConditionService) List(query *dto.MedicalConditionPaginationQuery) (*dto.MedicalConditionListResponse, error) {
	return s.inner.List(query)
}

func (s *eventedMedicalConditionService) FindByID(id uint) (*dto.MedicalConditionResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventedMedicalConditionService) Create(req *dto.CreateMedicalConditionRequest) (*dto.MedicalConditionResponse, error) {
	resp, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicalConditionCreated,
		events.NewMedicalConditionCreatedEvent(resp.ID, resp.PatientID, resp.ConditionName, resp.Status),
	)
	return resp, nil
}

func (s *eventedMedicalConditionService) Update(id uint, req *dto.UpdateMedicalConditionRequest) (*dto.MedicalConditionResponse, error) {
	resp, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicalConditionUpdated,
		events.NewMedicalConditionUpdatedEvent(resp.ID, resp.PatientID, resp.Status),
	)
	return resp, nil
}

func (s *eventedMedicalConditionService) Delete(id uint) error {
	if err := s.inner.Delete(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicalConditionDeleted,
		events.NewMedicalConditionDeletedEvent(id),
	)
	return nil
}
