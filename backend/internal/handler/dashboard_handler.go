package handler

import (
	"net/http"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	dashboard "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/dashboard"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service     dashboard.DashboardService
	doctorRepo  repository.DoctorRepository
	patientRepo repository.PatientRepository
}

func NewDashboardHandler(svc dashboard.DashboardService, dr repository.DoctorRepository, pr repository.PatientRepository) *DashboardHandler {
	return &DashboardHandler{service: svc, doctorRepo: dr, patientRepo: pr}
}

func (h *DashboardHandler) Overview(c *gin.Context) {
	var query dto.DashboardOverviewQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.GetOverview(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve dashboard overview", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard overview retrieved successfully", result)
}

func (h *DashboardHandler) AdminDashboard(c *gin.Context) {
	var query dto.DashboardPeriodQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.GetAdminDashboard(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve admin dashboard", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Admin dashboard retrieved successfully", result)
}

func (h *DashboardHandler) DoctorDashboard(c *gin.Context) {
	var query dto.DashboardDoctorQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	doctor, err := h.doctorRepo.FindByUserID(uid)
	if err != nil || doctor == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Doctor profile not found", "")
		return
	}

	result, svcErr := h.service.GetDoctorDashboard(doctor.ID, &query)
	if svcErr != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve doctor dashboard", svcErr.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Doctor dashboard retrieved successfully", result)
}

func (h *DashboardHandler) ReceptionistDashboard(c *gin.Context) {
	result, err := h.service.GetReceptionistDashboard()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve receptionist dashboard", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Receptionist dashboard retrieved successfully", result)
}

func (h *DashboardHandler) PatientDashboard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	patient, err := h.patientRepo.FindByUserID(uid)
	if err != nil || patient == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Patient profile not found", "")
		return
	}

	result, svcErr := h.service.GetPatientDashboard(patient.ID)
	if svcErr != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve patient dashboard", svcErr.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient dashboard retrieved successfully", result)
}

func (h *DashboardHandler) AppointmentReport(c *gin.Context) {
	var query dto.DashboardAppointmentReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Doctor can only see their own report
	role, _ := c.Get("role")
	if role == "doctor" {
		userID, _ := c.Get("user_id")
		uid := userID.(uint)
		doctor, err := h.doctorRepo.FindByUserID(uid)
		if err == nil && doctor != nil {
			query.DoctorID = &doctor.ID
		}
	}

	result, err := h.service.GetAppointmentReport(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve appointment report", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Appointment report retrieved successfully", result)
}

func (h *DashboardHandler) RevenueReport(c *gin.Context) {
	var query dto.DashboardRevenueReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.GetRevenueReport(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve revenue report", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Revenue report retrieved successfully", result)
}

func (h *DashboardHandler) PatientReport(c *gin.Context) {
	var query dto.DashboardPatientReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.GetPatientReport(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve patient report", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient report retrieved successfully", result)
}
