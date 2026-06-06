package doctorspecialization

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedRoomTypeService struct {
	inner RoomTypeService
	redis *cache.RedisClient
}

func NewCachedRoomTypeService(inner RoomTypeService, redisClient *cache.RedisClient) RoomTypeService {
	if redisClient == nil {
		return inner
	}
	return &cachedRoomTypeService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedRoomTypeService) List(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
	key := cache.RoomTypeListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.RoomTypeListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.List(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomTypeService) DeletedList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeDeletedListResponse, error) {
	key := cache.RoomTypeDeletedListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.RoomTypeDeletedListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.DeletedList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomTypeService) ActiveList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
	key := cache.RoomTypeActiveListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.RoomTypeListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.ActiveList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomTypeService) InactiveList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
	key := cache.RoomTypeInactiveListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.RoomTypeListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.InactiveList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomTypeService) FindByID(id uint) (*dto.RoomTypeResponse, error) {
	key := cache.RoomTypeKey(id)
	var resp dto.RoomTypeResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.FindByID(id)
	if err != nil {
		return nil, err
	}

	s.setCache(key, result)
	return result, nil
}
func (s *cachedRoomTypeService) Create(req *dto.CreateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}
func (s *cachedRoomTypeService) Update(id uint, req *dto.UpdateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}
func (s *cachedRoomTypeService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}
func (s *cachedRoomTypeService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}
func (s *cachedRoomTypeService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedRoomTypeService) Activate(id uint) error {
	if err := s.inner.Activate(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedRoomTypeService) Deactivate(id uint) error {
	if err := s.inner.Deactivate(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func normalizeCachePart(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *cachedRoomTypeService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedRoomTypeService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternRoomTypeAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternRoomTypeAll, err)
	}
}
