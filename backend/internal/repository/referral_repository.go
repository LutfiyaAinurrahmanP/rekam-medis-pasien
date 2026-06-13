package repository

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type ReferralRepository interface {
	List(query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error)
	DeletedList(query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error)
	FindMyReferrals(patientID uint, status string) ([]models.Referral, error)
	FindByID(id uint) (*models.Referral, error)
	FindByIDUnscoped(id uint) (*models.Referral, error)
	FindByPatientID(patientID uint, query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error)
	FindByDoctorID(doctorID uint, query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error)
	Create(referral *models.Referral) error
	Update(referral *models.Referral) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	GenerateReferralNumber() (string, error)
}

type referralRepository struct {
	db *gorm.DB
}

func NewReferralRepository(db *gorm.DB) ReferralRepository {
	return &referralRepository{db: db}
}

func (r *referralRepository) buildQuery(baseQuery *gorm.DB, query dto.ReferralPaginationQuery) *gorm.DB {
	q := baseQuery

	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		q = q.Where("LOWER(referrals.referral_number) LIKE ? OR LOWER(referrals.reason) LIKE ?", searchTerm, searchTerm)
	}
	if query.Status != "" {
		q = q.Where("referrals.status = ?", query.Status)
	}
	if query.Priority != "" {
		q = q.Where("referrals.priority = ?", query.Priority)
	}
	if query.ReferralType != "" {
		q = q.Where("referrals.referral_type = ?", query.ReferralType)
	}
	if query.ReferringDoctorID != nil {
		q = q.Where("referrals.referring_doctor_id = ?", *query.ReferringDoctorID)
	}
	if query.PatientID != nil {
		q = q.Where("referrals.patient_id = ?", *query.PatientID)
	}
	if query.DoctorID != nil {
		// Could be referring doctor or referred to doctor
		q = q.Where("referrals.referring_doctor_id = ? OR referrals.referred_to_doctor_id = ?", *query.DoctorID, *query.DoctorID)
	}

	return q
}

func (r *referralRepository) preloadAll(q *gorm.DB) *gorm.DB {
	return q.Preload("Patient").
		Preload("MedicalRecord").
		Preload("ReferringDoctor.Specialization").
		Preload("ReferredToDepartment").
		Preload("ReferredToDoctor.Specialization")
}

func (r *referralRepository) List(query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error) {
	var referrals []models.Referral
	var totalItems int64

	baseQuery := r.db.Model(&models.Referral{})
	q := r.buildQuery(baseQuery, query)

	if err := q.Count(&totalItems).Error; err != nil {
		return nil, dto.ReferralPaginationMeta{}, err
	}

	offset := (query.Page - 1) * query.PageSize
	orderClause := fmt.Sprintf("referrals.%s %s", query.SortBy, query.SortDir)

	q = r.preloadAll(q)
	if err := q.Order(orderClause).Limit(query.PageSize).Offset(offset).Find(&referrals).Error; err != nil {
		return nil, dto.ReferralPaginationMeta{}, err
	}

	meta := dto.ReferralPaginationMeta{
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalItems: totalItems,
		TotalPages: int(math.Ceil(float64(totalItems) / float64(query.PageSize))),
	}

	return referrals, meta, nil
}

func (r *referralRepository) DeletedList(query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error) {
	var referrals []models.Referral
	var totalItems int64

	baseQuery := r.db.Unscoped().Model(&models.Referral{}).Where("referrals.deleted_at IS NOT NULL")
	q := r.buildQuery(baseQuery, query)

	if err := q.Count(&totalItems).Error; err != nil {
		return nil, dto.ReferralPaginationMeta{}, err
	}

	offset := (query.Page - 1) * query.PageSize
	orderClause := fmt.Sprintf("referrals.%s %s", query.SortBy, query.SortDir)

	q = r.preloadAll(q)
	if err := q.Order(orderClause).Limit(query.PageSize).Offset(offset).Find(&referrals).Error; err != nil {
		return nil, dto.ReferralPaginationMeta{}, err
	}

	meta := dto.ReferralPaginationMeta{
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalItems: totalItems,
		TotalPages: int(math.Ceil(float64(totalItems) / float64(query.PageSize))),
	}

	return referrals, meta, nil
}

func (r *referralRepository) FindMyReferrals(patientID uint, status string) ([]models.Referral, error) {
	var referrals []models.Referral
	q := r.db.Where("patient_id = ?", patientID)
	if status != "" {
		q = q.Where("status = ?", status)
	}

	q = r.preloadAll(q)
	if err := q.Order("created_at desc").Find(&referrals).Error; err != nil {
		return nil, err
	}
	return referrals, nil
}

func (r *referralRepository) FindByID(id uint) (*models.Referral, error) {
	var referral models.Referral
	q := r.db.Where("id = ?", id)
	q = r.preloadAll(q)
	if err := q.First(&referral).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &referral, nil
}

func (r *referralRepository) FindByIDUnscoped(id uint) (*models.Referral, error) {
	var referral models.Referral
	q := r.db.Unscoped().Where("id = ?", id)
	if err := q.First(&referral).Error; err != nil {
		return nil, err
	}
	return &referral, nil
}

func (r *referralRepository) FindByPatientID(patientID uint, query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error) {
	query.PatientID = &patientID
	return r.List(query)
}

func (r *referralRepository) FindByDoctorID(doctorID uint, query dto.ReferralPaginationQuery) ([]models.Referral, dto.ReferralPaginationMeta, error) {
	query.DoctorID = &doctorID
	return r.List(query)
}

func (r *referralRepository) Create(referral *models.Referral) error {
	return r.db.Create(referral).Error
}

func (r *referralRepository) Update(referral *models.Referral) error {
	return r.db.Save(referral).Error
}

func (r *referralRepository) SoftDelete(id uint) error {
	result := r.db.Delete(&models.Referral{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("referral not found")
	}

	return nil
}

func (r *referralRepository) Restore(id uint) error {
	return r.db.Unscoped().Model(&models.Referral{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *referralRepository) HardDelete(id uint) error {
	return r.db.Unscoped().Delete(&models.Referral{}, id).Error
}

func (r *referralRepository) GenerateReferralNumber() (string, error) {
	var count int64
	if err := r.db.Unscoped().Model(&models.Referral{}).Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("REF-2024-%06d", count+1), nil
}
