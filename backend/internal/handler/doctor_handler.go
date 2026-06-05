package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	doctorservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/doctor"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type DoctorHandler struct {
	service doctorservice.DoctorService
}

func NewDoctorHandler(service doctorservice.DoctorService) *DoctorHandler {
	return &DoctorHandler{
		service: service,
	}
}

func (h *DoctorHandler) GetMyDoctorData(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	doctor, err := h.service.GetMyDoctorData(userID.(uint))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor data not found", err)
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Doctor data retrieved successfully", doctor)
}
func (h *DoctorHandler) UpdateMyDoctorData(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	var req dto.UpdateDoctorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	doctor, err := h.service.UpdateMyDoctorData(userID.(uint), &req)
	if err != nil {
		if err.Error() == "doctor not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor data not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update doctor data", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Doctor data updated successfully", doctor)
}
func (h *DoctorHandler) ListDoctors(ctx *gin.Context) {
	var query dto.DoctorPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	doctors, err := h.service.ListDoctors(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve doctors", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Doctors retrieved successfully", doctors)
}
func (h *DoctorHandler) DeletedListDoctors(ctx *gin.Context) {
	var query dto.DoctorPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	doctors, err := h.service.DeletedListDoctors(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve deleted doctors", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Deleted doctors retrieved successfully", doctors)
}
func (h *DoctorHandler) GetDoctorByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor ID", err.Error())
		return
	}

	doctor, err := h.service.GetDoctorByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor noy found", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Doctor retrieved successfully", doctor)
}
func (h *DoctorHandler) GetDoctorBySpecializationID(ctx *gin.Context) {
	specID, err := strconv.ParseUint(ctx.Param("spec_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid specialization ID", err.Error())
		return
	}

	doctor, err := h.service.GetDoctorBySpecializationID(uint(specID))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor not found", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Doctor retrieved successfully", doctor)
}
func (h *DoctorHandler) CreateDoctor(ctx *gin.Context) {
	var req dto.CreateDoctorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	doctor, err := h.service.CreateDoctor(&req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "employee already exists" {
			utils.ErrorResponse(ctx, http.StatusConflict, "Duplicate data", errMsg)
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to created doctor", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Doctor created successfully", doctor)
}
func (h *DoctorHandler) UpdateDoctor(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor ID", err.Error())
		return
	}

	var req dto.UpdateDoctorRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	doctor, err := h.service.UpdateDoctor(uint(id), &req)
	if err != nil {
		if err.Error() == "doctor not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update doctor", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Doctor updated successfully", doctor)
}

func (h *DoctorHandler) ActivateDoctor(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor ID", err.Error())
		return
	}

	doctor, err := h.service.ActivateDoctor(uint(id))
	if err != nil {
		if err.Error() == "doctor not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to activate doctor", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Doctor activated successfully", doctor)
}

func (h *DoctorHandler) DeactivateDoctor(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor ID", err.Error())
		return
	}

	doctor, err := h.service.DeactivateDoctor(uint(id))
	if err != nil {
		if err.Error() == "doctor not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to deactivate doctor", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Doctor deactivated successfully", doctor)
}

func (h *DoctorHandler) SoftDeleteDoctor(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor ID", err.Error())
		return
	}

	if err := h.service.SoftDeleteDoctor(uint(id)); err != nil {
		if err.Error() == "doctor not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete doctor", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Doctor deleted successfully", nil)
}
func (h *DoctorHandler) RestoreDoctor(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor ID", err.Error())
		return
	}

	if err := h.service.RestoreDoctor(uint(id)); err != nil {
		if err.Error() == "doctor not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore doctor", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Doctor restored successfully", nil)
}
func (h *DoctorHandler) HardDeleteDoctor(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor ID", err.Error())
		return
	}

	if err := h.service.HardDeleteDoctor(uint(id)); err != nil {
		if err.Error() == "doctor not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to permanently delete doctor", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Doctor permanently deleted", nil)
}
