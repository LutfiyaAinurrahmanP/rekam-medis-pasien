package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type MedicalHistoryRepository interface {
	List(query *dto.MedicalHistoryPaginationQuery) ([]models.Patient, int64, error)
	FindByID(id uint) (*models.Patient, error)
}

type medicalHistoryRepository struct {
	db *gorm.DB
}

func NewMedicalHistoryRepository(db *gorm.DB) MedicalHistoryRepository {
	return &medicalHistoryRepository{
		db: db,
	}
}

func applyMedicalHistoryListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "created_at"
	switch sortBy {
	case "updated_at", "patient_code", "full_name":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("patients.%s %s", column, direction))
}

func (r *medicalHistoryRepository) buildBaseQuery(query *dto.MedicalHistoryPaginationQuery) *gorm.DB {
	db := r.db.Model(&models.Patient{})

	if query.PatientID != nil {
		db = db.Where("patients.id = ?", *query.PatientID)
	}

	return db
}

func (r *medicalHistoryRepository) List(query *dto.MedicalHistoryPaginationQuery) ([]models.Patient, int64, error) {
	var patients []models.Patient
	var total int64

	db := r.buildBaseQuery(query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyMedicalHistoryListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("MedicalAllergies").
		Preload("MedicalConditions").
		Preload("SurgicalHistories").
		Preload("FamilyHistories").
		Offset(offset).Limit(query.PageSize).Find(&patients).Error; err != nil {
		return nil, 0, err
	}

	return patients, total, nil
}

func (r *medicalHistoryRepository) FindByID(id uint) (*models.Patient, error) {
	var patient models.Patient
	err := r.db.Preload("MedicalAllergies").
		Preload("MedicalConditions").
		Preload("SurgicalHistories").
		Preload("FamilyHistories").
		First(&patient, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("medical history for patient not found")
		}
		return nil, err
	}
	return &patient, nil
}
