package service

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type DoctorService interface {
	GetMyDoctorData(userID uint) (*dto.DoctorResponse, error)
	UpdateMyDoctorData(userID uint, req *dto.UpdateDoctorRequest) (*dto.DoctorResponse, error)
	ListDoctors(query *dto.DoctorPaginationQuery) (*dto.DoctorListResponse, error)
	DeletedListDoctors(query *dto.DoctorPaginationQuery) (*dto.DoctorDeletedListResponse, error)
	GetDoctorByID(id uint) (*dto.DoctorResponse, error)
	GetDoctorBySpecialization(spec string) (*dto.DoctorResponse, error)
	CreateDoctor(req *dto.CreateDoctorRequest) (*dto.DoctorResponse, error)
	UpdateDoctor(id uint, req *dto.UpdateDoctorRequest) (*dto.DoctorResponse, error)
	ActivateDoctor(id uint) (*dto.DoctorResponse, error)
	DeactivateDoctor(id uint) (*dto.DoctorResponse, error)
	SoftDeleteDoctor(id uint) error
	RestoreDoctor(id uint) error
	HardDeleteDoctor(id uint) error
}

type doctorService struct {
	repo   repository.DoctorRepository
	config *config.Config
}

func NewDoctorService(repo repository.DoctorRepository, config *config.Config) DoctorService {
	return &doctorService{
		repo:   repo,
		config: config,
	}
}

func (s doctorService) GetMyDoctorData(userID uint) (*dto.DoctorResponse, error) {
	doctor, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.toDoctorResponse(doctor), nil
}

func (s doctorService) UpdateMyDoctorData(userID uint, req *dto.UpdateDoctorRequest) (*dto.DoctorResponse, error) {
	doctor, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Update(doctor); err != nil {
		return nil, err
	}

	if req.FullName != nil {
		doctor.FullName = *req.FullName
	}

	if req.Specialization != nil {
		doctor.Specialization = *req.Specialization
	}

	if req.Phone != nil {
		doctor.Phone = *req.Phone
	}

	if req.Email != nil {
		doctor.Email = *req.Email
	}

	if req.DepartmentID != nil {
		doctor.DepartmentID = req.DepartmentID
	}

	if req.IsActive != nil {
		doctor.IsActive = *req.IsActive
	}

	return s.toDoctorResponse(doctor), nil
}

func (s doctorService) ListDoctors(query *dto.DoctorPaginationQuery) (*dto.DoctorListResponse, error) {
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

	doctors, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	doctorResponses := make([]dto.DoctorResponse, len(doctors))
	for i, doctor := range doctors {
		doctorResponses[i] = *s.toDoctorResponse(&doctor)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.DoctorListResponse{
		Data: doctorResponses,
		Meta: dto.DoctorPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s doctorService) DeletedListDoctors(query *dto.DoctorPaginationQuery) (*dto.DoctorDeletedListResponse, error) {
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

	doctors, total, err := s.repo.DeleteList(query)
	if err != nil {
		return nil, err
	}

	doctorResponses := make([]dto.DeletedDoctorResponse, len(doctors))
	for i, doctor := range doctors {
		doctorResponses[i] = *s.toDeleteDoctorResponse(&doctor)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.DoctorDeletedListResponse{
		Data: doctorResponses,
		Meta: dto.DoctorPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s doctorService) GetDoctorByID(id uint) (*dto.DoctorResponse, error) {
	doctor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toDoctorResponse(doctor), nil
}

func (s doctorService) GetDoctorBySpecialization(specialization string) (*dto.DoctorResponse, error) {
	spec, err := s.repo.FindBySpecialization(specialization)
	if err != nil {
		return nil, err
	}
	return s.toDoctorResponse(spec), nil
}

func (s doctorService) CreateDoctor(req *dto.CreateDoctorRequest) (*dto.DoctorResponse, error) {
	exists, err := s.repo.IsEmployeeIDExists(req.EmployeeID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("employee already exists")
	}

	doctors := &models.Doctor{
		UserID:         req.UserID,
		EmployeeID:     req.EmployeeID,
		FullName:       req.FullName,
		Specialization: req.Specialization,
		LicenseNumber:  req.LicenseNumber,
		Phone:          req.Phone,
		Email:          req.Email,
		DepartmentID:   req.DepartmentID,
		IsActive:       true,
	}

	if err := s.repo.Create(doctors); err != nil {
		return nil, err
	}

	return s.toDoctorResponse(doctors), nil
}

func (s doctorService) UpdateDoctor(id uint, req *dto.UpdateDoctorRequest) (*dto.DoctorResponse, error) {
	doctor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.FullName != nil {
		doctor.FullName = *req.FullName
	}

	if req.Specialization != nil {
		doctor.Specialization = *req.Specialization
	}

	if req.Phone != nil {
		doctor.Phone = *req.Phone
	}

	if req.Email != nil {
		doctor.Email = *req.Email
	}

	if req.DepartmentID != nil {
		doctor.DepartmentID = req.DepartmentID
	}

	if req.IsActive != nil {
		doctor.IsActive = *req.IsActive
	}

	if err := s.repo.Update(doctor); err != nil {
		return nil, err
	}

	return s.toDoctorResponse(doctor), nil
}

func (s doctorService) ActivateDoctor(id uint) (*dto.DoctorResponse, error) {
	doctor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	doctor.IsActive = true

	if err := s.repo.Update(doctor); err != nil {
		return nil, err
	}

	return s.toDoctorResponse(doctor), nil
}

func (s doctorService) DeactivateDoctor(id uint) (*dto.DoctorResponse, error) {
	doctor, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	doctor.IsActive = false

	if err := s.repo.Update(doctor); err != nil {
		return nil, err
	}

	return s.toDoctorResponse(doctor), nil
}

func (s doctorService) SoftDeleteDoctor(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(id)
}

func (s doctorService) RestoreDoctor(id uint) error {
	return s.repo.Restore(id)
}

func (s doctorService) HardDeleteDoctor(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *doctorService) toDoctorResponse(doctor *models.Doctor) *dto.DoctorResponse {
	return &dto.DoctorResponse{
		ID:             doctor.ID,
		UserID:         doctor.UserID,
		EmployeeID:     doctor.EmployeeID,
		FullName:       doctor.FullName,
		Specialization: doctor.Specialization,
		LicenseNumber:  doctor.LicenseNumber,
		Phone:          doctor.Phone,
		Email:          doctor.Email,
		DepartmentID:   doctor.DepartmentID,
		IsActive:       &doctor.IsActive,
		CreatedAt:      doctor.CreatedAt,
		UpdatedAt:      doctor.UpdatedAt,
	}
}

func (s *doctorService) toDeleteDoctorResponse(doctor *models.Doctor) *dto.DeletedDoctorResponse {
	return &dto.DeletedDoctorResponse{
		ID:             doctor.ID,
		UserID:         doctor.UserID,
		EmployeeID:     doctor.EmployeeID,
		FullName:       doctor.FullName,
		Specialization: doctor.Specialization,
		LicenseNumber:  doctor.LicenseNumber,
		Phone:          doctor.Phone,
		Email:          doctor.Email,
		DepartmentID:   doctor.DepartmentID,
		IsActive:       &doctor.IsActive,
		CreatedAt:      doctor.CreatedAt,
		UpdatedAt:      doctor.UpdatedAt,
		DeletedAt:      &doctor.DeletedAt.Time,
	}
}
