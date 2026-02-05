package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

// SetupPatientsRouter sets up patient-related routes
func SetupPatientsRouter(rg *gin.RouterGroup, cfg *RouteConfig, patientHandler *handler.PatientHandler) {
	patientGroup := rg.Group("/patients")

	patientGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// Patient self-service routes
		patientSelfRoutes := patientGroup.Group("/me")
		patientSelfRoutes.Use(middleware.RoleMiddleware(models.RolePatient))
		{
			patientSelfRoutes.GET("", patientHandler.GetMyPatientData)
			patientSelfRoutes.PUT("", patientHandler.UpdateMyPatientData)
		}

		// GET /:id - accessible by all authenticated users (with ownership check in handler)
		patientGroup.GET("/:id", patientHandler.GetPatientByID)

		// Staff can list patients (Doctor, Receptionist, Admin, SuperAdmin)
		staffRoutes := patientGroup.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			staffRoutes.GET("", patientHandler.ListPatients)
			staffRoutes.GET("/code/:code", patientHandler.GetPatientByCode)
			staffRoutes.GET("/search", patientHandler.SearchPatients)
		}

		// Receptionist, Admin, and SuperAdmin routes
		adminRoutes := patientGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.POST("", patientHandler.CreatePatient)
			adminRoutes.PUT("/:id", patientHandler.UpdatePatient)
			adminRoutes.DELETE("/:id", patientHandler.SoftDeletePatient)
			adminRoutes.GET("/deleted", patientHandler.DeleteListPatients)
			adminRoutes.PATCH("/:id/restore", patientHandler.RestorePatient)
		}

		// SuperAdmin only routes
		superAdminRoutes := patientGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", patientHandler.HardDeletePatient)
		}
	}
}
