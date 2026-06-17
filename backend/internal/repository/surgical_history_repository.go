package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type SurgicalHistoryRepository interface {
	List(query *dto.SurgicalHistoryPaginationQuery) ([]models.SurgicalHistory, int64, error)
	FindByID(id uint) (*models.SurgicalHistory, error)
	Create(history *models.SurgicalHistory) error
	Update(history *models.SurgicalHistory) error
	Delete(id uint) error
}

type surgicalHistoryRepository struct {
	db *gorm.DB
}

func NewSurgicalHistoryRepository(db *gorm.DB) SurgicalHistoryRepository {
	return &surgicalHistoryRepository{db: db}
}

func applySurgicalHistoryListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "surgery_date"
	switch sortBy {
	case "id", "created_at", "surgery_date":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("surgical_histories.%s %s", column, direction))
}

func (r *surgicalHistoryRepository) buildBaseQuery(query *dto.SurgicalHistoryPaginationQuery) *gorm.DB {
	db := r.db.Model(&models.SurgicalHistory{})

	if query.PatientID != nil {
		db = db.Where("patient_id = ?", *query.PatientID)
	}

	return db
}

func (r *surgicalHistoryRepository) List(query *dto.SurgicalHistoryPaginationQuery) ([]models.SurgicalHistory, int64, error) {
	var histories []models.SurgicalHistory
	var total int64

	db := r.buildBaseQuery(query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applySurgicalHistoryListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Patient").
		Offset(offset).Limit(query.PageSize).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *surgicalHistoryRepository) FindByID(id uint) (*models.SurgicalHistory, error) {
	var history models.SurgicalHistory
	err := r.db.Preload("Patient").First(&history, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("surgical history not found")
		}
		return nil, err
	}
	return &history, nil
}

func (r *surgicalHistoryRepository) Create(history *models.SurgicalHistory) error {
	return r.db.Create(history).Error
}

func (r *surgicalHistoryRepository) Update(history *models.SurgicalHistory) error {
	return r.db.Save(history).Error
}

func (r *surgicalHistoryRepository) Delete(id uint) error {
	result := r.db.Delete(&models.SurgicalHistory{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("surgical history not found")
	}

	return nil
}
