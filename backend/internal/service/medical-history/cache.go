package medicalhistory

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedMedicalHistoryService struct {
	inner MedicalHistoryService
	redis *cache.RedisClient
}

func NewCachedMedicalHistoryService(inner MedicalHistoryService, redisClient *cache.RedisClient) MedicalHistoryService {
	if redisClient == nil {
		return inner
	}
	return &cachedMedicalHistoryService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedMedicalHistoryService) List(query *dto.MedicalHistoryPaginationQuery) (*dto.MedicalHistoryListResponse, error) {
	key := cache.MedicalHistoriesListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicalHistoryListResponse
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

func (s *cachedMedicalHistoryService) FindByID(id uint) (*dto.MedicalHistoryDetailResponse, error) {
	key := cache.MedicalHistoriesKey(id)
	var resp dto.MedicalHistoryDetailResponse
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

func normalizeCachePart(value string) string {
	if value == "" {
		return "all"
	}
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *cachedMedicalHistoryService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedMedicalHistoryService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternMedicalHistoryAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternMedicalHistoryAll, err)
	}
}
