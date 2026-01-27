package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	RolePatient      = "patient"
	RoleDoctor       = "doctor"
	RoleReceptionist = "receptionist"
	RoleAdmin        = "admin"
	RoleSuperAdmin   = "super_admin"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"unique;not null;size:50;index" json:"username"`
	Email     string         `gorm:"unique;not null;size:100;index" json:"email"`
	Phone     string         `gorm:"unique;not null;size:15;index" json:"phone"`
	Password  string         `gorm:"not null;size:255" json:"password"`
	Role      string         `gorm:"type:varchar(20);not null;default:'patient';index" json:"role"`
	IsActive  bool           `gorm:"not null;default:true;index" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = RolePatient
	}

	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}
	return nil
}

func (u *User) IsPatient() bool {
	return u.Role == RolePatient
}

func (u *User) IsDoctor() bool {
	return u.Role == RoleDoctor
}

func (u *User) IsReceptionist() bool {
	return u.Role == RoleReceptionist
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

func ValidateRole(role string) bool {
	validRoles := GetAvailableRoles()
	for _, r := range validRoles {
		if r == role {
			return true
		}
	}
	return false
}

func GetAvailableRoles() []string {
	return []string{
		RolePatient,
		RoleDoctor,
		RoleReceptionist,
		RoleAdmin,
		RoleSuperAdmin,
	}
}
