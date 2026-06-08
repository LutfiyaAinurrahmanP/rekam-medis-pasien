package cache

import "fmt"

// Konvensi penamaan key: {domain}:{qualifier}:{identifier}
// Contoh: user:id:42, room:available:p1:s10

// ─── User ──────────────────────────────────────────────────────────────────

func UserKey(id uint) string            { return fmt.Sprintf("user:id:%d", id) }
func UserListKey(page, size int) string { return fmt.Sprintf("user:list:p%d:s%d", page, size) }
func UserListQueryKey(page, size int, search, role string, isActive string, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"user:list:p%d:s%d:q%s:r%s:a%s:sb%s:sd%s",
		page,
		size,
		search,
		role,
		isActive,
		sortBy,
		sortDir,
	)
}
func UserDeletedListKey(page, size int) string {
	return fmt.Sprintf("user:deleted:p%d:s%d", page, size)
}
func UserDeletedListQueryKey(page, size int, search, role string, isActive string, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"user:deleted:p%d:s%d:q%s:r%s:a%s:sb%s:sd%s",
		page,
		size,
		search,
		role,
		isActive,
		sortBy,
		sortDir,
	)
}

func PasswordResetChallengeKey(email string) string {
	return fmt.Sprintf("auth:password-reset:challenge:%s", email)
}

func PasswordResetTokenKey(token string) string {
	return fmt.Sprintf("auth:password-reset:token:%s", token)
}

// ─── Department ────────────────────────────────────────────────────────────

func DepartmentKey(id uint) string { return fmt.Sprintf("department:id:%d", id) }
func DepartmentListKey(page, size int) string {
	return fmt.Sprintf("department:list:p%d:s%d", page, size)
}
func DepartmentListQueryKey(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"department:list:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func DepartmentDeletedListKey(page, size int) string {
	return fmt.Sprintf("department:deleted:p%d:s%d", page, size)
}
func DepartmentDeletedListQueryKey(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"department:deleted:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

// ─── Patient ───────────────────────────────────────────────────────────────

func PatientKey(id uint) string            { return fmt.Sprintf("patient:id:%d", id) }
func PatientListKey(page, size int) string { return fmt.Sprintf("patient:list:p%d:s%d", page, size) }
func PatientListQueryKey(page, size int, gender, bloodType, insuranceProvider string, minAge, maxAge int, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"patient:list:p%d:s%d:g%s:b%s:ip%s:min%d:max%d:sb%s:sd%s",
		page,
		size,
		gender,
		bloodType,
		insuranceProvider,
		minAge,
		maxAge,
		sortBy,
		sortDir,
	)
}
func PatientDeletedListKey(page, size int) string {
	return fmt.Sprintf("patient:deleted:p%d:s%d", page, size)
}
func PatientDeletedListQueryKey(page, size int, gender, bloodType, insuranceProvider string, minAge, maxAge int, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"patient:deleted:p%d:s%d:g%s:b%s:ip%s:min%d:max%d:sb%s:sd%s",
		page,
		size,
		gender,
		bloodType,
		insuranceProvider,
		minAge,
		maxAge,
		sortBy,
		sortDir,
	)
}
func PatientByCodeKey(code string) string   { return fmt.Sprintf("patient:code:%s", code) }
func PatientByUserIDKey(userID uint) string { return fmt.Sprintf("patient:user:%d", userID) }

// ─── Doctor Specialization ─────────────────────────────────────────────────

func DoctorSpecializationKey(id uint) string {
	return fmt.Sprintf("doctor_specialization:id:%d", id)
}

func DoctorSpecializationListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"doctor_specialization:list:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func DoctorSpecializationDeletedListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"doctor_specialization:deleted:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func DoctorSpecializationActiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"doctor_specialization:active:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func DoctorSpecializationInactiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"doctor_specialization:inactive:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

// ─── Doctor ────────────────────────────────────────────────────────────────

func DoctorKey(id uint) string            { return fmt.Sprintf("doctor:id:%d", id) }
func DoctorListKey(page, size int) string { return fmt.Sprintf("doctor:list:p%d:s%d", page, size) }
func DoctorListQueryKey(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"doctor:list:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func DoctorActiveListQueryKey(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"doctor:active:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func DoctorInactiveListQueryKey(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"doctor:inactive:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}
func DoctorByUserIDKey(userID uint) string { return fmt.Sprintf("doctor:user:%d", userID) }

// ─── Room Types ────────────────────────────────────────────────────────────

func RoomTypeKey(id uint) string {
	return fmt.Sprintf("room_type:id:%d", id)
}

func RoomTypeListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"room_type:list:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func RoomTypeDeletedListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"room_type:deleted:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func RoomTypeActiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"room_type:active:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func RoomTypeInactiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"room_type:inactive:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

// ─── Room ──────────────────────────────────────────────────────────────────

func RoomKey(id uint) string { return fmt.Sprintf("room:id:%d", id) }

func RoomListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"room:list:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func RoomDeletedListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"room:deleted:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func RoomActiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"room:active:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func RoomInactiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"room:inactive:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func RoomAvailableListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"room:available:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func RoomOccupiedListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"room:occupied:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

// ─── Type Test Category ─────────────────────────────────────────────────

func TypeTestCategoryKey(id uint) string {
	return fmt.Sprintf("type_test_category:id:%d", id)
}

func TypeTestCategoryListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"type_test_category:list:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func TypeTestCategoryDeletedListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"type_test_category:deleted:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func TypeTestCategoryActiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"type_test_category:active:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func TypeTestCategoryInactiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"type_test_category:inactive:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

// ─── TypeTest ──────────────────────────────────────────────────────────────

func TypeTestKey(id uint) string {
	return fmt.Sprintf("type_test:id:%d", id)
}

func TypeTestListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"type_test:list:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func TypeTestDeletedListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"type_test:deleted:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func TypeTestActiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"type_test:active:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func TypeTestInactiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"type_test:inactive:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

// ─── Medicine Type ──────────────────────────────────────────────────────────────

func MedicineTypeKey(id uint) string {
	return fmt.Sprintf("medicine_type:id:%d", id)
}

func MedicineTypeListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"medicine_type:list:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func MedicineTypeDeletedListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"medicine_type:deleted:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func MedicineTypeActiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"medicine_type:active:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func MedicineTypeInactiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"medicine_type:inactive:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

// ─── Medicine ──────────────────────────────────────────────────────────────

func MedicineKey(id uint) string {
	return fmt.Sprintf("medicine:id:%d", id)
}

func MedicineListQuery(page, size int, search string, isActive *bool, typeId *uint, stockStatus, sortBy, sortDir string) string {
	// Konversi pointer ke string agar aman digabungkan
	activeStr := "all"
	if isActive != nil {
		activeStr = fmt.Sprintf("%t", *isActive)
	}

	typeStr := "all"
	if typeId != nil {
		typeStr = fmt.Sprintf("%d", *typeId)
	}

	if stockStatus == "" {
		stockStatus = "all"
	}
	if search == "" {
		search = "none"
	}

	return fmt.Sprintf(
		"medicine:list:p%d:s%d:q%s:act%s:type%s:stk%s:sb%s:sd%s",
		page, size, search, activeStr, typeStr, stockStatus, sortBy, sortDir,
	)
}

func MedicineDeletedListQuery(page, size int, search, sortBy, sortDir string) string {
	if search == "" {
		search = "none"
	}
	return fmt.Sprintf(
		"medicine:deleted:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func MedicineActiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"medicine:active:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

func MedicineInactiveListQuery(page, size int, search, sortBy, sortDir string) string {
	return fmt.Sprintf(
		"medicine:inactive:p%d:s%d:q%s:sb%s:sd%s",
		page,
		size,
		search,
		sortBy,
		sortDir,
	)
}

// ─── Appointment ──────────────────────────────────────────────────────────

func AppointmentKey(id uint) string {
	return fmt.Sprintf("appointment:id:%d", id)
}

func AppointmentListQueryKey(page, size int, patientID, doctorID, departmentID *uint, status, date, dateFrom, dateTo string, daysAhead, daysBack int, sortBy, sortDir string) string {
	pID, dID, depID := "all", "all", "all"
	if patientID != nil {
		pID = fmt.Sprintf("%d", *patientID)
	}
	if doctorID != nil {
		dID = fmt.Sprintf("%d", *doctorID)
	}
	if departmentID != nil {
		depID = fmt.Sprintf("%d", *departmentID)
	}
	if status == "" {
		status = "all"
	}
	return fmt.Sprintf(
		"appointment:list:p%d:s%d:pat%s:doc%s:dep%s:st%s:d%s:df%s:dt%s:da%d:db%d:sb%s:sd%s",
		page, size, pID, dID, depID, status, date, dateFrom, dateTo, daysAhead, daysBack, sortBy, sortDir,
	)
}

func AppointmentDeletedListQueryKey(page, size int, patientID, doctorID, departmentID *uint, status, date, dateFrom, dateTo string, daysAhead, daysBack int, sortBy, sortDir string) string {
	pID, dID, depID := "all", "all", "all"
	if patientID != nil {
		pID = fmt.Sprintf("%d", *patientID)
	}
	if doctorID != nil {
		dID = fmt.Sprintf("%d", *doctorID)
	}
	if departmentID != nil {
		depID = fmt.Sprintf("%d", *departmentID)
	}
	if status == "" {
		status = "all"
	}
	return fmt.Sprintf(
		"appointment:deleted:p%d:s%d:pat%s:doc%s:dep%s:st%s:d%s:df%s:dt%s:da%d:db%d:sb%s:sd%s",
		page, size, pID, dID, depID, status, date, dateFrom, dateTo, daysAhead, daysBack, sortBy, sortDir,
	)
}

// ─── Invalidation patterns ─────────────────────────────────────────────────

const (
	PatternUserAll                 = "user:*"
	PatternDepartmentAll           = "department:*"
	PatternPatientAll              = "patient:*"
	PatternDoctorSpecializationAll = "doctor_specialization:*"
	PatternDoctorAll               = "doctor:*"
	PatternRoomTypeAll             = "room_type:*"
	PatternRoomAll                 = "room:*"
	PatternTypeTestCategoryAll     = "type_test_category:*"
	PatternTypeTestAll             = "type_test:*"
	PatternMedicineTypeAll         = "medicine_type:*"
	PatternMedicineAll             = "medicine:*"
	PatternAppointmentAll          = "appointment:*"
	PatternMedicalRecordAll        = "medical_record:*"
	PatternHospitalizationAll      = "hospitalization:*"
)

// ─── Medical Record ────────────────────────────────────────────────────────

func MedicalRecordKey(id uint) string {
	return fmt.Sprintf("medical_record:id:%d", id)
}

func MedicalRecordListQueryKey(page, size int, patientID, doctorID, departmentID *uint, status, dateFrom, dateTo string, sortBy, sortDir string) string {
	pID, dID, depID := "all", "all", "all"
	if patientID != nil {
		pID = fmt.Sprintf("%d", *patientID)
	}
	if doctorID != nil {
		dID = fmt.Sprintf("%d", *doctorID)
	}
	if departmentID != nil {
		depID = fmt.Sprintf("%d", *departmentID)
	}
	if status == "" {
		status = "all"
	}
	return fmt.Sprintf(
		"medical_record:list:p%d:s%d:pat%s:doc%s:dep%s:st%s:df%s:dt%s:sb%s:sd%s",
		page, size, pID, dID, depID, status, dateFrom, dateTo, sortBy, sortDir,
	)
}

func MedicalRecordDeletedListQueryKey(page, size int, patientID, doctorID, departmentID *uint, status, dateFrom, dateTo string, sortBy, sortDir string) string {
	pID, dID, depID := "all", "all", "all"
	if patientID != nil {
		pID = fmt.Sprintf("%d", *patientID)
	}
	if doctorID != nil {
		dID = fmt.Sprintf("%d", *doctorID)
	}
	if departmentID != nil {
		depID = fmt.Sprintf("%d", *departmentID)
	}
	if status == "" {
		status = "all"
	}
	return fmt.Sprintf(
		"medical_record:deleted:p%d:s%d:pat%s:doc%s:dep%s:st%s:df%s:dt%s:sb%s:sd%s",
		page, size, pID, dID, depID, status, dateFrom, dateTo, sortBy, sortDir,
	)
}

// ─── Hospitalization ────────────────────────────────────────────────────────

func HospitalizationKey(id uint) string {
	return fmt.Sprintf("hospitalization:id:%d", id)
}

func HospitalizationListQueryKey(page, size int, patientID, doctorID, roomID *uint, search, status, notStatus, sortBy, sortDir string) string {
	pID, dID, rID := "all", "all", "all"
	if patientID != nil {
		pID = fmt.Sprintf("%d", *patientID)
	}
	if doctorID != nil {
		dID = fmt.Sprintf("%d", *doctorID)
	}
	if roomID != nil {
		rID = fmt.Sprintf("%d", *roomID)
	}
	if status == "" {
		status = "all"
	}
	if notStatus == "" {
		notStatus = "all"
	}
	if search == "" {
		search = "none"
	}
	return fmt.Sprintf(
		"hospitalization:list:p%d:s%d:pat%s:doc%s:rm%s:q%s:st%s:nst%s:sb%s:sd%s",
		page, size, pID, dID, rID, search, status, notStatus, sortBy, sortDir,
	)
}

func HospitalizationDeletedListQueryKey(page, size int, patientID, doctorID, roomID *uint, search, status, notStatus, sortBy, sortDir string) string {
	return fmt.Sprintf("hospitalization:deleted:%s", HospitalizationListQueryKey(page, size, patientID, doctorID, roomID, search, status, notStatus, sortBy, sortDir))
}
