package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type TypeTestHandler struct {
	service service.TypeTestService
}

func NewTypeTestHandler(s service.TypeTestService) *TypeTestHandler {
	return &TypeTestHandler{service: s}
}

func (h *TypeTestHandler) ListTypeTests(ctx *gin.Context) {
	var q dto.TypeTestPaginationQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.List(&q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve test types", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Test types retrieved successfully", res)
}

func (h *TypeTestHandler) ListActiveTypeTests(ctx *gin.Context) {
	var q dto.TypeTestPaginationQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.ListActive(&q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve active test types", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Active test types retrieved successfully", res)
}

func (h *TypeTestHandler) ListInactiveTypeTests(ctx *gin.Context) {
	var q dto.TypeTestPaginationQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.ListInactive(&q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve inactive test types", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Inactive test types retrieved successfully", res)
}

func (h *TypeTestHandler) DeletedListTypeTests(ctx *gin.Context) {
	var q dto.TypeTestPaginationQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.DeleteList(&q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve deleted test types", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Deleted test types retrieved successfully", res)
}

func (h *TypeTestHandler) SearchTypeTests(ctx *gin.Context) {
	var q dto.TypeTestSearchQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.Search(&q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to search test types", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Search results retrieved successfully", res)
}

func (h *TypeTestHandler) GetTypeTestByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Test type not found", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Test type retrieved successfully", res)
}

func (h *TypeTestHandler) GetTypeTestByCode(ctx *gin.Context) {
	code := ctx.Param("code")
	res, err := h.service.FindByCode(code)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Test type not found", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Test type retrieved successfully", res)
}

func (h *TypeTestHandler) GetTypeTestsByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	var q dto.TypeTestPaginationQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.FindByCategory(category, &q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve test types by category", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Test types by category retrieved successfully", res)
}

func (h *TypeTestHandler) CreateTypeTest(ctx *gin.Context) {
	var req dto.CreateTypeTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.Create(&req)
	if err != nil {
		if err.Error() == "code already exists" {
			utils.ErrorResponse(ctx, http.StatusConflict, "Duplicate data", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create test type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "Test type created successfully", res)
}

func (h *TypeTestHandler) UpdateTypeTest(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	var req dto.UpdateTypeTestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.Update(uint(id), &req)
	if err != nil {
		if err.Error() == "code already exists" {
			utils.ErrorResponse(ctx, http.StatusConflict, "Duplicate data", err.Error())
			return
		}
		if err.Error() == "test type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Test type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update test type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Test type updated successfully", res)
}

func (h *TypeTestHandler) ActivateTypeTest(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.Activate(uint(id)); err != nil {
		if err.Error() == "test type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Test type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to activate test type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Test type activated successfully", nil)
}

func (h *TypeTestHandler) DeactivateTypeTest(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.Deactivate(uint(id)); err != nil {
		if err.Error() == "test type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Test type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to deactivate test type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Test type deactivated successfully", nil)
}

func (h *TypeTestHandler) SoftDeleteTypeTest(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.SoftDelete(uint(id)); err != nil {
		if err.Error() == "test type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Test type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete test type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Test type deleted successfully", nil)
}

func (h *TypeTestHandler) RestoreTypeTest(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.Restore(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore test type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Test type restored successfully", nil)
}

func (h *TypeTestHandler) HardDeleteTypeTest(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.HardDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to hard delete test type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Test type permanently deleted", nil)
}
