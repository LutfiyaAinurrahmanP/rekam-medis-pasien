package department

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

// eventDepartmentService membungkus DepartmentService dengan Kafka event publishing.
type eventDepartmentService struct {
	inner     DepartmentService
	publisher kafka.EventPublisher
}

// NewEventDepartmentService mengembalikan DepartmentService dengan event publishing.
func NewEventDepartmentService(inner DepartmentService, publisher kafka.EventPublisher) DepartmentService {
	if publisher == nil {
		return inner
	}
	return &eventDepartmentService{inner: inner, publisher: publisher}
}

// ─── Read operations ──────────────────────────────────────────────────────────

func (s *eventDepartmentService) ListDepartments(query *dto.DepartmentPaginationQuery) (*dto.DepartmentListResponse, error) {
	return s.inner.ListDepartments(query)
}

func (s *eventDepartmentService) DeleteListDepartments(query *dto.DepartmentPaginationQuery) (*dto.DepartmentDeletedListResponse, error) {
	return s.inner.DeleteListDepartments(query)
}

func (s *eventDepartmentService) GetDepartmentByID(id uint) (*dto.DepartmentResponse, error) {
	return s.inner.GetDepartmentByID(id)
}

// ─── Write operations ─────────────────────────────────────────────────────────

func (s *eventDepartmentService) CreateDepartment(req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	result, err := s.inner.CreateDepartment(req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicDepartmentCreated,
		events.NewDepartmentCreatedEvent(result.ID, result.Name, result.Description, result.Code))
	return result, nil
}

func (s *eventDepartmentService) UpdateDepartment(id uint, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error) {
	result, err := s.inner.UpdateDepartment(id, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicDepartmentUpdated,
		events.NewDepartmentUpdatedEvent(result.ID, result.Name, result.Description, result.Code))
	return result, nil
}

func (s *eventDepartmentService) SoftDeleteDepartment(id uint) error {
	dept, _ := s.inner.GetDepartmentByID(id)
	if err := s.inner.SoftDeleteDepartment(id); err != nil {
		return err
	}
	name := ""
	if dept != nil {
		name = dept.Name
	}
	s.publisher.PublishAsync(kafka.TopicDepartmentDeleted,
		events.NewDepartmentDeletedEvent(id, name, "soft_delete"))
	return nil
}

func (s *eventDepartmentService) HardDeleteDepartment(id uint) error {
	dept, _ := s.inner.GetDepartmentByID(id)
	if err := s.inner.HardDeleteDepartment(id); err != nil {
		return err
	}
	name := ""
	if dept != nil {
		name = dept.Name
	}
	s.publisher.PublishAsync(kafka.TopicDepartmentDeleted,
		events.NewDepartmentDeletedEvent(id, name, "hard_delete"))
	return nil
}

func (s *eventDepartmentService) RestoreDepartment(id uint) error {
	if err := s.inner.RestoreDepartment(id); err != nil {
		return err
	}
	if dept, err := s.inner.GetDepartmentByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicDepartmentRestored,
			events.NewDepartmentRestoredEvent(dept.ID, dept.Name))
	}
	return nil
}
