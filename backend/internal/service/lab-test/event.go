package labtest

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventedLabTestService struct {
	inner     LabTestService
	publisher kafka.EventPublisher
}

func NewEventedLabTestService(inner LabTestService, publisher kafka.EventPublisher) LabTestService {
	if publisher == nil {
		return inner
	}
	return &eventedLabTestService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventedLabTestService) List(query *dto.LabTestPaginationQuery) (*dto.LabTestListResponse, error) {
	return s.inner.List(query)
}

func (s *eventedLabTestService) DeletedList(query *dto.LabTestPaginationQuery) (*dto.LabTestDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}

func (s *eventedLabTestService) FindByID(id uint) (*dto.LabTestResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventedLabTestService) FindByIDUnscoped(id uint) (*dto.LabTestResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *eventedLabTestService) Create(req *dto.CreateLabTestRequest) (*dto.LabTestResponse, error) {
	resp, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicLabTestCreated,
		events.NewLabTestCreatedEvent(
			resp.ID,
			resp.MedicalRecordID,
			resp.TestTypeID,
			resp.OrderedByDoctorID,
			resp.OrderDate,
			resp.Status,
		),
	)
	return resp, nil
}

func (s *eventedLabTestService) Update(id uint, req *dto.UpdateLabTestRequest) (*dto.LabTestResponse, error) {
	resp, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicLabTestUpdated,
		events.NewLabTestUpdatedEvent(resp.ID, "update"),
	)
	return resp, nil
}

func (s *eventedLabTestService) CollectSample(id uint) (*dto.LabTestResponse, error) {
	resp, err := s.inner.CollectSample(id)
	if err != nil {
		return nil, err
	}

	date := ""
	if resp.SampleCollectionDate != nil {
		date = *resp.SampleCollectionDate
	}
	s.publisher.PublishAsync(
		kafka.TopicLabTestSampleCollected,
		events.NewLabTestSampleCollectedEvent(resp.ID, date),
	)
	return resp, nil
}

func (s *eventedLabTestService) Start(id uint) (*dto.LabTestResponse, error) {
	resp, err := s.inner.Start(id)
	if err != nil {
		return nil, err
	}

	date := ""
	if resp.TestStartDate != nil {
		date = *resp.TestStartDate
	}
	s.publisher.PublishAsync(
		kafka.TopicLabTestStarted,
		events.NewLabTestStartedEvent(resp.ID, date),
	)
	return resp, nil
}

func (s *eventedLabTestService) Complete(id uint, req *dto.CompleteLabTestRequest) (*dto.LabTestResponse, error) {
	resp, err := s.inner.Complete(id, req)
	if err != nil {
		return nil, err
	}

	date := ""
	if resp.ResultDate != nil {
		date = *resp.ResultDate
	}
	s.publisher.PublishAsync(
		kafka.TopicLabTestCompleted,
		events.NewLabTestCompletedEvent(resp.ID, date),
	)
	return resp, nil
}

func (s *eventedLabTestService) Cancel(id uint) (*dto.LabTestResponse, error) {
	resp, err := s.inner.Cancel(id)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicLabTestCancelled,
		events.NewLabTestCancelledEvent(resp.ID),
	)
	return resp, nil
}

func (s *eventedLabTestService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicLabTestDeleted,
		events.NewLabTestDeletedEvent(id, "soft_delete"),
	)
	return nil
}

func (s *eventedLabTestService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicLabTestRestored,
		events.NewLabTestRestoredEvent(id),
	)
	return nil
}

func (s *eventedLabTestService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicLabTestDeleted,
		events.NewLabTestDeletedEvent(id, "hard_delete"),
	)
	return nil
}
