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

func TypeTestKey(id uint) string            { return fmt.Sprintf("typetest:id:%d", id) }
func TypeTestListKey(page, size int) string { return fmt.Sprintf("typetest:list:p%d:s%d", page, size) }
func TypeTestSearchKey(q string, page, size int) string {
	return fmt.Sprintf("typetest:search:%s:p%d:s%d", q, page, size)
}

// ─── Medicine ──────────────────────────────────────────────────────────────

func MedicineKey(id uint) string            { return fmt.Sprintf("medicine:id:%d", id) }
func MedicineNameKey(name string) string    { return fmt.Sprintf("medicine:name:%s", name) }
func MedicineListKey(page, size int) string { return fmt.Sprintf("medicine:list:p%d:s%d", page, size) }
func MedicineDeletedListKey(page, size int) string {
	return fmt.Sprintf("medicine:list:p%d:s%d", page, size)
}
func MedicineAvailableKey(page, size int) string {
	return fmt.Sprintf("medicine:available:p%d:s%d", page, size)
}
func MedicineLowStockKey(page, size int) string {
	return fmt.Sprintf("medicine:lowstock:p%d:s%d", page, size)
}

func MedicineOutOfStockKey(page, size int) string {
	return fmt.Sprintf("medicine:outofstock:p%d:s%d", page, size)
}

func MedicineInactiveKey(page, size int) string {
	return fmt.Sprintf("medicine:inactive:p%d:s%d", page, size)
}

func MedicineTypeKey(page, size int) string {
	return fmt.Sprintf("medicine:type:p%d:s%d", page, size)
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
	PatternTypeTestAll             = "typetest:*"
	PatternMedicineAll             = "medicine:*"
)
