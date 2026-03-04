package service

import (
	"context"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

// roomAvailabilityTTL adalah TTL khusus untuk data ketersediaan kamar.
// Data ini berubah cepat (bed occupancy), jadi TTL-nya dibuat lebih pendek
// dari TTL default (biasanya 5m) agar data tidak stale terlalu lama.
const roomAvailabilityTTL = 30 * time.Second

// cachedRoomService membungkus RoomService dengan Redis caching.
//
// Best practices yang diterapkan:
//   - GetAvailableRooms / GetOccupiedRooms → TTL pendek (30 detik) karena data volatile
//   - OccupyRoom / ReleaseRoom → invalidasi seluruh namespace room:* (ketersediaan berubah)
//   - GetRoomByID / ListRooms / GetByRoomNumber → TTL default
//   - Deleted list tidak di-cache
type cachedRoomService struct {
	inner RoomService
	redis *cache.RedisClient
}

// NewCachedRoomService returns a RoomService with Redis caching.
// Jika redisClient nil, langsung kembalikan inner service tanpa cache.
func NewCachedRoomService(inner RoomService, redisClient *cache.RedisClient) RoomService {
	if redisClient == nil {
		return inner
	}
	return &cachedRoomService{inner: inner, redis: redisClient}
}

// ─── Read operations (dengan cache) ────────────────────────────────────────

func (s *cachedRoomService) ListRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	key := cache.RoomListKey(query.Page, query.PageSize)
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
	// Deleted list tidak di-cache
	return s.inner.DeleteListRooms(query)
}

// GetAvailableRooms menggunakan TTL pendek karena data ketersediaan berubah cepat.
func (s *cachedRoomService) GetAvailableRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	key := cache.RoomAvailableKey(query.Page, query.PageSize)
	var resp dto.RoomListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetAvailableRooms(query)
	if err != nil {
		return nil, err
	}
	s.setCacheWithTTL(key, result, roomAvailabilityTTL)
	return result, nil
}

// GetOccupiedRooms menggunakan TTL pendek karena data occupancy berubah cepat.
func (s *cachedRoomService) GetOccupiedRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	key := cache.RoomOccupiedKey(query.Page, query.PageSize)
	var resp dto.RoomListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetOccupiedRooms(query)
	if err != nil {
		return nil, err
	}
	s.setCacheWithTTL(key, result, roomAvailabilityTTL)
	return result, nil
}

func (s *cachedRoomService) GetInactiveRooms(query *dto.RoomPaginationQuery) (*dto.RoomListResponse, error) {
	key := cache.RoomInactiveKey(query.Page, query.PageSize)
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

func (s *cachedRoomService) GetByRoomNumber(roomNumber string) (*dto.RoomResponse, error) {
	key := cache.RoomByNumberKey(roomNumber)
	var resp dto.RoomResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetByRoomNumber(roomNumber)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomService) GetByRoomType(roomType string) (*dto.RoomResponse, error) {
	key := cache.RoomByTypeKey(roomType)
	var resp dto.RoomResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetByRoomType(roomType)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedRoomService) GetByDepatymentID(deptID string) (*dto.RoomResponse, error) {
	key := cache.RoomByDeptKey(deptID)
	var resp dto.RoomResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetByDepatymentID(deptID)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

// ─── Write operations (invalidate cache) ───────────────────────────────────

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

// OccupyRoom dan ReleaseRoom sangat sering dipanggil dan langsung mengubah
// ketersediaan kamar → invalidasi seluruh namespace room:*
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

// ─── helpers ───────────────────────────────────────────────────────────────

func (s *cachedRoomService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedRoomService) setCacheWithTTL(key string, value any, ttl time.Duration) {
	if err := s.redis.Set(context.Background(), key, value, ttl); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedRoomService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternRoomAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternRoomAll, err)
	}
}

// compile-time interface check
var _ RoomService = (*cachedRoomService)(nil)
