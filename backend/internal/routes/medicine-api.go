package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupMedicineRouter(rg *gin.RouterGroup, cfg *RouteConfig, m *handler.MedicineHandler){
	g := rg.Group("/medicines")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		g.GET("", m.ListMedicines)
		g.GET("/available", m.ListByAvailable)
		g.GET("/name/:name", m.FindByName)
		g.GET("/:id", m.FindByID)

		// Stock routes (requires doctor, receptionist, admin, or superadmin)
		stockRoutes := g.Group("")
		stockRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			stockRoutes.GET("/low-stock", m.GetByLowStock)
		}
	}
}
