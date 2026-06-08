package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupMedicineRouter(rg *gin.RouterGroup, cfg *RouteConfig, m *handler.MedicineHandler) {
	g := rg.Group("/medicines")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		g.GET("", m.ListMedicines)
		g.GET("/active", m.ActiveListMedicines)
		g.GET("/available", m.AvailableListMedicines)
		g.GET("/:id", m.FindByID)

		// Stock routes (requires doctor, receptionist, admin, or superadmin)
		stockRoutes := g.Group("")
		stockRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			stockRoutes.GET("/low-stock", m.LowStockListMedicines)
			stockRoutes.GET("/out-of-stock", m.OutStockListMedicines)
		}

		// Receptionist routes
		receptionistRoutes := g.Group("")
		receptionistRoutes.Use(middleware.RoleMiddleware(models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			receptionistRoutes.GET("/inactive", m.InactiveListMedicines)
			receptionistRoutes.PATCH("/:id/add-stock", m.AddStock)
			receptionistRoutes.PATCH("/:id/reduce-stock", m.ReduceStock)
		}

		// Admin, Super Admin only
		adminRoutes := g.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.POST("", m.CreateMedicine)
			adminRoutes.PUT("/:id", m.UpdateMedicine)
			adminRoutes.PATCH("/:id/activate", m.Activate)
			adminRoutes.PATCH("/:id/deactivate", m.Deactivate)

			adminRoutes.GET("/deleted", m.DeletedListMedicines)
			adminRoutes.DELETE("/:id", m.SoftDelete)
			adminRoutes.PATCH("/:id/restore", m.Restore)
		}

		superAdminRoutes := g.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", m.HardDelete)
		}
	}
}
