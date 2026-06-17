package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

// SetupTypeTestRouter sets up test-types related routes
func SetupTypeTestRouter(rg *gin.RouterGroup, cfg *RouteConfig, th *handler.TypeTestHandler) {
	g := rg.Group("/type-tests")
	g.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// Public routes - all authenticated users
		g.GET("", th.List)
		g.GET("/active", th.ActiveList)

		// Doctor, Receptionist, Admin, Super Admin
		inactiveRoutes := g.Group("")
		inactiveRoutes.Use(middleware.RoleMiddleware(
			models.RoleDoctor, models.RoleReceptionist,
			models.RoleAdmin, models.RoleSuperAdmin,
		))
		{
			inactiveRoutes.GET("/inactive", th.InactiveList)
		}

		// Admin, Super Admin only
		adminRoutes := g.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/deleted", th.DeletedList)
			adminRoutes.POST("", th.Create)
			adminRoutes.PUT("/:id", th.Update)
			adminRoutes.PATCH("/:id/activate", th.Activate)
			adminRoutes.PATCH("/:id/deactivate", th.Deactivate)
			adminRoutes.DELETE("/:id", th.SoftDelete)
			adminRoutes.PATCH("/:id/restore", th.Restore)
		}

		// Super Admin only
		superAdmin := g.Group("")
		superAdmin.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdmin.DELETE("/:id/hard-delete", th.HardDelete)
		}

		// /:id must come last to avoid conflicting with named routes
		g.GET("/:id", th.FindByID)
	}
}
