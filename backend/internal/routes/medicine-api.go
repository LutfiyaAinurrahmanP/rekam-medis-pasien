package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupMedicineRouter(rg *gin.RouterGroup, cfg *RouteConfig, m *handler.MedicineHandler){
	g := rg.Group("/medicines")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		g.GET("", m.ListMedicines)
		g.GET("/available", m.GetByAvailable)
	}
}
