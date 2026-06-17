package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupTypeTestCategoryRouter(rg *gin.RouterGroup, cfg *RouteConfig, typeTestCategoryHandler *handler.TypeTestCategoryHandler) {
	ttcGroup := rg.Group("/type-test-categories")

	ttcGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		allRoutes := ttcGroup.Group("")
		allRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleReceptionist, models.RoleDoctor, models.RoleAdmin, models.RoleSuperAdmin))
		{
			allRoutes.GET("", typeTestCategoryHandler.List)
			allRoutes.GET("/active", typeTestCategoryHandler.ActiveList)
			allRoutes.GET("/:id", typeTestCategoryHandler.FindByID)
		}

		// admin and superadmin
		adminRoutes := ttcGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/inactive", typeTestCategoryHandler.InactiveList)
			adminRoutes.GET("/deleted", typeTestCategoryHandler.DeletedList)
			adminRoutes.POST("", typeTestCategoryHandler.Create)
			adminRoutes.PUT("/:id", typeTestCategoryHandler.Update)
			adminRoutes.DELETE("/:id", typeTestCategoryHandler.SoftDelete)
			adminRoutes.PATCH("/:id/restore", typeTestCategoryHandler.Restore)
			adminRoutes.PATCH("/:id/activate", typeTestCategoryHandler.Activate)
			adminRoutes.PATCH("/:id/deactivate", typeTestCategoryHandler.Deactivate)
		}

		// super admin only
		superAdminRoutes := ttcGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", typeTestCategoryHandler.HardDelete)
		}
	}
}
