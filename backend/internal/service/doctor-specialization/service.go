package doctorspecialization

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type DoctorSpecializationService interface {
	List(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error)
	DeletedList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationDeletedListResponse, error)
	FindByID(id uint) (*dto.DoctorSpecializationResponse, error)
	Create(req *dto.CreateDoctorSpecializationRequest) (*dto.DoctorSpecializationResponse, error)
	Update(id uint, req *dto.UpdateDoctorSpecializationRequest) (*dto.DoctorSpecializationResponse, error)
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	ActiveList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error)
	InactiveList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error)
	Activate(id uint) error
	Deactivate(id uint) error
}

type doctorSpecializationService struct {
	repo   repository.DoctorSpecializationRepository
	config *config.Config
}

func NewDoctorSpecializationService(repo repository.DoctorSpecializationRepository, config *config.Config) DoctorSpecializationService {
	return &doctorSpecializationService{
		repo:   repo,
		config: config,
	}
}

func (s *doctorSpecializationService) List(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error) {
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

	doctorSpecializations, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	doctorSpecializationsResponse := make([]dto.DoctorSpecializationResponse, len(doctorSpecializations))
	for i, ds := range doctorSpecializations {
		doctorSpecializationsResponse[i] = *s.toDoctorSpecializationResponse(&ds)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.DoctorSpecializationListResponse{
		Data: doctorSpecializationsResponse,
		Meta: dto.DoctorSpecializationPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *doctorSpecializationService) DeletedList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationDeletedListResponse, error) {
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

	deletedDoctorSpecialization, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	deletedDoctorSpecializationResponse := make([]dto.DeletedDoctorSpecializationResponse, len(deletedDoctorSpecialization))
	for i, dds := range deletedDoctorSpecialization {
		deletedDoctorSpecializationResponse[i] = *s.toDoctorSpecializationDeletedResponse(&dds)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.DoctorSpecializationDeletedListResponse{
		Data: deletedDoctorSpecializationResponse,
		Meta: dto.DoctorSpecializationPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}
func (s *doctorSpecializationService) FindByID(id uint) (*dto.DoctorSpecializationResponse, error) {
	doctorSpecialization, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toDoctorSpecializationResponse(doctorSpecialization), nil
}
func (s *doctorSpecializationService) Create(req *dto.CreateDoctorSpecializationRequest) (*dto.DoctorSpecializationResponse, error) {
	exists, err := s.repo.IsNameExists(req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("name already exists")
	}

	exists, err = s.repo.IsCodeExists(req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("code already exists")
	}

	res := &models.DoctorSpecialization{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    *req.IsActive,
	}

	if err := s.repo.Create(res); err != nil {
		return nil, err
	}
	return s.toDoctorSpecializationResponse(res), nil
}
func (s *doctorSpecializationService) Update(id uint, req *dto.UpdateDoctorSpecializationRequest) (*dto.DoctorSpecializationResponse, error) {
	res, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil && *req.Name != res.Name {
		exists, err := s.repo.IsNameExists(*req.Name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("name already exists")
		}
		res.Name = *req.Name
	}

	if req.Code != nil && *req.Code != res.Code {
		exists, err := s.repo.IsCodeExists(*req.Code, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("code already exists")
		}
		res.Code = *req.Code
	}

	if req.Description != nil {
		res.Description = *req.Description
	}

	if req.IsActive != nil {
		res.IsActive = *req.IsActive
	}

	if err := s.repo.Update(res); err != nil {
		return nil, err
	}

	return s.toDoctorSpecializationResponse(res), nil
}
func (s *doctorSpecializationService) SoftDelete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}
func (s *doctorSpecializationService) Restore(id uint) error {
	return s.repo.Restore(id)
}
func (s *doctorSpecializationService) HardDelete(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *doctorSpecializationService) ActiveList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error) {
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

	doctorSpecializations, total, err := s.repo.ActiveList(query)
	if err != nil {
		return nil, err
	}

	doctorSpecializationsResponse := make([]dto.DoctorSpecializationResponse, len(doctorSpecializations))
	for i, ds := range doctorSpecializations {
		doctorSpecializationsResponse[i] = *s.toDoctorSpecializationResponse(&ds)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.DoctorSpecializationListResponse{
		Data: doctorSpecializationsResponse,
		Meta: dto.DoctorSpecializationPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *doctorSpecializationService) InactiveList(query *dto.DoctorSpecializationPaginationQuery) (*dto.DoctorSpecializationListResponse, error) {
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

	doctorSpecializations, total, err := s.repo.InactiveList(query)
	if err != nil {
		return nil, err
	}

	doctorSpecializationsResponse := make([]dto.DoctorSpecializationResponse, len(doctorSpecializations))
	for i, ds := range doctorSpecializations {
		doctorSpecializationsResponse[i] = *s.toDoctorSpecializationResponse(&ds)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.DoctorSpecializationListResponse{
		Data: doctorSpecializationsResponse,
		Meta: dto.DoctorSpecializationPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *doctorSpecializationService) Activate(id uint) error {
	return s.repo.Activate(id)
}

func (s *doctorSpecializationService) Deactivate(id uint) error {
	return s.repo.Deactivate(id)
}

func (s *doctorSpecializationService) toDoctorSpecializationResponse(ds *models.DoctorSpecialization) *dto.DoctorSpecializationResponse {
	return &dto.DoctorSpecializationResponse{
		ID:          ds.ID,
		Name:        ds.Name,
		Code:        ds.Code,
		Description: ds.Description,
		IsActive:    ds.IsActive,
		CreatedAt:   ds.CreatedAt,
		UpdatedAt:   ds.UpdatedAt,
	}
}

func (s *doctorSpecializationService) toDoctorSpecializationDeletedResponse(ds *models.DoctorSpecialization) *dto.DeletedDoctorSpecializationResponse {
	return &dto.DeletedDoctorSpecializationResponse{
		ID:          ds.ID,
		Name:        ds.Name,
		Code:        ds.Code,
		Description: ds.Description,
		IsActive:    ds.IsActive,
		CreatedAt:   ds.CreatedAt,
		UpdatedAt:   ds.UpdatedAt,
		DeletedAt:   &ds.DeletedAt.Time,
	}
}
