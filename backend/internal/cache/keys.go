package cache

import "fmt"

// Konvensi penamaan key: {domain}:{qualifier}:{identifier}
// Contoh: user:id:42, room:available:p1:s10

// ─── User ──────────────────────────────────────────────────────────────────

func UserKey(id uint) string            { return fmt.Sprintf("user:id:%d", id) }
func UserListKey(page, size int) string { return fmt.Sprintf("user:list:p%d:s%d", page, size) }
func UserDeletedListKey(page, size int) string {
	return fmt.Sprintf("user:deleted:p%d:s%d", page, size)
}

// ─── Department ────────────────────────────────────────────────────────────

func DepartmentKey(id uint) string { return fmt.Sprintf("department:id:%d", id) }
func DepartmentListKey(page, size int) string {
	return fmt.Sprintf("department:list:p%d:s%d", page, size)
}

// ─── Patient ───────────────────────────────────────────────────────────────

func PatientKey(id uint) string             { return fmt.Sprintf("patient:id:%d", id) }
func PatientListKey(page, size int) string  { return fmt.Sprintf("patient:list:p%d:s%d", page, size) }
func PatientByCodeKey(code string) string   { return fmt.Sprintf("patient:code:%s", code) }
func PatientByUserIDKey(userID uint) string { return fmt.Sprintf("patient:user:%d", userID) }

// ─── Doctor ────────────────────────────────────────────────────────────────

func DoctorKey(id uint) string             { return fmt.Sprintf("doctor:id:%d", id) }
func DoctorListKey(page, size int) string  { return fmt.Sprintf("doctor:list:p%d:s%d", page, size) }
func DoctorBySpecKey(spec string) string   { return fmt.Sprintf("doctor:spec:%s", spec) }
func DoctorByUserIDKey(userID uint) string { return fmt.Sprintf("doctor:user:%d", userID) }

// ─── Room ──────────────────────────────────────────────────────────────────

func RoomKey(id uint) string            { return fmt.Sprintf("room:id:%d", id) }
func RoomListKey(page, size int) string { return fmt.Sprintf("room:list:p%d:s%d", page, size) }
func RoomAvailableKey(page, size int) string {
	return fmt.Sprintf("room:available:p%d:s%d", page, size)
}
func RoomOccupiedKey(page, size int) string { return fmt.Sprintf("room:occupied:p%d:s%d", page, size) }
func RoomInactiveKey(page, size int) string { return fmt.Sprintf("room:inactive:p%d:s%d", page, size) }
func RoomByNumberKey(number string) string  { return fmt.Sprintf("room:number:%s", number) }
func RoomByTypeKey(roomType string) string  { return fmt.Sprintf("room:type:%s", roomType) }
func RoomByDeptKey(deptID string) string    { return fmt.Sprintf("room:dept:%s", deptID) }

// ─── TypeTest ──────────────────────────────────────────────────────────────

func TypeTestKey(id uint) string            { return fmt.Sprintf("typetest:id:%d", id) }
func TypeTestListKey(page, size int) string { return fmt.Sprintf("typetest:list:p%d:s%d", page, size) }
func TypeTestSearchKey(q string, page, size int) string {
	return fmt.Sprintf("typetest:search:%s:p%d:s%d", q, page, size)
}

// ─── Medicine ──────────────────────────────────────────────────────────────

func MedicineKey(id uint) string            { return fmt.Sprintf("medicine:id:%d", id) }
func MedicineNameKey(name string) string            { return fmt.Sprintf("medicine:name:%d", name) }
func MedicineListKey(page, size int) string { return fmt.Sprintf("medicine:list:p%d:s%d", page, size) }
func MedicineDeletedListKey(page, size int) string { return fmt.Sprintf("medicine:list:p%d:s%d", page, size) }
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

// ─── Invalidation patterns ─────────────────────────────────────────────────

const (
	PatternUserAll       = "user:*"
	PatternDepartmentAll = "department:*"
	PatternPatientAll    = "patient:*"
	PatternDoctorAll     = "doctor:*"
	PatternRoomAll       = "room:*"
	PatternTypeTestAll   = "typetest:*"
	PatternMedicineAll   = "medicine:*"
)
