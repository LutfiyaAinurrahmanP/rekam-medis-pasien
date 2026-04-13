package medicine

import (
	"context"
	"log"

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
	key := cache.MedicineListKey(query.Page, query.PageSize)
	var resp dto.MedicineListResponse
	if err := s.redis.Get(context.Background(), key, resp); err == nil {
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
	key := cache.MedicineDeletedListKey(query.Page, query.PageSize)
	var resp dto.MedicineDeletedListResponse
	if err := s.redis.Get(context.Background(), key, resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.DeletedList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedMedicineService) ListByAvailable(query *dto.MedicinePaginationQuery) (*dto.MedicineAvailableResponse, error) {
	key := cache.MedicineAvailableKey(query.Page, query.PageSize)
	var resp dto.MedicineAvailableResponse
	if err := s.redis.Get(context.Background(), key, resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.ListByAvailable(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedMedicineService) ListByLowStock(query *dto.MedicinePaginationQuery) (*dto.MedicineLowStockResponse, error) {
	key := cache.MedicineLowStockKey(query.Page, query.PageSize)
	var resp dto.MedicineLowStockResponse

	if err := s.redis.Get(context.Background(), key, resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.ListByLowStock(query)
	if err != nil {
		return nil, err
	}

	s.setCache(key,result)
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
		return  nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedMedicineService) FindByName(name string) (*dto.MedicineResponse, error) {
	key := cache.MedicineNameKey(name)
	var resp dto.MedicineResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.FindByName(name)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedMedicineService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}