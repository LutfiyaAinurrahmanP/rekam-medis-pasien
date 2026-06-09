package models

import (
	"time"

	"gorm.io/gorm"
)

type LabTest struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	MedicalRecordID      uint           `gorm:"index;not null" json:"medical_record_id"`
	MedicalRecord        *MedicalRecord `gorm:"foreignKey:MedicalRecordID" json:"medical_record,omitempty"`
	TestTypeID           uint           `gorm:"index;not null" json:"test_type_id"`
	TestType             *TypeTest      `gorm:"foreignKey:TestTypeID" json:"test_type,omitempty"`
	OrderedByDoctorID    uint           `gorm:"index;not null" json:"ordered_by_doctor_id"`
	OrderedByDoctor      *Doctor        `gorm:"foreignKey:OrderedByDoctorID" json:"ordered_by_doctor,omitempty"`
	OrderDate            string         `gorm:"type:date;not null;index" json:"order_date"`
	SampleCollectionDate *string        `gorm:"type:date;index" json:"sample_collection_date,omitempty"`
	TestStartDate        *string        `gorm:"type:date;index" json:"test_start_date,omitempty"`
	ResultDate           *string        `gorm:"type:date" json:"result_date,omitempty"`
	ResultValue          *string        `gorm:"type:text" json:"result_value,omitempty"`
	ResultUnit           *string        `gorm:"type:varchar(50)" json:"result_unit,omitempty"`
	ReferenceRange       *string        `gorm:"type:varchar(200)" json:"reference_range,omitempty"`
	Status               string         `gorm:"type:varchar(20);not null;default:'ordered';index" json:"status"`
	Notes                string         `gorm:"type:text" json:"notes"`
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (LabTest) TableName() string {
	return "lab_tests"
}

func (m *LabTest) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}
	return nil
}
