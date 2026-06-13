package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupBillingRoutes(rg *gin.RouterGroup, cfg *RouteConfig, h *handler.BillingHandler) {
	g := rg.Group("/billing")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// ALL Authenticated Users
		g.GET("/:id", h.FindByID)
		g.GET("/patient/:patient_id", h.FindByPatientID)
		g.GET("/invoice/:invoice_number", h.FindByInvoiceNumber)
		g.GET("/:id/items", h.ListItems)
		g.GET("/:id/items/:item_id", h.FindItemByID)

		// Doctor, Receptionist, Admin, Super Admin
		staffRoutes := g.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			staffRoutes.GET("", h.List)
		}

		// Receptionist, Admin, Super Admin
		receptionistAdminRoutes := g.Group("")
		receptionistAdminRoutes.Use(middleware.RoleMiddleware(
			models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			receptionistAdminRoutes.POST("", h.Create)
			receptionistAdminRoutes.PUT("/:id", h.Update)
			receptionistAdminRoutes.PATCH("/:id/pay", h.RecordPayment)
			receptionistAdminRoutes.PATCH("/:id/cancel", h.Cancel)
			
			// Items
			receptionistAdminRoutes.POST("/:id/items", h.CreateItem)
			receptionistAdminRoutes.PUT("/:id/items/:item_id", h.UpdateItem)
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
			adminRoutes.DELETE("/:id/items/:item_id", h.DeleteItem)
		}

		// Super Admin Only
		superAdminRoutes := g.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", h.HardDelete)
		}
	}
}
