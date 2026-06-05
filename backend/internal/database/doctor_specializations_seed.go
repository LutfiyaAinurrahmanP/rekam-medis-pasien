package database

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func seedDoctorSpecializations(tx *gorm.DB, count int) ([]models.DoctorSpecialization, error) {
	specializations := []string{
		"Cardiology", "Neurology", "Pediatrics", "Obstetrics", "Orthopedics",
		"Radiology", "Internal Medicine", "General Surgery", "Pulmonology", "Dermatology",
		"Psychiatry", "Urology", "Gastroenterology", "Endocrinology", "Nephrology",
		"Rheumatology", "Hematology", "Infectious Disease", "Allergy and Immunology", "Geriatrics",
		"Otolaryngology", "Anesthesiology",
	}

	specs := make([]models.DoctorSpecialization, 0, len(specializations))
	for i, name := range specializations {
		code := fmt.Sprintf("SP-%03d", i+1)
		spec := models.DoctorSpecialization{
			Name:        name,
			Code:        code,
			Description: fmt.Sprintf("Specialization for %s", name),
			IsActive:    true,
		}
		specs = append(specs, spec)
	}

	if err := tx.Create(&specs).Error; err != nil {
		return nil, fmt.Errorf("failed to seed doctor specializations: %w", err)
	}

	return specs, nil
}

func seedDeletedDoctorSpecializations(tx *gorm.DB) error {
	deletedTime := time.Now().AddDate(0, -1, 0) // 1 month ago
	deletedAt := gorm.DeletedAt{Time: deletedTime, Valid: true}

	deletedSpecializations := []string{
		"Oncology", "Ophthalmology", "Pathology", "Preventive Medicine", "Physical Medicine and Rehabilitation",
		"Plastic Surgery", "Thoracic Surgery", "Vascular Surgery", "Medical Genetics", "Pain Medicine",
		"Sleep Medicine", "Sports Medicine", "Hospice and Palliative Medicine", "Addiction Medicine", "Occupational Medicine",
		"Medical Toxicology", "Undersea and Hyperbaric Medicine", "Clinical Neurophysiology", "Nuclear Medicine", "Interventional Radiology",
		"Neuroradiology", "Pediatric Surgery",
	}

	specs := make([]models.DoctorSpecialization, 0, len(deletedSpecializations))
	for i, name := range deletedSpecializations {
		code := fmt.Sprintf("SP-DEL-%03d", i+1)
		spec := models.DoctorSpecialization{
			Name:        name,
			Code:        code,
			Description: fmt.Sprintf("Deleted Specialization %s", name),
			IsActive:    false,
			DeletedAt:   deletedAt,
		}
		specs = append(specs, spec)
	}

	if err := tx.Create(&specs).Error; err != nil {
		return fmt.Errorf("failed to seed deleted doctor specializations: %w", err)
	}

	return nil
}
