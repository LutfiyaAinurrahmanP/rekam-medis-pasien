package appointment

import (
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type AppointmentService interface {
	List(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error)
	DeletedList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error)
	UpcomingList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error)
	PastList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error)
	TodayList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error)
	FindByID(id uint) (*dto.AppointmentResponse, error)
	FindByIDUnscoped(id uint) (*dto.AppointmentResponse, error)
	Create(req *dto.CreateAppointmentRequest) (*dto.AppointmentResponse, error)
	Update(id uint, req *dto.UpdateAppointmentRequest) (*dto.AppointmentResponse, error)
	Confirm(id uint) error
	Start(id uint) error
	Complete(id uint) error
	Cancel(id uint, req *dto.CancelAppointmentRequest) error
	Reschedule(id uint, req *dto.RescheduleAppointmentRequest) error
	NoShow(id uint) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type appointmentService struct {
	repo   repository.AppointmentRepository
	config *config.Config
}

func NewAppointmentService(repo repository.AppointmentRepository, config *config.Config) AppointmentService {
	return &appointmentService{
		repo:   repo,
		config: config,
	}
}

func (s *appointmentService) normalizeQuery(query *dto.AppointmentPaginationQuery, defaultSortBy, defaultSortDir string) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = s.config.Pagination.DefaultPageSize
	}

	if query.PageSize > s.config.Pagination.MaxPageSize {
		query.PageSize = s.config.Pagination.MaxPageSize
	}

	if query.SortBy == "" {
		query.SortBy = defaultSortBy
	}

	if query.SortDir == "" {
		query.SortDir = defaultSortDir
	}
}

func (s *appointmentService) toResponse(m *models.Appointment) *dto.AppointmentResponse {
	if m == nil {
		return nil
	}

	var patientResp *dto.AppointmentPatientResponse
	if m.Patient != nil {
		patientResp = &dto.AppointmentPatientResponse{
			ID:          m.Patient.ID,
			PatientCode: m.Patient.PatientCode,
			FullName:    m.Patient.FullName,
			Phone:       m.Patient.Phone,
		}
	}

	var doctorResp *dto.AppointmentDoctorResponse
	if m.Doctor != nil {
		specName := m.Doctor.Specialization.Name
		deptName := m.Doctor.Department.Name

		doctorResp = &dto.AppointmentDoctorResponse{
			ID:             m.Doctor.ID,
			FullName:       m.Doctor.FullName,
			Specialization: specName,
			Department:     deptName,
		}
	}

	return &dto.AppointmentResponse{
		ID:              m.ID,
		PatientID:       m.PatientID,
		Patient:         patientResp,
		DoctorID:        m.DoctorID,
		Doctor:          doctorResp,
		AppointmentDate: m.AppointmentDate,
		AppointmentTime: m.AppointmentTime,
		DurationMinutes: m.DurationMinutes,
		Status:          m.Status,
		Reason:          m.Reason,
		Notes:           m.Notes,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func (s *appointmentService) List(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	s.normalizeQuery(query, "appointment_date", "asc")

	appointments, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AppointmentResponse, len(appointments))
	for i, a := range appointments {
		responses[i] = *s.toResponse(&a)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.AppointmentListResponse{
		Data: responses,
		Meta: dto.AppointmentPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *appointmentService) DeletedList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	s.normalizeQuery(query, "appointment_date", "desc")

	appointments, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AppointmentResponse, len(appointments))
	for i, a := range appointments {
		responses[i] = *s.toResponse(&a)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.AppointmentListResponse{
		Data: responses,
		Meta: dto.AppointmentPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *appointmentService) UpcomingList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	s.normalizeQuery(query, "appointment_date", "asc")

	appointments, total, err := s.repo.UpcomingList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AppointmentResponse, len(appointments))
	for i, a := range appointments {
		responses[i] = *s.toResponse(&a)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.AppointmentListResponse{
		Data: responses,
		Meta: dto.AppointmentPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *appointmentService) PastList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	s.normalizeQuery(query, "appointment_date", "desc")

	appointments, total, err := s.repo.PastList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AppointmentResponse, len(appointments))
	for i, a := range appointments {
		responses[i] = *s.toResponse(&a)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.AppointmentListResponse{
		Data: responses,
		Meta: dto.AppointmentPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *appointmentService) TodayList(query *dto.AppointmentPaginationQuery) (*dto.AppointmentListResponse, error) {
	s.normalizeQuery(query, "appointment_date", "asc")

	appointments, total, err := s.repo.TodayList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AppointmentResponse, len(appointments))
	for i, a := range appointments {
		responses[i] = *s.toResponse(&a)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.AppointmentListResponse{
		Data: responses,
		Meta: dto.AppointmentPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *appointmentService) FindByID(id uint) (*dto.AppointmentResponse, error) {
	appointment, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(appointment), nil
}

func (s *appointmentService) FindByIDUnscoped(id uint) (*dto.AppointmentResponse, error) {
	appointment, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(appointment), nil
}

func (s *appointmentService) Create(req *dto.CreateAppointmentRequest) (*dto.AppointmentResponse, error) {
	appointment := &models.Appointment{
		PatientID:       req.PatientID,
		DoctorID:        req.DoctorID,
		AppointmentDate: req.AppointmentDate,
		AppointmentTime: req.AppointmentTime,
		DurationMinutes: req.DurationMinutes,
		Reason:          req.Reason,
		Notes:           req.Notes,
		Status:          "scheduled",
	}

	if err := s.repo.Create(appointment); err != nil {
		return nil, err
	}

	// Fetch again to get preloaded relations for response
	created, _ := s.repo.FindByID(appointment.ID)
	if created == nil {
		created = appointment
	}

	return s.toResponse(created), nil
}

func (s *appointmentService) Update(id uint, req *dto.UpdateAppointmentRequest) (*dto.AppointmentResponse, error) {
	appointment, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.AppointmentDate != nil {
		appointment.AppointmentDate = *req.AppointmentDate
	}
	if req.AppointmentTime != nil {
		appointment.AppointmentTime = *req.AppointmentTime
	}
	if req.DurationMinutes != nil {
		appointment.DurationMinutes = *req.DurationMinutes
	}
	if req.Reason != nil {
		appointment.Reason = *req.Reason
	}
	if req.Notes != nil {
		appointment.Notes = *req.Notes
	}

	if err := s.repo.Update(appointment); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(appointment.ID)
	if updated == nil {
		updated = appointment
	}

	return s.toResponse(updated), nil
}

func (s *appointmentService) Confirm(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Confirm(id)
}

func (s *appointmentService) Start(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Start(id)
}

func (s *appointmentService) Complete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Complete(id)
}

func (s *appointmentService) Cancel(id uint, req *dto.CancelAppointmentRequest) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Cancel(id, req.Reason)
}

func (s *appointmentService) Reschedule(id uint, req *dto.RescheduleAppointmentRequest) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Reschedule(id, req.AppointmentDate, req.AppointmentTime)
}

func (s *appointmentService) NoShow(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.NoShow(id)
}

func (s *appointmentService) SoftDelete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}

func (s *appointmentService) Restore(id uint) error {
	_, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	return s.repo.Restore(id)
}

func (s *appointmentService) HardDelete(id uint) error {
	_, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	return s.repo.HardDelete(id)
}
