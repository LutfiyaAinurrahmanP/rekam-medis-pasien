package handler

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	roomservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/room"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	service roomservice.RoomService
}

func NewRoomHandler(service roomservice.RoomService) *RoomHandler {
	return &RoomHandler{
		service: service,
	}
}

func (h *RoomHandler) ListRooms(ctx *gin.Context) {
	var query dto.RoomPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	rooms, err := h.service.ListRooms(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve rooms", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Rooms retrieved successfully", rooms)
}

func (h *RoomHandler) GetAvailableRoom(ctx *gin.Context) {
	var query dto.RoomPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	rooms, err := h.service.GetAvailableRooms(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve available rooms", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Available rooms retrieved successfully", rooms)
}

func (h *RoomHandler) GetByOccupiedRoom(ctx *gin.Context) {
	var query dto.RoomPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	rooms, err := h.service.GetOccupiedRooms(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve occupied rooms", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Occupied rooms retrieved successfully", rooms)
}

func (h *RoomHandler) GetByActiveRoom(ctx *gin.Context) {
	var query dto.RoomPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	rooms, err := h.service.GetActiveRooms(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve active rooms", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Active rooms retrieved successfully", rooms)
}

func (h *RoomHandler) GetByInactiveRoom(ctx *gin.Context) {
	var query dto.RoomPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	rooms, err := h.service.GetInactiveRooms(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve inactive rooms", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Inactive rooms retrieved successfully", rooms)
}

func (h *RoomHandler) DeletedListRooms(ctx *gin.Context) {
	var query dto.RoomPaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	rooms, err := h.service.DeleteListRooms(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve rooms", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Rooms retrieved successfully", rooms)
}

func (h *RoomHandler) GetRoomByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room ID", err.Error())
		return
	}

	room, err := h.service.GetRoomByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Room not found", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Room retrieved successfully", room)
}

func (h *RoomHandler) CreateRoom(ctx *gin.Context) {
	var req dto.CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	room, err := h.service.CreateRoom(&req)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "room number already exists" {
			utils.ErrorResponse(ctx, http.StatusConflict, "Duplicate data", errMsg)
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to created room", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Room created successfully", room)
}

func (h *RoomHandler) UpdateRoom(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room ID", err.Error())
		return
	}

	var req dto.UpdateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	room, err := h.service.UpdateRoom(uint(id), &req)
	if err != nil {
		if err.Error() == "room not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update room", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Room updated successfully", room)
}

func (h *RoomHandler) ActivateRoom(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room id", err.Error())
		return
	}

	room, err := h.service.ActivateRoom(uint(id))
	if err != nil {
		if err.Error() == "room not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to activate room", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Room activate successfully", room)
}

func (h *RoomHandler) DeactivateRoom(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room id", err.Error())
		return
	}

	room, err := h.service.DeactivateRoom(uint(id))
	if err != nil {
		if err.Error() == "room not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to activate room", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Room activate successfully", room)
}

func (h *RoomHandler) OccupyRoom(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room ID", err.Error())
		return
	}

	var req dto.OccupyRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	room, err := h.service.OccupyRoom(uint(id), req.Beds)
	if err != nil {
		if err.Error() == "room not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to occupy room", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Room occupied successfully", room)
}

func (h *RoomHandler) ReleaseRoom(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room ID", err.Error())
		return
	}

	var req dto.ReleaseRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	room, err := h.service.ReleaseRoom(uint(id), req.Beds)
	if err != nil {
		if err.Error() == "room not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to release room", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Room released successfully", room)
}

func (h *RoomHandler) SoftDeleteRoom(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room ID", err.Error())
		return
	}

	if err := h.service.SoftDeleteRoom(uint(id)); err != nil {
		if err.Error() == "room not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to soft delete room", err.Error())
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Room deleted successfully", nil)
}

func (h *RoomHandler) RestoreRoom(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room ID", err.Error())
		return
	}

	if err := h.service.RestoreRoom(uint(id)); err != nil {
		if err.Error() == "room not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room not foound", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retosre room", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Room restore successfully", nil)
}

func (h *RoomHandler) HardDeleteRoom(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid room id", err.Error())
		return
	}

	if err := h.service.HardDeleteRoom(uint(id)); err != nil {
		if err.Error() == "room not found" {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Room not found", err.Error())
			return
		}
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to hard delete room", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Room hard delete successfully", nil)
}
