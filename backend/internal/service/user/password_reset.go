package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
)

type passwordResetChallenge struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	ResetCode string    `json:"reset_code"`
	ExpiresAt time.Time `json:"expires_at"`
}

type passwordResetSession struct {
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	ResetToken string    `json:"reset_token"`
	ExpiresAt  time.Time `json:"expires_at"`
}

func (s *userService) RequestPasswordReset(email string) (*dto.PasswordResetRequestResponse, error) {
	return s.createPasswordResetChallenge(email)
}

func (s *userService) ResendPasswordResetCode(email string) (*dto.PasswordResetRequestResponse, error) {
	return s.createPasswordResetChallenge(email)
}

func (s *userService) VerifyResetCode(email, code string) (*dto.PasswordResetTokenResponse, error) {
	if s.redis == nil {
		return nil, errors.New("password reset service is unavailable")
	}

	normalizedEmail := normalizeResetEmail(email)
	key := cache.PasswordResetChallengeKey(normalizedEmail)

	var challenge passwordResetChallenge
	if err := s.redis.Get(context.Background(), key, &challenge); err != nil {
		return nil, errors.New("reset code is invalid or expired")
	}

	if challenge.ResetCode != code {
		return nil, errors.New("reset code is invalid or expired")
	}

	resetToken, err := utils.GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}

	session := passwordResetSession{
		UserID:     challenge.UserID,
		Username:   challenge.Username,
		Email:      challenge.Email,
		ResetToken: resetToken,
		ExpiresAt:  challenge.ExpiresAt,
	}

	remainingTTL := time.Until(challenge.ExpiresAt)
	if remainingTTL <= 0 {
		return nil, errors.New("reset code is invalid or expired")
	}

	if err := s.redis.Set(context.Background(), cache.PasswordResetTokenKey(resetToken), session, remainingTTL); err != nil {
		return nil, err
	}

	_ = s.redis.Delete(context.Background(), key)

	return &dto.PasswordResetTokenResponse{
		ResetToken: resetToken,
		ExpiresIn:  int(remainingTTL.Seconds()),
	}, nil
}

func (s *userService) ResetPasswordWithToken(resetToken, newPassword string) error {
	if s.redis == nil {
		return errors.New("password reset service is unavailable")
	}

	var session passwordResetSession
	if err := s.redis.Get(context.Background(), cache.PasswordResetTokenKey(resetToken), &session); err != nil {
		return errors.New("reset token is invalid or expired")
	}

	user, err := s.repo.FindByID(session.UserID)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	if err := s.repo.Update(user); err != nil {
		return err
	}

	_ = s.redis.Delete(context.Background(), cache.UserKey(user.ID))
	_ = s.redis.Delete(context.Background(), cache.PasswordResetTokenKey(resetToken))
	return nil
}

func (s *userService) createPasswordResetChallenge(email string) (*dto.PasswordResetRequestResponse, error) {
	if s.redis == nil {
		return nil, errors.New("password reset service is unavailable")
	}

	normalizedEmail := normalizeResetEmail(email)
	user, err := s.repo.FindByEmail(normalizedEmail)
	if err != nil {
		// Hindari user enumeration: tetap return sukses tanpa membocorkan status akun.
		return nil, nil
	}

	resetCode, err := utils.GenerateNumericCode(6)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().UTC().Add(s.passwordResetTTL())
	challenge := passwordResetChallenge{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		ResetCode: resetCode,
		ExpiresAt: expiresAt,
	}

	if err := s.redis.Set(context.Background(), cache.PasswordResetChallengeKey(normalizedEmail), challenge, s.passwordResetTTL()); err != nil {
		return nil, err
	}

	return &dto.PasswordResetRequestResponse{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		ResetCode: resetCode,
		ExpiresIn: int(s.passwordResetTTL().Seconds()),
	}, nil
}

func (s *userService) marshalPasswordResetChallenge(challenge passwordResetChallenge) ([]byte, error) {
	return json.Marshal(challenge)
}

func (s *userService) unmarshalPasswordResetChallenge(data []byte, challenge *passwordResetChallenge) error {
	return json.Unmarshal(data, challenge)
}

func (s *userService) unmarshalPasswordResetSession(data []byte, session *passwordResetSession) error {
	return json.Unmarshal(data, session)
}

func (s *userService) debugPasswordResetChallenge(challenge passwordResetChallenge) string {
	return fmt.Sprintf("user=%d email=%s expires_at=%s", challenge.UserID, challenge.Email, challenge.ExpiresAt.Format(time.RFC3339))
}