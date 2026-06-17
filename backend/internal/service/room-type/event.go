package doctorspecialization

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventRoomTypeService struct {
	inner     RoomTypeService
	publisher kafka.EventPublisher
}

func NewEventRoomTypeService(inner RoomTypeService, publisher kafka.EventPublisher) RoomTypeService {
	if publisher == nil {
		return inner
	}
	return &eventRoomTypeService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventRoomTypeService) List(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
	return s.inner.List(query)
}
func (s *eventRoomTypeService) ActiveList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
	return s.inner.ActiveList(query)
}
func (s *eventRoomTypeService) InactiveList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
	return s.inner.InactiveList(query)
}
func (s *eventRoomTypeService) DeletedList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}
func (s *eventRoomTypeService) FindByID(id uint) (*dto.RoomTypeResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventRoomTypeService) Create(req *dto.CreateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicRoomTypeCreated,
		events.NewRoomTypeCreatedEvent(result.ID, result.Name, result.Code, result.Description, result.IsActive),
	)
	return result, nil
}
func (s *eventRoomTypeService) Update(id uint, req *dto.UpdateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(
		kafka.TopicRoomTypeUpdated,
		events.NewRoomTypeUpdatedEvent(result.ID, result.Name, result.Code, result.Description, result.IsActive, "update"),
	)
	return result, nil
}
func (s *eventRoomTypeService) SoftDelete(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	name := ""
	if ds != nil {
		name = ds.Name
	}

	s.publisher.PublishAsync(
		kafka.TopicRoomTypeDeleted,
		events.NewRoomTypeDeletedEvent(id, name, "soft_delete"),
	)
	return nil
}
func (s *eventRoomTypeService) Restore(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	name := ""
	if ds != nil {
		name = ds.Name
	}

	s.publisher.PublishAsync(
		kafka.TopicRoomTypeRestored,
		events.NewRoomTypeRestoredEvent(id, name),
	)
	return nil
}
func (s *eventRoomTypeService) HardDelete(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	name := ""
	if ds != nil {
		name = ds.Name
	}

	s.publisher.PublishAsync(
		kafka.TopicRoomTypeDeleted,
		events.NewRoomTypeDeletedEvent(id, name, "hard_delete"),
	)
	return nil
}

func (s *eventRoomTypeService) Activate(id uint) error {
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
		kafka.TopicRoomTypeUpdated,
		events.NewRoomTypeUpdatedEvent(id, name, code, ds.Description, true, "activate"),
	)
	return nil
}

func (s *eventRoomTypeService) Deactivate(id uint) error {
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
		kafka.TopicRoomTypeUpdated,
		events.NewRoomTypeUpdatedEvent(id, name, code, ds.Description, false, "deactivate"),
	)
	return nil
}
