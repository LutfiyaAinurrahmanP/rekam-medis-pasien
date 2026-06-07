package typetest

import (
	"context"
	"fmt"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

// cachedTypeTestService membungkus TypeTestService dengan Redis caching.
// Read operations di-cache; write/delete operations menginvalidasi cache.
type cachedTypeTestService struct {
	inner TypeTestService
	redis *cache.RedisClient
}

// NewCachedTypeTestService returns a TypeTestService with Redis caching.
// Jika redisClient nil, langsung kembalikan inner service tanpa cache.
func NewCachedTypeTestService(inner TypeTestService, redisClient *cache.RedisClient) TypeTestService {
	if redisClient == nil {
		return inner
	}
	return &cachedTypeTestService{inner: inner, redis: redisClient}
}

// ─── Read operations (dengan cache) ────────────────────────────────────────

func (s *cachedTypeTestService) List(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	key := cache.TypeTestListKey(query.Page, query.PageSize)
	var resp dto.TypeTestListResponse
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

func (s *cachedTypeTestService) ListActive(query *dto.TypeTestPaginationQuery) (*dto.ActiveTypeTestListResponse, error) {
	key := fmt.Sprintf("typetest:active:p%d:s%d", query.Page, query.PageSize)
	var resp dto.ActiveTypeTestListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.ListActive(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedTypeTestService) ListInactive(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	key := fmt.Sprintf("typetest:inactive:p%d:s%d", query.Page, query.PageSize)
	var resp dto.TypeTestListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.ListInactive(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedTypeTestService) DeleteList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestDeletedListResponse, error) {
	// Deleted list tidak di-cache karena berubah dinamis
	return s.inner.DeleteList(query)
}

func (s *cachedTypeTestService) Search(query *dto.TypeTestSearchQuery) (*dto.TypeTestSearchResponse, error) {
	key := cache.TypeTestSearchKey(query.Keyword+"|"+query.Category, query.Page, query.PageSize)
	var resp dto.TypeTestSearchResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.Search(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedTypeTestService) FindByID(id uint) (*dto.TypeTestResponse, error) {
	key := cache.TypeTestKey(id)
	var resp dto.TypeTestResponse
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

func (s *cachedTypeTestService) FindByCode(code string) (*dto.TypeTestResponse, error) {
	key := fmt.Sprintf("typetest:code:%s", code)
	var resp dto.TypeTestResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.FindByCode(code)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

// ─── Write operations (invalidate cache) ───────────────────────────────────

func (s *cachedTypeTestService) Create(req *dto.CreateTypeTestRequest) (*dto.TypeTestResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedTypeTestService) Update(id uint, req *dto.UpdateTypeTestRequest) (*dto.TypeTestResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedTypeTestService) Activate(id uint) error {
	if err := s.inner.Activate(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedTypeTestService) Deactivate(id uint) error {
	if err := s.inner.Deactivate(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedTypeTestService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedTypeTestService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedTypeTestService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

// ─── helpers ───────────────────────────────────────────────────────────────

func (s *cachedTypeTestService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedTypeTestService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternTypeTestAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternTypeTestAll, err)
	}
}
