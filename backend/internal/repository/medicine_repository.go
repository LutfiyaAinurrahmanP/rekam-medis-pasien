package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type MedicineRepository interface {
	List(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	DeletedList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	AvailableList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	LowStockList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	OutStockList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	ActiveList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	InactiveList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error)
	FindByID(id uint) (*models.Medicine, error)
	FindByIDUnscoped(id uint) (*models.Medicine, error)
	Create(medicine *models.Medicine) error
	Update(medicine *models.Medicine) error
	AddStock(id uint, quantity int) error
	ReduceStock(id uint, quantity int) error
	Activate(id uint) error
	Deactivate(id uint) error
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

// --- Helper Functions ---

func escapeSearchMedicinePattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func applyMedicineListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "created_at"
	switch sortBy {
	case "name", "generic_name", "brand_name", "manufacturer", "stock_quantity", "price", "created_at", "updated_at":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("%s %s", column, direction))
}

func (r *medicineRepository) buildBaseQuery(query *dto.MedicinePaginationQuery, isDeleted bool) *gorm.DB {
	var db *gorm.DB
	if isDeleted {
		db = r.db.Unscoped().Model(&models.Medicine{}).Where("deleted_at IS NOT NULL")
	} else {
		db = r.db.Model(&models.Medicine{})
	}

	// 1. Search
	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchMedicinePattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? OR generic_name ILIKE ? OR brand_name ILIKE ? OR manufacturer ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// 2. Filters
	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if query.MedicineTypeID != nil {
		db = db.Where("medicine_type_id = ?", *query.MedicineTypeID)
	}


	if query.StockStatus != "" {
		switch query.StockStatus {
		case "available":
			db = db.Where("stock_quantity > ?", 0)
		case "low_stock":
			// Misalnya low stock itu > 0 dan <= 5
			db = db.Where("stock_quantity > ? AND stock_quantity <= ?", 0, 5)
		case "out_of_stock":
			db = db.Where("stock_quantity = ?", 0)
		}
	}

	return db
}

// --- Interface Implementations ---

func (r *medicineRepository) List(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	var medicines []models.Medicine
	var total int64

	db := r.buildBaseQuery(query, false)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyMedicineListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&medicines).Error; err != nil {
		return nil, 0, err
	}

	return medicines, total, nil
}

func (r *medicineRepository) DeletedList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	var medicines []models.Medicine
	var total int64

	db := r.buildBaseQuery(query, true)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyMedicineListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&medicines).Error; err != nil {
		return nil, 0, err
	}

	return medicines, total, nil
}

// Convenience methods that use List() under the hood
func (r *medicineRepository) AvailableList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	q := *query
	q.StockStatus = "available"
	return r.List(&q)
}

func (r *medicineRepository) LowStockList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	q := *query
	q.StockStatus = "low_stock"
	return r.List(&q)
}

func (r *medicineRepository) OutStockList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	q := *query
	q.StockStatus = "out_of_stock"
	return r.List(&q)
}

func (r *medicineRepository) ActiveList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	q := *query
	active := true
	q.IsActive = &active
	return r.List(&q)
}

func (r *medicineRepository) InactiveList(query *dto.MedicinePaginationQuery) ([]models.Medicine, int64, error) {
	q := *query
	active := false
	q.IsActive = &active
	return r.List(&q)
}

func (r *medicineRepository) FindByID(id uint) (*models.Medicine, error) {
	var m models.Medicine
	err := r.db.First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("medicine not found")
		}
		return nil, err
	}
	return &m, nil
}

func (r *medicineRepository) FindByIDUnscoped(id uint) (*models.Medicine, error) {
	var m models.Medicine
	err := r.db.Unscoped().First(&m, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("medicine not found")
		}
		return nil, err
	}
	return &m, nil
}

func (r *medicineRepository) Create(medicine *models.Medicine) error {
	return r.db.Create(medicine).Error
}

func (r *medicineRepository) Update(medicine *models.Medicine) error {
	return r.db.Save(medicine).Error
}

func (r *medicineRepository) AddStock(id uint, quantity int) error {
	result := r.db.Model(&models.Medicine{}).Where("id = ?", id).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity + ?", quantity))

	if result.Error != nil {
		return result.Error
	}	

	if result.RowsAffected == 0 {
		return errors.New("medicine not found")
	}

	return nil
}

func (r *medicineRepository) ReduceStock(id uint, quantity int) error {
	result := r.db.Model(&models.Medicine{}).Where("id = ?", id).
		UpdateColumn("stock_quantity", gorm.Expr("stock_quantity - ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medicine not found")
	}

	return nil
}

func (r *medicineRepository) Activate(id uint) error {
	result := r.db.Model(&models.Medicine{}).Where("id = ?", id).Update("is_active", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medicine not found")
	}

	return nil
}

func (r *medicineRepository) Deactivate(id uint) error {
	result := r.db.Model(&models.Medicine{}).Where("id = ?", id).Update("is_active", false)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medicine not found")
	}

	return nil
}

func (r *medicineRepository) SoftDelete(id uint) error {
	result := r.db.Delete(&models.Medicine{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medicine not found")
	}

	return nil
}

func (r *medicineRepository) Restore(id uint) error {
	result := r.db.Unscoped().Model(&models.Medicine{}).Where("id = ?", id).Update("deleted_at", nil)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medicine not found")
	}

	return nil
}

func (r *medicineRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.Medicine{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medicine not found")
	}

	return nil
}
