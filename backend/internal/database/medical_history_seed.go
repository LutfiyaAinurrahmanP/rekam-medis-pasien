package database

import (
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func SeedMedicalHistory(db *gorm.DB) error {
	var patients []models.Patient
	if err := db.Find(&patients).Error; err != nil {
		return err
	}

	if len(patients) == 0 {
		log.Println("⚠️  No patients found, skipping medical history seeding")
		return nil
	}

	log.Println("🌱 Seeding medical history...")

	var allergies []models.Allergy
	var conditions []models.MedicalCondition
	var surgeries []models.SurgicalHistory
	var familyHistories []models.FamilyHistory

	now := time.Now()
	pastDate := now.AddDate(-5, 0, 0) // 5 years ago

	for i, patient := range patients {
		// Seed Allergies for some patients
		if i%3 == 0 {
			allergies = append(allergies, models.Allergy{
				PatientID:    patient.ID,
				AllergenType: "drug",
				AllergenName: "Penicillin",
				Reaction:     "Rash and shortness of breath",
				Severity:     "severe",
				Notes:        "Avoid all beta-lactams",
				CreatedAt:    now,
				UpdatedAt:    now,
			})
		}
		if i%5 == 0 {
			allergies = append(allergies, models.Allergy{
				PatientID:    patient.ID,
				AllergenType: "food",
				AllergenName: "Peanuts",
				Reaction:     "Hives",
				Severity:     "moderate",
				Notes:        "",
				CreatedAt:    now,
				UpdatedAt:    now,
			})
		}

		// Seed Medical Conditions
		if i%2 == 0 {
			conditions = append(conditions, models.MedicalCondition{
				PatientID:     patient.ID,
				ConditionName: "Hypertension",
				ICDCode:       "I10",
				DiagnosedDate: &pastDate,
				Status:        "ongoing",
				Notes:         "Controlled with medication",
				CreatedAt:     now,
				UpdatedAt:     now,
			})
		}

		// Seed Surgical History
		if i%4 == 0 {
			surgeries = append(surgeries, models.SurgicalHistory{
				PatientID:     patient.ID,
				ProcedureName: "Appendectomy",
				SurgeryDate:   pastDate,
				SurgeonName:   "Dr. Surgeon",
				Hospital:      "General Hospital",
				Notes:         "No complications",
				CreatedAt:     now,
				UpdatedAt:     now,
			})
		}

		// Seed Family History
		if i%3 == 0 {
			familyHistories = append(familyHistories, models.FamilyHistory{
				PatientID:     patient.ID,
				FamilyMember:  "father",
				ConditionName: "Diabetes Type 2",
				Notes:         "Diagnosed at age 50",
				CreatedAt:     now,
				UpdatedAt:     now,
			})
		}
	}

	if len(allergies) > 0 {
		if err := db.Create(&allergies).Error; err != nil {
			log.Printf("❌ Failed to seed allergies: %v", err)
			return err
		}
	}
	if len(conditions) > 0 {
		if err := db.Create(&conditions).Error; err != nil {
			log.Printf("❌ Failed to seed medical conditions: %v", err)
			return err
		}
	}
	if len(surgeries) > 0 {
		if err := db.Create(&surgeries).Error; err != nil {
			log.Printf("❌ Failed to seed surgical history: %v", err)
			return err
		}
	}
	if len(familyHistories) > 0 {
		if err := db.Create(&familyHistories).Error; err != nil {
			log.Printf("❌ Failed to seed family history: %v", err)
			return err
		}
	}

	return nil
}
