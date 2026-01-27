package middleware

import (
	"net/http"
	"strconv"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

func OwnershipMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authUserID, exists := ctx.Get("user_id")
		if !exists {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			ctx.Abort()
			return
		}

		role, roleExists := ctx.Get("role")
		if !roleExists {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			ctx.Abort()
			return
		}

		userRole := role.(string)

		if userRole == models.RoleAdmin || userRole == models.RoleSuperAdmin {
			ctx.Next()
			return
		}

		targetUserIDStr := ctx.Param("id")
		targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid user ID", err.Error())
			ctx.Abort()
			return
		}

		if authUserID.(uint) != uint(targetUserID) {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Access denied: You can only access your own data", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func OptionalOwnershipMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authUserID, exists := ctx.Get("user_id")
		if !exists {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			ctx.Abort()
			return
		}

		role, roleExists := ctx.Get("role")
		if !roleExists {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			ctx.Abort()
			return
		}

		userRole := role.(string)

		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				ctx.Next()
				return
			}
		}

		targetUserIDStr := ctx.Param("id")
		if targetUserIDStr != "" {
			targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalied user ID", err.Error())
				ctx.Abort()
				return
			}

			if authUserID.(uint) == uint(targetUserID) {
				ctx.Next()
				return
			}
		}

		utils.ErrorResponse(ctx, http.StatusForbidden, "Access denied: insufficient permissions", nil)
		ctx.Abort()
	}
}

func SelfOnlyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authUserID, exists := ctx.Get("user_id")
		if !exists {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			ctx.Abort()
			return
		}

		targetUserIDStr := ctx.Param("id")
		targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid User ID", err.Error())
			ctx.Abort()
			return
		}

		if authUserID.(uint) != uint(targetUserID) {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Access denied: You can only access your own data", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
