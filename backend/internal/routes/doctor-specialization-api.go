package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupDoctorSpecializationRouter(rg *gin.RouterGroup, cfg *RouteConfig, doctorSpecializationHandler *handler.DoctorSpecializationHandler) {
	dsGroup := rg.Group("/doctor-specializations")

	dsGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// all roles
		allRoutes := dsGroup.Group("")
		allRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleReceptionist, models.RoleDoctor, models.RoleAdmin, models.RoleSuperAdmin))
		{
			allRoutes.GET("", doctorSpecializationHandler.List)
			allRoutes.GET("/:id", doctorSpecializationHandler.FindByID)
		}

		// admin and superadmin
		adminRoutes := dsGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/deleted", doctorSpecializationHandler.DeletedList)
			adminRoutes.POST("", doctorSpecializationHandler.Create)
			adminRoutes.PUT("/:id", doctorSpecializationHandler.Update)
			adminRoutes.DELETE("/:id", doctorSpecializationHandler.SoftDelete)
			adminRoutes.PATCH("/:id/restore", doctorSpecializationHandler.Restore)
		}

		// super admin only
		superAdminRoutes := dsGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", doctorSpecializationHandler.HardDelete)
		}
	}
}
