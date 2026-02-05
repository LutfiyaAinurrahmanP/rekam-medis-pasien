package service

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type DepartmentService interface {
	ListDepartments(query *dto.DepartmentPaginationQuery) (*dto.DepartmentListResponse, error)
	DeleteListDepartments(query *dto.DepartmentPaginationQuery) (*dto.DepartmentDeletedListResponse, error)
	GetDepartmentByID(id uint) (*dto.DepartmentResponse, error)
	CreateDepartment(req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error)
	UpdateDepartment(id uint, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error)
	SoftDeleteDepartment(id uint) error
	RestoreDepartment(id uint) error
	HardDeleteDepartment(id uint) error
}

type departmentService struct {
	repo   repository.DepartmentRepository
	config *config.Config
}

func NewDepartmentService(repo repository.DepartmentRepository, config *config.Config) DepartmentService {
	return &departmentService{
		repo:   repo,
		config: config,
	}
}

func (s departmentService) ListDepartments(query *dto.DepartmentPaginationQuery) (*dto.DepartmentListResponse, error) {
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

	departments, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	departmentResponses := make([]dto.DepartmentResponse, len(departments))
	for i, department := range departments {
		departmentResponses[i] = *s.toDepartmentResponse(&department)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.DepartmentListResponse{
		Data: departmentResponses,
		Meta: dto.DepartmentPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s departmentService) DeleteListDepartments(query *dto.DepartmentPaginationQuery) (*dto.DepartmentDeletedListResponse, error) {
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

	deletedDepartments, total, err := s.repo.DeleteList(query)
	if err != nil {
		return nil, err
	}

	deletedDepartmentsResponses := make([]dto.DeletedDepartmentResponse, len(deletedDepartments))
	for i, deletedDepartment := range deletedDepartments {
		deletedDepartmentsResponses[i] = *s.toDeleteDepartmentResponse(&deletedDepartment)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.DepartmentDeletedListResponse{
		Data: deletedDepartmentsResponses,
		Meta: dto.DepartmentPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s departmentService) GetDepartmentByID(id uint) (*dto.DepartmentResponse, error) {
	department, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return s.toDepartmentResponse(department), nil
}

func (s departmentService) CreateDepartment(req *dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	exists, err := s.repo.IsCodeExists(req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("code already exists")
	}

	department := &models.Department{
		Name:          req.Name,
		Code:          req.Code,
		Description:   req.Description,
		FloorLocation: req.FloorLocation,
	}

	if err := s.repo.Create(department); err != nil {
		return nil, err
	}
	return s.toDepartmentResponse(department), nil
}

func (s departmentService) UpdateDepartment(id uint, req *dto.UpdateDepartmentRequest) (*dto.DepartmentResponse, error) {
	department, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	// Update Name if provided
	if req.Name != nil {
		department.Name = *req.Name
	}

	// Update Code if provided and different
	if req.Code != nil && *req.Code != department.Code {
		exists, err := s.repo.IsCodeExists(*req.Code, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("code already exists")
		}
		department.Code = *req.Code
	}

	// Update Description if provided
	if req.Description != nil {
		department.Description = *req.Description
	}

	// Update FloorLocation if provided
	if req.FloorLocation != nil {
		department.FloorLocation = *req.FloorLocation
	}

	if err := s.repo.Update(department); err != nil {
		return nil, err
	}

	return s.toDepartmentResponse(department), nil
}

func (s departmentService) SoftDeleteDepartment(id uint) error {
	_, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(id)
}

func (s departmentService) RestoreDepartment(id uint) error {
	return s.repo.Restore(id)
}

func (s departmentService) HardDeleteDepartment(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *departmentService) toDepartmentResponse(department *models.Department) *dto.DepartmentResponse {
	return &dto.DepartmentResponse{
		ID:            department.ID,
		Name:          department.Name,
		Code:          department.Code,
		Description:   department.Description,
		FloorLocation: department.FloorLocation,
		CreatedAt:     department.CreatedAt,
		UpdatedAt:     department.UpdatedAt,
	}
}

func (s *departmentService) toDeleteDepartmentResponse(department *models.Department) *dto.DeletedDepartmentResponse {
	return &dto.DeletedDepartmentResponse{
		ID:            department.ID,
		Name:          department.Name,
		Code:          department.Code,
		Description:   department.Description,
		FloorLocation: department.FloorLocation,
		CreatedAt:     department.CreatedAt,
		UpdatedAt:     department.UpdatedAt,
		DeletedAt:     &department.DeletedAt.Time,
	}
}
