package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	List(query *dto.RoomPaginationQuery) ([]models.Room, int64, error)
	DeleteList(query *dto.RoomPaginationQuery) ([]models.Room, int64, error)
	FindAvailableRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error)
	FindOccupiedRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error)
	FindActiveRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error)
	FindInactiveRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error)
	FindByID(id uint) (*models.Room, error)
	Create(room *models.Room) error
	Update(room *models.Room) error
	Activate(id uint) error
	Deactivate(id uint) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsRoomNumberExists(roomNumber string, excludeID ...uint) (bool, error)
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{
		db: db,
	}
}

func escapeSearchRoomPattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func applyRoomListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
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

func (r *roomRepository) List(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	var (
		rooms []models.Room
		total int64
	)

	db := r.db.Model(&models.Room{})

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchRoomPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"room_number ILIKE ? ESCAPE '\\'",
			searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyRoomListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}
	return rooms, total, nil
}

func (r *roomRepository) DeleteList(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	var (
		rooms []models.Room
		total int64
	)

	db := r.db.Unscoped().Model(&models.Room{}).Where("deleted_at IS NOT NULL")

	if searchTerm := strings.TrimSpace(query.Search); searchTerm != "" {
		escaped := escapeSearchRoomPattern(searchTerm)
		searchPattern := "%" + escaped + "%"

		db = db.Where(
			"room_number ILIKE ? ESCAPE '\\'",
			searchPattern,
		)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyRoomListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}
	return rooms, total, nil
}

func (r *roomRepository) FindAvailableRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	var (
		rooms []models.Room
		total int64
	)

	db := r.db.Model(&models.Room{}).Where("available_beds > 0 AND is_active = ?", true)

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("room_number ILIKE ? OR room_type ILIKE ?", searchPattern, searchPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := query.SortBy
	if query.SortDir == "desc" {
		orderClause += " desc"
	} else {
		orderClause += " asc"
	}
	db = db.Order(orderClause)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}
	return rooms, total, nil
}

func (r *roomRepository) FindOccupiedRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	var (
		rooms []models.Room
		total int64
	)

	db := r.db.Model(&models.Room{}).Where("available_beds = 0")

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("room_number ILIKE ? OR room_type ILIKE ?", searchPattern, searchPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := query.SortBy
	if query.SortDir == "desc" {
		orderClause += " desc"
	} else {
		orderClause += " asc"
	}
	db = db.Order(orderClause)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}
	return rooms, total, nil
}

func (r *roomRepository) FindActiveRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	var (
		rooms []models.Room
		total int64
	)

	db := r.db.Model(&models.Room{}).Where("is_active = ?", true)

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("room_number ILIKE ? OR room_type ILIKE ?", searchPattern, searchPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := query.SortBy
	if query.SortDir == "desc" {
		orderClause += " desc"
	} else {
		orderClause += " asc"
	}
	db = db.Order(orderClause)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}
	return rooms, total, nil
}

func (r *roomRepository) FindInactiveRooms(query *dto.RoomPaginationQuery) ([]models.Room, int64, error) {
	var (
		rooms []models.Room
		total int64
	)

	db := r.db.Model(&models.Room{}).Where("is_active = ?", false)

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("room_number ILIKE ? OR room_type ILIKE ?", searchPattern, searchPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := query.SortBy
	if query.SortDir == "desc" {
		orderClause += " desc"
	} else {
		orderClause += " asc"
	}
	db = db.Order(orderClause)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&rooms).Error; err != nil {
		return nil, 0, err
	}
	return rooms, total, nil
}

func (r *roomRepository) FindByID(id uint) (*models.Room, error) {
	var room models.Room
	err := r.db.First(&room, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("room not found")
		}
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) Create(room *models.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) Update(room *models.Room) error {
	return r.db.Save(room).Error
}

func (r *roomRepository) Activate(id uint) error {
	result := r.db.Model(&models.Room{}).Where("id = ?", id).Update("is_active", true)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("room not found")
	}
	return nil
}

func (r *roomRepository) Deactivate(id uint) error {
	result := r.db.Model(&models.Room{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("room not found")
	}
	return nil
}

func (r *roomRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.Room{}, id).Error
}

func (r *roomRepository) Restore(id uint) error {
	result := r.db.Model(&models.Room{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("room not found")
	}
	return nil
}

func (r *roomRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.Room{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("room not found")
	}
	return nil
}

func (r *roomRepository) IsRoomNumberExists(roomNumber string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Room{}).Where("room_number = ?", roomNumber)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}
