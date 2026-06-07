package medicinetype

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type MedicineTypeService interface {
	List(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error)
	DeletedList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeDeletedListResponse, error)
	FindByID(id uint) (*dto.MedicineTypeResponse, error)
	Create(req *dto.CreateMedicineTypeRequest) (*dto.MedicineTypeResponse, error)
	Update(id uint, req *dto.UpdateMedicineTypeRequest) (*dto.MedicineTypeResponse, error)
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	ActiveList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error)
	InactiveList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error)
	Activate(id uint) error
	Deactivate(id uint) error
}

type medicineTypeService struct {
	repo   repository.MedicineTypeRepository
	config *config.Config
}

func NewMedicineTypeService(repo repository.MedicineTypeRepository, config *config.Config) MedicineTypeService {
	return &medicineTypeService{
		repo:   repo,
		config: config,
	}
}

func (s *medicineTypeService) List(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
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

	medicineTypes, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	medicineTypesResponse := make([]dto.MedicineTypeResponse, len(medicineTypes))
	for i, mt := range medicineTypes {
		medicineTypesResponse[i] = *s.toMedicineTypeResponse(&mt)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.MedicineTypeListResponse{
		Data: medicineTypesResponse,
		Meta: dto.MedicineTypePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineTypeService) DeletedList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeDeletedListResponse, error) {
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

	deletedMedicineTypes, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	deletedMedicineTypesResponse := make([]dto.DeletedMedicineTypeResponse, len(deletedMedicineTypes))
	for i, mt := range deletedMedicineTypes {
		deletedMedicineTypesResponse[i] = *s.toMedicineTypeDeletedResponse(&mt)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineTypeDeletedListResponse{
		Data: deletedMedicineTypesResponse,
		Meta: dto.MedicineTypePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}
func (s *medicineTypeService) FindByID(id uint) (*dto.MedicineTypeResponse, error) {
	medicineType, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toMedicineTypeResponse(medicineType), nil
}
func (s *medicineTypeService) Create(req *dto.CreateMedicineTypeRequest) (*dto.MedicineTypeResponse, error) {
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

	res := &models.MedicineType{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    *req.IsActive,
	}

	if err := s.repo.Create(res); err != nil {
		return nil, err
	}
	return s.toMedicineTypeResponse(res), nil
}
func (s *medicineTypeService) Update(id uint, req *dto.UpdateMedicineTypeRequest) (*dto.MedicineTypeResponse, error) {
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

	return s.toMedicineTypeResponse(res), nil
}
func (s *medicineTypeService) SoftDelete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}
func (s *medicineTypeService) Restore(id uint) error {
	return s.repo.Restore(id)
}
func (s *medicineTypeService) HardDelete(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *medicineTypeService) ActiveList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
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

	medicineTypes, total, err := s.repo.ActiveList(query)
	if err != nil {
		return nil, err
	}

	medicineTypesResponse := make([]dto.MedicineTypeResponse, len(medicineTypes))
	for i, mt := range medicineTypes {
		medicineTypesResponse[i] = *s.toMedicineTypeResponse(&mt)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.MedicineTypeListResponse{
		Data: medicineTypesResponse,
		Meta: dto.MedicineTypePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineTypeService) InactiveList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
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

	medicineTypes, total, err := s.repo.InactiveList(query)
	if err != nil {
		return nil, err
	}

	medicineTypesResponse := make([]dto.MedicineTypeResponse, len(medicineTypes))
	for i, mt := range medicineTypes {
		medicineTypesResponse[i] = *s.toMedicineTypeResponse(&mt)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.MedicineTypeListResponse{
		Data: medicineTypesResponse,
		Meta: dto.MedicineTypePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineTypeService) Activate(id uint) error {
	return s.repo.Activate(id)
}

func (s *medicineTypeService) Deactivate(id uint) error {
	return s.repo.Deactivate(id)
}

func (s *medicineTypeService) toMedicineTypeResponse(mt *models.MedicineType) *dto.MedicineTypeResponse {
	return &dto.MedicineTypeResponse{
		ID:          mt.ID,
		Name:        mt.Name,
		Code:        mt.Code,
		Description: mt.Description,
		IsActive:    mt.IsActive,
		CreatedAt:   mt.CreatedAt,
		UpdatedAt:   mt.UpdatedAt,
	}
}

func (s *medicineTypeService) toMedicineTypeDeletedResponse(mt *models.MedicineType) *dto.DeletedMedicineTypeResponse {
	return &dto.DeletedMedicineTypeResponse{
		ID:          mt.ID,
		Name:        mt.Name,
		Code:        mt.Code,
		Description: mt.Description,
		IsActive:    mt.IsActive,
		CreatedAt:   mt.CreatedAt,
		UpdatedAt:   mt.UpdatedAt,
		DeletedAt:   &mt.DeletedAt.Time,
	}
}
