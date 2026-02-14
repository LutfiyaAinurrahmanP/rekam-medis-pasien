package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupDoctorRouter(rg *gin.RouterGroup, cfg *RouteConfig, doctorHandler *handler.DoctorHandler) {
	doctorGroup := rg.Group("/doctors")

	doctorGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		doctorSelfRoutes := doctorGroup.Group("/me")
		doctorSelfRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor))
		{
			doctorSelfRoutes.GET("", doctorHandler.GetMyDoctorData)
			doctorSelfRoutes.PUT("", doctorHandler.UpdateMyDoctorData)
		}

		doctorGroup.GET("", doctorHandler.ListDoctors)
		doctorGroup.GET("/:id", doctorHandler.GetDoctorByID)
		doctorGroup.GET("/specialization/:spec", doctorHandler.GetDoctorBySpecialization)

		adminRoutes := doctorGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.GET("/deleted", doctorHandler.DeletedListDoctors)
			adminRoutes.POST("", doctorHandler.CreateDoctor)
			adminRoutes.PUT("/:id", doctorHandler.UpdateDoctor)
			adminRoutes.PATCH("/:id/activate", doctorHandler.ActivateDoctor)
			adminRoutes.PATCH("/:id/deactivate", doctorHandler.DeactivateDoctor)
			adminRoutes.DELETE("/:id", doctorHandler.SoftDeleteDoctor)
			adminRoutes.PATCH("/:id/restore", doctorHandler.RestoreDoctor)
		}

		superAdminRoutes := doctorGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", doctorHandler.HardDeleteDoctor)
		}
	}

}
