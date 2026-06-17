package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type AllergyRepository interface {
	List(query *dto.AllergyPaginationQuery) ([]models.Allergy, int64, error)
	FindByID(id uint) (*models.Allergy, error)
	Create(allergy *models.Allergy) error
	Update(allergy *models.Allergy) error
	Delete(id uint) error
}

type allergyRepository struct {
	db *gorm.DB
}

func NewAllergyRepository(db *gorm.DB) AllergyRepository {
	return &allergyRepository{db: db}
}

func applyAllergyListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "created_at"
	switch sortBy {
	case "id", "created_at", "allergen_name":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("allergies.%s %s", column, direction))
}

func (r *allergyRepository) buildBaseQuery(query *dto.AllergyPaginationQuery) *gorm.DB {
	db := r.db.Model(&models.Allergy{})

	if query.PatientID != nil {
		db = db.Where("patient_id = ?", *query.PatientID)
	}

	return db
}

func (r *allergyRepository) List(query *dto.AllergyPaginationQuery) ([]models.Allergy, int64, error) {
	var allergies []models.Allergy
	var total int64

	db := r.buildBaseQuery(query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyAllergyListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Patient").
		Offset(offset).Limit(query.PageSize).Find(&allergies).Error; err != nil {
		return nil, 0, err
	}

	return allergies, total, nil
}

func (r *allergyRepository) FindByID(id uint) (*models.Allergy, error) {
	var allergy models.Allergy
	err := r.db.Preload("Patient").First(&allergy, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("allergy not found")
		}
		return nil, err
	}
	return &allergy, nil
}

func (r *allergyRepository) Create(allergy *models.Allergy) error {
	return r.db.Create(allergy).Error
}

func (r *allergyRepository) Update(allergy *models.Allergy) error {
	return r.db.Save(allergy).Error
}

func (r *allergyRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Allergy{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("allergy not found")
	}

	return nil
}
