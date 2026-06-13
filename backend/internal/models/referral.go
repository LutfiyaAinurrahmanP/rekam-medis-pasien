package models

import (
	"time"

	"gorm.io/gorm"
)

type Referral struct {
	ID                     uint           `gorm:"primaryKey" json:"id"`
	ReferralNumber         string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"referral_number"`
	PatientID              uint           `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"patient_id"`
	Patient                *Patient       `gorm:"foreignKey:PatientID;references:ID" json:"-"`
	MedicalRecordID        uint           `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"medical_record_id"`
	MedicalRecord          *MedicalRecord `gorm:"foreignKey:MedicalRecordID;references:ID" json:"-"`
	ReferringDoctorID      uint           `gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"referring_doctor_id"`
	ReferringDoctor        *Doctor        `gorm:"foreignKey:ReferringDoctorID;references:ID" json:"-"`
	ReferralDate           string         `gorm:"type:date;not null;index" json:"referral_date"`
	ReferralType           string         `gorm:"type:varchar(20);not null;index" json:"referral_type"` // internal, external
	ReferredToDepartmentID *uint          `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"referred_to_department_id"`
	ReferredToDepartment   *Department    `gorm:"foreignKey:ReferredToDepartmentID;references:ID" json:"-"`
	ReferredToDoctorID     *uint          `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"referred_to_doctor_id"`
	ReferredToDoctor       *Doctor        `gorm:"foreignKey:ReferredToDoctorID;references:ID" json:"-"`
	ReferredToFacility     string         `gorm:"type:varchar(255)" json:"referred_to_facility"`
	Reason                 string         `gorm:"type:text;not null" json:"reason"`
	Diagnosis              string         `gorm:"type:text" json:"diagnosis"`
	Priority               string         `gorm:"type:varchar(20);not null;default:'routine';index" json:"priority"` // routine, urgent, emergency
	Status                 string         `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`   // pending, accepted, rejected, completed, cancelled
	AcceptedAt             *time.Time     `json:"accepted_at"`
	CompletedAt            *time.Time     `json:"completed_at"`
	RejectionReason        string         `gorm:"type:text" json:"rejection_reason"`
	CancellationReason     string         `gorm:"type:text" json:"cancellation_reason"`
	Notes                  string         `gorm:"type:text" json:"notes"`
	CreatedAt              time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Referral) TableName() string {
	return "referrals"
}

// BeforeCreate is a GORM hook that can be used to set default values
func (r *Referral) BeforeCreate(tx *gorm.DB) (err error) {
	if r.Status == "" {
		r.Status = "pending"
	}
	if r.Priority == "" {
		r.Priority = "routine"
	}
	return
}
