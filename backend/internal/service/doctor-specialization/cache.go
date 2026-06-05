package doctorspecialization

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedDoctorSpecializationService struct {
	inner DoctorSpecializationService
	redis *cache.RedisClient
}

func NewCachedDoctorSpecializationService(inner DoctorSpecializationService, redisClient *cache.RedisClient) DoctorSpecializationService {
	if redisClient == nil {
		return inner
	}
	return &cachedDoctorSpecializationService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedDoctorSpecializationService) List(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error) {
	key := cache.DoctorSpecializationListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.DoctorSpecializationListResponse
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

func (s *cachedDoctorSpecializationService) DeletedList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationDeletedListResponse, error) {
	key := cache.DoctorSpecializationDeletedListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.DoctorSpecializationDeletedListResponse
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

func (s *cachedDoctorSpecializationService) FindByID(id uint) (*dto.DoctorSpecializationResponse, error) {
	key := cache.DoctorSpecializationKey(id)
	var resp dto.DoctorSpecializationResponse
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
func (s *cachedDoctorSpecializationService) Create(req *dto.CreateDoctorSpecializationRequest) (*dto.DoctorSpecializationResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}
func (s *cachedDoctorSpecializationService) Update(id uint, req *dto.UpdateDoctorSpecializationRequest) (*dto.DoctorSpecializationResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}
func (s *cachedDoctorSpecializationService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}
func (s *cachedDoctorSpecializationService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}
func (s *cachedDoctorSpecializationService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func normalizeCachePart(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *cachedDoctorSpecializationService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedDoctorSpecializationService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternDoctorSpecializationAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternDoctorSpecializationAll, err)
	}
}
