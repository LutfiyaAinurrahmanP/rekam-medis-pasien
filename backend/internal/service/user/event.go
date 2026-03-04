package user

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

// eventUserService membungkus UserService dan mempublikasikan event Kafka
// setelah setiap operasi mutasi berhasil. Pola ini identik dengan cachedUserService.
type eventUserService struct {
	inner     UserService
	publisher kafka.EventPublisher
}

// NewEventUserService mengembalikan UserService yang dilengkapi event publishing.
// Jika publisher nil, inner service dikembalikan langsung tanpa dekorasi.
func NewEventUserService(inner UserService, publisher kafka.EventPublisher) UserService {
	if publisher == nil {
		return inner
	}
	return &eventUserService{inner: inner, publisher: publisher}
}

// ─── Read operations (delegasi langsung) ─────────────────────────────────────

func (s *eventUserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	return s.inner.GetUserByID(id)
}

func (s *eventUserService) ListUsers(query *dto.UserPaginationQuery) (*dto.UserListResponse, error) {
	return s.inner.ListUsers(query)
}

func (s *eventUserService) DeleteListUsers(query *dto.UserPaginationQuery) (*dto.UserDeletedListResponse, error) {
	return s.inner.DeleteListUsers(query)
}

func (s *eventUserService) VerifyPasswordForDeletion(id uint, password string) error {
	return s.inner.VerifyPasswordForDeletion(id, password)
}

// ─── Write operations (delegasi + publish event) ──────────────────────────────

func (s *eventUserService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	result, err := s.inner.Register(req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicUserRegistered,
		events.NewUserRegisteredEvent(result.ID, result.Username, result.Email, result.Phone, result.Role))
	return result, nil
}

func (s *eventUserService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	result, err := s.inner.Login(req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicUserLogin,
		events.NewUserLoginEvent(result.User.ID, result.User.Username, result.User.Email, result.User.Role))
	return result, nil
}

func (s *eventUserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	result, err := s.inner.CreateUser(req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicUserCreated,
		events.NewUserCreatedEvent(result.ID, result.Username, result.Email, result.Phone, result.Role, result.IsActive))
	return result, nil
}

func (s *eventUserService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	result, err := s.inner.UpdateUser(id, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicUserUpdated,
		events.NewUserUpdatedEvent(result.ID, result.Username, result.Email, result.Phone, result.Role, result.IsActive, "profile_update"))
	return result, nil
}

func (s *eventUserService) ChangePassword(id uint, req *dto.ChangePasswordRequest) error {
	if err := s.inner.ChangePassword(id, req); err != nil {
		return err
	}
	// Ambil data user untuk event (password baru tidak disertakan di event)
	if user, err := s.inner.GetUserByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicUserUpdated,
			events.NewUserUpdatedEvent(user.ID, user.Username, user.Email, user.Phone, user.Role, user.IsActive, "password_change"))
	}
	return nil
}

func (s *eventUserService) ResetPassword(id uint, newPassword string) error {
	if err := s.inner.ResetPassword(id, newPassword); err != nil {
		return err
	}
	if user, err := s.inner.GetUserByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicUserUpdated,
			events.NewUserUpdatedEvent(user.ID, user.Username, user.Email, user.Phone, user.Role, user.IsActive, "password_reset"))
	}
	return nil
}

func (s *eventUserService) ActivateUser(id uint) error {
	if err := s.inner.ActivateUser(id); err != nil {
		return err
	}
	if user, err := s.inner.GetUserByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicUserUpdated,
			events.NewUserUpdatedEvent(user.ID, user.Username, user.Email, user.Phone, user.Role, user.IsActive, "activate"))
	}
	return nil
}

func (s *eventUserService) DeactivateUser(id uint) error {
	if err := s.inner.DeactivateUser(id); err != nil {
		return err
	}
	if user, err := s.inner.GetUserByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicUserUpdated,
			events.NewUserUpdatedEvent(user.ID, user.Username, user.Email, user.Phone, user.Role, user.IsActive, "deactivate"))
	}
	return nil
}

func (s *eventUserService) SoftDeleteUser(id uint) error {
	// Ambil data sebelum dihapus agar event tetap mengandung informasi user
	user, _ := s.inner.GetUserByID(id)
	if err := s.inner.SoftDeleteUser(id); err != nil {
		return err
	}
	username := ""
	if user != nil {
		username = user.Username
	}
	s.publisher.PublishAsync(kafka.TopicUserDeleted,
		events.NewUserDeletedEvent(id, username, "soft_delete"))
	return nil
}

func (s *eventUserService) HardDeleteUser(id uint) error {
	// Ambil data sebelum hard delete karena setelah itu record hilang
	user, _ := s.inner.GetUserByID(id)
	if err := s.inner.HardDeleteUser(id); err != nil {
		return err
	}
	username := ""
	if user != nil {
		username = user.Username
	}
	s.publisher.PublishAsync(kafka.TopicUserDeleted,
		events.NewUserDeletedEvent(id, username, "hard_delete"))
	return nil
}

func (s *eventUserService) RestoreUser(id uint) error {
	if err := s.inner.RestoreUser(id); err != nil {
		return err
	}
	if user, err := s.inner.GetUserByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicUserRestored,
			events.NewUserRestoredEvent(user.ID, user.Username))
	}
	return nil
}
