package familyhistory

import (
	"context"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedFamilyHistoryService struct {
	inner FamilyHistoryService
	redis *cache.RedisClient
}

func NewCachedFamilyHistoryService(inner FamilyHistoryService, redisClient *cache.RedisClient) FamilyHistoryService {
	if redisClient == nil {
		return inner
	}
	return &cachedFamilyHistoryService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedFamilyHistoryService) List(query *dto.FamilyHistoryPaginationQuery) (*dto.FamilyHistoryListResponse, error) {
	key := cache.FamilyHistoryListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.FamilyHistoryListResponse
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

func (s *cachedFamilyHistoryService) FindByID(id uint) (*dto.FamilyHistoryResponse, error) {
	key := cache.FamilyHistoryKey(id)
	var resp dto.FamilyHistoryResponse
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

func (s *cachedFamilyHistoryService) Create(req *dto.CreateFamilyHistoryRequest) (*dto.FamilyHistoryResponse, error) {
	result, err := s.inner.Create(req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedFamilyHistoryService) Update(id uint, req *dto.UpdateFamilyHistoryRequest) (*dto.FamilyHistoryResponse, error) {
	result, err := s.inner.Update(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedFamilyHistoryService) Delete(id uint) error {
	if err := s.inner.Delete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedFamilyHistoryService) setCache(key string, value interface{}) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("Failed to set cache for %s: %v", key, err)
	}
}

func (s *cachedFamilyHistoryService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternFamilyHistoryAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternFamilyHistoryAll, err)
	}
}
