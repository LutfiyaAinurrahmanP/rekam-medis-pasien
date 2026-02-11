package models

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	UserID         *uint          `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_id"`
	User           User           `gorm:"foreignKey:UserID;references:ID" json:"-"`
	EmployeeID     string         `gorm:"unique;not null;size:50;index" json:"employee_id"`
	FullName       string         `gorm:"not null;size:100" json:"full_name"`
	Specialization string         `gorm:"not null;size:100" json:"specialization"`
	LicenseNumber  string         `gorm:"unique;not null;size:50" json:"license_number"`
	Phone          string         `gorm:"size:15" json:"phone"`
	Email          string         `gorm:"size:100" json:"email"`
	DepartmentID   *uint          `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"department_id"`
	Department     Department     `gorm:"foreignKey:DepartmentID;references:ID" json:"-"`
	IsActive       bool           `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Doctor) TableName() string {
	return "doctors"
}

func (p *Doctor) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}
	return nil
}
