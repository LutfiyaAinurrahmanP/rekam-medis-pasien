package patient

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

// cachedPatientService membungkus PatientService dengan Redis caching.
//
// Best practices yang diterapkan:
//   - GetMyPatientData di-cache per userID (private data)
//   - UpdateMyPatientData invalidasi hanya entry milik userID yang bersangkutan + seluruh list
//   - Semua write admin → invalidasi seluruh namespace patient:*
type cachedPatientService struct {
	inner PatientService
	redis *cache.RedisClient
}

// NewCachedPatientService returns a PatientService with Redis caching.
// Jika redisClient nil, langsung kembalikan inner service tanpa cache.
func NewCachedPatientService(inner PatientService, redisClient *cache.RedisClient) PatientService {
	if redisClient == nil {
		return inner
	}
	return &cachedPatientService{inner: inner, redis: redisClient}
}

// ─── Read operations (dengan cache) ────────────────────────────────────────

func (s *cachedPatientService) ListPatients(query *dto.PatientPaginationQuery) (*dto.PatientListResponse, error) {
	key := cache.PatientListQueryKey(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.Gender),
		normalizeCachePart(query.BloodType),
		normalizeCachePart(query.InsuranceProvider),
		query.MinAge,
		query.MaxAge,
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.PatientListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.ListPatients(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedPatientService) DeleteListPatients(query *dto.PatientPaginationQuery) (*dto.PatientDeletedListResponse, error) {
	key := cache.PatientDeletedListQueryKey(
		query.Page,
		query.PageSize,
		normalizeCachePart(query.Search),
		normalizeCachePart(query.Gender),
		normalizeCachePart(query.BloodType),
		normalizeCachePart(query.InsuranceProvider),
		query.MinAge,
		query.MaxAge,
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.PatientDeletedListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.DeleteListPatients(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedPatientService) GetPatientByID(id uint) (*dto.PatientResponse, error) {
	key := cache.PatientKey(id)
	var resp dto.PatientResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetPatientByID(id)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedPatientService) GetPatientByCode(code string) (*dto.PatientResponse, error) {
	key := cache.PatientByCodeKey(code)
	var resp dto.PatientResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetPatientByCode(code)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedPatientService) GetMyPatientData(userID uint) (*dto.PatientResponse, error) {
	key := cache.PatientByUserIDKey(userID)
	var resp dto.PatientResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetMyPatientData(userID)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

// ─── Write operations (invalidate cache) ───────────────────────────────────

func (s *cachedPatientService) UpdateMyPatientData(userID uint, req *dto.UpdatePatientRequest) (*dto.PatientResponse, error) {
	result, err := s.inner.UpdateMyPatientData(userID, req)
	if err != nil {
		return nil, err
	}
	// Invalidasi entry milik user ini + semua list (data bisa muncul di list)
	_ = s.redis.Delete(context.Background(), cache.PatientByUserIDKey(userID))
	s.invalidateListAndID(result.ID)
	return result, nil
}

func (s *cachedPatientService) CreatePatient(req *dto.CreatePatientRequest) (*dto.PatientResponse, error) {
	result, err := s.inner.CreatePatient(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedPatientService) UpdatePatient(id uint, req *dto.UpdatePatientRequest) (*dto.PatientResponse, error) {
	result, err := s.inner.UpdatePatient(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedPatientService) SoftDeletePatient(id uint) error {
	if err := s.inner.SoftDeletePatient(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedPatientService) RestorePatient(id uint) error {
	if err := s.inner.RestorePatient(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedPatientService) HardDeletePatient(id uint) error {
	if err := s.inner.HardDeletePatient(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

// ─── helpers ───────────────────────────────────────────────────────────────

func (s *cachedPatientService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedPatientService) invalidateListAndID(id uint) {
	_ = s.redis.Delete(context.Background(), cache.PatientKey(id))
	if err := s.redis.DeleteByPattern(context.Background(), "patient:list:*"); err != nil {
		log.Printf("⚠️  Redis invalidate patient list failed: %v", err)
	}
}

func (s *cachedPatientService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternPatientAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternPatientAll, err)
	}
}

func normalizeCachePart(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

// compile-time interface check
var _ PatientService = (*cachedPatientService)(nil)
