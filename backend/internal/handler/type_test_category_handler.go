package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	typetestcategory "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/type-test-category"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type TypeTestCategoryHandler struct {
	service typetestcategory.TypeTestCategoryService
}

func NewTypeTestCategoryHandler(service typetestcategory.TypeTestCategoryService) *TypeTestCategoryHandler {
	return &TypeTestCategoryHandler{
		service: service,
	}
}

func (h *TypeTestCategoryHandler) List(ctx *gin.Context) {
	var query dto.TypeTestCategoryPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve type test categories", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Type test categories retrieve successfully", res)
}
func (h *TypeTestCategoryHandler) DeletedList(ctx *gin.Context) {
	var query dto.TypeTestCategoryPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to retrieve type test categories", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Deleted type test categories retrieve successfully", res)
}

func (h *TypeTestCategoryHandler) ActiveList(ctx *gin.Context) {
	var query dto.TypeTestCategoryPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.ActiveList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve active type test categories", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Active type test categories retrieved successfully", res)
}

func (h *TypeTestCategoryHandler) InactiveList(ctx *gin.Context) {
	var query dto.TypeTestCategoryPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.InactiveList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve inactive type test categories", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Inactive type test categories retrieved successfully", res)
}

func (h *TypeTestCategoryHandler) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid type test category ID", err.Error())
		return
	}
	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Type test category not found", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test category retrieved successfully", res)
}
func (h *TypeTestCategoryHandler) Create(ctx *gin.Context) {
	var req dto.CreateTypeTestCategoryRequest
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

		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create type test category", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Type test category created successfully", res)
}
func (h *TypeTestCategoryHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid type test category id", err.Error())
		return
	}

	var req dto.UpdateTypeTestCategoryRequest
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

		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update type test category", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Type test category update successfully", res)
}
func (h *TypeTestCategoryHandler) SoftDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid type test category id", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		if err.Error() == "type test category not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Type test category not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete type test category", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test category deleted successfully", nil)
}
func (h *TypeTestCategoryHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid type test category id", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		if err.Error() == "type test category not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Type test category not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore type test category", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test category restored successfully", nil)
}
func (h *TypeTestCategoryHandler) HardDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid type test category id", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		if err.Error() == "type test category not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Type test category not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to hard delete type test category", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test category hard deleted successfully", nil)
}

func (h *TypeTestCategoryHandler) Activate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid type test category id", err.Error())
		return
	}

	if err := h.service.Activate(uint(id)); err != nil {
		if err.Error() == "type test category not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Type test category not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to activate type test category", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test category activated successfully", nil)
}

func (h *TypeTestCategoryHandler) Deactivate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid type test category id", err.Error())
		return
	}

	if err := h.service.Deactivate(uint(id)); err != nil {
		if err.Error() == "type test category not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Type test category not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to deactivate type test category", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test category deactivated successfully", nil)
}
