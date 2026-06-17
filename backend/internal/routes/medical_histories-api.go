package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupMedicalHistoryRoutes(rg *gin.RouterGroup, cfg *RouteConfig, m *handler.MedicalHistoryHandler) {
	g := rg.Group("/medical-history")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// ALL Authenticated Users
		g.GET("/:id", m.FindByID)
		g.GET("/patient/:pid", m.FindByPatientID)

		// Doctor, Receptionist, Admin, Super Admin
		staffRoutes := g.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			staffRoutes.GET("", m.List)
		}
	}
}
