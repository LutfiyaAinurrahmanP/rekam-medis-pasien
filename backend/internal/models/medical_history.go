package models

import (
	"time"
)

type Allergy struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PatientID    uint      `gorm:"not null;index" json:"patient_id"`
	AllergenType string    `gorm:"type:varchar(50);not null" json:"allergen_type"`
	AllergenName string    `gorm:"type:varchar(100);not null" json:"allergen_name"`
	Reaction     string    `gorm:"type:text;not null" json:"reaction"`
	Severity     string    `gorm:"type:varchar(20);not null" json:"severity"`
	Notes        string    `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Patient Patient `gorm:"foreignKey:PatientID;references:ID" json:"-"`
}

func (Allergy) TableName() string {
	return "allergies"
}

type MedicalCondition struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	PatientID     uint       `gorm:"not null;index" json:"patient_id"`
	ConditionName string     `gorm:"type:varchar(200);not null;index" json:"condition_name"`
	ICDCode       string     `gorm:"type:varchar(20)" json:"icd_code"`
	DiagnosedDate *time.Time `gorm:"type:date" json:"diagnosed_date"`
	Status        string     `gorm:"type:varchar(20);default:'ongoing';index" json:"status"`
	Notes         string     `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Patient Patient `gorm:"foreignKey:PatientID;references:ID" json:"-"`
}

func (MedicalCondition) TableName() string {
	return "medical_conditions"
}

type SurgicalHistory struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PatientID     uint      `gorm:"not null;index" json:"patient_id"`
	ProcedureName string    `gorm:"type:varchar(200);not null" json:"procedure_name"`
	SurgeryDate   time.Time `gorm:"type:date;not null;index" json:"surgery_date"`
	SurgeonName   string    `gorm:"type:varchar(100)" json:"surgeon_name"`
	Hospital      string    `gorm:"type:varchar(200)" json:"hospital"`
	Complication  string    `gorm:"type:text" json:"complication"`
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Patient Patient `gorm:"foreignKey:PatientID;references:ID" json:"-"`
}

func (SurgicalHistory) TableName() string {
	return "surgical_histories"
}

type FamilyHistory struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	PatientID     uint      `gorm:"not null;index" json:"patient_id"`
	FamilyMember  string    `gorm:"type:varchar(50);not null;index" json:"family_member"`
	ConditionName string    `gorm:"type:varchar(200);not null" json:"condition_name"`
	Relation      string    `gorm:"type:varchar(100)" json:"relation"`
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Patient Patient `gorm:"foreignKey:PatientID;references:ID" json:"-"`
}

func (FamilyHistory) TableName() string {
	return "family_histories"
}