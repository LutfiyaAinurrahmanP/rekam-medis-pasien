package models

import (
	"time"

	"gorm.io/gorm"
)

type TypeTestCategory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;unique;size:100;index" json:"name"`
	Code        string         `gorm:"unique;size:20;index" json:"code"`
	Description string         `gorm:"type:text" json:"description"`
	IsActive    bool           `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (TypeTestCategory) TableName() string {
	return "type_test_categories"
}

func (t *TypeTestCategory) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	if t.UpdatedAt.IsZero() {
		t.UpdatedAt = now
	}
	return nil
}
