package dto

import "time"

type DoctorListResponse struct {
	Data []DoctorResponse     `json:"data"`
	Meta DoctorPaginationMeta `json:"meta"`
}

type DoctorDeletedListResponse struct {
	Data []DeletedDoctorResponse `json:"data"`
	Meta DoctorPaginationMeta    `json:"meta"`
}

type CreateDoctorRequest struct {
	UserID           *uint  `json:"user_id" binding:"omitempty"`
	EmployeeID       string `json:"employee_id" binding:"required,min=1,max=50"`
	FullName         string `json:"full_name" binding:"required,min=1,max=100"`
	SpecializationID *uint  `json:"doctor_specialization_id" binding:"required"`
	LicenseNumber    string `json:"license_number" binding:"required,min=1,max=50"`
	Phone            string `json:"phone" binding:"omitempty,max=15"`
	Email            string `json:"email" binding:"omitempty,max=100"`
	DepartmentID     *uint  `json:"department_id" binding:"required"`
	IsActive         *bool  `json:"is_active" binding:"omitempty"`
}

type UpdateDoctorRequest struct {
	UserID           *uint   `json:"user_id" binding:"omitempty"`
	EmployeeID       *string `json:"employee_id" binding:"omitempty,min=1,max=50"`
	FullName         *string `json:"full_name" binding:"omitempty,min=1,max=100"`
	SpecializationID *uint   `json:"doctor_specialization_id" binding:"omitempty"`
	LicenseNumber    *string `json:"license_number" binding:"omitempty,min=1,max=50"`
	Phone            *string `json:"phone" binding:"omitempty,max=15"`
	Email            *string `json:"email" binding:"omitempty,max=100"`
	DepartmentID     *uint   `json:"department_id" binding:"omitempty"`
	IsActive         *bool   `json:"is_active" binding:"omitempty"`
}

type DoctorResponse struct {
	ID               uint      `json:"id"`
	UserID           *uint     `json:"user_id"`
	EmployeeID       string    `json:"employee_id"`
	FullName         string    `json:"full_name"`
	SpecializationID uint      `json:"doctor_specialization_id"`
	LicenseNumber    string    `json:"license_number"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email"`
	DepartmentID     *uint     `json:"department_id"`
	IsActive         *bool     `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type DeletedDoctorResponse struct {
	ID               uint       `json:"id"`
	UserID           *uint      `json:"user_id"`
	EmployeeID       string     `json:"employee_id"`
	FullName         string     `json:"full_name"`
	SpecializationID uint       `json:"doctor_specialization_id"`
	LicenseNumber    string     `json:"license_number"`
	Phone            string     `json:"phone"`
	Email            string     `json:"email"`
	DepartmentID     *uint      `json:"department_id"`
	IsActive         *bool      `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at"`
}

type DoctorPaginationQuery struct {
	Page             int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize         int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	Search           string `form:"search" binding:"omitempty"`
	SpecializationID *uint  `form:"doctor_specialization_id" binding:"omitempty"`
	IsActive         *bool  `form:"is_active" binding:"omitempty"`
	SortBy           string `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at full_name employee_id license_number specialization_id"`
	SortDir          string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type DoctorPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
