package patient

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type PatientService interface {
	ListPatients(query *dto.PatientPaginationQuery) (*dto.PatientListResponse, error)
	DeleteListPatients(query *dto.PatientPaginationQuery) (*dto.PatientDeletedListResponse, error)
	GetPatientByID(id uint) (*dto.PatientResponse, error)
	GetPatientByCode(code string) (*dto.PatientResponse, error)
	GetMyPatientData(userID uint) (*dto.PatientResponse, error)
	UpdateMyPatientData(userID uint, req *dto.UpdatePatientRequest) (*dto.PatientResponse, error)
	CreatePatient(req *dto.CreatePatientRequest) (*dto.PatientResponse, error)
	UpdatePatient(id uint, req *dto.UpdatePatientRequest) (*dto.PatientResponse, error)
	SoftDeletePatient(id uint) error
	RestorePatient(id uint) error
	HardDeletePatient(id uint) error
}

type patientService struct {
	repo   repository.PatientRepository
	config *config.Config
}

func NewPatientService(repo repository.PatientRepository, config *config.Config) PatientService {
	return &patientService{
		repo:   repo,
		config: config,
	}
}

func (s *patientService) ListPatients(query *dto.PatientPaginationQuery) (*dto.PatientListResponse, error) {
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
		query.SortBy = "created_at"
	}

	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	patients, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	patientResponses := make([]dto.PatientResponse, len(patients))
	for i, patient := range patients {
		patientResponses[i] = *s.toPatientResponse(&patient)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.PatientListResponse{
		Data: patientResponses,
		Meta: dto.PatientPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *patientService) DeleteListPatients(query *dto.PatientPaginationQuery) (*dto.PatientDeletedListResponse, error) {
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
		query.SortBy = "deleted_at"
	}

	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	deletedPatients, total, err := s.repo.DeleteList(query)
	if err != nil {
		return nil, err
	}

	deletedPatientResponses := make([]dto.DeletedPatientResponse, len(deletedPatients))
	for i, deletedPatient := range deletedPatients {
		deletedPatientResponses[i] = *s.toDeletePatientResponse(&deletedPatient)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))

	return &dto.PatientDeletedListResponse{
		Data: deletedPatientResponses,
		Meta: dto.PatientPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *patientService) GetPatientByID(id uint) (*dto.PatientResponse, error) {
	patient, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return s.toPatientResponse(patient), nil
}

func (s *patientService) GetPatientByCode(code string) (*dto.PatientResponse, error) {
	patient, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}
	return s.toPatientResponse(patient), nil
}

func (s *patientService) GetMyPatientData(userID uint) (*dto.PatientResponse, error) {
	patient, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.toPatientResponse(patient), nil
}

func (s *patientService) UpdateMyPatientData(userID uint, req *dto.UpdatePatientRequest) (*dto.PatientResponse, error) {
	patient, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Patient can only update certain fields
	if req.FullName != nil {
		patient.FullName = *req.FullName
	}

	if req.Phone != nil {
		patient.Phone = *req.Phone
	}

	if req.Email != nil {
		patient.Email = *req.Email
	}

	if req.Address != nil {
		patient.Address = *req.Address
	}

	if req.EmergencyContactName != nil {
		patient.EmergencyContactName = *req.EmergencyContactName
	}

	if req.EmergencyContactPhone != nil {
		patient.EmergencyContactPhone = *req.EmergencyContactPhone
	}

	if req.InsuranceNumber != nil {
		patient.InsuranceNumber = *req.InsuranceNumber
	}

	if req.InsuranceProvider != nil {
		patient.InsuranceProvider = *req.InsuranceProvider
	}

	if req.Allergies != nil {
		patient.Allergies = *req.Allergies
	}

	// Note: Patient cannot update patient_code, date_of_birth, gender, blood_type

	if err := s.repo.Update(patient); err != nil {
		return nil, err
	}

	return s.toPatientResponse(patient), nil
}

func (s *patientService) CreatePatient(req *dto.CreatePatientRequest) (*dto.PatientResponse, error) {
	exists, err := s.repo.IsCodeExists(req.PatientCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("patient code already exists")
	}

	patient := &models.Patient{
		UserID:                req.UserID,
		PatientCode:           req.PatientCode,
		FullName:              req.FullName,
		DateOfBirth:           req.DateOfBirth,
		Gender:                req.Gender,
		BloodType:             req.BloodType,
		Phone:                 req.Phone,
		Email:                 req.Email,
		Address:               req.Address,
		EmergencyContactName:  req.EmergencyContactName,
		EmergencyContactPhone: req.EmergencyContactPhone,
		InsuranceNumber:       req.InsuranceNumber,
		InsuranceProvider:     req.InsuranceProvider,
		Allergies:             req.Allergies,
	}

	if err := s.repo.Create(patient); err != nil {
		return nil, err
	}
	return s.toPatientResponse(patient), nil
}

func (s *patientService) UpdatePatient(id uint, req *dto.UpdatePatientRequest) (*dto.PatientResponse, error) {
	patient, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	if req.FullName != nil {
		patient.FullName = *req.FullName
	}

	if req.DateOfBirth != nil {
		patient.DateOfBirth = *req.DateOfBirth
	}

	if req.Gender != nil {
		patient.Gender = *req.Gender
	}

	if req.BloodType != nil {
		patient.BloodType = *req.BloodType
	}

	if req.Phone != nil {
		patient.Phone = *req.Phone
	}

	if req.Email != nil {
		patient.Email = *req.Email
	}

	if req.Address != nil {
		patient.Address = *req.Address
	}

	if req.EmergencyContactName != nil {
		patient.EmergencyContactName = *req.EmergencyContactName
	}

	if req.EmergencyContactPhone != nil {
		patient.EmergencyContactPhone = *req.EmergencyContactPhone
	}

	if req.InsuranceNumber != nil {
		patient.InsuranceNumber = *req.InsuranceNumber
	}

	if req.InsuranceProvider != nil {
		patient.InsuranceProvider = *req.InsuranceProvider
	}

	if req.Allergies != nil {
		patient.Allergies = *req.Allergies
	}

	if err := s.repo.Update(patient); err != nil {
		return nil, err
	}

	return s.toPatientResponse(patient), nil
}

func (s *patientService) SoftDeletePatient(id uint) error {
	_, err := s.repo.FindById(id)
	if err != nil {
		return err
	}

	return s.repo.SoftDelete(id)
}

func (s *patientService) RestorePatient(id uint) error {
	return s.repo.Restore(id)
}

func (s *patientService) HardDeletePatient(id uint) error {
	return s.repo.HardDelete(id)
}

func (s *patientService) toPatientResponse(patient *models.Patient) *dto.PatientResponse {
	return &dto.PatientResponse{
		ID:                    patient.ID,
		UserID:                patient.UserID,
		PatientCode:           patient.PatientCode,
		FullName:              patient.FullName,
		DateOfBirth:           patient.DateOfBirth,
		Age:                   patient.Age(),
		Gender:                patient.Gender,
		BloodType:             patient.BloodType,
		Phone:                 patient.Phone,
		Email:                 patient.Email,
		Address:               patient.Address,
		EmergencyContactName:  patient.EmergencyContactName,
		EmergencyContactPhone: patient.EmergencyContactPhone,
		InsuranceNumber:       patient.InsuranceNumber,
		InsuranceProvider:     patient.InsuranceProvider,
		Allergies:             patient.Allergies,
		CreatedAt:             patient.CreatedAt,
		UpdatedAt:             patient.UpdatedAt,
	}
}

func (s *patientService) toDeletePatientResponse(patient *models.Patient) *dto.DeletedPatientResponse {
	return &dto.DeletedPatientResponse{
		ID:                    patient.ID,
		UserID:                patient.UserID,
		PatientCode:           patient.PatientCode,
		FullName:              patient.FullName,
		DateOfBirth:           patient.DateOfBirth,
		Age:                   patient.Age(),
		Gender:                patient.Gender,
		BloodType:             patient.BloodType,
		Phone:                 patient.Phone,
		Email:                 patient.Email,
		Address:               patient.Address,
		EmergencyContactName:  patient.EmergencyContactName,
		EmergencyContactPhone: patient.EmergencyContactPhone,
		InsuranceNumber:       patient.InsuranceNumber,
		InsuranceProvider:     patient.InsuranceProvider,
		Allergies:             patient.Allergies,
		CreatedAt:             patient.CreatedAt,
		UpdatedAt:             patient.UpdatedAt,
		DeletedAt:             &patient.DeletedAt.Time,
	}
}
