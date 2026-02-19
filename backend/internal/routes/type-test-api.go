package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

// SetupTypeTestRouter sets up test-types related routes
func SetupTypeTestRouter(rg *gin.RouterGroup, cfg *RouteConfig, th *handler.TypeTestHandler) {
	g := rg.Group("/test-types")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// Public routes - all authenticated users
		g.GET("", th.ListTypeTests)
		g.GET("/active", th.ListActiveTypeTests)
		g.GET("/search", th.SearchTypeTests)
		g.GET("/code/:code", th.GetTypeTestByCode)
		g.GET("/category/:category", th.GetTypeTestsByCategory)

		// Doctor, Receptionist, Admin, Super Admin
		inactiveRoutes := g.Group("")
		inactiveRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleReceptionist,
			models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			inactiveRoutes.GET("/inactive", th.ListInactiveTypeTests)
		}

		// Admin, Super Admin only
		adminRoutes := g.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/deleted", th.DeletedListTypeTests)
			adminRoutes.POST("", th.CreateTypeTest)
			adminRoutes.PUT("/:id", th.UpdateTypeTest)
			adminRoutes.PATCH("/:id/activate", th.ActivateTypeTest)
			adminRoutes.PATCH("/:id/deactivate", th.DeactivateTypeTest)
			adminRoutes.DELETE("/:id", th.SoftDeleteTypeTest)
			adminRoutes.PATCH("/:id/restore", th.RestoreTypeTest)
		}

		// Super Admin only
		superAdmin := g.Group("")
		superAdmin.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdmin.DELETE("/:id/hard-delete", th.HardDeleteTypeTest)
		}

		// /:id must come last to avoid conflicting with named routes
		g.GET("/:id", th.GetTypeTestByID)
	}
}

