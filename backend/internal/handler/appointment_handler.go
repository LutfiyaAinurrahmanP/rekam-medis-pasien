package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/appointment"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	service     appointment.AppointmentService
	doctorRepo  repository.DoctorRepository
	patientRepo repository.PatientRepository
}

func NewAppointmentHandler(s appointment.AppointmentService, dr repository.DoctorRepository, pr repository.PatientRepository) *AppointmentHandler {
	return &AppointmentHandler{
		service:     s,
		doctorRepo:  dr,
		patientRepo: pr,
	}
}

// Helper to inject DoctorID or PatientID based on logged-in user
func (h *AppointmentHandler) injectOwnership(ctx *gin.Context, query *dto.AppointmentPaginationQuery) error {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil // if no user id, we just pass
	}

	role, _ := ctx.Get("role")
	uid := userID.(uint)

	if role == "doctor" {
		doctor, err := h.doctorRepo.FindByUserID(uid)
		if err == nil && doctor != nil {
			query.DoctorID = &doctor.ID
		}
	} else if role == "patient" {
		patient, err := h.patientRepo.FindByUserID(uid)
		if err == nil && patient != nil {
			query.PatientID = &patient.ID
		}
	}
	return nil
}

func (h *AppointmentHandler) MyAppointments(ctx *gin.Context) {
	var query dto.AppointmentPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	h.injectOwnership(ctx, &query)

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve my appointments", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "My appointments retrieved successfully", res)
}

func (h *AppointmentHandler) MySchedule(ctx *gin.Context) {
	var query dto.AppointmentPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	h.injectOwnership(ctx, &query)

	// If a specific date is given it acts like List, if not we can just default to TodayList or List
	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve my schedule", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "My schedule retrieved successfully", res)
}

func (h *AppointmentHandler) List(ctx *gin.Context) {
	var query dto.AppointmentPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve appointments", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointments retrieved successfully", res)
}

func (h *AppointmentHandler) TodayList(ctx *gin.Context) {
	var query dto.AppointmentPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.TodayList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve today's appointments", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Today's appointments retrieved successfully", res)
}

func (h *AppointmentHandler) UpcomingList(ctx *gin.Context) {
	var query dto.AppointmentPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	// For regular users, limit to their own upcoming
	role, _ := ctx.Get("role")
	if role == "patient" || role == "doctor" {
		h.injectOwnership(ctx, &query)
	}

	res, err := h.service.UpcomingList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve upcoming appointments", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Upcoming appointments retrieved successfully", res)
}

func (h *AppointmentHandler) PastList(ctx *gin.Context) {
	var query dto.AppointmentPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	// For regular users, limit to their own past
	role, _ := ctx.Get("role")
	if role == "patient" || role == "doctor" {
		h.injectOwnership(ctx, &query)
	}

	res, err := h.service.PastList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve past appointments", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Past appointments retrieved successfully", res)
}

func (h *AppointmentHandler) CancelledList(ctx *gin.Context) {
	var query dto.AppointmentPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	
	query.Status = "cancelled"

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve cancelled appointments", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Cancelled appointments retrieved successfully", res)
}

func (h *AppointmentHandler) DeletedList(ctx *gin.Context) {
	var query dto.AppointmentPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve deleted appointments", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Deleted appointments retrieved successfully", res)
}

func (h *AppointmentHandler) FindByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Appointment not found", err.Error())
		return
	}
	
	// Ownership check
	role, _ := ctx.Get("role")
	userID, _ := ctx.Get("user_id")
	uid := userID.(uint)
	
	if role == "doctor" {
		doctor, err := h.doctorRepo.FindByUserID(uid)
		if err == nil && doctor != nil && res.DoctorID != doctor.ID {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Forbidden", "You do not have access to this appointment")
			return
		}
	} else if role == "patient" {
		patient, err := h.patientRepo.FindByUserID(uid)
		if err == nil && patient != nil && res.PatientID != patient.ID {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Forbidden", "You do not have access to this appointment")
			return
		}
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Appointment retrieved successfully", res)
}

func (h *AppointmentHandler) Create(ctx *gin.Context) {
	var req dto.CreateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	// If patient, force patient_id to self
	role, _ := ctx.Get("role")
	userID, _ := ctx.Get("user_id")
	if role == "patient" {
		patient, err := h.patientRepo.FindByUserID(userID.(uint))
		if err == nil && patient != nil {
			req.PatientID = patient.ID
		}
	}

	res, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create appointment", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "Appointment created successfully", res)
}

func (h *AppointmentHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update appointment", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointment updated successfully", res)
}

func (h *AppointmentHandler) Confirm(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Confirm(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to confirm appointment", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointment confirmed successfully", nil)
}

func (h *AppointmentHandler) Start(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Start(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to start appointment", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointment started successfully", nil)
}

func (h *AppointmentHandler) Complete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Complete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to complete appointment", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointment completed successfully", nil)
}

func (h *AppointmentHandler) Cancel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.CancelAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if err := h.service.Cancel(uint(id), &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to cancel appointment", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointment cancelled successfully", nil)
}

func (h *AppointmentHandler) Reschedule(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.RescheduleAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if err := h.service.Reschedule(uint(id), &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to reschedule appointment", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointment rescheduled successfully", nil)
}

func (h *AppointmentHandler) NoShow(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.NoShow(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to mark appointment as no-show", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointment marked as no-show successfully", nil)
}

func (h *AppointmentHandler) SoftDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete appointment", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointment deleted successfully", nil)
}

func (h *AppointmentHandler) Restore(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to restore appointment", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointment restored successfully", nil)
}

func (h *AppointmentHandler) HardDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to hard delete appointment", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Appointment permanently deleted", nil)
}
