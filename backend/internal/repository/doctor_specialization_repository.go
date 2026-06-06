package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type DoctorSpecializationRepository interface {
	List(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error)
	DeletedList(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error)
	FindByID(id uint) (*models.DoctorSpecialization, error)
	Create(doctorSpecialization *models.DoctorSpecialization) error
	Update(doctorSpecialization *models.DoctorSpecialization) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsNameExists(name string, excludeID ...uint) (bool, error)
	IsCodeExists(code string, excludeID ...uint) (bool, error)
	ActiveList(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error)
	InactiveList(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error)
	Activate(id uint) error
	Deactivate(id uint) error
}

type doctorSpecializationRepository struct {
	db *gorm.DB
}

func NewDoctorSpecializationRepository(db *gorm.DB) DoctorSpecializationRepository {
	return &doctorSpecializationRepository{
		db: db,
	}
}
func escapeSearchDoctorSpecializationPattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func applyDoctorSpecializationListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
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

func (r *doctorSpecializationRepository) List(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error) {
	// Define variables
	var (
		doctorSpecializations []models.DoctorSpecialization
		total                 int64
	)

	// Base query
	db := r.db.Model(&models.DoctorSpecialization{})

	// Search functionality
	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchDoctorSpecializationPattern(searchTerm)
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
	db = applyDoctorSpecializationListOrder(db, query.SortBy, query.SortDir)

	// Apply pagination
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&doctorSpecializations).Error; err != nil {
		return nil, 0, err
	}

	// Return results
	return doctorSpecializations, total, nil
}

func (r *doctorSpecializationRepository) DeletedList(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error) {
	var (
		doctorSpecializations []models.DoctorSpecialization
		total                 int64
	)

	db := r.db.Unscoped().Model(&models.DoctorSpecialization{}).Where("deleted_at IS NOT NULL")

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchDoctorSpecializationPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyDoctorSpecializationListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&doctorSpecializations).Error; err != nil {
		return nil, 0, err
	}

	return doctorSpecializations, total, nil
}

func (r *doctorSpecializationRepository) FindByID(id uint) (*models.DoctorSpecialization, error) {
	var doctorSpecialization models.DoctorSpecialization
	err := r.db.First(&doctorSpecialization, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("doctor specialization not found")
		}
		return nil, err
	}
	return &doctorSpecialization, nil
}

func (r *doctorSpecializationRepository) Create(doctorSpecialization *models.DoctorSpecialization) error {
	return r.db.Create(doctorSpecialization).Error
}

func (r *doctorSpecializationRepository) Update(doctorSpecialization *models.DoctorSpecialization) error {
	return r.db.Save(doctorSpecialization).Error
}

func (r *doctorSpecializationRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.DoctorSpecialization{}, id).Error
}

func (r *doctorSpecializationRepository) Restore(id uint) error {
	result := r.db.Model(&models.DoctorSpecialization{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("doctor specialization not found")
	}
	return nil
}

func (r *doctorSpecializationRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.DoctorSpecialization{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("doctor specialization not found")
	}
	return nil
}

func (r *doctorSpecializationRepository) IsNameExists(name string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.DoctorSpecialization{}).Where("name = ?", name)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *doctorSpecializationRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.DoctorSpecialization{}).Where("code = ?", code)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *doctorSpecializationRepository) ActiveList(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error) {
	var (
		doctorSpecializations []models.DoctorSpecialization
		total                 int64
	)

	db := r.db.Model(&models.DoctorSpecialization{}).Where("is_active = ?", true)

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchDoctorSpecializationPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"(name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\')",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyDoctorSpecializationListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&doctorSpecializations).Error; err != nil {
		return nil, 0, err
	}

	return doctorSpecializations, total, nil
}

func (r *doctorSpecializationRepository) InactiveList(query *dto.DoctorSpecializationPaginationQuery) ([]models.DoctorSpecialization, int64, error) {
	var (
		doctorSpecializations []models.DoctorSpecialization
		total                 int64
	)

	db := r.db.Model(&models.DoctorSpecialization{}).Where("is_active = ?", false)

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchDoctorSpecializationPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"(name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\')",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyDoctorSpecializationListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&doctorSpecializations).Error; err != nil {
		return nil, 0, err
	}

	return doctorSpecializations, total, nil
}

func (r *doctorSpecializationRepository) Activate(id uint) error {
	result := r.db.Model(&models.DoctorSpecialization{}).Where("id = ?", id).Update("is_active", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("doctor specialization not found")
	}
	return nil
}

func (r *doctorSpecializationRepository) Deactivate(id uint) error {
	result := r.db.Model(&models.DoctorSpecialization{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("doctor specialization not found")
	}
	return nil
}
