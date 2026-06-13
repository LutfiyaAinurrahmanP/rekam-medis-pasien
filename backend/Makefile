.PHONY: help build run seed test test-users test-departments test-doctor-specializations test-patients test-reset-code test-doctors test-rooms test-type-tests test-unit test-unit-handler test-unit-service test-unit-repository test-unit-verbose test-unit-handler-verbose test-unit-service-verbose test-unit-repository-verbose test-unit-race test-unit-failfast test-unit-timeout test-unit-run test-unit-list test-unit-coverage test-unit-coverage-html test-unit-coverage-handler test-unit-coverage-service test-unit-coverage-repository test-coverage test-ci test-all-cleanup clean clean-coverage docker-build docker-up docker-up-deps docker-down docker-logs docker-clean migrate lint fmt mod-download test-all mod-tidy dev prod-build check-docker wait-db stop-app-container db-refresh

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@go build -o main cmd/api/main.go

run: check-docker docker-up-deps stop-app-container wait-db ## Run local app with auto dependency setup
	@echo "Running application..."
	@go run cmd/api/main.go

seed: ## Seed database with sample data
	@echo "Seeding database..."
	@go run cmd/seed/main.go

test: ## Run all tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter ""; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter ""; \
	fi

test-users: ## Run user API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(TestGetMyProfile_Success|TestGetMyProfile_Unauthorized|TestGetMyProfile_UserNotFound|TestUpdateMyProfile_Success|TestUpdateMyProfile_ValidationError|TestUpdateMyProfile_CannotChangeRole|TestChangeMyPassword_Success|TestChangeMyPassword_InvalidOldPassword|TestDeleteMyAccount_Success|TestDeactivateMyAccount_Success|TestCreateUser_Success|TestCreateUser_DuplicateUsername|TestListUsers_Success|TestDeleteListUsers_Success|TestGetUserByID_Success|TestGetUserByID_NotFound|TestUpdateUser_Success|TestSoftDeleteUser_Success|TestRestoreUser_Success|TestResetPassword_Success|TestActivateUser_Success|TestDeactivateUser_Success|TestHardDeleteUser_Success|TestUserRepository_Create_Success|TestUserRepository_FindByID_Success|TestUserRepository_FindByID_NotFound|TestUserRepository_FindByUsername_Success|TestUserRepository_FindByUsername_NotFound|TestUserRepository_FindByEmail_Success|TestUserRepository_FindByEmail_NotFound|TestUserRepository_FindByPhone_Success|TestUserRepository_FindByPhone_NotFound|TestUserRepository_List_Success|TestUserRepository_List_Empty|TestUserRepository_DeleteList_Success|TestUserRepository_Update_Success|TestUserRepository_Update_NotFound|TestUserRepository_SoftDelete_Success|TestUserRepository_SoftDelete_NotFound|TestUserRepository_HardDelete_Success|TestUserRepository_HardDelete_NotFound|TestUserRepository_Restore_Success|TestUserRepository_Restore_NotFound|TestUserRepository_IsUsernameExists_True|TestUserRepository_IsUsernameExists_False|TestUserRepository_IsEmailExists_True|TestUserRepository_IsEmailExists_False|TestUserRepository_IsPhoneExists_True|TestUserRepository_IsPhoneExists_False|TestUserRepository_FindByUsernameOrEmail_ByUsername|TestUserRepository_FindByUsernameOrEmail_ByEmail|TestUserRepository_FindByUsernameOrEmail_NotFound|TestRequestPasswordReset_Success|TestRequestPasswordReset_EmailNotFound|TestRequestPasswordReset_RedisUnavailable|TestRequestPasswordReset_EmailNormalization|TestResendPasswordResetCode_Success|TestResendPasswordResetCode_EmailNotFound|TestResetPasswordWithToken_Success|TestResetPasswordWithToken_InvalidToken|TestResetPasswordWithToken_UserNotFoundAfterToken|TestResetPasswordWithToken_RedisUnavailable|TestForgotPasswordFlow_FullHappyPath|TestRequestPasswordReset_MultipleCalls|TestResetPasswordWithToken_PasswordIsHashed|TestGetUserByID_Success|TestGetUserByID_UserNotFound|TestListUsers_Success|TestListUsers_EmptyResult|TestCreateUser_Success|TestCreateUser_EmptyPassword|TestCreateUser_UsernameExists|TestUpdateUser_Success|TestSoftDeleteUser_Success|TestSoftDeleteUser_NotFound|TestRestoreUser_Success|TestHardDeleteUser_Success|TestChangePassword_Success|TestChangePassword_InvalidOldPassword|TestResetPassword_Success|TestResetPassword_UserNotFound|TestActivateUser_Success|TestActivateUser_UserNotFound|TestDeactivateUser_Success|TestDeleteListUsers_Success|TestCreateUser_InvalidRole|TestCreateUser_EmailExists|TestCreateUser_PhoneExists|TestCreateUser_DefaultRole|TestVerifyPasswordForDeletion_Success|TestUserModel_IsDoctor|TestUserModel_IsReceptionist|TestUserModel_IsAdmin|TestUserModel_IsSuperAdmin|TestValidateRole_AllValidRoles|TestValidateRole_InvalidRoles|TestGetAvailableRoles)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(TestGetMyProfile_Success|TestGetMyProfile_Unauthorized|TestGetMyProfile_UserNotFound|TestUpdateMyProfile_Success|TestUpdateMyProfile_ValidationError|TestUpdateMyProfile_CannotChangeRole|TestChangeMyPassword_Success|TestChangeMyPassword_InvalidOldPassword|TestDeleteMyAccount_Success|TestDeactivateMyAccount_Success|TestCreateUser_Success|TestCreateUser_DuplicateUsername|TestListUsers_Success|TestDeleteListUsers_Success|TestGetUserByID_Success|TestGetUserByID_NotFound|TestUpdateUser_Success|TestSoftDeleteUser_Success|TestRestoreUser_Success|TestResetPassword_Success|TestActivateUser_Success|TestDeactivateUser_Success|TestHardDeleteUser_Success|TestUserRepository_Create_Success|TestUserRepository_FindByID_Success|TestUserRepository_FindByID_NotFound|TestUserRepository_FindByUsername_Success|TestUserRepository_FindByUsername_NotFound|TestUserRepository_FindByEmail_Success|TestUserRepository_FindByEmail_NotFound|TestUserRepository_FindByPhone_Success|TestUserRepository_FindByPhone_NotFound|TestUserRepository_List_Success|TestUserRepository_List_Empty|TestUserRepository_DeleteList_Success|TestUserRepository_Update_Success|TestUserRepository_Update_NotFound|TestUserRepository_SoftDelete_Success|TestUserRepository_SoftDelete_NotFound|TestUserRepository_HardDelete_Success|TestUserRepository_HardDelete_NotFound|TestUserRepository_Restore_Success|TestUserRepository_Restore_NotFound|TestUserRepository_IsUsernameExists_True|TestUserRepository_IsUsernameExists_False|TestUserRepository_IsEmailExists_True|TestUserRepository_IsEmailExists_False|TestUserRepository_IsPhoneExists_True|TestUserRepository_IsPhoneExists_False|TestUserRepository_FindByUsernameOrEmail_ByUsername|TestUserRepository_FindByUsernameOrEmail_ByEmail|TestUserRepository_FindByUsernameOrEmail_NotFound|TestRequestPasswordReset_Success|TestRequestPasswordReset_EmailNotFound|TestRequestPasswordReset_RedisUnavailable|TestRequestPasswordReset_EmailNormalization|TestResendPasswordResetCode_Success|TestResendPasswordResetCode_EmailNotFound|TestResetPasswordWithToken_Success|TestResetPasswordWithToken_InvalidToken|TestResetPasswordWithToken_UserNotFoundAfterToken|TestResetPasswordWithToken_RedisUnavailable|TestForgotPasswordFlow_FullHappyPath|TestRequestPasswordReset_MultipleCalls|TestResetPasswordWithToken_PasswordIsHashed|TestGetUserByID_Success|TestGetUserByID_UserNotFound|TestListUsers_Success|TestListUsers_EmptyResult|TestCreateUser_Success|TestCreateUser_EmptyPassword|TestCreateUser_UsernameExists|TestUpdateUser_Success|TestSoftDeleteUser_Success|TestSoftDeleteUser_NotFound|TestRestoreUser_Success|TestHardDeleteUser_Success|TestChangePassword_Success|TestChangePassword_InvalidOldPassword|TestResetPassword_Success|TestResetPassword_UserNotFound|TestActivateUser_Success|TestActivateUser_UserNotFound|TestDeactivateUser_Success|TestDeleteListUsers_Success|TestCreateUser_InvalidRole|TestCreateUser_EmailExists|TestCreateUser_PhoneExists|TestCreateUser_DefaultRole|TestVerifyPasswordForDeletion_Success|TestUserModel_IsDoctor|TestUserModel_IsReceptionist|TestUserModel_IsAdmin|TestUserModel_IsSuperAdmin|TestValidateRole_AllValidRoles|TestValidateRole_InvalidRoles|TestGetAvailableRoles)"; \
	fi

test-departments: ## Run department API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(TestListDepartments_Success|TestListDepartments_InvalidPagination|TestListDepartments_WithSearch|TestListDepartments_ServiceError|TestGetDepartmentByID_Success|TestGetDepartmentByID_InvalidID|TestGetDepartmentByID_DepartmentNotFound|TestCreateDepartment_Success|TestCreateDepartment_InvalidInput|TestCreateDepartment_CodeAlreadyExists|TestUpdateDepartment_Success|TestUpdateDepartment_InvalidID|TestUpdateDepartment_DepartmentNotFound|TestUpdateDepartment_CodeAlreadyExists|TestSoftDeleteDepartment_Success|TestSoftDeleteDepartment_InvalidID|TestSoftDeleteDepartment_NotFound|TestRestoreDepartment_Success|TestRestoreDepartment_InvalidID|TestRestoreDepartment_NotFound|TestHardDeleteDepartment_Success|TestHardDeleteDepartment_InvalidID|TestHardDeleteDepartment_NotFound|TestDepartmentRepository_Create_Success|TestDepartmentRepository_FindByID_Success|TestDepartmentRepository_FindByID_NotFound|TestDepartmentRepository_FindByName_Success|TestDepartmentRepository_FindByName_NotFound|TestDepartmentRepository_FindByCode_Success|TestDepartmentRepository_FindByCode_NotFound|TestDepartmentRepository_IsCodeExists_True|TestDepartmentRepository_IsCodeExists_False|TestDepartmentRepository_Update_Success|TestDepartmentRepository_List_Success|TestDepartmentRepository_SoftDelete_Success|TestDepartmentRepository_Restore_Success|TestDepartmentRepository_HardDelete_Success|TestGetDepartmentByID_Success|TestGetDepartmentByID_DepartmentNotFound|TestListDepartments_Success|TestListDepartments_WithSearch|TestListDepartments_Empty|TestListDepartments_DefaultPagination|TestCreateDepartment_Success|TestCreateDepartment_CodeAlreadyExists|TestCreateDepartment_CreateFailed|TestUpdateDepartment_Success|TestUpdateDepartment_DepartmentNotFound|TestUpdateDepartment_ChangeCode|TestUpdateDepartment_CodeAlreadyExists|TestSoftDeleteDepartment_Success|TestSoftDeleteDepartment_DepartmentNotFound|TestRestoreDepartment_Success|TestHardDeleteDepartment_Success|TestDeleteListDepartments_Success)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(TestListDepartments_Success|TestListDepartments_InvalidPagination|TestListDepartments_WithSearch|TestListDepartments_ServiceError|TestGetDepartmentByID_Success|TestGetDepartmentByID_InvalidID|TestGetDepartmentByID_DepartmentNotFound|TestCreateDepartment_Success|TestCreateDepartment_InvalidInput|TestCreateDepartment_CodeAlreadyExists|TestUpdateDepartment_Success|TestUpdateDepartment_InvalidID|TestUpdateDepartment_DepartmentNotFound|TestUpdateDepartment_CodeAlreadyExists|TestSoftDeleteDepartment_Success|TestSoftDeleteDepartment_InvalidID|TestSoftDeleteDepartment_NotFound|TestRestoreDepartment_Success|TestRestoreDepartment_InvalidID|TestRestoreDepartment_NotFound|TestHardDeleteDepartment_Success|TestHardDeleteDepartment_InvalidID|TestHardDeleteDepartment_NotFound|TestDepartmentRepository_Create_Success|TestDepartmentRepository_FindByID_Success|TestDepartmentRepository_FindByID_NotFound|TestDepartmentRepository_FindByName_Success|TestDepartmentRepository_FindByName_NotFound|TestDepartmentRepository_FindByCode_Success|TestDepartmentRepository_FindByCode_NotFound|TestDepartmentRepository_IsCodeExists_True|TestDepartmentRepository_IsCodeExists_False|TestDepartmentRepository_Update_Success|TestDepartmentRepository_List_Success|TestDepartmentRepository_SoftDelete_Success|TestDepartmentRepository_Restore_Success|TestDepartmentRepository_HardDelete_Success|TestGetDepartmentByID_Success|TestGetDepartmentByID_DepartmentNotFound|TestListDepartments_Success|TestListDepartments_WithSearch|TestListDepartments_Empty|TestListDepartments_DefaultPagination|TestCreateDepartment_Success|TestCreateDepartment_CodeAlreadyExists|TestCreateDepartment_CreateFailed|TestUpdateDepartment_Success|TestUpdateDepartment_DepartmentNotFound|TestUpdateDepartment_ChangeCode|TestUpdateDepartment_CodeAlreadyExists|TestSoftDeleteDepartment_Success|TestSoftDeleteDepartment_DepartmentNotFound|TestRestoreDepartment_Success|TestHardDeleteDepartment_Success|TestDeleteListDepartments_Success)"; \
	fi

test-doctor-specializations: ## Run doctor specialization API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(TestListDoctorSpecializations_Success|TestListDoctorSpecializations_InvalidQuery|TestListDoctorSpecializations_ServiceError|TestFindDoctorSpecializationByID_Success|TestFindDoctorSpecializationByID_NotFound|TestCreateDoctorSpecialization_Success|TestUpdateDoctorSpecialization_Success|TestUpdateDoctorSpecialization_Conflict|TestSoftDeleteDoctorSpecialization_Success|TestRestoreDoctorSpecialization_Success|TestHardDeleteDoctorSpecialization_Success|TestDeletedListDoctorSpecializations_Success|TestActiveListDoctorSpecializations_Success|TestInactiveListDoctorSpecializations_Success|TestActivateDoctorSpecialization_Success|TestDeactivateDoctorSpecialization_Success|TestDoctorSpecializationService_FindByID_Success|TestDoctorSpecializationService_FindByID_NotFound|TestDoctorSpecializationService_List_Success|TestDoctorSpecializationService_Create_Success|TestDoctorSpecializationService_Create_NameAlreadyExists|TestDoctorSpecializationService_Update_Success|TestDoctorSpecializationService_SoftDelete_Success|TestDoctorSpecializationService_Restore_Success|TestDoctorSpecializationService_HardDelete_Success|TestDoctorSpecializationService_ActiveList_Success|TestDoctorSpecializationService_InactiveList_Success|TestDoctorSpecializationService_Activate_Success|TestDoctorSpecializationService_Deactivate_Success|TestDoctorSpecializationRepository_Create_Success|TestDoctorSpecializationRepository_FindByID_Success|TestDoctorSpecializationRepository_FindByID_NotFound|TestDoctorSpecializationRepository_IsNameExists_True|TestDoctorSpecializationRepository_IsCodeExists_True|TestDoctorSpecializationRepository_IsCodeExists_False|TestDoctorSpecializationRepository_Update_Success|TestDoctorSpecializationRepository_List_Success|TestDoctorSpecializationRepository_ActiveList_Success|TestDoctorSpecializationRepository_SoftDelete_Success|TestDoctorSpecializationRepository_Restore_Success|TestDoctorSpecializationRepository_HardDelete_Success)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(TestListDoctorSpecializations_Success|TestListDoctorSpecializations_InvalidQuery|TestListDoctorSpecializations_ServiceError|TestFindDoctorSpecializationByID_Success|TestFindDoctorSpecializationByID_NotFound|TestCreateDoctorSpecialization_Success|TestUpdateDoctorSpecialization_Success|TestUpdateDoctorSpecialization_Conflict|TestSoftDeleteDoctorSpecialization_Success|TestRestoreDoctorSpecialization_Success|TestHardDeleteDoctorSpecialization_Success|TestDeletedListDoctorSpecializations_Success|TestActiveListDoctorSpecializations_Success|TestInactiveListDoctorSpecializations_Success|TestActivateDoctorSpecialization_Success|TestDeactivateDoctorSpecialization_Success|TestDoctorSpecializationService_FindByID_Success|TestDoctorSpecializationService_FindByID_NotFound|TestDoctorSpecializationService_List_Success|TestDoctorSpecializationService_Create_Success|TestDoctorSpecializationService_Create_NameAlreadyExists|TestDoctorSpecializationService_Update_Success|TestDoctorSpecializationService_SoftDelete_Success|TestDoctorSpecializationService_Restore_Success|TestDoctorSpecializationService_HardDelete_Success|TestDoctorSpecializationService_ActiveList_Success|TestDoctorSpecializationService_InactiveList_Success|TestDoctorSpecializationService_Activate_Success|TestDoctorSpecializationService_Deactivate_Success|TestDoctorSpecializationRepository_Create_Success|TestDoctorSpecializationRepository_FindByID_Success|TestDoctorSpecializationRepository_FindByID_NotFound|TestDoctorSpecializationRepository_IsNameExists_True|TestDoctorSpecializationRepository_IsCodeExists_True|TestDoctorSpecializationRepository_IsCodeExists_False|TestDoctorSpecializationRepository_Update_Success|TestDoctorSpecializationRepository_List_Success|TestDoctorSpecializationRepository_ActiveList_Success|TestDoctorSpecializationRepository_SoftDelete_Success|TestDoctorSpecializationRepository_Restore_Success|TestDoctorSpecializationRepository_HardDelete_Success)"; \
	fi

test-patients: ## Run patient API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(TestGetMyPatientData_Success|TestGetMyPatientData_Unauthorized|TestUpdateMyPatientData_Success|TestUpdateMyPatientData_NotFound|TestListPatients_Success|TestListPatients_InvalidPagination|TestGetPatientByID_Success|TestGetPatientByID_InvalidID|TestGetPatientByID_OwnershipForbidden|TestGetPatientByCode_Success|TestCreatePatient_Success|TestCreatePatient_DuplicateCode|TestUpdatePatient_Success|TestSoftDeletePatient_Success|TestRestorePatient_Success|TestHardDeletePatient_Success|TestPatientRepository_Create_Success|TestPatientRepository_FindByID_Success|TestPatientRepository_FindByID_NotFound|TestPatientRepository_FindByUserID_Success|TestPatientRepository_FindByUserID_NotFound|TestPatientRepository_FindByCode_Success|TestPatientRepository_FindByCode_NotFound|TestPatientRepository_List_Success|TestPatientRepository_DeleteList_Success|TestPatientRepository_Update_Success|TestPatientRepository_SoftDelete_Success|TestPatientRepository_Restore_Success|TestPatientRepository_HardDelete_Success|TestPatientRepository_IsCodeExists_True|TestPatientRepository_IsCodeExists_False|TestGetPatientByID_Success|TestGetPatientByID_PatientNotFound|TestGetPatientByCode_Success|TestGetPatientByCode_PatientNotFound|TestGetMyPatientData_Success|TestListPatients_Success|TestListPatients_Empty|TestListPatients_DefaultPagination|TestDeleteListPatients_Success|TestCreatePatient_Success|TestCreatePatient_CodeAlreadyExists|TestUpdatePatient_Success|TestUpdatePatient_PatientNotFound|TestUpdateMyPatientData_Success|TestSoftDeletePatient_Success|TestSoftDeletePatient_NotFound|TestRestorePatient_Success|TestHardDeletePatient_Success|TestUserModel_IsPatient)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(TestGetMyPatientData_Success|TestGetMyPatientData_Unauthorized|TestUpdateMyPatientData_Success|TestUpdateMyPatientData_NotFound|TestListPatients_Success|TestListPatients_InvalidPagination|TestGetPatientByID_Success|TestGetPatientByID_InvalidID|TestGetPatientByID_OwnershipForbidden|TestGetPatientByCode_Success|TestCreatePatient_Success|TestCreatePatient_DuplicateCode|TestUpdatePatient_Success|TestSoftDeletePatient_Success|TestRestorePatient_Success|TestHardDeletePatient_Success|TestPatientRepository_Create_Success|TestPatientRepository_FindByID_Success|TestPatientRepository_FindByID_NotFound|TestPatientRepository_FindByUserID_Success|TestPatientRepository_FindByUserID_NotFound|TestPatientRepository_FindByCode_Success|TestPatientRepository_FindByCode_NotFound|TestPatientRepository_List_Success|TestPatientRepository_DeleteList_Success|TestPatientRepository_Update_Success|TestPatientRepository_SoftDelete_Success|TestPatientRepository_Restore_Success|TestPatientRepository_HardDelete_Success|TestPatientRepository_IsCodeExists_True|TestPatientRepository_IsCodeExists_False|TestGetPatientByID_Success|TestGetPatientByID_PatientNotFound|TestGetPatientByCode_Success|TestGetPatientByCode_PatientNotFound|TestGetMyPatientData_Success|TestListPatients_Success|TestListPatients_Empty|TestListPatients_DefaultPagination|TestDeleteListPatients_Success|TestCreatePatient_Success|TestCreatePatient_CodeAlreadyExists|TestUpdatePatient_Success|TestUpdatePatient_PatientNotFound|TestUpdateMyPatientData_Success|TestSoftDeletePatient_Success|TestSoftDeletePatient_NotFound|TestRestorePatient_Success|TestHardDeletePatient_Success|TestUserModel_IsPatient)"; \
	fi

test-reset-code: ## Run reset-code verification tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(TestVerifyResetCode_Success|TestVerifyResetCode_WrongCode|TestVerifyResetCode_ChallengeNotFound|TestVerifyResetCode_ExpiredChallenge|TestVerifyResetCode_RedisUnavailable|TestVerifyResetCode_EmailNormalization|TestVerifyResetCode_CannotReuse)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(TestVerifyResetCode_Success|TestVerifyResetCode_WrongCode|TestVerifyResetCode_ChallengeNotFound|TestVerifyResetCode_ExpiredChallenge|TestVerifyResetCode_RedisUnavailable|TestVerifyResetCode_EmailNormalization|TestVerifyResetCode_CannotReuse)"; \
	fi

test-doctors: ## Run doctor API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(TestDoctorHandler_GetMyDoctorData_Success|TestDoctorHandler_UpdateMyDoctorData_Success|TestDoctorHandler_ListDoctors_Success|TestDoctorHandler_DeletedListDoctors_Success|TestDoctorHandler_ActiveList_Success|TestDoctorHandler_InactiveList_Success|TestDoctorHandler_GetDoctorByID_Success|TestDoctorHandler_CreateDoctor_Success|TestDoctorHandler_CreateDoctor_Conflict|TestDoctorHandler_UpdateDoctor_Success|TestDoctorHandler_ActivateDoctor_Success|TestDoctorHandler_DeactivateDoctor_Success|TestDoctorHandler_SoftDeleteDoctor_Success|TestDoctorHandler_RestoreDoctor_Success|TestDoctorHandler_HardDeleteDoctor_Success|TestDoctorService_GetMyDoctorData_Success|TestDoctorService_GetMyDoctorData_NotFound|TestDoctorService_UpdateMyDoctorData_Success|TestDoctorService_GetDoctorByID_Success|TestDoctorService_GetDoctorByID_NotFound|TestDoctorService_List_Success|TestDoctorService_DeletedList_Success|TestDoctorService_ActiveList_Success|TestDoctorService_InactiveList_Success|TestDoctorService_Create_Success|TestDoctorService_Create_EmployeeIDAlreadyExists|TestDoctorService_Update_Success|TestDoctorService_SoftDelete_Success|TestDoctorService_Restore_Success|TestDoctorService_HardDelete_Success|TestDoctorService_Activate_Success|TestDoctorService_Deactivate_Success|TestDoctorRepository_Create_Success|TestDoctorRepository_FindByID_Success|TestDoctorRepository_FindByID_NotFound|TestDoctorRepository_FindByUserID_Success|TestDoctorRepository_FindByDepartmentID_Success|TestDoctorRepository_IsEmployeeIDExists_True|TestDoctorRepository_IsEmployeeIDExists_False|TestDoctorRepository_Update_Success|TestDoctorRepository_List_Success|TestDoctorRepository_DeleteList_Success|TestDoctorRepository_ActiveList_Success|TestDoctorRepository_InactiveList_Success|TestDoctorRepository_SoftDelete_Success|TestDoctorRepository_Restore_Success|TestDoctorRepository_HardDelete_Success)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(TestDoctorHandler_GetMyDoctorData_Success|TestDoctorHandler_UpdateMyDoctorData_Success|TestDoctorHandler_ListDoctors_Success|TestDoctorHandler_DeletedListDoctors_Success|TestDoctorHandler_ActiveList_Success|TestDoctorHandler_InactiveList_Success|TestDoctorHandler_GetDoctorByID_Success|TestDoctorHandler_CreateDoctor_Success|TestDoctorHandler_CreateDoctor_Conflict|TestDoctorHandler_UpdateDoctor_Success|TestDoctorHandler_ActivateDoctor_Success|TestDoctorHandler_DeactivateDoctor_Success|TestDoctorHandler_SoftDeleteDoctor_Success|TestDoctorHandler_RestoreDoctor_Success|TestDoctorHandler_HardDeleteDoctor_Success|TestDoctorService_GetMyDoctorData_Success|TestDoctorService_GetMyDoctorData_NotFound|TestDoctorService_UpdateMyDoctorData_Success|TestDoctorService_GetDoctorByID_Success|TestDoctorService_GetDoctorByID_NotFound|TestDoctorService_List_Success|TestDoctorService_DeletedList_Success|TestDoctorService_ActiveList_Success|TestDoctorService_InactiveList_Success|TestDoctorService_Create_Success|TestDoctorService_Create_EmployeeIDAlreadyExists|TestDoctorService_Update_Success|TestDoctorService_SoftDelete_Success|TestDoctorService_Restore_Success|TestDoctorService_HardDelete_Success|TestDoctorService_Activate_Success|TestDoctorService_Deactivate_Success|TestDoctorRepository_Create_Success|TestDoctorRepository_FindByID_Success|TestDoctorRepository_FindByID_NotFound|TestDoctorRepository_FindByUserID_Success|TestDoctorRepository_FindByDepartmentID_Success|TestDoctorRepository_IsEmployeeIDExists_True|TestDoctorRepository_IsEmployeeIDExists_False|TestDoctorRepository_Update_Success|TestDoctorRepository_List_Success|TestDoctorRepository_DeleteList_Success|TestDoctorRepository_ActiveList_Success|TestDoctorRepository_InactiveList_Success|TestDoctorRepository_SoftDelete_Success|TestDoctorRepository_Restore_Success|TestDoctorRepository_HardDelete_Success)"; \
	fi

test-rooms: ## Run room API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(Test.*Room|Test.*RoomType)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(Test.*Room|Test.*RoomType)"; \
	fi

test-type-tests: ## Run type test API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(Test.*TypeTest.*|Test.*TypeTestCategory.*)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(Test.*TypeTest.*|Test.*TypeTestCategory.*)"; \
	fi

test-medicines: ## Run medicine API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(Test.*Medicine.*|Test.*MedicineType.*)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(Test.*Medicine.*|Test.*MedicineType.*)"; \
	fi

test-appointments: ## Run appointment API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(Test.*Appointment.*)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(Test.*Appointment.*)"; \
	fi

test-medical-records: ## Run medical record API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(Test.*MedicalRecord.*)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(Test.*MedicalRecord.*)"; \
	fi

test-hospitalizations: ## Run hospitalization API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(Test.*Hospitalization.*)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(Test.*Hospitalization.*)"; \
	fi

test-lab-tests: ## Run lab test API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(Test.*LabTest.*)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(Test.*LabTest.*)"; \
	fi

test-prescriptions: ## Run prescription API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(Test.*Prescription.*)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(Test.*Prescription.*)"; \
	fi

test-vital-signs: ## Run vital signs API tests with Jest-like output
	@if command -v powershell >/dev/null 2>&1; then \
		powershell -File scripts/test-runner.ps1 -Path "./tests/..." -Filter "^(Test.*VitalSign.*)"; \
	else \
		bash ./scripts/test-runner.sh -Path "./tests/..." -Filter "^(Test.*VitalSign.*)"; \
	fi

# ============================================================================
# UNIT TEST COMMANDS (NEW - Comprehensive Unit Testing)
# ============================================================================

test-unit: ## Run all unit tests (handler, service, repository)
	@echo "Running all unit tests..."
	@go test ./tests/unit/... --count=1

test-unit-handler: ## Run only handler layer unit tests
	@echo "Running handler unit tests..."
	@go test ./tests/unit/handler --count=1

test-unit-service: ## Run only service layer unit tests
	@echo "Running service unit tests..."
	@go test ./tests/unit/service --count=1

test-unit-repository: ## Run only repository layer unit tests
	@echo "Running repository unit tests..."
	@go test ./tests/unit/repository --count=1

test-unit-verbose: ## Run all unit tests with verbose output
	@echo "Running all unit tests with verbose output..."
	@go test ./tests/unit/... -v --count=1

test-unit-handler-verbose: ## Run handler unit tests with verbose output
	@echo "Running handler unit tests with verbose output..."
	@go test ./tests/unit/handler -v --count=1

test-unit-service-verbose: ## Run service unit tests with verbose output
	@echo "Running service unit tests with verbose output..."
	@go test ./tests/unit/service -v --count=1

test-unit-repository-verbose: ## Run repository unit tests with verbose output
	@echo "Running repository unit tests with verbose output..."
	@go test ./tests/unit/repository -v --count=1

test-unit-race: ## Run unit tests with race detector
	@echo "Running unit tests with race detector..."
	@go test ./tests/unit/... -race --count=1

test-unit-failfast: ## Run unit tests and stop at first failure
	@echo "Running unit tests in failfast mode..."
	@go test ./tests/unit/... -failfast --count=1

test-unit-timeout: ## Run unit tests with 60s timeout
	@echo "Running unit tests with custom timeout..."
	@go test ./tests/unit/... --timeout=60s --count=1

test-unit-run: ## Run specific unit test (usage: make test-unit-run TEST=TestGetMyProfile)
	@echo "Running specific test: $(TEST)"
	@go test ./tests/unit/... -run $(TEST) -v --count=1

test-unit-list: ## List all available unit tests
	@echo "Available unit tests:"
	@go test ./tests/unit/... -list .

test-unit-coverage: ## Generate unit test coverage report
	@echo "Generating unit test coverage report..."
	@go test ./tests/unit/... -coverprofile=coverage_unit.out --count=1
	@go tool cover -func=coverage_unit.out | tail -30
	@echo "Full coverage report saved to: coverage_unit.out"

test-unit-coverage-html: ## Generate HTML unit test coverage report
	@echo "Generating HTML unit test coverage report..."
	@go test ./tests/unit/... -coverprofile=coverage_unit.out --count=1
	@go tool cover -html=coverage_unit.out -o coverage_unit.html
	@echo "HTML coverage report generated: coverage_unit.html"

test-unit-coverage-handler: ## Coverage report for handler layer only
	@echo "Generating handler layer coverage..."
	@go test ./tests/unit/handler -coverprofile=coverage_handler.out --count=1
	@go tool cover -func=coverage_handler.out

test-unit-coverage-service: ## Coverage report for service layer only
	@echo "Generating service layer coverage..."
	@go test ./tests/unit/service -coverprofile=coverage_service.out --count=1
	@go tool cover -func=coverage_service.out

test-unit-coverage-repository: ## Coverage report for repository layer only
	@echo "Generating repository layer coverage..."
	@go test ./tests/unit/repository -coverprofile=coverage_repository.out --count=1
	@go tool cover -func=coverage_repository.out

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build files
	@echo "Cleaning..."
	@rm -f main
	@go clean

docker-build: ## Build docker images
	@echo "Building docker images..."
	@docker compose build

docker-up: ## Start docker containers
	@echo "Starting docker containers..."
	@docker compose up -d

docker-up-deps: ## Start only dependency containers (without app)
	@echo "Starting dependency containers (postgres, redis, kafka, zookeeper)..."
	@docker compose up -d postgres redis zookeeper kafka kafka-ui

docker-down: ## Stop docker containers
	@echo "Stopping docker containers..."
	@docker compose down

docker-logs: ## View docker logs
	@docker compose logs -f

docker-clean: ## Remove docker containers and volumes
	@echo "Removing docker containers and volumes..."
	@docker compose down -v

# ─── Helper targets ────────────────────────────────────────────────────────

check-docker: ## Check if docker and docker-compose are available
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "❌ Docker is not installed. Please install Docker first."; \
		exit 1; \
	fi
	@if ! docker compose version >/dev/null 2>&1; then \
		echo "❌ Docker Compose is not available. Please install Docker Compose first."; \
		exit 1; \
	fi
	@echo "✅ Docker and Docker Compose are available"

wait-db: ## Wait for database to be ready
	@echo "⏳ Waiting for database to be ready..."
	@for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30; do \
		if nc -z 127.0.0.1 5432 2>/dev/null; then \
			echo "✅ Database is ready!"; \
			sleep 1; \
			exit 0; \
		fi; \
		echo "Attempt $$i/30 - Database not ready yet, waiting..."; \
		sleep 2; \
	done; \
	echo "❌ Database failed to start after 60 seconds"; \
	exit 1

stop-app-container: ## Stop docker app container to avoid port conflicts with local go run
	@docker compose stop app >/dev/null 2>&1 || true
	@docker compose rm -f app >/dev/null 2>&1 || true

migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run cmd/api/main.go

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

mod-download: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

test-all: ## Test all API
	@echo "Test all API"
	@go test -v ./internal/routes/... -count=1

mod-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

test-ci: ## Run tests in CI mode (unit tests with coverage and race detection)
	@echo "Running CI test suite..."
	@go test ./tests/unit/... -race -coverprofile=coverage_unit.out --count=1
	@echo "✅ CI tests passed!"

test-all-cleanup: clean-coverage ## Clean all test artifacts
	@echo "All test artifacts cleaned!"

clean-coverage: ## Clean coverage files
	@echo "Cleaning coverage files..."
	@rm -f coverage.out coverage.html coverage_unit.out coverage_unit.html coverage_handler.out coverage_service.out coverage_repository.out

dev: docker-up-deps ## Start development dependencies only
	@echo "Development environment started!"

prod-build: ## Build for production
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

db-refresh: check-docker docker-up-deps wait-db ## Wipe database, auto-migrate, and re-seed
	@echo "Wiping database schema..."
	@docker compose exec postgres psql -U postgres -d sirekam_medis -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "Flushing Redis cache..."
	@docker compose exec redis redis-cli FLUSHALL || true
	@echo "Running migrations and seeder..."
	@go run cmd/seed/main.go
	@echo "✅ Database successfully refreshed and seeded!"
