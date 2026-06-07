package typetestcategory

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventTypeTestCategoryService struct {
	inner     TypeTestCategoryService
	publisher kafka.EventPublisher
}

func NewEventTypeTestCategoryService(inner TypeTestCategoryService, publisher kafka.EventPublisher) TypeTestCategoryService {
	if publisher == nil {
		return inner
	}
	return &eventTypeTestCategoryService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventTypeTestCategoryService) List(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
	return s.inner.List(query)
}
func (s *eventTypeTestCategoryService) ActiveList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
	return s.inner.ActiveList(query)
}
func (s *eventTypeTestCategoryService) InactiveList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
	return s.inner.InactiveList(query)
}
func (s *eventTypeTestCategoryService) DeletedList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}
func (s *eventTypeTestCategoryService) FindByID(id uint) (*dto.TypeTestCategoryResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventTypeTestCategoryService) Create(req *dto.CreateTypeTestCategoryRequest) (*dto.TypeTestCategoryResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicTypeTestCategoryCreated,
		events.NewTypeTestCategoryCreatedEvent(result.ID, result.Name, result.Code, result.Description, result.IsActive),
	)
	return result, nil
}
func (s *eventTypeTestCategoryService) Update(id uint, req *dto.UpdateTypeTestCategoryRequest) (*dto.TypeTestCategoryResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(
		kafka.TopicTypeTestCategoryUpdated,
		events.NewTypeTestCategoryUpdatedEvent(result.ID, result.Name, result.Code, result.Description, result.IsActive, "update"),
	)
	return result, nil
}
func (s *eventTypeTestCategoryService) SoftDelete(id uint) error {
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
		kafka.TopicTypeTestCategoryDeleted,
		events.NewTypeTestCategoryDeletedEvent(id, name, code, "soft_delete"),
	)
	return nil
}
func (s *eventTypeTestCategoryService) Restore(id uint) error {
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
		kafka.TopicTypeTestCategoryRestored,
		events.NewTypeTestCategoryRestoredEvent(id, name, code),
	)
	return nil
}
func (s *eventTypeTestCategoryService) HardDelete(id uint) error {
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
		kafka.TopicTypeTestCategoryDeleted,
		events.NewTypeTestCategoryDeletedEvent(id, name, code, "hard_delete"),
	)
	return nil
}

func (s *eventTypeTestCategoryService) Activate(id uint) error {
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
		kafka.TopicTypeTestCategoryUpdated,
		events.NewTypeTestCategoryUpdatedEvent(id, name, code, ds.Description, true, "activate"),
	)
	return nil
}

func (s *eventTypeTestCategoryService) Deactivate(id uint) error {
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
		kafka.TopicTypeTestCategoryUpdated,
		events.NewTypeTestCategoryUpdatedEvent(id, name, code, ds.Description, false, "deactivate"),
	)
	return nil
}
