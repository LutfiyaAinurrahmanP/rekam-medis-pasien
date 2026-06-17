package medicine

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type MedicineService interface {
	List(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error)
	DeletedList(query *dto.MedicinePaginationQuery) (*dto.MedicineDeletedListResponse, error)
	AvailableList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error)
	LowStockList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error)
	OutStockList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error)
	ActiveList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error)
	InactiveList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error)
	FindByID(id uint) (*dto.MedicineResponse, error)
	FindByIDUnscoped(id uint) (*dto.DeletedMedicineResponse, error)
	Create(req *dto.CreateMedicineRequest) (*dto.MedicineResponse, error)
	Update(id uint, req *dto.UpdateMedicineRequest) (*dto.MedicineResponse, error)
	AddStock(id uint, req *dto.AddStockRequest) error
	ReduceStock(id uint, req *dto.ReduceStockRequest) error
	Activate(id uint) error
	Deactivate(id uint, req *dto.DeactivateMedicineRequest) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type medicineService struct {
	repo   repository.MedicineRepository
	config *config.Config
}

func NewMedicineService(repo repository.MedicineRepository, config *config.Config) MedicineService {
	return &medicineService{
		repo:   repo,
		config: config,
	}
}

func (s *medicineService) normalizeQuery(query *dto.MedicinePaginationQuery, defaultSortBy, defaultSortDir string) {
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

func (s *medicineService) List(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")

	medicines, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineResponse, len(medicines))
	for i, m := range medicines {
		responses[i] = *s.toResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineListResponse{
		Data: responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) DeletedList(query *dto.MedicinePaginationQuery) (*dto.MedicineDeletedListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")

	medicines, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DeletedMedicineResponse, len(medicines))
	for i, m := range medicines {
		responses[i] = *s.toDeletedResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineDeletedListResponse{
		Data: responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) AvailableList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")
	medicines, total, err := s.repo.AvailableList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineResponse, len(medicines))
	for i, m := range medicines {
		responses[i] = *s.toResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineListResponse{
		Data: responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) LowStockList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	s.normalizeQuery(query, "stock", "asc")
	medicines, total, err := s.repo.LowStockList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineResponse, len(medicines))
	for i, m := range medicines {
		responses[i] = *s.toResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineListResponse{
		Data: responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) OutStockList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	s.normalizeQuery(query, "stock", "asc")
	medicines, total, err := s.repo.OutStockList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineResponse, len(medicines))
	for i, m := range medicines {
		responses[i] = *s.toResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineListResponse{
		Data: responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) ActiveList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")
	medicines, total, err := s.repo.ActiveList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineResponse, len(medicines))
	for i, m := range medicines {
		responses[i] = *s.toResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineListResponse{
		Data: responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) InactiveList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")
	medicines, total, err := s.repo.InactiveList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineResponse, len(medicines))
	for i, m := range medicines {
		responses[i] = *s.toResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineListResponse{
		Data: responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) FindByID(id uint) (*dto.MedicineResponse, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(m), nil
}

func (s *medicineService) FindByIDUnscoped(id uint) (*dto.DeletedMedicineResponse, error) {
	m, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return nil, err
	}
	return s.toDeletedResponse(m), nil
}

func (s *medicineService) Create(req *dto.CreateMedicineRequest) (*dto.MedicineResponse, error) {
	m := &models.Medicine{
		Name:           req.Name,
		GenericName:    req.GenericName,
		BrandName:      req.BrandName,
		MedicineTypeID: *req.MedicineTypeID,
		Strength:       req.Strength,
		Manufacturer:   req.Manufacturer,
		Unit:           req.Unit,
		StockQuantity:  *req.StockQuantity,
	}

	if req.Price != nil {
		m.Price = *req.Price
	}
	if req.IsActive != nil {
		m.IsActive = *req.IsActive
	} else {
		m.IsActive = true
	}

	if err := s.repo.Create(m); err != nil {
		return nil, err
	}

	return s.toResponse(m), nil
}

func (s *medicineService) Update(id uint, req *dto.UpdateMedicineRequest) (*dto.MedicineResponse, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.GenericName != nil {
		m.GenericName = *req.GenericName
	}
	if req.BrandName != nil {
		m.BrandName = *req.BrandName
	}
	if req.MedicineTypeID != nil {
		m.MedicineTypeID = *req.MedicineTypeID
	}
	if req.Strength != nil {
		m.Strength = *req.Strength
	}
	if req.Manufacturer != nil {
		m.Manufacturer = *req.Manufacturer
	}
	if req.Unit != nil {
		m.Unit = *req.Unit
	}
	if req.Price != nil {
		m.Price = *req.Price
	}

	if err := s.repo.Update(m); err != nil {
		return nil, err
	}

	return s.toResponse(m), nil
}

func (s *medicineService) AddStock(id uint, req *dto.AddStockRequest) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !m.IsActive {
		return errors.New("medicine is not active")
	}
	return s.repo.AddStock(id, req.Quantity)
}

func (s *medicineService) ReduceStock(id uint, req *dto.ReduceStockRequest) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if !m.IsActive {
		return errors.New("medicine is not active")
	}

	if m.StockQuantity < req.Quantity {
		return errors.New("insufficient stock")
	}

	return s.repo.ReduceStock(id, req.Quantity)
}

func (s *medicineService) Activate(id uint) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if m.IsActive {
		return errors.New("medicine is already active")
	}
	return s.repo.Activate(id)
}

func (s *medicineService) Deactivate(id uint, req *dto.DeactivateMedicineRequest) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !m.IsActive {
		return errors.New("medicine is already inactive")
	}
	return s.repo.Deactivate(id)
}

func (s *medicineService) SoftDelete(id uint) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if m.DeletedAt.Valid {
		return errors.New("medicine is already deleted")
	}
	return s.repo.SoftDelete(id)
}

func (s *medicineService) Restore(id uint) error {
	m, _ := s.repo.FindByIDUnscoped(id)
	if m == nil {
		return errors.New("medicine not found")
	}
	if !m.DeletedAt.Valid {
		return errors.New("medicine is not deleted")
	}
	return s.repo.Restore(id)
}

func (s *medicineService) HardDelete(id uint) error {
	m, _ := s.repo.FindByIDUnscoped(id)
	if m == nil {
		return errors.New("medicine not found")
	}
	if !m.DeletedAt.Valid {
		return errors.New("medicine is not deleted")
	}
	return s.repo.HardDelete(id)
}

func (s *medicineService) toResponse(m *models.Medicine) *dto.MedicineResponse {
	return &dto.MedicineResponse{
		ID:             m.ID,
		Name:           m.Name,
		GenericName:    m.GenericName,
		BrandName:      m.BrandName,
		MedicineTypeID: m.MedicineTypeID,
		Strength:       m.Strength,
		Manufacturer:   m.Manufacturer,
		Unit:           m.Unit,
		StockQuantity:  m.StockQuantity,
		Price:          m.Price,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func (s *medicineService) toDeletedResponse(m *models.Medicine) *dto.DeletedMedicineResponse {
	return &dto.DeletedMedicineResponse{
		ID:             m.ID,
		Name:           m.Name,
		GenericName:    m.GenericName,
		BrandName:      m.BrandName,
		MedicineTypeID: m.MedicineTypeID,
		Strength:       m.Strength,
		Manufacturer:   m.Manufacturer,
		Unit:           m.Unit,
		StockQuantity:  m.StockQuantity,
		Price:          m.Price,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      &m.DeletedAt.Time,
	}
}
