package room

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

// eventRoomService membungkus RoomService dengan Kafka event publishing.
type eventRoomService struct {
	inner     RoomService
	publisher kafka.EventPublisher
}

// NewEventRoomService mengembalikan RoomService dengan event publishing.
func NewEventRoomService(inner RoomService, publisher kafka.EventPublisher) RoomService {
	if publisher == nil {
		return inner
	}
	return &eventRoomService{inner: inner, publisher: publisher}
}

// ─── Read operations ──────────────────────────────────────────────────────────

func (s *eventRoomService) ListRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	return s.inner.ListRooms(query)
}

func (s *eventRoomService) DeleteListRooms(query *dto.RoomPaginationQuery) (*dto.RoomDeletedListResponse, error) {
	return s.inner.DeleteListRooms(query)
}

func (s *eventRoomService) GetAvailableRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	return s.inner.GetAvailableRooms(query)
}

func (s *eventRoomService) GetOccupiedRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	return s.inner.GetOccupiedRooms(query)
}

func (s *eventRoomService) GetActiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	return s.inner.GetActiveRooms(query)
}

func (s *eventRoomService) GetInactiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	return s.inner.GetInactiveRooms(query)
}

func (s *eventRoomService) GetRoomByID(id uint) (*dto.RoomResponse, error) {
	return s.inner.GetRoomByID(id)
}

// ─── Write operations ─────────────────────────────────────────────────────────

func (s *eventRoomService) CreateRoom(req *dto.CreateRoomRequest) (*dto.RoomResponse, error) {
	result, err := s.inner.CreateRoom(req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicRoomCreated,
		events.NewRoomCreatedEvent(
			result.ID, result.RoomNumber, result.RoomTypeID, result.DepartmentID,
			result.BedCapacity, result.AvailableBeds, result.PricePerDay, result.IsActive,
		))
	return result, nil
}

func (s *eventRoomService) UpdateRoom(id uint, req *dto.UpdateRoomRequest) (*dto.RoomResponse, error) {
	result, err := s.inner.UpdateRoom(id, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicRoomUpdated,
		events.NewRoomUpdatedEvent(
			result.ID, result.RoomNumber, result.RoomTypeID, result.DepartmentID,
			result.BedCapacity, result.AvailableBeds, result.PricePerDay, result.IsActive, "update",
		))
	return result, nil
}

func (s *eventRoomService) ActivateRoom(id uint) (*dto.RoomResponse, error) {
	result, err := s.inner.ActivateRoom(id)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicRoomUpdated,
		events.NewRoomUpdatedEvent(
			result.ID, result.RoomNumber, result.RoomTypeID, result.DepartmentID,
			result.BedCapacity, result.AvailableBeds, result.PricePerDay, result.IsActive, "activate",
		))
	return result, nil
}

func (s *eventRoomService) DeactivateRoom(id uint) (*dto.RoomResponse, error) {
	result, err := s.inner.DeactivateRoom(id)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicRoomUpdated,
		events.NewRoomUpdatedEvent(
			result.ID, result.RoomNumber, result.RoomTypeID, result.DepartmentID,
			result.BedCapacity, result.AvailableBeds, result.PricePerDay, result.IsActive, "deactivate",
		))
	return result, nil
}

func (s *eventRoomService) OccupyRoom(id uint, beds int) (*dto.RoomResponse, error) {
	result, err := s.inner.OccupyRoom(id, beds)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicRoomUpdated,
		events.NewRoomUpdatedEvent(
			result.ID, result.RoomNumber, result.RoomTypeID, result.DepartmentID,
			result.BedCapacity, result.AvailableBeds, result.PricePerDay, result.IsActive, "occupy",
		))
	return result, nil
}

func (s *eventRoomService) ReleaseRoom(id uint, beds int) (*dto.RoomResponse, error) {
	result, err := s.inner.ReleaseRoom(id, beds)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicRoomUpdated,
		events.NewRoomUpdatedEvent(
			result.ID, result.RoomNumber, result.RoomTypeID, result.DepartmentID,
			result.BedCapacity, result.AvailableBeds, result.PricePerDay, result.IsActive, "release",
		))
	return result, nil
}

func (s *eventRoomService) SoftDeleteRoom(id uint) error {
	room, _ := s.inner.GetRoomByID(id)
	if err := s.inner.SoftDeleteRoom(id); err != nil {
		return err
	}
	roomNumber := ""
	if room != nil {
		roomNumber = room.RoomNumber
	}
	s.publisher.PublishAsync(kafka.TopicRoomDeleted,
		events.NewRoomDeletedEvent(id, roomNumber, "soft_delete"))
	return nil
}

func (s *eventRoomService) HardDeleteRoom(id uint) error {
	room, _ := s.inner.GetRoomByID(id)
	if err := s.inner.HardDeleteRoom(id); err != nil {
		return err
	}
	roomNumber := ""
	if room != nil {
		roomNumber = room.RoomNumber
	}
	s.publisher.PublishAsync(kafka.TopicRoomDeleted,
		events.NewRoomDeletedEvent(id, roomNumber, "hard_delete"))
	return nil
}

func (s *eventRoomService) RestoreRoom(id uint) error {
	if err := s.inner.RestoreRoom(id); err != nil {
		return err
	}
	if room, err := s.inner.GetRoomByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicRoomRestored,
			events.NewRoomRestoredEvent(room.ID, room.RoomNumber))
	}
	return nil
}
