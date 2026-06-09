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
	appointmentservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/appointment"
	departmentservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/department"
	doctorservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/doctor"
	doctorspecializationservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/doctor-specialization"
	hospitalizationservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/hospitalization"
	labtestservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/lab-test"
	medicalrecordservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medical-record"
	medicineservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medicine"
	medicinetypeservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/medicine-type"
	patientservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/patient"
	prescriptionservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/prescription"
	roomservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/room"
	roomtypeservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/room-type"
	typetestservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/type-test"
	typetestcategorieservice "github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/service/type-test-category"
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
		Config:                      cfg,
		UserHandler:                 dependencies.UserHandler,
		DepartmentHandler:           dependencies.DepartmentHandler,
		PatientHandler:              dependencies.PatientHandler,
		DoctorSpecializationHandler: dependencies.DoctorSpecializationHandler,
		DoctorHandler:               dependencies.DoctorHandler,
		RoomTypeHandler:             dependencies.RoomTypeHandler,
		RoomHandler:                 dependencies.RoomHandler,
		TypeTestCategoryHandler:     dependencies.TypeTestCategoryHandler,
		TypeTestHandler:             dependencies.TypeTestHandler,
		MedicineHandler:             dependencies.MedicineHandler,
		MedicineTypeHandler:         dependencies.MedicineTypeHandler,
		AppointmentHandler:          dependencies.AppointmentHandler,
		MedicalRecordHandler:        dependencies.MedicalRecordHandler,
		HospitalizationHandler:      dependencies.HospitalizationHandler,
		LabTestHandler:              dependencies.LabTestHandler,
		PrescriptionHandler:         dependencies.PrescriptionHandler,
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
	UserRepository                 repository.UserRepository
	DepartmentRepository           repository.DepartmentRepository
	PatientRepository              repository.PatientRepository
	DoctorSpecializationRepository repository.DoctorSpecializationRepository
	DoctorRepository               repository.DoctorRepository
	RoomTypeRepository             repository.RoomTypeRepository
	RoomRepository                 repository.RoomRepository
	TypeTestCategoryRepository     repository.TypeTestCategoryRepository
	TypeTestRepository             repository.TypeTestRepository
	MedicineRepository             repository.MedicineRepository
	MedicineTypeRepository         repository.MedicineTypeRepository
	AppointmentRepository          repository.AppointmentRepository
	MedicalRecordRepository        repository.MedicalRecordRepository
	HospitalizationRepository      repository.HospitalizationRepository
	LabTestRepository              repository.LabTestRepository
	PrescriptionRepository         repository.PrescriptionRepository

	// Services
	UserService                 userservice.UserService
	DepartmentService           departmentservice.DepartmentService
	PatientService              patientservice.PatientService
	DoctorSpecializationService doctorspecializationservice.DoctorSpecializationService
	DoctorService               doctorservice.DoctorService
	RoomTypeService             roomtypeservice.RoomTypeService
	RoomService                 roomservice.RoomService
	TypeTestCategoryService     typetestcategorieservice.TypeTestCategoryService
	TypeTestService             typetestservice.TypeTestService
	MedicineService             medicineservice.MedicineService
	MedicineTypeService         medicinetypeservice.MedicineTypeService
	AppointmentService          appointmentservice.AppointmentService
	MedicalRecordService        medicalrecordservice.MedicalRecordService
	HospitalizationService      hospitalizationservice.HospitalizationService
	LabTestService              labtestservice.LabTestService
	PrescriptionService         prescriptionservice.PrescriptionService

	// Handlers
	UserHandler                 *handler.UserHandler
	DepartmentHandler           *handler.DepartmentHandler
	PatientHandler              *handler.PatientHandler
	DoctorSpecializationHandler *handler.DoctorSpecializationHandler
	DoctorHandler               *handler.DoctorHandler
	RoomTypeHandler             *handler.RoomTypeHandler
	RoomHandler                 *handler.RoomHandler
	TypeTestCategoryHandler     *handler.TypeTestCategoryHandler
	TypeTestHandler             *handler.TypeTestHandler
	MedicineHandler             *handler.MedicineHandler
	MedicineTypeHandler         *handler.MedicineTypeHandler
	AppointmentHandler          *handler.AppointmentHandler
	MedicalRecordHandler        *handler.MedicalRecordHandler
	HospitalizationHandler      *handler.HospitalizationHandler
	LabTestHandler              *handler.LabTestHandler
	PrescriptionHandler         *handler.PrescriptionHandler
}

// initDependencies initializes all application dependencies
func initDependencies(db *gorm.DB, cfg *config.Config, redisClient *cache.RedisClient, publisher kafka.EventPublisher) *Dependencies {
	// Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)
	patientRepo := repository.NewPatientRepository(db)
	doctorSpecializationRepo := repository.NewDoctorSpecializationRepository(db)
	doctorRepo := repository.NewDoctorRepository(db)
	roomTypeRepo := repository.NewRoomTypeRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	typeTestCategoryRepo := repository.NewTypeTestCategoryRepository(db)
	typeTestRepo := repository.NewTypeTestRepository(db)
	medicineRepo := repository.NewMedicineRepository(db)
	medicineTypeRepo := repository.NewMedicineTypeRepository(db)
	appointmentRepo := repository.NewAppointmentRepository(db)
	medicalRecordRepo := repository.NewMedicalRecordRepository(db)
	hospitalizationRepo := repository.NewHospitalizationRepository(db)
	labTestRepo := repository.NewLabTestRepository(db)
	prescriptionRepo := repository.NewPrescriptionRepository(db)

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

	doctorSpecializationService := doctorspecializationservice.NewEventDoctorSpecializationService(
		doctorspecializationservice.NewCachedDoctorSpecializationService(doctorspecializationservice.NewDoctorSpecializationService(doctorSpecializationRepo, cfg), redisClient),
		publisher,
	)

	doctorService := doctorservice.NewEventDoctorService(
		doctorservice.NewCachedDoctorService(doctorservice.NewDoctorService(doctorRepo, cfg), redisClient),
		publisher,
	)
	roomTypeService := roomtypeservice.NewEventRoomTypeService(
		roomtypeservice.NewCachedRoomTypeService(roomtypeservice.NewRoomTypeService(roomTypeRepo, cfg), redisClient),
		publisher,
	)
	roomService := roomservice.NewEventRoomService(
		roomservice.NewCachedRoomService(roomservice.NewRoomService(roomRepo, cfg), redisClient),
		publisher,
	)
	typeTestCategoryService := typetestcategorieservice.NewEventTypeTestCategoryService(
		typetestcategorieservice.NewCachedTypeTestCategoryService(typetestcategorieservice.NewTypeTestCategoryService(typeTestCategoryRepo, cfg), redisClient),
		publisher,
	)
	typeTestService := typetestservice.NewEventTypeTestService(
		typetestservice.NewCachedTypeTestService(typetestservice.NewTypeTestService(typeTestRepo, cfg), redisClient),
		publisher,
	)
	medicineService := medicineservice.NewEventMedicineService(
		medicineservice.NewCachedMedicineService(medicineservice.NewMedicineService(medicineRepo, cfg), redisClient),
		publisher,
	)
	medicineTypeService := medicinetypeservice.NewEventMedicineTypeService(
		medicinetypeservice.NewCachedMedicineTypeService(medicinetypeservice.NewMedicineTypeService(medicineTypeRepo, cfg), redisClient),
		publisher,
	)
	appointmentService := appointmentservice.NewEventAppointmentService(
		appointmentservice.NewCachedAppointmentService(appointmentservice.NewAppointmentService(appointmentRepo, cfg), redisClient),
		publisher,
	)
	medicalRecordService := medicalrecordservice.NewEventedMedicalRecordService(
		medicalrecordservice.NewCachedMedicalRecordService(medicalrecordservice.NewMedicalRecordService(medicalRecordRepo, cfg), redisClient),
		publisher,
	)
	hospitalizationService := hospitalizationservice.NewEventedHospitalizationService(
		hospitalizationservice.NewCachedHospitalizationService(hospitalizationservice.NewHospitalizationService(hospitalizationRepo, cfg), redisClient),
		publisher,
	)
	labTestService := labtestservice.NewEventedLabTestService(
		labtestservice.NewCachedLabTestService(labtestservice.NewLabTestService(labTestRepo, cfg), redisClient),
		publisher,
	)
	prescriptionService := prescriptionservice.NewEventPrescriptionService(
		prescriptionservice.NewCachedPrescriptionService(prescriptionservice.NewPrescriptionService(prescriptionRepo, cfg), redisClient),
		publisher,
	)

	// Initialize Handlers
	userHandler := handler.NewUserHandler(userService)
	departmentHandler := handler.NewDepartmentHandler(departmentService)
	patientHandler := handler.NewPatientHandler(patientService)
	doctorSpecializationHandler := handler.NewDoctorSpecializationHandler(doctorSpecializationService)
	doctorHandler := handler.NewDoctorHandler(doctorService)
	roomTypeHandler := handler.NewRoomTypeHandler(roomTypeService)
	roomHandler := handler.NewRoomHandler(roomService)
	typeTestCategoryHandler := handler.NewTypeTestCategoryHandler(typeTestCategoryService)
	typeTestHandler := handler.NewTypeTestHandler(typeTestService)
	medicineHandler := handler.NewMedicineHandler(medicineService)
	medicineTypeHandler := handler.NewMedicineTypeHandler(medicineTypeService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService, doctorRepo, patientRepo)
	medicalRecordHandler := handler.NewMedicalRecordHandler(medicalRecordService, doctorRepo, patientRepo)
	hospitalizationHandler := handler.NewHospitalizationHandler(hospitalizationService)
	labTestHandler := handler.NewLabTestHandler(labTestService, doctorRepo)
	prescriptionHandler := handler.NewPrescriptionHandler(prescriptionService, doctorRepo, patientRepo)

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

		DoctorSpecializationRepository: doctorSpecializationRepo,
		DoctorSpecializationService:    doctorSpecializationService,
		DoctorSpecializationHandler:    doctorSpecializationHandler,

		DoctorRepository: doctorRepo,
		DoctorService:    doctorService,
		DoctorHandler:    doctorHandler,

		RoomTypeRepository: roomTypeRepo,
		RoomTypeService:    roomTypeService,
		RoomTypeHandler:    roomTypeHandler,

		RoomRepository: roomRepo,
		RoomService:    roomService,
		RoomHandler:    roomHandler,

		TypeTestCategoryRepository: typeTestCategoryRepo,
		TypeTestCategoryService:    typeTestCategoryService,
		TypeTestCategoryHandler:    typeTestCategoryHandler,

		TypeTestRepository: typeTestRepo,
		TypeTestService:    typeTestService,
		TypeTestHandler:    typeTestHandler,

		MedicineRepository: medicineRepo,
		MedicineService:    medicineService,
		MedicineHandler:    medicineHandler,

		MedicineTypeRepository: medicineTypeRepo,
		MedicineTypeService:    medicineTypeService,
		MedicineTypeHandler:    medicineTypeHandler,

		AppointmentRepository: appointmentRepo,
		AppointmentService:    appointmentService,
		AppointmentHandler:    appointmentHandler,

		MedicalRecordRepository: medicalRecordRepo,
		MedicalRecordService:    medicalRecordService,
		MedicalRecordHandler:    medicalRecordHandler,

		HospitalizationRepository: hospitalizationRepo,
		HospitalizationService:    hospitalizationService,
		HospitalizationHandler:    hospitalizationHandler,

		LabTestRepository: labTestRepo,
		LabTestService:    labTestService,
		LabTestHandler:    labTestHandler,

		PrescriptionRepository: prescriptionRepo,
		PrescriptionService:    prescriptionService,
		PrescriptionHandler:    prescriptionHandler,
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
