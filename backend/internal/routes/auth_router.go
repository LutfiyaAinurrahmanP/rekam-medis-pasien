package routes

import "github.com/gin-gonic/gin"

func SetupAuthRouter(rg *gin.RouterGroup, cfg *RouteConfig) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/register", cfg.UserHandler.Register)
		authGroup.POST("/login", cfg.UserHandler.Login)
	}
}
