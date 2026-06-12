package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type MedicalConditionRepository interface {
	List(query *dto.MedicalConditionPaginationQuery) ([]models.MedicalCondition, int64, error)
	FindByID(id uint) (*models.MedicalCondition, error)
	Create(condition *models.MedicalCondition) error
	Update(condition *models.MedicalCondition) error
	Delete(id uint) error
}

type medicalConditionRepository struct {
	db *gorm.DB
}

func NewMedicalConditionRepository(db *gorm.DB) MedicalConditionRepository {
	return &medicalConditionRepository{db: db}
}

func applyMedicalConditionListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "created_at"
	switch sortBy {
	case "id", "created_at", "condition_name":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("medical_conditions.%s %s", column, direction))
}

func (r *medicalConditionRepository) buildBaseQuery(query *dto.MedicalConditionPaginationQuery) *gorm.DB {
	db := r.db.Model(&models.MedicalCondition{})

	if query.PatientID != nil {
		db = db.Where("patient_id = ?", *query.PatientID)
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	return db
}

func (r *medicalConditionRepository) List(query *dto.MedicalConditionPaginationQuery) ([]models.MedicalCondition, int64, error) {
	var conditions []models.MedicalCondition
	var total int64

	db := r.buildBaseQuery(query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyMedicalConditionListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Patient").
		Offset(offset).Limit(query.PageSize).Find(&conditions).Error; err != nil {
		return nil, 0, err
	}

	return conditions, total, nil
}

func (r *medicalConditionRepository) FindByID(id uint) (*models.MedicalCondition, error) {
	var condition models.MedicalCondition
	err := r.db.Preload("Patient").First(&condition, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("medical condition not found")
		}
		return nil, err
	}
	return &condition, nil
}

func (r *medicalConditionRepository) Create(condition *models.MedicalCondition) error {
	return r.db.Create(condition).Error
}

func (r *medicalConditionRepository) Update(condition *models.MedicalCondition) error {
	return r.db.Save(condition).Error
}

func (r *medicalConditionRepository) Delete(id uint) error {
	result := r.db.Delete(&models.MedicalCondition{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("medical condition not found")
	}

	return nil
}
