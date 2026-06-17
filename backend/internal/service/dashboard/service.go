package dashboard

import (
	"fmt"
	"strings"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type DashboardService interface {
	GetOverview(query *dto.DashboardOverviewQuery) (*dto.DashboardOverviewResponse, error)
	GetAdminDashboard(query *dto.DashboardPeriodQuery) (*dto.DashboardAdminResponse, error)
	GetDoctorDashboard(doctorID uint, query *dto.DashboardDoctorQuery) (*dto.DashboardDoctorResponse, error)
	GetReceptionistDashboard() (*dto.DashboardReceptionistResponse, error)
	GetPatientDashboard(patientID uint) (*dto.DashboardPatientResponse, error)
	GetAppointmentReport(query *dto.DashboardAppointmentReportQuery) (*dto.DashboardAppointmentReportResponse, error)
	GetRevenueReport(query *dto.DashboardRevenueReportQuery) (*dto.DashboardRevenueReportResponse, error)
	GetPatientReport(query *dto.DashboardPatientReportQuery) (*dto.DashboardPatientReportResponse, error)
}

type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) resolveDateRange(period, startDate, endDate string) (string, string) {
	now := time.Now()
	today := now.Format("2006-01-02")

	switch period {
	case "today":
		return today, today
	case "this_week":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := now.AddDate(0, 0, -(weekday - 1))
		end := start.AddDate(0, 0, 6)
		return start.Format("2006-01-02"), end.Format("2006-01-02")
	case "this_month":
		return fmt.Sprintf("%d-%02d-01", now.Year(), now.Month()),
			fmt.Sprintf("%d-%02d-%02d", now.Year(), now.Month(), daysInMonth(now.Year(), now.Month()))
	case "last_month":
		last := now.AddDate(0, -1, 0)
		return fmt.Sprintf("%d-%02d-01", last.Year(), last.Month()),
			fmt.Sprintf("%d-%02d-%02d", last.Year(), last.Month(), daysInMonth(last.Year(), last.Month()))
	case "custom":
		if startDate != "" && endDate != "" {
			return startDate, endDate
		}
	}
	return today, today
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func normalizeCachePart(s string) string {
	if s == "" {
		return "all"
	}
	return strings.ToLower(strings.TrimSpace(s))
}

func (s *dashboardService) GetOverview(query *dto.DashboardOverviewQuery) (*dto.DashboardOverviewResponse, error) {
	date := query.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return s.repo.GetOverview(date)
}

func (s *dashboardService) GetAdminDashboard(query *dto.DashboardPeriodQuery) (*dto.DashboardAdminResponse, error) {
	if query.Period == "" {
		query.Period = "today"
	}
	startDate, endDate := s.resolveDateRange(query.Period, query.StartDate, query.EndDate)
	return s.repo.GetAdminDashboard(query, startDate, endDate)
}

func (s *dashboardService) GetDoctorDashboard(doctorID uint, query *dto.DashboardDoctorQuery) (*dto.DashboardDoctorResponse, error) {
	date := query.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	return s.repo.GetDoctorDashboard(doctorID, date)
}

func (s *dashboardService) GetReceptionistDashboard() (*dto.DashboardReceptionistResponse, error) {
	date := time.Now().Format("2006-01-02")
	return s.repo.GetReceptionistDashboard(date)
}

func (s *dashboardService) GetPatientDashboard(patientID uint) (*dto.DashboardPatientResponse, error) {
	return s.repo.GetPatientDashboard(patientID)
}

func (s *dashboardService) GetAppointmentReport(query *dto.DashboardAppointmentReportQuery) (*dto.DashboardAppointmentReportResponse, error) {
	if query.Period == "" {
		query.Period = "this_month"
	}
	startDate, endDate := s.resolveDateRange(query.Period, query.StartDate, query.EndDate)
	return s.repo.GetAppointmentReport(query, startDate, endDate)
}

func (s *dashboardService) GetRevenueReport(query *dto.DashboardRevenueReportQuery) (*dto.DashboardRevenueReportResponse, error) {
	if query.Period == "" {
		query.Period = "this_month"
	}
	startDate, endDate := s.resolveDateRange(query.Period, query.StartDate, query.EndDate)
	return s.repo.GetRevenueReport(query, startDate, endDate)
}

func (s *dashboardService) GetPatientReport(query *dto.DashboardPatientReportQuery) (*dto.DashboardPatientReportResponse, error) {
	if query.Period == "" {
		query.Period = "this_month"
	}
	startDate, endDate := s.resolveDateRange(query.Period, query.StartDate, query.EndDate)
	return s.repo.GetPatientReport(query, startDate, endDate)
}
