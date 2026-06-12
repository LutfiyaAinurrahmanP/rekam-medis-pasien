package medicalhistory

import (
	"math"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type MedicalHistoryService interface {
	List(query *dto.MedicalHistoryPaginationQuery) (*dto.MedicalHistoryListResponse, error)
	FindByID(id uint) (*dto.MedicalHistoryDetailResponse, error)
}

type medicalHistoryService struct {
	repo   repository.MedicalHistoryRepository
	config *config.Config
}

func NewMedicalHistoryService(repo repository.MedicalHistoryRepository, cfg *config.Config) MedicalHistoryService {
	return &medicalHistoryService{
		repo:   repo,
		config: cfg,
	}
}

func (s *medicalHistoryService) normalizeQuery(query *dto.MedicalHistoryPaginationQuery, defaultSortBy, defaultSortDir string) {
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

func (s *medicalHistoryService) List(query *dto.MedicalHistoryPaginationQuery) (*dto.MedicalHistoryListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")

	patients, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicalHistoryOverviewResponse, len(patients))
	for i, p := range patients {
		responses[i] = dto.MedicalHistoryOverviewResponse{
			ID:                 p.ID, // Treating patient ID as the medical history ID
			PatientID:          p.ID,
			PatientName:        p.FullName,
			AllergiesCount:     len(p.MedicalAllergies),
			ConditionsCount:    len(p.MedicalConditions),
			SurgeriesCount:     len(p.SurgicalHistories),
			FamilyHistoryCount: len(p.FamilyHistories),
			LastUpdated:        p.UpdatedAt,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicalHistoryListResponse{
		Data: responses,
		Meta: dto.MedicalHistoryPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicalHistoryService) FindByID(id uint) (*dto.MedicalHistoryDetailResponse, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toDetailResponse(p), nil
}

func (s *medicalHistoryService) toDetailResponse(p *models.Patient) *dto.MedicalHistoryDetailResponse {
	res := &dto.MedicalHistoryDetailResponse{
		ID:                p.ID,
		PatientID:         p.ID,
		Allergies:         make([]dto.AllergyResponse, 0),
		MedicalConditions: make([]dto.MedicalConditionResponse, 0),
		SurgicalHistories: make([]dto.SurgicalHistoryResponse, 0),
		FamilyHistories:   make([]dto.FamilyHistoryResponse, 0),
	}

	for _, a := range p.MedicalAllergies {
		res.Allergies = append(res.Allergies, dto.AllergyResponse{
			ID:           a.ID,
			PatientID:    a.PatientID,
			AllergenType: a.AllergenType,
			AllergenName: a.AllergenName,
			Reaction:     a.Reaction,
			Severity:     a.Severity,
			Notes:        a.Notes,
			CreatedAt:    a.CreatedAt,
			UpdatedAt:    a.UpdatedAt,
		})
	}
	for _, c := range p.MedicalConditions {
		res.MedicalConditions = append(res.MedicalConditions, dto.MedicalConditionResponse{
			ID:            c.ID,
			PatientID:     c.PatientID,
			ConditionName: c.ConditionName,
			ICDCode:       c.ICDCode,
			DiagnosedDate: diagnosedDateStr(c.DiagnosedDate),
			Status:        c.Status,
			Notes:         c.Notes,
			CreatedAt:     c.CreatedAt,
			UpdatedAt:     c.UpdatedAt,
		})
	}
	for _, sh := range p.SurgicalHistories {
		res.SurgicalHistories = append(res.SurgicalHistories, dto.SurgicalHistoryResponse{
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
		})
	}
	for _, f := range p.FamilyHistories {
		res.FamilyHistories = append(res.FamilyHistories, dto.FamilyHistoryResponse{
			ID:            f.ID,
			PatientID:     f.PatientID,
			FamilyMember:  f.FamilyMember,
			ConditionName: f.ConditionName,
			Relation:      f.Relation,
			Notes:         f.Notes,
			CreatedAt:     f.CreatedAt,
			UpdatedAt:     f.UpdatedAt,
		})
	}

	return res
}

func diagnosedDateStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}
