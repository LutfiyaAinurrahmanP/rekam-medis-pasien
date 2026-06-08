package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupMedicalRecordRouter(rg *gin.RouterGroup, cfg *RouteConfig, h *handler.MedicalRecordHandler) {
	// Protected endpoints
	mrGroup := rg.Group("/medical-records")
	mrGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// 1. /my-records (Patient & Doctor)
		patientDoctorRoutes := mrGroup.Group("")
		patientDoctorRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleDoctor))
		{
			patientDoctorRoutes.GET("/my-records", h.MyMedicalRecords)
		}

		// 2. /patient/:patientID and /:id (All Authenticated)
		allAuthRoutes := mrGroup.Group("")
		allAuthRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			allAuthRoutes.GET("/patient/:patientID", h.MedicalRecordsByPatientID)
			allAuthRoutes.GET("/:id", h.FindByID) // Ownership checked inside handler
		}

		// 3. /medical-records (List) (Doctor, Receptionist, Admin, Super Admin)
		staffRoutes := mrGroup.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			staffRoutes.GET("", h.List)
		}

		// 4. Lifecycle & management (Doctor Only)
		doctorRoutes := mrGroup.Group("")
		doctorRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor))
		{
			doctorRoutes.POST("", h.Create)
			doctorRoutes.PATCH("/:id/finalize", h.Finalize)
		}

		doctorAdminRoutes := mrGroup.Group("")
		doctorAdminRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleAdmin, models.RoleSuperAdmin))
		{
			doctorAdminRoutes.PUT("/:id", h.Update)
		}

		// 5. Deletion & Restoration (Admin & Super Admin)
		adminRoutes := mrGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/deleted", h.DeletedList)
			adminRoutes.DELETE("/:id", h.SoftDelete)
			adminRoutes.PATCH("/:id/restore", h.Restore)
		}

		// 6. Hard Delete (Super Admin Only)
		superAdminRoutes := mrGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", h.HardDelete)
		}
	}
}
