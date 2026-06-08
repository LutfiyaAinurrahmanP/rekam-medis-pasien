package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medicine"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type MedicineHandler struct {
	service medicine.MedicineService
}

func NewMedicineHandler(s medicine.MedicineService) *MedicineHandler {
	return &MedicineHandler{
		service: s,
	}
}

func (h *MedicineHandler) ListMedicines(ctx *gin.Context) {
	var query dto.MedicinePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve medicines", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicines retrieved successfully", res)
}

func (h *MedicineHandler) DeletedListMedicines(ctx *gin.Context) {
	var query dto.MedicinePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve deleted medicines", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Deleted medicines retrieved successfully", res)
}

func (h *MedicineHandler) ActiveListMedicines(ctx *gin.Context) {
	var query dto.MedicinePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.ActiveList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve active medicines", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Active medicines retrieved successfully", res)
}

func (h *MedicineHandler) InactiveListMedicines(ctx *gin.Context) {
	var query dto.MedicinePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.InactiveList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve inactive medicines", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Inactive medicines retrieved successfully", res)
}

func (h *MedicineHandler) AvailableListMedicines(ctx *gin.Context) {
	var query dto.MedicinePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.AvailableList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve available medicines", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Available medicines retrieved successfully", res)
}

func (h *MedicineHandler) LowStockListMedicines(ctx *gin.Context) {
	var query dto.MedicinePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.LowStockList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve low stock medicines", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Low stock medicines retrieved successfully", res)
}

func (h *MedicineHandler) OutStockListMedicines(ctx *gin.Context) {
	var query dto.MedicinePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.OutStockList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve out of stock medicines", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Out of stock medicines retrieved successfully", res)
}

func (h *MedicineHandler) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Medicine not found", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine retrieved successfully", res)
}

func (h *MedicineHandler) CreateMedicine(ctx *gin.Context) {
	var req dto.CreateMedicineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create medicine", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "Medicine created successfully", res)
}

func (h *MedicineHandler) UpdateMedicine(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	var req dto.UpdateMedicineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update medicine", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine updated successfully", res)
}

func (h *MedicineHandler) AddStock(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	var req dto.AddStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	if err := h.service.AddStock(uint(id), &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add stock", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Stock added successfully", nil)
}

func (h *MedicineHandler) ReduceStock(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	var req dto.ReduceStockRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	if err := h.service.ReduceStock(uint(id), &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to reduce stock", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Stock reduced successfully", nil)
}

func (h *MedicineHandler) Activate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.Activate(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to activate medicine", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine activated successfully", nil)
}

func (h *MedicineHandler) Deactivate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	var req dto.DeactivateMedicineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	if err := h.service.Deactivate(uint(id), &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to deactivate medicine", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine deactivated successfully", nil)
}

func (h *MedicineHandler) SoftDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.SoftDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete medicine", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine deleted successfully", nil)
}

func (h *MedicineHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.Restore(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore medicine", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine restored successfully", nil)
}

func (h *MedicineHandler) HardDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.HardDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to permanently delete medicine", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine permanently deleted successfully", nil)
}
