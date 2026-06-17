package models

import (
	"time"

	"gorm.io/gorm"
)

type Hospitalization struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	PatientID          uint           `gorm:"index;not null" json:"patient_id"`
	Patient            *Patient       `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	DoctorID           uint           `gorm:"index;not null" json:"doctor_id"`
	Doctor             *Doctor        `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	RoomID             uint           `gorm:"index;not null" json:"room_id"`
	Room               *Room          `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	AdmissionDate      string         `gorm:"type:date;not null;index" json:"admission_date"`
	AdmissionTime      string         `gorm:"type:time without time zone;not null" json:"admission_time"`
	DischargeDate      *string        `gorm:"type:date;index" json:"discharge_date"`
	DischargeTime      *string        `gorm:"type:time without time zone" json:"discharge_time"`
	ReasonForAdmission string         `gorm:"type:text;not null" json:"reason_for_admission"`
	Status             string         `gorm:"type:varchar(20);not null;default:'admitted';index" json:"status"`
	Notes              string         `gorm:"type:text" json:"notes"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Hospitalization) TableName() string {
	return "hospitalizations"
}

func (h *Hospitalization) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if h.CreatedAt.IsZero() {
		h.CreatedAt = now
	}
	if h.UpdatedAt.IsZero() {
		h.UpdatedAt = now
	}
	if h.Status == "" {
		h.Status = "admitted"
	}
	return nil
}
