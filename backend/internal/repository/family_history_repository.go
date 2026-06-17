package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type FamilyHistoryRepository interface {
	List(query *dto.FamilyHistoryPaginationQuery) ([]models.FamilyHistory, int64, error)
	FindByID(id uint) (*models.FamilyHistory, error)
	Create(history *models.FamilyHistory) error
	Update(history *models.FamilyHistory) error
	Delete(id uint) error
}

type familyHistoryRepository struct {
	db *gorm.DB
}

func NewFamilyHistoryRepository(db *gorm.DB) FamilyHistoryRepository {
	return &familyHistoryRepository{db: db}
}

func applyFamilyHistoryListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "created_at"
	switch sortBy {
	case "id", "created_at", "family_member":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("family_histories.%s %s", column, direction))
}

func (r *familyHistoryRepository) buildBaseQuery(query *dto.FamilyHistoryPaginationQuery) *gorm.DB {
	db := r.db.Model(&models.FamilyHistory{})

	if query.PatientID != nil {
		db = db.Where("patient_id = ?", *query.PatientID)
	}

	return db
}

func (r *familyHistoryRepository) List(query *dto.FamilyHistoryPaginationQuery) ([]models.FamilyHistory, int64, error) {
	var histories []models.FamilyHistory
	var total int64

	db := r.buildBaseQuery(query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyFamilyHistoryListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Patient").
		Offset(offset).Limit(query.PageSize).Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *familyHistoryRepository) FindByID(id uint) (*models.FamilyHistory, error) {
	var history models.FamilyHistory
	err := r.db.Preload("Patient").First(&history, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("family history not found")
		}
		return nil, err
	}
	return &history, nil
}

func (r *familyHistoryRepository) Create(history *models.FamilyHistory) error {
	return r.db.Create(history).Error
}

func (r *familyHistoryRepository) Update(history *models.FamilyHistory) error {
	return r.db.Save(history).Error
}

func (r *familyHistoryRepository) Delete(id uint) error {
	result := r.db.Delete(&models.FamilyHistory{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("family history not found")
	}

	return nil
}
