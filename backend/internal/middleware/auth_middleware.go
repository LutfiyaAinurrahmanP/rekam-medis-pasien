package middleware

import (
	"net/http"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHandler := ctx.GetHeader("Authorization")
		if authHandler == "" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Authorization header is required", nil)
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHandler, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid authorization header format", nil)
			ctx.Abort()
			return
		}

		token := parts[1]

		claims, err := utils.ValidateToken(token, cfg.JWT.Secret)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			ctx.Abort()
			return
		}

		userRole := role.(string)
		isAllowed := false

		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			utils.ErrorResponse(ctx, http.StatusForbidden, "Access denied: insufficent permissions", nil)
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
