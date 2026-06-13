package models

import (
	"time"

	"gorm.io/gorm"
)

type Billing struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	PatientID         uint            `gorm:"index;not null" json:"patient_id"`
	Patient           *Patient        `gorm:"foreignKey:PatientID" json:"patient,omitempty"`
	MedicalRecordID   *uint           `gorm:"index" json:"medical_record_id,omitempty"`
	MedicalRecord     *MedicalRecord  `gorm:"foreignKey:MedicalRecordID" json:"medical_record,omitempty"`
	HospitalizationID *uint           `gorm:"index" json:"hospitalization_id,omitempty"`
	Hospitalization   *Hospitalization`gorm:"foreignKey:HospitalizationID" json:"hospitalization,omitempty"`
	InvoiceNumber     string          `gorm:"type:varchar(50);uniqueIndex;not null" json:"invoice_number"`
	BillingDate       string          `gorm:"type:date;not null;index" json:"billing_date"`
	DueDate           string          `gorm:"type:date;not null;index" json:"due_date"`
	TotalAmount       float64         `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	PaidAmount        float64         `gorm:"type:decimal(12,2);default:0" json:"paid_amount"`
	RemainingAmount   float64         `gorm:"type:decimal(12,2);default:0" json:"remaining_amount"`
	Status            string          `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	PaymentMethod     *string         `gorm:"type:varchar(50)" json:"payment_method,omitempty"`
	Notes             *string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`
	
	Items             []BillingItem   `gorm:"foreignKey:BillingID" json:"items,omitempty"`
}

func (Billing) TableName() string {
	return "billing"
}

func (b *Billing) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = now
	}
	return nil
}

type BillingItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	BillingID   uint           `gorm:"index;not null" json:"billing_id"`
	Billing     *Billing       `gorm:"foreignKey:BillingID" json:"billing,omitempty"`
	Description string         `gorm:"type:varchar(200);not null" json:"description"`
	Quantity    int            `gorm:"not null" json:"quantity"`
	UnitPrice   float64        `gorm:"type:decimal(12,2);not null" json:"unit_price"`
	TotalPrice  float64        `gorm:"type:decimal(12,2);not null" json:"total_price"`
	ItemType    *string        `gorm:"type:varchar(50)" json:"item_type,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (BillingItem) TableName() string {
	return "billing_items"
}

func (bi *BillingItem) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if bi.CreatedAt.IsZero() {
		bi.CreatedAt = now
	}
	if bi.UpdatedAt.IsZero() {
		bi.UpdatedAt = now
	}
	return nil
}
