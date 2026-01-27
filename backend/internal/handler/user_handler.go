package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) getUserIDFromContext(ctx *gin.Context) (uint, error) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("user not authenticated")
	}
	return userID.(uint), nil
}

func (h *UserHandler) getUserRoleFromContext(ctx *gin.Context) (string, error) {
	role, exists := ctx.Get("role")
	if !exists {
		return "", fmt.Errorf("user role not found")
	}
	return role.(string), nil
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to register user", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "User registered successfully", user)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Login failed", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Login successful", response)
}

func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if err := h.service.ResetPassword(uint(id), req.NewPassword); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to reset password", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Password reset successfully", nil)
}

func (h *UserHandler) CreateUser(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	response, err := h.service.CreateUser(&req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to created user", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User created successfully", response)
}

func (h *UserHandler) GetUserByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "User not found", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "User retrieved successfully", user)
}

func (h *UserHandler) ListUsers(ctx *gin.Context) {
	var query dto.PaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	users, err := h.service.ListUsers(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve users", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Users retrieved successfully", users)
}

func (h *UserHandler) DeleteListUsers(ctx *gin.Context) {
	var query dto.PaginationQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	deletedUsers, err := h.service.DeleteListUsers(&query)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to retrieve users", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Deleted users retrieved successfully", deletedUsers)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user id", err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	user, err := h.service.UpdateUser(uint(id), &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to updateuser", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User updated successfuly", user)
}

func (h *UserHandler) SoftDeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	if err := h.service.SoftDeleteUser(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Failed to delete user", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User deleted successfully", nil)
}

func (h *UserHandler) HardDeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	if err := h.service.HardDeleteUser(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to permanently delete user", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User permanently delete successfully", nil)
}

func (h *UserHandler) RestoreUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	if err := h.service.RestoreUser(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to restore user", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User restored successfully", nil)
}

func (h *UserHandler) ActivateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	if err := h.service.ActivateUser(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to activated user", err.Error())
		return
	}
}

func (h *UserHandler) DeactivateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
		return
	}

	if err := h.service.DeactivateUser(uint(id)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to deactivate user", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "User deactivated successfully", nil)

}

// Personal handler
func (h *UserHandler) GetMyProfile(ctx *gin.Context) {
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "User not found", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Profile retrieved successfully", user)
}

func (h *UserHandler) UpdateMyProfile(ctx *gin.Context) {
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if req.Role != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, "Cannot change role via profile update", nil)
		return
	}

	if req.IsActive != nil {
		utils.ErrorResponse(ctx, http.StatusForbidden, "Cannot change account status via profile update", nil)
		return
	}

	user, err := h.service.UpdateUser(userID, &req)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update profile", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Profile updated successfully", user)
}

func (h *UserHandler) ChangeMyPassword(ctx *gin.Context) {
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	var req dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if err := h.service.ChangePassword(userID, &req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to change password", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Password changed successfully", nil)
}

func (h *UserHandler) DeleteMyAccount(ctx *gin.Context) {
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	var req dto.DeleteAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if err := h.service.VerifyPasswordForDeletion(userID, req.Password); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid password", err.Error())
		return
	}

	if err := h.service.SoftDeleteUser(userID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to deleted account", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Account deleted successfully", nil)
}

func (h *UserHandler) DeactivateMyAccount(ctx *gin.Context) {
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	var req dto.DeactivateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(ctx, err)
		return
	}

	if err := h.service.VerifyPasswordForDeletion(userID, req.Password); err != nil {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid password", err.Error())
		return
	}

	if err := h.service.DeactivateUser(userID); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Failed to deactivate account", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Account deactivate successfully. You can reactivate by contacting admin.", nil)
}
