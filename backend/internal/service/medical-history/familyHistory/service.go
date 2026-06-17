package familyhistory

import (
	"errors"
	"math"
	"strings"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type FamilyHistoryService interface {
	List(query *dto.FamilyHistoryPaginationQuery) (*dto.FamilyHistoryListResponse, error)
	FindByID(id uint) (*dto.FamilyHistoryResponse, error)
	Create(req *dto.CreateFamilyHistoryRequest) (*dto.FamilyHistoryResponse, error)
	Update(id uint, req *dto.UpdateFamilyHistoryRequest) (*dto.FamilyHistoryResponse, error)
	Delete(id uint) error
}

type familyHistoryService struct {
	repo   repository.FamilyHistoryRepository
	config *config.Config
}

func NewFamilyHistoryService(repo repository.FamilyHistoryRepository, config *config.Config) FamilyHistoryService {
	return &familyHistoryService{
		repo:   repo,
		config: config,
	}
}

func (s *familyHistoryService) normalizeQuery(query *dto.FamilyHistoryPaginationQuery, defaultSortBy, defaultSortDir string) {
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

func (s *familyHistoryService) toResponse(fh *models.FamilyHistory) *dto.FamilyHistoryResponse {
	if fh == nil {
		return nil
	}

	resp := &dto.FamilyHistoryResponse{
		ID:            fh.ID,
		PatientID:     fh.PatientID,
		FamilyMember:  fh.FamilyMember,
		ConditionName: fh.ConditionName,
		Relation:      fh.Relation,
		Notes:         fh.Notes,
		CreatedAt:     fh.CreatedAt,
		UpdatedAt:     fh.UpdatedAt,
	}

	if fh.Patient.ID != 0 {
		resp.Patient = &dto.FamilyHistoryPatientResponse{
			ID:          fh.Patient.ID,
			PatientCode: fh.Patient.PatientCode,
			FullName:    fh.Patient.FullName,
		}
	}

	return resp
}

func (s *familyHistoryService) List(query *dto.FamilyHistoryPaginationQuery) (*dto.FamilyHistoryListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")

	histories, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	var data []dto.FamilyHistoryResponse
	for i := range histories {
		data = append(data, *s.toResponse(&histories[i]))
	}

	if data == nil {
		data = []dto.FamilyHistoryResponse{}
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.FamilyHistoryListResponse{
		Data: data,
		Meta: dto.FamilyHistoryPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *familyHistoryService) FindByID(id uint) (*dto.FamilyHistoryResponse, error) {
	history, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(history), nil
}

func (s *familyHistoryService) Create(req *dto.CreateFamilyHistoryRequest) (*dto.FamilyHistoryResponse, error) {
	history := &models.FamilyHistory{
		PatientID:     req.PatientID,
		FamilyMember:  req.FamilyMember,
		ConditionName: req.ConditionName,
		Relation:      req.Relation,
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

func (s *familyHistoryService) Update(id uint, req *dto.UpdateFamilyHistoryRequest) (*dto.FamilyHistoryResponse, error) {
	history, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.FamilyMember != "" {
		history.FamilyMember = req.FamilyMember
	}
	if req.ConditionName != "" {
		history.ConditionName = req.ConditionName
	}
	if req.Relation != "" {
		history.Relation = req.Relation
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

func (s *familyHistoryService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("family history not found")
	}

	return s.repo.Delete(id)
}

func normalizeCachePart(s string) string {
	if s == "" {
		return "all"
	}
	return strings.ToLower(s)
}
