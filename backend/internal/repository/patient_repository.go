package repository

import (
	"errors"
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type PatientRepository interface {
	List(query *dto.PatientPaginationQuery) ([]models.Patient, int64, error)
	DeleteList(query *dto.PatientPaginationQuery) ([]models.Patient, int64, error)
	FindById(id uint) (*models.Patient, error)
	FindByUserID(userID uint) (*models.Patient, error)
	FindByCode(code string) (*models.Patient, error)
	Create(patient *models.Patient) error
	Update(patient *models.Patient) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsCodeExists(code string, excludeID ...uint) (bool, error)
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{
		db: db,
	}
}

func (r patientRepository) List(query *dto.PatientPaginationQuery) ([]models.Patient, int64, error) {
	var (
		patients []models.Patient
		total    int64
	)

	db := r.db.Model(&models.Patient{})

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("full_name ILIKE ? OR patient_code ILIKE ? OR phone ILIKE ? OR email ILIKE ?", searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if query.Gender != "" {
		db = db.Where("gender = ?", query.Gender)
	}

	if query.BloodType != "" {
		db = db.Where("blood_type = ?", query.BloodType)
	}

	if query.InsuranceProvider != "" {
		db = db.Where("insurance_provider ILIKE ?", fmt.Sprintf("%%%s%%", query.InsuranceProvider))
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
	if err := db.Offset(offset).Limit(query.PageSize).Find(&patients).Error; err != nil {
		return nil, 0, err
	}
	return patients, total, nil
}

func (r patientRepository) DeleteList(query *dto.PatientPaginationQuery) ([]models.Patient, int64, error) {
	var (
		patients []models.Patient
		total    int64
	)

	db := r.db.Unscoped().Model(&models.Patient{}).Where("deleted_at IS NOT NULL")

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("full_name ILIKE ? OR patient_code ILIKE ? OR phone ILIKE ? OR email ILIKE ?", searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if query.Gender != "" {
		db = db.Where("gender = ?", query.Gender)
	}

	if query.BloodType != "" {
		db = db.Where("blood_type = ?", query.BloodType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&patients).Error; err != nil {
		return nil, 0, err
	}
	return patients, total, nil
}

func (r patientRepository) FindById(id uint) (*models.Patient, error) {
	var patient models.Patient
	err := r.db.First(&patient, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, err
	}
	return &patient, nil
}

func (r patientRepository) FindByUserID(userID uint) (*models.Patient, error) {
	var patient models.Patient
	err := r.db.Where("user_id = ?", &userID).First(&patient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, err
	}
	return &patient, nil
}

func (r patientRepository) FindByCode(code string) (*models.Patient, error) {
	var patient models.Patient
	err := r.db.Where("patient_code = ?", code).First(&patient).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("patient not found")
		}
		return nil, err
	}
	return &patient, nil
}

func (r patientRepository) Create(patient *models.Patient) error {
	return r.db.Create(patient).Error
}

func (r patientRepository) Update(patient *models.Patient) error {
	return r.db.Save(patient).Error
}

func (r patientRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.Patient{}, id).Error
}

func (r patientRepository) Restore(id uint) error {
	result := r.db.Model(&models.Patient{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("patient not found")
	}
	return nil
}

func (r patientRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.Patient{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("patient not found")
	}
	return nil
}

func (r patientRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Patient{}).Where("patient_code = ?", code)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}
