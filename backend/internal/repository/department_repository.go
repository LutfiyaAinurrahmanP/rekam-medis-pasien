package repository

import (
	"errors"
	"fmt"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type DepartmentRepository interface {
	List(query *dto.DepartmentPaginationQuery) ([]models.Department, int64, error)
	DeleteList(query *dto.DepartmentPaginationQuery) ([]models.Department, int64, error)
	FindById(id uint) (*models.Department, error)
	FindByName(name string) (*models.Department, error)
	FindByCode(code string) (*models.Department, error)
	Create(department *models.Department) error
	Update(department *models.Department) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsCodeExists(code string, excludeID ...uint) (bool, error)
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{
		db: db,
	}
}

func (r *departmentRepository) List(query *dto.DepartmentPaginationQuery) ([]models.Department, int64, error) {
	var (
		departments []models.Department
		total       int64
	)

	db := r.db.Model(&models.Department{})

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("name ILIKE ?  OR code ILIKE ?", searchPattern, searchPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&departments).Error; err != nil {
		return nil, 0, err
	}
	return departments, total, nil
}

func (r *departmentRepository) DeleteList(query *dto.DepartmentPaginationQuery) ([]models.Department, int64, error) {
	var (
		departments []models.Department
		total       int64
	)

	db := r.db.Unscoped().Model(&models.Department{}).Where("deleted_at IS NOT NULL")

	if query.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", query.Search)
		db = db.Where("name ILIKE ? OR code ILIKE ?", searchPattern, searchPattern)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&departments).Error; err != nil {
		return nil, 0, err
	}
	return departments, total, nil
}

func (r *departmentRepository) FindById(id uint) (*models.Department, error) {
	var department models.Department
	err := r.db.First(&department, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) FindByName(name string) (*models.Department, error) {
	var department models.Department
	err := r.db.Where("name = ?", name).First(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("name not found")
		}
		return nil, err
	}

	return &department, nil
}

func (r *departmentRepository) FindByCode(code string) (*models.Department, error) {
	var department models.Department
	err := r.db.Where("code = ?", code).First(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("code not found")
		}
		return nil, err
	}
	return &department, nil
}

func (r *departmentRepository) Create(department *models.Department) error {
	return r.db.Create(department).Error
}

func (r *departmentRepository) Update(department *models.Department) error {
	return r.db.Save(department).Error
}

func (r *departmentRepository) SoftDelete(id uint) error {
	return r.db.Delete(&models.Department{}, id).Error
}

func (r *departmentRepository) Restore(id uint) error {
	return r.db.Model(&models.Department{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *departmentRepository) HardDelete(id uint) error {
	return r.db.Unscoped().Delete(&models.Department{}, id).Error
}

func (r *departmentRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
	var count int64
	query := r.db.Model(&models.Department{}).Where("code = ?", code)

	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	err := query.Count(&count).Error
	return count > 0, err
}
