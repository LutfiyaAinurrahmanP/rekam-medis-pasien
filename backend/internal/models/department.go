package models

import (
	"time"

	"gorm.io/gorm"
)

type Department struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"not null;size:100;index" json:"name"`
	Code          string         `gorm:"unique;not null;size:20;index" json:"code"`
	Description   string         `gorm:"null" json:"description"`
	FloorLocation string         `gorm:"null;size:50" json:"floor_location"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Department) TableName() string {
	return "departments"
}

func (d *Department) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	if d.UpdatedAt.IsZero() {
		d.UpdatedAt = now
	}
	return nil
}
