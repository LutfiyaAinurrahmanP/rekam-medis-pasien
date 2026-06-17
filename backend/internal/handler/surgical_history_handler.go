package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	surgicalhistory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/surgicalHistory"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type SurgicalHistoryHandler struct {
	service surgicalhistory.SurgicalHistoryService
}

func NewSurgicalHistoryHandler(service surgicalhistory.SurgicalHistoryService) *SurgicalHistoryHandler {
	return &SurgicalHistoryHandler{service: service}
}

func (h *SurgicalHistoryHandler) List(c *gin.Context) {
	var query dto.SurgicalHistoryPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve surgical histories", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Surgical histories retrieved successfully", result)
}

func (h *SurgicalHistoryHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	result, err := h.service.FindByID(uint(id))
	if err != nil {
		if err.Error() == "surgical history not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Surgical history not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve surgical history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Surgical history retrieved successfully", result)
}

func (h *SurgicalHistoryHandler) Create(c *gin.Context) {
	var req dto.CreateSurgicalHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create surgical history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Surgical history created successfully", result)
}

func (h *SurgicalHistoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdateSurgicalHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.Update(uint(id), &req)
	if err != nil {
		if err.Error() == "surgical history not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Surgical history not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update surgical history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Surgical history updated successfully", result)
}

func (h *SurgicalHistoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if err.Error() == "surgical history not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Surgical history not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete surgical history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Surgical history deleted successfully", nil)
}
