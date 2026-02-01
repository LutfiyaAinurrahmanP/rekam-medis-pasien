package handler

import (
	"net/http"

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