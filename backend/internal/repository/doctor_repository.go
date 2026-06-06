package repository

import (
	"errors"
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type DoctorRepository interface {
	List(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error)
	DeleteList(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error)
	FindByID(id uint) (*models.Doctor, error)
	FindByUserID(userID uint) (*models.Doctor, error)
	FindByDepartmentID(departmentID uint) (*models.Doctor, error)
	ActiveList(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error)
	InactiveList(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error)
	Create(doctor *models.Doctor) error
	Update(doctor *models.Doctor) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsEmployeeIDExists(employeeID string, excludeID ...uint) (bool, error)
}

type doctorRepository struct {
	db *gorm.DB
}

func NewDoctorRepository(db *gorm.DB) DoctorRepository {
	return &doctorRepository{
		db: db,
	}
}

func (r *doctorRepository) List(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error) {
	var (
		doctors []models.Doctor
		total   int64
	)

	db := r.db.Model(&models.Doctor{})

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("full_name ILIKE ?", searchPattern)
	}

	if query.SpecializationID != nil {
		db = db.Where("specialization_id = ?", query.SpecializationID)
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
	if err := db.Offset(offset).Limit(query.PageSize).Find(&doctors).Error; err != nil {
		return nil, 0, err
	}
	return doctors, total, nil
}

func (r *doctorRepository) DeleteList(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error) {
	var (
		deletedDoctors []models.Doctor
		total          int64
	)

	db := r.db.Unscoped().Model(&models.Doctor{}).Where("deleted_at IS NOT NULL")

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("full_name ILIKE ?", searchPattern)
	}

	if query.SpecializationID != nil {
		db = db.Where("specialization_id = ?", query.SpecializationID)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, nil
	}

	orderClause := query.SortBy
	if query.SortDir == "desc" {
		orderClause += " desc"
	} else {
		orderClause += " asc"
	}

	db = db.Order(orderClause)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&deletedDoctors).Error; err != nil {
		return nil, 0, err
	}
	return deletedDoctors, total, nil
}

func (r *doctorRepository) FindByID(id uint) (*models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.First(&doctor, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("doctor not found")
		}
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorRepository) FindByUserID(userID uint) (*models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.Where("user_id = ?", &userID).First(&doctor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("doctor not found")
		}
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorRepository) FindByDepartmentID(departmentID uint) (*models.Doctor, error) {
	var doctor models.Doctor
	err := r.db.Where("department_id = ?", departmentID).First(&doctor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("doctor not found")
		}
		return nil, err
	}
	return &doctor, nil
}

func (r *doctorRepository) ActiveList(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error) {
	var (
		doctors []models.Doctor
		total   int64
	)

	db := r.db.Model(&models.Doctor{}).Where("is_active = ?", true)

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("full_name ILIKE ?", searchPattern)
	}

	if query.SpecializationID != nil {
		db = db.Where("specialization_id = ?", query.SpecializationID)
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
	if err := db.Offset(offset).Limit(query.PageSize).Find(&doctors).Error; err != nil {
		return nil, 0, err
	}
	return doctors, total, nil
}

func (r *doctorRepository) InactiveList(query *dto.DoctorPaginationQuery) ([]models.Doctor, int64, error) {
	var (
		doctors []models.Doctor
		total   int64
	)

	db := r.db.Model(&models.Doctor{}).Where("is_active = ?", false)

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("full_name ILIKE ?", searchPattern)
	}

	if query.SpecializationID != nil {
		db = db.Where("specialization_id = ?", query.SpecializationID)
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
	if err := db.Offset(offset).Limit(query.PageSize).Find(&doctors).Error; err != nil {
		return nil, 0, err
	}
	return doctors, total, nil
}

func (r *doctorRepository) Create(doctor *models.Doctor) error {
	return r.db.Create(doctor).Error
}

func (r *doctorRepository) Update(doctor *models.Doctor) error {
	return r.db.Save(doctor).Error
}

func (r *doctorRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.Doctor{}, id).Error
}

func (r *doctorRepository) Restore(id uint) error {
	result := r.db.Model(&models.Doctor{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("doctor not found")
	}
	return nil
}

func (r *doctorRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.Doctor{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("doctor not found")
	}
	return nil
}

func (r *doctorRepository) IsEmployeeIDExists(employeeID string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Doctor{}).Where("employee_id = ?", employeeID)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	err := query.Count(&count).Error
	return count > 0, err
}
