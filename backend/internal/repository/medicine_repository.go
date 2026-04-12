package repository

import (
	"errors"
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type MedicineRepository interface {
	List(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	DeletedList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	ListByAvailable(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	GetTotalStockValueByAvailable(query *dto.MedicinePaginationQuery) (float64, error)
	GetTotalStockQuantityByAvailable(query *dto.MedicinePaginationQuery) (int64, error)
	ListByLowStock(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	GetTotalStockValueByLowStock(query *dto.MedicinePaginationQuery) (float64, error)
	GetTotalStockQuantityByLowStock(query *dto.MedicinePaginationQuery) (int64, error)
	ListByOutStock(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	ListByInactive(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	FindByID(id uint) (*models.Medicine, error)
	FindByName(name string) (*models.Medicine, error)
	FindByType(types string) (*models.Medicine, error)
	Search(criteria *dto.MedicineSearchCriteria, item *dto.MedicineSearchItem, response *dto.MedicineSearchResponse) ([]models.Medicine, int64, error)
	Create(medicine *models.Medicine) error
	Update(medicine *models.Medicine) error
	UpdateStock(id uint) error
	UpdateReduceStock(id uint) error
	UpdateActivate(id uint) error
	UpdateDeactivate(id uint) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type medicineRepository struct {
	db *gorm.DB
}

func NewMedicineRepository(db *gorm.DB) MedicineRepository {
	return &medicineRepository{
		db: db,
	}
}

func (r *medicineRepository) buildBaseQuery(query *dto.MedicinePaginationQuery) *gorm.DB {
	db := r.db.Model(&models.Medicine{})

	if query.Search != "" {
		pattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where(
			"name ILIKE ? OR generic_name ILIKE ? OR brand_name ILIKE ? OR type ILIKE ?",
			pattern, pattern, pattern, pattern)
	}

	if query.HasStock != nil {
		if *query.HasStock {
			db = db.Where("stock_quantity > ?", 0)
		} else {
			db = db.Where("stock_quantity = ?", 0)
		}
	}

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if query.Manufacturer != "" {
		db = db.Where("manufacturer = ?", query.Manufacturer)
	}

	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	if query.MinStock > 0 {
		db = db.Where("stock_quantity >= ?", query.MinStock)
	}

	if query.MaxStock > 0 {
		db = db.Where("stock_quantity <= ?", query.MaxStock)
	}

	return db
}

func (r *medicineRepository) List(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	var (
		medicines []models.Medicine
		total     int64
	)

	db := r.buildBaseQuery(query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "name"
	if query.SortBy != "" {
		sortBy = query.SortBy
	}

	sortDir := "asc"
	if query.SortDir == "desc" {
		sortDir = "desc"
	}

	db = db.Order(fmt.Sprintf("%s %s", sortBy, sortDir))

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&medicines).Error; err != nil {
		return nil, 0, err
	}

	return medicines, total, nil
}

func (r *medicineRepository) DeletedList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) ListByAvailable(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	queryForFilter := *query
	available := true
	queryForFilter.HasStock = &available
	return r.List(&queryForFilter)
}

func (r *medicineRepository) GetTotalStockValueByAvailable(query *dto.MedicinePaginationQuery) (float64, error) {
	var totalValue float64

	queryForFilter := *query
	available := true
	queryForFilter.HasStock = &available
	db := r.buildBaseQuery(&queryForFilter)

	// Calculate SUM(stock_quantity * price) for all available medicines
	if err := db.Select("COALESCE(SUM(stock_quantity * price), 0)").Row().Scan(&totalValue); err != nil {
		return 0, err
	}

	return totalValue, nil
}

func (r *medicineRepository) GetTotalStockQuantityByAvailable(query *dto.MedicinePaginationQuery) (int64, error) {
	var totalQuantity int64

	queryForFilter := *query
	available := true
	queryForFilter.HasStock = &available
	db := r.buildBaseQuery(&queryForFilter)

	// Calculate SUM(stock_quantity) for all available medicines
	if err := db.Select("COALESCE(SUM(stock_quantity), 0)").Row().Scan(&totalQuantity); err != nil {
		return 0, err
	}

	return totalQuantity, nil
}



func (r *medicineRepository) ListByLowStock(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	queryForFilter := *query
	
	// Set low stock threshold
	lowStockThreshold := 10
	queryForFilter.MaxStock = lowStockThreshold
	queryForFilter.MinStock = 0
	
	return r.List(&queryForFilter)
}

func (r *medicineRepository) GetTotalStockValueByLowStock(query *dto.MedicinePaginationQuery) (float64, error) {
	var totalValue float64

	queryForFilter := *query
	lowStockThreshold := 10
	queryForFilter.MaxStock = lowStockThreshold
	queryForFilter.MinStock = 0
	db := r.buildBaseQuery(&queryForFilter)

	// Calculate SUM(stock_quantity * price) for all low stock medicines
	if err := db.Select("COALESCE(SUM(stock_quantity * price), 0)").Row().Scan(&totalValue); err != nil {
		return 0, err
	}

	return totalValue, nil
}

func (r *medicineRepository) GetTotalStockQuantityByLowStock(query *dto.MedicinePaginationQuery) (int64, error) {
	var totalQuantity int64

	queryForFilter := *query
	lowStockThreshold := 10
	queryForFilter.MaxStock = lowStockThreshold
	queryForFilter.MinStock = 0
	db := r.buildBaseQuery(&queryForFilter)

	// Calculate SUM(stock_quantity) for all low stock medicines
	if err := db.Select("COALESCE(SUM(stock_quantity), 0)").Row().Scan(&totalQuantity); err != nil {
		return 0, err
	}

	return totalQuantity, nil
}

func (r *medicineRepository) ListByOutStock(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) ListByInactive(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) FindByID(id uint) (*models.Medicine, error) {
	var m models.Medicine
	err := r.db.First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, errors.New("medicine not found")
		}
		return nil, err
	}
	return &m, nil
}

func (r *medicineRepository) FindByName(name string) (*models.Medicine, error) {
	var m models.Medicine
	err := r.db.Where("name = ?", name).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			return nil, errors.New("medicine not found")
		}
		return nil, err
	}
	return &m, nil
}

func (r *medicineRepository) FindByType(types string) (*models.Medicine, error) {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) Search(criteria *dto.MedicineSearchCriteria, item *dto.MedicineSearchItem, response *dto.MedicineSearchResponse) ([]models.Medicine, int64, error) {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) Create(medicine *models.Medicine) error {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) Update(medicine *models.Medicine) error {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) UpdateStock(id uint) error {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) UpdateReduceStock(id uint) error {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) UpdateActivate(id uint) error {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) UpdateDeactivate(id uint) error {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) SoftDelete(id uint) error {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) Restore(id uint) error {
	panic("not implemented") // TODO: Implement
}

func (r *medicineRepository) HardDelete(id uint) error {
	panic("not implemented") // TODO: Implement
}
