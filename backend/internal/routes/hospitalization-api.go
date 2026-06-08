package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupHospitalizationRouter(router *gin.RouterGroup, cfg *RouteConfig, h *handler.HospitalizationHandler) {
	api := router.Group("")
	
	api.Use(middleware.AuthMiddleware(cfg.Config))

	hospGroup := api.Group("/hospitalizations")
	{
		// 1. Read (Doctor, Receptionist, Admin, Super Admin)
		readRoutes := hospGroup.Group("")
		readRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			readRoutes.GET("", h.List)
			readRoutes.GET("/active", h.ActiveList)
			readRoutes.GET("/inactive", h.InactiveList)
		}

		// 2. Read All Authenticated
		allAuthRoutes := hospGroup.Group("")
		{
			allAuthRoutes.GET("/:id", h.FindByID)
		}

		// 3. Create (Receptionist, Admin, Super Admin)
		createRoutes := hospGroup.Group("")
		createRoutes.Use(middleware.RoleMiddleware(models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			createRoutes.POST("", h.Create)
		}

		// 4. Update (Doctor, Receptionist, Admin, Super Admin)
		updateRoutes := hospGroup.Group("")
		updateRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			updateRoutes.PUT("/:id", h.Update)
			updateRoutes.PATCH("/:id/discharge", h.Discharge)
			updateRoutes.PATCH("/:id/transfer", h.Transfer)
		}

		// 5. Deactivate & Activate (Admin, Super Admin)
		adminRoutes := hospGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/deleted", h.DeletedList)
			adminRoutes.PATCH("/:id/activate", h.Activate)
			adminRoutes.PATCH("/:id/deactivate", h.Deactivate)
			adminRoutes.DELETE("/:id", h.SoftDelete)
			adminRoutes.PATCH("/:id/restore", h.Restore)
		}

		// 6. Hard Delete (Super Admin)
		superAdminRoutes := hospGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", h.HardDelete)
		}
	}
}
