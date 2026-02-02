package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
	service service.DepartmentService
}

func NewDepartmentHandler(service service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		service: service,
	}
}

func (h *DepartmentHandler) ListDepartments(ctx *gin.Context)  {
	var query dto.DepartmentPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	departments, err := h.service.ListDepartments(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve department", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Departments retrieved successfully", departments)
}
func (h *DepartmentHandler) GetDepartmentByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid department ID", err.Error())
		return
	}

	department, err := h.service.GetDepartmentByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Department not found", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Department retrieved successfully", department)
}

func (h *DepartmentHandler) CreateDepartment(ctx *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	department, err := h.service.CreateDepartment(&req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "code already exists" {
			utils.ErrorResponse(ctx, http.StatusConflict, "Duplicate data", errMsg)
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create department", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Department created successfully", department)
}

func (h *DepartmentHandler) UpdateDepartment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid department ID", err.Error())
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	department, err := h.service.UpdateDepartment(uint(id), &req)
	if err != nil {
		if err.Error() == "department not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Department not found", err.Error())
			return
		}

		errMsg := err.Error()
		if errMsg == "code already exists" {
			utils.ErrorResponse(ctx, http.StatusConflict, "Duplicate data", errMsg)
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update department", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Department updated successfully", department)
}

func (h *DepartmentHandler) SoftDeleteDepartment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid department ID", err.Error())
		return
	}

	err = h.service.SoftDeleteDepartment(uint(id))
	if err != nil {
		if err.Error() == "department not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Department not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete department", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Department deleted successfully", nil)
}

func (h *DepartmentHandler) RestoreDepartment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid department ID", err.Error())
		return
	}

	err = h.service.RestoreDepartment(uint(id))
	if err != nil {
		if err.Error() == "department not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Department not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore department", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Department restored successfully", nil)
}

func (h *DepartmentHandler) HardDeleteDepartment(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid department ID", err.Error())
		return
	}

	err = h.service.HardDeleteDepartment(uint(id))
	if err != nil {
		if err.Error() == "department not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Department not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to permanently delete department", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Department permanently deleted", nil)
}
func (h *DepartmentHandler) DeleteListDepartments(ctx *gin.Context) {
	var query dto.DepartmentPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	deletedDepartments, err := h.service.DeleteListDepartments(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve departments", deletedDepartments)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Deleted departments retrieved successfully", deletedDepartments)
}