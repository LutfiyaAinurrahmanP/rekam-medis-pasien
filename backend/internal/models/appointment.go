package models

import (
	"time"

	"gorm.io/gorm"
)

type Appointment struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	PatientID       uint           `gorm:"index;not null" json:"patient_id"`
	Patient         *Patient       `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	DoctorID        uint           `gorm:"index;not null" json:"doctor_id"`
	Doctor          *Doctor        `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	AppointmentDate string         `gorm:"type:date;not null;index" json:"appointment_date"`
	AppointmentTime string         `gorm:"type:time without time zone;not null" json:"appointment_time"`
	DurationMinutes int            `gorm:"default:30" json:"duration_minutes"`
	Status          string         `gorm:"type:varchar(20);not null;default:'scheduled';index" json:"status"`
	Reason          string         `gorm:"type:varchar(255)" json:"reason"`
	Notes           string         `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Appointment) TableName() string {
	return "appointments"
}

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = now
	}
	return nil
}
