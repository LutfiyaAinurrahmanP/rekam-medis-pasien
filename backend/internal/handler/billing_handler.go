package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/billing"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type BillingHandler struct {
	service billing.BillingService
}

func NewBillingHandler(service billing.BillingService) *BillingHandler {
	return &BillingHandler{service: service}
}

func (h *BillingHandler) List(c *gin.Context) {
	var query dto.BillingPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	res, err := h.service.List(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve billings", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, res.Message, res)
}

func (h *BillingHandler) DeletedList(c *gin.Context) {
	var query dto.BillingPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	res, err := h.service.DeletedList(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve deleted billings", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, res.Message, res)
}

func (h *BillingHandler) FindByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Billing not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing retrieved successfully", res)
}

func (h *BillingHandler) FindByPatientID(c *gin.Context) {
	patientID, err := strconv.ParseUint(c.Param("patient_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid patient ID", nil)
		return
	}

	var query dto.BillingPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	res, err := h.service.FindByPatientID(uint(patientID), query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve patient billings", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Patient billing records retrieved successfully", res)
}

func (h *BillingHandler) FindByInvoiceNumber(c *gin.Context) {
	invoiceNumber := c.Param("invoice_number")

	res, err := h.service.FindByInvoiceNumber(invoiceNumber)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Billing not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing retrieved successfully", res)
}

func (h *BillingHandler) Create(c *gin.Context) {
	var req dto.CreateBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.Create(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create billing", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Billing created successfully", res)
}

func (h *BillingHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	var req dto.UpdateBillingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.Update(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update billing", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing updated successfully", res)
}

func (h *BillingHandler) RecordPayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	var req dto.RecordPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.RecordPayment(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to record payment", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment recorded successfully", res)
}

func (h *BillingHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	res, err := h.service.Cancel(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to cancel billing", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing cancelled successfully", res)
}

func (h *BillingHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete billing", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing deleted successfully", nil)
}

func (h *BillingHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	res, err := h.service.Restore(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to restore billing", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing restored successfully", res)
}

func (h *BillingHandler) HardDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to permanently delete billing", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing permanently deleted successfully", nil)
}

func (h *BillingHandler) ListItems(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	var query dto.BillingItemPaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", err.Error())
		return
	}

	res, err := h.service.ListItems(uint(id), query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve billing items", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing items retrieved successfully", res)
}

func (h *BillingHandler) FindItemByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID", nil)
		return
	}

	res, err := h.service.FindItemByID(uint(id), uint(itemID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Billing item not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing item retrieved successfully", res)
}

func (h *BillingHandler) CreateItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	var req dto.CreateBillingItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.CreateItem(uint(id), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create billing item", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Billing item added successfully", res)
}

func (h *BillingHandler) UpdateItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID", nil)
		return
	}

	var req dto.UpdateBillingItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.UpdateItem(uint(id), uint(itemID), req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update billing item", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing item updated successfully", res)
}

func (h *BillingHandler) DeleteItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid billing ID", nil)
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid item ID", nil)
		return
	}

	if err := h.service.DeleteItem(uint(id), uint(itemID)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete billing item", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Billing item deleted successfully", nil)
}
