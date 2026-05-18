package models

import (
	"time"

	"gorm.io/gorm"
)

type DoctorSpecialization struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;unique;size:100;index" json:"name"`
	Code        string         `gorm:"size:20;unique;index" json:"code,omitempty"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool           `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (DoctorSpecialization) TableName() string {
	return "doctor_specializations"
}

func (d *DoctorSpecialization) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	if d.UpdatedAt.IsZero() {
		d.UpdatedAt = now
	}
	return nil
}
