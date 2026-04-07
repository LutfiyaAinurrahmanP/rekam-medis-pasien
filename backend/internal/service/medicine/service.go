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
	// DeletedList()
	GetByAvailable(query *dto.MedicinePaginationQuery) (*dto.MedicineAvailableResponse, error)
	// GetByLowStock()
	// GetByOutStock()
	// GetByInactive()
	// GetByID()
	// GetByName()
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

func (s *medicineService) GetByAvailable(query *dto.MedicinePaginationQuery) (*dto.MedicineAvailableResponse, error) {
	s.normalizeQuery(query, "name", "asc")

	medicines, total, err := s.repo.ListByAvailable(query)
	if err != nil {
		return nil, err
	}

	// Convert medicines to response format
	responses := make([]dto.MedicineResponse, len(medicines))
	medicineTypeMap := make(map[string]int64)
	for i, m := range medicines {
		responses[i] = *s.toResponse(&m)
		medicineTypeMap[m.Type]++
	}

	// Get total stock value from ALL available medicines (not just current page)
	totalStockValue, err := s.repo.GetTotalStockValueByAvailable(query)
	if err != nil {
		return nil, err
	}

	// Get total stock quantity from ALL available medicines (not just current page)
	totalStockQuantity, err := s.repo.GetTotalStockQuantityByAvailable(query)
	if err != nil {
		return nil, err
	}

	// Convert map to slice for response
	medicineTypes := make([]dto.MedicineTypeCount, 0, len(medicineTypeMap))
	for t, count := range medicineTypeMap {
		medicineTypes = append(medicineTypes, dto.MedicineTypeCount{
			Type:  t,
			Count: count,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicineAvailableResponse{
		TotalAvailable:     total,
		TotalStockQuantity: totalStockQuantity,
		TotalStockValue:    totalStockValue,
		MedicineTypes:      medicineTypes,
		Data:               responses,
		Meta: dto.MedicinePaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
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