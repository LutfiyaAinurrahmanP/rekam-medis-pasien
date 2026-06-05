package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	doctorspecialization "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/doctor-specialization"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type DoctorSpecializationHandler struct {
	service doctorspecialization.DoctorSpecializationService
}

func NewDoctorSpecializationHandler(service doctorspecialization.DoctorSpecializationService) *DoctorSpecializationHandler {
	return &DoctorSpecializationHandler{
		service: service,
	}
}

func (h *DoctorSpecializationHandler) List(ctx *gin.Context) {
	var query dto.DoctorSpecializationPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve doctor specializations", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Doctor specialization retrieve successfully", res)
}
func (h *DoctorSpecializationHandler) DeletedList(ctx *gin.Context) {
	var query dto.DoctorSpecializationPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to retrieve doctor specialization", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Deleted doctor specialization retrieve successfully", res)
}
func (h *DoctorSpecializationHandler) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor specialization ID", err.Error())
		return
	}
	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor specialization not found", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Doctor specialization retrieved successfully", res)
}
func (h *DoctorSpecializationHandler) Create(ctx *gin.Context) {
	var req dto.CreateDoctorSpecializationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Create(&req)
	if err != nil {
		errMsg := err.Error()
		switch errMsg {
		case "name already exists":
			utils.ErrorResponse(ctx, http.StatusConflict, "name already exists", err.Error())
			return
		case "code already exists":
			utils.ErrorResponse(ctx, http.StatusConflict, "code already exists", err.Error())
			return
		}

		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create doctor specialization", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Doctor specialization created successfully", res)
}
func (h *DoctorSpecializationHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor specialization id", err.Error())
		return
	}

	var req dto.UpdateDoctorSpecializationRequest
	if err := ctx.ShouldBind(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		errMsg := err.Error()
		switch errMsg {
		case "name already exists":
			utils.ErrorResponse(ctx, http.StatusConflict, "name already exists", err.Error())
			return
		case "code already exists":
			utils.ErrorResponse(ctx, http.StatusConflict, "code already exists", err.Error())
			return
		}

		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create doctor specialization", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Doctor specialization update successfully", res)
}
func (h *DoctorSpecializationHandler) SoftDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor specialization id", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		if err.Error() == "doctor specialization not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor specialization not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete doctor specialization", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Doctor specialization deleted successfully", nil)
}
func (h *DoctorSpecializationHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor specialization id", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		if err.Error() == "doctor specialization not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor specialization not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore doctor specialization", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Doctor specialization restored successfully", nil)
}
func (h *DoctorSpecializationHandler) HardDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid doctor specialization id", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		if err.Error() == "doctor specialization not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Doctor specialization not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to hard delete doctor specialization", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Doctor specialization hard deleted successfully", nil)
}
