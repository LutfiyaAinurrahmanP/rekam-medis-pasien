package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type TypeTestRepository interface {
	List(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error)
	ActiveList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error)
	InactiveList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error)
	DeletedList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error)
	FindByID(id uint) (*models.TypeTest, error)
	Create(typeTest *models.TypeTest) error
	Update(typeTest *models.TypeTest) error
	Activate(id uint) error
	Deactivate(id uint) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsCodeExists(code string, excludeID ...uint) (bool, error)
}

type typeTestRepository struct {
	db *gorm.DB
}

func NewTypeTestRepository(db *gorm.DB) TypeTestRepository {
	return &typeTestRepository{db: db}
}

func escapeSearchTypeTestPattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func applyTypeTestListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
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

func (r *typeTestRepository) List(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	var (
		typeTests []models.TypeTest
		total     int64
	)

	db := r.db.Model(&models.TypeTest{})

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchTypeTestPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	db = applyTypeTestListOrder(db, query.SortBy, query.SortDir)

	// Apply pagination
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTests).Error; err != nil {
		return nil, 0, err
	}

	// Return results
	return typeTests, total, nil
}

func (r *typeTestRepository) ActiveList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	var (
		typeTests []models.TypeTest
		total     int64
	)

	db := r.db.Model(&models.TypeTest{}).Where("is_active = ?", true)

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchTypeTestPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyTypeTestListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTests).Error; err != nil {
		return nil, 0, err
	}

	return typeTests, total, nil
}

func (r *typeTestRepository) InactiveList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	var (
		typeTests []models.TypeTest
		total     int64
	)

	db := r.db.Model(&models.TypeTest{}).Where("is_active = ?", false)

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchTypeTestPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyTypeTestListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&typeTests).Error; err != nil {
		return nil, 0, err
	}
	return typeTests, total, nil
}

func (r *typeTestRepository) DeletedList(query *dto.TypeTestPaginationQuery) ([]models.TypeTest, int64, error) {
	var (
		typeTests []models.TypeTest
		total     int64
	)

	db := r.db.Unscoped().Model(&models.TypeTest{}).Where("deleted_at IS NOT NULL")

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchTypeTestPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyTypeTestListOrder(db, query.SortBy, query.SortDir)

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
