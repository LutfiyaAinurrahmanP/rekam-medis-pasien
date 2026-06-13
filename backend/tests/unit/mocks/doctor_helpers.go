package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

// ============= Doctor Model Builders =============

func NewTestDoctorWithData(id uint, userID uint, employeeID, fullName, licenseNumber, phone, email string, deptID, specID uint, isActive bool) *models.Doctor {
	return &models.Doctor{
		ID:               id,
		UserID:           &userID,
		EmployeeID:       employeeID,
		FullName:         fullName,
		SpecializationID: specID,
		LicenseNumber:    licenseNumber,
		Phone:            phone,
		Email:            email,
		DepartmentID:     &deptID,
		IsActive:         isActive,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

func NewTestDoctorList(count int) []models.Doctor {
	doctors := make([]models.Doctor, count)
	for i := 0; i < count; i++ {
		doctors[i] = *NewTestDoctorWithData(
			uint(i+1),
			uint(i+100),
			"DOC00"+string(rune(i+1)),
			"Doctor "+string(rune(i+1)),
			"LIC00"+string(rune(i+1)),
			"0812345678"+string(rune(i+1)),
			"doc"+string(rune(i+1))+"@hospital.com",
			uint(1),
			uint(1),
			true,
		)
	}
	return doctors
}

// ============= Doctor DTO Builders =============

func NewTestDoctorResponseWithData(id uint, userID uint, employeeID, fullName, licenseNumber, phone, email string, deptID, specID uint, isActive bool) *dto.DoctorResponse {
	return &dto.DoctorResponse{
		ID:               id,
		UserID:           &userID,
		EmployeeID:       employeeID,
		FullName:         fullName,
		LicenseNumber:    licenseNumber,
		Phone:            phone,
		Email:                  email,
		DepartmentID:           &deptID,
		SpecializationID:       specID,
		IsActive:               &isActive,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
}

func NewDoctorPaginationQuery(page, pageSize int) *dto.DoctorPaginationQuery {
	return &dto.DoctorPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

func NewCreateDoctorRequest(userID uint, employeeID, fullName, licenseNumber, phone, email string, deptID, specID uint, isActive bool) *dto.CreateDoctorRequest {
	return &dto.CreateDoctorRequest{
		UserID:           &userID,
		EmployeeID:       employeeID,
		FullName:         fullName,
		SpecializationID: &specID,
		LicenseNumber:    licenseNumber,
		Phone:            phone,
		Email:            email,
		DepartmentID:     &deptID,
		IsActive:         &isActive,
	}
}

func NewUpdateDoctorRequest(fullName, phone, email string, deptID, specID uint, isActive bool) *dto.UpdateDoctorRequest {
	return &dto.UpdateDoctorRequest{
		FullName:         PtrString(fullName),
		Phone:            PtrString(phone),
		Email:            PtrString(email),
		DepartmentID:     PtrUint(deptID),
		SpecializationID: PtrUint(specID),
		IsActive:         PtrBool(isActive),
	}
}
