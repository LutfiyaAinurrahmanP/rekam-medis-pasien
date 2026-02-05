package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitTestDB() *gorm.DB {
	// Coba ambil dari environment variable dulu
	dbHost := os.Getenv("TEST_DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbUser := os.Getenv("TEST_DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}

	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "sirekam_test"
	}

	dbPort := os.Getenv("TEST_DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	// PostgreSQL DSN untuk testing
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		// Jika database test tidak ada, coba buat databasenya
		log.Printf("Trying to create test database: %s", dbName)

		// Connect ke postgres default database untuk membuat database test
		defaultDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
			dbHost, dbUser, dbPassword, dbPort)

		defaultDB, err2 := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if err2 != nil {
			log.Printf("Cannot connect to PostgreSQL. Make sure PostgreSQL is running.")
			log.Printf("You can start it with: docker-compose up -d postgres")
			log.Fatalf("Failed to connect to test database: %v", err)
		}

		// Buat database test
		sqlDB, _ := defaultDB.DB()
		_, execErr := sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if execErr != nil {
			log.Printf("Database may already exist or error creating: %v", execErr)
		}

		// Coba connect lagi
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})

		if err != nil {
			log.Fatalf("Failed to connect to test database after creation: %v", err)
		}
	}

	fmt.Println("[INFO] Test database connected successfully")
	return db
}
