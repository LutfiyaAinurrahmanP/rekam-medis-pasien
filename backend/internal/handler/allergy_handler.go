package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	allergy "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/allergy"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type AllergyHandler struct {
	service allergy.AllergyService
}

func NewAllergyHandler(service allergy.AllergyService) *AllergyHandler {
	return &AllergyHandler{service: service}
}

func (h *AllergyHandler) List(c *gin.Context) {
	var query dto.AllergyPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve allergies", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Allergies retrieved successfully", result)
}

func (h *AllergyHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	result, err := h.service.FindByID(uint(id))
	if err != nil {
		if err.Error() == "allergy not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Allergy not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve allergy", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Allergy retrieved successfully", result)
}

func (h *AllergyHandler) Create(c *gin.Context) {
	var req dto.CreateAllergyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create allergy", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Allergy created successfully", result)
}

func (h *AllergyHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdateAllergyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.Update(uint(id), &req)
	if err != nil {
		if err.Error() == "allergy not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Allergy not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update allergy", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Allergy updated successfully", result)
}

func (h *AllergyHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if err.Error() == "allergy not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Allergy not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete allergy", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Allergy deleted successfully", nil)
}
