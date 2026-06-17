package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	medicalcondition "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/medicalCondition"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type MedicalConditionHandler struct {
	service medicalcondition.MedicalConditionService
}

func NewMedicalConditionHandler(service medicalcondition.MedicalConditionService) *MedicalConditionHandler {
	return &MedicalConditionHandler{service: service}
}

func (h *MedicalConditionHandler) List(c *gin.Context) {
	var query dto.MedicalConditionPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve medical conditions", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Medical conditions retrieved successfully", result)
}

func (h *MedicalConditionHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	result, err := h.service.FindByID(uint(id))
	if err != nil {
		if err.Error() == "medical condition not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Medical condition not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve medical condition", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Medical condition retrieved successfully", result)
}

func (h *MedicalConditionHandler) Create(c *gin.Context) {
	var req dto.CreateMedicalConditionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create medical condition", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Medical condition created successfully", result)
}

func (h *MedicalConditionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdateMedicalConditionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.Update(uint(id), &req)
	if err != nil {
		if err.Error() == "medical condition not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Medical condition not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update medical condition", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Medical condition updated successfully", result)
}

func (h *MedicalConditionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if err.Error() == "medical condition not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Medical condition not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete medical condition", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Medical condition deleted successfully", nil)
}
