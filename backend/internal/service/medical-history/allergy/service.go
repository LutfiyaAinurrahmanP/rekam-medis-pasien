package allergy

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type AllergyService interface {
	List(query *dto.AllergyPaginationQuery) (*dto.AllergyListResponse, error)
	FindByID(id uint) (*dto.AllergyResponse, error)
	Create(req *dto.CreateAllergyRequest) (*dto.AllergyResponse, error)
	Update(id uint, req *dto.UpdateAllergyRequest) (*dto.AllergyResponse, error)
	Delete(id uint) error
}

type allergyService struct {
	repo   repository.AllergyRepository
	config *config.Config
}

func NewAllergyService(repo repository.AllergyRepository, config *config.Config) AllergyService {
	return &allergyService{
		repo:   repo,
		config: config,
	}
}

func (s *allergyService) normalizeQuery(query *dto.AllergyPaginationQuery, defaultSortBy, defaultSortDir string) {
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

func (s *allergyService) toResponse(a *models.Allergy) *dto.AllergyResponse {
	if a == nil {
		return nil
	}

	resp := &dto.AllergyResponse{
		ID:           a.ID,
		PatientID:    a.PatientID,
		AllergenType: a.AllergenType,
		AllergenName: a.AllergenName,
		Reaction:     a.Reaction,
		Severity:     a.Severity,
		Notes:        a.Notes,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}

	if a.Patient.ID != 0 {
		resp.Patient = &dto.AllergyPatientResponse{
			ID:          a.Patient.ID,
			PatientCode: a.Patient.PatientCode,
			FullName:    a.Patient.FullName,
		}
	}

	return resp
}

func (s *allergyService) List(query *dto.AllergyPaginationQuery) (*dto.AllergyListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")

	allergies, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	var data []dto.AllergyResponse
	for i := range allergies {
		data = append(data, *s.toResponse(&allergies[i]))
	}

	if data == nil {
		data = []dto.AllergyResponse{}
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.AllergyListResponse{
		Data: data,
		Meta: dto.AllergyPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *allergyService) FindByID(id uint) (*dto.AllergyResponse, error) {
	allergy, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(allergy), nil
}

func (s *allergyService) Create(req *dto.CreateAllergyRequest) (*dto.AllergyResponse, error) {
	allergy := &models.Allergy{
		PatientID:    req.PatientID,
		AllergenType: req.AllergenType,
		AllergenName: req.AllergenName,
		Reaction:     req.Reaction,
		Severity:     req.Severity,
		Notes:        req.Notes,
	}

	if err := s.repo.Create(allergy); err != nil {
		return nil, err
	}

	created, err := s.repo.FindByID(allergy.ID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(created), nil
}

func (s *allergyService) Update(id uint, req *dto.UpdateAllergyRequest) (*dto.AllergyResponse, error) {
	allergy, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.AllergenType != "" {
		allergy.AllergenType = req.AllergenType
	}
	if req.AllergenName != "" {
		allergy.AllergenName = req.AllergenName
	}
	if req.Reaction != "" {
		allergy.Reaction = req.Reaction
	}
	if req.Severity != "" {
		allergy.Severity = req.Severity
	}
	if req.Notes != "" {
		allergy.Notes = req.Notes
	}

	if err := s.repo.Update(allergy); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(updated), nil
}

func (s *allergyService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("allergy not found")
	}

	return s.repo.Delete(id)
}


