package service

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
)

type UserService interface {
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
	ListUsers(query *dto.PaginationQuery) (*dto.UserListResponse, error)
	DeleteListUsers(query *dto.PaginationQuery) (*dto.DeletedUserListResponse, error)
	UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	SoftDeleteUser(id uint) error
	HardDeleteUser(id uint) error
	RestoreUser(id uint) error
	ChangePassword(id uint, req *dto.ChangePasswordRequest) error
	ResetPassword(id uint, newPassword string) error
	ActivateUser(id uint) error
	DeactivateUser(id uint) error
	VerifyPasswordForDeletion(id uint, password string) error
}

type userService struct {
	repo   repository.UserRepository
	config *config.Config
}

func NewUserService(repo repository.UserRepository, config *config.Config) UserService {
	return &userService{
		repo:   repo,
		config: config,
	}
}

func (s *userService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	exists, err := s.repo.IsUsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	exists, err = s.repo.IsEmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	exists, err = s.repo.IsPhoneExists(req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("phone already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = models.RolePatient
	}

	if !models.ValidateRole(role) {
		return nil, errors.New("invalid role")
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     role,
		IsActive: true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *userService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByUsernameOrEmail(req.UsernameOrEmail)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("user account is inactive")
	}

	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	token, expiresAt, err := utils.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		user.Role,
		s.config.JWT.Secret,
		s.config.JWT.ExpiredTime,
	)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *s.toUserResponse(user),
	}, nil
}

func (s *userService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	exists, err := s.repo.IsUsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	exists, err = s.repo.IsEmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	exists, err = s.repo.IsPhoneExists(req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("phone already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = models.RolePatient
	}

	if !models.ValidateRole(role) {
		return nil, errors.New("invalid role")
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     role,
		IsActive: true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *userService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toUserResponse(user), nil
}

func (s *userService) ListUsers(query *dto.PaginationQuery) (*dto.UserListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = s.config.Pagination.DefaultPageSize
	}

	if query.PageSize > s.config.Pagination.MaxPageSize {
		query.PageSize = s.config.Pagination.MaxPageSize
	}

	if query.SortBy == "" {
		query.SortBy = "created_at"
	}

	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	if query.Role != "" && !models.ValidateRole(query.Role) {
		return nil, errors.New("invalid role filter")
	}

	users, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.toUserResponse(&user)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.UserListResponse{
		Data: userResponses,
		Meta: dto.PaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *userService) DeleteListUsers(query *dto.PaginationQuery) (*dto.DeletedUserListResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = s.config.Pagination.DefaultPageSize
	}

	if query.PageSize > s.config.Pagination.MaxPageSize {
		query.PageSize = s.config.Pagination.MaxPageSize
	}

	if query.SortBy == "" {
		query.SortBy = "deleted_at"
	}

	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	if query.Role != "" && !models.ValidateRole(query.Role) {
		return nil, errors.New("invalid role filter")
	}

	deletedUsers, total, err := s.repo.DeleteList(query)
	if err != nil {
		return nil, err
	}

	deletedUsersReponses := make([]dto.DeletedUserResponse, len(deletedUsers))
	for i, deletedUser := range deletedUsers {
		deletedUsersReponses[i] = *s.toDeletedUserResponse(&deletedUser)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.DeletedUserListResponse{
		Data: deletedUsersReponses,
		Meta: dto.PaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil

}

func (s *userService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Username != nil && *req.Username != user.Username {
		exists, err := s.repo.IsUsernameExists(*req.Username, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("username already exists")
		}
		user.Username = *req.Username
	}

	if req.Email != nil && *req.Email != user.Email {
		exists, err := s.repo.IsEmailExists(*req.Email, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already exists")
		}
		user.Email = *req.Email
	}

	if req.Phone != nil && *req.Phone != user.Phone {
		exists, err := s.repo.IsPhoneExists(*req.Phone, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("phone already exists")
		}
		user.Phone = *req.Phone
	}

	if req.Password != nil {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if req.Role != nil {
		if !models.ValidateRole(*req.Role) {
			return nil, errors.New("invalid role")
		}
		user.Role = *req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil

}

func (s *userService) ChangePassword(id uint, req *dto.ChangePasswordRequest) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(user.Password, req.OldPassword) {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.Update(user)
}

func (s *userService) SoftDeleteUser(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(id)
}

func (s *userService) HardDeleteUser(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *userService) RestoreUser(id uint) error {
	return s.repo.Restore(id)
}
func (s *userService) ResetPassword(id uint, newPassword string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.Update(user)
}

func (s *userService) ActivateUser(id uint) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	user.IsActive = true
	return s.repo.Update(user)
}

func (s *userService) DeactivateUser(id uint) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	user.IsActive = false
	return s.repo.Update(user)
}

func (s *userService) VerifyPasswordForDeletion(id uint, password string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(user.Password, password) {
		return errors.New("invalid password")
	}

	return nil
}

func (s *userService) toUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (s *userService) toDeletedUserResponse(user *models.User) *dto.DeletedUserResponse {
	return &dto.DeletedUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: &user.DeletedAt.Time,
	}
}
