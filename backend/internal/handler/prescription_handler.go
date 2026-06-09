package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/prescription"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type PrescriptionHandler struct {
	service     prescription.PrescriptionService
	doctorRepo  repository.DoctorRepository
	patientRepo repository.PatientRepository
}

func NewPrescriptionHandler(s prescription.PrescriptionService, dr repository.DoctorRepository, pr repository.PatientRepository) *PrescriptionHandler {
	return &PrescriptionHandler{
		service:     s,
		doctorRepo:  dr,
		patientRepo: pr,
	}
}

func (h *PrescriptionHandler) injectDoctorOwnership(ctx *gin.Context, query *dto.PrescriptionPaginationQuery) error {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil
	}
	role, _ := ctx.Get("role")
	uid := userID.(uint)

	if role == "doctor" {
		doctor, err := h.doctorRepo.FindByUserID(uid)
		if err == nil && doctor != nil {
			query.DoctorID = &doctor.ID
		}
	}
	return nil
}

func (h *PrescriptionHandler) List(ctx *gin.Context) {
	var query dto.PrescriptionPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	h.injectDoctorOwnership(ctx, &query)

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve prescriptions", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescriptions retrieved successfully", res)
}

func (h *PrescriptionHandler) DeletedList(ctx *gin.Context) {
	var query dto.PrescriptionPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve deleted prescriptions", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Deleted prescriptions retrieved successfully", res)
}

func (h *PrescriptionHandler) FindByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Prescription not found", err.Error())
		return
	}

	role, _ := ctx.Get("role")
	userID, _ := ctx.Get("user_id")
	uid := userID.(uint)

	if role == "doctor" {
		doctor, err := h.doctorRepo.FindByUserID(uid)
		if err == nil && doctor != nil && res.DoctorID != doctor.ID {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Forbidden", "You do not have access to this prescription")
			return
		}
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Prescription retrieved successfully", res)
}

func (h *PrescriptionHandler) Create(ctx *gin.Context) {
	var req dto.CreatePrescriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	role, _ := ctx.Get("role")
	userID, _ := ctx.Get("user_id")
	if role == "doctor" {
		doctor, err := h.doctorRepo.FindByUserID(userID.(uint))
		if err == nil && doctor != nil {
			req.DoctorID = doctor.ID
		}
	}

	res, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create prescription", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "Prescription created successfully", res)
}

func (h *PrescriptionHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdatePrescriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update prescription", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescription updated successfully", res)
}

func (h *PrescriptionHandler) Dispense(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.Dispense(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to dispense prescription", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescription dispensed successfully", res)
}

func (h *PrescriptionHandler) Cancel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.Cancel(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to cancel prescription", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescription cancelled successfully", res)
}

func (h *PrescriptionHandler) SoftDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete prescription", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescription deleted successfully", nil)
}

func (h *PrescriptionHandler) Restore(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to restore prescription", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescription restored successfully", nil)
}

func (h *PrescriptionHandler) HardDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to hard delete prescription", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescription permanently deleted", nil)
}

func (h *PrescriptionHandler) PrescriptionsByMedicalRecordID(ctx *gin.Context) {
	recordIDStr := ctx.Param("recordID")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid medical record ID format", err.Error())
		return
	}

	var query dto.PrescriptionPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	rid := uint(recordID)
	query.MedicalRecordID = &rid

	h.injectDoctorOwnership(ctx, &query)

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve prescriptions for medical record", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescriptions retrieved successfully", res)
}

func (h *PrescriptionHandler) ListItems(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.ListItems(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve prescription items", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescription items retrieved successfully", res)
}

func (h *PrescriptionHandler) FindItemByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid Prescription ID format", err.Error())
		return
	}

	itemIDStr := ctx.Param("itemID")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid Item ID format", err.Error())
		return
	}

	res, err := h.service.FindItemByID(uint(id), uint(itemID))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Prescription item not found", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Prescription item retrieved successfully", res)
}

func (h *PrescriptionHandler) CreateItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.CreatePrescriptionItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.CreateItem(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create prescription item", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "Prescription item created successfully", res)
}

func (h *PrescriptionHandler) UpdateItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid Prescription ID format", err.Error())
		return
	}

	itemIDStr := ctx.Param("itemID")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid Item ID format", err.Error())
		return
	}

	var req dto.UpdatePrescriptionItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.UpdateItem(uint(id), uint(itemID), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update prescription item", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescription item updated successfully", res)
}

func (h *PrescriptionHandler) DeleteItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid Prescription ID format", err.Error())
		return
	}

	itemIDStr := ctx.Param("itemID")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid Item ID format", err.Error())
		return
	}

	if err := h.service.DeleteItem(uint(id), uint(itemID)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete prescription item", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Prescription item deleted successfully", nil)
}
