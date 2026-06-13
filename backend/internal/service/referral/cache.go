package referral

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedReferralService struct {
	inner ReferralService
	redis *cache.RedisClient
}

func NewCachedReferralService(inner ReferralService, redisClient *cache.RedisClient) ReferralService {
	if redisClient == nil {
		return inner
	}
	return &cachedReferralService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedReferralService) List(query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	key := cache.ReferralListQueryKey(
		query.Page, query.PageSize, query.PatientID, query.ReferringDoctorID, query.DoctorID,
		normalizeCachePart(query.ReferralType), normalizeCachePart(query.Status), normalizeCachePart(query.Priority), normalizeCachePart(query.Search), normalizeCachePart(query.SortBy), normalizeCachePart(query.SortDir),
	)

	var resp dto.ReferralListResponse
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

func (s *cachedReferralService) DeletedList(query dto.ReferralPaginationQuery) (*dto.ReferralDeletedListResponse, error) {
	key := cache.ReferralDeletedListQueryKey(
		query.Page, query.PageSize, query.PatientID, query.ReferringDoctorID, query.DoctorID,
		normalizeCachePart(query.ReferralType), normalizeCachePart(query.Status), normalizeCachePart(query.Priority), normalizeCachePart(query.Search), normalizeCachePart(query.SortBy), normalizeCachePart(query.SortDir),
	)

	var resp dto.ReferralDeletedListResponse
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

func (s *cachedReferralService) FindMyReferrals(patientID uint, status string) (*dto.ReferralMyListResponse, error) {
	return s.inner.FindMyReferrals(patientID, status)
}

func (s *cachedReferralService) FindByID(id uint) (*dto.ReferralResponse, error) {
	key := cache.ReferralKey(id)
	var resp dto.ReferralResponse
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

func (s *cachedReferralService) FindByIDUnscoped(id uint) (*dto.ReferralResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *cachedReferralService) FindByPatientID(patientID uint, query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	query.PatientID = &patientID
	return s.List(query)
}

func (s *cachedReferralService) FindByDoctorID(doctorID uint, query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	query.DoctorID = &doctorID
	return s.List(query)
}

func (s *cachedReferralService) Create(req dto.CreateReferralRequest) (*dto.ReferralResponse, error) {
	result, err := s.inner.Create(req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedReferralService) Update(id uint, req dto.UpdateReferralRequest) (*dto.ReferralResponse, error) {
	result, err := s.inner.Update(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedReferralService) Accept(id uint, req dto.AcceptReferralRequest) (*dto.ReferralResponse, error) {
	result, err := s.inner.Accept(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedReferralService) Reject(id uint, req dto.RejectReferralRequest) (*dto.ReferralResponse, error) {
	result, err := s.inner.Reject(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedReferralService) Complete(id uint, req dto.CompleteReferralRequest) (*dto.ReferralResponse, error) {
	result, err := s.inner.Complete(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedReferralService) Cancel(id uint, req dto.CancelReferralRequest) (*dto.ReferralResponse, error) {
	result, err := s.inner.Cancel(id, req)
	if err == nil {
		s.invalidateAll()
	}
	return result, err
}

func (s *cachedReferralService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedReferralService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedReferralService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedReferralService) setCache(key string, value interface{}) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("Failed to set cache for %s: %v", key, err)
	}
}

func (s *cachedReferralService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternReferralAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternReferralAll, err)
	}
}

func normalizeCachePart(s string) string {
	if s == "" {
		return "all"
	}
	return strings.ToLower(s)
}
