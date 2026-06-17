package typetest

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

// eventTypeTestService membungkus TypeTestService dengan Kafka event publishing.
type eventTypeTestService struct {
	inner     TypeTestService
	publisher kafka.EventPublisher
}

// NewEventTypeTestService mengembalikan TypeTestService dengan event publishing.
func NewEventTypeTestService(inner TypeTestService, publisher kafka.EventPublisher) TypeTestService {
	if publisher == nil {
		return inner
	}
	return &eventTypeTestService{
		inner:     inner,
		publisher: publisher,
	}
}

// ─── Read operations ──────────────────────────────────────────────────────────

func (s *eventTypeTestService) List(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	return s.inner.List(query)
}

func (s *eventTypeTestService) ActiveList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	return s.inner.ActiveList(query)
}

func (s *eventTypeTestService) InactiveList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	return s.inner.InactiveList(query)
}

func (s *eventTypeTestService) DeletedList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}

func (s *eventTypeTestService) FindByID(id uint) (*dto.TypeTestResponse, error) {
	return s.inner.FindByID(id)
}

// ─── Write operations ─────────────────────────────────────────────────────────

func (s *eventTypeTestService) Create(req *dto.CreateTypeTestRequest) (*dto.TypeTestResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicTypeTestCreated,
		events.NewTypeTestCreatedEvent(result.ID, result.Code, result.Name, result.TypeTestCategoryID, result.Description, result.Price, result.IsActive))
	return result, nil
}

func (s *eventTypeTestService) Update(id uint, req *dto.UpdateTypeTestRequest) (*dto.TypeTestResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicTypeTestUpdated,
		events.NewTypeTestUpdatedEvent(result.ID, result.Code, result.Name, result.TypeTestCategoryID, result.Description, result.Price, result.IsActive, "update"))
	return result, nil
}

func (s *eventTypeTestService) Activate(id uint) error {
	if err := s.inner.Activate(id); err != nil {
		return err
	}
	if t, err := s.inner.FindByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicTypeTestUpdated,
			events.NewTypeTestUpdatedEvent(t.ID, t.Code, t.Name, t.TypeTestCategoryID, t.Description, t.Price, true, "activate"))
	}
	return nil
}

func (s *eventTypeTestService) Deactivate(id uint) error {
	if err := s.inner.Deactivate(id); err != nil {
		return err
	}
	if t, err := s.inner.FindByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicTypeTestUpdated,
			events.NewTypeTestUpdatedEvent(t.ID, t.Code, t.Name, t.TypeTestCategoryID, t.Description, t.Price, false, "deactivate"))
	}
	return nil
}

func (s *eventTypeTestService) SoftDelete(id uint) error {
	// Ambil data sebelum dihapus
	t, _ := s.inner.FindByID(id)
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	code, name := "", ""
	if t != nil {
		code = t.Code
		name = t.Name
	}
	s.publisher.PublishAsync(kafka.TopicTypeTestDeleted,
		events.NewTypeTestDeletedEvent(id, code, name, "soft_delete"))
	return nil
}

func (s *eventTypeTestService) HardDelete(id uint) error {
	t, _ := s.inner.FindByID(id)
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	code, name := "", ""
	if t != nil {
		code = t.Code
		name = t.Name
	}
	s.publisher.PublishAsync(kafka.TopicTypeTestDeleted,
		events.NewTypeTestDeletedEvent(id, code, name, "hard_delete"))
	return nil
}

func (s *eventTypeTestService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	if t, err := s.inner.FindByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicTypeTestRestored,
			events.NewTypeTestRestoredEvent(t.ID, t.Code, t.Name))
	}
	return nil
}
