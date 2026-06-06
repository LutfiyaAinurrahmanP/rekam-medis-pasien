package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	roomtypeservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/room-type"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type RoomTypeHandler struct {
	service roomtypeservice.RoomTypeService
}

func NewRoomTypeHandler(service roomtypeservice.RoomTypeService) *RoomTypeHandler {
	return &RoomTypeHandler{
		service: service,
	}
}

func (h *RoomTypeHandler) List(ctx *gin.Context) {
	var query dto.RoomTypePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.List(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve room types", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Room type retrieve successfully", res)
}
func (h *RoomTypeHandler) DeletedList(ctx *gin.Context) {
	var query dto.RoomTypePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}
	res, err := h.service.DeletedList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to retrieve room type", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Deleted room type retrieve successfully", res)
}

func (h *RoomTypeHandler) ActiveList(ctx *gin.Context) {
	var query dto.RoomTypePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.ActiveList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve active room types", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Active room types retrieved successfully", res)
}

func (h *RoomTypeHandler) InactiveList(ctx *gin.Context) {
	var query dto.RoomTypePaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	res, err := h.service.InactiveList(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve inactive room types", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Inactive room types retrieved successfully", res)
}

func (h *RoomTypeHandler) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room type ID", err.Error())
		return
	}
	res, err := h.service.FindByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Room type not found", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Room type retrieved successfully", res)
}
func (h *RoomTypeHandler) Create(ctx *gin.Context) {
	var req dto.CreateRoomTypeRequest
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

		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create room type", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Room type created successfully", res)
}
func (h *RoomTypeHandler) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room type id", err.Error())
		return
	}

	var req dto.UpdateRoomTypeRequest
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

		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to create room type", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Room type update successfully", res)
}
func (h *RoomTypeHandler) SoftDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room type id", err.Error())
		return
	}

	if err := h.service.SoftDelete(uint(id)); err != nil {
		if err.Error() == "room type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to delete room type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Room type deleted successfully", nil)
}
func (h *RoomTypeHandler) Restore(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room type id", err.Error())
		return
	}

	if err := h.service.Restore(uint(id)); err != nil {
		if err.Error() == "room type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore room type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Room type restored successfully", nil)
}
func (h *RoomTypeHandler) HardDelete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room type id", err.Error())
		return
	}

	if err := h.service.HardDelete(uint(id)); err != nil {
		if err.Error() == "room type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to hard delete room type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Room type hard deleted successfully", nil)
}

func (h *RoomTypeHandler) Activate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room type id", err.Error())
		return
	}

	if err := h.service.Activate(uint(id)); err != nil {
		if err.Error() == "room type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to activate room type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Room type activated successfully", nil)
}

func (h *RoomTypeHandler) Deactivate(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room type id", err.Error())
		return
	}

	if err := h.service.Deactivate(uint(id)); err != nil {
		if err.Error() == "room type not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room type not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to deactivate room type", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Room type deactivated successfully", nil)
}
