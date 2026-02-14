package routes

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/gin-gonic/gin"
)

// SetupRoomsRouter sets up room-related routes
func SetupRoomsRouter(rg *gin.RouterGroup, cfg *RouteConfig, roomHandler *handler.RoomHandler) {
	roomGroup := rg.Group("/rooms")

	roomGroup.Use(middleware.AuthMiddleware(cfg.Config))
	{
		// All authenticated users can view rooms
		roomGroup.GET("", roomHandler.ListRooms)
		roomGroup.GET("/available", roomHandler.GetAvailableRoom)
		roomGroup.GET("/:id", roomHandler.GetRoomByID)
		roomGroup.GET("/number/:room_number", roomHandler.GetRoomByRoomNumber)
		roomGroup.GET("/type/:room_type", roomHandler.GetRoomByRoomType)
		roomGroup.GET("/department/:dept_id", roomHandler.GetRoomByDepartmentID)

		// Staff routes (Doctor, Receptionist, Admin, SuperAdmin) - can view occupied rooms
		staffRoutes := roomGroup.Group("")
		staffRoutes.Use(middleware.RoleMiddleware(models.RoleDoctor, models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			staffRoutes.GET("/occupied", roomHandler.GetByOccupiedRoom)
		}

		// Admin routes (Receptionist, Admin, SuperAdmin) - for occupy/release
		receptionistRoutes := roomGroup.Group("")
		receptionistRoutes.Use(middleware.RoleMiddleware(models.RoleReceptionist, models.RoleAdmin, models.RoleSuperAdmin))
		{
			receptionistRoutes.GET("/inactive", roomHandler.GetByInactiveRoom)
			receptionistRoutes.PATCH("/:id/occupy", roomHandler.OccupyRoom)
			receptionistRoutes.PATCH("/:id/release", roomHandler.ReleaseRoom)
		}

		// Admin routes (Admin, SuperAdmin only) - for CRUD operations
		adminRoutes := roomGroup.Group("")
		adminRoutes.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleSuperAdmin))
		{
			adminRoutes.POST("", roomHandler.CreateRoom)
			adminRoutes.PUT("/:id", roomHandler.UpdateRoom)
			adminRoutes.PATCH("/:id/activate", roomHandler.ActivateRoom)
			adminRoutes.PATCH("/:id/deactivate", roomHandler.DeactivateRoom)
			adminRoutes.DELETE("/:id", roomHandler.SoftDeleteRoom)
			adminRoutes.GET("/deleted", roomHandler.DeletedListRooms)
			adminRoutes.PATCH("/:id/restore", roomHandler.RestoreRoom)
		}

		// SuperAdmin only routes
		superAdminRoutes := roomGroup.Group("")
		superAdminRoutes.Use(middleware.RoleMiddleware(models.RoleSuperAdmin))
		{
			superAdminRoutes.DELETE("/:id/hard-delete", roomHandler.HardDeleteRoom)
		}
	}
}
