package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type MedicineTypeRepository interface {
	List(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error)
	DeletedList(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error)
	FindByID(id uint) (*models.MedicineType, error)
	Create(medicineType *models.MedicineType) error
	Update(medicineType *models.MedicineType) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsNameExists(name string, excludeID ...uint) (bool, error)
	IsCodeExists(code string, excludeID ...uint) (bool, error)
	ActiveList(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error)
	InactiveList(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error)
	Activate(id uint) error
	Deactivate(id uint) error
}

type medicineTypeRepository struct {
	db *gorm.DB
}

func NewMedicineTypeRepository(db *gorm.DB) MedicineTypeRepository {
	return &medicineTypeRepository{
		db: db,
	}
}
func escapeSearchMedicineTypePattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func applyMedicineTypeListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "created_at"
	switch sortBy {
	case "name", "code", "created_at":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("%s %s", column, direction))
}

func (r *medicineTypeRepository) List(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error) {
	// Define variables
	var (
		medicineTypes []models.MedicineType
		total         int64
	)

	// Base query
	db := r.db.Model(&models.MedicineType{})

	// Search functionality
	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchMedicineTypePattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	// Count total records for pagination
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	db = applyMedicineTypeListOrder(db, query.SortBy, query.SortDir)

	// Apply pagination
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&medicineTypes).Error; err != nil {
		return nil, 0, err
	}

	// Return results
	return medicineTypes, total, nil
}

func (r *medicineTypeRepository) DeletedList(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error) {
	var (
		medicineTypes []models.MedicineType
		total         int64
	)

	db := r.db.Unscoped().Model(&models.MedicineType{}).Where("deleted_at IS NOT NULL")

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchMedicineTypePattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyMedicineTypeListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&medicineTypes).Error; err != nil {
		return nil, 0, err
	}

	return medicineTypes, total, nil
}

func (r *medicineTypeRepository) FindByID(id uint) (*models.MedicineType, error) {
	var medicineType models.MedicineType
	err := r.db.First(&medicineType, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("medicine type not found")
		}
		return nil, err
	}
	return &medicineType, nil
}

func (r *medicineTypeRepository) Create(medicineType *models.MedicineType) error {
	return r.db.Create(medicineType).Error
}

func (r *medicineTypeRepository) Update(medicineType *models.MedicineType) error {
	return r.db.Save(medicineType).Error
}

func (r *medicineTypeRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.MedicineType{}, id).Error
}

func (r *medicineTypeRepository) Restore(id uint) error {
	result := r.db.Model(&models.MedicineType{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medicine type not found")
	}
	return nil
}

func (r *medicineTypeRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.MedicineType{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medicine type not found")
	}
	return nil
}

func (r *medicineTypeRepository) IsNameExists(name string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.MedicineType{}).Where("name = ?", name)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *medicineTypeRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.MedicineType{}).Where("code = ?", code)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *medicineTypeRepository) ActiveList(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error) {
	var (
		medicineTypes []models.MedicineType
		total         int64
	)

	db := r.db.Model(&models.MedicineType{}).Where("is_active = ?", true)

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchMedicineTypePattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyMedicineTypeListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&medicineTypes).Error; err != nil {
		return nil, 0, err
	}

	return medicineTypes, total, nil
}

func (r *medicineTypeRepository) InactiveList(query *dto.MedicineTypePaginationQuery) ([]models.MedicineType, int64, error) {
	var (
		medicineTypes []models.MedicineType
		total         int64
	)

	db := r.db.Model(&models.MedicineType{}).Where("is_active = ?", false)

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchMedicineTypePattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyMedicineTypeListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&medicineTypes).Error; err != nil {
		return nil, 0, err
	}

	return medicineTypes, total, nil
}

func (r *medicineTypeRepository) Activate(id uint) error {
	result := r.db.Model(&models.MedicineType{}).Where("id = ?", id).Update("is_active", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("medicine type not found")
	}
	return nil
}

func (r *medicineTypeRepository) Deactivate(id uint) error {
	result := r.db.Model(&models.MedicineType{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("medicine type not found")
	}
	return nil
}
