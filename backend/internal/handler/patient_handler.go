package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	patientservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/patient"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	service patientservice.PatientService
}

func NewPatientHandler(service patientservice.PatientService) *PatientHandler {
	return &PatientHandler{
		service: service,
	}
}

func (h *PatientHandler) GetMyPatientData(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	patient, err := h.service.GetMyPatientData(userID.(uint))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Patient data not found", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Patient data retrieved successfully", patient)
}

func (h *PatientHandler) UpdateMyPatientData(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	var req dto.UpdatePatientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	patient, err := h.service.UpdateMyPatientData(userID.(uint), &req)
	if err != nil {
		if err.Error() == "patient not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Patient data not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update patient data", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Patient data updated successfully", patient)
}

func (h *PatientHandler) ListPatients(ctx *gin.Context) {
	var query dto.PatientPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	patients, err := h.service.ListPatients(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve patients", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Patients retrieved successfully", patients)
}

func (h *PatientHandler) SearchPatients(ctx *gin.Context) {
	var query dto.PatientPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	patients, err := h.service.ListPatients(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to search patients", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Patients search results retrieved successfully", patients)
}

func (h *PatientHandler) DeleteListPatients(ctx *gin.Context) {
	var query dto.PatientPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	patients, err := h.service.DeleteListPatients(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve deleted patients", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Deleted patients retrieved successfully", patients)
}

func (h *PatientHandler) GetPatientByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid patient ID", err.Error())
		return
	}

	patient, err := h.service.GetPatientByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Patient not found", err.Error())
		return
	}

	// Ownership check for Patient role
	role, roleExists := ctx.Get("role")
	userID, userIDExists := ctx.Get("user_id")

	if roleExists && userIDExists && role == models.RolePatient {
		// Patient can only view their own data
		if patient.UserID == nil || *patient.UserID != userID.(uint) {
			utils.ErrorResponse(ctx, http.StatusForbidden, "You can only view your own patient data", nil)
			return
		}
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Patient retrieved successfully", patient)
}

func (h *PatientHandler) GetPatientByCode(ctx *gin.Context) {
	code := ctx.Param("code")

	patient, err := h.service.GetPatientByCode(code)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Patient not found", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Patient retrieved successfully", patient)
}

func (h *PatientHandler) CreatePatient(ctx *gin.Context) {
	var req dto.CreatePatientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	patient, err := h.service.CreatePatient(&req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "patient code already exists" {
			utils.ErrorResponse(ctx, http.StatusConflict, "Duplicate data", errMsg)
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create patient", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Patient created successfully", patient)
}

func (h *PatientHandler) UpdatePatient(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid patient ID", err.Error())
		return
	}

	var req dto.UpdatePatientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	patient, err := h.service.UpdatePatient(uint(id), &req)
	if err != nil {
		if err.Error() == "patient not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Patient not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update patient", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Patient updated successfully", patient)
}

func (h *PatientHandler) SoftDeletePatient(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid patient ID", err.Error())
		return
	}

	if err := h.service.SoftDeletePatient(uint(id)); err != nil {
		if err.Error() == "patient not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Patient not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete patient", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Patient deleted successfully", nil)
}

func (h *PatientHandler) RestorePatient(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid patient ID", err.Error())
		return
	}

	if err := h.service.RestorePatient(uint(id)); err != nil {
		if err.Error() == "patient not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Patient not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore patient", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Patient restored successfully", nil)
}

func (h *PatientHandler) HardDeletePatient(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid patient ID", err.Error())
		return
	}

	if err := h.service.HardDeletePatient(uint(id)); err != nil {
		if err.Error() == "patient not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Patient not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to permanently delete patient", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Patient permanently deleted", nil)
}
