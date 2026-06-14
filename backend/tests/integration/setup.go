package integration

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/utils"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDatabase() *gorm.DB {
	dbName := uuid.New().String()
	db, err := gorm.Open(sqlite.Open("file:"+dbName+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to in-memory database: " + err.Error())
	}

	// Migrate schemas
	err = db.AutoMigrate(
		&models.User{},
		&models.Patient{},
		&models.Allergy{},
		&models.MedicalCondition{},
		&models.SurgicalHistory{},
		&models.FamilyHistory{},
		&models.VitalSign{},
		&models.Doctor{},
		&models.Appointment{},
		&models.MedicalRecord{},
		&models.Hospitalization{},
		&models.Billing{},
		&models.BillingItem{},
		&models.Department{},
		&models.DoctorSpecialization{},
		&models.RoomType{},
		&models.Room{},
		&models.TypeTestCategory{},
		&models.TypeTest{},
		&models.MedicineType{},
		&models.Medicine{},
		&models.LabTest{},
		&models.Prescription{},
		&models.PrescriptionItem{},
	)
	if err != nil {
		panic("Failed to migrate test database: " + err.Error())
	}

	return db
}

func SetupTestConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			Env: "test",
		},
		JWT: config.JWTConfig{
			Secret: "supersecretkey123",
		},
		Database: config.DatabaseConfig{
			User: "test",
			Password: "password",
			Name: "testdb",
		},
	}
}

func GenerateTestToken(userID uint, role string, cfg *config.Config) string {
	tokenString, _, _ := utils.GenerateToken(userID, "testuser", "test@example.com", role, cfg.JWT.Secret, 24*time.Hour)
	return tokenString
}
