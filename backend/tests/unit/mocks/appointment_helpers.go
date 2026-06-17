package mocks

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func NewTestAppointmentWithData(id uint, patientID uint, doctorID uint, date, timeStr, status, reason, notes string, duration int, isDeleted bool) *models.Appointment {
	now := time.Now()
	a := &models.Appointment{
		ID:              id,
		PatientID:       patientID,
		DoctorID:        doctorID,
		AppointmentDate: date,
		AppointmentTime: timeStr,
		DurationMinutes: duration,
		Status:          status,
		Reason:          reason,
		Notes:           notes,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if isDeleted {
		a.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	}

	return a
}

func NewTestAppointmentList(count int) []models.Appointment {
	var list []models.Appointment
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestAppointmentWithData(
			uint(i),
			uint(1),
			uint(1),
			"2023-12-01",
			"10:00",
			"scheduled",
			fmt.Sprintf("Reason %d", i),
			fmt.Sprintf("Notes %d", i),
			30,
			false,
		))
	}
	return list
}

func NewTestAppointmentResponse(a *models.Appointment) *dto.AppointmentResponse {
	return &dto.AppointmentResponse{
		ID:              a.ID,
		PatientID:       a.PatientID,
		DoctorID:        a.DoctorID,
		AppointmentDate: a.AppointmentDate,
		AppointmentTime: a.AppointmentTime,
		DurationMinutes: a.DurationMinutes,
		Status:          a.Status,
		Reason:          a.Reason,
		Notes:           a.Notes,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
	}
}

func NewAppointmentPaginationQuery(page, pageSize int) *dto.AppointmentPaginationQuery {
	return &dto.AppointmentPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "appointment_date",
		SortDir:  "desc",
	}
}

func NewCreateAppointmentRequest(patientID, doctorID uint, date, timeStr string, duration int, reason, notes string) *dto.CreateAppointmentRequest {
	return &dto.CreateAppointmentRequest{
		PatientID:       patientID,
		DoctorID:        doctorID,
		AppointmentDate: date,
		AppointmentTime: timeStr,
		DurationMinutes: duration,
		Reason:          reason,
		Notes:           notes,
	}
}

func NewUpdateAppointmentRequest(date, timeStr string, duration int, reason, notes string) *dto.UpdateAppointmentRequest {
	return &dto.UpdateAppointmentRequest{
		AppointmentDate: &date,
		AppointmentTime: &timeStr,
		DurationMinutes: &duration,
		Reason:          &reason,
		Notes:           &notes,
	}
}
