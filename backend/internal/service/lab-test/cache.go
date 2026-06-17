package labtest

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedLabTestService struct {
	inner LabTestService
	redis *cache.RedisClient
}

func NewCachedLabTestService(inner LabTestService, redisClient *cache.RedisClient) LabTestService {
	if redisClient == nil {
		return inner
	}
	return &cachedLabTestService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedLabTestService) List(query *dto.LabTestPaginationQuery) (*dto.LabTestListResponse, error) {
	key := cache.LabTestListQueryKey(
		query.Page,
		query.PageSize,
		query.MedicalRecordID,
		query.TestTypeID,
		query.OrderedByDoctorID,
		normalizeCachePart(query.Status),
		normalizeCachePart(query.NotStatus),
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)

	var cached dto.LabTestListResponse
	if err := s.redis.Get(context.Background(), key, &cached); err == nil {
		return &cached, nil
	}

	res, err := s.inner.List(query)
	if err == nil {
		s.setCache(key, res)
	}

	return res, err
}

func (s *cachedLabTestService) DeletedList(query *dto.LabTestPaginationQuery) (*dto.LabTestDeletedListResponse, error) {
	key := cache.LabTestDeletedListQueryKey(
		query.Page,
		query.PageSize,
		query.MedicalRecordID,
		query.TestTypeID,
		query.OrderedByDoctorID,
		normalizeCachePart(query.Status),
		normalizeCachePart(query.NotStatus),
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)

	var cached dto.LabTestDeletedListResponse
	if err := s.redis.Get(context.Background(), key, &cached); err == nil {
		return &cached, nil
	}

	res, err := s.inner.DeletedList(query)
	if err == nil {
		s.setCache(key, res)
	}

	return res, err
}

func (s *cachedLabTestService) FindByID(id uint) (*dto.LabTestResponse, error) {
	key := cache.LabTestKey(id)

	var cached dto.LabTestResponse
	if err := s.redis.Get(context.Background(), key, &cached); err == nil {
		return &cached, nil
	}

	res, err := s.inner.FindByID(id)
	if err == nil {
		s.setCache(key, res)
	}

	return res, err
}

func (s *cachedLabTestService) FindByIDUnscoped(id uint) (*dto.LabTestResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *cachedLabTestService) Create(req *dto.CreateLabTestRequest) (*dto.LabTestResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedLabTestService) Update(id uint, req *dto.UpdateLabTestRequest) (*dto.LabTestResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedLabTestService) CollectSample(id uint) (*dto.LabTestResponse, error) {
	result, err := s.inner.CollectSample(id)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedLabTestService) Start(id uint) (*dto.LabTestResponse, error) {
	result, err := s.inner.Start(id)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedLabTestService) Complete(id uint, req *dto.CompleteLabTestRequest) (*dto.LabTestResponse, error) {
	result, err := s.inner.Complete(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedLabTestService) Cancel(id uint) (*dto.LabTestResponse, error) {
	result, err := s.inner.Cancel(id)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedLabTestService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedLabTestService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedLabTestService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedLabTestService) setCache(key string, value interface{}) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("Failed to set cache for %s: %v", key, err)
	}
}

func (s *cachedLabTestService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternLabTestAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternLabTestAll, err)
	}
}

func normalizeCachePart(s string) string {
	if s == "" {
		return "none"
	}
	return strings.ToLower(s)
}
