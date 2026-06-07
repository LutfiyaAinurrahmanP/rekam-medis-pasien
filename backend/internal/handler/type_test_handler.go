package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	typetestservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/type-test"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type TypeTestHandler struct {
	service typetestservice.TypeTestService
}

func NewTypeTestHandler(s typetestservice.TypeTestService) *TypeTestHandler {
	return &TypeTestHandler{service: s}
}

func (h *TypeTestHandler) List(ctx *gin.Context) {
	var q dto.TypeTestPaginationQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.List(&q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve type tests", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type tests retrieved successfully", res)
}

func (h *TypeTestHandler) ActiveList(ctx *gin.Context) {
	var q dto.TypeTestPaginationQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.ActiveList(&q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve active type tests", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Active type tests retrieved successfully", res)
}

func (h *TypeTestHandler) InactiveList(ctx *gin.Context) {
	var q dto.TypeTestPaginationQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.InactiveList(&q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve inactive type tests", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Inactive type tests retrieved successfully", res)
}

func (h *TypeTestHandler) DeletedList(ctx *gin.Context) {
	var q dto.TypeTestPaginationQuery
	if err := ctx.ShouldBindQuery(&q); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.DeletedList(&q)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve deleted type tests", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Deleted type tests retrieved successfully", res)
}

func (h *TypeTestHandler) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Type test not found", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test retrieved successfully", res)
}

func (h *TypeTestHandler) Create(ctx *gin.Context) {
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
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create type test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusCreated, "Type test created successfully", res)
}

func (h *TypeTestHandler) Update(ctx *gin.Context) {
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
		if err.Error() == "type test not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Type test not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update type test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test updated successfully", res)
}

func (h *TypeTestHandler) SoftDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.SoftDelete(uint(id)); err != nil {
		if err.Error() == "type test not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Type test not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete type test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test deleted successfully", nil)
}

func (h *TypeTestHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.Restore(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore type test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test restored successfully", nil)
}

func (h *TypeTestHandler) HardDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.HardDelete(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to hard delete type test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test permanently deleted", nil)
}

func (h *TypeTestHandler) Activate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.Activate(uint(id)); err != nil {
		if err.Error() == "type test not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Type test not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to activate type test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test activated successfully", nil)
}

func (h *TypeTestHandler) Deactivate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid id", err.Error())
		return
	}
	if err := h.service.Deactivate(uint(id)); err != nil {
		if err.Error() == "type test not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Type test not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to deactivate type test", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Type test deactivated successfully", nil)
}
