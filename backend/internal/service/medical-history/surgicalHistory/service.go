package surgicalhistory

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

type SurgicalHistoryService interface {
	List(query *dto.SurgicalHistoryPaginationQuery) (*dto.SurgicalHistoryListResponse, error)
	FindByID(id uint) (*dto.SurgicalHistoryResponse, error)
	Create(req *dto.CreateSurgicalHistoryRequest) (*dto.SurgicalHistoryResponse, error)
	Update(id uint, req *dto.UpdateSurgicalHistoryRequest) (*dto.SurgicalHistoryResponse, error)
	Delete(id uint) error
}

type surgicalHistoryService struct {
	repo   repository.SurgicalHistoryRepository
	config *config.Config
}

func NewSurgicalHistoryService(repo repository.SurgicalHistoryRepository, config *config.Config) SurgicalHistoryService {
	return &surgicalHistoryService{
		repo:   repo,
		config: config,
	}
}

func (s *surgicalHistoryService) normalizeQuery(query *dto.SurgicalHistoryPaginationQuery, defaultSortBy, defaultSortDir string) {
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

func (s *surgicalHistoryService) toResponse(sh *models.SurgicalHistory) *dto.SurgicalHistoryResponse {
	if sh == nil {
		return nil
	}

	resp := &dto.SurgicalHistoryResponse{
		ID:            sh.ID,
		PatientID:     sh.PatientID,
		ProcedureName: sh.ProcedureName,
		SurgeryDate:   sh.SurgeryDate.Format("2006-01-02"),
		SurgeonName:   sh.SurgeonName,
		Hospital:      sh.Hospital,
		Complication:  sh.Complication,
		Notes:         sh.Notes,
		CreatedAt:     sh.CreatedAt,
		UpdatedAt:     sh.UpdatedAt,
	}

	if sh.Patient.ID != 0 {
		resp.Patient = &dto.SurgicalHistoryPatientResponse{
			ID:          sh.Patient.ID,
			PatientCode: sh.Patient.PatientCode,
			FullName:    sh.Patient.FullName,
		}
	}

	return resp
}

func (s *surgicalHistoryService) List(query *dto.SurgicalHistoryPaginationQuery) (*dto.SurgicalHistoryListResponse, error) {
	s.normalizeQuery(query, "surgery_date", "desc")

	histories, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	var data []dto.SurgicalHistoryResponse
	for i := range histories {
		data = append(data, *s.toResponse(&histories[i]))
	}

	if data == nil {
		data = []dto.SurgicalHistoryResponse{}
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.SurgicalHistoryListResponse{
		Data: data,
		Meta: dto.SurgicalHistoryPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *surgicalHistoryService) FindByID(id uint) (*dto.SurgicalHistoryResponse, error) {
	history, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(history), nil
}

func (s *surgicalHistoryService) Create(req *dto.CreateSurgicalHistoryRequest) (*dto.SurgicalHistoryResponse, error) {
	surgeryDate, err := time.Parse("2006-01-02", req.SurgeryDate)
	if err != nil {
		return nil, errors.New("invalid surgery_date format, expected YYYY-MM-DD")
	}

	history := &models.SurgicalHistory{
		PatientID:     req.PatientID,
		ProcedureName: req.ProcedureName,
		SurgeryDate:   surgeryDate,
		SurgeonName:   req.SurgeonName,
		Hospital:      req.Hospital,
		Complication:  req.Complication,
		Notes:         req.Notes,
	}

	if err := s.repo.Create(history); err != nil {
		return nil, err
	}

	created, err := s.repo.FindByID(history.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(created), nil
}

func (s *surgicalHistoryService) Update(id uint, req *dto.UpdateSurgicalHistoryRequest) (*dto.SurgicalHistoryResponse, error) {
	history, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.ProcedureName != "" {
		history.ProcedureName = req.ProcedureName
	}
	if req.SurgeryDate != "" {
		t, err := time.Parse("2006-01-02", req.SurgeryDate)
		if err == nil {
			history.SurgeryDate = t
		}
	}
	if req.SurgeonName != "" {
		history.SurgeonName = req.SurgeonName
	}
	if req.Hospital != "" {
		history.Hospital = req.Hospital
	}
	if req.Complication != "" {
		history.Complication = req.Complication
	}
	if req.Notes != "" {
		history.Notes = req.Notes
	}

	if err := s.repo.Update(history); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updated), nil
}

func (s *surgicalHistoryService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("surgical history not found")
	}

	return s.repo.Delete(id)
}

func normalizeCachePart(s string) string {
	if s == "" {
		return "all"
	}
	return strings.ToLower(s)
}
