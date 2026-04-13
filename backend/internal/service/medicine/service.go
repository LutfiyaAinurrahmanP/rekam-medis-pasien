package medicine

import (
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type MedicineService interface {
	List(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error)
	DeletedList(query *dto.MedicinePaginationQuery) (*dto.MedicineDeletedListResponse, error)
	ListByAvailable(query *dto.MedicinePaginationQuery) (*dto.MedicineAvailableResponse, error)
	ListByLowStock(query *dto.MedicinePaginationQuery) (*dto.MedicineLowStockResponse, error)
	ListByOutStock(query *dto.MedicinePaginationQuery) (*dto.MedicineOutOfStockResponse, error)
	ListByInactive(query *dto.MedicinePaginationQuery) (*dto.MedicineInactiveResponse, error)
	FindByID(id uint) (*dto.MedicineResponse, error)
	FindByName(name string) (*dto.MedicineResponse, error)
	// GetByType()
	// Search()
	// Create()
	// Update()
	// UpdateStock()
	// UpdateMedicineReduceStock()
	// UpdateActivate()
	// UpdateDeactivate()
	// SoftDelete()
	// Restore()
	// HardDelete()
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
	s.normalizeQuery(query, "name", "asc")

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
			Page: query.Page,
			PageSize: query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) DeletedList(query *dto.MedicinePaginationQuery) (*dto.MedicineDeletedListResponse, error) {
	s.normalizeQuery(query, "name", "asc")

	medicines, total, err :=s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineDeletedResponse, len(medicines))
	for i, m := range medicines{
		responses[i] = *s.toDeletedResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineDeletedListResponse{
		Data: responses,
		Meta: dto.MedicinePaginationMeta{
			Page: query.Page,
			PageSize: query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) ListByAvailable(query *dto.MedicinePaginationQuery) (*dto.MedicineAvailableResponse, error) {
	s.normalizeQuery(query, "name", "asc")

	medicines, total, err := s.repo.ListByAvailable(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineResponse, len(medicines))
	for i, m := range medicines {
		responses[i] = *s.toResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineAvailableResponse{
		TotalAvailable:     total,
		Data:               responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) ListByLowStock(query *dto.MedicinePaginationQuery) (*dto.MedicineLowStockResponse, error) {
	s.normalizeQuery(query, "name", "asc")

	medicines, total, err := s.repo.ListByLowStock(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineResponse, len(medicines))
	for i, m := range medicines{
		responses[i] = *s.toResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineLowStockResponse{
		Threshold:          10, 
		TotalLowStock:      total,
		Data:               responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
} 

func (s *medicineService) ListByOutStock(query *dto.MedicinePaginationQuery) (*dto.MedicineOutOfStockResponse, error) {
	s.normalizeQuery(query, "name", "asc")

	medicines, total, err := s.repo.ListByOutStock(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineResponse, len(medicines))
	for i, m := range medicines {
		responses[i] = *s.toResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineOutOfStockResponse{
		TotalOutOfStock: total,
		Data:               responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicineService) ListByInactive(query *dto.MedicinePaginationQuery) (*dto.MedicineInactiveResponse, error) {
	s.normalizeQuery(query, "name", "asc")

	medicines, total, err := s.repo.ListByInactive(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicineResponse, len(medicines))
	for i, m := range medicines{
		responses[i] = *s.toResponse(&m)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineInactiveResponse{
		TotalInactive: total,
		Data:               responses,
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

func (s *medicineService) FindByName(name string) (*dto.MedicineResponse, error) {
	m, err := s.repo.FindByName(name)
	if err != nil {
		return nil, err
	}
	return s.toResponse(m), nil
}


func (s *medicineService) toResponse(m *models.Medicine) *dto.MedicineResponse{
	return &dto.MedicineResponse{
		ID: m.ID,
		Name: m.Name,
		GenericName: m.GenericName,
		BrandName: m.BrandName,
		Type: m.Type,
		Strength: m.Strength,
		Manufacturer: m.Manufacturer,
		Unit: m.Unit,
		StockQuantity: m.StockQuantity,
		Price: m.Price,
		IsActive: m.IsActive,
		IsLowStock: m.StockQuantity == 10,
		IsOutOfStock: m.StockQuantity == 0,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
} 

func (s *medicineService) toDeletedResponse(m *models.Medicine) *dto.MedicineDeletedResponse{
	return &dto.MedicineDeletedResponse{
		ID: m.ID,
		Name: m.Name,
		GenericName: m.GenericName,
		BrandName: m.BrandName,
		Type: m.Type,
		Strength: m.Strength,
		Manufacturer: m.Manufacturer,
		Unit: m.Unit,
		StockQuantity: m.StockQuantity,
		Price: m.Price,
		IsActive: m.IsActive,
		IsLowStock: m.StockQuantity == 10,
		IsOutOfStock: m.StockQuantity == 0,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: &m.DeletedAt.Time,
	}
} 