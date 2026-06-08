package dto

import "time"

type AppointmentPatientResponse struct {
	ID          uint   `json:"id"`
	PatientCode string `json:"patient_code"`
	FullName    string `json:"full_name"`
	Phone       string `json:"phone,omitempty"`
	Age         int    `json:"age,omitempty"`
	Gender      string `json:"gender,omitempty"`
}

type AppointmentDoctorResponse struct {
	ID             uint   `json:"id"`
	FullName       string `json:"full_name"`
	Specialization string `json:"specialization,omitempty"`
	Department     string `json:"department,omitempty"`
}

type AppointmentResponse struct {
	ID              uint                        `json:"id"`
	PatientID       uint                        `json:"patient_id"`
	Patient         *AppointmentPatientResponse `json:"patient,omitempty"`
	DoctorID        uint                        `json:"doctor_id"`
	Doctor          *AppointmentDoctorResponse  `json:"doctor,omitempty"`
	AppointmentDate string                      `json:"appointment_date"`
	AppointmentTime string                      `json:"appointment_time"`
	DurationMinutes int                         `json:"duration_minutes"`
	Status          string                      `json:"status"`
	Reason          string                      `json:"reason"`
	Notes           string                      `json:"notes"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
}

type DeletedAppointmentResponse struct {
	ID              uint                        `json:"id"`
	PatientID       uint                        `json:"patient_id"`
	Patient         *AppointmentPatientResponse `json:"patient,omitempty"`
	DoctorID        uint                        `json:"doctor_id"`
	Doctor          *AppointmentDoctorResponse  `json:"doctor,omitempty"`
	AppointmentDate string                      `json:"appointment_date"`
	AppointmentTime string                      `json:"appointment_time"`
	DurationMinutes int                         `json:"duration_minutes"`
	Status          string                      `json:"status"`
	Reason          string                      `json:"reason"`
	Notes           string                      `json:"notes"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
	DeletedAt       *time.Time                  `json:"deleted_at"`
}

type AppointmentListResponse struct {
	Data []AppointmentResponse     `json:"data"`
	Meta AppointmentPaginationMeta `json:"meta"`
}

type AppointmentDeletedListResponse struct {
	Data []DeletedAppointmentResponse `json:"data"`
	Meta AppointmentPaginationMeta    `json:"meta"`
}

type CreateAppointmentRequest struct {
	PatientID       uint   `json:"patient_id" binding:"required"`
	DoctorID        uint   `json:"doctor_id" binding:"required"`
	AppointmentDate string `json:"appointment_date" binding:"required,datetime=2006-01-02"`
	AppointmentTime string `json:"appointment_time" binding:"required"`
	DurationMinutes int    `json:"duration_minutes" binding:"required,min=15"`
	Reason          string `json:"reason" binding:"omitempty,max=255"`
	Notes           string `json:"notes" binding:"omitempty"`
}

type UpdateAppointmentRequest struct {
	AppointmentDate *string `json:"appointment_date" binding:"omitempty,datetime=2006-01-02"`
	AppointmentTime *string `json:"appointment_time" binding:"omitempty"`
	DurationMinutes *int    `json:"duration_minutes" binding:"omitempty,min=15"`
	Reason          *string `json:"reason" binding:"omitempty,max=255"`
	Notes           *string `json:"notes" binding:"omitempty"`
}

type RescheduleAppointmentRequest struct {
	AppointmentDate string `json:"appointment_date" binding:"required,datetime=2006-01-02"`
	AppointmentTime string `json:"appointment_time" binding:"required"`
	Reason          string `json:"reason" binding:"omitempty,max=255"`
}

type CancelAppointmentRequest struct {
	Reason string `json:"reason" binding:"required,max=255"`
}

type AppointmentPaginationQuery struct {
	Page         int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize     int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	PatientID    *uint  `form:"patient_id" binding:"omitempty"`
	DoctorID     *uint  `form:"doctor_id" binding:"omitempty"`
	DepartmentID *uint  `form:"department_id" binding:"omitempty"`
	Status       string `form:"status" binding:"omitempty,oneof=scheduled confirmed in_progress completed cancelled no_show"`
	Date         string `form:"date" binding:"omitempty,datetime=2006-01-02"`
	DateFrom     string `form:"date_from" binding:"omitempty,datetime=2006-01-02"`
	DateTo       string `form:"date_to" binding:"omitempty,datetime=2006-01-02"`
	DaysAhead    int    `form:"days_ahead" binding:"omitempty,min=1"`
	DaysBack     int    `form:"days_back" binding:"omitempty,min=1"`
	SortBy       string `form:"sort_by,default=appointment_date" binding:"omitempty,oneof=created_at appointment_date status duration_minutes"`
	SortDir      string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type AppointmentPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
