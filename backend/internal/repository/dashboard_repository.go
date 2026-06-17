package repository

import (
	"fmt"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetOverview(date string) (*dto.DashboardOverviewResponse, error)
	GetAdminDashboard(query *dto.DashboardPeriodQuery, startDate, endDate string) (*dto.DashboardAdminResponse, error)
	GetDoctorDashboard(doctorID uint, date string) (*dto.DashboardDoctorResponse, error)
	GetReceptionistDashboard(date string) (*dto.DashboardReceptionistResponse, error)
	GetPatientDashboard(patientID uint) (*dto.DashboardPatientResponse, error)
	GetAppointmentReport(query *dto.DashboardAppointmentReportQuery, startDate, endDate string) (*dto.DashboardAppointmentReportResponse, error)
	GetRevenueReport(query *dto.DashboardRevenueReportQuery, startDate, endDate string) (*dto.DashboardRevenueReportResponse, error)
	GetPatientReport(query *dto.DashboardPatientReportQuery, startDate, endDate string) (*dto.DashboardPatientReportResponse, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetOverview(date string) (*dto.DashboardOverviewResponse, error) {
	var totalPatients, totalDoctors, totalDepartments, totalRooms, totalMedicines int64
	r.db.Model(&models.Patient{}).Count(&totalPatients)
	r.db.Model(&models.Doctor{}).Count(&totalDoctors)
	r.db.Model(&models.Room{}).Count(&totalRooms)
	r.db.Model(&models.Medicine{}).Count(&totalMedicines)
	r.db.Table("departments").Count(&totalDepartments)

	// Today appointment counts
	var apptTotal, apptPending, apptConfirmed, apptCompleted, apptCancelled int64
	r.db.Model(&models.Appointment{}).Where("appointment_date = ?", date).Count(&apptTotal)
	r.db.Model(&models.Appointment{}).Where("appointment_date = ? AND status = ?", date, "scheduled").Count(&apptPending)
	r.db.Model(&models.Appointment{}).Where("appointment_date = ? AND status = ?", date, "confirmed").Count(&apptConfirmed)
	r.db.Model(&models.Appointment{}).Where("appointment_date = ? AND status = ?", date, "completed").Count(&apptCompleted)
	r.db.Model(&models.Appointment{}).Where("appointment_date = ? AND status = ?", date, "cancelled").Count(&apptCancelled)

	var newPatients, newRecords, activeHosp int64
	r.db.Model(&models.Patient{}).Where("DATE(created_at) = ?", date).Count(&newPatients)
	r.db.Model(&models.MedicalRecord{}).Where("visit_date = ?", date).Count(&newRecords)
	r.db.Model(&models.Hospitalization{}).Where("status = ?", "admitted").Count(&activeHosp)

	// Room stats
	var availableRooms, occupiedRooms int64
	r.db.Model(&models.Room{}).Where("available_beds > 0 AND is_active = ?", true).Count(&availableRooms)
	r.db.Model(&models.Room{}).Where("available_beds = 0 AND is_active = ?", true).Count(&occupiedRooms)
	occupancyRate := 0.0
	if totalRooms > 0 {
		occupancyRate = float64(occupiedRooms) / float64(totalRooms) * 100
	}

	return &dto.DashboardOverviewResponse{
		SummaryDate: date,
		MasterData: dto.DashboardMasterData{
			TotalPatients:    totalPatients,
			TotalDoctors:     totalDoctors,
			TotalDepartments: totalDepartments,
			TotalRooms:       totalRooms,
			TotalMedicines:   totalMedicines,
		},
		Today: dto.DashboardTodayStats{
			Appointments: dto.DashboardAppointmentCounts{
				Total:     apptTotal,
				Pending:   apptPending,
				Confirmed: apptConfirmed,
				Completed: apptCompleted,
				Cancelled: apptCancelled,
			},
			NewPatients:            newPatients,
			NewMedicalRecords:      newRecords,
			ActiveHospitalizations: activeHosp,
		},
		Rooms: dto.DashboardRoomStats{
			Total:         totalRooms,
			Available:     availableRooms,
			Occupied:      occupiedRooms,
			OccupancyRate: occupancyRate,
		},
		Billing: dto.DashboardOverviewBilling{},
	}, nil
}

func (r *dashboardRepository) GetAdminDashboard(query *dto.DashboardPeriodQuery, startDate, endDate string) (*dto.DashboardAdminResponse, error) {
	var total, completed, cancelled, noShow int64
	r.db.Model(&models.Appointment{}).Where("appointment_date BETWEEN ? AND ?", startDate, endDate).Count(&total)
	r.db.Model(&models.Appointment{}).Where("appointment_date BETWEEN ? AND ? AND status = ?", startDate, endDate, "completed").Count(&completed)
	r.db.Model(&models.Appointment{}).Where("appointment_date BETWEEN ? AND ? AND status = ?", startDate, endDate, "cancelled").Count(&cancelled)
	r.db.Model(&models.Appointment{}).Where("appointment_date BETWEEN ? AND ? AND status = ?", startDate, endDate, "no_show").Count(&noShow)
	completionRate := 0.0
	if total > 0 {
		completionRate = float64(completed) / float64(total) * 100
	}

	var totalActive, newReg int64
	r.db.Model(&models.Patient{}).Count(&totalActive)
	r.db.Model(&models.Patient{}).Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate).Count(&newReg)

	var totalAdmissions, currentlyHosp, totalDischarged int64
	var availableBeds int64
	r.db.Model(&models.Hospitalization{}).Where("admission_date BETWEEN ? AND ?", startDate, endDate).Count(&totalAdmissions)
	r.db.Model(&models.Hospitalization{}).Where("status = ?", "admitted").Count(&currentlyHosp)
	r.db.Model(&models.Hospitalization{}).Where("status = ? AND admission_date BETWEEN ? AND ?", "discharged", startDate, endDate).Count(&totalDischarged)
	r.db.Model(&models.Room{}).Where("available_beds > 0 AND is_active = ?", true).Count(&availableBeds)

	// Top departments
	type deptCount struct {
		DepartmentID   uint   `gorm:"column:department_id"`
		DepartmentName string `gorm:"column:department_name"`
		Count          int64  `gorm:"column:count"`
	}
	var topDepts []deptCount
	r.db.Table("appointments a").
		Select("d.id as department_id, d.name as department_name, COUNT(a.id) as count").
		Joins("JOIN doctors doc ON a.doctor_id = doc.id").
		Joins("JOIN departments d ON doc.department_id = d.id").
		Where("a.appointment_date BETWEEN ? AND ?", startDate, endDate).
		Group("d.id, d.name").
		Order("count DESC").
		Limit(5).
		Scan(&topDepts)

	depts := make([]dto.DashboardTopDepartment, len(topDepts))
	for i, d := range topDepts {
		depts[i] = dto.DashboardTopDepartment{DepartmentID: d.DepartmentID, DepartmentName: d.DepartmentName, Count: d.Count}
	}

	// Appointment trend
	type trendRow struct {
		Date  string `gorm:"column:date"`
		Count int64  `gorm:"column:count"`
	}
	var trendRows []trendRow
	r.db.Table("appointments").
		Select("appointment_date as date, COUNT(id) as count").
		Where("appointment_date BETWEEN ? AND ?", startDate, endDate).
		Group("appointment_date").
		Order("appointment_date ASC").
		Scan(&trendRows)

	trend := make([]dto.DashboardTrendPoint, len(trendRows))
	for i, t := range trendRows {
		trend[i] = dto.DashboardTrendPoint{Date: t.Date, Count: t.Count}
	}

	return &dto.DashboardAdminResponse{
		Period:      query.Period,
		PeriodRange: dto.DashboardPeriodRange{Start: startDate, End: endDate},
		Appointments: dto.DashboardAdminAppointments{
			Total:          total,
			Completed:      completed,
			Cancelled:      cancelled,
			NoShow:         noShow,
			CompletionRate: completionRate,
		},
		Patients: dto.DashboardAdminPatients{
			TotalActive:      totalActive,
			NewRegistrations: newReg,
			Returning:        total - newReg,
		},
		Hospitalization: dto.DashboardAdminHospitalization{
			TotalAdmissions:       totalAdmissions,
			CurrentlyHospitalized: currentlyHosp,
			TotalDischarged:       totalDischarged,
			AvailableBeds:         availableBeds,
		},
		Billing:          dto.DashboardAdminBilling{},
		Referrals:        dto.DashboardAdminReferrals{},
		TopDepartments:   depts,
		AppointmentTrend: trend,
		RevenueTrend:     []dto.DashboardRevenueTrendPoint{},
	}, nil
}

func (r *dashboardRepository) GetDoctorDashboard(doctorID uint, date string) (*dto.DashboardDoctorResponse, error) {
	// Doctor info
	var doctor models.Doctor
	r.db.Preload("Specialization").First(&doctor, doctorID)

	var total, completed, pending, cancelled int64
	r.db.Model(&models.Appointment{}).Where("doctor_id = ? AND appointment_date = ?", doctorID, date).Count(&total)
	r.db.Model(&models.Appointment{}).Where("doctor_id = ? AND appointment_date = ? AND status = ?", doctorID, date, "completed").Count(&completed)
	r.db.Model(&models.Appointment{}).Where("doctor_id = ? AND appointment_date = ? AND status IN (?)", doctorID, date, []string{"scheduled", "confirmed"}).Count(&pending)
	r.db.Model(&models.Appointment{}).Where("doctor_id = ? AND appointment_date = ? AND status = ?", doctorID, date, "cancelled").Count(&cancelled)

	// Next patient
	var nextAppt models.Appointment
	var nextPatient *dto.DashboardDoctorNextPatient
	err := r.db.Preload("Patient").
		Where("doctor_id = ? AND appointment_date = ? AND status IN (?)", doctorID, date, []string{"scheduled", "confirmed"}).
		Order("appointment_time ASC").
		First(&nextAppt).Error
	if err == nil && nextAppt.Patient != nil {
		nextPatient = &dto.DashboardDoctorNextPatient{
			AppointmentID:   nextAppt.ID,
			AppointmentTime: nextAppt.AppointmentTime,
			PatientName:     nextAppt.Patient.FullName,
		}
	}

	// Upcoming
	var upcomingAppts []models.Appointment
	r.db.Preload("Patient").
		Where("doctor_id = ? AND appointment_date = ? AND status IN (?)", doctorID, date, []string{"scheduled", "confirmed"}).
		Order("appointment_time ASC").Limit(10).Find(&upcomingAppts)
	upcoming := make([]dto.DashboardDoctorAppointment, len(upcomingAppts))
	for i, a := range upcomingAppts {
		item := dto.DashboardDoctorAppointment{AppointmentID: a.ID, AppointmentTime: a.AppointmentTime}
		if a.Patient != nil {
			item.PatientName = a.Patient.FullName
			item.PatientID = a.Patient.ID
		}
		upcoming[i] = item
	}

	// Stats
	now := time.Now()
	firstOfMonth := fmt.Sprintf("%d-%02d-01", now.Year(), now.Month())
	var totalThisMonth, pendingLab, pendingRx int64
	r.db.Model(&models.Appointment{}).Where("doctor_id = ? AND appointment_date >= ?", doctorID, firstOfMonth).Count(&totalThisMonth)
	r.db.Table("lab_tests").Where("ordered_by_doctor_id = ? AND status = ?", doctorID, "ordered").Count(&pendingLab)
	r.db.Model(&models.Prescription{}).Where("doctor_id = ? AND status = ?", doctorID, "pending").Count(&pendingRx)

	// Recent records
	var recentRecords []models.MedicalRecord
	r.db.Preload("Patient").Where("doctor_id = ?", doctorID).Order("created_at DESC").Limit(5).Find(&recentRecords)
	records := make([]dto.DashboardDoctorRecord, len(recentRecords))
	for i, mr := range recentRecords {
		rec := dto.DashboardDoctorRecord{ID: mr.ID, VisitDate: mr.VisitDate, ChiefComplaint: mr.ChiefComplaint, Status: mr.Status}
		if mr.Patient != nil {
			rec.PatientName = mr.Patient.FullName
		}
		records[i] = rec
	}

	specName := ""
	if doctor.Specialization.ID != 0 {
		specName = doctor.Specialization.Name
	}

	return &dto.DashboardDoctorResponse{
		Doctor:               dto.DashboardDoctorInfo{ID: doctor.ID, Name: doctor.FullName, Specialization: specName},
		Date:                 date,
		TodaySchedule:        dto.DashboardDoctorSchedule{TotalAppointments: total, Completed: completed, Pending: pending, Cancelled: cancelled, NextPatient: nextPatient},
		UpcomingAppointments: upcoming,
		Statistics:           dto.DashboardDoctorStats{TotalPatientsThisMonth: totalThisMonth, TotalPatientsToday: total, PendingLabResults: pendingLab, PendingPrescriptions: pendingRx, PendingReferrals: 0},
		RecentMedicalRecords: records,
	}, nil
}

func (r *dashboardRepository) GetReceptionistDashboard(date string) (*dto.DashboardReceptionistResponse, error) {
	var total, pending, confirmed, inProgress, completed, cancelled int64
	r.db.Model(&models.Appointment{}).Where("appointment_date = ?", date).Count(&total)
	r.db.Model(&models.Appointment{}).Where("appointment_date = ? AND status = ?", date, "scheduled").Count(&pending)
	r.db.Model(&models.Appointment{}).Where("appointment_date = ? AND status = ?", date, "confirmed").Count(&confirmed)
	r.db.Model(&models.Appointment{}).Where("appointment_date = ? AND status = ?", date, "in_progress").Count(&inProgress)
	r.db.Model(&models.Appointment{}).Where("appointment_date = ? AND status = ?", date, "completed").Count(&completed)
	r.db.Model(&models.Appointment{}).Where("appointment_date = ? AND status = ?", date, "cancelled").Count(&cancelled)

	// Queue
	var queueAppts []models.Appointment
	r.db.Preload("Patient").Preload("Doctor").
		Where("appointment_date = ? AND status IN (?)", date, []string{"scheduled", "confirmed", "in_progress"}).
		Order("appointment_time ASC").Limit(20).Find(&queueAppts)
	queue := make([]dto.DashboardAppointmentQueueItem, len(queueAppts))
	for i, a := range queueAppts {
		item := dto.DashboardAppointmentQueueItem{AppointmentID: a.ID, ScheduledTime: a.AppointmentTime, Status: a.Status}
		if a.Patient != nil {
			item.PatientName = a.Patient.FullName
			item.PatientID = a.Patient.ID
		}
		if a.Doctor != nil {
			item.DoctorName = a.Doctor.FullName
		}
		queue[i] = item
	}

	var currentHosp, availBeds, newAdmissions, discharged int64
	r.db.Model(&models.Hospitalization{}).Where("status = ?", "admitted").Count(&currentHosp)
	r.db.Model(&models.Room{}).Where("available_beds > 0 AND is_active = ?", true).Count(&availBeds)
	r.db.Model(&models.Hospitalization{}).Where("admission_date = ?", date).Count(&newAdmissions)
	r.db.Model(&models.Hospitalization{}).Where("discharge_date = ? AND status = ?", date, "discharged").Count(&discharged)

	var newPatients int64
	r.db.Model(&models.Patient{}).Where("DATE(created_at) = ?", date).Count(&newPatients)

	return &dto.DashboardReceptionistResponse{
		Date: date,
		AppointmentsToday: dto.DashboardReceptionistAppointments{
			Total: total, Pending: pending, Confirmed: confirmed,
			InProgress: inProgress, Completed: completed, Cancelled: cancelled,
			AppointmentQueue: queue,
		},
		Hospitalization: dto.DashboardReceptionistHospitalization{
			CurrentlyHospitalized: currentHosp, AvailableBeds: availBeds,
			NewAdmissionsToday: newAdmissions, DischargedToday: discharged,
		},
		Billing:                      dto.DashboardReceptionistBilling{},
		NewPatientRegistrationsToday: newPatients,
	}, nil
}

func (r *dashboardRepository) GetPatientDashboard(patientID uint) (*dto.DashboardPatientResponse, error) {
	var patient models.Patient
	r.db.First(&patient, patientID)

	today := time.Now().Format("2006-01-02")

	// Upcoming appointments
	var upcomingAppts []models.Appointment
	r.db.Preload("Doctor").Preload("Doctor.Specialization").Preload("Doctor.Department").
		Where("patient_id = ? AND appointment_date >= ? AND status IN (?)", patientID, today, []string{"scheduled", "confirmed"}).
		Order("appointment_date ASC, appointment_time ASC").Limit(5).Find(&upcomingAppts)
	upcoming := make([]dto.DashboardPatientAppointment, len(upcomingAppts))
	for i, a := range upcomingAppts {
		item := dto.DashboardPatientAppointment{ID: a.ID, ScheduledDate: a.AppointmentDate, ScheduledTime: a.AppointmentTime, Status: a.Status}
		if a.Doctor != nil {
			item.DoctorName = a.Doctor.FullName
			if a.Doctor.Department.ID != 0 {
				item.DepartmentName = a.Doctor.Department.Name
			}
		}
		upcoming[i] = item
	}

	// Active prescriptions
	var prescriptions []models.Prescription
	r.db.Preload("Doctor").
		Joins("JOIN medical_records ON prescriptions.medical_record_id = medical_records.id").
		Where("medical_records.patient_id = ? AND prescriptions.status IN (?)", patientID, []string{"active", "pending"}).
		Order("prescriptions.created_at DESC").Limit(5).Find(&prescriptions)
	rxList := make([]dto.DashboardPatientPrescription, len(prescriptions))
	for i, rx := range prescriptions {
		item := dto.DashboardPatientPrescription{ID: rx.ID, Status: rx.Status, PrescribedDate: rx.CreatedAt.Format("2006-01-02")}
		if rx.Doctor != nil {
			item.DoctorName = rx.Doctor.FullName
		}
		rxList[i] = item
	}

	// Pending lab results
	var labTests []models.LabTest
	r.db.Joins("JOIN medical_records ON lab_tests.medical_record_id = medical_records.id").
		Where("medical_records.patient_id = ? AND lab_tests.status IN (?)", patientID, []string{"ordered", "processing"}).
		Order("lab_tests.created_at DESC").Limit(5).Find(&labTests)
	labs := make([]dto.DashboardPatientLabResult, len(labTests))
	for i, lt := range labTests {
		labs[i] = dto.DashboardPatientLabResult{ID: lt.ID, Status: lt.Status, OrderedDate: lt.CreatedAt.Format("2006-01-02")}
	}

	// Recent visits
	var recentRecords []models.MedicalRecord
	r.db.Preload("Doctor").Where("patient_id = ?", patientID).Order("visit_date DESC").Limit(5).Find(&recentRecords)
	visits := make([]dto.DashboardPatientVisit, len(recentRecords))
	for i, mr := range recentRecords {
		v := dto.DashboardPatientVisit{ID: mr.ID, VisitDate: mr.VisitDate, ChiefComplaint: mr.ChiefComplaint, Status: mr.Status}
		if mr.Doctor != nil {
			v.DoctorName = mr.Doctor.FullName
		}
		visits[i] = v
	}

	age := 0
	if patient.DateOfBirth != "" {
		if dob, err := time.Parse("2006-01-02", patient.DateOfBirth); err == nil {
			age = int(time.Since(dob).Hours() / 24 / 365)
		}
	}

	return &dto.DashboardPatientResponse{
		Patient:              dto.DashboardPatientInfo{ID: patient.ID, Name: patient.FullName, MedicalRecordNumber: patient.PatientCode, DateOfBirth: patient.DateOfBirth, Age: age},
		UpcomingAppointments: upcoming,
		Billing:              dto.DashboardPatientBilling{UnpaidBills: []dto.DashboardPatientUnpaidBill{}},
		ActivePrescriptions:  rxList,
		PendingLabResults:    labs,
		ActiveReferrals:      []dto.DashboardPatientReferral{},
		RecentVisits:         visits,
	}, nil
}

func (r *dashboardRepository) GetAppointmentReport(query *dto.DashboardAppointmentReportQuery, startDate, endDate string) (*dto.DashboardAppointmentReportResponse, error) {
	baseQuery := r.db.Model(&models.Appointment{}).Where("appointment_date BETWEEN ? AND ?", startDate, endDate)
	if query.DoctorID != nil {
		baseQuery = baseQuery.Where("doctor_id = ?", *query.DoctorID)
	}
	if query.DepartmentID != nil {
		baseQuery = baseQuery.Joins("JOIN doctors ON appointments.doctor_id = doctors.id").Where("doctors.department_id = ?", *query.DepartmentID)
	}

	var total, scheduled, confirmed, completed, cancelled, noShow int64
	baseQuery.Count(&total)
	baseQuery.Where("status = ?", "scheduled").Count(&scheduled)
	baseQuery.Where("status = ?", "confirmed").Count(&confirmed)
	baseQuery.Where("status = ?", "completed").Count(&completed)
	baseQuery.Where("status = ?", "cancelled").Count(&cancelled)
	baseQuery.Where("status = ?", "no_show").Count(&noShow)

	completionRate, cancellationRate := 0.0, 0.0
	if total > 0 {
		completionRate = float64(completed) / float64(total) * 100
		cancellationRate = float64(cancelled) / float64(total) * 100
	}

	return &dto.DashboardAppointmentReportResponse{
		Period:      query.Period,
		PeriodRange: dto.DashboardPeriodRange{Start: startDate, End: endDate},
		Totals: dto.DashboardAppointmentReportTotals{
			Total: total, Scheduled: scheduled, Confirmed: confirmed,
			Completed: completed, Cancelled: cancelled, NoShow: noShow,
			CompletionRate: completionRate, CancellationRate: cancellationRate,
		},
		ByDepartment: []dto.DashboardDeptStats{},
		ByDoctor:     []dto.DashboardDoctorStatItem{},
		Trend:        []dto.DashboardAppointmentTrendItem{},
		PeakHours:    []dto.DashboardPeakHour{},
	}, nil
}

func (r *dashboardRepository) GetRevenueReport(query *dto.DashboardRevenueReportQuery, startDate, endDate string) (*dto.DashboardRevenueReportResponse, error) {
	return &dto.DashboardRevenueReportResponse{
		Period:      query.Period,
		PeriodRange: dto.DashboardPeriodRange{Start: startDate, End: endDate},
		Revenue:     dto.DashboardRevenueTotals{},
		ByCategory:  []dto.DashboardRevenueCategory{},
		Trend:       []dto.DashboardRevenueTrendPoint{},
		Comparison:  dto.DashboardRevenueComparison{},
	}, nil
}

func (r *dashboardRepository) GetPatientReport(query *dto.DashboardPatientReportQuery, startDate, endDate string) (*dto.DashboardPatientReportResponse, error) {
	var newPatients, totalVisits int64
	r.db.Model(&models.Patient{}).Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate).Count(&newPatients)
	r.db.Model(&models.Appointment{}).Where("appointment_date BETWEEN ? AND ? AND status = ?", startDate, endDate, "completed").Count(&totalVisits)

	// Gender distribution
	type genderCount struct {
		Gender string `gorm:"column:gender"`
		Count  int64  `gorm:"column:count"`
	}
	var genderRows []genderCount
	r.db.Model(&models.Patient{}).Select("gender, COUNT(id) as count").Group("gender").Scan(&genderRows)
	genders := make([]dto.DashboardGenderCount, len(genderRows))
	for i, g := range genderRows {
		genders[i] = dto.DashboardGenderCount{Gender: g.Gender, Count: g.Count}
	}

	return &dto.DashboardPatientReportResponse{
		Period:      query.Period,
		PeriodRange: dto.DashboardPeriodRange{Start: startDate, End: endDate},
		Registrations: dto.DashboardPatientRegistrations{
			NewPatients: newPatients, TotalVisits: totalVisits,
		},
		Demographics:      dto.DashboardPatientDemographics{ByGender: genders, ByAgeGroup: []dto.DashboardAgeGroupCount{}},
		RegistrationTrend: []dto.DashboardPatientTrendPoint{},
		Comparison:        dto.DashboardPatientComparison{},
	}, nil
}
