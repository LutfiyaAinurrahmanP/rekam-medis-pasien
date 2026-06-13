package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupAPIRouter sets up all API routes by calling individual router setups
func SetupAPIRouter(rg *gin.RouterGroup, cfg *RouteConfig) {
	// Setup Auth routes
	SetupAuthRouter(rg, cfg)

	// Setup Users routes
	SetupUsersRouter(rg, cfg)

	// Setup Departments routes
	if cfg.DepartmentHandler != nil {
		SetupDepartmentsRouter(rg, cfg, cfg.DepartmentHandler)
	}

	if cfg.DoctorSpecializationHandler != nil {
		SetupDoctorSpecializationRouter(rg, cfg, cfg.DoctorSpecializationHandler)
	}

	// Setup Doctors routes
	if cfg.DoctorHandler != nil {
		SetupDoctorRouter(rg, cfg, cfg.DoctorHandler)
	}

	// Setup Patients routes
	if cfg.PatientHandler != nil {
		SetupPatientsRouter(rg, cfg, cfg.PatientHandler)
	}

	// Setup Room Types routes
	if cfg.RoomTypeHandler != nil {
		SetupRoomTypeRouter(rg, cfg, cfg.RoomTypeHandler)
	}

	// Setup Rooms routes
	if cfg.RoomHandler != nil {
		SetupRoomsRouter(rg, cfg, cfg.RoomHandler)
	}

	// Setup Type Test Categories routes
	if cfg.TypeTestCategoryHandler != nil {
		SetupTypeTestCategoryRouter(rg, cfg, cfg.TypeTestCategoryHandler)
	}

	// Setup Type Tests routes
	if cfg.TypeTestHandler != nil {
		SetupTypeTestRouter(rg, cfg, cfg.TypeTestHandler)
	}

	if cfg.MedicineHandler != nil {
		SetupMedicineRouter(rg, cfg, cfg.MedicineHandler)
	}

	if cfg.MedicineTypeHandler != nil {
		SetupMedicineTypeRouter(rg, cfg, cfg.MedicineTypeHandler)
	}

	if cfg.AppointmentHandler != nil {
		SetupAppointmentRouter(rg, cfg, cfg.AppointmentHandler)
	}

	if cfg.MedicalRecordHandler != nil {
		SetupMedicalRecordRouter(rg, cfg, cfg.MedicalRecordHandler)
	}

	if cfg.HospitalizationHandler != nil {
		SetupHospitalizationRouter(rg, cfg, cfg.HospitalizationHandler)
	}

	if cfg.LabTestHandler != nil {
		SetupLabTestRouter(rg, cfg, cfg.LabTestHandler)
	}

	if cfg.PrescriptionHandler != nil {
		SetupPrescriptionRoutes(rg, cfg, cfg.PrescriptionHandler)
	}

	if cfg.VitalSignHandler != nil {
		SetupVitalSignRoutes(rg, cfg, cfg.VitalSignHandler)
	}

	if cfg.AllergyHandler != nil {
		SetupAllergyRoutes(rg, cfg, cfg.AllergyHandler)
	}

	if cfg.MedicalConditionHandler != nil {
		SetupMedicalConditionRoutes(rg, cfg, cfg.MedicalConditionHandler)
	}

	if cfg.SurgicalHistoryHandler != nil {
		SetupSurgicalHistoryRoutes(rg, cfg, cfg.SurgicalHistoryHandler)
	}

	if cfg.FamilyHistoryHandler != nil {
		SetupFamilyHistoryRoutes(rg, cfg, cfg.FamilyHistoryHandler)
	}

	if cfg.ReferralHandler != nil {
		SetupReferralRoutes(rg, cfg, cfg.ReferralHandler)
	}

	if cfg.MedicalHistoryHandler != nil {
		SetupMedicalHistoryRoutes(rg, cfg, cfg.MedicalHistoryHandler)
	}

	if cfg.DashboardHandler != nil {
		SetupDashboardRouter(rg, cfg, cfg.DashboardHandler)
	}

	if cfg.BillingHandler != nil {
		SetupBillingRoutes(rg, cfg, cfg.BillingHandler)
	}
}
