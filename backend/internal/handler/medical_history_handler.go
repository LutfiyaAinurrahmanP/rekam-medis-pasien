package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	medicalhistory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type MedicalHistoryHandler struct {
	service medicalhistory.MedicalHistoryService
}

func NewMedicalHistoryHandler(service medicalhistory.MedicalHistoryService) *MedicalHistoryHandler {
	return &MedicalHistoryHandler{service: service}
}

func (h *MedicalHistoryHandler) List(c *gin.Context) {
	var query dto.MedicalHistoryPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve medical histories", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Medical histories retrieved successfully", result)
}

// FindByID returns the medical history aggregate for the specified patient ID.
func (h *MedicalHistoryHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	result, err := h.service.FindByID(uint(id))
	if err != nil {
		if err.Error() == "medical history for patient not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Medical history not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve medical history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Medical history retrieved successfully", result)
}

// FindByPatientID serves the same purpose as FindByID for the Medical History Aggregate.
func (h *MedicalHistoryHandler) FindByPatientID(c *gin.Context) {
	pid, err := strconv.ParseUint(c.Param("pid"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Patient ID format", err.Error())
		return
	}

	result, err := h.service.FindByID(uint(pid))
	if err != nil {
		if err.Error() == "medical history for patient not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Patient medical history not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve patient medical history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient medical history retrieved successfully", result)
}
