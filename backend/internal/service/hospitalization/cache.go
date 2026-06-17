package hospitalization

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedHospitalizationService struct {
	inner HospitalizationService
	redis *cache.RedisClient
}

func NewCachedHospitalizationService(inner HospitalizationService, redisClient *cache.RedisClient) HospitalizationService {
	if redisClient == nil {
		return inner
	}
	return &cachedHospitalizationService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedHospitalizationService) List(query *dto.HospitalizationPaginationQuery) (*dto.HospitalizationListResponse, error) {
	key := cache.HospitalizationListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		query.DoctorID,
		query.RoomID,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.Status),
		normalizeCachePart(query.NotStatus),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.HospitalizationListResponse
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

func (s *cachedHospitalizationService) DeletedList(query *dto.HospitalizationPaginationQuery) (*dto.HospitalizationDeletedListResponse, error) {
	key := cache.HospitalizationDeletedListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		query.DoctorID,
		query.RoomID,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.Status),
		normalizeCachePart(query.NotStatus),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.HospitalizationDeletedListResponse
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

func (s *cachedHospitalizationService) FindByID(id uint) (*dto.HospitalizationResponse, error) {
	key := cache.HospitalizationKey(id)
	var resp dto.HospitalizationResponse
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

func (s *cachedHospitalizationService) FindByIDUnscoped(id uint) (*dto.HospitalizationResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *cachedHospitalizationService) Create(req *dto.CreateHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	result, err := s.inner.Create(req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedHospitalizationService) Update(id uint, req *dto.UpdateHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	result, err := s.inner.Update(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedHospitalizationService) Discharge(id uint, req *dto.DischargeHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	result, err := s.inner.Discharge(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedHospitalizationService) Transfer(id uint, req *dto.TransferHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	result, err := s.inner.Transfer(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedHospitalizationService) Activate(id uint) (*dto.HospitalizationResponse, error) {
	result, err := s.inner.Activate(id)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedHospitalizationService) Deactivate(id uint) (*dto.HospitalizationResponse, error) {
	result, err := s.inner.Deactivate(id)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedHospitalizationService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedHospitalizationService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedHospitalizationService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedHospitalizationService) setCache(key string, value interface{}) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("Failed to set cache for %s: %v", key, err)
	}
}

func (s *cachedHospitalizationService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternHospitalizationAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternHospitalizationAll, err)
	}
}

func normalizeCachePart(s string) string {
	if s == "" {
		return "all"
	}
	return strings.ToLower(s)
}
