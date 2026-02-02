package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

// SetupDepartmentsRouter sets up department-related routes
func SetupDepartmentsRouter(rg *gin.RouterGroup, cfg *RouteConfig, deptHandler *handler.DepartmentHandler) {
	deptGroup := rg.Group("/departments")

	deptGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// All authenticated users can list and get departments
		deptGroup.GET("", deptHandler.ListDepartments)
		deptGroup.GET("/:id", deptHandler.GetDepartmentByID)

		// Admin and SuperAdmin routes
		adminRoutes := deptGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.POST("", deptHandler.CreateDepartment)
			adminRoutes.PUT("/:id", deptHandler.UpdateDepartment)
			adminRoutes.DELETE("/:id", deptHandler.SoftDeleteDepartment)
			adminRoutes.GET("/deleted", deptHandler.DeleteListDepartments)
			adminRoutes.PATCH("/:id/restore", deptHandler.RestoreDepartment)
		}

		// SuperAdmin only routes
		superAdminRoutes := deptGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", deptHandler.HardDeleteDepartment)
		}
	}
}
