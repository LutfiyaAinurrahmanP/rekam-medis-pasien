package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func setupAPIRouter(rg *gin.RouterGroup, cfg *RouteConfig) {
	apiGroup := rg.Group("/users")

	apiGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		meRoutes := apiGroup.Group("/me")
		{
			meRoutes.GET("", cfg.UserHandler.GetMyProfile)
			meRoutes.PUT("", cfg.UserHandler.UpdateMyProfile)
			meRoutes.PATCH("/change-password", cfg.UserHandler.ChangeMyPassword)
			meRoutes.DELETE("", cfg.UserHandler.DeleteMyAccount)
			meRoutes.PATCH("/deactivate", cfg.UserHandler.DeactivateMyAccount)
		}

		// apiGroup.PATCH("/:id/change-password", cfg.UserHandler.ChangePassword)
		adminRoutes := apiGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			// Users
			adminRoutes.POST("", cfg.UserHandler.CreateUser)
			adminRoutes.GET("", cfg.UserHandler.ListUsers)
			adminRoutes.GET("/deleted", cfg.UserHandler.DeleteListUsers)
			adminRoutes.GET("/:id", cfg.UserHandler.GetUserByID)
			adminRoutes.PUT("/:id", cfg.UserHandler.UpdateUser)
			adminRoutes.DELETE("/:id", cfg.UserHandler.SoftDeleteUser)
			adminRoutes.PATCH("/:id/restore", cfg.UserHandler.RestoreUser)
			adminRoutes.PATCH("/:id/reset-password", cfg.UserHandler.ResetPassword)
			adminRoutes.PATCH("/:id/activate", cfg.UserHandler.ActivateUser)
			adminRoutes.PATCH("/:id/deactivate", cfg.UserHandler.DeactivateUser)
		}

		superAdminRoutes := apiGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", cfg.UserHandler.HardDeleteUser)
		}
	}
}
