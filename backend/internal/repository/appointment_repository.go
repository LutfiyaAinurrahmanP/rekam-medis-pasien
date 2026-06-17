package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	List(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error)
	DeletedList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error)
	UpcomingList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error)
	PastList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error)
	TodayList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error)
	FindByID(id uint) (*models.Appointment, error)
	FindByIDUnscoped(id uint) (*models.Appointment, error)
	Create(appointment *models.Appointment) error
	Update(appointment *models.Appointment) error
	Confirm(id uint) error
	Start(id uint) error
	Complete(id uint) error
	Cancel(id uint, reason string) error
	Reschedule(id uint, newDate, newTime string) error
	NoShow(id uint) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{
		db: db,
	}
}

func applyAppointmentListOrder(db *gorm.DB, sortBy, sortDir string) *gorm.DB {
	column := "appointment_date"
	switch sortBy {
	case "created_at", "appointment_date", "status", "duration_minutes":
		column = sortBy
	}

	direction := "DESC"
	if strings.EqualFold(sortDir, "asc") {
		direction = "ASC"
	}

	return db.Order(fmt.Sprintf("appointments.%s %s", column, direction))
}

func (r *appointmentRepository) buildBaseQuery(query *dto.AppointmentPaginationQuery, isDeleted bool) *gorm.DB {
	var db *gorm.DB
	if isDeleted {
		db = r.db.Unscoped().Model(&models.Appointment{}).Where("appointments.deleted_at IS NOT NULL")
	} else {
		db = r.db.Model(&models.Appointment{})
	}

	// 1. Date Range
	if query.Date != "" {
		db = db.Where("appointment_date = ?", query.Date)
	}
	if query.DateFrom != "" {
		db = db.Where("appointment_date >= ?", query.DateFrom)
	}
	if query.DateTo != "" {
		db = db.Where("appointment_date <= ?", query.DateTo)
	}
	if query.DaysAhead > 0 {
		today := time.Now().Format("2006-01-02")
		future := time.Now().AddDate(0, 0, query.DaysAhead).Format("2006-01-02")
		db = db.Where("appointment_date >= ? AND appointment_date <= ?", today, future)
	}
	if query.DaysBack > 0 {
		today := time.Now().Format("2006-01-02")
		past := time.Now().AddDate(0, 0, -query.DaysBack).Format("2006-01-02")
		db = db.Where("appointment_date <= ? AND appointment_date >= ?", today, past)
	}

	// 2. Filters
	if query.PatientID != nil {
		db = db.Where("patient_id = ?", *query.PatientID)
	}
	if query.DoctorID != nil {
		db = db.Where("doctor_id = ?", *query.DoctorID)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// 3. Department Join Filter
	if query.DepartmentID != nil {
		db = db.Joins("JOIN doctors ON doctors.id = appointments.doctor_id").
			Where("doctors.department_id = ?", *query.DepartmentID)
	}

	return db
}

func (r *appointmentRepository) List(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	db := r.buildBaseQuery(query, false)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyAppointmentListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Patient").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Doctor.Department").
		Offset(offset).Limit(query.PageSize).Find(&appointments).Error; err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

func (r *appointmentRepository) DeletedList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	db := r.buildBaseQuery(query, true)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = applyAppointmentListOrder(db, query.SortBy, query.SortDir)

	offset := (query.Page - 1) * query.PageSize
	if err := db.Preload("Patient").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Doctor.Department").
		Offset(offset).Limit(query.PageSize).Find(&appointments).Error; err != nil {
		return nil, 0, err
	}

	return appointments, total, nil
}

func (r *appointmentRepository) UpcomingList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error) {
	q := *query
	if q.DaysAhead == 0 {
		q.DaysAhead = 7 // Default 7 days
	}
	return r.List(&q)
}

func (r *appointmentRepository) PastList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error) {
	q := *query
	if q.DaysBack == 0 {
		q.DaysBack = 30 // Default 30 days
	}
	return r.List(&q)
}

func (r *appointmentRepository) TodayList(query *dto.AppointmentPaginationQuery) ([]models.Appointment, int64, error) {
	q := *query
	q.Date = time.Now().Format("2006-01-02")
	return r.List(&q)
}

func (r *appointmentRepository) FindByID(id uint) (*models.Appointment, error) {
	var a models.Appointment
	err := r.db.Preload("Patient").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Doctor.Department").
		First(&a, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *appointmentRepository) FindByIDUnscoped(id uint) (*models.Appointment, error) {
	var a models.Appointment
	err := r.db.Unscoped().Preload("Patient").
		Preload("Doctor").
		Preload("Doctor.Specialization").
		Preload("Doctor.Department").
		First(&a, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("appointment not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *appointmentRepository) Create(appointment *models.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *appointmentRepository) Update(appointment *models.Appointment) error {
	return r.db.Save(appointment).Error
}

func (r *appointmentRepository) Confirm(id uint) error {
	result := r.db.Model(&models.Appointment{}).Where("id = ?", id).Update("status", "confirmed")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("appointment not found")
	}
	return nil
}

func (r *appointmentRepository) Start(id uint) error {
	result := r.db.Model(&models.Appointment{}).Where("id = ?", id).Update("status", "in_progress")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("appointment not found")
	}
	return nil
}

func (r *appointmentRepository) Complete(id uint) error {
	result := r.db.Model(&models.Appointment{}).Where("id = ?", id).Update("status", "completed")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("appointment not found")
	}
	return nil
}

func (r *appointmentRepository) Cancel(id uint, reason string) error {
	result := r.db.Model(&models.Appointment{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": "cancelled",
		"reason": reason,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("appointment not found")
	}
	return nil
}

func (r *appointmentRepository) Reschedule(id uint, newDate, newTime string) error {
	result := r.db.Model(&models.Appointment{}).Where("id = ?", id).Updates(map[string]interface{}{
		"appointment_date": newDate,
		"appointment_time": newTime,
		"status":           "scheduled",
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("appointment not found")
	}
	return nil
}

func (r *appointmentRepository) NoShow(id uint) error {
	result := r.db.Model(&models.Appointment{}).Where("id = ?", id).Update("status", "no_show")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("appointment not found")
	}
	return nil
}

func (r *appointmentRepository) SoftDelete(id uint) error {
	result := r.db.Delete(&models.Appointment{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("appointment not found")
	}

	return nil
}

func (r *appointmentRepository) Restore(id uint) error {
	result := r.db.Unscoped().Model(&models.Appointment{}).Where("id = ?", id).Update("deleted_at", nil)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("appointment not found")
	}

	return nil
}

func (r *appointmentRepository) HardDelete(id uint) error {
	result := r.db.Unscoped().Delete(&models.Appointment{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("appointment not found")
	}

	return nil
}
