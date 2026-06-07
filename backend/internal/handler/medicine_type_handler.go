package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	medicinetype "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medicine-type"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type MedicineTypeHandler struct {
	service medicinetype.MedicineTypeService
}

func NewMedicineTypeHandler(service medicinetype.MedicineTypeService) *MedicineTypeHandler {
	return &MedicineTypeHandler{
		service: service,
	}
}

func (h *MedicineTypeHandler) List(ctx *gin.Context) {
	var query dto.MedicineTypePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve medicine types", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Medicine types retrieve successfully", res)
}
func (h *MedicineTypeHandler) DeletedList(ctx *gin.Context) {
	var query dto.MedicineTypePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to retrieve medicine types", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Deleted medicine types retrieve successfully", res)
}

func (h *MedicineTypeHandler) ActiveList(ctx *gin.Context) {
	var query dto.MedicineTypePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.ActiveList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve active medicine types", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Active medicine types retrieved successfully", res)
}

func (h *MedicineTypeHandler) InactiveList(ctx *gin.Context) {
	var query dto.MedicineTypePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.InactiveList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve inactive medicine types", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Inactive medicine types retrieved successfully", res)
}

func (h *MedicineTypeHandler) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid medicine type ID", err.Error())
		return
	}
	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Medicine type not found", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine type retrieved successfully", res)
}
func (h *MedicineTypeHandler) Create(ctx *gin.Context) {
	var req dto.CreateMedicineTypeRequest
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

		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create medicine type", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Medicine type created successfully", res)
}
func (h *MedicineTypeHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid medicine type id", err.Error())
		return
	}

	var req dto.UpdateMedicineTypeRequest
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

		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update medicine type", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Medicine type update successfully", res)
}
func (h *MedicineTypeHandler) SoftDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid medicine type id", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		if err.Error() == "medicine type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Medicine type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete medicine type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine type deleted successfully", nil)
}
func (h *MedicineTypeHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid medicine type id", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		if err.Error() == "medicine type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Medicine type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore medicine type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine type restored successfully", nil)
}
func (h *MedicineTypeHandler) HardDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid medicine type id", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		if err.Error() == "medicine type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Medicine type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to hard delete medicine type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine type hard deleted successfully", nil)
}

func (h *MedicineTypeHandler) Activate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid medicine type id", err.Error())
		return
	}

	if err := h.service.Activate(uint(id)); err != nil {
		if err.Error() == "medicine type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Medicine type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to activate medicine type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine type activated successfully", nil)
}

func (h *MedicineTypeHandler) Deactivate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid medicine type id", err.Error())
		return
	}

	if err := h.service.Deactivate(uint(id)); err != nil {
		if err.Error() == "medicine type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Medicine type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to deactivate medicine type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Medicine type deactivated successfully", nil)
}
