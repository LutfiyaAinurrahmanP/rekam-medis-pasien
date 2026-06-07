package models

import (
	"time"

	"gorm.io/gorm"
)

type TypeTest struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"not null;size:200;index" json:"name"`
	Code               string         `gorm:"not null;unique;size:50;index" json:"code"`
	TypeTestCategoryID uint           `gorm:"not null;index" json:"type_test_category_id"`
	Description        string         `gorm:"type:text" json:"description"`
	Price              float64        `gorm:"type:decimal(10,2)" json:"price"`
	IsActive           bool           `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (TypeTest) TableName() string {
	return "type_tests"
}

func (t *TypeTest) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = now
	}
	return nil
}
