package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type HospitalizationRepository interface {
	List(query *dto.HospitalizationPaginationQuery) ([]models.Hospitalization, int64, error)
	DeletedList(query *dto.HospitalizationPaginationQuery) ([]models.Hospitalization, int64, error)
	FindByID(id uint) (*models.Hospitalization, error)
	FindByIDUnscoped(id uint) (*models.Hospitalization, error)
	Create(hospitalization *models.Hospitalization) error
	Update(id uint, updates map[string]interface{}) error
	Discharge(id uint, updates map[string]interface{}) error
	Transfer(id uint, updates map[string]interface{}) error
	Activate(id uint, updates map[string]interface{}) error
	Deactivate(id uint, updates map[string]interface{}) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type hospitalizationRepository struct {
	db *gorm.DB
}

func NewHospitalizationRepository(db *gorm.DB) HospitalizationRepository {
	return &hospitalizationRepository{db: db}
}

func (r *hospitalizationRepository) applyFilters(db *gorm.DB, query *dto.HospitalizationPaginationQuery) *gorm.DB {
	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(admission_reason) LIKE ? OR LOWER(notes) LIKE ?", searchTerm, searchTerm)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if query.NotStatus != "" {
		db = db.Where("status != ?", query.NotStatus)
	}

	if query.PatientID != nil {
		db = db.Where("patient_id = ?", *query.PatientID)
	}

	if query.DoctorID != nil {
		db = db.Where("doctor_id = ?", *query.DoctorID)
	}

	if query.RoomID != nil {
		db = db.Where("room_id = ?", *query.RoomID)
	}

	return db
}

func (r *hospitalizationRepository) List(query *dto.HospitalizationPaginationQuery) ([]models.Hospitalization, int64, error) {
	var hospitalizations []models.Hospitalization
	var total int64

	db := r.db.Model(&models.Hospitalization{})
	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := fmt.Sprintf("%s %s", query.SortBy, query.SortDir)
	offset := (query.Page - 1) * query.PageSize

	err := db.Preload("Patient").
		Preload("Doctor.Specialization").
		Preload("Room.RoomType").
		Order(orderClause).
		Limit(query.PageSize).
		Offset(offset).
		Find(&hospitalizations).Error

	if err != nil {
		return nil, 0, err
	}

	return hospitalizations, total, nil
}

func (r *hospitalizationRepository) DeletedList(query *dto.HospitalizationPaginationQuery) ([]models.Hospitalization, int64, error) {
	var hospitalizations []models.Hospitalization
	var total int64

	db := r.db.Unscoped().Model(&models.Hospitalization{}).Where("deleted_at IS NOT NULL")
	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := fmt.Sprintf("%s %s", query.SortBy, query.SortDir)
	offset := (query.Page - 1) * query.PageSize

	err := db.Preload("Patient").
		Preload("Doctor.Specialization").
		Preload("Room.RoomType").
		Order(orderClause).
		Limit(query.PageSize).
		Offset(offset).
		Find(&hospitalizations).Error

	if err != nil {
		return nil, 0, err
	}

	return hospitalizations, total, nil
}

func (r *hospitalizationRepository) FindByID(id uint) (*models.Hospitalization, error) {
	var hospitalization models.Hospitalization
	err := r.db.Preload("Patient").
		Preload("Doctor.Specialization").
		Preload("Room.RoomType").
		First(&hospitalization, id).Error
	if err != nil {
		return nil, err
	}
	return &hospitalization, nil
}

func (r *hospitalizationRepository) FindByIDUnscoped(id uint) (*models.Hospitalization, error) {
	var hospitalization models.Hospitalization
	err := r.db.Unscoped().
		Preload("Patient").
		Preload("Doctor.Specialization").
		Preload("Room.RoomType").
		First(&hospitalization, id).Error
	if err != nil {
		return nil, err
	}
	return &hospitalization, nil
}

func (r *hospitalizationRepository) Create(hospitalization *models.Hospitalization) error {
	return r.db.Create(hospitalization).Error
}

func (r *hospitalizationRepository) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&models.Hospitalization{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("hospitalization not found")
	}
	return nil
}

func (r *hospitalizationRepository) Discharge(id uint, updates map[string]interface{}) error {
	return r.Update(id, updates)
}

func (r *hospitalizationRepository) Transfer(id uint, updates map[string]interface{}) error {
	return r.Update(id, updates)
}

func (r *hospitalizationRepository) Activate(id uint, updates map[string]interface{}) error {
	return r.Update(id, updates)
}

func (r *hospitalizationRepository) Deactivate(id uint, updates map[string]interface{}) error {
	return r.Update(id, updates)
}

func (r *hospitalizationRepository) SoftDelete(id uint) error {
	result := r.db.Delete(&models.Hospitalization{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("hospitalization not found")
	}
	return nil
}

func (r *hospitalizationRepository) Restore(id uint) error {
	result := r.db.Unscoped().Model(&models.Hospitalization{}).Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("hospitalization not found")
	}
	return nil
}

func (r *hospitalizationRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.Hospitalization{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("hospitalization not found")
	}
	return nil
}


