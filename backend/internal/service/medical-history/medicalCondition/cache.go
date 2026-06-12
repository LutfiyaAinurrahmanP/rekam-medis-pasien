package medicalcondition

import (
	"context"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedMedicalConditionService struct {
	inner MedicalConditionService
	redis *cache.RedisClient
}

func NewCachedMedicalConditionService(inner MedicalConditionService, redisClient *cache.RedisClient) MedicalConditionService {
	if redisClient == nil {
		return inner
	}
	return &cachedMedicalConditionService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedMedicalConditionService) List(query *dto.MedicalConditionPaginationQuery) (*dto.MedicalConditionListResponse, error) {
	key := cache.MedicalConditionListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		normalizeCachePart(query.Status),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicalConditionListResponse
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

func (s *cachedMedicalConditionService) FindByID(id uint) (*dto.MedicalConditionResponse, error) {
	key := cache.MedicalConditionKey(id)
	var resp dto.MedicalConditionResponse
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

func (s *cachedMedicalConditionService) Create(req *dto.CreateMedicalConditionRequest) (*dto.MedicalConditionResponse, error) {
	result, err := s.inner.Create(req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedMedicalConditionService) Update(id uint, req *dto.UpdateMedicalConditionRequest) (*dto.MedicalConditionResponse, error) {
	result, err := s.inner.Update(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedMedicalConditionService) Delete(id uint) error {
	if err := s.inner.Delete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicalConditionService) setCache(key string, value interface{}) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("Failed to set cache for %s: %v", key, err)
	}
}

func (s *cachedMedicalConditionService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternMedicalConditionAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternMedicalConditionAll, err)
	}
}
