package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupRoomTypeRouter(rg *gin.RouterGroup, cfg *RouteConfig, roomTypeHandler *handler.RoomTypeHandler) {
	dsGroup := rg.Group("/room-types")

	dsGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		allRoutes := dsGroup.Group("")
		allRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleReceptionist, models.RoleDoctor, models.RoleAdmin, models.RoleSuperAdmin))
		{
			allRoutes.GET("", roomTypeHandler.List)
			allRoutes.GET("/active", roomTypeHandler.ActiveList)
			allRoutes.GET("/:id", roomTypeHandler.FindByID)
		}

		// admin and superadmin
		adminRoutes := dsGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			allRoutes.GET("/inactive", roomTypeHandler.InactiveList)
			adminRoutes.GET("/deleted", roomTypeHandler.DeletedList)
			adminRoutes.POST("", roomTypeHandler.Create)
			adminRoutes.PUT("/:id", roomTypeHandler.Update)
			adminRoutes.DELETE("/:id", roomTypeHandler.SoftDelete)
			adminRoutes.PATCH("/:id/restore", roomTypeHandler.Restore)
			adminRoutes.PATCH("/:id/activate", roomTypeHandler.Activate)
			adminRoutes.PATCH("/:id/deactivate", roomTypeHandler.Deactivate)
		}

		// super admin only
		superAdminRoutes := dsGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", roomTypeHandler.HardDelete)
		}
	}
}
