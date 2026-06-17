package medicine

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedMedicineService struct {
	inner MedicineService
	redis *cache.RedisClient
}

func NewCachedMedicineService(inner MedicineService, redisClient *cache.RedisClient) MedicineService {
	if redisClient == nil {
		return inner
	}
	return &cachedMedicineService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedMedicineService) List(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	// isActive mapped, stockStatus mapped, type mapped
	key := cache.MedicineListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		query.IsActive,
		query.MedicineTypeID,
		normalizeCachePart(query.StockStatus),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineListResponse
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

func (s *cachedMedicineService) DeletedList(query *dto.MedicinePaginationQuery) (*dto.MedicineDeletedListResponse, error) {
	key := cache.MedicineDeletedListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineDeletedListResponse
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

func (s *cachedMedicineService) AvailableList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	key := cache.MedicineListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		query.IsActive,
		query.MedicineTypeID,
		"available",
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.AvailableList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedMedicineService) LowStockList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	key := cache.MedicineListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		query.IsActive,
		query.MedicineTypeID,
		"low_stock",
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.LowStockList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedMedicineService) OutStockList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	key := cache.MedicineListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		query.IsActive,
		query.MedicineTypeID,
		"out_of_stock",
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.OutStockList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedMedicineService) ActiveList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	isActive := true
	key := cache.MedicineListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		&isActive,
		query.MedicineTypeID,
		normalizeCachePart(query.StockStatus),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineListResponse
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

func (s *cachedMedicineService) InactiveList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	isActive := false
	key := cache.MedicineListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		&isActive,
		query.MedicineTypeID,
		normalizeCachePart(query.StockStatus),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.MedicineListResponse
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

func (s *cachedMedicineService) FindByID(id uint) (*dto.MedicineResponse, error) {
	key := cache.MedicineKey(id)
	var resp dto.MedicineResponse
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

func (s *cachedMedicineService) FindByIDUnscoped(id uint) (*dto.DeletedMedicineResponse, error) {
	key := cache.MedicineKey(id)
	var resp dto.DeletedMedicineResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.FindByIDUnscoped(id)
	if err != nil {
		return nil, err
	}

	s.setCache(key, result)
	return result, nil
}

func (s *cachedMedicineService) Create(req *dto.CreateMedicineRequest) (*dto.MedicineResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedMedicineService) Update(id uint, req *dto.UpdateMedicineRequest) (*dto.MedicineResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedMedicineService) AddStock(id uint, req *dto.AddStockRequest) error {
	if err := s.inner.AddStock(id, req); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicineService) ReduceStock(id uint, req *dto.ReduceStockRequest) error {
	if err := s.inner.ReduceStock(id, req); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicineService) Activate(id uint) error {
	if err := s.inner.Activate(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicineService) Deactivate(id uint, req *dto.DeactivateMedicineRequest) error {
	if err := s.inner.Deactivate(id, req); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicineService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicineService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedMedicineService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func normalizeCachePart(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *cachedMedicineService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedMedicineService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternMedicineAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternMedicineAll, err)
	}
}
