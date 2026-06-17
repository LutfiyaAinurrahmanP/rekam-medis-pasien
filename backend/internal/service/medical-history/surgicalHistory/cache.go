package surgicalhistory

import (
	"context"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedSurgicalHistoryService struct {
	inner SurgicalHistoryService
	redis *cache.RedisClient
}

func NewCachedSurgicalHistoryService(inner SurgicalHistoryService, redisClient *cache.RedisClient) SurgicalHistoryService {
	if redisClient == nil {
		return inner
	}
	return &cachedSurgicalHistoryService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedSurgicalHistoryService) List(query *dto.SurgicalHistoryPaginationQuery) (*dto.SurgicalHistoryListResponse, error) {
	key := cache.SurgicalHistoryListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.SurgicalHistoryListResponse
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

func (s *cachedSurgicalHistoryService) FindByID(id uint) (*dto.SurgicalHistoryResponse, error) {
	key := cache.SurgicalHistoryKey(id)
	var resp dto.SurgicalHistoryResponse
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

func (s *cachedSurgicalHistoryService) Create(req *dto.CreateSurgicalHistoryRequest) (*dto.SurgicalHistoryResponse, error) {
	result, err := s.inner.Create(req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedSurgicalHistoryService) Update(id uint, req *dto.UpdateSurgicalHistoryRequest) (*dto.SurgicalHistoryResponse, error) {
	result, err := s.inner.Update(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedSurgicalHistoryService) Delete(id uint) error {
	if err := s.inner.Delete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedSurgicalHistoryService) setCache(key string, value interface{}) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("Failed to set cache for %s: %v", key, err)
	}
}

func (s *cachedSurgicalHistoryService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternSurgicalHistoryAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternSurgicalHistoryAll, err)
	}
}
