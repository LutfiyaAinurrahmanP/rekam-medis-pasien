package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type LabTestRepository interface {
	List(query *dto.LabTestPaginationQuery) ([]models.LabTest, int64, error)
	DeletedList(query *dto.LabTestPaginationQuery) ([]models.LabTest, int64, error)
	FindByID(id uint) (*models.LabTest, error)
	FindByIDUnscoped(id uint) (*models.LabTest, error)
	Create(labTest *models.LabTest) error
	Update(labTest *models.LabTest) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type labTestRepository struct {
	db *gorm.DB
}

func NewLabTestRepository(db *gorm.DB) LabTestRepository {
	return &labTestRepository{
		db: db,
	}
}

func applyLabTestListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "created_at"
	switch sortBy {
	case "id", "order_date", "result_date", "status", "created_at":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("lab_tests.%s %s", column, direction))
}

func (r *labTestRepository) buildBaseQuery(query *dto.LabTestPaginationQuery, isDeleted bool) *gorm.DB {
	var db *gorm.DB
	if isDeleted {
		db = r.db.Unscoped().Model(&models.LabTest{}).Where("lab_tests.deleted_at IS NOT NULL")
	} else {
		db = r.db.Model(&models.LabTest{})
	}

	if query.MedicalRecordID != nil {
		db = db.Where("lab_tests.medical_record_id = ?", *query.MedicalRecordID)
	}
	if query.TestTypeID != nil {
		db = db.Where("lab_tests.test_type_id = ?", *query.TestTypeID)
	}
	if query.OrderedByDoctorID != nil {
		db = db.Where("lab_tests.ordered_by_doctor_id = ?", *query.OrderedByDoctorID)
	}
	if query.Status != "" {
		db = db.Where("lab_tests.status = ?", query.Status)
	}
	if query.NotStatus != "" {
		db = db.Where("lab_tests.status != ?", query.NotStatus)
	}

	if query.Search != "" {
		search := "%" + query.Search + "%"
		db = db.Where("lab_tests.notes ILIKE ? OR lab_tests.result_value ILIKE ?", search, search)
	}

	return db
}

func (r *labTestRepository) List(query *dto.LabTestPaginationQuery) ([]models.LabTest, int64, error) {
	var labTests []models.LabTest
	var total int64

	db := r.buildBaseQuery(query, false)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyLabTestListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("MedicalRecord").
		Preload("TestType").
		Preload("OrderedByDoctor").
		Preload("OrderedByDoctor.Specialization").
		Offset(offset).Limit(query.PageSize).Find(&labTests).Error; err != nil {
		return nil, 0, err
	}

	return labTests, total, nil
}

func (r *labTestRepository) DeletedList(query *dto.LabTestPaginationQuery) ([]models.LabTest, int64, error) {
	var labTests []models.LabTest
	var total int64

	db := r.buildBaseQuery(query, true)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyLabTestListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("MedicalRecord").
		Preload("TestType").
		Preload("OrderedByDoctor").
		Preload("OrderedByDoctor.Specialization").
		Offset(offset).Limit(query.PageSize).Find(&labTests).Error; err != nil {
		return nil, 0, err
	}

	return labTests, total, nil
}

func (r *labTestRepository) FindByID(id uint) (*models.LabTest, error) {
	var labTest models.LabTest
	err := r.db.Preload("MedicalRecord").
		Preload("TestType").
		Preload("OrderedByDoctor").
		Preload("OrderedByDoctor.Specialization").
		First(&labTest, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lab test not found")
		}
		return nil, err
	}
	return &labTest, nil
}

func (r *labTestRepository) FindByIDUnscoped(id uint) (*models.LabTest, error) {
	var labTest models.LabTest
	err := r.db.Unscoped().Preload("MedicalRecord").
		Preload("TestType").
		Preload("OrderedByDoctor").
		Preload("OrderedByDoctor.Specialization").
		First(&labTest, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lab test not found")
		}
		return nil, err
	}
	return &labTest, nil
}

func (r *labTestRepository) Create(labTest *models.LabTest) error {
	return r.db.Create(labTest).Error
}

func (r *labTestRepository) Update(labTest *models.LabTest) error {
	return r.db.Save(labTest).Error
}

func (r *labTestRepository) SoftDelete(id uint) error {
	result := r.db.Delete(&models.LabTest{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("lab test not found")
	}

	return nil
}

func (r *labTestRepository) Restore(id uint) error {
	result := r.db.Unscoped().Model(&models.LabTest{}).Where("id = ?", id).Update("deleted_at", nil)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("lab test not found")
	}

	return nil
}

func (r *labTestRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.LabTest{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("lab test not found")
	}

	return nil
}
