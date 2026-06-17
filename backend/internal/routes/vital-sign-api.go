package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupVitalSignRoutes(router *gin.RouterGroup, cfg *RouteConfig, h *handler.VitalSignHandler) {
	vsGroup := router.Group("/vital-signs")
	vsGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// 1. All Authenticated Routes (All Users)
		allAuthRoutes := vsGroup.Group("")
		allAuthRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			allAuthRoutes.GET("/:id", h.FindByID)
		}

		// 2. Staff Routes (Doctor, Receptionist, Admin, Super Admin)
		staffRoutes := vsGroup.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			staffRoutes.GET("", h.List)
			staffRoutes.POST("", h.Create)
			staffRoutes.PUT("/:id", h.Update)
		}

		// 3. Admin & Super Admin Routes
		adminRoutes := vsGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/deleted", h.DeletedList)
			adminRoutes.DELETE("/:id", h.SoftDelete)
			adminRoutes.PATCH("/:id/restore", h.Restore)
		}

		// 4. Super Admin Only Routes
		superAdminRoutes := vsGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", h.HardDelete)
		}
	}
}
