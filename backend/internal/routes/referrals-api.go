package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupReferralRoutes(rg *gin.RouterGroup, cfg *RouteConfig, h *handler.ReferralHandler) {
	g := rg.Group("/referrals")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// ALL Authenticated Users
		g.GET("/:id", h.FindByID)
		g.GET("/patient/:patient_id", h.FindByPatientID)

		// Patient
		patientRoutes := g.Group("")
		patientRoutes.Use(middleware.RoleMiddleware(models.RolePatient))
		{
			patientRoutes.GET("/my-referrals", h.FindMyReferrals)
		}

		// Doctor, Receptionist, Admin, Super Admin
		staffRoutes := g.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			staffRoutes.GET("", h.List)
			staffRoutes.GET("/doctor/:doctor_id", h.FindByDoctorID)
			staffRoutes.PATCH("/:id/complete", h.Complete)
			staffRoutes.PATCH("/:id/cancel", h.Cancel)
		}

		// Doctor, Admin, Super Admin
		doctorAdminRoutes := g.Group("")
		doctorAdminRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			doctorAdminRoutes.POST("", h.Create)
			doctorAdminRoutes.PUT("/:id", h.Update)
			doctorAdminRoutes.PATCH("/:id/accept", h.Accept)
			doctorAdminRoutes.PATCH("/:id/reject", h.Reject)
		}

		// Admin, Super Admin
		adminRoutes := g.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(
			models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			adminRoutes.GET("/deleted", h.DeletedList)
			adminRoutes.DELETE("/:id", h.Delete)
			adminRoutes.PATCH("/:id/restore", h.Restore)
		}

		// Super Admin Only
		superAdminRoutes := g.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", h.HardDelete)
		}
	}
}
