package surgicalhistory

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventedSurgicalHistoryService struct {
	inner     SurgicalHistoryService
	publisher kafka.EventPublisher
}

func NewEventedSurgicalHistoryService(inner SurgicalHistoryService, publisher kafka.EventPublisher) SurgicalHistoryService {
	if publisher == nil {
		return inner
	}
	return &eventedSurgicalHistoryService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventedSurgicalHistoryService) List(query *dto.SurgicalHistoryPaginationQuery) (*dto.SurgicalHistoryListResponse, error) {
	return s.inner.List(query)
}

func (s *eventedSurgicalHistoryService) FindByID(id uint) (*dto.SurgicalHistoryResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventedSurgicalHistoryService) Create(req *dto.CreateSurgicalHistoryRequest) (*dto.SurgicalHistoryResponse, error) {
	resp, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicSurgicalHistoryCreated,
		events.NewSurgicalHistoryCreatedEvent(resp.ID, resp.PatientID, resp.ProcedureName, resp.SurgeryDate),
	)
	return resp, nil
}

func (s *eventedSurgicalHistoryService) Update(id uint, req *dto.UpdateSurgicalHistoryRequest) (*dto.SurgicalHistoryResponse, error) {
	resp, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicSurgicalHistoryUpdated,
		events.NewSurgicalHistoryUpdatedEvent(resp.ID, resp.PatientID),
	)
	return resp, nil
}

func (s *eventedSurgicalHistoryService) Delete(id uint) error {
	if err := s.inner.Delete(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicSurgicalHistoryDeleted,
		events.NewSurgicalHistoryDeletedEvent(id),
	)
	return nil
}
