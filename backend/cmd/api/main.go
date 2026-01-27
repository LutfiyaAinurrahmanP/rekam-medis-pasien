package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/database"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/routes"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigration(db); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	// Seed database (optional - only for development)
	// if cfg.App.Env == "development" {
	// 	if err := database.SeedDatabase(db); err != nil {
	// 		log.Printf("⚠️  Failed to seed database: %v", err)
	// 	}
	// }

	// Initialize dependencies
	dependencies := initDependencies(db, cfg)

	// Setup router with all routes
	router := routes.SetupRouter(&routes.RouteConfig{
		Config:      cfg,
		UserHandler: dependencies.UserHandler,
		// Tambahkan handler lain di sini
		// PatientHandler:  dependencies.PatientHandler,
		// DoctorHandler:   dependencies.DoctorHandler,
	})

	// Setup HTTP server
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.App.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in goroutine
	go func() {
		log.Println("╔════════════════════════════════════════════════════╗")
		log.Println("║        Sirekam Medis API Server Started          ║")
		log.Println("╚════════════════════════════════════════════════════╝")
		log.Printf("🚀 Server running on port %s", cfg.App.Port)
		log.Printf("📝 Environment: %s", cfg.App.Env)
		log.Printf("🔗 API Base URL: http://localhost:%s/api/v1", cfg.App.Port)
		log.Printf("🏥 Health Check: http://localhost:%s/health", cfg.App.Port)
		log.Println("════════════════════════════════════════════════════")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(srv, db)
}

// Dependencies holds all application dependencies
type Dependencies struct {
	// Repositories
	UserRepository repository.UserRepository

	// Services
	UserService service.UserService

	// Handlers
	UserHandler *handler.UserHandler

	// Tambahkan dependencies lain di sini
	// PatientRepository repository.PatientRepository
	// PatientService    service.PatientService
	// PatientHandler    *handler.PatientHandler
}

// initDependencies initializes all application dependencies
func initDependencies(db *gorm.DB, cfg *config.Config) *Dependencies {
	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize Services
	userService := service.NewUserService(userRepo, cfg)

	// Initialize Handlers
	userHandler := handler.NewUserHandler(userService)

	return &Dependencies{
		UserRepository: userRepo,
		UserService:    userService,
		UserHandler:    userHandler,
	}
}

// gracefulShutdown handles graceful shutdown of the server
func gracefulShutdown(srv *http.Server, db *gorm.DB) {
	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	// Close database connection
	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("❌ Failed to close database connection: %v", err)
		} else {
			log.Println("✅ Database connection closed")
		}
	}

	log.Println("✅ Server exited gracefully")
}
