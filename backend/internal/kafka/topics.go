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

	// ── Medicine Topics ────────────────────────────────────────────────────────────
	TopicMedicineCreated      = "medicine.created"
	TopicMedicineUpdated      = "medicine.updated"
	TopicMedicineDeleted      = "medicine.deleted"
	TopicMedicineRestored     = "medicine.restored"
	TopicMedicineActivated    = "medicine.activated"
	TopicMedicineDeactivated  = "medicine.deactivated"
	TopicMedicineStockAdded   = "medicine.stock_added"
	TopicMedicineStockReduced = "medicine.stock_reduced"

	// ── Appointment Topics ─────────────────────────────────────────────────────────
	TopicAppointmentCreated     = "appointment.created"
	TopicAppointmentUpdated     = "appointment.updated"
	TopicAppointmentDeleted     = "appointment.deleted"
	TopicAppointmentRestored    = "appointment.restored"
	TopicAppointmentConfirmed   = "appointment.confirmed"
	TopicAppointmentStarted     = "appointment.started"
	TopicAppointmentCompleted   = "appointment.completed"
	TopicAppointmentCancelled   = "appointment.cancelled"
	TopicAppointmentRescheduled = "appointment.rescheduled"
	TopicAppointmentNoShow      = "appointment.no_show"

	// ── Medical Record Topics ──────────────────────────────────────────────────
	TopicMedicalRecordCreated   = "medical_record.created"
	TopicMedicalRecordUpdated   = "medical_record.updated"
	TopicMedicalRecordDeleted   = "medical_record.deleted"
	TopicMedicalRecordRestored  = "medical_record.restored"
	TopicMedicalRecordFinalized = "medical_record.finalized"

	// ── Hospitalization Topics ─────────────────────────────────────────────────
	TopicHospitalizationCreated     = "hospitalization.created"
	TopicHospitalizationUpdated     = "hospitalization.updated"
	TopicHospitalizationDeleted     = "hospitalization.deleted"
	TopicHospitalizationRestored    = "hospitalization.restored"
	TopicHospitalizationDischarged  = "hospitalization.discharged"
	TopicHospitalizationTransferred = "hospitalization.transferred"

	// ── Lab Test Topics ────────────────────────────────────────────────────────
	TopicLabTestCreated         = "lab_test.created"
	TopicLabTestUpdated         = "lab_test.updated"
	TopicLabTestSampleCollected = "lab_test.sample_collected"
	TopicLabTestStarted         = "lab_test.started"
	TopicLabTestCompleted       = "lab_test.completed"
	TopicLabTestCancelled       = "lab_test.cancelled"
	TopicLabTestDeleted         = "lab_test.deleted"
	TopicLabTestRestored        = "lab_test.restored"

	// ── Prescription Topics ────────────────────────────────────────────────────
	TopicPrescriptionCreated   = "prescription.created"
	TopicPrescriptionUpdated   = "prescription.updated"
	TopicPrescriptionDispensed = "prescription.dispensed"
	TopicPrescriptionCancelled = "prescription.cancelled"
	TopicPrescriptionDeleted   = "prescription.deleted"
	TopicPrescriptionRestored  = "prescription.restored"

	// ── Prescription Item Topics ────────────────────────────────────────────────
	TopicPrescriptionItemCreated = "prescription_item.created"
	TopicPrescriptionItemUpdated = "prescription_item.updated"
	TopicPrescriptionItemDeleted = "prescription_item.deleted"

	// ── Vital Sign Topics ──────────────────────────────────────────────────────
	TopicVitalSignCreated  = "vital_sign.created"
	TopicVitalSignUpdated  = "vital_sign.updated"
	TopicVitalSignDeleted  = "vital_sign.deleted"
	TopicVitalSignRestored = "vital_sign.restored"

	// ── Allergy Topics ─────────────────────────────────────────────────────────
	TopicAllergyCreated = "allergy.created"
	TopicAllergyUpdated = "allergy.updated"
	TopicAllergyDeleted = "allergy.deleted"

	// ── Medical Condition Topics ───────────────────────────────────────────────
	TopicMedicalConditionCreated = "medical_condition.created"
	TopicMedicalConditionUpdated = "medical_condition.updated"
	TopicMedicalConditionDeleted = "medical_condition.deleted"

	// ── Surgical History Topics ────────────────────────────────────────────────
	TopicSurgicalHistoryCreated = "surgical_history.created"
	TopicSurgicalHistoryUpdated = "surgical_history.updated"
	TopicSurgicalHistoryDeleted = "surgical_history.deleted"

	// ── Family History Topics ──────────────────────────────────────────────────
	TopicFamilyHistoryCreated = "family_history.created"
	TopicFamilyHistoryUpdated = "family_history.updated"
	TopicFamilyHistoryDeleted = "family_history.deleted"

	// ── Referral Topics ────────────────────────────────────────────────────────
	TopicReferralCreated   = "referral.created"
	TopicReferralUpdated   = "referral.updated"
	TopicReferralDeleted   = "referral.deleted"
	TopicReferralRestored  = "referral.restored"
	TopicReferralAccepted  = "referral.accepted"
	TopicReferralRejected  = "referral.rejected"
	TopicReferralCompleted = "referral.completed"
	TopicReferralCancelled = "referral.cancelled"
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

		TopicMedicineCreated,
		TopicMedicineUpdated,
		TopicMedicineDeleted,
		TopicMedicineRestored,
		TopicMedicineActivated,
		TopicMedicineDeactivated,
		TopicMedicineStockAdded,
		TopicMedicineStockReduced,

		TopicAppointmentCreated,
		TopicAppointmentUpdated,
		TopicAppointmentDeleted,
		TopicAppointmentRestored,
		TopicAppointmentConfirmed,
		TopicAppointmentStarted,
		TopicAppointmentCompleted,
		TopicAppointmentCancelled,
		TopicAppointmentRescheduled,
		TopicAppointmentNoShow,

		TopicMedicalRecordCreated,
		TopicMedicalRecordUpdated,
		TopicMedicalRecordDeleted,
		TopicMedicalRecordRestored,
		TopicMedicalRecordFinalized,

		TopicHospitalizationCreated,
		TopicHospitalizationUpdated,
		TopicHospitalizationDeleted,
		TopicHospitalizationRestored,
		TopicHospitalizationDischarged,
		TopicHospitalizationTransferred,

		TopicLabTestCreated,
		TopicLabTestUpdated,
		TopicLabTestSampleCollected,
		TopicLabTestStarted,
		TopicLabTestCompleted,
		TopicLabTestCancelled,
		TopicLabTestDeleted,
		TopicLabTestRestored,

		TopicPrescriptionCreated,
		TopicPrescriptionUpdated,
		TopicPrescriptionDispensed,
		TopicPrescriptionCancelled,
		TopicPrescriptionDeleted,
		TopicPrescriptionRestored,

		TopicPrescriptionItemCreated,
		TopicPrescriptionItemUpdated,
		TopicPrescriptionItemDeleted,

		TopicVitalSignCreated,
		TopicVitalSignUpdated,
		TopicVitalSignDeleted,
		TopicVitalSignRestored,

		TopicAllergyCreated,
		TopicAllergyUpdated,
		TopicAllergyDeleted,

		TopicMedicalConditionCreated,
		TopicMedicalConditionUpdated,
		TopicMedicalConditionDeleted,

		TopicSurgicalHistoryCreated,
		TopicSurgicalHistoryUpdated,
		TopicSurgicalHistoryDeleted,

		TopicFamilyHistoryCreated,
		TopicFamilyHistoryUpdated,
		TopicFamilyHistoryDeleted,

		TopicReferralCreated,
		TopicReferralUpdated,
		TopicReferralDeleted,
		TopicReferralRestored,
		TopicReferralAccepted,
		TopicReferralRejected,
		TopicReferralCompleted,
		TopicReferralCancelled,
	}
}
