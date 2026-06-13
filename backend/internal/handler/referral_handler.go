package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/referral"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type ReferralHandler struct {
	service referral.ReferralService
}

func NewReferralHandler(service referral.ReferralService) *ReferralHandler {
	return &ReferralHandler{service: service}
}

func (h *ReferralHandler) List(c *gin.Context) {
	var query dto.ReferralPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	res, err := h.service.List(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve referrals", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Referrals retrieved successfully", res)
}

func (h *ReferralHandler) DeletedList(c *gin.Context) {
	var query dto.ReferralPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	res, err := h.service.DeletedList(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve deleted referrals", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Deleted referrals retrieved successfully", res)
}

func (h *ReferralHandler) FindMyReferrals(c *gin.Context) {
	patientID, exists := c.Get("patient_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "Patient ID not found in token")
		return
	}

	status := c.Query("status")
	res, err := h.service.FindMyReferrals(patientID.(uint), status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve my referrals", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "My referrals retrieved successfully", res)
}

func (h *ReferralHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	res, err := h.service.FindByID(uint(id))
	if err != nil {
		if err.Error() == "referral not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Referral not found", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve referral", err.Error())
		return
	}

	// For patient role, verify ownership
	role, _ := c.Get("role")
	if role == "patient" {
		patientID, _ := c.Get("patient_id")
		if res.PatientID != patientID.(uint) {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", "You can only access your own referrals")
			return
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Referral retrieved successfully", res)
}

func (h *ReferralHandler) FindByPatientID(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("patient_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid patient ID format", err.Error())
		return
	}

	// For patient role, verify ownership
	role, _ := c.Get("role")
	if role == "patient" {
		loggedInPatientID, _ := c.Get("patient_id")
		if uint(patientID) != loggedInPatientID.(uint) {
			utils.ErrorResponse(c, http.StatusForbidden, "Forbidden", "You can only access your own referrals")
			return
		}
	}

	var query dto.ReferralPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	res, err := h.service.FindByPatientID(uint(patientID), query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve patient referrals", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient referrals retrieved successfully", res)
}

func (h *ReferralHandler) FindByDoctorID(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Param("doctor_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid doctor ID format", err.Error())
		return
	}

	var query dto.ReferralPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	res, err := h.service.FindByDoctorID(uint(doctorID), query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve doctor referrals", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Doctor referrals retrieved successfully", res)
}

func (h *ReferralHandler) Create(c *gin.Context) {
	var req dto.CreateReferralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err.Error())
		return
	}

	if req.ReferralType == "internal" && req.ReferredToDepartmentID == nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", "referred_to_department_id is required for internal referrals")
		return
	}
	if req.ReferralType == "external" && req.ReferredToFacility == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validation error", "referred_to_facility is required for external referrals")
		return
	}

	res, err := h.service.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create referral", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Referral created successfully", res)
}

func (h *ReferralHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.UpdateReferralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err.Error())
		return
	}

	res, err := h.service.Update(uint(id), req)
	if err != nil {
		if err.Error() == "referral not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Referral not found", nil)
			return
		}
		if err.Error() == "cannot update referral that is not in pending status" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid state", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update referral", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Referral updated successfully", res)
}

func (h *ReferralHandler) Accept(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.AcceptReferralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err.Error())
		return
	}

	res, err := h.service.Accept(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to accept referral", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Referral accepted successfully", res)
}

func (h *ReferralHandler) Reject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.RejectReferralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err.Error())
		return
	}

	res, err := h.service.Reject(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to reject referral", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Referral rejected", res)
}

func (h *ReferralHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.CompleteReferralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err.Error())
		return
	}

	res, err := h.service.Complete(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to complete referral", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Referral completed successfully", res)
}

func (h *ReferralHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	var req dto.CancelReferralRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid input data", err.Error())
		return
	}

	res, err := h.service.Cancel(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cancel referral", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Referral cancelled successfully", res)
}

func (h *ReferralHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete referral", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Referral deleted successfully", nil)
}

func (h *ReferralHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to restore referral", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Referral restored successfully", nil)
}

func (h *ReferralHandler) HardDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to permanently delete referral", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Referral permanently deleted successfully", nil)
}
