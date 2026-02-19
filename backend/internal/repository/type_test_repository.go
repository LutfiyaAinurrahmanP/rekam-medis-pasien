package repository

import (
	"errors"
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type TypeTestRepository interface {
	List(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error)
	ListActive(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error)
	ListInactive(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error)
	DeleteList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error)
	Search(query *dto.TypeTestSearchQuery) ([]models.TypeTest, int64, error)
	FindByID(id uint) (*models.TypeTest, error)
	FindByCode(code string) (*models.TypeTest, error)
	FindByCategory(category string, query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error)
	Create(typeTest *models.TypeTest) error
	Update(typeTest *models.TypeTest) error
	Activate(id uint) error
	Deactivate(id uint) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsCodeExists(code string, excludeID ...uint) (bool, error)
	CountActiveByCategory() ([]dto.ActiveTypeTestCategorySummary, error)
}

type typeTestRepository struct {
	db *gorm.DB
}

func NewTypeTestRepository(db *gorm.DB) TypeTestRepository {
	return &typeTestRepository{db: db}
}

func (r *typeTestRepository) buildBaseQuery(query *dto.TypeTestPaginationQuery) *gorm.DB {
	db := r.db.Model(&models.TypeTest{})

	if query.Search != "" {
		pattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?", pattern, pattern, pattern)
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}
	if query.MinPrice != nil {
		db = db.Where("price >= ?", *query.MinPrice)
	}
	if query.MaxPrice != nil {
		db = db.Where("price <= ?", *query.MaxPrice)
	}
	return db
}

func (r *typeTestRepository) List(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	var (
		typeTests []models.TypeTest
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
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTests).Error; err != nil {
		return nil, 0, err
	}
	return typeTests, total, nil
}

func (r *typeTestRepository) ListActive(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	var (
		typeTests []models.TypeTest
		total     int64
	)

	active := true
	query.IsActive = &active
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
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTests).Error; err != nil {
		return nil, 0, err
	}
	return typeTests, total, nil
}

func (r *typeTestRepository) ListInactive(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	var (
		typeTests []models.TypeTest
		total     int64
	)

	inactive := false
	query.IsActive = &inactive
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
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTests).Error; err != nil {
		return nil, 0, err
	}
	return typeTests, total, nil
}

func (r *typeTestRepository) DeleteList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	var (
		typeTests []models.TypeTest
		total     int64
	)

	db := r.db.Unscoped().Model(&models.TypeTest{}).Where("deleted_at IS NOT NULL")

	if query.Search != "" {
		pattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("name ILIKE ? OR code ILIKE ?", pattern, pattern)
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortBy := "deleted_at"
	sortDir := "desc"
	db = db.Order(fmt.Sprintf("%s %s", sortBy, sortDir))

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTests).Error; err != nil {
		return nil, 0, err
	}
	return typeTests, total, nil
}

func (r *typeTestRepository) Search(query *dto.TypeTestSearchQuery) ([]models.TypeTest, int64, error) {
	var (
		typeTests []models.TypeTest
		total     int64
	)

	db := r.db.Model(&models.TypeTest{})

	if query.Keyword != "" {
		pattern := fmt.Sprintf("%%%s%%", query.Keyword)
		db = db.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?", pattern, pattern, pattern)
	}
	if query.Category != "" {
		db = db.Where("category = ?", query.Category)
	}
	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}
	if query.MinPrice != nil {
		db = db.Where("price >= ?", *query.MinPrice)
	}
	if query.MaxPrice != nil {
		db = db.Where("price <= ?", *query.MaxPrice)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("name asc")
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTests).Error; err != nil {
		return nil, 0, err
	}
	return typeTests, total, nil
}

func (r *typeTestRepository) FindByID(id uint) (*models.TypeTest, error) {
	var t models.TypeTest
	err := r.db.First(&t, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("test type not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *typeTestRepository) FindByCode(code string) (*models.TypeTest, error) {
	var t models.TypeTest
	err := r.db.Where("code = ?", code).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("test type not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *typeTestRepository) FindByCategory(category string, query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	var (
		typeTests []models.TypeTest
		total     int64
	)

	db := r.db.Model(&models.TypeTest{}).Where("category = ?", category)

	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("name asc")
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTests).Error; err != nil {
		return nil, 0, err
	}
	return typeTests, total, nil
}

func (r *typeTestRepository) Create(typeTest *models.TypeTest) error {
	return r.db.Create(typeTest).Error
}

func (r *typeTestRepository) Update(typeTest *models.TypeTest) error {
	return r.db.Save(typeTest).Error
}

func (r *typeTestRepository) Activate(id uint) error {
	result := r.db.Model(&models.TypeTest{}).Where("id = ?", id).Update("is_active", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("test type not found")
	}
	return nil
}

func (r *typeTestRepository) Deactivate(id uint) error {
	result := r.db.Model(&models.TypeTest{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("test type not found")
	}
	return nil
}

func (r *typeTestRepository) SoftDelete(id uint) error {
	result := r.db.Delete(&models.TypeTest{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("test type not found")
	}
	return nil
}

func (r *typeTestRepository) Restore(id uint) error {
	result := r.db.Model(&models.TypeTest{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("test type not found")
	}
	return nil
}

func (r *typeTestRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.TypeTest{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("test type not found")
	}
	return nil
}

func (r *typeTestRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Unscoped().Model(&models.TypeTest{}).Where("code = ?", code)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *typeTestRepository) CountActiveByCategory() ([]dto.ActiveTypeTestCategorySummary, error) {
	var results []dto.ActiveTypeTestCategorySummary
	err := r.db.Model(&models.TypeTest{}).
		Select("category, COUNT(*) as count").
		Where("is_active = ?", true).
		Group("category").
		Order("category asc").
		Scan(&results).Error
	return results, err
}
