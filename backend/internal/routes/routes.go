package routes

import (
	"net/http"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/middleware"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Config      *config.Config
	UserHandler *handler.UserHandler
}

func SetupRouter(cfg *RouteConfig) *gin.Engine {
	if cfg.Config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/health", healthCheck)
	router.GET("/", welcomePage)

	v1 := router.Group("/api/v1")
	{
		setupAuthRouter(v1, cfg)
		setupAPIRouter(v1, cfg)
	}

	return router
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":    "OK",
		"message":   "Sirekam Medis API is running",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}

func welcomePage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Sirekam Medis API",
		"version": "1.0.0",
		"docs":    "https://github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien",
		"endpoints": gin.H{
			"health": "/health",
			"api":    "/api/v1",
			"auth":   "/api/v1/auth",
			"users":  "/api/v1/users",
		},
	})
}
