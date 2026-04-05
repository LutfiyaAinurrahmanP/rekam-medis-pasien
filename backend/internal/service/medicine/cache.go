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

func (s *cachedMedicineService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}