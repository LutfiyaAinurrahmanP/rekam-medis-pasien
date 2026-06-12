package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupMedicalConditionRoutes(router *gin.RouterGroup, cfg *RouteConfig, h *handler.MedicalConditionHandler) {
	mcGroup := router.Group("/medical-conditions")
	mcGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// 1. All Authenticated Routes
		allAuthRoutes := mcGroup.Group("")
		allAuthRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			allAuthRoutes.GET("/:id", h.FindByID)
		}

		// 2. Staff Routes (Doctor, Receptionist, Admin, Super Admin)
		staffRoutes := mcGroup.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			staffRoutes.GET("", h.List)
			staffRoutes.POST("", h.Create)
			staffRoutes.PUT("/:id", h.Update)
			staffRoutes.DELETE("/:id", h.Delete)
		}
	}
}
