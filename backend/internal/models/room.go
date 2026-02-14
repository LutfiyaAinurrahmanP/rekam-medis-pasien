package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	RoomNumber    string         `gorm:"unique;not null;size:20;index" json:"room_number"`
	RoomType      string         `gorm:"type:varchar(20);not null;index" json:"room_type"`
	DepartmentID  *uint          `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"department_id"`
	Department    *Department    `gorm:"foreignKey:DepartmentID;references:ID" json:"department,omitempty"`
	BedCapacity   int            `gorm:"not null;default:1" json:"bed_capacity"`
	AvailableBeds int            `gorm:"not null;default:1" json:"available_beds"`
	PricePerDay   float64        `gorm:"type:decimal(10,2);default:0" json:"price_per_day"`
	IsActive      bool           `gorm:"default:true;index" json:"is_active"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Room) TableName() string {
	return "rooms"
}

func (r *Room) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	if r.CreatedAt.IsZero() {
		r.CreatedAt = now
	}
	if r.UpdatedAt.IsZero() {
		r.UpdatedAt = now
	}
	// Set available_beds to bed_capacity if not set
	if r.AvailableBeds == 0 && r.BedCapacity > 0 {
		r.AvailableBeds = r.BedCapacity
	}
	return nil
}

// Helper methods
func (r *Room) OccupancyRate() float64 {
	if r.BedCapacity == 0 {
		return 0
	}
	occupied := r.BedCapacity - r.AvailableBeds
	return float64(occupied) / float64(r.BedCapacity) * 100
}

func (r *Room) IsAvailable() bool {
	return r.IsActive && r.AvailableBeds > 0
}

func (r *Room) IsFull() bool {
	return r.AvailableBeds == 0
}
