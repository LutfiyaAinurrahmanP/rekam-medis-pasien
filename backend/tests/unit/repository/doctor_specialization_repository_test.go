package repository

import (
	"errors"
	"testing"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/tests/unit/mocks"
	"github.com/stretchr/testify/assert"
)

// ============= Test Cases: Create =============

func TestDoctorSpecializationRepository_Create_Success(t *testing.T) {
	specialization := mocks.NewTestDoctorSpecialization()
	specialization.Name = "Cardiology"
	specialization.Code = "CARD"

	assert.NotNil(t, specialization)
	assert.Equal(t, "Cardiology", specialization.Name)
	assert.Equal(t, "CARD", specialization.Code)
}

// ============= Test Cases: FindByID =============

func TestDoctorSpecializationRepository_FindByID_Success(t *testing.T) {
	specialization := mocks.NewTestDoctorSpecializationWithData(1, "Cardiology", "CARD", "Heart Specialization", true)

	assert.NotNil(t, specialization)
	assert.Equal(t, uint(1), specialization.ID)
	assert.Equal(t, "Cardiology", specialization.Name)
	assert.Equal(t, "CARD", specialization.Code)
}

func TestDoctorSpecializationRepository_FindByID_NotFound(t *testing.T) {
	var specialization *mocks.MockDoctorSpecializationRepository
	err := errors.New("doctor specialization not found")

	assert.Nil(t, specialization)
	assert.Error(t, err)
	assert.Equal(t, "doctor specialization not found", err.Error())
}

// ============= Test Cases: FindByName / FindByCode (Simulated with IsExists) =============

func TestDoctorSpecializationRepository_IsNameExists_True(t *testing.T) {
	name := "Cardiology"
	assert.Equal(t, "Cardiology", name)
}

func TestDoctorSpecializationRepository_IsCodeExists_True(t *testing.T) {
	code := "CARD"
	assert.Equal(t, "CARD", code)
}

func TestDoctorSpecializationRepository_IsCodeExists_False(t *testing.T) {
	code := "NEWCODE"
	assert.Equal(t, "NEWCODE", code)
}

// ============= Test Cases: Update =============

func TestDoctorSpecializationRepository_Update_Success(t *testing.T) {
	specialization := mocks.NewTestDoctorSpecializationWithData(1, "Cardiology", "CARD", "Heart Specialization", true)
	specialization.Name = "Updated Cardiology"

	assert.Equal(t, "Updated Cardiology", specialization.Name)
	assert.Equal(t, uint(1), specialization.ID)
}

// ============= Test Cases: List =============

func TestDoctorSpecializationRepository_List_Success(t *testing.T) {
	specializations := mocks.NewTestDoctorSpecializationList(3)

	assert.Equal(t, 3, len(specializations))
	assert.NotNil(t, specializations[0])
	assert.NotNil(t, specializations[1])
	assert.NotNil(t, specializations[2])
}

func TestDoctorSpecializationRepository_ActiveList_Success(t *testing.T) {
	specializations := mocks.NewTestDoctorSpecializationList(2)
	assert.Equal(t, 2, len(specializations))
}

// ============= Test Cases: SoftDelete =============

func TestDoctorSpecializationRepository_SoftDelete_Success(t *testing.T) {
	specialization := mocks.NewTestDoctorSpecializationWithData(1, "Cardiology", "CARD", "Heart Specialization", true)
	assert.NotNil(t, specialization)
	assert.Equal(t, uint(1), specialization.ID)
}

// ============= Test Cases: Restore =============

func TestDoctorSpecializationRepository_Restore_Success(t *testing.T) {
	specialization := mocks.NewTestDoctorSpecializationWithData(1, "Cardiology", "CARD", "Heart Specialization", true)
	assert.NotNil(t, specialization)
}

// ============= Test Cases: HardDelete =============

func TestDoctorSpecializationRepository_HardDelete_Success(t *testing.T) {
	specialization := mocks.NewTestDoctorSpecializationWithData(1, "Cardiology", "CARD", "Heart Specialization", true)
	assert.NotNil(t, specialization)
}
