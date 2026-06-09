package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupPrescriptionRoutes(rg *gin.RouterGroup, cfg *RouteConfig, h *handler.PrescriptionHandler) {
	// Protected endpoints
	prGroup := rg.Group("/prescriptions")
	prGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// 1. All Authenticated Routes
		allAuthRoutes := prGroup.Group("")
		allAuthRoutes.Use(middleware.RoleMiddleware(models.RolePatient, models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			allAuthRoutes.GET("/:id", h.FindByID) // access checked in handler
			allAuthRoutes.GET("/:id/items", h.ListItems)
			allAuthRoutes.GET("/:id/items/:itemID", h.FindItemByID)
		}

		// 2. Staff Routes (Doctor, Receptionist, Admin, Super Admin)
		staffRoutes := prGroup.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			staffRoutes.GET("", h.List)
			staffRoutes.GET("/medical-record/:recordID", h.PrescriptionsByMedicalRecordID)
			staffRoutes.PATCH("/:id/cancel", h.Cancel)
		}

		// 3. Doctor & Admin Routes (Doctor, Admin, Super Admin)
		doctorAdminRoutes := prGroup.Group("")
		doctorAdminRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleAdmin, models.RoleSuperAdmin))
		{
			doctorAdminRoutes.POST("", h.Create)
			doctorAdminRoutes.PUT("/:id", h.Update)
			doctorAdminRoutes.POST("/:id/items", h.CreateItem)
			doctorAdminRoutes.PUT("/:id/items/:itemID", h.UpdateItem)
			doctorAdminRoutes.DELETE("/:id/items/:itemID", h.DeleteItem)
		}

		// 4. Dispense Route (Receptionist, Admin, Super Admin)
		receptionistAdminRoutes := prGroup.Group("")
		receptionistAdminRoutes.Use(middleware.RoleMiddleware(models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			receptionistAdminRoutes.PATCH("/:id/dispense", h.Dispense)
		}

		// 5. Deletion & Restoration (Admin & Super Admin)
		adminRoutes := prGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/deleted", h.DeletedList)
			adminRoutes.DELETE("/:id", h.SoftDelete)
			adminRoutes.PATCH("/:id/restore", h.Restore)
		}

		// 6. Hard Delete (Super Admin Only)
		superAdminRoutes := prGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", h.HardDelete)
		}
	}
}
