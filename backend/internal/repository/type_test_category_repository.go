package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type TypeTestCategoryRepository interface {
	List(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error)
	DeletedList(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error)
	FindByID(id uint) (*models.TypeTestCategory, error)
	Create(typeTestCategory *models.TypeTestCategory) error
	Update(typeTestCategory *models.TypeTestCategory) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsNameExists(name string, excludeID ...uint) (bool, error)
	IsCodeExists(code string, excludeID ...uint) (bool, error)
	ActiveList(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error)
	InactiveList(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error)
	Activate(id uint) error
	Deactivate(id uint) error
}

type typeTestCategoryRepository struct {
	db *gorm.DB
}

func NewTypeTestCategoryRepository(db *gorm.DB) TypeTestCategoryRepository {
	return &typeTestCategoryRepository{
		db: db,
	}
}
func escapeSearchTypeTestCategoryPattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func applyTypeTestCategoryListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
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

func (r *typeTestCategoryRepository) List(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error) {
	// Define variables
	var (
		typeTestCategories []models.TypeTestCategory
		total              int64
	)

	// Base query
	db := r.db.Model(&models.TypeTestCategory{})

	// Search functionality
	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchTypeTestCategoryPattern(searchTerm)
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
	db = applyTypeTestCategoryListOrder(db, query.SortBy, query.SortDir)

	// Apply pagination
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTestCategories).Error; err != nil {
		return nil, 0, err
	}

	// Return results
	return typeTestCategories, total, nil
}

func (r *typeTestCategoryRepository) DeletedList(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error) {
	var (
		typeTestCategories []models.TypeTestCategory
		total              int64
	)

	db := r.db.Unscoped().Model(&models.TypeTestCategory{}).Where("deleted_at IS NOT NULL")

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchTypeTestCategoryPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyTypeTestCategoryListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTestCategories).Error; err != nil {
		return nil, 0, err
	}

	return typeTestCategories, total, nil
}

func (r *typeTestCategoryRepository) FindByID(id uint) (*models.TypeTestCategory, error) {
	var typeTestCategory models.TypeTestCategory
	err := r.db.First(&typeTestCategory, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("type test category not found")
		}
		return nil, err
	}
	return &typeTestCategory, nil
}

func (r *typeTestCategoryRepository) Create(typeTestCategory *models.TypeTestCategory) error {
	return r.db.Create(typeTestCategory).Error
}

func (r *typeTestCategoryRepository) Update(typeTestCategory *models.TypeTestCategory) error {
	return r.db.Save(typeTestCategory).Error
}

func (r *typeTestCategoryRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.TypeTestCategory{}, id).Error
}

func (r *typeTestCategoryRepository) Restore(id uint) error {
	result := r.db.Model(&models.TypeTestCategory{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("type test category not found")
	}
	return nil
}

func (r *typeTestCategoryRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.TypeTestCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("type test category not found")
	}
	return nil
}

func (r *typeTestCategoryRepository) IsNameExists(name string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.TypeTestCategory{}).Where("name = ?", name)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *typeTestCategoryRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.TypeTestCategory{}).Where("code = ?", code)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *typeTestCategoryRepository) ActiveList(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error) {
	var (
		typeTestCategories []models.TypeTestCategory
		total              int64
	)

	db := r.db.Model(&models.TypeTestCategory{}).Where("is_active = ?", true)

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchTypeTestCategoryPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"(name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\')",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyTypeTestCategoryListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTestCategories).Error; err != nil {
		return nil, 0, err
	}

	return typeTestCategories, total, nil
}

func (r *typeTestCategoryRepository) InactiveList(query *dto.TypeTestCategoryPaginationQuery) ([]models.TypeTestCategory, int64, error) {
	var (
		typeTestCategories []models.TypeTestCategory
		total              int64
	)

	db := r.db.Model(&models.TypeTestCategory{}).Where("is_active = ?", false)

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchTypeTestCategoryPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"(name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\')",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyTypeTestCategoryListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTestCategories).Error; err != nil {
		return nil, 0, err
	}

	return typeTestCategories, total, nil
}

func (r *typeTestCategoryRepository) Activate(id uint) error {
	result := r.db.Model(&models.TypeTestCategory{}).Where("id = ?", id).Update("is_active", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("type test category not found")
	}
	return nil
}

func (r *typeTestCategoryRepository) Deactivate(id uint) error {
	result := r.db.Model(&models.TypeTestCategory{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("type test category not found")
	}
	return nil
}
