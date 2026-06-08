package hospitalization

import (
	"errors"
	"math"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type HospitalizationService interface {
	List(query *dto.HospitalizationPaginationQuery) (*dto.HospitalizationListResponse, error)
	DeletedList(query *dto.HospitalizationPaginationQuery) (*dto.HospitalizationDeletedListResponse, error)
	FindByID(id uint) (*dto.HospitalizationResponse, error)
	FindByIDUnscoped(id uint) (*dto.HospitalizationResponse, error)
	Create(req *dto.CreateHospitalizationRequest) (*dto.HospitalizationResponse, error)
	Update(id uint, req *dto.UpdateHospitalizationRequest) (*dto.HospitalizationResponse, error)
	Discharge(id uint, req *dto.DischargeHospitalizationRequest) (*dto.HospitalizationResponse, error)
	Transfer(id uint, req *dto.TransferHospitalizationRequest) (*dto.HospitalizationResponse, error)
	Activate(id uint) (*dto.HospitalizationResponse, error)
	Deactivate(id uint) (*dto.HospitalizationResponse, error)
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type hospitalizationService struct {
	repo   repository.HospitalizationRepository
	config *config.Config
}

func NewHospitalizationService(repo repository.HospitalizationRepository, config *config.Config) HospitalizationService {
	return &hospitalizationService{
		repo:   repo,
		config: config,
	}
}

func (s *hospitalizationService) normalizeQuery(query *dto.HospitalizationPaginationQuery, defaultSortBy, defaultSortDir string) {
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

func (s *hospitalizationService) toResponse(m *models.Hospitalization) *dto.HospitalizationResponse {
	resp := &dto.HospitalizationResponse{
		ID:                 m.ID,
		PatientID:          m.PatientID,
		DoctorID:           m.DoctorID,
		RoomID:             m.RoomID,
		AdmissionDate:      m.AdmissionDate,
		AdmissionTime:      m.AdmissionTime,
		DischargeDate:      m.DischargeDate,
		DischargeTime:      m.DischargeTime,
		AdmissionReason:    m.ReasonForAdmission,
		Status:             m.Status,
		Notes:              m.Notes,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}

	if m.Patient != nil {
		resp.Patient = &dto.HospitalizationPatientResponse{
			ID:                  m.Patient.ID,
			Name:                m.Patient.FullName,
			MedicalRecordNumber: m.Patient.PatientCode,
		}
	}

	if m.Doctor != nil {
		resp.Doctor = &dto.HospitalizationDoctorResponse{
			ID:             m.Doctor.ID,
			Name:           m.Doctor.FullName,
			Specialization: "",
		}
		if m.Doctor.Specialization.Name != "" {
			resp.Doctor.Specialization = m.Doctor.Specialization.Name
		}
	}

	if m.Room != nil {
		resp.Room = &dto.HospitalizationRoomResponse{
			ID:         m.Room.ID,
			RoomNumber: m.Room.RoomNumber,
			RoomType:   "",
		}
		if m.Room.RoomType != nil {
			resp.Room.RoomType = m.Room.RoomType.Name
		}
	}

	return resp
}

func (s *hospitalizationService) toDeletedResponse(m *models.Hospitalization) *dto.DeletedHospitalizationResponse {
	resp := s.toResponse(m)
	deletedResp := &dto.DeletedHospitalizationResponse{
		HospitalizationResponse: *resp,
	}

	if m.DeletedAt.Valid {
		deletedResp.DeletedAt = &m.DeletedAt.Time
	}

	return deletedResp
}

func (s *hospitalizationService) List(query *dto.HospitalizationPaginationQuery) (*dto.HospitalizationListResponse, error) {
	s.normalizeQuery(query, "admission_date", "desc")

	records, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.HospitalizationResponse, len(records))
	for i, r := range records {
		responses[i] = *s.toResponse(&r)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.HospitalizationListResponse{
		Data: responses,
		Meta: dto.HospitalizationPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *hospitalizationService) DeletedList(query *dto.HospitalizationPaginationQuery) (*dto.HospitalizationDeletedListResponse, error) {
	s.normalizeQuery(query, "admission_date", "desc")

	records, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DeletedHospitalizationResponse, len(records))
	for i, r := range records {
		responses[i] = *s.toDeletedResponse(&r)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.HospitalizationDeletedListResponse{
		Data: responses,
		Meta: dto.HospitalizationPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *hospitalizationService) FindByID(id uint) (*dto.HospitalizationResponse, error) {
	record, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(record), nil
}

func (s *hospitalizationService) FindByIDUnscoped(id uint) (*dto.HospitalizationResponse, error) {
	record, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(record), nil
}

func (s *hospitalizationService) Create(req *dto.CreateHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	parsedDate, err := time.Parse("2006-01-02T15:04:05Z", req.AdmissionDate)
	if err != nil {
		return nil, errors.New("invalid admission_date format, must be YYYY-MM-DDTHH:MM:SSZ")
	}

	record := &models.Hospitalization{
		PatientID:          req.PatientID,
		RoomID:             req.RoomID,
		DoctorID:           req.AttendingDoctorID,
		AdmissionDate:      parsedDate.Format("2006-01-02"),
		AdmissionTime:      parsedDate.Format("15:04:05"),
		ReasonForAdmission: req.AdmissionReason,
		Status:             "admitted",
	}

	if req.Status != "" {
		record.Status = req.Status
	}

	if err := s.repo.Create(record); err != nil {
		return nil, err
	}

	created, _ := s.repo.FindByID(record.ID)
	if created == nil {
		created = record
	}

	return s.toResponse(created), nil
}

func (s *hospitalizationService) Update(id uint, req *dto.UpdateHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	record, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if record.Status == "discharged" {
		return nil, errors.New("cannot update discharged hospitalization")
	}

	updates := make(map[string]interface{})
	if req.RoomID != nil {
		updates["room_id"] = *req.RoomID
	}
	if req.AttendingDoctorID != nil {
		updates["doctor_id"] = *req.AttendingDoctorID
	}
	if req.AdmissionReason != nil {
		updates["reason_for_admission"] = *req.AdmissionReason
	}

	if len(updates) > 0 {
		if err := s.repo.Update(id, updates); err != nil {
			return nil, err
		}
	}

	updated, _ := s.repo.FindByID(record.ID)
	if updated == nil {
		updated = record
	}

	return s.toResponse(updated), nil
}

func (s *hospitalizationService) Discharge(id uint, req *dto.DischargeHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	record, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if record.Status == "discharged" {
		return nil, errors.New("patient has already been discharged")
	}

	now := time.Now()
	dateStr := now.Format("2006-01-02")
	timeStr := now.Format("15:04:05")

	updates := map[string]interface{}{
		"status":         "discharged",
		"discharge_date": dateStr,
		"discharge_time": timeStr,
	}

	if req.DischargeSummary != "" {
		if record.Notes != "" {
			updates["notes"] = record.Notes + "\n\nDischarge Summary: " + req.DischargeSummary
		} else {
			updates["notes"] = "Discharge Summary: " + req.DischargeSummary
		}
	}

	if err := s.repo.Discharge(id, updates); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(record.ID)
	if updated == nil {
		updated = record
	}

	return s.toResponse(updated), nil
}

func (s *hospitalizationService) Transfer(id uint, req *dto.TransferHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	record, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if record.Status == "discharged" {
		return nil, errors.New("cannot transfer discharged patient")
	}

	updates := map[string]interface{}{
		"status": "transferred",
	}

	if req.Notes != "" {
		if record.Notes != "" {
			updates["notes"] = record.Notes + "\n\nTransfer Notes: " + req.Notes
		} else {
			updates["notes"] = "Transfer Notes: " + req.Notes
		}
	}

	if err := s.repo.Transfer(id, updates); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(record.ID)
	if updated == nil {
		updated = record
	}

	return s.toResponse(updated), nil
}

func (s *hospitalizationService) Activate(id uint) (*dto.HospitalizationResponse, error) {
	record, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"status":         "admitted",
		"discharge_date": nil,
		"discharge_time": nil,
	}

	if err := s.repo.Activate(id, updates); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(record.ID)
	if updated == nil {
		updated = record
	}

	return s.toResponse(updated), nil
}

func (s *hospitalizationService) Deactivate(id uint) (*dto.HospitalizationResponse, error) {
	record, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dateStr := now.Format("2006-01-02")
	timeStr := now.Format("15:04:05")

	updates := map[string]interface{}{
		"status": "discharged",
	}
	if record.DischargeDate == nil {
		updates["discharge_date"] = dateStr
		updates["discharge_time"] = timeStr
	}

	if err := s.repo.Deactivate(id, updates); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(record.ID)
	if updated == nil {
		updated = record
	}

	return s.toResponse(updated), nil
}

func (s *hospitalizationService) SoftDelete(id uint) error {
	h, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if h.DeletedAt.Valid {
		return errors.New("hospitalization already deleted")
	}
	if h.Status == "admitted" {
		return errors.New("cannot delete admitted hospitalization")
	}
	if h.Status == "transferred" {
		return errors.New("cannot delete transferred hospitalization")
	}
	return s.repo.SoftDelete(id)
}

func (s *hospitalizationService) Restore(id uint) error {
	h, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	if !h.DeletedAt.Valid {
		return errors.New("hospitalization not deleted")
	}
	return s.repo.Restore(id)
}

func (s *hospitalizationService) HardDelete(id uint) error {
	h, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	if !h.DeletedAt.Valid {
		return errors.New("hospitalization not deleted")
	}
	return s.repo.HardDelete(id)
}