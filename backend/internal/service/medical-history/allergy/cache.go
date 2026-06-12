package allergy

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedAllergyService struct {
	inner AllergyService
	redis *cache.RedisClient
}

func NewCachedAllergyService(inner AllergyService, redisClient *cache.RedisClient) AllergyService {
	if redisClient == nil {
		return inner
	}
	return &cachedAllergyService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedAllergyService) List(query *dto.AllergyPaginationQuery) (*dto.AllergyListResponse, error) {
	key := cache.AllergyListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.AllergyListResponse
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

func (s *cachedAllergyService) FindByID(id uint) (*dto.AllergyResponse, error) {
	key := cache.AllergyKey(id)
	var resp dto.AllergyResponse
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

func (s *cachedAllergyService) Create(req *dto.CreateAllergyRequest) (*dto.AllergyResponse, error) {
	result, err := s.inner.Create(req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedAllergyService) Update(id uint, req *dto.UpdateAllergyRequest) (*dto.AllergyResponse, error) {
	result, err := s.inner.Update(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedAllergyService) Delete(id uint) error {
	if err := s.inner.Delete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func normalizeCachePart(value string) string {
	if value == "" {
		return "all"
	}
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *cachedAllergyService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedAllergyService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternAllergyAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternAllergyAll, err)
	}
}
