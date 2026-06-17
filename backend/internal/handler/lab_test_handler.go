package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/lab-test"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type LabTestHandler struct {
	service    labtest.LabTestService
	doctorRepo repository.DoctorRepository
}

func NewLabTestHandler(s labtest.LabTestService, dr repository.DoctorRepository) *LabTestHandler {
	return &LabTestHandler{
		service:    s,
		doctorRepo: dr,
	}
}

func (h *LabTestHandler) injectOwnership(ctx *gin.Context, query *dto.LabTestPaginationQuery) error {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil
	}

	role, _ := ctx.Get("role")
	uid := userID.(uint)

	if role == "doctor" {
		doctor, err := h.doctorRepo.FindByUserID(uid)
		if err == nil && doctor != nil {
			query.OrderedByDoctorID = &doctor.ID
		}
	}
	return nil
}

func (h *LabTestHandler) List(ctx *gin.Context) {
	var query dto.LabTestPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve lab tests", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Lab tests retrieved successfully", res)
}

func (h *LabTestHandler) DeletedList(ctx *gin.Context) {
	var query dto.LabTestPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve deleted lab tests", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Deleted lab tests retrieved successfully", res)
}

func (h *LabTestHandler) FindByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.FindByID(uint(id))
	if err != nil {
		if err.Error() == "lab test not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Not Found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve lab test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Lab test retrieved successfully", res)
}

func (h *LabTestHandler) FindByMedicalRecordID(ctx *gin.Context) {
	recordIDStr := ctx.Param("record_id")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid record ID format", err.Error())
		return
	}

	var query dto.LabTestPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	recID := uint(recordID)
	query.MedicalRecordID = &recID

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve medical record lab tests", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Medical record lab tests retrieved successfully", res.Data)
}

func (h *LabTestHandler) Create(ctx *gin.Context) {
	var req dto.CreateLabTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to order lab test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "Lab test ordered successfully", res)
}

func (h *LabTestHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdateLabTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update lab test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Lab test updated successfully", res)
}

func (h *LabTestHandler) CollectSample(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.CollectSample(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Failed to collect sample", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Sample collected successfully", res)
}

func (h *LabTestHandler) Start(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.Start(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Failed to start lab test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Lab test started successfully", res)
}

func (h *LabTestHandler) Complete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.CompleteLabTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Complete(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Failed to complete lab test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Lab test completed successfully", res)
}

func (h *LabTestHandler) Cancel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.Cancel(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Failed to cancel lab test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Lab test cancelled successfully", res)
}

func (h *LabTestHandler) SoftDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete lab test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Lab test deleted successfully", nil)
}

func (h *LabTestHandler) Restore(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to restore lab test", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Lab test restored successfully", nil)
}

func (h *LabTestHandler) HardDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to permanently delete lab test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Lab test permanently deleted successfully", nil)
}
