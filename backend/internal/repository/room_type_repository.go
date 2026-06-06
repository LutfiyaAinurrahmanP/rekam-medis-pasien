package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type RoomTypeRepository interface {
	List(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error)
	DeletedList(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error)
	FindByID(id uint) (*models.RoomType, error)
	Create(roomType *models.RoomType) error
	Update(roomType *models.RoomType) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsNameExists(name string, excludeID ...uint) (bool, error)
	IsCodeExists(code string, excludeID ...uint) (bool, error)
	ActiveList(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error)
	InactiveList(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error)
	Activate(id uint) error
	Deactivate(id uint) error
}

type roomTypeRepository struct {
	db *gorm.DB
}

func NewRoomTypeRepository(db *gorm.DB) RoomTypeRepository {
	return &roomTypeRepository{
		db: db,
	}
}
func escapeSearchRoomTypePattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func applyRoomTypeListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
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

func (r *roomTypeRepository) List(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error) {
	// Define variables
	var (
		roomTypes []models.RoomType
		total                 int64
	)

	// Base query
	db := r.db.Model(&models.RoomType{})

	// Search functionality
	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchRoomTypePattern(searchTerm)
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
	db = applyRoomTypeListOrder(db, query.SortBy, query.SortDir)

	// Apply pagination
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&roomTypes).Error; err != nil {
		return nil, 0, err
	}

	// Return results
	return roomTypes, total, nil
}

func (r *roomTypeRepository) DeletedList(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error) {
	var (
		roomTypes []models.RoomType
		total                 int64
	)

	db := r.db.Unscoped().Model(&models.RoomType{}).Where("deleted_at IS NOT NULL")

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchRoomTypePattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\'",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyRoomTypeListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&roomTypes).Error; err != nil {
		return nil, 0, err
	}

	return roomTypes, total, nil
}

func (r *roomTypeRepository) FindByID(id uint) (*models.RoomType, error) {
	var roomType models.RoomType
	err := r.db.First(&roomType, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("room type not found")
		}
		return nil, err
	}
	return &roomType, nil
}

func (r *roomTypeRepository) Create(roomType *models.RoomType) error {
	return r.db.Create(roomType).Error
}

func (r *roomTypeRepository) Update(roomType *models.RoomType) error {
	return r.db.Save(roomType).Error
}

func (r *roomTypeRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.RoomType{}, id).Error
}

func (r *roomTypeRepository) Restore(id uint) error {
	result := r.db.Model(&models.RoomType{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("room type not found")
	}
	return nil
}

func (r *roomTypeRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.RoomType{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("room type not found")
	}
	return nil
}

func (r *roomTypeRepository) IsNameExists(name string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.RoomType{}).Where("name = ?", name)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *roomTypeRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.RoomType{}).Where("code = ?", code)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *roomTypeRepository) ActiveList(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error) {
	var (
		roomTypes []models.RoomType
		total                 int64
	)

	db := r.db.Model(&models.RoomType{}).Where("is_active = ?", true)

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchRoomTypePattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"(name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\')",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyRoomTypeListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&roomTypes).Error; err != nil {
		return nil, 0, err
	}

	return roomTypes, total, nil
}

func (r *roomTypeRepository) InactiveList(query *dto.RoomTypePaginationQuery) ([]models.RoomType, int64, error) {
	var (
		roomTypes []models.RoomType
		total                 int64
	)

	db := r.db.Model(&models.RoomType{}).Where("is_active = ?", false)

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchRoomTypePattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"(name ILIKE ? ESCAPE '\\' OR code ILIKE ? ESCAPE '\\')",
			searchPattern, searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyRoomTypeListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&roomTypes).Error; err != nil {
		return nil, 0, err
	}

	return roomTypes, total, nil
}

func (r *roomTypeRepository) Activate(id uint) error {
	result := r.db.Model(&models.RoomType{}).Where("id = ?", id).Update("is_active", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("room type not found")
	}
	return nil
}

func (r *roomTypeRepository) Deactivate(id uint) error {
	result := r.db.Model(&models.RoomType{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("room type not found")
	}
	return nil
}
