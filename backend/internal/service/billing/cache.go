package billing

import (
	"context"
	"strings"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedBillingService struct {
	inner BillingService
	redis *cache.RedisClient
}

func NewCachedBillingService(inner BillingService, redis *cache.RedisClient) BillingService {
	return &cachedBillingService{
		inner: inner,
		redis: redis,
	}
}

func normalizeCachePart(part string) string {
	if part == "" {
		return "all"
	}
	return strings.ToLower(part)
}

func (s *cachedBillingService) setCache(key string, data interface{}, expiration time.Duration) {
	s.redis.Set(context.Background(), key, data, expiration)
}

func (s *cachedBillingService) invalidateAll() {
	s.redis.DeleteByPattern(context.Background(), cache.PatternBillingAll)
}

func (s *cachedBillingService) List(query dto.BillingPaginationQuery) (*dto.BillingListResponse, error) {
	key := cache.BillingListQueryKey(
		query.Page, query.PageSize, query.PatientID, normalizeCachePart(query.Status), normalizeCachePart(query.PaymentMethod), normalizeCachePart(query.Search), normalizeCachePart(query.SortBy), normalizeCachePart(query.SortDir),
	)

	var resp dto.BillingListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.List(query)
	if err != nil {
		return nil, err
	}

	s.setCache(key, result, 5*time.Minute)
	return result, nil
}

func (s *cachedBillingService) DeletedList(query dto.BillingPaginationQuery) (*dto.BillingDeletedListResponse, error) {
	key := cache.BillingDeletedListQueryKey(
		query.Page, query.PageSize, query.PatientID, normalizeCachePart(query.Status), normalizeCachePart(query.PaymentMethod), normalizeCachePart(query.Search), normalizeCachePart(query.SortBy), normalizeCachePart(query.SortDir),
	)

	var resp dto.BillingDeletedListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.DeletedList(query)
	if err != nil {
		return nil, err
	}

	s.setCache(key, result, 5*time.Minute)
	return result, nil
}

func (s *cachedBillingService) FindByID(id uint) (*dto.BillingResponse, error) {
	key := cache.BillingKey(id)

	var resp dto.BillingResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.FindByID(id)
	if err != nil {
		return nil, err
	}

	s.setCache(key, result, 10*time.Minute)
	return result, nil
}

func (s *cachedBillingService) FindByIDUnscoped(id uint) (*dto.BillingResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *cachedBillingService) FindByInvoiceNumber(invoice string) (*dto.BillingResponse, error) {
	return s.inner.FindByInvoiceNumber(invoice)
}

func (s *cachedBillingService) FindByPatientID(patientID uint, query dto.BillingPaginationQuery) (*dto.BillingListResponse, error) {
	query.PatientID = &patientID
	return s.List(query)
}

func (s *cachedBillingService) Create(req dto.CreateBillingRequest) (*dto.BillingResponse, error) {
	result, err := s.inner.Create(req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedBillingService) Update(id uint, req dto.UpdateBillingRequest) (*dto.BillingResponse, error) {
	result, err := s.inner.Update(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedBillingService) RecordPayment(id uint, req dto.RecordPaymentRequest) (*dto.BillingResponse, error) {
	result, err := s.inner.RecordPayment(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedBillingService) Cancel(id uint) (*dto.BillingResponse, error) {
	result, err := s.inner.Cancel(id)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedBillingService) Delete(id uint) error {
	err := s.inner.Delete(id)
	if err == nil {
		s.invalidateAll()
	}
	return err
}

func (s *cachedBillingService) Restore(id uint) (*dto.BillingResponse, error) {
	result, err := s.inner.Restore(id)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedBillingService) HardDelete(id uint) error {
	err := s.inner.HardDelete(id)
	if err == nil {
		s.invalidateAll()
	}
	return err
}

func (s *cachedBillingService) ListItems(billingID uint, query dto.BillingItemPaginationQuery) ([]dto.BillingItemResponse, error) {
	return s.inner.ListItems(billingID, query)
}

func (s *cachedBillingService) FindItemByID(billingID, itemID uint) (*dto.BillingItemResponse, error) {
	return s.inner.FindItemByID(billingID, itemID)
}

func (s *cachedBillingService) CreateItem(billingID uint, req dto.CreateBillingItemRequest) (*dto.BillingItemResponse, error) {
	result, err := s.inner.CreateItem(billingID, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedBillingService) UpdateItem(billingID, itemID uint, req dto.UpdateBillingItemRequest) (*dto.BillingItemResponse, error) {
	result, err := s.inner.UpdateItem(billingID, itemID, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedBillingService) DeleteItem(billingID, itemID uint) error {
	err := s.inner.DeleteItem(billingID, itemID)
	if err == nil {
		s.invalidateAll()
	}
	return err
}
