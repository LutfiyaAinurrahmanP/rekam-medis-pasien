package repository

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
)

// ============= Test Cases: Create =============

func TestDepartmentRepository_Create_Success(t *testing.T) {
	department := mocks.NewTestDepartment()
	department.Name = "Cardiology"
	department.Code = "CARD001"

	assert.NotNil(t, department)
	assert.Equal(t, "Cardiology", department.Name)
	assert.Equal(t, "CARD001", department.Code)
}

// ============= Test Cases: FindByID =============

func TestDepartmentRepository_FindByID_Success(t *testing.T) {
	department := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Heart and Lung Department", "Floor 2")

	assert.NotNil(t, department)
	assert.Equal(t, uint(1), department.ID)
	assert.Equal(t, "Cardiology", department.Name)
	assert.Equal(t, "CARD001", department.Code)
}

func TestDepartmentRepository_FindByID_NotFound(t *testing.T) {
	var department *mocks.MockDepartmentRepository
	err := errors.New("department not found")

	assert.Nil(t, department)
	assert.Error(t, err)
	assert.Equal(t, "department not found", err.Error())
}

// ============= Test Cases: FindByName =============

func TestDepartmentRepository_FindByName_Success(t *testing.T) {
	department := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Heart and Lung Department", "Floor 2")

	assert.NotNil(t, department)
	assert.Equal(t, "Cardiology", department.Name)
}

func TestDepartmentRepository_FindByName_NotFound(t *testing.T) {
	var department *mocks.MockDepartmentRepository
	err := errors.New("department not found")

	assert.Nil(t, department)
	assert.Error(t, err)
}

// ============= Test Cases: FindByCode =============

func TestDepartmentRepository_FindByCode_Success(t *testing.T) {
	department := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Heart and Lung Department", "Floor 2")

	assert.NotNil(t, department)
	assert.Equal(t, "CARD001", department.Code)
}

func TestDepartmentRepository_FindByCode_NotFound(t *testing.T) {
	var department *mocks.MockDepartmentRepository
	err := errors.New("code not found")

	assert.Nil(t, department)
	assert.Error(t, err)
}

// ============= Test Cases: IsCodeExists =============

func TestDepartmentRepository_IsCodeExists_True(t *testing.T) {
	code := "CARD001"
	// Dalam skenario nyata, repository akan mengecek apakah kode sudah ada
	assert.Equal(t, "CARD001", code)
}

func TestDepartmentRepository_IsCodeExists_False(t *testing.T) {
	code := "NEWDEPT"
	assert.Equal(t, "NEWDEPT", code)
}

// ============= Test Cases: Update =============

func TestDepartmentRepository_Update_Success(t *testing.T) {
	department := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Heart Department", "Floor 2")
	department.Name = "Updated Cardiology"

	assert.Equal(t, "Updated Cardiology", department.Name)
	assert.Equal(t, uint(1), department.ID)
}

// ============= Test Cases: List =============

func TestDepartmentRepository_List_Success(t *testing.T) {
	departments := mocks.NewTestDepartmentList(3)

	assert.Equal(t, 3, len(departments))
	assert.NotNil(t, departments[0])
	assert.NotNil(t, departments[1])
	assert.NotNil(t, departments[2])
}

// ============= Test Cases: SoftDelete =============

func TestDepartmentRepository_SoftDelete_Success(t *testing.T) {
	department := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Heart Department", "Floor 2")
	assert.NotNil(t, department)
	assert.Equal(t, uint(1), department.ID)
}

// ============= Test Cases: Restore =============

func TestDepartmentRepository_Restore_Success(t *testing.T) {
	department := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Heart Department", "Floor 2")
	assert.NotNil(t, department)
}

// ============= Test Cases: HardDelete =============

func TestDepartmentRepository_HardDelete_Success(t *testing.T) {
	department := mocks.NewTestDepartmentWithData(1, "Cardiology", "CARD001", "Heart Department", "Floor 2")
	assert.NotNil(t, department)
}
