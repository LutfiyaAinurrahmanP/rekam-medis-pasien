package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	vitalsign "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/vital-sign"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type VitalSignHandler struct {
	service vitalsign.VitalSignService
}

func NewVitalSignHandler(service vitalsign.VitalSignService) *VitalSignHandler {
	return &VitalSignHandler{service: service}
}

func (h *VitalSignHandler) List(c *gin.Context) {
	var query dto.VitalSignPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve vital signs", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Vital signs retrieved successfully", result)
}

func (h *VitalSignHandler) DeletedList(c *gin.Context) {
	var query dto.VitalSignPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve deleted vital signs", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Deleted vital signs retrieved successfully", result)
}

func (h *VitalSignHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	result, err := h.service.FindByID(uint(id))
	if err != nil {
		if err.Error() == "vital signs not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Vital signs not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve vital signs", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Vital signs retrieved successfully", result)
}

func (h *VitalSignHandler) Create(c *gin.Context) {
	var req dto.CreateVitalSignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.Create(&req)
	if err != nil {
		if err.Error() == "vital signs already recorded for this medical record" {
			utils.ErrorResponse(c, http.StatusConflict, "Vital signs already recorded", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create vital signs", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Vital signs recorded successfully", result)
}

func (h *VitalSignHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdateVitalSignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.Update(uint(id), &req)
	if err != nil {
		if err.Error() == "vital signs not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Vital signs not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update vital signs", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Vital signs updated successfully", result)
}

func (h *VitalSignHandler) SoftDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		if err.Error() == "vital signs not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Vital signs not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete vital signs", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Vital signs deleted successfully", nil)
}

func (h *VitalSignHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		if err.Error() == "vital signs not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Vital signs not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to restore vital signs", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Vital signs restored successfully", nil)
}

func (h *VitalSignHandler) HardDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		if err.Error() == "vital signs not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Vital signs not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to hard delete vital signs", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Vital signs permanently deleted successfully", nil)
}
