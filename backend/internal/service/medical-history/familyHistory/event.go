package familyhistory

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventedFamilyHistoryService struct {
	inner     FamilyHistoryService
	publisher kafka.EventPublisher
}

func NewEventedFamilyHistoryService(inner FamilyHistoryService, publisher kafka.EventPublisher) FamilyHistoryService {
	if publisher == nil {
		return inner
	}
	return &eventedFamilyHistoryService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventedFamilyHistoryService) List(query *dto.FamilyHistoryPaginationQuery) (*dto.FamilyHistoryListResponse, error) {
	return s.inner.List(query)
}

func (s *eventedFamilyHistoryService) FindByID(id uint) (*dto.FamilyHistoryResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventedFamilyHistoryService) Create(req *dto.CreateFamilyHistoryRequest) (*dto.FamilyHistoryResponse, error) {
	resp, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicFamilyHistoryCreated,
		events.NewFamilyHistoryCreatedEvent(resp.ID, resp.PatientID, resp.FamilyMember, resp.ConditionName),
	)
	return resp, nil
}

func (s *eventedFamilyHistoryService) Update(id uint, req *dto.UpdateFamilyHistoryRequest) (*dto.FamilyHistoryResponse, error) {
	resp, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicFamilyHistoryUpdated,
		events.NewFamilyHistoryUpdatedEvent(resp.ID, resp.PatientID),
	)
	return resp, nil
}

func (s *eventedFamilyHistoryService) Delete(id uint) error {
	if err := s.inner.Delete(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicFamilyHistoryDeleted,
		events.NewFamilyHistoryDeletedEvent(id),
	)
	return nil
}
