package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-record"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type MedicalRecordHandler struct {
	service     medicalrecord.MedicalRecordService
	doctorRepo  repository.DoctorRepository
	patientRepo repository.PatientRepository
}

func NewMedicalRecordHandler(s medicalrecord.MedicalRecordService, dr repository.DoctorRepository, pr repository.PatientRepository) *MedicalRecordHandler {
	return &MedicalRecordHandler{
		service:     s,
		doctorRepo:  dr,
		patientRepo: pr,
	}
}

func (h *MedicalRecordHandler) injectOwnership(ctx *gin.Context, query *dto.MedicalRecordPaginationQuery) error {
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
	} else if role == "patient" {
		patient, err := h.patientRepo.FindByUserID(uid)
		if err == nil && patient != nil {
			query.PatientID = &patient.ID
		}
	}
	return nil
}

func (h *MedicalRecordHandler) MyMedicalRecords(ctx *gin.Context) {
	var query dto.MedicalRecordPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	h.injectOwnership(ctx, &query)

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve my medical records", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "My medical records retrieved successfully", res)
}

func (h *MedicalRecordHandler) MedicalRecordsByPatientID(ctx *gin.Context) {
	patientIDStr := ctx.Param("patientID")
	patientID, err := strconv.ParseUint(patientIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid patient ID format", err.Error())
		return
	}

	var query dto.MedicalRecordPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	pid := uint(patientID)
	query.PatientID = &pid

	role, _ := ctx.Get("role")
	userID, _ := ctx.Get("user_id")
	if role == "patient" {
		uid := userID.(uint)
		patient, err := h.patientRepo.FindByUserID(uid)
		if err != nil || patient == nil || patient.ID != pid {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Forbidden", "You can only access your own medical records")
			return
		}
	}

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve patient medical records", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Patient medical records retrieved successfully", res)
}

func (h *MedicalRecordHandler) List(ctx *gin.Context) {
	var query dto.MedicalRecordPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve medical records", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medical records retrieved successfully", res)
}

func (h *MedicalRecordHandler) DeletedList(ctx *gin.Context) {
	var query dto.MedicalRecordPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve deleted medical records", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Deleted medical records retrieved successfully", res)
}

func (h *MedicalRecordHandler) FindByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Medical record not found", err.Error())
		return
	}

	role, _ := ctx.Get("role")
	userID, _ := ctx.Get("user_id")
	uid := userID.(uint)

	if role == "doctor" {
		doctor, err := h.doctorRepo.FindByUserID(uid)
		if err == nil && doctor != nil && res.DoctorID != doctor.ID {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Forbidden", "You do not have access to this medical record")
			return
		}
	} else if role == "patient" {
		patient, err := h.patientRepo.FindByUserID(uid)
		if err == nil && patient != nil && res.PatientID != patient.ID {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Forbidden", "You do not have access to this medical record")
			return
		}
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Medical record retrieved successfully", res)
}

func (h *MedicalRecordHandler) Create(ctx *gin.Context) {
	var req dto.CreateMedicalRecordRequest
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
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create medical record", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "Medical record created successfully", res)
}

func (h *MedicalRecordHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdateMedicalRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update medical record", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medical record updated successfully", res)
}

func (h *MedicalRecordHandler) Finalize(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Finalize(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to finalize medical record", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medical record finalized successfully", nil)
}

func (h *MedicalRecordHandler) SoftDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete medical record", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medical record deleted successfully", nil)
}

func (h *MedicalRecordHandler) Restore(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to restore medical record", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medical record restored successfully", nil)
}

func (h *MedicalRecordHandler) HardDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to hard delete medical record", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medical record permanently deleted", nil)
}
