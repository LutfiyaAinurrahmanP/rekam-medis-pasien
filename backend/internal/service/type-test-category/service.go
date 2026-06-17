package typetestcategory

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type TypeTestCategoryService interface {
	List(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error)
	DeletedList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryDeletedListResponse, error)
	FindByID(id uint) (*dto.TypeTestCategoryResponse, error)
	Create(req *dto.CreateTypeTestCategoryRequest) (*dto.TypeTestCategoryResponse, error)
	Update(id uint, req *dto.UpdateTypeTestCategoryRequest) (*dto.TypeTestCategoryResponse, error)
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	ActiveList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error)
	InactiveList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error)
	Activate(id uint) error
	Deactivate(id uint) error
}

type typeTestCategoryService struct {
	repo   repository.TypeTestCategoryRepository
	config *config.Config
}

func NewTypeTestCategoryService(repo repository.TypeTestCategoryRepository, config *config.Config) TypeTestCategoryService {
	return &typeTestCategoryService{
		repo:   repo,
		config: config,
	}
}

func (s *typeTestCategoryService) List(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
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

	typeTestCategories, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	typeTestCategoriesResponse := make([]dto.TypeTestCategoryResponse, len(typeTestCategories))
	for i, ttc := range typeTestCategories {
		typeTestCategoriesResponse[i] = *s.toTypeTestCategoryResponse(&ttc)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.TypeTestCategoryListResponse{
		Data: typeTestCategoriesResponse,
		Meta: dto.TypeTestCategoryPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *typeTestCategoryService) DeletedList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryDeletedListResponse, error) {
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

	deletedTypeTestCategories, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	deletedTypeTestCategoriesResponse := make([]dto.DeletedTypeTestCategoryResponse, len(deletedTypeTestCategories))
	for i, dttc := range deletedTypeTestCategories {
		deletedTypeTestCategoriesResponse[i] = *s.toTypeTestCategoryDeletedResponse(&dttc)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.TypeTestCategoryDeletedListResponse{
		Data: deletedTypeTestCategoriesResponse,
		Meta: dto.TypeTestCategoryPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}
func (s *typeTestCategoryService) FindByID(id uint) (*dto.TypeTestCategoryResponse, error) {
	typeTestCategory, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toTypeTestCategoryResponse(typeTestCategory), nil
}
func (s *typeTestCategoryService) Create(req *dto.CreateTypeTestCategoryRequest) (*dto.TypeTestCategoryResponse, error) {
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

	res := &models.TypeTestCategory{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    *req.IsActive,
	}

	if err := s.repo.Create(res); err != nil {
		return nil, err
	}
	return s.toTypeTestCategoryResponse(res), nil
}
func (s *typeTestCategoryService) Update(id uint, req *dto.UpdateTypeTestCategoryRequest) (*dto.TypeTestCategoryResponse, error) {
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

	return s.toTypeTestCategoryResponse(res), nil
}
func (s *typeTestCategoryService) SoftDelete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}
func (s *typeTestCategoryService) Restore(id uint) error {
	return s.repo.Restore(id)
}
func (s *typeTestCategoryService) HardDelete(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *typeTestCategoryService) ActiveList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
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

	typeTestCategories, total, err := s.repo.ActiveList(query)
	if err != nil {
		return nil, err
	}

	typeTestCategoriesResponse := make([]dto.TypeTestCategoryResponse, len(typeTestCategories))
	for i, ttc := range typeTestCategories {
		typeTestCategoriesResponse[i] = *s.toTypeTestCategoryResponse(&ttc)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.TypeTestCategoryListResponse{
		Data: typeTestCategoriesResponse,
		Meta: dto.TypeTestCategoryPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *typeTestCategoryService) InactiveList(query *dto.TypeTestCategoryPaginationQuery) (*dto.TypeTestCategoryListResponse, error) {
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

	typeTestCategories, total, err := s.repo.InactiveList(query)
	if err != nil {
		return nil, err
	}

	typeTestCategoriesResponse := make([]dto.TypeTestCategoryResponse, len(typeTestCategories))
	for i, ttc := range typeTestCategories {
		typeTestCategoriesResponse[i] = *s.toTypeTestCategoryResponse(&ttc)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.TypeTestCategoryListResponse{
		Data: typeTestCategoriesResponse,
		Meta: dto.TypeTestCategoryPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *typeTestCategoryService) Activate(id uint) error {
	return s.repo.Activate(id)
}

func (s *typeTestCategoryService) Deactivate(id uint) error {
	return s.repo.Deactivate(id)
}

func (s *typeTestCategoryService) toTypeTestCategoryResponse(ttc *models.TypeTestCategory) *dto.TypeTestCategoryResponse {
	return &dto.TypeTestCategoryResponse{
		ID:          ttc.ID,
		Name:        ttc.Name,
		Code:        ttc.Code,
		Description: ttc.Description,
		IsActive:    ttc.IsActive,
		CreatedAt:   ttc.CreatedAt,
		UpdatedAt:   ttc.UpdatedAt,
	}
}

func (s *typeTestCategoryService) toTypeTestCategoryDeletedResponse(ttc *models.TypeTestCategory) *dto.DeletedTypeTestCategoryResponse {
	return &dto.DeletedTypeTestCategoryResponse{
		ID:          ttc.ID,
		Name:        ttc.Name,
		Code:        ttc.Code,
		Description: ttc.Description,
		IsActive:    ttc.IsActive,
		CreatedAt:   ttc.CreatedAt,
		UpdatedAt:   ttc.UpdatedAt,
		DeletedAt:   &ttc.DeletedAt.Time,
	}
}
