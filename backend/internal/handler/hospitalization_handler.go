package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/hospitalization"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type HospitalizationHandler struct {
	service hospitalization.HospitalizationService
}

func NewHospitalizationHandler(service hospitalization.HospitalizationService) *HospitalizationHandler {
	return &HospitalizationHandler{service: service}
}

func (h *HospitalizationHandler) List(ctx *gin.Context) {
	var query dto.HospitalizationPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve hospitalizations", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Hospitalizations retrieved successfully", res)
}

func (h *HospitalizationHandler) DeletedList(ctx *gin.Context) {
	var query dto.HospitalizationPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve deleted hospitalizations", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Deleted hospitalizations retrieved successfully", res)
}

func (h *HospitalizationHandler) ActiveList(ctx *gin.Context) {
	var query dto.HospitalizationPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	query.Status = "admitted"
	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve active hospitalizations", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Active hospitalizations retrieved successfully", res)
}

func (h *HospitalizationHandler) InactiveList(ctx *gin.Context) {
	var query dto.HospitalizationPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	query.NotStatus = "admitted"
	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve inactive hospitalizations", err.Error())
		return
	}
	
	utils.SuccessResponse(ctx, http.StatusOK, "Inactive hospitalizations retrieved successfully", res)
}

func (h *HospitalizationHandler) FindByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Hospitalization not found", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Hospitalization retrieved successfully", res)
}

func (h *HospitalizationHandler) Create(ctx *gin.Context) {
	var req dto.CreateHospitalizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create hospitalization", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "Hospitalization created successfully", res)
}

func (h *HospitalizationHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdateHospitalizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update hospitalization", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Hospitalization updated successfully", res)
}

func (h *HospitalizationHandler) Discharge(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.DischargeHospitalizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Discharge(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to discharge patient", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Patient discharged successfully", res)
}

func (h *HospitalizationHandler) Transfer(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.TransferHospitalizationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Transfer(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to transfer patient", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Patient transferred successfully", res)
}

func (h *HospitalizationHandler) Activate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.Activate(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to activate hospitalization", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Hospitalization activated successfully", res)
}

func (h *HospitalizationHandler) Deactivate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.Deactivate(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to deactivate hospitalization", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Hospitalization deactivated successfully", res)
}

func (h *HospitalizationHandler) SoftDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete hospitalization", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Hospitalization deleted successfully", nil)
}

func (h *HospitalizationHandler) Restore(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to restore hospitalization", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Hospitalization restored successfully", nil)
}

func (h *HospitalizationHandler) HardDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to permanently delete hospitalization", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Hospitalization permanently deleted successfully", nil)
}
