package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func NewTestLabTestWithData(id uint, medicalRecordID, testTypeID, orderedByDoctorID uint, status string, isDeleted bool) *models.LabTest {
	now := time.Now()
	val := "Normal"
	unit := "mg/dL"
	ref := "10-20"
	lt := &models.LabTest{
		ID:                id,
		MedicalRecordID:   medicalRecordID,
		TestTypeID:        testTypeID,
		OrderedByDoctorID: orderedByDoctorID,
		OrderDate:         "2023-12-01",
		Status:            status,
		ResultValue:       &val,
		ResultUnit:        &unit,
		ReferenceRange:    &ref,
		Notes:             "Test Notes",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if isDeleted {
		lt.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	}

	return lt
}

func NewTestLabTestList(count int) []models.LabTest {
	var list []models.LabTest
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestLabTestWithData(uint(i), 1, 1, 1, "ordered", false))
	}
	return list
}

func NewTestLabTestResponse(lt *models.LabTest) *dto.LabTestResponse {
	return &dto.LabTestResponse{
		ID:                lt.ID,
		MedicalRecordID:   lt.MedicalRecordID,
		TestTypeID:        lt.TestTypeID,
		OrderedByDoctorID: lt.OrderedByDoctorID,
		OrderDate:         lt.OrderDate,
		Status:            lt.Status,
		ResultValue:       lt.ResultValue,
		ResultUnit:        lt.ResultUnit,
		ReferenceRange:    lt.ReferenceRange,
		Notes:             lt.Notes,
		CreatedAt:         lt.CreatedAt,
		UpdatedAt:         lt.UpdatedAt,
	}
}

func NewLabTestPaginationQuery(page, pageSize int) *dto.LabTestPaginationQuery {
	return &dto.LabTestPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewCreateLabTestRequest(medicalRecordID, testTypeID, doctorID uint) *dto.CreateLabTestRequest {
	return &dto.CreateLabTestRequest{
		MedicalRecordID:   medicalRecordID,
		TestTypeID:        testTypeID,
		OrderedByDoctorID: doctorID,
		OrderDate:         "2023-12-01",
		Status:            "ordered",
		Notes:             "Test Notes",
	}
}

func NewUpdateLabTestRequest(notes, ref string) *dto.UpdateLabTestRequest {
	return &dto.UpdateLabTestRequest{
		Notes:          &notes,
		ReferenceRange: &ref,
	}
}
