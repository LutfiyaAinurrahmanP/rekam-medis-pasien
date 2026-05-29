package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
)

// ============= User Model Builders =============

// NewTestUser creates a test user with default values
func NewTestUser() *models.User {
	return &models.User{
		ID:        1,
		Username:  "test_user",
		Email:     "test@example.com",
		Phone:     "08123456789",
		Password:  "hashed_password",
		Role:      models.RolePatient,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewTestUserWithData creates a test user with custom data
func NewTestUserWithData(id uint, username, email, phone, role string, isActive bool) *models.User {
	return &models.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Phone:     phone,
		Password:  "hashed_password",
		Role:      role,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ============= DTO Builders =============

// NewTestUserResponse creates a test user response with default values
func NewTestUserResponse() *dto.UserResponse {
	return &dto.UserResponse{
		ID:        1,
		Username:  "test_user",
		Email:     "test@example.com",
		Phone:     "08123456789",
		Role:      models.RolePatient,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewTestUserResponseWithData creates a test user response with custom data
func NewTestUserResponseWithData(id uint, username, email, phone, role string, isActive bool) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        id,
		Username:  username,
		Email:     email,
		Phone:     phone,
		Role:      role,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewRegisterRequest creates a test register request
func NewRegisterRequest(username, email, phone, password string) *dto.RegisterRequest {
	return &dto.RegisterRequest{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: password,
		Role:     models.RolePatient,
	}
}

// NewCreateUserRequest creates a test create user request
func NewCreateUserRequest(username, email, phone, password, role string, isActive bool) *dto.CreateUserRequest {
	return &dto.CreateUserRequest{
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: password,
		Role:     role,
		IsActive: &isActive,
	}
}

// NewLoginRequest creates a test login request
func NewLoginRequest(usernameOrEmail, password string) *dto.LoginRequest {
	return &dto.LoginRequest{
		UsernameOrEmail: usernameOrEmail,
		Password:        password,
	}
}

// NewUpdateUserRequest creates a test update user request
func NewUpdateUserRequest() *dto.UpdateUserRequest {
	username := "updated_user"
	email := "updated@example.com"
	return &dto.UpdateUserRequest{
		Username: &username,
		Email:    &email,
	}
}

// NewChangePasswordRequest creates a test change password request
func NewChangePasswordRequest(oldPass, newPass string) *dto.ChangePasswordRequest {
	return &dto.ChangePasswordRequest{
		OldPassword: oldPass,
		NewPassword: newPass,
	}
}

// NewUserPaginationQuery creates a test pagination query
func NewUserPaginationQuery(page, pageSize int) *dto.UserPaginationQuery {
	return &dto.UserPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

// ============= Pointer Helpers =============

// PtrString returns a pointer to a string
func PtrString(s string) *string {
	return &s
}

// PtrBool returns a pointer to a bool
func PtrBool(b bool) *bool {
	return &b
}

// PtrUint returns a pointer to an uint
func PtrUint(u uint) *uint {
	return &u
}

// ============= Test List Builders =============

// NewTestUserList creates a list of test users
func NewTestUserList(count int) []models.User {
	users := make([]models.User, count)
	for i := 0; i < count; i++ {
		users[i] = *NewTestUserWithData(
			uint(i+1),
			"user_"+string(rune(i+1)),
			"user"+string(rune(i+1))+"@example.com",
			"0812345678"+string(rune(i+1)),
			models.RolePatient,
			true,
		)
	}
	return users
}

// NewTestUserResponseList creates a list of test user responses
func NewTestUserResponseList(count int) []dto.UserResponse {
	responses := make([]dto.UserResponse, count)
	for i := 0; i < count; i++ {
		responses[i] = *NewTestUserResponseWithData(
			uint(i+1),
			"user_"+string(rune(i+1)),
			"user"+string(rune(i+1))+"@example.com",
			"0812345678"+string(rune(i+1)),
			models.RolePatient,
			true,
		)
	}
	return responses
}

// ============= Department Model Builders =============

// NewTestDepartment creates a test department with default values
func NewTestDepartment() *models.Department {
	return &models.Department{
		ID:            1,
		Name:          "General Medicine",
		Code:          "GM001",
		Description:   "General Medicine Department",
		FloorLocation: "Floor 1",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// NewTestDepartmentWithData creates a test department with custom data
func NewTestDepartmentWithData(id uint, name, code, description, floorLocation string) *models.Department {
	return &models.Department{
		ID:            id,
		Name:          name,
		Code:          code,
		Description:   description,
		FloorLocation: floorLocation,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// ============= Department DTO Builders =============

// NewTestDepartmentResponse creates a test department response with default values
func NewTestDepartmentResponse() *dto.DepartmentResponse {
	return &dto.DepartmentResponse{
		ID:            1,
		Name:          "General Medicine",
		Code:          "GM001",
		Description:   "General Medicine Department",
		FloorLocation: "Floor 1",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// NewTestDepartmentResponseWithData creates a test department response with custom data
func NewTestDepartmentResponseWithData(id uint, name, code, description, floorLocation string) *dto.DepartmentResponse {
	return &dto.DepartmentResponse{
		ID:            id,
		Name:          name,
		Code:          code,
		Description:   description,
		FloorLocation: floorLocation,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// NewCreateDepartmentRequest creates a test create department request
func NewCreateDepartmentRequest(name, code, description, floorLocation string) *dto.CreateDepartmentRequest {
	return &dto.CreateDepartmentRequest{
		Name:          name,
		Code:          code,
		Description:   description,
		FloorLocation: floorLocation,
	}
}

// NewUpdateDepartmentRequest creates a test update department request
func NewUpdateDepartmentRequest(name, code, description, floorLocation string) *dto.UpdateDepartmentRequest {
	return &dto.UpdateDepartmentRequest{
		Name:          PtrString(name),
		Code:          PtrString(code),
		Description:   PtrString(description),
		FloorLocation: PtrString(floorLocation),
	}
}

// NewDepartmentPaginationQuery creates a test pagination query for departments
func NewDepartmentPaginationQuery(page, pageSize int) *dto.DepartmentPaginationQuery {
	return &dto.DepartmentPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

// ============= Department Test List Builders =============

// NewTestDepartmentList creates a list of test departments
func NewTestDepartmentList(count int) []models.Department {
	departments := make([]models.Department, count)
	for i := 0; i < count; i++ {
		departments[i] = *NewTestDepartmentWithData(
			uint(i+1),
			"Department "+string(rune(i+1)),
			"DEP"+string(rune(i+1)),
			"Description "+string(rune(i+1)),
			"Floor "+string(rune(i+1)),
		)
	}
	return departments
}

// NewTestDepartmentResponseList creates a list of test department responses
func NewTestDepartmentResponseList(count int) []dto.DepartmentResponse {
	responses := make([]dto.DepartmentResponse, count)
	for i := 0; i < count; i++ {
		responses[i] = *NewTestDepartmentResponseWithData(
			uint(i+1),
			"Department "+string(rune(i+1)),
			"DEP"+string(rune(i+1)),
			"Description "+string(rune(i+1)),
			"Floor "+string(rune(i+1)),
		)
	}
	return responses
}

// ============= Patient Model Builders =============

// NewTestPatient creates a test patient with default values
func NewTestPatient() *models.Patient {
	userID := uint(1)
	return &models.Patient{
		ID:                    1,
		UserID:                &userID,
		PatientCode:           "PAT001",
		FullName:              "John Doe",
		DateOfBirth:           "1990-01-01",
		Gender:                "male",
		BloodType:             "O",
		Phone:                 "08123456789",
		Email:                 "john@example.com",
		Address:               "Jakarta",
		EmergencyContactName:  "Jane Doe",
		EmergencyContactPhone: "08129876543",
		InsuranceNumber:       "INS-001",
		InsuranceProvider:     "BPJS",
		Allergies:             "None",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}

// NewTestPatientWithData creates a test patient with custom data
func NewTestPatientWithData(id uint, patientCode, fullName, dateOfBirth, gender, bloodType string) *models.Patient {
	userID := id
	return &models.Patient{
		ID:                    id,
		UserID:                &userID,
		PatientCode:           patientCode,
		FullName:              fullName,
		DateOfBirth:           dateOfBirth,
		Gender:                gender,
		BloodType:             bloodType,
		Phone:                 "08123456789",
		Email:                 "patient@example.com",
		Address:               "Jakarta",
		EmergencyContactName:  "Emergency Contact",
		EmergencyContactPhone: "08129876543",
		InsuranceNumber:       "INS-001",
		InsuranceProvider:     "BPJS",
		Allergies:             "None",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}

// ============= Patient DTO Builders =============

// NewTestPatientResponse creates a test patient response with default values
func NewTestPatientResponse() *dto.PatientResponse {
	userID := uint(1)
	return &dto.PatientResponse{
		ID:                    1,
		UserID:                &userID,
		PatientCode:           "PAT001",
		FullName:              "John Doe",
		DateOfBirth:           "1990-01-01",
		Age:                   35,
		Gender:                "male",
		BloodType:             "O",
		Phone:                 "08123456789",
		Email:                 "john@example.com",
		Address:               "Jakarta",
		EmergencyContactName:  "Jane Doe",
		EmergencyContactPhone: "08129876543",
		InsuranceNumber:       "INS-001",
		InsuranceProvider:     "BPJS",
		Allergies:             "None",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}

// NewTestPatientResponseWithData creates a test patient response with custom data
func NewTestPatientResponseWithData(id uint, patientCode, fullName, dateOfBirth, gender, bloodType string) *dto.PatientResponse {
	userID := id
	return &dto.PatientResponse{
		ID:                    id,
		UserID:                &userID,
		PatientCode:           patientCode,
		FullName:              fullName,
		DateOfBirth:           dateOfBirth,
		Age:                   35,
		Gender:                gender,
		BloodType:             bloodType,
		Phone:                 "08123456789",
		Email:                 "patient@example.com",
		Address:               "Jakarta",
		EmergencyContactName:  "Emergency Contact",
		EmergencyContactPhone: "08129876543",
		InsuranceNumber:       "INS-001",
		InsuranceProvider:     "BPJS",
		Allergies:             "None",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}

// NewDeletedPatientResponse creates a test deleted patient response with default values
func NewDeletedPatientResponse() *dto.DeletedPatientResponse {
	patient := NewTestPatientResponse()
	deletedAt := time.Now()
	return &dto.DeletedPatientResponse{
		ID:                    patient.ID,
		UserID:                patient.UserID,
		PatientCode:           patient.PatientCode,
		FullName:              patient.FullName,
		DateOfBirth:           patient.DateOfBirth,
		Age:                   patient.Age,
		Gender:                patient.Gender,
		BloodType:             patient.BloodType,
		Phone:                 patient.Phone,
		Email:                 patient.Email,
		Address:               patient.Address,
		EmergencyContactName:  patient.EmergencyContactName,
		EmergencyContactPhone: patient.EmergencyContactPhone,
		InsuranceNumber:       patient.InsuranceNumber,
		InsuranceProvider:     patient.InsuranceProvider,
		Allergies:             patient.Allergies,
		CreatedAt:             patient.CreatedAt,
		UpdatedAt:             patient.UpdatedAt,
		DeletedAt:             &deletedAt,
	}
}

// NewPatientPaginationQuery creates a test pagination query for patients
func NewPatientPaginationQuery(page, pageSize int) *dto.PatientPaginationQuery {
	return &dto.PatientPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "desc",
	}
}

// NewTestPatientList creates a list of test patients
func NewTestPatientList(count int) []models.Patient {
	patients := make([]models.Patient, count)
	for i := 0; i < count; i++ {
		patients[i] = *NewTestPatientWithData(
			uint(i+1),
			"PAT00"+string(rune(i+1)),
			"Patient "+string(rune(i+1)),
			"1990-01-01",
			"male",
			"O",
		)
	}
	return patients
}

// NewTestDeletedPatientList creates a list of test deleted patients
func NewTestDeletedPatientList(count int) []models.Patient {
	patients := make([]models.Patient, count)
	for i := 0; i < count; i++ {
		patients[i] = *NewTestPatientWithData(
			uint(i+1),
			"PATDEL00"+string(rune(i+1)),
			"Deleted Patient "+string(rune(i+1)),
			"1990-01-01",
			"female",
			"A",
		)
	}
	return patients
}
