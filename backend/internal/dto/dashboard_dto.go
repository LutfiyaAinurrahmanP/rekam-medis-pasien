package dto

// ─── Dashboard Query ────────────────────────────────────────────────────────

type DashboardOverviewQuery struct {
	Date string `form:"date" binding:"omitempty,datetime=2006-01-02"`
}

type DashboardPeriodQuery struct {
	Period    string `form:"period,default=today" binding:"omitempty,oneof=today this_week this_month last_month custom"`
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
}

type DashboardDoctorQuery struct {
	Date string `form:"date" binding:"omitempty,datetime=2006-01-02"`
}

type DashboardAppointmentReportQuery struct {
	Period       string `form:"period,default=this_month" binding:"omitempty,oneof=today this_week this_month last_month custom"`
	StartDate    string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate      string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	DoctorID     *uint  `form:"doctor_id" binding:"omitempty"`
	DepartmentID *uint  `form:"department_id" binding:"omitempty"`
	GroupBy      string `form:"group_by,default=day" binding:"omitempty,oneof=day week month"`
}

type DashboardRevenueReportQuery struct {
	Period    string `form:"period,default=this_month" binding:"omitempty,oneof=today this_week this_month last_month custom"`
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	GroupBy   string `form:"group_by,default=day" binding:"omitempty,oneof=day week month"`
}

type DashboardPatientReportQuery struct {
	Period    string `form:"period,default=this_month" binding:"omitempty,oneof=today this_week this_month last_month custom"`
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
}

// ─── Shared / Embedded ─────────────────────────────────────────────────────

type DashboardPeriodRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type DashboardTopDepartment struct {
	DepartmentID   uint   `json:"department_id"`
	DepartmentName string `json:"department_name"`
	Count          int64  `json:"appointment_count"`
}

type DashboardTrendPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type DashboardRevenueTrendPoint struct {
	Date   string  `json:"date"`
	Billed float64 `json:"billed"`
	Paid   float64 `json:"paid"`
}

type DashboardPatientTrendPoint struct {
	Date        string `json:"date"`
	NewPatients int64  `json:"new_patients"`
	TotalVisits int64  `json:"total_visits"`
}

// ─── Overview Response ─────────────────────────────────────────────────────

type DashboardOverviewResponse struct {
	SummaryDate string                     `json:"summary_date"`
	MasterData  DashboardMasterData        `json:"master_data"`
	Today       DashboardTodayStats        `json:"today"`
	Rooms       DashboardRoomStats         `json:"rooms"`
	Billing     DashboardOverviewBilling   `json:"billing"`
}

type DashboardMasterData struct {
	TotalPatients    int64 `json:"total_patients"`
	TotalDoctors     int64 `json:"total_doctors"`
	TotalDepartments int64 `json:"total_departments"`
	TotalRooms       int64 `json:"total_rooms"`
	TotalMedicines   int64 `json:"total_medicines"`
}

type DashboardTodayStats struct {
	Appointments          DashboardAppointmentCounts `json:"appointments"`
	NewPatients           int64                      `json:"new_patients"`
	NewMedicalRecords     int64                      `json:"new_medical_records"`
	ActiveHospitalizations int64                     `json:"active_hospitalizations"`
}

type DashboardAppointmentCounts struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Confirmed int64 `json:"confirmed"`
	Completed int64 `json:"completed"`
	Cancelled int64 `json:"cancelled"`
}

type DashboardRoomStats struct {
	Total          int64   `json:"total"`
	Available      int64   `json:"available"`
	Occupied       int64   `json:"occupied"`
	OutOfService   int64   `json:"out_of_service"`
	OccupancyRate  float64 `json:"occupancy_rate"`
}

type DashboardOverviewBilling struct {
	TodayRevenue float64 `json:"today_revenue"`
	UnpaidCount  int64   `json:"unpaid_count"`
	UnpaidTotal  float64 `json:"unpaid_total"`
}

// ─── Admin Dashboard Response ──────────────────────────────────────────────

type DashboardAdminResponse struct {
	Period      string                `json:"period"`
	PeriodRange DashboardPeriodRange  `json:"period_range"`
	Appointments DashboardAdminAppointments `json:"appointments"`
	Patients    DashboardAdminPatients    `json:"patients"`
	Hospitalization DashboardAdminHospitalization `json:"hospitalization"`
	Billing     DashboardAdminBilling     `json:"billing"`
	Referrals   DashboardAdminReferrals   `json:"referrals"`
	TopDepartments []DashboardTopDepartment `json:"top_departments"`
	AppointmentTrend []DashboardTrendPoint `json:"appointment_trend"`
	RevenueTrend []DashboardRevenueTrendPoint `json:"revenue_trend"`
}

type DashboardAdminAppointments struct {
	Total          int64   `json:"total"`
	Completed      int64   `json:"completed"`
	Cancelled      int64   `json:"cancelled"`
	NoShow         int64   `json:"no_show"`
	CompletionRate float64 `json:"completion_rate"`
}

type DashboardAdminPatients struct {
	TotalActive       int64 `json:"total_active"`
	NewRegistrations  int64 `json:"new_registrations"`
	Returning         int64 `json:"returning"`
}

type DashboardAdminHospitalization struct {
	TotalAdmissions          int64   `json:"total_admissions"`
	CurrentlyHospitalized    int64   `json:"currently_hospitalized"`
	TotalDischarged          int64   `json:"total_discharged"`
	AverageLengthOfStayDays  float64 `json:"average_length_of_stay_days"`
	AvailableBeds            int64   `json:"available_beds"`
}

type DashboardAdminBilling struct {
	TotalRevenue float64 `json:"total_revenue"`
	PaidCount    int64   `json:"paid_count"`
	UnpaidCount  int64   `json:"unpaid_count"`
	UnpaidTotal  float64 `json:"unpaid_total"`
	AverageBill  float64 `json:"average_bill"`
}

type DashboardAdminReferrals struct {
	TotalIssued int64 `json:"total_issued"`
	Internal    int64 `json:"internal"`
	External    int64 `json:"external"`
	Pending     int64 `json:"pending"`
	Completed   int64 `json:"completed"`
}

// ─── Doctor Dashboard Response ─────────────────────────────────────────────

type DashboardDoctorResponse struct {
	Doctor               DashboardDoctorInfo          `json:"doctor"`
	Date                 string                       `json:"date"`
	TodaySchedule        DashboardDoctorSchedule      `json:"today_schedule"`
	UpcomingAppointments []DashboardDoctorAppointment `json:"upcoming_appointments"`
	Statistics           DashboardDoctorStats         `json:"statistics"`
	RecentMedicalRecords []DashboardDoctorRecord      `json:"recent_medical_records"`
}

type DashboardDoctorInfo struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
}

type DashboardDoctorSchedule struct {
	TotalAppointments int64                        `json:"total_appointments"`
	Completed         int64                        `json:"completed"`
	Pending           int64                        `json:"pending"`
	Cancelled         int64                        `json:"cancelled"`
	NextPatient       *DashboardDoctorNextPatient  `json:"next_patient,omitempty"`
}

type DashboardDoctorNextPatient struct {
	AppointmentID   uint   `json:"appointment_id"`
	AppointmentTime string `json:"appointment_time"`
	PatientName     string `json:"patient_name"`
	ChiefComplaint  string `json:"chief_complaint"`
}

type DashboardDoctorAppointment struct {
	AppointmentID   uint   `json:"appointment_id"`
	AppointmentTime string `json:"appointment_time"`
	PatientName     string `json:"patient_name"`
	PatientID       uint   `json:"patient_id"`
}

type DashboardDoctorStats struct {
	TotalPatientsThisMonth int64 `json:"total_patients_this_month"`
	TotalPatientsToday     int64 `json:"total_patients_today"`
	PendingLabResults      int64 `json:"pending_lab_results"`
	PendingPrescriptions   int64 `json:"pending_prescriptions"`
	PendingReferrals       int64 `json:"pending_referrals"`
}

type DashboardDoctorRecord struct {
	ID             uint   `json:"id"`
	PatientName    string `json:"patient_name"`
	VisitDate      string `json:"visit_date"`
	ChiefComplaint string `json:"chief_complaint"`
	Status         string `json:"status"`
}

// ─── Receptionist Dashboard Response ───────────────────────────────────────

type DashboardReceptionistResponse struct {
	Date                        string                            `json:"date"`
	AppointmentsToday           DashboardReceptionistAppointments `json:"appointments_today"`
	Hospitalization             DashboardReceptionistHospitalization `json:"hospitalization"`
	Billing                     DashboardReceptionistBilling      `json:"billing"`
	NewPatientRegistrationsToday int64                             `json:"new_patient_registrations_today"`
}

type DashboardReceptionistAppointments struct {
	Total             int64                          `json:"total"`
	Pending           int64                          `json:"pending"`
	Confirmed         int64                          `json:"confirmed"`
	InProgress        int64                          `json:"in_progress"`
	Completed         int64                          `json:"completed"`
	Cancelled         int64                          `json:"cancelled"`
	AppointmentQueue  []DashboardAppointmentQueueItem `json:"appointment_queue"`
}

type DashboardAppointmentQueueItem struct {
	AppointmentID uint   `json:"appointment_id"`
	ScheduledTime string `json:"scheduled_time"`
	PatientName   string `json:"patient_name"`
	PatientID     uint   `json:"patient_id"`
	DoctorName    string `json:"doctor_name"`
	Status        string `json:"status"`
}

type DashboardReceptionistHospitalization struct {
	CurrentlyHospitalized int64 `json:"currently_hospitalized"`
	AvailableBeds         int64 `json:"available_beds"`
	NewAdmissionsToday    int64 `json:"new_admissions_today"`
	DischargedToday       int64 `json:"discharged_today"`
}

type DashboardReceptionistBilling struct {
	UnpaidCount int64   `json:"unpaid_count"`
	UnpaidTotal float64 `json:"unpaid_total"`
}

// ─── Patient Dashboard Response ─────────────────────────────────────────────

type DashboardPatientResponse struct {
	Patient             DashboardPatientInfo         `json:"patient"`
	UpcomingAppointments []DashboardPatientAppointment `json:"upcoming_appointments"`
	Billing             DashboardPatientBilling      `json:"billing"`
	ActivePrescriptions []DashboardPatientPrescription `json:"active_prescriptions"`
	PendingLabResults   []DashboardPatientLabResult  `json:"pending_lab_results"`
	ActiveReferrals     []DashboardPatientReferral   `json:"active_referrals"`
	RecentVisits        []DashboardPatientVisit      `json:"recent_visits"`
}

type DashboardPatientInfo struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	MedicalRecordNumber string `json:"medical_record_number"`
	DateOfBirth         string `json:"date_of_birth"`
	Age                 int    `json:"age"`
}

type DashboardPatientAppointment struct {
	ID             uint   `json:"id"`
	ScheduledDate  string `json:"scheduled_date"`
	ScheduledTime  string `json:"scheduled_time"`
	DoctorName     string `json:"doctor_name"`
	DepartmentName string `json:"department_name"`
	Status         string `json:"status"`
}

type DashboardPatientBilling struct {
	UnpaidCount int64                        `json:"unpaid_count"`
	UnpaidTotal float64                      `json:"unpaid_total"`
	UnpaidBills []DashboardPatientUnpaidBill `json:"unpaid_bills"`
}

type DashboardPatientUnpaidBill struct {
	ID            uint    `json:"id"`
	InvoiceNumber string  `json:"invoice_number"`
	TotalAmount   float64 `json:"total_amount"`
	DueDate       string  `json:"due_date"`
	Status        string  `json:"status"`
}

type DashboardPatientPrescription struct {
	ID             uint   `json:"id"`
	PrescribedDate string `json:"prescribed_date"`
	DoctorName     string `json:"doctor_name"`
	MedicinesCount int    `json:"medicines_count"`
	Status         string `json:"status"`
}

type DashboardPatientLabResult struct {
	ID          uint   `json:"id"`
	TestName    string `json:"test_name"`
	OrderedDate string `json:"ordered_date"`
	Status      string `json:"status"`
}

type DashboardPatientReferral struct {
	ID             uint   `json:"id"`
	ReferralNumber string `json:"referral_number"`
	ReferredTo     string `json:"referred_to"`
	ReferralDate   string `json:"referral_date"`
	Status         string `json:"status"`
}

type DashboardPatientVisit struct {
	ID             uint   `json:"id"`
	VisitDate      string `json:"visit_date"`
	DoctorName     string `json:"doctor_name"`
	ChiefComplaint string `json:"chief_complaint"`
	Status         string `json:"status"`
}

// ─── Appointment Report Response ────────────────────────────────────────────

type DashboardAppointmentReportResponse struct {
	Period      string               `json:"period"`
	PeriodRange DashboardPeriodRange `json:"period_range"`
	Totals      DashboardAppointmentReportTotals `json:"totals"`
	ByDepartment []DashboardDeptStats `json:"by_department"`
	ByDoctor    []DashboardDoctorStatItem `json:"by_doctor"`
	Trend       []DashboardAppointmentTrendItem `json:"trend"`
	PeakHours   []DashboardPeakHour  `json:"peak_hours"`
}

type DashboardAppointmentReportTotals struct {
	Total            int64   `json:"total"`
	Scheduled        int64   `json:"scheduled"`
	Confirmed        int64   `json:"confirmed"`
	Completed        int64   `json:"completed"`
	Cancelled        int64   `json:"cancelled"`
	NoShow           int64   `json:"no_show"`
	CompletionRate   float64 `json:"completion_rate"`
	CancellationRate float64 `json:"cancellation_rate"`
}

type DashboardDeptStats struct {
	DepartmentID   uint    `json:"department_id"`
	DepartmentName string  `json:"department_name"`
	Total          int64   `json:"total"`
	Completed      int64   `json:"completed"`
	CompletionRate float64 `json:"completion_rate"`
}

type DashboardDoctorStatItem struct {
	DoctorID       uint    `json:"doctor_id"`
	DoctorName     string  `json:"doctor_name"`
	Total          int64   `json:"total"`
	Completed      int64   `json:"completed"`
	CompletionRate float64 `json:"completion_rate"`
}

type DashboardAppointmentTrendItem struct {
	Date      string `json:"date"`
	Total     int64  `json:"total"`
	Completed int64  `json:"completed"`
	Cancelled int64  `json:"cancelled"`
}

type DashboardPeakHour struct {
	Hour  string `json:"hour"`
	Count int64  `json:"count"`
}

// ─── Revenue Report Response ────────────────────────────────────────────────

type DashboardRevenueReportResponse struct {
	Period      string               `json:"period"`
	PeriodRange DashboardPeriodRange `json:"period_range"`
	Revenue     DashboardRevenueTotals `json:"revenue"`
	ByCategory  []DashboardRevenueCategory `json:"by_category"`
	Trend       []DashboardRevenueTrendPoint `json:"trend"`
	Comparison  DashboardRevenueComparison `json:"comparison"`
}

type DashboardRevenueTotals struct {
	TotalBilled        float64 `json:"total_billed"`
	TotalPaid          float64 `json:"total_paid"`
	TotalUnpaid        float64 `json:"total_unpaid"`
	TotalBills         int64   `json:"total_bills"`
	PaidBills          int64   `json:"paid_bills"`
	UnpaidBills        int64   `json:"unpaid_bills"`
	AverageBillAmount  float64 `json:"average_bill_amount"`
	CollectionRate     float64 `json:"collection_rate"`
}

type DashboardRevenueCategory struct {
	Category    string  `json:"category"`
	TotalAmount float64 `json:"total_amount"`
	BillCount   int64   `json:"bill_count"`
}

type DashboardRevenueComparison struct {
	VsPreviousPeriod DashboardRevenueChanges `json:"vs_previous_period"`
}

type DashboardRevenueChanges struct {
	RevenueChangePercent   float64 `json:"revenue_change_percent"`
	BillCountChangePercent float64 `json:"bill_count_change_percent"`
}

// ─── Patient Report Response ─────────────────────────────────────────────────

type DashboardPatientReportResponse struct {
	Period            string                         `json:"period"`
	PeriodRange       DashboardPeriodRange            `json:"period_range"`
	Registrations     DashboardPatientRegistrations   `json:"registrations"`
	Demographics      DashboardPatientDemographics    `json:"demographics"`
	RegistrationTrend []DashboardPatientTrendPoint    `json:"registration_trend"`
	Comparison        DashboardPatientComparison      `json:"comparison"`
}

type DashboardPatientRegistrations struct {
	NewPatients       int64 `json:"new_patients"`
	ReturningPatients int64 `json:"returning_patients"`
	TotalVisits       int64 `json:"total_visits"`
	UniqueVisitors    int64 `json:"unique_visitors"`
}

type DashboardPatientDemographics struct {
	ByGender   []DashboardGenderCount   `json:"by_gender"`
	ByAgeGroup []DashboardAgeGroupCount `json:"by_age_group"`
}

type DashboardGenderCount struct {
	Gender string `json:"gender"`
	Count  int64  `json:"count"`
}

type DashboardAgeGroupCount struct {
	Group string `json:"group"`
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type DashboardPatientComparison struct {
	VsPreviousPeriod DashboardPatientChanges `json:"vs_previous_period"`
}

type DashboardPatientChanges struct {
	NewPatientsChangePercent  float64 `json:"new_patients_change_percent"`
	TotalVisitsChangePercent  float64 `json:"total_visits_change_percent"`
}
