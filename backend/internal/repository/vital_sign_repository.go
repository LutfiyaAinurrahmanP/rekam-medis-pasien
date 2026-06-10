package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type VitalSignRepository interface {
	List(query *dto.VitalSignPaginationQuery) ([]models.VitalSign, int64, error)
	DeletedList(query *dto.VitalSignPaginationQuery) ([]models.VitalSign, int64, error)
	FindByID(id uint) (*models.VitalSign, error)
	FindByIDUnscoped(id uint) (*models.VitalSign, error)
	Create(vitalSign *models.VitalSign) error
	Update(vitalSign *models.VitalSign) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type vitalSignRepository struct {
	db *gorm.DB
}

func NewVitalSignRepository(db *gorm.DB) VitalSignRepository {
	return &vitalSignRepository{
		db: db,
	}
}

func applyVitalSignListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "created_at"
	switch sortBy {
	case "id", "recorded_at", "created_at": // Note: mapping sortby
		if sortBy == "recorded_at" {
			column = "measurement_date"
		} else {
			column = sortBy
		}
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("vital_signs.%s %s", column, direction))
}

func (r *vitalSignRepository) buildBaseQuery(query *dto.VitalSignPaginationQuery, isDeleted bool) *gorm.DB {
	var db *gorm.DB
	if isDeleted {
		db = r.db.Unscoped().Model(&models.VitalSign{}).Where("vital_signs.deleted_at IS NOT NULL")
	} else {
		db = r.db.Model(&models.VitalSign{})
	}

	if query.MedicalRecordID != nil {
		db = db.Where("medical_record_id = ?", *query.MedicalRecordID)
	}

	return db
}

func (r *vitalSignRepository) List(query *dto.VitalSignPaginationQuery) ([]models.VitalSign, int64, error) {
	var vitalSigns []models.VitalSign
	var total int64

	db := r.buildBaseQuery(query, false)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyVitalSignListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("MedicalRecord").
		Offset(offset).Limit(query.PageSize).Find(&vitalSigns).Error; err != nil {
		return nil, 0, err
	}

	return vitalSigns, total, nil
}

func (r *vitalSignRepository) DeletedList(query *dto.VitalSignPaginationQuery) ([]models.VitalSign, int64, error) {
	var vitalSigns []models.VitalSign
	var total int64

	db := r.buildBaseQuery(query, true)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyVitalSignListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("MedicalRecord").
		Offset(offset).Limit(query.PageSize).Find(&vitalSigns).Error; err != nil {
		return nil, 0, err
	}

	return vitalSigns, total, nil
}

func (r *vitalSignRepository) FindByID(id uint) (*models.VitalSign, error) {
	var v models.VitalSign
	err := r.db.Preload("MedicalRecord").First(&v, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("vital signs not found")
		}
		return nil, err
	}
	return &v, nil
}

func (r *vitalSignRepository) FindByIDUnscoped(id uint) (*models.VitalSign, error) {
	var v models.VitalSign
	err := r.db.Unscoped().Preload("MedicalRecord").First(&v, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("vital signs not found")
		}
		return nil, err
	}
	return &v, nil
}

func (r *vitalSignRepository) FindByMedicalRecordID(recordID uint) (*models.VitalSign, error) {
	var v models.VitalSign
	err := r.db.Preload("MedicalRecord").Where("medical_record_id = ?", recordID).First(&v).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("vital signs not found")
		}
		return nil, err
	}
	return &v, nil
}

func (r *vitalSignRepository) Create(vitalSign *models.VitalSign) error {
	return r.db.Create(vitalSign).Error
}

func (r *vitalSignRepository) Update(vitalSign *models.VitalSign) error {
	return r.db.Save(vitalSign).Error
}

func (r *vitalSignRepository) SoftDelete(id uint) error {
	result := r.db.Delete(&models.VitalSign{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("vital signs not found")
	}
	return nil
}

func (r *vitalSignRepository) Restore(id uint) error {
	result := r.db.Unscoped().Model(&models.VitalSign{}).Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("vital signs not found")
	}
	return nil
}

func (r *vitalSignRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.VitalSign{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("vital signs not found")
	}
	return nil
}
