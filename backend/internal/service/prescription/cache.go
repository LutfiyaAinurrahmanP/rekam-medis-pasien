package prescription

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedPrescriptionService struct {
	inner PrescriptionService
	redis *cache.RedisClient
}

func NewCachedPrescriptionService(inner PrescriptionService, redisClient *cache.RedisClient) PrescriptionService {
	if redisClient == nil {
		return inner
	}
	return &cachedPrescriptionService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedPrescriptionService) List(query *dto.PrescriptionPaginationQuery) (*dto.PrescriptionListResponse, error) {
	key := cache.PrescriptionListQueryKey(
		query.Page,
		query.PageSize,
		query.DoctorID,
		query.MedicalRecordID,
		normalizeCachePart(query.Status),
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.PrescriptionListResponse
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

func (s *cachedPrescriptionService) DeletedList(query *dto.PrescriptionPaginationQuery) (*dto.PrescriptionDeletedListResponse, error) {
	key := cache.PrescriptionDeletedListQueryKey(
		query.Page,
		query.PageSize,
		query.DoctorID,
		query.MedicalRecordID,
		normalizeCachePart(query.Status),
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.PrescriptionDeletedListResponse
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

func (s *cachedPrescriptionService) FindByID(id uint) (*dto.PrescriptionResponse, error) {
	key := cache.PrescriptionKey(id)
	var resp dto.PrescriptionResponse
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

func (s *cachedPrescriptionService) FindByIDUnscoped(id uint) (*dto.PrescriptionResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *cachedPrescriptionService) Create(req *dto.CreatePrescriptionRequest) (*dto.PrescriptionResponse, error) {
	result, err := s.inner.Create(req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedPrescriptionService) Update(id uint, req *dto.UpdatePrescriptionRequest) (*dto.PrescriptionResponse, error) {
	result, err := s.inner.Update(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedPrescriptionService) Dispense(id uint) (*dto.PrescriptionResponse, error) {
	result, err := s.inner.Dispense(id)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedPrescriptionService) Cancel(id uint) (*dto.PrescriptionResponse, error) {
	result, err := s.inner.Cancel(id)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedPrescriptionService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedPrescriptionService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedPrescriptionService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedPrescriptionService) ListItems(prescriptionID uint) ([]dto.PrescriptionItemResponse, error) {
	return s.inner.ListItems(prescriptionID)
}

func (s *cachedPrescriptionService) FindItemByID(prescriptionID, itemID uint) (*dto.PrescriptionItemResponse, error) {
	return s.inner.FindItemByID(prescriptionID, itemID)
}

func (s *cachedPrescriptionService) CreateItem(prescriptionID uint, req *dto.CreatePrescriptionItemRequest) (*dto.PrescriptionItemResponse, error) {
	result, err := s.inner.CreateItem(prescriptionID, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedPrescriptionService) UpdateItem(prescriptionID, itemID uint, req *dto.UpdatePrescriptionItemRequest) (*dto.PrescriptionItemResponse, error) {
	result, err := s.inner.UpdateItem(prescriptionID, itemID, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedPrescriptionService) DeleteItem(prescriptionID, itemID uint) error {
	if err := s.inner.DeleteItem(prescriptionID, itemID); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedPrescriptionService) setCache(key string, value interface{}) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("Failed to set cache for %s: %v", key, err)
	}
}

func (s *cachedPrescriptionService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternPrescriptionAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternPrescriptionAll, err)
	}
}

func normalizeCachePart(s string) string {
	if s == "" {
		return "all"
	}
	return strings.ToLower(s)
}
