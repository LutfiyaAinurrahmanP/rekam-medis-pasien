package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupLabTestRouter(rg *gin.RouterGroup, cfg *RouteConfig, m *handler.LabTestHandler) {
	g := rg.Group("/lab-tests")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// ALL Authenticated Users
		g.GET("/:id", m.FindByID)
		g.GET("/medical-record/:record_id", m.FindByMedicalRecordID)

		// Any staff + doctor can cancel (ownership checked in service or handler)
		g.PATCH("/:id/cancel", middleware.RoleMiddleware(models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin), m.Cancel)

		// Receptionist, Admin, Super Admin
		receptionistAdminRoutes := g.Group("")
		receptionistAdminRoutes.Use(middleware.RoleMiddleware(
			models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			receptionistAdminRoutes.PATCH("/:id/collect-sample", m.CollectSample)
			receptionistAdminRoutes.PATCH("/:id/start", m.Start)
			receptionistAdminRoutes.PATCH("/:id/complete", m.Complete)
		}

		// Doctor, Admin, Super Admin
		doctorAdminRoutes := g.Group("")
		doctorAdminRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			doctorAdminRoutes.POST("", m.Create)
			doctorAdminRoutes.PUT("/:id", m.Update)
		}

		// Doctor, Receptionist, Admin, Super Admin
		staffRoutes := g.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			staffRoutes.GET("", m.List)
		}

		// Admin, Super Admin only
		adminRoutes := g.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/deleted", m.DeletedList)
			adminRoutes.DELETE("/:id", m.SoftDelete)
			adminRoutes.PATCH("/:id/restore", m.Restore)
		}

		// Super Admin Only
		superAdminRoutes := g.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", m.HardDelete)
		}
	}
}
