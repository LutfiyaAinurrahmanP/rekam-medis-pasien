package vitalsign

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedVitalSignService struct {
	inner VitalSignService
	redis *cache.RedisClient
}

func NewCachedVitalSignService(inner VitalSignService, redisClient *cache.RedisClient) VitalSignService {
	if redisClient == nil {
		return inner
	}
	return &cachedVitalSignService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedVitalSignService) List(query *dto.VitalSignPaginationQuery) (*dto.VitalSignListResponse, error) {
	key := cache.VitalSignListQueryKey(
		query.Page,
		query.PageSize,
		query.MedicalRecordID,
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.VitalSignListResponse
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

func (s *cachedVitalSignService) DeletedList(query *dto.VitalSignPaginationQuery) (*dto.VitalSignDeletedListResponse, error) {
	key := cache.VitalSignDeletedListQueryKey(
		query.Page,
		query.PageSize,
		query.MedicalRecordID,
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.VitalSignDeletedListResponse
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

func (s *cachedVitalSignService) FindByID(id uint) (*dto.VitalSignResponse, error) {
	key := cache.VitalSignKey(id)
	var resp dto.VitalSignResponse
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

func (s *cachedVitalSignService) FindByIDUnscoped(id uint) (*dto.VitalSignResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *cachedVitalSignService) Create(req *dto.CreateVitalSignRequest) (*dto.VitalSignResponse, error) {
	result, err := s.inner.Create(req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedVitalSignService) Update(id uint, req *dto.UpdateVitalSignRequest) (*dto.VitalSignResponse, error) {
	result, err := s.inner.Update(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedVitalSignService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedVitalSignService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedVitalSignService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedVitalSignService) setCache(key string, value interface{}) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("Failed to set cache for %s: %v", key, err)
	}
}

func (s *cachedVitalSignService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternVitalSignAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternVitalSignAll, err)
	}
}

func normalizeCachePart(s string) string {
	if s == "" {
		return "all"
	}
	return strings.ToLower(s)
}
