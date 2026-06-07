package medicinetype

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedMedicineTypeService struct {
	inner MedicineTypeService
	redis *cache.RedisClient
}

func NewCachedMedicineTypeService(inner MedicineTypeService, redisClient *cache.RedisClient) MedicineTypeService {
	if redisClient == nil {
		return inner
	}
	return &cachedMedicineTypeService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedMedicineTypeService) List(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
	key := cache.MedicineTypeListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineTypeListResponse
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

func (s *cachedMedicineTypeService) DeletedList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeDeletedListResponse, error) {
	key := cache.MedicineTypeDeletedListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineTypeDeletedListResponse
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

func (s *cachedMedicineTypeService) ActiveList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
	key := cache.MedicineTypeActiveListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineTypeListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.ActiveList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedMedicineTypeService) InactiveList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
	key := cache.MedicineTypeInactiveListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineTypeListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.InactiveList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedMedicineTypeService) FindByID(id uint) (*dto.MedicineTypeResponse, error) {
	key := cache.MedicineTypeKey(id)
	var resp dto.MedicineTypeResponse
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
func (s *cachedMedicineTypeService) Create(req *dto.CreateMedicineTypeRequest) (*dto.MedicineTypeResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}
func (s *cachedMedicineTypeService) Update(id uint, req *dto.UpdateMedicineTypeRequest) (*dto.MedicineTypeResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}
func (s *cachedMedicineTypeService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}
func (s *cachedMedicineTypeService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}
func (s *cachedMedicineTypeService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicineTypeService) Activate(id uint) error {
	if err := s.inner.Activate(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicineTypeService) Deactivate(id uint) error {
	if err := s.inner.Deactivate(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func normalizeCachePart(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *cachedMedicineTypeService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedMedicineTypeService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternMedicineTypeAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternMedicineTypeAll, err)
	}
}
