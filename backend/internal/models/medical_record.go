package models

import (
	"time"

	"gorm.io/gorm"
)

type MedicalRecord struct {
	ID                       uint           `gorm:"primaryKey" json:"id"`
	PatientID                uint           `gorm:"index;not null" json:"patient_id"`
	Patient                  *Patient       `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	DoctorID                 uint           `gorm:"index;not null" json:"doctor_id"`
	Doctor                   *Doctor        `gorm:"foreignKey:DoctorID" json:"doctor,omitempty"`
	AppointmentID            *uint          `gorm:"index" json:"appointment_id,omitempty"`
	Appointment              *Appointment   `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`
	VisitDate                string         `gorm:"type:date;not null;index" json:"visit_date"`
	ChiefComplaint           string         `gorm:"type:text;not null" json:"chief_complaint"`
	HistoryOfIllness         string         `gorm:"type:text" json:"history_of_illness"`
	PhysicalExamination      string         `gorm:"type:text" json:"physical_examination"`
	Diagnosis                string         `gorm:"type:text;not null" json:"diagnosis"`
	TreatmentPlan            string         `gorm:"type:text;not null" json:"treatment_plan"`
	Notes                    string         `gorm:"type:text" json:"notes"`
	Status                   string         `gorm:"type:varchar(20);not null;default:'draft';index" json:"status"`
	CreatedAt                time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt                gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	VitalSign                *VitalSign     `gorm:"foreignKey:MedicalRecordID" json:"vital_sign,omitempty"`
}

func (MedicalRecord) TableName() string {
	return "medical_records"
}

func (m *MedicalRecord) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}
	return nil
}
