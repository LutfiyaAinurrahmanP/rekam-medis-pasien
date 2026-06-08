package medicalrecord

import (
	"errors"
	"fmt"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type MedicalRecordService interface {
	List(query *dto.MedicalRecordPaginationQuery) (*dto.MedicalRecordListResponse, error)
	DeletedList(query *dto.MedicalRecordPaginationQuery) (*dto.MedicalRecordDeletedListResponse, error)
	FindByID(id uint) (*dto.MedicalRecordResponse, error)
	FindByIDUnscoped(id uint) (*dto.MedicalRecordResponse, error)
	Create(req *dto.CreateMedicalRecordRequest) (*dto.MedicalRecordResponse, error)
	Update(id uint, req *dto.UpdateMedicalRecordRequest) (*dto.MedicalRecordResponse, error)
	Finalize(id uint) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type medicalRecordService struct {
	repo   repository.MedicalRecordRepository
	config *config.Config
}

func NewMedicalRecordService(repo repository.MedicalRecordRepository, config *config.Config) MedicalRecordService {
	return &medicalRecordService{
		repo:   repo,
		config: config,
	}
}

func (s *medicalRecordService) normalizeQuery(query *dto.MedicalRecordPaginationQuery, defaultSortBy, defaultSortDir string) {
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

func (s *medicalRecordService) toResponse(m *models.MedicalRecord) *dto.MedicalRecordResponse {
	if m == nil {
		return nil
	}

	var patientResp *dto.MedicalRecordPatientResponse
	if m.Patient != nil {
		patientResp = &dto.MedicalRecordPatientResponse{
			ID:          m.Patient.ID,
			PatientCode: m.Patient.PatientCode,
			FullName:    m.Patient.FullName,
			Age:         m.Patient.Age(),
			Gender:      m.Patient.Gender,
			BloodType:   m.Patient.BloodType,
		}
	}

	var doctorResp *dto.MedicalRecordDoctorResponse
	if m.Doctor != nil {
		specName := m.Doctor.Specialization.Name
		deptName := m.Doctor.Department.Name

		doctorResp = &dto.MedicalRecordDoctorResponse{
			ID:             m.Doctor.ID,
			FullName:       m.Doctor.FullName,
			Specialization: specName,
			Department:     deptName,
		}
	}

	// Map Vital Signs
	vitalSigns := &dto.VitalSignsResponse{}
	if m.VitalSign != nil {
		vitalSigns.BloodPressure = ""
		if m.VitalSign.SystolicBP != nil && m.VitalSign.DiastolicBP != nil {
			vitalSigns.BloodPressure = fmt.Sprintf("%d/%d", *m.VitalSign.SystolicBP, *m.VitalSign.DiastolicBP)
		}
		if m.VitalSign.HeartRate != nil {
			vitalSigns.HeartRate = *m.VitalSign.HeartRate
		}
		if m.VitalSign.RespiratoryRate != nil {
			vitalSigns.RespiratoryRate = *m.VitalSign.RespiratoryRate
		}
		if m.VitalSign.BodyTemperature != nil {
			vitalSigns.Temperature = *m.VitalSign.BodyTemperature
		}
		if m.VitalSign.OxygenSaturation != nil {
			vitalSigns.OxygenSaturation = int(*m.VitalSign.OxygenSaturation)
		}
		if m.VitalSign.WeightKg != nil {
			vitalSigns.Weight = *m.VitalSign.WeightKg
		}
		if m.VitalSign.HeightCm != nil {
			vitalSigns.Height = float64(*m.VitalSign.HeightCm)
		}
		if m.VitalSign.BMI != nil {
			vitalSigns.BMI = *m.VitalSign.BMI
		}
	}

	// Map Medical History
	medicalHistory := &dto.MedicalHistorySummaryResponse{
		Allergies:         make([]any, 0),
		MedicalConditions: make([]any, 0),
		SurgicalHistory:   make([]any, 0),
		FamilyHistory:     make([]any, 0),
	}

	if m.Patient != nil {
		for _, a := range m.Patient.MedicalAllergies {
			medicalHistory.Allergies = append(medicalHistory.Allergies, a)
		}
		for _, c := range m.Patient.MedicalConditions {
			medicalHistory.MedicalConditions = append(medicalHistory.MedicalConditions, c)
		}
		for _, s := range m.Patient.SurgicalHistories {
			medicalHistory.SurgicalHistory = append(medicalHistory.SurgicalHistory, s)
		}
		for _, f := range m.Patient.FamilyHistories {
			medicalHistory.FamilyHistory = append(medicalHistory.FamilyHistory, f)
		}
	}

	return &dto.MedicalRecordResponse{
		ID:                  m.ID,
		PatientID:           m.PatientID,
		Patient:             patientResp,
		DoctorID:            m.DoctorID,
		Doctor:              doctorResp,
		AppointmentID:       m.AppointmentID,
		VisitDate:           m.VisitDate,
		ChiefComplaint:      m.ChiefComplaint,
		HistoryOfIllness:    m.HistoryOfIllness,
		PhysicalExamination: m.PhysicalExamination,
		VitalSigns:          vitalSigns,
		MedicalHistory:      medicalHistory,
		Diagnosis:           m.Diagnosis,
		TreatmentPlan:       m.TreatmentPlan,
		Notes:               m.Notes,
		Status:              m.Status,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func (s *medicalRecordService) List(query *dto.MedicalRecordPaginationQuery) (*dto.MedicalRecordListResponse, error) {
	s.normalizeQuery(query, "visit_date", "desc")

	records, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MedicalRecordResponse, len(records))
	for i, r := range records {
		responses[i] = *s.toResponse(&r)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicalRecordListResponse{
		Data: responses,
		Meta: dto.MedicalRecordPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicalRecordService) DeletedList(query *dto.MedicalRecordPaginationQuery) (*dto.MedicalRecordDeletedListResponse, error) {
	s.normalizeQuery(query, "visit_date", "desc")

	records, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DeletedMedicalRecordResponse, len(records))
	for i, r := range records {
		responses[i] = *s.toDeletedResponse(&r)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.MedicalRecordDeletedListResponse{
		Data: responses,
		Meta: dto.MedicalRecordPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *medicalRecordService) FindByID(id uint) (*dto.MedicalRecordResponse, error) {
	record, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(record), nil
}

func (s *medicalRecordService) FindByIDUnscoped(id uint) (*dto.MedicalRecordResponse, error) {
	record, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(record), nil
}

func (s *medicalRecordService) Create(req *dto.CreateMedicalRecordRequest) (*dto.MedicalRecordResponse, error) {
	record := &models.MedicalRecord{
		PatientID:           req.PatientID,
		DoctorID:            req.DoctorID,
		AppointmentID:       req.AppointmentID,
		VisitDate:           req.VisitDate,
		ChiefComplaint:      req.ChiefComplaint,
		HistoryOfIllness:    req.HistoryOfIllness,
		PhysicalExamination: req.PhysicalExamination,
		Diagnosis:           req.Diagnosis,
		TreatmentPlan:       req.TreatmentPlan,
		Notes:               req.Notes,
		Status:              "draft",
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

func (s *medicalRecordService) Update(id uint, req *dto.UpdateMedicalRecordRequest) (*dto.MedicalRecordResponse, error) {
	record, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if record.Status != "draft" {
		return nil, errors.New("medical record is not in draft status")
	}

	if req.ChiefComplaint != nil {
		record.ChiefComplaint = *req.ChiefComplaint
	}
	if req.HistoryOfIllness != nil {
		record.HistoryOfIllness = *req.HistoryOfIllness
	}
	if req.PhysicalExamination != nil {
		record.PhysicalExamination = *req.PhysicalExamination
	}
	if req.Diagnosis != nil {
		record.Diagnosis = *req.Diagnosis
	}
	if req.TreatmentPlan != nil {
		record.TreatmentPlan = *req.TreatmentPlan
	}
	if req.Notes != nil {
		record.Notes = *req.Notes
	}

	if err := s.repo.Update(record); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(record.ID)
	if updated == nil {
		updated = record
	}

	return s.toResponse(updated), nil
}

func (s *medicalRecordService) Finalize(id uint) error {
	mr, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if mr.Status != "draft" {
		return errors.New("medical record is not in draft status")
	}
	return s.repo.Finalize(id)
}

func (s *medicalRecordService) SoftDelete(id uint) error {
	mr, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if mr.Status != "draft" {
		return errors.New("medical record is not in draft status")
	}
	return s.repo.SoftDelete(id)
}

func (s *medicalRecordService) Restore(id uint) error {
	mr, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	if mr.Status != "draft" {
		return errors.New("medical record is not in draft status")
	}
	return s.repo.Restore(id)
}

func (s *medicalRecordService) HardDelete(id uint) error {
	mr, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	if mr.Status != "draft" {
		return errors.New("medical record is not in draft status")
	}
	return s.repo.HardDelete(id)
}

func (s *medicalRecordService) toDeletedResponse(m *models.MedicalRecord) *dto.DeletedMedicalRecordResponse {
	resp := s.toResponse(m)
	deletedResp := &dto.DeletedMedicalRecordResponse{
		MedicalRecordResponse: *resp,
	}

	if m.DeletedAt.Valid {
		deletedResp.DeletedAt = &m.DeletedAt.Time
	}

	return deletedResp
}
