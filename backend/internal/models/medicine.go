package models

import (
	"time"

	"gorm.io/gorm"
)

type Medicine struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"not null;size:200;index" json:"name"`
	GenericName    string         `gorm:"size:200" json:"generic_name"`
	BrandName      string         `gorm:"size:200" json:"brand_name"`
	MedicineTypeID uint           `gorm:"not null" json:"medicine_type_id"`
	Strength       string         `gorm:"size:50" json:"strength"`
	Manufacturer   string         `gorm:"size:100" json:"manufacturer"`
	Unit           string         `gorm:"size:20" json:"unit"`
	StockQuantity  int            `gorm:"not null;default:0" json:"stock_quantity"`
	Price          float64        `gorm:"type:decimal(10,2)" json:"price"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Medicine) TableName() string {
	return "medicines"
}

func (m *Medicine) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	if m.UpdatedAt.IsZero() {
		m.UpdatedAt = now
	}
	return nil
}
