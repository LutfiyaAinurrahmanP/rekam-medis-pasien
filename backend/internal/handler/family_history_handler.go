package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	familyhistory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-history/familyHistory"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type FamilyHistoryHandler struct {
	service familyhistory.FamilyHistoryService
}

func NewFamilyHistoryHandler(service familyhistory.FamilyHistoryService) *FamilyHistoryHandler {
	return &FamilyHistoryHandler{service: service}
}

func (h *FamilyHistoryHandler) List(c *gin.Context) {
	var query dto.FamilyHistoryPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve family histories", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Family histories retrieved successfully", result)
}

func (h *FamilyHistoryHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	result, err := h.service.FindByID(uint(id))
	if err != nil {
		if err.Error() == "family history not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Family history not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve family history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Family history retrieved successfully", result)
}

func (h *FamilyHistoryHandler) Create(c *gin.Context) {
	var req dto.CreateFamilyHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create family history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Family history created successfully", result)
}

func (h *FamilyHistoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdateFamilyHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	result, err := h.service.Update(uint(id), &req)
	if err != nil {
		if err.Error() == "family history not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Family history not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update family history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Family history updated successfully", result)
}

func (h *FamilyHistoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		if err.Error() == "family history not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Family history not found", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete family history", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Family history deleted successfully", nil)
}
