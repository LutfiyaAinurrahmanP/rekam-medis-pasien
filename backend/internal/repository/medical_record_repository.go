package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type MedicalRecordRepository interface {
	List(query *dto.MedicalRecordPaginationQuery) ([]models.MedicalRecord, int64, error)
	DeletedList(query *dto.MedicalRecordPaginationQuery) ([]models.MedicalRecord, int64, error)
	FindByID(id uint) (*models.MedicalRecord, error)
	FindByIDUnscoped(id uint) (*models.MedicalRecord, error)
	Create(record *models.MedicalRecord) error
	Update(record *models.MedicalRecord) error
	Finalize(id uint) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type medicalRecordRepository struct {
	db *gorm.DB
}

func NewMedicalRecordRepository(db *gorm.DB) MedicalRecordRepository {
	return &medicalRecordRepository{
		db: db,
	}
}

func applyMedicalRecordListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "visit_date"
	switch sortBy {
	case "created_at", "visit_date":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("medical_records.%s %s", column, direction))
}

func (r *medicalRecordRepository) buildBaseQuery(query *dto.MedicalRecordPaginationQuery, isDeleted bool) *gorm.DB {
	var db *gorm.DB
	if isDeleted {
		db = r.db.Unscoped().Model(&models.MedicalRecord{}).Where("medical_records.deleted_at IS NOT NULL")
	} else {
		db = r.db.Model(&models.MedicalRecord{})
	}

	// 1. Date Range
	if query.DateFrom != "" {
		db = db.Where("visit_date >= ?", query.DateFrom)
	}
	if query.DateTo != "" {
		db = db.Where("visit_date <= ?", query.DateTo)
	}

	// 2. Filters
	if query.PatientID != nil {
		db = db.Where("patient_id = ?", *query.PatientID)
	}
	if query.DoctorID != nil {
		db = db.Where("doctor_id = ?", *query.DoctorID)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// 3. Department Join Filter
	if query.DepartmentID != nil {
		db = db.Joins("JOIN doctors ON doctors.id = medical_records.doctor_id").
			Where("doctors.department_id = ?", *query.DepartmentID)
	}

	return db
}

func (r *medicalRecordRepository) List(query *dto.MedicalRecordPaginationQuery) ([]models.MedicalRecord, int64, error) {
	var records []models.MedicalRecord
	var total int64

	db := r.buildBaseQuery(query, false)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyMedicalRecordListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Patient").
		Preload("Patient.MedicalAllergies").
		Preload("Patient.MedicalConditions").
		Preload("Patient.SurgicalHistories").
		Preload("Patient.FamilyHistories").
		Preload("VitalSign").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Doctor.Department").
		Offset(offset).Limit(query.PageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *medicalRecordRepository) DeletedList(query *dto.MedicalRecordPaginationQuery) ([]models.MedicalRecord, int64, error) {
	var records []models.MedicalRecord
	var total int64

	db := r.buildBaseQuery(query, true)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyMedicalRecordListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Patient").
		Preload("Patient.MedicalAllergies").
		Preload("Patient.MedicalConditions").
		Preload("Patient.SurgicalHistories").
		Preload("Patient.FamilyHistories").
		Preload("VitalSign").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Doctor.Department").
		Offset(offset).Limit(query.PageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *medicalRecordRepository) FindByID(id uint) (*models.MedicalRecord, error) {
	var m models.MedicalRecord
	err := r.db.Preload("Patient").
		Preload("Patient.MedicalAllergies").
		Preload("Patient.MedicalConditions").
		Preload("Patient.SurgicalHistories").
		Preload("Patient.FamilyHistories").
		Preload("VitalSign").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Doctor.Department").
		First(&m, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("medical record not found")
		}
		return nil, err
	}
	return &m, nil
}

func (r *medicalRecordRepository) FindByIDUnscoped(id uint) (*models.MedicalRecord, error) {
	var m models.MedicalRecord
	err := r.db.Unscoped().Preload("Patient").
		Preload("Patient.MedicalAllergies").
		Preload("Patient.MedicalConditions").
		Preload("Patient.SurgicalHistories").
		Preload("Patient.FamilyHistories").
		Preload("VitalSign").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Doctor.Department").
		First(&m, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("medical record not found")
		}
		return nil, err
	}
	return &m, nil
}

func (r *medicalRecordRepository) Create(record *models.MedicalRecord) error {
	return r.db.Create(record).Error
}

func (r *medicalRecordRepository) Update(record *models.MedicalRecord) error {
	return r.db.Save(record).Error
}

func (r *medicalRecordRepository) Finalize(id uint) error {
	result := r.db.Model(&models.MedicalRecord{}).Where("id = ? AND status = ?", id, "draft").Updates(map[string]interface{}{
		"status": "finalized",
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("medical record not found or already finalized")
	}
	return nil
}

func (r *medicalRecordRepository) SoftDelete(id uint) error {
	result := r.db.Delete(&models.MedicalRecord{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medical record not found")
	}

	return nil
}

func (r *medicalRecordRepository) Restore(id uint) error {
	result := r.db.Unscoped().Model(&models.MedicalRecord{}).Where("id = ?", id).Update("deleted_at", nil)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medical record not found")
	}

	return nil
}

func (r *medicalRecordRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.MedicalRecord{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medical record not found")
	}

	return nil
}
