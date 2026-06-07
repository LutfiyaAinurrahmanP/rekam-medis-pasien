package kafka

// Topic constants — semua Kafka topic yang digunakan oleh aplikasi.
// Penamaan: <domain>.<action>
const (
	// ── User Topics ──────────────────────────────────────────────────────────
	TopicUserRegistered             = "user.registered"
	TopicUserLogin                  = "user.login"
	TopicUserCreated                = "user.created"
	TopicUserUpdated                = "user.updated"
	TopicUserDeleted                = "user.deleted"
	TopicUserRestored               = "user.restored"
	TopicUserPasswordResetRequested = "user.password_reset_requested"

	// ── Patient Topics ────────────────────────────────────────────────────────
	TopicPatientCreated  = "patient.created"
	TopicPatientUpdated  = "patient.updated"
	TopicPatientDeleted  = "patient.deleted"
	TopicPatientRestored = "patient.restored"

	// ── Doctor Specialization Topics ─────────────────────────────────────────────────────────
	TopicDoctorSpecializationCreated  = "doctor_specialization.created"
	TopicDoctorSpecializationUpdated  = "doctor_specialization.updated"
	TopicDoctorSpecializationDeleted  = "doctor_specialization.deleted"
	TopicDoctorSpecializationRestored = "doctor_specialization.restored"

	// ── Doctor Topics ─────────────────────────────────────────────────────────
	TopicDoctorCreated  = "doctor.created"
	TopicDoctorUpdated  = "doctor.updated"
	TopicDoctorDeleted  = "doctor.deleted"
	TopicDoctorRestored = "doctor.restored"

	// ── Department Topics ─────────────────────────────────────────────────────
	TopicDepartmentCreated  = "department.created"
	TopicDepartmentUpdated  = "department.updated"
	TopicDepartmentDeleted  = "department.deleted"
	TopicDepartmentRestored = "department.restored"

	// ── Room Type Topics ───────────────────────────────────────────────────────────
	TopicRoomTypeCreated  = "room_type.created"
	TopicRoomTypeUpdated  = "room_type.updated"
	TopicRoomTypeDeleted  = "room_type.deleted"
	TopicRoomTypeRestored = "room_type.restored"

	// ── Room Topics ───────────────────────────────────────────────────────────
	TopicRoomCreated  = "room.created"
	TopicRoomUpdated  = "room.updated"
	TopicRoomDeleted  = "room.deleted"
	TopicRoomRestored = "room.restored"

	// ── TypeTestCategory Topics ───────────────────────────────────────────────
	TopicTypeTestCategoryCreated  = "type_test_category.created"
	TopicTypeTestCategoryUpdated  = "type_test_category.updated"
	TopicTypeTestCategoryDeleted  = "type_test_category.deleted"
	TopicTypeTestCategoryRestored = "type_test_category.restored"

	// ── TypeTest Topics ───────────────────────────────────────────────────────
	TopicTypeTestCreated  = "type_test.created"
	TopicTypeTestUpdated  = "type_test.updated"
	TopicTypeTestDeleted  = "type_test.deleted"
	TopicTypeTestRestored = "type_test.restored"

	// ── Medicine Type Topics ───────────────────────────────────────────────────────
	TopicMedicineTypeCreated  = "medicine_type.created"
	TopicMedicineTypeUpdated  = "medicine_type.updated"
	TopicMedicineTypeDeleted  = "medicine_type.deleted"
	TopicMedicineTypeRestored = "medicine_type.restored"
)

// AllTopics mengembalikan semua topic yang perlu dibuat/dipastikan ada di Kafka.
func AllTopics() []string {
	return []string{
		TopicUserRegistered,
		TopicUserLogin,
		TopicUserCreated,
		TopicUserUpdated,
		TopicUserDeleted,
		TopicUserRestored,
		TopicUserPasswordResetRequested,

		TopicPatientCreated,
		TopicPatientUpdated,
		TopicPatientDeleted,
		TopicPatientRestored,

		TopicDoctorSpecializationCreated,
		TopicDoctorSpecializationUpdated,
		TopicDoctorSpecializationDeleted,
		TopicDoctorSpecializationRestored,

		TopicDoctorCreated,
		TopicDoctorUpdated,
		TopicDoctorDeleted,
		TopicDoctorRestored,

		TopicDepartmentCreated,
		TopicDepartmentUpdated,
		TopicDepartmentDeleted,
		TopicDepartmentRestored,

		TopicRoomTypeCreated,
		TopicRoomTypeUpdated,
		TopicRoomTypeDeleted,
		TopicRoomTypeRestored,

		TopicRoomCreated,
		TopicRoomUpdated,
		TopicRoomDeleted,
		TopicRoomRestored,

		TopicTypeTestCategoryCreated,
		TopicTypeTestCategoryUpdated,
		TopicTypeTestCategoryDeleted,
		TopicTypeTestCategoryRestored,

		TopicTypeTestCreated,
		TopicTypeTestUpdated,
		TopicTypeTestDeleted,
		TopicTypeTestRestored,

		TopicMedicineTypeCreated,
		TopicMedicineTypeUpdated,
		TopicMedicineTypeDeleted,
		TopicMedicineTypeRestored,
	}
}
