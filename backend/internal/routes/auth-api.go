package routes

import "github.com/gin-gonic/gin"

func SetupAuthRouter(rg *gin.RouterGroup, cfg *RouteConfig) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", cfg.UserHandler.Register)
		authGroup.POST("/login", cfg.UserHandler.Login)
		authGroup.POST("/forgot-password", cfg.UserHandler.ForgotPassword)
		authGroup.POST("/verify-reset-code", cfg.UserHandler.VerifyResetCode)
		authGroup.PATCH("/reset-password", cfg.UserHandler.ResetPasswordWithToken)
		authGroup.POST("/resend-reset-code", cfg.UserHandler.ResendResetCode)
	}
}
