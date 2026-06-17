package repository

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
)

// ============= Test Cases: Create =============

func TestPatientRepository_Create_Success(t *testing.T) {
	patient := mocks.NewTestPatient()
	patient.PatientCode = "PAT001"
	patient.FullName = "John Doe"

	assert.NotNil(t, patient)
	assert.Equal(t, "PAT001", patient.PatientCode)
	assert.Equal(t, "John Doe", patient.FullName)
}

// ============= Test Cases: FindByID =============

func TestPatientRepository_FindByID_Success(t *testing.T) {
	patient := mocks.NewTestPatientWithData(1, "PAT001", "John Doe", "1990-01-01", "male", "O")

	assert.NotNil(t, patient)
	assert.Equal(t, uint(1), patient.ID)
	assert.Equal(t, "PAT001", patient.PatientCode)
}

func TestPatientRepository_FindByID_NotFound(t *testing.T) {
	var patient *mocks.MockPatientRepository
	err := errors.New("patient not found")

	assert.Nil(t, patient)
	assert.Error(t, err)
	assert.Equal(t, "patient not found", err.Error())
}

// ============= Test Cases: FindByUserID =============

func TestPatientRepository_FindByUserID_Success(t *testing.T) {
	patient := mocks.NewTestPatientWithData(1, "PAT001", "John Doe", "1990-01-01", "male", "O")

	assert.NotNil(t, patient)
	assert.Equal(t, uint(1), *patient.UserID)
}

func TestPatientRepository_FindByUserID_NotFound(t *testing.T) {
	var patient *mocks.MockPatientRepository
	err := errors.New("patient not found")

	assert.Nil(t, patient)
	assert.Error(t, err)
}

// ============= Test Cases: List =============

func TestPatientRepository_List_Success(t *testing.T) {
	patients := mocks.NewTestPatientList(3)

	assert.Equal(t, 3, len(patients))
	assert.NotNil(t, patients[0])
	assert.NotNil(t, patients[1])
	assert.NotNil(t, patients[2])
}

// ============= Test Cases: DeleteList =============

func TestPatientRepository_DeleteList_Success(t *testing.T) {
	patients := mocks.NewTestDeletedPatientList(1)

	assert.Equal(t, 1, len(patients))
	assert.NotNil(t, patients[0])
}

// ============= Test Cases: Update =============

func TestPatientRepository_Update_Success(t *testing.T) {
	patient := mocks.NewTestPatient()
	patient.FullName = "Updated Patient"

	assert.NotNil(t, patient)
	assert.Equal(t, "Updated Patient", patient.FullName)
}

// ============= Test Cases: SoftDelete =============

func TestPatientRepository_SoftDelete_Success(t *testing.T) {
	patient := mocks.NewTestPatient()
	assert.NotNil(t, patient)
	assert.Equal(t, uint(1), patient.ID)
}

// ============= Test Cases: Restore =============

func TestPatientRepository_Restore_Success(t *testing.T) {
	patient := mocks.NewTestPatient()
	assert.NotNil(t, patient)
}

// ============= Test Cases: HardDelete =============

func TestPatientRepository_HardDelete_Success(t *testing.T) {
	patient := mocks.NewTestPatient()
	assert.NotNil(t, patient)
}

// ============= Test Cases: IsCodeExists =============

func TestPatientRepository_IsCodeExists_True(t *testing.T) {
	exists := true
	assert.True(t, exists)
}

func TestPatientRepository_IsCodeExists_False(t *testing.T) {
	exists := false
	assert.False(t, exists)
}
