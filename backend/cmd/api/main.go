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

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
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

	// Initialize Redis (optional — app tetap berjalan bila Redis tidak tersedia)
	var redisClient *cache.RedisClient
	if rc, err := cache.NewRedisClient(&cfg.Redis); err != nil {
		log.Printf("⚠️  Redis tidak tersedia, cache dinonaktifkan: %v", err)
	} else {
		redisClient = rc
	}

	// Initialize dependencies
	dependencies := initDependencies(db, cfg, redisClient)

	// Setup router with all routes
	router := routes.SetupRouter(&routes.RouteConfig{
		Config:            cfg,
		UserHandler:       dependencies.UserHandler,
		DepartmentHandler: dependencies.DepartmentHandler,
		PatientHandler:    dependencies.PatientHandler,
		DoctorHandler:     dependencies.DoctorHandler,
		RoomHandler:       dependencies.RoomHandler,
		TypeTestHandler:   dependencies.TypeTestHandler,
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
	gracefulShutdown(srv, db, redisClient)
}

// Dependencies holds all application dependencies
type Dependencies struct {
	// Repositories
	UserRepository       repository.UserRepository
	DepartmentRepository repository.DepartmentRepository
	PatientRepository    repository.PatientRepository
	DoctorRepository     repository.DoctorRepository
	RoomRepository       repository.RoomRepository
	TypeTestRepository   repository.TypeTestRepository

	// Services
	UserService       service.UserService
	DepartmentService service.DepartmentService
	PatientService    service.PatientService
	DoctorService     service.DoctorService
	RoomService       service.RoomService
	TypeTestService   service.TypeTestService

	// Handlers
	UserHandler       *handler.UserHandler
	DepartmentHandler *handler.DepartmentHandler
	PatientHandler    *handler.PatientHandler
	DoctorHandler     *handler.DoctorHandler
	RoomHandler       *handler.RoomHandler
	TypeTestHandler   *handler.TypeTestHandler
}

// initDependencies initializes all application dependencies
func initDependencies(db *gorm.DB, cfg *config.Config, redisClient *cache.RedisClient) *Dependencies {
	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	doctorRepo := repository.NewDoctorRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	typeTestRepo := repository.NewTypeTestRepository(db)

	// Initialize Services
	userService := service.NewCachedUserService(
		service.NewUserService(userRepo, cfg), redisClient,
	)
	departmentService := service.NewCachedDepartmentService(
		service.NewDepartmentService(departmentRepo, cfg), redisClient,
	)
	patientService := service.NewCachedPatientService(
		service.NewPatientService(patientRepo, cfg), redisClient,
	)
	doctorService := service.NewCachedDoctorService(
		service.NewDoctorService(doctorRepo, cfg), redisClient,
	)
	roomService := service.NewCachedRoomService(
		service.NewRoomService(roomRepo, cfg), redisClient,
	)
	typeTestService := service.NewCachedTypeTestService(
		service.NewTypeTestService(typeTestRepo, cfg), redisClient,
	)

	// Initialize Handlers
	userHandler := handler.NewUserHandler(userService)
	departmentHandler := handler.NewDepartmentHandler(departmentService)
	patientHandler := handler.NewPatientHandler(patientService)
	doctorHandler := handler.NewDoctorHandler(doctorService)
	roomHandler := handler.NewRoomHandler(roomService)
	typeTestHandler := handler.NewTypeTestHandler(typeTestService)

	return &Dependencies{
		UserRepository: userRepo,
		UserService:    userService,
		UserHandler:    userHandler,

		DepartmentRepository: departmentRepo,
		DepartmentService:    departmentService,
		DepartmentHandler:    departmentHandler,

		PatientRepository: patientRepo,
		PatientService:    patientService,
		PatientHandler:    patientHandler,

		DoctorRepository: doctorRepo,
		DoctorService:    doctorService,
		DoctorHandler:    doctorHandler,

		RoomRepository: roomRepo,
		RoomService:    roomService,
		RoomHandler:    roomHandler,

		TypeTestRepository: typeTestRepo,
		TypeTestService:    typeTestService,
		TypeTestHandler:    typeTestHandler,
	}
}

// gracefulShutdown handles graceful shutdown of the server
func gracefulShutdown(srv *http.Server, db *gorm.DB, redisClient *cache.RedisClient) {
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

	// Close Redis connection
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			log.Printf("❌ Failed to close Redis connection: %v", err)
		} else {
			log.Println("✅ Redis connection closed")
		}
	}

	log.Println("✅ Server exited gracefully")
}
