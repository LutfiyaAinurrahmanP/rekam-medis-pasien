package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupAllergyRoutes(router *gin.RouterGroup, cfg *RouteConfig, h *handler.AllergyHandler) {
	allergyGroup := router.Group("/allergies")
	allergyGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// 1. All Authenticated Routes
		allAuthRoutes := allergyGroup.Group("")
		allAuthRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			allAuthRoutes.GET("/:id", h.FindByID)
		}

		// 2. Staff Routes (Doctor, Receptionist, Admin, Super Admin)
		staffRoutes := allergyGroup.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			staffRoutes.GET("", h.List)
			staffRoutes.POST("", h.Create)
			staffRoutes.PUT("/:id", h.Update)
			staffRoutes.DELETE("/:id", h.Delete)
		}
	}
}
