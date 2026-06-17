package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupMedicineTypeRouter(rg *gin.RouterGroup, cfg *RouteConfig, medicineTypeHandler *handler.MedicineTypeHandler) {
	medicineTypeGroup := rg.Group("/medicine-types")

	medicineTypeGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		allRoutes := medicineTypeGroup.Group("")
		allRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleReceptionist, models.RoleDoctor, models.RoleAdmin, models.RoleSuperAdmin))
		{
			allRoutes.GET("", medicineTypeHandler.List)
			allRoutes.GET("/active", medicineTypeHandler.ActiveList)
			allRoutes.GET("/:id", medicineTypeHandler.FindByID)
		}

		// admin and superadmin
		adminRoutes := medicineTypeGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			allRoutes.GET("/inactive", medicineTypeHandler.InactiveList)
			adminRoutes.GET("/deleted", medicineTypeHandler.DeletedList)
			adminRoutes.POST("", medicineTypeHandler.Create)
			adminRoutes.PUT("/:id", medicineTypeHandler.Update)
			adminRoutes.DELETE("/:id", medicineTypeHandler.SoftDelete)
			adminRoutes.PATCH("/:id/restore", medicineTypeHandler.Restore)
			adminRoutes.PATCH("/:id/activate", medicineTypeHandler.Activate)
			adminRoutes.PATCH("/:id/deactivate", medicineTypeHandler.Deactivate)
		}

		// super admin only
		superAdminRoutes := medicineTypeGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", medicineTypeHandler.HardDelete)
		}
	}
}
