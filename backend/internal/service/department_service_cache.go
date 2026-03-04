package service

import (
	"context"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

// cachedDepartmentService membungkus DepartmentService dengan Redis caching.
type cachedDepartmentService struct {
	inner DepartmentService
	redis *cache.RedisClient
}

// NewCachedDepartmentService returns a DepartmentService with Redis caching.
// Jika redisClient nil, langsung kembalikan inner service tanpa cache.
func NewCachedDepartmentService(inner DepartmentService, redisClient *cache.RedisClient) DepartmentService {
	if redisClient == nil {
		return inner
	}
	return &cachedDepartmentService{inner: inner, redis: redisClient}
}

// ─── Read operations (dengan cache) ────────────────────────────────────────

func (s *cachedDepartmentService) ListDepartments(query *dto.DepartmentPaginationQuery) (*dto.DepartmentListResponse, error) {
	key := cache.DepartmentListKey(query.Page, query.PageSize)
	var resp dto.DepartmentListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.ListDepartments(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedDepartmentService) DeleteListDepartments(query *dto.DepartmentPaginationQuery) (*dto.DepartmentDeletedListResponse, error) {
	// Deleted list tidak di-cache karena berubah dinamis
	return s.inner.DeleteListDepartments(query)
}

func (s *cachedDepartmentService) GetDepartmentByID(id uint) (*dto.DepartmentResponse, error) {
	key := cache.DepartmentKey(id)
	var resp dto.DepartmentResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetDepartmentByID(id)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

// ─── Write operations (invalidate cache) ───────────────────────────────────

func (s *cachedDepartmentService) CreateDepartment(req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	result, err := s.inner.CreateDepartment(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedDepartmentService) UpdateDepartment(id uint, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error) {
	result, err := s.inner.UpdateDepartment(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedDepartmentService) SoftDeleteDepartment(id uint) error {
	if err := s.inner.SoftDeleteDepartment(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedDepartmentService) RestoreDepartment(id uint) error {
	if err := s.inner.RestoreDepartment(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedDepartmentService) HardDeleteDepartment(id uint) error {
	if err := s.inner.HardDeleteDepartment(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

// ─── helpers ───────────────────────────────────────────────────────────────

func (s *cachedDepartmentService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedDepartmentService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternDepartmentAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternDepartmentAll, err)
	}
}

// ─── compile-time interface check ──────────────────────────────────────────
var _ DepartmentService = (*cachedDepartmentService)(nil)
