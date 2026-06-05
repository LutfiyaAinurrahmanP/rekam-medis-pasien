package doctor

import (
	"context"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

// cachedDoctorService membungkus DoctorService dengan Redis caching.
//
// Best practices yang diterapkan:
//   - GetMyDoctorData di-cache per userID
//   - UpdateMyDoctorData invalidasi hanya doctor:user:{userID} + seluruh list
//   - ActivateDoctor / DeactivateDoctor merubah status → invalidasi seluruh namespace
//   - OccupyRoom/ReleaseRoom analog tidak ada di doctor, tapi doctor availability bisa berubah
type cachedDoctorService struct {
	inner DoctorService
	redis *cache.RedisClient
}

// NewCachedDoctorService returns a DoctorService with Redis caching.
// Jika redisClient nil, langsung kembalikan inner service tanpa cache.
func NewCachedDoctorService(inner DoctorService, redisClient *cache.RedisClient) DoctorService {
	if redisClient == nil {
		return inner
	}
	return &cachedDoctorService{inner: inner, redis: redisClient}
}

// ─── Read operations (dengan cache) ────────────────────────────────────────

func (s *cachedDoctorService) GetMyDoctorData(userID uint) (*dto.DoctorResponse, error) {
	key := cache.DoctorByUserIDKey(userID)
	var resp dto.DoctorResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetMyDoctorData(userID)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedDoctorService) ListDoctors(query *dto.DoctorPaginationQuery) (*dto.DoctorListResponse, error) {
	key := cache.DoctorListKey(query.Page, query.PageSize)
	var resp dto.DoctorListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.ListDoctors(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedDoctorService) DeletedListDoctors(query *dto.DoctorPaginationQuery) (*dto.DoctorDeletedListResponse, error) {
	// Deleted list tidak di-cache karena berubah sangat dinamis
	return s.inner.DeletedListDoctors(query)
}

func (s *cachedDoctorService) GetDoctorByID(id uint) (*dto.DoctorResponse, error) {
	key := cache.DoctorKey(id)
	var resp dto.DoctorResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetDoctorByID(id)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedDoctorService) GetDoctorBySpecializationID(specID uint) (*dto.DoctorResponse, error) {
	key := cache.DoctorBySpecKey(specID)
	var resp dto.DoctorResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetDoctorBySpecializationID(specID)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

// ─── Write operations — dokter memperbarui data sendiri ────────────────────

func (s *cachedDoctorService) UpdateMyDoctorData(userID uint, req *dto.UpdateDoctorRequest) (*dto.DoctorResponse, error) {
	result, err := s.inner.UpdateMyDoctorData(userID, req)
	if err != nil {
		return nil, err
	}
	// Invalidasi entry personal + list
	_ = s.redis.Delete(context.Background(), cache.DoctorByUserIDKey(userID))
	if result != nil {
		_ = s.redis.Delete(context.Background(), cache.DoctorKey(result.ID))
	}
	_ = s.redis.DeleteByPattern(context.Background(), "doctor:list:*")
	return result, nil
}

// ─── Write operations — admin ───────────────────────────────────────────────

func (s *cachedDoctorService) CreateDoctor(req *dto.CreateDoctorRequest) (*dto.DoctorResponse, error) {
	result, err := s.inner.CreateDoctor(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedDoctorService) UpdateDoctor(id uint, req *dto.UpdateDoctorRequest) (*dto.DoctorResponse, error) {
	result, err := s.inner.UpdateDoctor(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedDoctorService) ActivateDoctor(id uint) (*dto.DoctorResponse, error) {
	result, err := s.inner.ActivateDoctor(id)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedDoctorService) DeactivateDoctor(id uint) (*dto.DoctorResponse, error) {
	result, err := s.inner.DeactivateDoctor(id)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedDoctorService) SoftDeleteDoctor(id uint) error {
	if err := s.inner.SoftDeleteDoctor(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedDoctorService) RestoreDoctor(id uint) error {
	if err := s.inner.RestoreDoctor(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedDoctorService) HardDeleteDoctor(id uint) error {
	if err := s.inner.HardDeleteDoctor(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

// ─── helpers ───────────────────────────────────────────────────────────────

func (s *cachedDoctorService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedDoctorService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternDoctorAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternDoctorAll, err)
	}
}

// compile-time interface check
var _ DoctorService = (*cachedDoctorService)(nil)
