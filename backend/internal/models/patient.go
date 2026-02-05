package models

import (
	"time"

	"gorm.io/gorm"
)

type Patient struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	UserID                *uint          `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user_id"`
	User                  User           `gorm:"foreignKey:UserID;references:ID" json:"-"`
	PatientCode           string         `gorm:"unique;not null;size:20;index" json:"patient_code"`
	FullName              string         `gorm:"not null;size:100" json:"full_name"`
	DateOfBirth           string         `gorm:"type:date;not null" json:"date_of_birth"`
	Gender                string         `gorm:"type:varchar(10);not null" json:"gender"`
	BloodType             string         `gorm:"type:varchar(5);not null" json:"blood_type"`
	Phone                 string         `gorm:"size:15" json:"phone"`
	Email                 string         `gorm:"size:100" json:"email"`
	Address               string         `gorm:"type:text" json:"address"`
	EmergencyContactName  string         `gorm:"size:100" json:"emergency_contact_name"`
	EmergencyContactPhone string         `gorm:"size:15" json:"emergency_contact_phone"`
	InsuranceNumber       string         `gorm:"size:50" json:"insurance_number"`
	InsuranceProvider     string         `gorm:"size:100" json:"insurance_provider"`
	Allergies             string         `gorm:"type:text" json:"allergies"`
	CreatedAt             time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Patient) TableName() string {
	return "patients"
}

func (p *Patient) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}
	return nil
}
