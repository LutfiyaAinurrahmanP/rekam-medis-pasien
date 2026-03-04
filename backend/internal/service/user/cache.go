package user

import (
	"context"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

// cachedUserService membungkus UserService dengan Redis caching.
//
// Best practices yang diterapkan:
//   - Login & Register TIDAK di-cache (operasi auth sensitif, token unik tiap request)
//   - ChangePassword / ResetPassword / Verify → invalidasi cache user tersebut
//   - Semua write operation → invalidasi seluruh namespace user:*
type cachedUserService struct {
	inner UserService
	redis *cache.RedisClient
}

// NewCachedUserService returns a UserService with Redis caching.
// Jika redisClient nil, langsung kembalikan inner service tanpa cache.
func NewCachedUserService(inner UserService, redisClient *cache.RedisClient) UserService {
	if redisClient == nil {
		return inner
	}
	return &cachedUserService{inner: inner, redis: redisClient}
}

// ─── Auth — TIDAK di-cache ─────────────────────────────────────────────────

func (s *cachedUserService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	result, err := s.inner.Register(req)
	if err != nil {
		return nil, err
	}
	// Invalidasi list setelah user baru dibuat
	s.invalidateAll()
	return result, nil
}

func (s *cachedUserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Tidak di-cache: token JWT bersifat unik per sesi dan mengandung expiry
	return s.inner.Login(req)
}

// ─── Read operations (dengan cache) ────────────────────────────────────────

func (s *cachedUserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	key := cache.UserKey(id)
	var resp dto.UserResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedUserService) ListUsers(query *dto.UserPaginationQuery) (*dto.UserListResponse, error) {
	key := cache.UserListKey(query.Page, query.PageSize)
	var resp dto.UserListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.ListUsers(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedUserService) DeleteListUsers(query *dto.UserPaginationQuery) (*dto.UserDeletedListResponse, error) {
	key := cache.UserDeletedListKey(query.Page, query.PageSize)
	var resp dto.UserDeletedListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.DeleteListUsers(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

// ─── Write operations (invalidate cache) ───────────────────────────────────

func (s *cachedUserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	result, err := s.inner.CreateUser(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedUserService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	result, err := s.inner.UpdateUser(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedUserService) SoftDeleteUser(id uint) error {
	if err := s.inner.SoftDeleteUser(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedUserService) HardDeleteUser(id uint) error {
	if err := s.inner.HardDeleteUser(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedUserService) RestoreUser(id uint) error {
	if err := s.inner.RestoreUser(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedUserService) ActivateUser(id uint) error {
	if err := s.inner.ActivateUser(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedUserService) DeactivateUser(id uint) error {
	if err := s.inner.DeactivateUser(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedUserService) ChangePassword(id uint, req *dto.ChangePasswordRequest) error {
	if err := s.inner.ChangePassword(id, req); err != nil {
		return err
	}
	// Invalidasi entry spesifik — password berubah, data user berubah
	_ = s.redis.Delete(context.Background(), cache.UserKey(id))
	return nil
}

func (s *cachedUserService) ResetPassword(id uint, newPassword string) error {
	if err := s.inner.ResetPassword(id, newPassword); err != nil {
		return err
	}
	_ = s.redis.Delete(context.Background(), cache.UserKey(id))
	return nil
}

func (s *cachedUserService) VerifyPasswordForDeletion(id uint, password string) error {
	// Operasi verifikasi — tidak mengubah data, tidak perlu invalidasi
	return s.inner.VerifyPasswordForDeletion(id, password)
}

// ─── helpers ───────────────────────────────────────────────────────────────

func (s *cachedUserService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedUserService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternUserAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternUserAll, err)
	}
}

// compile-time interface check
var _ UserService = (*cachedUserService)(nil)
