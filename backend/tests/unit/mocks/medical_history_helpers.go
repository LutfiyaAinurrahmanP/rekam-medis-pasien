package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

func NewTestMedicalHistoryPatientWithData(id uint) *models.Patient {
	now := time.Now()
	p := &models.Patient{
		ID:       id,
		FullName: "Budi Santoso",
	}

	p.MedicalAllergies = []models.Allergy{
		{ID: 1, PatientID: id, AllergenType: "drug", AllergenName: "Penisilin", Reaction: "Ruam", Severity: "severe", CreatedAt: now, UpdatedAt: now},
	}
	p.MedicalConditions = []models.MedicalCondition{
		{ID: 1, PatientID: id, ConditionName: "Diabetes Mellitus Tipe 2", ICDCode: "E11", Status: "ongoing", CreatedAt: now, UpdatedAt: now},
	}
	p.SurgicalHistories = []models.SurgicalHistory{
		{ID: 1, PatientID: id, ProcedureName: "Appendektomi", SurgeryDate: now, SurgeonName: "dr. Bima", Hospital: "RS Dr. Soetomo", CreatedAt: now, UpdatedAt: now},
	}
	p.FamilyHistories = []models.FamilyHistory{
		{ID: 1, PatientID: id, FamilyMember: "father", ConditionName: "Hipertensi", Relation: "Biological Father", CreatedAt: now, UpdatedAt: now},
	}

	return p
}

func NewTestMedicalHistoryPatientList(count int) []models.Patient {
	var list []models.Patient
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestMedicalHistoryPatientWithData(uint(i)))
	}
	return list
}

func NewMedicalHistoryPaginationQuery(page, pageSize int) *dto.MedicalHistoryPaginationQuery {
	return &dto.MedicalHistoryPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}
