package models

import (
	"time"

	"gorm.io/gorm"
)

type VitalSign struct {
	ID               uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	MedicalRecordID  uint       `gorm:"not null;index" json:"medical_record_id"`
	MeasurementDate  time.Time  `gorm:"type:date;not null;index" json:"measurement_date"`
	MeasurementTime  string     `gorm:"type:time without time zone;not null" json:"measurement_time"`
	SystolicBP       *int       `gorm:"type:int" json:"systolic_bp"`
	DiastolicBP      *int       `gorm:"type:int" json:"diastolic_bp"`
	HeartRate        *int       `gorm:"type:int" json:"heart_rate"`
	BodyTemperature  *float64   `gorm:"type:decimal(4,2)" json:"body_temperature"`
	RespiratoryRate  *int       `gorm:"type:int" json:"respiratory_rate"`
	OxygenSaturation *float64   `gorm:"type:decimal(5,2)" json:"oxygen_saturation"`
	WeightKg         *float64   `gorm:"type:decimal(5,2)" json:"weight_kg"`
	HeightCm         *int       `gorm:"type:int" json:"height_cm"`
	BMI              *float64   `gorm:"type:decimal(5,2)" json:"bmi"`
	Notes            string     `gorm:"type:text" json:"notes"`
	CreatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	MedicalRecord MedicalRecord `gorm:"foreignKey:MedicalRecordID;references:ID" json:"-"`
}

func (VitalSign) TableName() string {
	return "vital_signs"
}
