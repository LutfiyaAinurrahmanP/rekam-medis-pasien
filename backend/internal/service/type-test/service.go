package typetest

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type TypeTestService interface {
	List(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error)
	DeletedList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestDeletedListResponse, error)
	FindByID(id uint) (*dto.TypeTestResponse, error)
	Create(req *dto.CreateTypeTestRequest) (*dto.TypeTestResponse, error)
	Update(id uint, req *dto.UpdateTypeTestRequest) (*dto.TypeTestResponse, error)
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	ActiveList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error)
	InactiveList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error)
	Activate(id uint) error
	Deactivate(id uint) error
}

type typeTestService struct {
	repo   repository.TypeTestRepository
	config *config.Config
}

func NewTypeTestService(repo repository.TypeTestRepository, config *config.Config) TypeTestService {
	return &typeTestService{
		repo:   repo,
		config: config,
	}
}

func (s *typeTestService) normalizeQuery(query *dto.TypeTestPaginationQuery, defaultSortBy, defaultSortDir string) {
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
		query.SortBy = defaultSortBy
	}

	if query.SortDir == "" {
		query.SortDir = defaultSortDir
	}
}

func (s *typeTestService) List(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	s.normalizeQuery(query, "name", "asc")

	typeTests, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	typeTestsResponse := make([]dto.TypeTestResponse, len(typeTests))
	for i, tt := range typeTests {
		typeTestsResponse[i] = *s.toResponse(&tt)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.TypeTestListResponse{
		Data: typeTestsResponse,
		Meta: dto.TypeTestPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *typeTestService) DeletedList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestDeletedListResponse, error) {
	s.normalizeQuery(query, "name", "asc")

	deletedTypeTests, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	deletedTypeTestsResponse := make([]dto.DeletedTypeTestResponse, len(deletedTypeTests))
	for i, dtt := range deletedTypeTests {
		deletedTypeTestsResponse[i] = *s.toDeletedResponse(&dtt)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.TypeTestDeletedListResponse{
		Data: deletedTypeTestsResponse,
		Meta: dto.TypeTestPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *typeTestService) FindByID(id uint) (*dto.TypeTestResponse, error) {
	typeTest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(typeTest), nil
}

func (s *typeTestService) Create(req *dto.CreateTypeTestRequest) (*dto.TypeTestResponse, error) {
	exists, err := s.repo.IsCodeExists(req.Code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("code already exists")
	}

	res := &models.TypeTest{
		Name:               req.Name,
		Code:               req.Code,
		TypeTestCategoryID: req.TypeTestCategoryID,
		Description:        req.Description,
	}

	if req.Price != nil {
		res.Price = *req.Price
	}

	if req.IsActive != nil {
		res.IsActive = *req.IsActive
	} else {
		res.IsActive = true // default
	}

	if err := s.repo.Create(res); err != nil {
		return nil, err
	}
	return s.toResponse(res), nil
}

func (s *typeTestService) Update(id uint, req *dto.UpdateTypeTestRequest) (*dto.TypeTestResponse, error) {
	res, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil && *req.Name != res.Name {
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

	if req.TypeTestCategoryID != nil {
		res.TypeTestCategoryID = *req.TypeTestCategoryID
	}

	if req.Description != nil {
		res.Description = *req.Description
	}

	if req.Price != nil {
		res.Price = *req.Price
	}

	if req.IsActive != nil {
		res.IsActive = *req.IsActive
	}

	if err := s.repo.Update(res); err != nil {
		return nil, err
	}

	return s.toResponse(res), nil
}

func (s *typeTestService) SoftDelete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}

func (s *typeTestService) Restore(id uint) error {
	return s.repo.Restore(id)
}

func (s *typeTestService) HardDelete(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *typeTestService) ActiveList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	s.normalizeQuery(query, "name", "asc")

	typeTests, total, err := s.repo.ActiveList(query)
	if err != nil {
		return nil, err
	}

	typeTestsResponse := make([]dto.TypeTestResponse, len(typeTests))
	for i, tt := range typeTests {
		typeTestsResponse[i] = *s.toResponse(&tt)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.TypeTestListResponse{
		Data: typeTestsResponse,
		Meta: dto.TypeTestPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *typeTestService) InactiveList(query *dto.TypeTestPaginationQuery) (*dto.TypeTestListResponse, error) {
	s.normalizeQuery(query, "name", "asc")

	typeTests, total, err := s.repo.InactiveList(query)
	if err != nil {
		return nil, err
	}

	typeTestsResponse := make([]dto.TypeTestResponse, len(typeTests))
	for i, tt := range typeTests {
		typeTestsResponse[i] = *s.toResponse(&tt)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.TypeTestListResponse{
		Data: typeTestsResponse,
		Meta: dto.TypeTestPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *typeTestService) Activate(id uint) error {
	return s.repo.Activate(id)
}

func (s *typeTestService) Deactivate(id uint) error {
	return s.repo.Deactivate(id)
}

func (s *typeTestService) toResponse(tt *models.TypeTest) *dto.TypeTestResponse {
	return &dto.TypeTestResponse{
		ID:                 tt.ID,
		Name:               tt.Name,
		Code:               tt.Code,
		TypeTestCategoryID: tt.TypeTestCategoryID,
		Description:        tt.Description,
		Price:              tt.Price,
		IsActive:           tt.IsActive,
		CreatedAt:          tt.CreatedAt,
		UpdatedAt:          tt.UpdatedAt,
	}
}

func (s *typeTestService) toDeletedResponse(tt *models.TypeTest) *dto.DeletedTypeTestResponse {
	return &dto.DeletedTypeTestResponse{
		ID:                 tt.ID,
		Name:               tt.Name,
		Code:               tt.Code,
		TypeTestCategoryID: tt.TypeTestCategoryID,
		Description:        tt.Description,
		Price:              tt.Price,
		IsActive:           tt.IsActive,
		CreatedAt:          tt.CreatedAt,
		UpdatedAt:          tt.UpdatedAt,
		DeletedAt:          &tt.DeletedAt.Time,
	}
}
