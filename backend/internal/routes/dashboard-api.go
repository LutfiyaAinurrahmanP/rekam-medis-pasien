package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupDashboardRouter(rg *gin.RouterGroup, cfg *RouteConfig, m *handler.DashboardHandler) {
	g := rg.Group("/dashboard")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// Admin & Super Admin
		adminRoutes := g.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/overview", m.Overview)
			adminRoutes.GET("/admin", m.AdminDashboard)
			adminRoutes.GET("/reports/revenue", m.RevenueReport)
		}

		// Doctor Only
		doctorRoutes := g.Group("")
		doctorRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor))
		{
			doctorRoutes.GET("/doctor", m.DoctorDashboard)
		}

		// Receptionist Only
		receptionistRoutes := g.Group("")
		receptionistRoutes.Use(middleware.RoleMiddleware(models.RoleReceptionist))
		{
			receptionistRoutes.GET("/receptionist", m.ReceptionistDashboard)
		}

		// Patient Only
		patientRoutes := g.Group("")
		patientRoutes.Use(middleware.RoleMiddleware(models.RolePatient))
		{
			patientRoutes.GET("/patient", m.PatientDashboard)
		}

		// Doctor, Receptionist, Admin, Super Admin
		staffRoutes := g.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			staffRoutes.GET("/reports/appointments", m.AppointmentReport)
		}

		// Receptionist, Admin, Super Admin
		receptionistAdminRoutes := g.Group("")
		receptionistAdminRoutes.Use(middleware.RoleMiddleware(
			models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			receptionistAdminRoutes.GET("/reports/patients", m.PatientReport)
		}
	}
}
