package room

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedRoomService struct {
	inner RoomService
	redis *cache.RedisClient
}

func NewCachedRoomService(inner RoomService, redisClient *cache.RedisClient) RoomService {
	if redisClient == nil {
		return inner
	}
	return &cachedRoomService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedRoomService) ListRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	key := cache.RoomListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.RoomListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.ListRooms(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomService) DeleteListRooms(query *dto.RoomPaginationQuery) (*dto.RoomDeletedListResponse, error) {
	key := cache.RoomDeletedListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.RoomDeletedListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.DeleteListRooms(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomService) GetAvailableRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	key := cache.RoomAvailableListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.RoomListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.GetAvailableRooms(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomService) GetOccupiedRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	key := cache.RoomOccupiedListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.RoomListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.GetOccupiedRooms(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomService) GetActiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	key := cache.RoomActiveListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.RoomListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.GetActiveRooms(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomService) GetInactiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	key := cache.RoomInactiveListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.RoomListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.GetInactiveRooms(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomService) GetRoomByID(id uint) (*dto.RoomResponse, error) {
	key := cache.RoomKey(id)
	var resp dto.RoomResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.GetRoomByID(id)
	if err != nil {
		return nil, err
	}

	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomService) CreateRoom(req *dto.CreateRoomRequest) (*dto.RoomResponse, error) {
	result, err := s.inner.CreateRoom(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedRoomService) UpdateRoom(id uint, req *dto.UpdateRoomRequest) (*dto.RoomResponse, error) {
	result, err := s.inner.UpdateRoom(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedRoomService) ActivateRoom(id uint) (*dto.RoomResponse, error) {
	result, err := s.inner.ActivateRoom(id)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedRoomService) DeactivateRoom(id uint) (*dto.RoomResponse, error) {
	result, err := s.inner.DeactivateRoom(id)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedRoomService) OccupyRoom(id uint, beds int) (*dto.RoomResponse, error) {
	result, err := s.inner.OccupyRoom(id, beds)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedRoomService) ReleaseRoom(id uint, beds int) (*dto.RoomResponse, error) {
	result, err := s.inner.ReleaseRoom(id, beds)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedRoomService) SoftDeleteRoom(id uint) error {
	if err := s.inner.SoftDeleteRoom(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedRoomService) RestoreRoom(id uint) error {
	if err := s.inner.RestoreRoom(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedRoomService) HardDeleteRoom(id uint) error {
	if err := s.inner.HardDeleteRoom(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func normalizeCachePart(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *cachedRoomService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedRoomService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternRoomAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternRoomAll, err)
	}
}
