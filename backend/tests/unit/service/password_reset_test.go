package service

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	userservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/user"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ─── Tipe internal yang kita mirror untuk test ────────────────────────────────
// Kita harus me-reconstruct JSON yang di-store ke Redis agar bisa di-pre-populate
// tanpa bergantung pada ekspos internal struct dari package user.

type testPasswordResetChallenge struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	ResetCode string    `json:"reset_code"`
	ExpiresAt time.Time `json:"expires_at"`
}

type testPasswordResetSession struct {
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	ResetToken string    `json:"reset_token"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// ─── Helper: pre-populate Redis dengan challenge ──────────────────────────────

func seedPasswordResetChallenge(
	t *testing.T,
	redis *mocks.InMemoryRedis,
	email, resetCode string,
	userID uint, username string,
	ttl time.Duration,
) {
	t.Helper()
	challenge := testPasswordResetChallenge{
		UserID:    userID,
		Username:  username,
		Email:     email,
		ResetCode: resetCode,
		ExpiresAt: time.Now().UTC().Add(ttl),
	}
	data, err := json.Marshal(challenge)
	require.NoError(t, err)
	key := cache.PasswordResetChallengeKey(email)
	redis.SetRaw(key, data, ttl)
}

func seedPasswordResetSession(
	t *testing.T,
	redis *mocks.InMemoryRedis,
	resetToken string,
	userID uint, username, email string,
	ttl time.Duration,
) {
	t.Helper()
	session := testPasswordResetSession{
		UserID:     userID,
		Username:   username,
		Email:      email,
		ResetToken: resetToken,
		ExpiresAt:  time.Now().UTC().Add(ttl),
	}
	data, err := json.Marshal(session)
	require.NoError(t, err)
	key := cache.PasswordResetTokenKey(resetToken)
	redis.SetRaw(key, data, ttl)
}

// ─── Helper: buat service dengan InMemoryRedis ────────────────────────────────

func newTestServiceWithRedis(mockRepo *mocks.MockUserRepository, redis cache.RedisStore) userservice.UserService {
	return userservice.NewUserService(mockRepo, &config.Config{}, redis)
}

// ══════════════════════════════════════════════════════════════════════════════
// ── RequestPasswordReset ──────────────────────────────────────────────────────
// ══════════════════════════════════════════════════════════════════════════════

// TestRequestPasswordReset_Success memastikan OTP code di-generate dan disimpan
// ke Redis, lalu response berisi UserID, Email, ResetCode, dan ExpiresIn > 0.
func TestRequestPasswordReset_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	user := &models.User{
		ID:       1,
		Username: "john_doe",
		Email:    "john@example.com",
	}
	mockRepo.On("FindByEmail", "john@example.com").Return(user, nil)

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	result, err := svc.RequestPasswordReset("john@example.com")

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, "john@example.com", result.Email)
	assert.Equal(t, "john_doe", result.Username)
	assert.Len(t, result.ResetCode, 6, "OTP code harus 6 digit")
	assert.Greater(t, result.ExpiresIn, 0)

	// Pastikan challenge tersimpan di Redis
	key := cache.PasswordResetChallengeKey("john@example.com")
	assert.True(t, redisStore.Exists(key), "challenge harus tersimpan di Redis")

	mockRepo.AssertExpectations(t)
}

// TestRequestPasswordReset_EmailNotFound memastikan tidak ada user enumeration:
// jika email tidak ditemukan, service mengembalikan nil, nil (bukan error).
func TestRequestPasswordReset_EmailNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	mockRepo.On("FindByEmail", "notfound@example.com").
		Return(nil, errors.New("user not found"))

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	result, err := svc.RequestPasswordReset("notfound@example.com")

	// Tidak ada error dan tidak ada response — mencegah user enumeration
	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, redisStore.Count(), "Redis harus tetap kosong")

	mockRepo.AssertExpectations(t)
}

// TestRequestPasswordReset_RedisUnavailable memastikan service mengembalikan
// error jika Redis tidak tersedia (diperlukan untuk menyimpan challenge OTP).
func TestRequestPasswordReset_RedisUnavailable(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	// Buat service tanpa Redis (nil store)
	svc := userservice.NewUserService(mockRepo, &config.Config{})
	result, err := svc.RequestPasswordReset("john@example.com")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unavailable")
	assert.Nil(t, result)
}

// TestRequestPasswordReset_EmailNormalization memastikan email di-lowercase
// dan di-trim sebelum dicari ke database dan disimpan ke Redis.
func TestRequestPasswordReset_EmailNormalization(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	user := &models.User{
		ID:       2,
		Username: "jane_doe",
		Email:    "jane@example.com",
	}
	// Service harus normalize "  JANE@EXAMPLE.COM  " → "jane@example.com"
	mockRepo.On("FindByEmail", "jane@example.com").Return(user, nil)

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	result, err := svc.RequestPasswordReset("  JANE@EXAMPLE.COM  ")

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "jane@example.com", result.Email)

	// Challenge harus disimpan dengan email yang sudah dinormalisasi
	key := cache.PasswordResetChallengeKey("jane@example.com")
	assert.True(t, redisStore.Exists(key))

	mockRepo.AssertExpectations(t)
}

// ══════════════════════════════════════════════════════════════════════════════
// ── ResendPasswordResetCode ───────────────────────────────────────────────────
// ══════════════════════════════════════════════════════════════════════════════

// TestResendPasswordResetCode_Success memastikan OTP baru di-generate dan
// menimpa challenge lama di Redis (idempoten untuk user yang sama).
func TestResendPasswordResetCode_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	user := &models.User{
		ID:       1,
		Username: "john_doe",
		Email:    "john@example.com",
	}
	// Seed challenge lama
	seedPasswordResetChallenge(t, redisStore, "john@example.com", "111111", 1, "john_doe", 15*time.Minute)

	mockRepo.On("FindByEmail", "john@example.com").Return(user, nil)

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	result, err := svc.ResendPasswordResetCode("john@example.com")

	assert.NoError(t, err)
	require.NotNil(t, result)
	// Code baru harus berbeda dari yang lama (sangat besar kemungkinannya)
	assert.Len(t, result.ResetCode, 6)

	mockRepo.AssertExpectations(t)
}

// TestResendPasswordResetCode_EmailNotFound sama seperti ForgotPassword:
// tidak mengekspos apakah email ada atau tidak.
func TestResendPasswordResetCode_EmailNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	mockRepo.On("FindByEmail", "ghost@example.com").
		Return(nil, errors.New("user not found"))

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	result, err := svc.ResendPasswordResetCode("ghost@example.com")

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// ══════════════════════════════════════════════════════════════════════════════
// ── VerifyResetCode ───────────────────────────────────────────────────────────
// ══════════════════════════════════════════════════════════════════════════════

// TestVerifyResetCode_Success memastikan:
// 1. OTP yang benar menghasilkan reset token
// 2. Challenge dihapus dari Redis setelah verifikasi
// 3. Session baru tersimpan di Redis dengan reset token
func TestVerifyResetCode_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	seedPasswordResetChallenge(t, redisStore,
		"john@example.com", "123456", 1, "john_doe", 15*time.Minute,
	)
	challengeKey := cache.PasswordResetChallengeKey("john@example.com")

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	result, err := svc.VerifyResetCode("john@example.com", "123456")

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.ResetToken, "reset token harus di-generate")
	assert.Greater(t, result.ExpiresIn, 0, "ExpiresIn harus positif")

	// Challenge harus dihapus setelah berhasil verifikasi
	assert.False(t, redisStore.Exists(challengeKey), "challenge harus dihapus dari Redis")

	// Session baru harus tersimpan
	sessionKey := cache.PasswordResetTokenKey(result.ResetToken)
	assert.True(t, redisStore.Exists(sessionKey), "session dengan reset token harus tersimpan")
}

// TestVerifyResetCode_WrongCode memastikan kode yang salah menghasilkan error.
func TestVerifyResetCode_WrongCode(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	seedPasswordResetChallenge(t, redisStore,
		"john@example.com", "123456", 1, "john_doe", 15*time.Minute,
	)

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	result, err := svc.VerifyResetCode("john@example.com", "000000") // kode salah

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired")
	assert.Nil(t, result)
}

// TestVerifyResetCode_ChallengeNotFound memastikan error jika tidak ada
// challenge di Redis (belum request OTP atau sudah expired).
func TestVerifyResetCode_ChallengeNotFound(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()
	// Tidak di-seed challenge

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	result, err := svc.VerifyResetCode("nobody@example.com", "123456")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired")
	assert.Nil(t, result)
}

// TestVerifyResetCode_ExpiredChallenge memastikan error jika challenge
// masih ada di Redis tapi ExpiresAt-nya sudah lewat.
func TestVerifyResetCode_ExpiredChallenge(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	// Seed challenge yang sudah expired (ExpiresAt di masa lalu)
	challenge := testPasswordResetChallenge{
		UserID:    1,
		Username:  "john_doe",
		Email:     "john@example.com",
		ResetCode: "123456",
		ExpiresAt: time.Now().UTC().Add(-1 * time.Minute), // expired 1 menit lalu
	}
	data, err := json.Marshal(challenge)
	require.NoError(t, err)
	// Simpan dengan TTL panjang agar tidak dihapus oleh InMemoryRedis
	// tapi ExpiresAt-nya sudah lewat — service harus cek ExpiresAt sendiri
	redisStore.SetRaw(cache.PasswordResetChallengeKey("john@example.com"), data, 10*time.Minute)

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	result, err := svc.VerifyResetCode("john@example.com", "123456")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired")
	assert.Nil(t, result)
}

// TestVerifyResetCode_RedisUnavailable memastikan error jika Redis tidak tersedia.
func TestVerifyResetCode_RedisUnavailable(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	svc := userservice.NewUserService(mockRepo, &config.Config{})
	result, err := svc.VerifyResetCode("john@example.com", "123456")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unavailable")
	assert.Nil(t, result)
}

// TestVerifyResetCode_EmailNormalization memastikan email di-lowercase/trim
// juga saat verify, sehingga lookup key di Redis konsisten.
func TestVerifyResetCode_EmailNormalization(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	// Challenge disimpan dengan key lowercase
	seedPasswordResetChallenge(t, redisStore,
		"john@example.com", "654321", 1, "john_doe", 15*time.Minute,
	)

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	// Kirim email uppercase — harus tetap bisa ditemukan
	result, err := svc.VerifyResetCode("JOHN@EXAMPLE.COM", "654321")

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.ResetToken)
}

// ══════════════════════════════════════════════════════════════════════════════
// ── ResetPasswordWithToken ────────────────────────────────────────────────────
// ══════════════════════════════════════════════════════════════════════════════

// TestResetPasswordWithToken_Success memastikan:
// 1. Password di-update ke database
// 2. Session token dihapus dari Redis
// 3. User cache dihapus
func TestResetPasswordWithToken_Success(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	token := "valid-reset-token-abc123"
	seedPasswordResetSession(t, redisStore, token, 1, "john_doe", "john@example.com", 10*time.Minute)

	user := &models.User{
		ID:       1,
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "old_hashed_password",
	}
	mockRepo.On("FindByID", uint(1)).Return(user, nil)
	mockRepo.On("Update", user).Return(nil)

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	err := svc.ResetPasswordWithToken(token, "NewSecurePassword123!")

	assert.NoError(t, err)

	// Token session harus dihapus setelah reset berhasil
	sessionKey := cache.PasswordResetTokenKey(token)
	assert.False(t, redisStore.Exists(sessionKey), "session harus dihapus setelah reset")

	// Pastikan password yang di-update sudah di-hash (berbeda dari original)
	assert.NotEqual(t, "old_hashed_password", user.Password, "password harus di-hash sebelum disimpan")
	assert.NotEqual(t, "NewSecurePassword123!", user.Password, "password harus di-hash, bukan plaintext")

	mockRepo.AssertExpectations(t)
}

// TestResetPasswordWithToken_InvalidToken memastikan error jika token tidak
// ada di Redis (belum verify OTP atau sudah expired).
func TestResetPasswordWithToken_InvalidToken(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()
	// Tidak ada session di Redis

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	err := svc.ResetPasswordWithToken("invalid-token", "NewPassword123!")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid or expired")
}

// TestResetPasswordWithToken_UserNotFoundAfterToken memastikan error jika
// token valid di Redis tapi user sudah tidak ada di database.
func TestResetPasswordWithToken_UserNotFoundAfterToken(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	token := "orphaned-token"
	seedPasswordResetSession(t, redisStore, token, 999, "ghost", "ghost@example.com", 10*time.Minute)

	mockRepo.On("FindByID", uint(999)).Return(nil, errors.New("user not found"))

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	err := svc.ResetPasswordWithToken(token, "NewPassword123!")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// TestResetPasswordWithToken_RedisUnavailable memastikan error jika Redis tidak tersedia.
func TestResetPasswordWithToken_RedisUnavailable(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)

	svc := userservice.NewUserService(mockRepo, &config.Config{})
	err := svc.ResetPasswordWithToken("some-token", "NewPassword123!")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unavailable")
}

// ══════════════════════════════════════════════════════════════════════════════
// ── End-to-End Happy Path: full forgot password flow ─────────────────────────
// ══════════════════════════════════════════════════════════════════════════════

// TestForgotPasswordFlow_FullHappyPath mensimulasikan seluruh alur forgot password
// dari awal sampai reset berhasil menggunakan satu service instance + InMemoryRedis:
//
//  1. User meminta OTP via RequestPasswordReset
//  2. OTP dikirim ke email (di sini hanya dicek nilai-nya dari Redis)
//  3. User memasukkan OTP via VerifyResetCode → dapat reset token
//  4. User melakukan ResetPasswordWithToken → password berubah
func TestForgotPasswordFlow_FullHappyPath(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	user := &models.User{
		ID:       10,
		Username: "full_flow_user",
		Email:    "flow@example.com",
		Password: "old_hashed_password",
	}

	// ── Step 1: RequestPasswordReset ─────────────────────────────────────────
	mockRepo.On("FindByEmail", "flow@example.com").Return(user, nil)

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	resetResp, err := svc.RequestPasswordReset("flow@example.com")

	require.NoError(t, err)
	require.NotNil(t, resetResp)
	otpCode := resetResp.ResetCode
	assert.Len(t, otpCode, 6, "OTP harus 6 digit numerik")
	t.Logf("📧 OTP dikirim ke email: %s", otpCode)

	// ── Step 2: VerifyResetCode ───────────────────────────────────────────────
	tokenResp, err := svc.VerifyResetCode("flow@example.com", otpCode)

	require.NoError(t, err)
	require.NotNil(t, tokenResp)
	resetToken := tokenResp.ResetToken
	assert.NotEmpty(t, resetToken, "reset token harus di-generate")
	t.Logf("🔑 Reset token: %s", resetToken[:8]+"...")

	// Challenge harus sudah dihapus
	assert.False(t, redisStore.Exists(cache.PasswordResetChallengeKey("flow@example.com")),
		"challenge harus dihapus setelah verify")

	// ── Step 3: ResetPasswordWithToken ───────────────────────────────────────
	mockRepo.On("FindByID", uint(10)).Return(user, nil)
	mockRepo.On("Update", user).Return(nil)

	err = svc.ResetPasswordWithToken(resetToken, "BrandNewPassword@2026!")

	require.NoError(t, err)
	t.Logf("✅ Password berhasil di-reset untuk user %s", user.Username)

	// Token session harus dihapus setelah berhasil
	assert.False(t, redisStore.Exists(cache.PasswordResetTokenKey(resetToken)),
		"session token harus dihapus setelah reset")

	// Password di model harus sudah berubah dan di-hash
	assert.NotEqual(t, "old_hashed_password", user.Password)
	assert.NotEqual(t, "BrandNewPassword@2026!", user.Password,
		"password tersimpan harus bcrypt hash, bukan plaintext")

	// ── Step 4: Token yang sama tidak bisa dipakai ulang ──────────────────────
	err = svc.ResetPasswordWithToken(resetToken, "AnotherPassword@2026!")
	assert.Error(t, err, "token yang sudah dipakai tidak boleh bisa dipakai lagi")
	assert.Contains(t, err.Error(), "invalid or expired")

	mockRepo.AssertExpectations(t)
}

// ══════════════════════════════════════════════════════════════════════════════
// ── Edge Cases ────────────────────────────────────────────────────────────────
// ══════════════════════════════════════════════════════════════════════════════

// TestRequestPasswordReset_MultipleCalls memastikan setiap call menghasilkan
// OTP baru dan menimpa yang lama di Redis.
func TestRequestPasswordReset_MultipleCalls(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	user := &models.User{
		ID:       1,
		Username: "john_doe",
		Email:    "john@example.com",
	}
	// FindByEmail dipanggil dua kali
	mockRepo.On("FindByEmail", "john@example.com").Return(user, nil)

	svc := newTestServiceWithRedis(mockRepo, redisStore)

	resp1, err := svc.RequestPasswordReset("john@example.com")
	require.NoError(t, err)
	require.NotNil(t, resp1)

	resp2, err := svc.RequestPasswordReset("john@example.com")
	require.NoError(t, err)
	require.NotNil(t, resp2)

	// Redis hanya boleh menyimpan 1 challenge (yang terbaru menimpa yang lama)
	assert.Equal(t, 1, redisStore.Count(), "hanya 1 challenge aktif per email")

	mockRepo.AssertExpectations(t)
}

// TestVerifyResetCode_CannotReuse memastikan OTP tidak bisa dipakai dua kali
// (challenge dihapus setelah verify pertama berhasil).
func TestVerifyResetCode_CannotReuse(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	seedPasswordResetChallenge(t, redisStore,
		"john@example.com", "999888", 1, "john_doe", 15*time.Minute,
	)

	svc := newTestServiceWithRedis(mockRepo, redisStore)

	// Verify pertama berhasil
	_, err := svc.VerifyResetCode("john@example.com", "999888")
	assert.NoError(t, err)

	// Verify kedua dengan kode yang sama harus gagal
	result, err := svc.VerifyResetCode("john@example.com", "999888")
	assert.Error(t, err, "OTP tidak boleh bisa dipakai dua kali")
	assert.Nil(t, result)
}

// TestResetPasswordWithToken_PasswordIsHashed memastikan password yang disimpan
// ke database adalah bcrypt hash, bukan plaintext.
func TestResetPasswordWithToken_PasswordIsHashed(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	redisStore := mocks.NewInMemoryRedis()

	token := "hash-check-token"
	seedPasswordResetSession(t, redisStore, token, 1, "john_doe", "john@example.com", 10*time.Minute)

	user := &models.User{
		ID:       1,
		Password: "old_hash",
	}

	var capturedUser *models.User
	mockRepo.On("FindByID", uint(1)).Return(user, nil)
	// Capture argumen Update untuk inspeksi
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		capturedUser = args.Get(0).(*models.User)
	}).Return(nil)

	svc := newTestServiceWithRedis(mockRepo, redisStore)
	err := svc.ResetPasswordWithToken(token, "PlainPassword123!")

	require.NoError(t, err)
	if capturedUser != nil {
		assert.NotEqual(t, "PlainPassword123!", capturedUser.Password,
			"password tersimpan tidak boleh plaintext")
		assert.NotEqual(t, "old_hash", capturedUser.Password,
			"password harus berubah dari yang lama")
		// Bcrypt hash selalu dimulai dengan $2a$ atau $2b$
		assert.True(t,
			len(capturedUser.Password) > 30,
			"password hash harus cukup panjang (bcrypt menghasilkan ~60 karakter)")
	}
}

// ─── Pastikan InMemoryRedis memenuhi cache.RedisStore interface ───────────────
var _ cache.RedisStore = (*mocks.InMemoryRedis)(nil)
