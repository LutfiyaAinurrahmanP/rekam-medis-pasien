package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type PrescriptionRepository interface {
	List(query *dto.PrescriptionPaginationQuery) ([]models.Prescription, int64, error)
	DeletedList(query *dto.PrescriptionPaginationQuery) ([]models.Prescription, int64, error)
	FindByID(id uint) (*models.Prescription, error)
	FindByIDUnscoped(id uint) (*models.Prescription, error)
	Create(prescription *models.Prescription) error
	Update(prescription *models.Prescription) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error

	FindItemByID(id uint) (*models.PrescriptionItem, error)
	CreateItem(item *models.PrescriptionItem) error
	UpdateItem(item *models.PrescriptionItem) error
	DeleteItem(id uint) error
}

type prescriptionRepository struct {
	db *gorm.DB
}

func NewPrescriptionRepository(db *gorm.DB) PrescriptionRepository {
	return &prescriptionRepository{
		db: db,
	}
}

func applyPrescriptionListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "created_at"
	switch sortBy {
	case "id", "prescription_date", "status", "created_at":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("prescriptions.%s %s", column, direction))
}

func (r *prescriptionRepository) buildBaseQuery(query *dto.PrescriptionPaginationQuery, isDeleted bool) *gorm.DB {
	var db *gorm.DB
	if isDeleted {
		db = r.db.Unscoped().Model(&models.Prescription{}).Where("prescriptions.deleted_at IS NOT NULL")
	} else {
		db = r.db.Model(&models.Prescription{})
	}

	if query.DoctorID != nil {
		db = db.Where("doctor_id = ?", *query.DoctorID)
	}
	if query.MedicalRecordID != nil {
		db = db.Where("medical_record_id = ?", *query.MedicalRecordID)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Search != "" {
		db = db.Where("notes ILIKE ?", "%"+query.Search+"%")
	}

	return db
}

func (r *prescriptionRepository) List(query *dto.PrescriptionPaginationQuery) ([]models.Prescription, int64, error) {
	var prescriptions []models.Prescription
	var total int64

	db := r.buildBaseQuery(query, false)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyPrescriptionListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("MedicalRecord").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Items").
		Preload("Items.Medicine").
		Offset(offset).Limit(query.PageSize).Find(&prescriptions).Error; err != nil {
		return nil, 0, err
	}

	return prescriptions, total, nil
}

func (r *prescriptionRepository) DeletedList(query *dto.PrescriptionPaginationQuery) ([]models.Prescription, int64, error) {
	var prescriptions []models.Prescription
	var total int64

	db := r.buildBaseQuery(query, true)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyPrescriptionListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("MedicalRecord").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Items").
		Preload("Items.Medicine").
		Offset(offset).Limit(query.PageSize).Find(&prescriptions).Error; err != nil {
		return nil, 0, err
	}

	return prescriptions, total, nil
}

func (r *prescriptionRepository) FindByID(id uint) (*models.Prescription, error) {
	var p models.Prescription
	err := r.db.Preload("MedicalRecord").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Items").
		Preload("Items.Medicine").
		First(&p, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("prescription not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *prescriptionRepository) FindByIDUnscoped(id uint) (*models.Prescription, error) {
	var p models.Prescription
	err := r.db.Unscoped().Preload("MedicalRecord").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Items").
		Preload("Items.Medicine").
		First(&p, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("prescription not found")
		}
		return nil, err
	}
	return &p, nil
}

func (r *prescriptionRepository) Create(prescription *models.Prescription) error {
	return r.db.Create(prescription).Error
}

func (r *prescriptionRepository) Update(prescription *models.Prescription) error {
	return r.db.Save(prescription).Error
}

func (r *prescriptionRepository) SoftDelete(id uint) error {
	result := r.db.Delete(&models.Prescription{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("prescription not found")
	}
	return nil
}

func (r *prescriptionRepository) Restore(id uint) error {
	result := r.db.Unscoped().Model(&models.Prescription{}).Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("prescription not found")
	}
	return nil
}

func (r *prescriptionRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.Prescription{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("prescription not found")
	}
	return nil
}

func (r *prescriptionRepository) FindItemByID(id uint) (*models.PrescriptionItem, error) {
	var item models.PrescriptionItem
	if err := r.db.Preload("Medicine").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("prescription item not found")
		}
		return nil, err
	}
	return &item, nil
}

func (r *prescriptionRepository) CreateItem(item *models.PrescriptionItem) error {
	return r.db.Create(item).Error
}

func (r *prescriptionRepository) UpdateItem(item *models.PrescriptionItem) error {
	return r.db.Save(item).Error
}

func (r *prescriptionRepository) DeleteItem(id uint) error {
	result := r.db.Unscoped().Delete(&models.PrescriptionItem{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("prescription item not found")
	}
	return nil
}
