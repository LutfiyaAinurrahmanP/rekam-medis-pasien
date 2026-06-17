package labtest

import (
	"errors"
	"math"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type LabTestService interface {
	List(query *dto.LabTestPaginationQuery) (*dto.LabTestListResponse, error)
	DeletedList(query *dto.LabTestPaginationQuery) (*dto.LabTestDeletedListResponse, error)
	FindByID(id uint) (*dto.LabTestResponse, error)
	FindByIDUnscoped(id uint) (*dto.LabTestResponse, error)
	Create(req *dto.CreateLabTestRequest) (*dto.LabTestResponse, error)
	Update(id uint, req *dto.UpdateLabTestRequest) (*dto.LabTestResponse, error)
	CollectSample(id uint) (*dto.LabTestResponse, error)
	Start(id uint) (*dto.LabTestResponse, error)
	Complete(id uint, req *dto.CompleteLabTestRequest) (*dto.LabTestResponse, error)
	Cancel(id uint) (*dto.LabTestResponse, error)
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type labTestService struct {
	repo   repository.LabTestRepository
	config *config.Config
}

func NewLabTestService(repo repository.LabTestRepository, config *config.Config) LabTestService {
	return &labTestService{
		repo:   repo,
		config: config,
	}
}

func (s *labTestService) normalizeQuery(query *dto.LabTestPaginationQuery, defaultSortBy, defaultSortDir string) {
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

func (s *labTestService) toResponse(m *models.LabTest) *dto.LabTestResponse {
	if m == nil {
		return nil
	}

	var mrResp *dto.LabTestMedicalRecordResponse
	if m.MedicalRecord != nil {
		mrResp = &dto.LabTestMedicalRecordResponse{
			ID:             m.MedicalRecord.ID,
			VisitDate:      m.MedicalRecord.VisitDate,
			ChiefComplaint: m.MedicalRecord.ChiefComplaint,
		}
	}

	var ttResp *dto.LabTestTypeResponse
	if m.TestType != nil {
		ttResp = &dto.LabTestTypeResponse{
			ID:   m.TestType.ID,
			Name: m.TestType.Name,
			Code: m.TestType.Code,
		}
	}

	var drResp *dto.LabTestDoctorResponse
	if m.OrderedByDoctor != nil {
		specName := "Unspecified"
		if m.OrderedByDoctor.Specialization.Name != "" {
			specName = m.OrderedByDoctor.Specialization.Name
		}
		drResp = &dto.LabTestDoctorResponse{
			ID:             m.OrderedByDoctor.ID,
			Name:           m.OrderedByDoctor.FullName,
			Specialization: specName,
		}
	}

	return &dto.LabTestResponse{
		ID:                   m.ID,
		MedicalRecordID:      m.MedicalRecordID,
		MedicalRecord:        mrResp,
		TestTypeID:           m.TestTypeID,
		TestType:             ttResp,
		OrderedByDoctorID:    m.OrderedByDoctorID,
		OrderedByDoctor:      drResp,
		OrderDate:            m.OrderDate,
		SampleCollectionDate: m.SampleCollectionDate,
		TestStartDate:        m.TestStartDate,
		ResultDate:           m.ResultDate,
		ResultValue:          m.ResultValue,
		ResultUnit:           m.ResultUnit,
		ReferenceRange:       m.ReferenceRange,
		Status:               m.Status,
		Notes:                m.Notes,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}

func (s *labTestService) toDeletedResponse(m *models.LabTest) *dto.DeletedLabTestResponse {
	resp := s.toResponse(m)
	deletedResp := &dto.DeletedLabTestResponse{
		LabTestResponse: *resp,
	}

	if m.DeletedAt.Valid {
		deletedResp.DeletedAt = &m.DeletedAt.Time
	}

	return deletedResp
}

func (s *labTestService) List(query *dto.LabTestPaginationQuery) (*dto.LabTestListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")

	labTests, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.LabTestResponse, len(labTests))
	for i, r := range labTests {
		responses[i] = *s.toResponse(&r)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.LabTestListResponse{
		Data: responses,
		Meta: dto.LabTestPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *labTestService) DeletedList(query *dto.LabTestPaginationQuery) (*dto.LabTestDeletedListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")

	labTests, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DeletedLabTestResponse, len(labTests))
	for i, r := range labTests {
		responses[i] = *s.toDeletedResponse(&r)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.LabTestDeletedListResponse{
		Data: responses,
		Meta: dto.LabTestPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: int(total),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *labTestService) FindByID(id uint) (*dto.LabTestResponse, error) {
	labTest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(labTest), nil
}

func (s *labTestService) FindByIDUnscoped(id uint) (*dto.LabTestResponse, error) {
	labTest, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(labTest), nil
}

func (s *labTestService) Create(req *dto.CreateLabTestRequest) (*dto.LabTestResponse, error) {
	status := "ordered"
	if req.Status != "" {
		status = req.Status
	}

	labTest := &models.LabTest{
		MedicalRecordID:   req.MedicalRecordID,
		TestTypeID:        req.TestTypeID,
		OrderedByDoctorID: req.OrderedByDoctorID,
		OrderDate:         req.OrderDate,
		Status:            status,
		Notes:             req.Notes,
	}

	if err := s.repo.Create(labTest); err != nil {
		return nil, err
	}

	created, _ := s.repo.FindByID(labTest.ID)
	if created == nil {
		created = labTest
	}

	return s.toResponse(created), nil
}

func (s *labTestService) Update(id uint, req *dto.UpdateLabTestRequest) (*dto.LabTestResponse, error) {
	labTest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Notes != nil {
		labTest.Notes = *req.Notes
	}
	if req.ReferenceRange != nil {
		labTest.ReferenceRange = req.ReferenceRange
	}

	if err := s.repo.Update(labTest); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(id)
	if updated == nil {
		updated = labTest
	}

	return s.toResponse(updated), nil
}

func (s *labTestService) CollectSample(id uint) (*dto.LabTestResponse, error) {
	labTest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if labTest.Status == "cancelled" {
		return nil, errors.New("cannot collect sample for a cancelled lab test")
	}

	now := time.Now().Format("2006-01-02")
	labTest.Status = "sample_collected"
	labTest.SampleCollectionDate = &now

	if err := s.repo.Update(labTest); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(id)
	return s.toResponse(updated), nil
}

func (s *labTestService) Start(id uint) (*dto.LabTestResponse, error) {
	labTest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if labTest.Status == "cancelled" {
		return nil, errors.New("cannot start a cancelled lab test")
	}

	now := time.Now().Format("2006-01-02")
	labTest.Status = "in_progress"
	labTest.TestStartDate = &now

	if err := s.repo.Update(labTest); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(id)
	return s.toResponse(updated), nil
}

func (s *labTestService) Complete(id uint, req *dto.CompleteLabTestRequest) (*dto.LabTestResponse, error) {
	labTest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if labTest.Status == "cancelled" {
		return nil, errors.New("cannot complete a cancelled lab test")
	}

	now := time.Now().Format("2006-01-02")
	labTest.Status = "completed"
	labTest.ResultDate = &now

	if req.ResultValue != nil {
		labTest.ResultValue = req.ResultValue
	}
	if req.ResultUnit != nil {
		labTest.ResultUnit = req.ResultUnit
	}
	if req.ReferenceRange != nil {
		labTest.ReferenceRange = req.ReferenceRange
	}
	if req.Notes != nil {
		labTest.Notes = *req.Notes
	}

	if err := s.repo.Update(labTest); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(id)
	return s.toResponse(updated), nil
}

func (s *labTestService) Cancel(id uint) (*dto.LabTestResponse, error) {
	labTest, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if labTest.Status == "completed" {
		return nil, errors.New("cannot cancel a completed lab test")
	}

	labTest.Status = "cancelled"

	if err := s.repo.Update(labTest); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(id)
	return s.toResponse(updated), nil
}

func (s *labTestService) SoftDelete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}

func (s *labTestService) Restore(id uint) error {
	_, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	return s.repo.Restore(id)
}

func (s *labTestService) HardDelete(id uint) error {
	_, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	return s.repo.HardDelete(id)
}
