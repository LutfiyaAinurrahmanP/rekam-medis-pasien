package typetestcategory

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedTypeTestCategoryService struct {
	inner TypeTestCategoryService
	redis *cache.RedisClient
}

func NewCachedTypeTestCategoryService(inner TypeTestCategoryService, redisClient *cache.RedisClient) TypeTestCategoryService {
	if redisClient == nil {
		return inner
	}
	return &cachedTypeTestCategoryService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedTypeTestCategoryService) List(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
	key := cache.TypeTestCategoryListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.TypeTestCategoryListResponse
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

func (s *cachedTypeTestCategoryService) DeletedList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryDeletedListResponse, error) {
	key := cache.TypeTestCategoryDeletedListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.TypeTestCategoryDeletedListResponse
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

func (s *cachedTypeTestCategoryService) ActiveList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
	key := cache.TypeTestCategoryActiveListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.TypeTestCategoryListResponse
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

func (s *cachedTypeTestCategoryService) InactiveList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
	key := cache.TypeTestCategoryInactiveListQuery(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.TypeTestCategoryListResponse
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

func (s *cachedTypeTestCategoryService) FindByID(id uint) (*dto.TypeTestCategoryResponse, error) {
	key := cache.TypeTestCategoryKey(id)
	var resp dto.TypeTestCategoryResponse
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

func (s *cachedTypeTestCategoryService) Create(req *dto.CreateTypeTestCategoryRequest) (*dto.TypeTestCategoryResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}
func (s *cachedTypeTestCategoryService) Update(id uint, req *dto.UpdateTypeTestCategoryRequest) (*dto.TypeTestCategoryResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}
func (s *cachedTypeTestCategoryService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}
func (s *cachedTypeTestCategoryService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}
func (s *cachedTypeTestCategoryService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedTypeTestCategoryService) Activate(id uint) error {
	if err := s.inner.Activate(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedTypeTestCategoryService) Deactivate(id uint) error {
	if err := s.inner.Deactivate(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func normalizeCachePart(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *cachedTypeTestCategoryService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedTypeTestCategoryService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternTypeTestCategoryAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternTypeTestCategoryAll, err)
	}
}
