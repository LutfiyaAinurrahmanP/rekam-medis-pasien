package appointment

import (
	"context"
	"log"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/cache"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
)

type cachedAppointmentService struct {
	inner AppointmentService
	redis *cache.RedisClient
}

func NewCachedAppointmentService(inner AppointmentService, redisClient *cache.RedisClient) AppointmentService {
	if redisClient == nil {
		return inner
	}
	return &cachedAppointmentService{
		inner: inner,
		redis: redisClient,
	}
}

func (s *cachedAppointmentService) List(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	key := cache.AppointmentListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		query.DoctorID,
		query.DepartmentID,
		normalizeCachePart(query.Status),
		normalizeCachePart(query.Date),
		normalizeCachePart(query.DateFrom),
		normalizeCachePart(query.DateTo),
		query.DaysAhead,
		query.DaysBack,
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.AppointmentListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.List(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedAppointmentService) DeletedList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	key := cache.AppointmentDeletedListQueryKey(
		query.Page,
		query.PageSize,
		query.PatientID,
		query.DoctorID,
		query.DepartmentID,
		normalizeCachePart(query.Status),
		normalizeCachePart(query.Date),
		normalizeCachePart(query.DateFrom),
		normalizeCachePart(query.DateTo),
		query.DaysAhead,
		query.DaysBack,
		normalizeCachePart(query.SortBy),
		normalizeCachePart(query.SortDir),
	)
	var resp dto.AppointmentListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.DeletedList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedAppointmentService) UpcomingList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	q := *query
	if q.DaysAhead == 0 {
		q.DaysAhead = 7
	}
	key := cache.AppointmentListQueryKey(
		q.Page,
		q.PageSize,
		q.PatientID,
		q.DoctorID,
		q.DepartmentID,
		normalizeCachePart(q.Status),
		normalizeCachePart(q.Date),
		normalizeCachePart(q.DateFrom),
		normalizeCachePart(q.DateTo),
		q.DaysAhead,
		q.DaysBack,
		normalizeCachePart(q.SortBy),
		normalizeCachePart(q.SortDir),
	)
	var resp dto.AppointmentListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.UpcomingList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedAppointmentService) PastList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	q := *query
	if q.DaysBack == 0 {
		q.DaysBack = 30
	}
	key := cache.AppointmentListQueryKey(
		q.Page,
		q.PageSize,
		q.PatientID,
		q.DoctorID,
		q.DepartmentID,
		normalizeCachePart(q.Status),
		normalizeCachePart(q.Date),
		normalizeCachePart(q.DateFrom),
		normalizeCachePart(q.DateTo),
		q.DaysAhead,
		q.DaysBack,
		normalizeCachePart(q.SortBy),
		normalizeCachePart(q.SortDir),
	)
	var resp dto.AppointmentListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.PastList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedAppointmentService) TodayList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	q := *query
	// Note: We use a placeholder "today" for cache key so it naturally invalidates daily or doesn't cause key explosion
	key := cache.AppointmentListQueryKey(
		q.Page,
		q.PageSize,
		q.PatientID,
		q.DoctorID,
		q.DepartmentID,
		normalizeCachePart(q.Status),
		"today",
		normalizeCachePart(q.DateFrom),
		normalizeCachePart(q.DateTo),
		q.DaysAhead,
		q.DaysBack,
		normalizeCachePart(q.SortBy),
		normalizeCachePart(q.SortDir),
	)
	var resp dto.AppointmentListResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.TodayList(query)
	if err != nil {
		return nil, err
	}
	s.setCache(key, result)
	return result, nil
}

func (s *cachedAppointmentService) FindByID(id uint) (*dto.AppointmentResponse, error) {
	key := cache.AppointmentKey(id)
	var resp dto.AppointmentResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.FindByID(id)
	if err != nil {
		return nil, err
	}

	s.setCache(key, result)
	return result, nil
}

func (s *cachedAppointmentService) FindByIDUnscoped(id uint) (*dto.AppointmentResponse, error) {
	key := cache.AppointmentKey(id)
	var resp dto.AppointmentResponse
	if err := s.redis.Get(context.Background(), key, &resp); err == nil {
		return &resp, nil
	}

	result, err := s.inner.FindByIDUnscoped(id)
	if err != nil {
		return nil, err
	}

	s.setCache(key, result)
	return result, nil
}

func (s *cachedAppointmentService) Create(req *dto.CreateAppointmentRequest) (*dto.AppointmentResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedAppointmentService) Update(id uint, req *dto.UpdateAppointmentRequest) (*dto.AppointmentResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.invalidateAll()
	return result, nil
}

func (s *cachedAppointmentService) Confirm(id uint) error {
	if err := s.inner.Confirm(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedAppointmentService) Start(id uint) error {
	if err := s.inner.Start(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedAppointmentService) Complete(id uint) error {
	if err := s.inner.Complete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedAppointmentService) Cancel(id uint, req *dto.CancelAppointmentRequest) error {
	if err := s.inner.Cancel(id, req); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedAppointmentService) Reschedule(id uint, req *dto.RescheduleAppointmentRequest) error {
	if err := s.inner.Reschedule(id, req); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedAppointmentService) NoShow(id uint) error {
	if err := s.inner.NoShow(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedAppointmentService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedAppointmentService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func (s *cachedAppointmentService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	s.invalidateAll()
	return nil
}

func normalizeCachePart(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func (s *cachedAppointmentService) setCache(key string, value any) {
	if err := s.redis.Set(context.Background(), key, value, 0); err != nil {
		log.Printf("⚠️  Redis set failed for key %q: %v", key, err)
	}
}

func (s *cachedAppointmentService) invalidateAll() {
	if err := s.redis.DeleteByPattern(context.Background(), cache.PatternAppointmentAll); err != nil {
		log.Printf("⚠️  Redis invalidate failed for pattern %q: %v", cache.PatternAppointmentAll, err)
	}
}
