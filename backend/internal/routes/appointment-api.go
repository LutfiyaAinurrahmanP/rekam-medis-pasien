package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupAppointmentRouter(rg *gin.RouterGroup, cfg *RouteConfig, m *handler.AppointmentHandler) {
	g := rg.Group("/appointments")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// ALL Authenticated Users
		g.GET("/upcoming", m.UpcomingList)
		g.GET("/past", m.PastList)
		g.GET("/:id", m.FindByID)
		
		// Any staff + patient can cancel (ownership checked in service or handler)
		g.PATCH("/:id/cancel", m.Cancel)

		// Patient & Doctor Only
		patientDoctorRoutes := g.Group("")
		patientDoctorRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleDoctor))
		{
			patientDoctorRoutes.GET("/my-appointments", m.MyAppointments)
		}

		// Doctor Only
		doctorRoutes := g.Group("")
		doctorRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor))
		{
			doctorRoutes.GET("/my-schedule", m.MySchedule)
			doctorRoutes.PATCH("/:id/start", m.Start)
			doctorRoutes.PATCH("/:id/complete", m.Complete)
		}

		// Patient, Receptionist, Admin, Super Admin
		patientStaffRoutes := g.Group("")
		patientStaffRoutes.Use(middleware.RoleMiddleware(
			models.RolePatient, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			patientStaffRoutes.POST("", m.Create)
			patientStaffRoutes.PUT("/:id", m.Update)
			patientStaffRoutes.PATCH("/:id/reschedule", m.Reschedule)
		}

		// Doctor, Receptionist, Admin, Super Admin
		staffRoutes := g.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			staffRoutes.GET("", m.List)
			staffRoutes.GET("/today", m.TodayList)
			staffRoutes.PATCH("/:id/confirm", m.Confirm)
			staffRoutes.PATCH("/:id/no-show", m.NoShow)
		}

		// Receptionist, Admin, Super Admin
		receptionistAdminRoutes := g.Group("")
		receptionistAdminRoutes.Use(middleware.RoleMiddleware(
			models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			receptionistAdminRoutes.GET("/cancelled", m.CancelledList)
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
