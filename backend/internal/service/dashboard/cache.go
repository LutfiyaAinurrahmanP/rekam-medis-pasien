package dashboard

import (
	"context"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedDashboardService struct {
	inner DashboardService
	redis *cache.RedisClient
}

func NewCachedDashboardService(inner DashboardService, redisClient *cache.RedisClient) DashboardService {
	if redisClient == nil {
		return inner
	}
	return &cachedDashboardService{
		inner: inner,
		redis: redisClient,
	}
}

const dashboardCacheTTL = 2 * time.Minute

func (s *cachedDashboardService) GetOverview(query *dto.DashboardOverviewQuery) (*dto.DashboardOverviewResponse, error) {
	date := normalizeCachePart(query.Date)
	key := cache.DashboardOverviewKey(date)
	var resp dto.DashboardOverviewResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetOverview(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result, dashboardCacheTTL)
	return result, nil
}

func (s *cachedDashboardService) GetAdminDashboard(query *dto.DashboardPeriodQuery) (*dto.DashboardAdminResponse, error) {
	key := cache.DashboardAdminKey(
		normalizeCachePart(query.Period),
		normalizeCachePart(query.StartDate),
		normalizeCachePart(query.EndDate),
	)
	var resp dto.DashboardAdminResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetAdminDashboard(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result, dashboardCacheTTL)
	return result, nil
}

func (s *cachedDashboardService) GetDoctorDashboard(doctorID uint, query *dto.DashboardDoctorQuery) (*dto.DashboardDoctorResponse, error) {
	key := cache.DashboardDoctorKey(doctorID, normalizeCachePart(query.Date))
	var resp dto.DashboardDoctorResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetDoctorDashboard(doctorID, query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result, dashboardCacheTTL)
	return result, nil
}

func (s *cachedDashboardService) GetReceptionistDashboard() (*dto.DashboardReceptionistResponse, error) {
	key := cache.DashboardReceptionistKey("today")
	var resp dto.DashboardReceptionistResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetReceptionistDashboard()
	if err != nil {
		return nil, err
	}
	s.setCache(key, result, dashboardCacheTTL)
	return result, nil
}

func (s *cachedDashboardService) GetPatientDashboard(patientID uint) (*dto.DashboardPatientResponse, error) {
	key := cache.DashboardPatientKey(patientID)
	var resp dto.DashboardPatientResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetPatientDashboard(patientID)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result, dashboardCacheTTL)
	return result, nil
}

func (s *cachedDashboardService) GetAppointmentReport(query *dto.DashboardAppointmentReportQuery) (*dto.DashboardAppointmentReportResponse, error) {
	key := cache.DashboardAppointmentReportKey(
		normalizeCachePart(query.Period),
		normalizeCachePart(query.StartDate),
		normalizeCachePart(query.EndDate),
		query.DoctorID,
		query.DepartmentID,
		normalizeCachePart(query.GroupBy),
	)
	var resp dto.DashboardAppointmentReportResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetAppointmentReport(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result, dashboardCacheTTL)
	return result, nil
}

func (s *cachedDashboardService) GetRevenueReport(query *dto.DashboardRevenueReportQuery) (*dto.DashboardRevenueReportResponse, error) {
	key := cache.DashboardRevenueReportKey(
		normalizeCachePart(query.Period),
		normalizeCachePart(query.StartDate),
		normalizeCachePart(query.EndDate),
		normalizeCachePart(query.GroupBy),
	)
	var resp dto.DashboardRevenueReportResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetRevenueReport(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result, dashboardCacheTTL)
	return result, nil
}

func (s *cachedDashboardService) GetPatientReport(query *dto.DashboardPatientReportQuery) (*dto.DashboardPatientReportResponse, error) {
	key := cache.DashboardPatientReportKey(
		normalizeCachePart(query.Period),
		normalizeCachePart(query.StartDate),
		normalizeCachePart(query.EndDate),
	)
	var resp dto.DashboardPatientReportResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}
	result, err := s.inner.GetPatientReport(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result, dashboardCacheTTL)
	return result, nil
}

func (s *cachedDashboardService) setCache(key string, value any, ttl time.Duration) {
	if err := s.redis.Set(context.Background(), key, value, ttl); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}
