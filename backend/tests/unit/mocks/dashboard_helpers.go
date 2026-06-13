package mocks

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

func NewTestDashboardOverviewResponse() *dto.DashboardOverviewResponse {
	return &dto.DashboardOverviewResponse{
		SummaryDate: "2024-01-01",
		MasterData: dto.DashboardMasterData{
			TotalPatients:    100,
			TotalDoctors:     10,
			TotalDepartments: 5,
			TotalRooms:       20,
			TotalMedicines:   50,
		},
	}
}

func NewTestDashboardAdminResponse() *dto.DashboardAdminResponse {
	return &dto.DashboardAdminResponse{
		Period: "today",
		Appointments: dto.DashboardAdminAppointments{
			Total:     50,
			Completed: 40,
		},
	}
}

func NewTestDashboardDoctorResponse() *dto.DashboardDoctorResponse {
	return &dto.DashboardDoctorResponse{
		Doctor: dto.DashboardDoctorInfo{
			ID:   1,
			Name: "Dr. Andi",
		},
		TodaySchedule: dto.DashboardDoctorSchedule{
			TotalAppointments: 5,
		},
	}
}

func NewTestDashboardReceptionistResponse() *dto.DashboardReceptionistResponse {
	return &dto.DashboardReceptionistResponse{
		Date: "2024-01-01",
		AppointmentsToday: dto.DashboardReceptionistAppointments{
			Total: 20,
		},
	}
}

func NewTestDashboardPatientResponse() *dto.DashboardPatientResponse {
	return &dto.DashboardPatientResponse{
		Patient: dto.DashboardPatientInfo{
			ID:   1,
			Name: "Budi",
		},
	}
}

func NewTestDashboardAppointmentReportResponse() *dto.DashboardAppointmentReportResponse {
	return &dto.DashboardAppointmentReportResponse{
		Totals: dto.DashboardAppointmentReportTotals{
			Total: 100,
		},
	}
}

func NewTestDashboardRevenueReportResponse() *dto.DashboardRevenueReportResponse {
	return &dto.DashboardRevenueReportResponse{
		Revenue: dto.DashboardRevenueTotals{
			TotalBilled: 1000000,
		},
	}
}

func NewTestDashboardPatientReportResponse() *dto.DashboardPatientReportResponse {
	return &dto.DashboardPatientReportResponse{
		Registrations: dto.DashboardPatientRegistrations{
			NewPatients: 10,
		},
	}
}
