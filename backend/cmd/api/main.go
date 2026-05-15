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
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/eventhandler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/handler"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/routes"
	departmentservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/department"
	doctorservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/doctor"
	medicineservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medicine"
	patientservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/patient"
	roomservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/room"
	typetestservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/typetest"
	userservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/user"
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

	// Initialize Kafka producer (opsional — app tetap berjalan bila Kafka tidak tersedia)
	var kafkaProducer kafka.EventPublisher
	if cfg.Kafka.Enabled {
		producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.ClientID)
		kafkaProducer = producer
		log.Printf("✅ Kafka producer aktif — brokers: %v", cfg.Kafka.Brokers)
	} else {
		log.Println("⚠️  Kafka dinonaktifkan — event publishing tidak aktif")
	}

	// Context untuk consumer goroutines (dibatalkan saat shutdown)
	consumerCtx, cancelConsumers := context.WithCancel(context.Background())

	// Start Kafka consumers jika Kafka aktif
	var auditHandler *eventhandler.AuditEventHandler
	var notifHandler *eventhandler.NotificationEventHandler
	if cfg.Kafka.Enabled {
		auditHandler = eventhandler.NewAuditEventHandler(cfg.Kafka.Brokers, kafka.AllTopics())
		notifHandler = eventhandler.NewNotificationEventHandler(cfg.Kafka.Brokers, cfg)

		go auditHandler.Start(consumerCtx)
		go notifHandler.Start(consumerCtx)
		log.Println("✅ Kafka consumers started (audit, notification)")
	}

	// Initialize dependencies
	dependencies := initDependencies(db, cfg, redisClient, kafkaProducer)

	// Setup router with all routes
	router := routes.SetupRouter(&routes.RouteConfig{
		Config:            cfg,
		UserHandler:       dependencies.UserHandler,
		DepartmentHandler: dependencies.DepartmentHandler,
		PatientHandler:    dependencies.PatientHandler,
		DoctorHandler:     dependencies.DoctorHandler,
		RoomHandler:       dependencies.RoomHandler,
		TypeTestHandler:   dependencies.TypeTestHandler,
		MedicineHandler: dependencies.MedicineHandler,
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
	gracefulShutdown(srv, db, redisClient, kafkaProducer, cancelConsumers, auditHandler, notifHandler)
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
	MedicineRepository repository.MedicineRepository

	// Services
	UserService       userservice.UserService
	DepartmentService departmentservice.DepartmentService
	PatientService    patientservice.PatientService
	DoctorService     doctorservice.DoctorService
	RoomService       roomservice.RoomService
	TypeTestService   typetestservice.TypeTestService
	MedicineService medicineservice.MedicineService

	// Handlers
	UserHandler       *handler.UserHandler
	DepartmentHandler *handler.DepartmentHandler
	PatientHandler    *handler.PatientHandler
	DoctorHandler     *handler.DoctorHandler
	RoomHandler       *handler.RoomHandler
	TypeTestHandler   *handler.TypeTestHandler
	MedicineHandler *handler.MedicineHandler
}

// initDependencies initializes all application dependencies
func initDependencies(db *gorm.DB, cfg *config.Config, redisClient *cache.RedisClient, publisher kafka.EventPublisher) *Dependencies {
	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	doctorRepo := repository.NewDoctorRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	typeTestRepo := repository.NewTypeTestRepository(db)
	medicineRepo := repository.NewMedicineRepository(db)

	// Konversi *RedisClient ke interface RedisStore hanya jika non-nil.
	// Tanpa ini, passing typed-nil ke parameter interface menghasilkan non-nil
	// interface yang membungkus nil pointer — menyebabkan nil-check di service gagal.
	var redisStore cache.RedisStore
	if redisClient != nil {
		redisStore = redisClient
	}

	// Initialize Services
	// Layer order: base → cache → event
	userService := userservice.NewEventUserService(
		userservice.NewCachedUserService(userservice.NewUserService(userRepo, cfg, redisStore), redisStore),
		publisher,
	)
	departmentService := departmentservice.NewEventDepartmentService(
		departmentservice.NewCachedDepartmentService(departmentservice.NewDepartmentService(departmentRepo, cfg), redisClient),
		publisher,
	)
	patientService := patientservice.NewEventPatientService(
		patientservice.NewCachedPatientService(patientservice.NewPatientService(patientRepo, cfg), redisClient),
		publisher,
	)
	doctorService := doctorservice.NewEventDoctorService(
		doctorservice.NewCachedDoctorService(doctorservice.NewDoctorService(doctorRepo, cfg), redisClient),
		publisher,
	)
	roomService := roomservice.NewEventRoomService(
		roomservice.NewCachedRoomService(roomservice.NewRoomService(roomRepo, cfg), redisClient),
		publisher,
	)
	typeTestService := typetestservice.NewEventTypeTestService(
		typetestservice.NewCachedTypeTestService(typetestservice.NewTypeTestService(typeTestRepo, cfg), redisClient),
		publisher,
	)
	medicineService := medicineservice.NewMedicineEventService(
		medicineservice.NewCachedMedicineService(medicineservice.NewMedicineService(medicineRepo, cfg), redisClient),
		publisher,
	)

	// Initialize Handlers
	userHandler := handler.NewUserHandler(userService)
	departmentHandler := handler.NewDepartmentHandler(departmentService)
	patientHandler := handler.NewPatientHandler(patientService)
	doctorHandler := handler.NewDoctorHandler(doctorService)
	roomHandler := handler.NewRoomHandler(roomService)
	typeTestHandler := handler.NewTypeTestHandler(typeTestService)
	medicineHandler := handler.NewMedicineHandler(medicineService)

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

		MedicineRepository: medicineRepo,
		MedicineService: medicineService,
		MedicineHandler: medicineHandler,
	}
}

// gracefulShutdown handles graceful shutdown of the server
func gracefulShutdown(
	srv *http.Server,
	db *gorm.DB,
	redisClient *cache.RedisClient,
	kafkaProducer kafka.EventPublisher,
	cancelConsumers context.CancelFunc,
	auditHandler *eventhandler.AuditEventHandler,
	notifHandler *eventhandler.NotificationEventHandler,
) {
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

	// Stop Kafka consumers
	cancelConsumers()
	time.Sleep(500 * time.Millisecond) // beri waktu goroutine selesai
	if auditHandler != nil {
		if err := auditHandler.Close(); err != nil {
			log.Printf("❌ Failed to close audit consumer: %v", err)
		} else {
			log.Println("✅ Audit consumer closed")
		}
	}
	if notifHandler != nil {
		if err := notifHandler.Close(); err != nil {
			log.Printf("❌ Failed to close notification consumer: %v", err)
		} else {
			log.Println("✅ Notification consumer closed")
		}
	}

	// Close Kafka producer
	if kafkaProducer != nil {
		if err := kafkaProducer.Close(); err != nil {
			log.Printf("❌ Failed to close Kafka producer: %v", err)
		} else {
			log.Println("✅ Kafka producer closed")
		}
	}

	log.Println("✅ Server exited gracefully")
}
