package medicalcondition

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type MedicalConditionService interface {
	List(query *dto.MedicalConditionPaginationQuery) (*dto.MedicalConditionListResponse, error)
	FindByID(id uint) (*dto.MedicalConditionResponse, error)
	Create(req *dto.CreateMedicalConditionRequest) (*dto.MedicalConditionResponse, error)
	Update(id uint, req *dto.UpdateMedicalConditionRequest) (*dto.MedicalConditionResponse, error)
	Delete(id uint) error
}

type medicalConditionService struct {
	repo   repository.MedicalConditionRepository
	config *config.Config
}

func NewMedicalConditionService(repo repository.MedicalConditionRepository, config *config.Config) MedicalConditionService {
	return &medicalConditionService{
		repo:   repo,
		config: config,
	}
}

func (s *medicalConditionService) normalizeQuery(query *dto.MedicalConditionPaginationQuery, defaultSortBy, defaultSortDir string) {
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

func (s *medicalConditionService) toResponse(m *models.MedicalCondition) *dto.MedicalConditionResponse {
	if m == nil {
		return nil
	}

	resp := &dto.MedicalConditionResponse{
		ID:            m.ID,
		PatientID:     m.PatientID,
		ConditionName: m.ConditionName,
		ICDCode:       m.ICDCode,
		Status:        m.Status,
		Notes:         m.Notes,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}

	if m.DiagnosedDate != nil {
		resp.DiagnosedDate = m.DiagnosedDate.Format("2006-01-02")
	}

	if m.Patient.ID != 0 {
		resp.Patient = &dto.MedicalConditionPatientResponse{
			ID:          m.Patient.ID,
			PatientCode: m.Patient.PatientCode,
			FullName:    m.Patient.FullName,
		}
	}

	return resp
}

func (s *medicalConditionService) List(query *dto.MedicalConditionPaginationQuery) (*dto.MedicalConditionListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")

	conditions, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	var data []dto.MedicalConditionResponse
	for i := range conditions {
		data = append(data, *s.toResponse(&conditions[i]))
	}

	if data == nil {
		data = []dto.MedicalConditionResponse{}
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.MedicalConditionListResponse{
		Data: data,
		Meta: dto.MedicalConditionPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicalConditionService) FindByID(id uint) (*dto.MedicalConditionResponse, error) {
	condition, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(condition), nil
}

func (s *medicalConditionService) Create(req *dto.CreateMedicalConditionRequest) (*dto.MedicalConditionResponse, error) {
	status := req.Status
	if status == "" {
		status = "ongoing"
	}

	condition := &models.MedicalCondition{
		PatientID:     req.PatientID,
		ConditionName: req.ConditionName,
		ICDCode:       req.ICDCode,
		Status:        status,
		Notes:         req.Notes,
	}

	if req.DiagnosedDate != "" {
		t, err := time.Parse("2006-01-02", req.DiagnosedDate)
		if err == nil {
			condition.DiagnosedDate = &t
		}
	}

	if err := s.repo.Create(condition); err != nil {
		return nil, err
	}

	created, err := s.repo.FindByID(condition.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(created), nil
}

func (s *medicalConditionService) Update(id uint, req *dto.UpdateMedicalConditionRequest) (*dto.MedicalConditionResponse, error) {
	condition, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.ConditionName != "" {
		condition.ConditionName = req.ConditionName
	}
	if req.ICDCode != "" {
		condition.ICDCode = req.ICDCode
	}
	if req.Status != "" {
		condition.Status = req.Status
	}
	if req.Notes != "" {
		condition.Notes = req.Notes
	}
	if req.DiagnosedDate != "" {
		t, err := time.Parse("2006-01-02", req.DiagnosedDate)
		if err == nil {
			condition.DiagnosedDate = &t
		}
	}

	if err := s.repo.Update(condition); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updated), nil
}

func (s *medicalConditionService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("medical condition not found")
	}

	return s.repo.Delete(id)
}

func normalizeCachePart(s string) string {
	if s == "" {
		return "all"
	}
	return strings.ToLower(s)
}
