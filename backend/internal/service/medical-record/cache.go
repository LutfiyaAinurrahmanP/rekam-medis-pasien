package medicalrecord

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedMedicalRecordService struct {
	inner MedicalRecordService
	redis *cache.RedisClient
}

func NewCachedMedicalRecordService(inner MedicalRecordService, redisClient *cache.RedisClient) MedicalRecordService {
	if redisClient == nil {
		return inner
	}
	return &cachedMedicalRecordService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedMedicalRecordService) List(query *dto.MedicalRecordPaginationQuery) (*dto.MedicalRecordListResponse, error) {
	key := cache.MedicalRecordListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		query.DoctorID,
		query.DepartmentID,
		normalizeCachePart(query.Status),
		normalizeCachePart(query.DateFrom),
		normalizeCachePart(query.DateTo),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicalRecordListResponse
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

func (s *cachedMedicalRecordService) DeletedList(query *dto.MedicalRecordPaginationQuery) (*dto.MedicalRecordDeletedListResponse, error) {
	key := cache.MedicalRecordDeletedListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		query.DoctorID,
		query.DepartmentID,
		normalizeCachePart(query.Status),
		normalizeCachePart(query.DateFrom),
		normalizeCachePart(query.DateTo),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicalRecordDeletedListResponse
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

func (s *cachedMedicalRecordService) FindByID(id uint) (*dto.MedicalRecordResponse, error) {
	key := cache.MedicalRecordKey(id)
	var resp dto.MedicalRecordResponse
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

func (s *cachedMedicalRecordService) FindByIDUnscoped(id uint) (*dto.MedicalRecordResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *cachedMedicalRecordService) Create(req *dto.CreateMedicalRecordRequest) (*dto.MedicalRecordResponse, error) {
	result, err := s.inner.Create(req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedMedicalRecordService) Update(id uint, req *dto.UpdateMedicalRecordRequest) (*dto.MedicalRecordResponse, error) {
	result, err := s.inner.Update(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedMedicalRecordService) Finalize(id uint) error {
	if err := s.inner.Finalize(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicalRecordService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicalRecordService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicalRecordService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicalRecordService) setCache(key string, value interface{}) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("Failed to set cache for %s: %v", key, err)
	}
}

func (s *cachedMedicalRecordService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternMedicalRecordAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternMedicalRecordAll, err)
	}
}

func normalizeCachePart(s string) string {
	if s == "" {
		return "all"
	}
	return strings.ToLower(s)
}
