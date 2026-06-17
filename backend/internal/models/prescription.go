package models

import (
	"time"

	"gorm.io/gorm"
)

type Prescription struct {
	ID               uint               `gorm:"primaryKey" json:"id"`
	MedicalRecordID  uint               `gorm:"index;not null" json:"medical_record_id"`
	MedicalRecord    *MedicalRecord     `gorm:"foreignKey:MedicalRecordID" json:"medical_record,omitempty"`
	DoctorID         uint               `gorm:"index;not null" json:"doctor_id"`
	Doctor           *Doctor            `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	PrescriptionDate string             `gorm:"type:date;not null;index" json:"prescription_date"`
	Notes            string             `gorm:"type:text" json:"notes"`
	Status           string             `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	CreatedAt        time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt     `gorm:"index" json:"deleted_at,omitempty"`
	Items            []PrescriptionItem `gorm:"foreignKey:PrescriptionID" json:"items,omitempty"`
}

func (Prescription) TableName() string {
	return "prescriptions"
}

func (p *Prescription) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}
	return nil
}

type PrescriptionItem struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	PrescriptionID uint           `gorm:"index;not null" json:"prescription_id"`
	Prescription   *Prescription  `gorm:"foreignKey:PrescriptionID" json:"-"`
	MedicineID     uint           `gorm:"index;not null" json:"medicine_id"`
	Medicine       *Medicine      `gorm:"foreignKey:MedicineID" json:"medicine,omitempty"`
	Dosage         string         `gorm:"type:varchar(100);not null" json:"dosage"`
	Frequency      string         `gorm:"type:varchar(100);not null" json:"frequency"`
	DurationDays   int            `gorm:"not null" json:"duration_days"`
	Quantity       int            `gorm:"not null" json:"quantity"`
	Instructions   string         `gorm:"type:text" json:"instructions"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (PrescriptionItem) TableName() string {
	return "prescription_items"
}

func (p *PrescriptionItem) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}
	return nil
}
