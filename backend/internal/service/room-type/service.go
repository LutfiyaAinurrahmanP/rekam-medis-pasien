package doctorspecialization

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type RoomTypeService interface {
	List(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error)
	DeletedList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeDeletedListResponse, error)
	FindByID(id uint) (*dto.RoomTypeResponse, error)
	Create(req *dto.CreateRoomTypeRequest) (*dto.RoomTypeResponse, error)
	Update(id uint, req *dto.UpdateRoomTypeRequest) (*dto.RoomTypeResponse, error)
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	ActiveList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error)
	InactiveList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error)
	Activate(id uint) error
	Deactivate(id uint) error
}

type roomTypeService struct {
	repo   repository.RoomTypeRepository
	config *config.Config
}

func NewRoomTypeService(repo repository.RoomTypeRepository, config *config.Config) RoomTypeService {
	return &roomTypeService{
		repo:   repo,
		config: config,
	}
}

func (s *roomTypeService) List(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
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

	roomTypes, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	roomTypesResponse := make([]dto.RoomTypeResponse, len(roomTypes))
	for i, ds := range roomTypes {
		roomTypesResponse[i] = *s.toRoomTypeResponse(&ds)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.RoomTypeListResponse{
		Data: roomTypesResponse,
		Meta: dto.RoomTypePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *roomTypeService) DeletedList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeDeletedListResponse, error) {
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

	deletedRoomType, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	deletedRoomTypeResponse := make([]dto.DeletedRoomTypeResponse, len(deletedRoomType))
	for i, dds := range deletedRoomType {
		deletedRoomTypeResponse[i] = *s.toRoomTypeDeletedResponse(&dds)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.RoomTypeDeletedListResponse{
		Data: deletedRoomTypeResponse,
		Meta: dto.RoomTypePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}
func (s *roomTypeService) FindByID(id uint) (*dto.RoomTypeResponse, error) {
	roomType, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toRoomTypeResponse(roomType), nil
}
func (s *roomTypeService) Create(req *dto.CreateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
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

	res := &models.RoomType{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IsActive:    *req.IsActive,
	}

	if err := s.repo.Create(res); err != nil {
		return nil, err
	}
	return s.toRoomTypeResponse(res), nil
}
func (s *roomTypeService) Update(id uint, req *dto.UpdateRoomTypeRequest) (*dto.RoomTypeResponse, error) {
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

	return s.toRoomTypeResponse(res), nil
}
func (s *roomTypeService) SoftDelete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}
func (s *roomTypeService) Restore(id uint) error {
	return s.repo.Restore(id)
}
func (s *roomTypeService) HardDelete(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *roomTypeService) ActiveList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
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

	roomTypes, total, err := s.repo.ActiveList(query)
	if err != nil {
		return nil, err
	}

	roomTypesResponse := make([]dto.RoomTypeResponse, len(roomTypes))
	for i, ds := range roomTypes {
		roomTypesResponse[i] = *s.toRoomTypeResponse(&ds)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.RoomTypeListResponse{
		Data: roomTypesResponse,
		Meta: dto.RoomTypePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *roomTypeService) InactiveList(query *dto.RoomTypePaginationQuery) (*dto.RoomTypeListResponse, error) {
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

	roomTypes, total, err := s.repo.InactiveList(query)
	if err != nil {
		return nil, err
	}

	roomTypesResponse := make([]dto.RoomTypeResponse, len(roomTypes))
	for i, ds := range roomTypes {
		roomTypesResponse[i] = *s.toRoomTypeResponse(&ds)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.RoomTypeListResponse{
		Data: roomTypesResponse,
		Meta: dto.RoomTypePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *roomTypeService) Activate(id uint) error {
	return s.repo.Activate(id)
}

func (s *roomTypeService) Deactivate(id uint) error {
	return s.repo.Deactivate(id)
}

func (s *roomTypeService) toRoomTypeResponse(ds *models.RoomType) *dto.RoomTypeResponse {
	return &dto.RoomTypeResponse{
		ID:          ds.ID,
		Name:        ds.Name,
		Code:        ds.Code,
		Description: ds.Description,
		IsActive:    ds.IsActive,
		CreatedAt:   ds.CreatedAt,
		UpdatedAt:   ds.UpdatedAt,
	}
}

func (s *roomTypeService) toRoomTypeDeletedResponse(ds *models.RoomType) *dto.DeletedRoomTypeResponse {
	return &dto.DeletedRoomTypeResponse{
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
