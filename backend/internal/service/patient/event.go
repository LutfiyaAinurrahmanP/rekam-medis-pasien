package patient

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

// eventPatientService membungkus PatientService dengan Kafka event publishing.
type eventPatientService struct {
	inner     PatientService
	publisher kafka.EventPublisher
}

// NewEventPatientService mengembalikan PatientService dengan event publishing.
// Jika publisher nil, inner service dikembalikan langsung.
func NewEventPatientService(inner PatientService, publisher kafka.EventPublisher) PatientService {
	if publisher == nil {
		return inner
	}
	return &eventPatientService{inner: inner, publisher: publisher}
}

// ─── Read operations ──────────────────────────────────────────────────────────

func (s *eventPatientService) ListPatients(query *dto.PatientPaginationQuery) (*dto.PatientListResponse, error) {
	return s.inner.ListPatients(query)
}

func (s *eventPatientService) DeleteListPatients(query *dto.PatientPaginationQuery) (*dto.PatientDeletedListResponse, error) {
	return s.inner.DeleteListPatients(query)
}

func (s *eventPatientService) GetPatientByID(id uint) (*dto.PatientResponse, error) {
	return s.inner.GetPatientByID(id)
}

func (s *eventPatientService) GetMyPatientData(userID uint) (*dto.PatientResponse, error) {
	return s.inner.GetMyPatientData(userID)
}

// ─── Write operations ─────────────────────────────────────────────────────────

func (s *eventPatientService) CreatePatient(req *dto.CreatePatientRequest) (*dto.PatientResponse, error) {
	result, err := s.inner.CreatePatient(req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicPatientCreated,
		events.NewPatientCreatedEvent(
			result.ID, result.PatientCode, result.FullName,
			result.DateOfBirth, result.Gender, result.BloodType,
			result.Phone, result.Email, result.InsuranceProvider,
		))
	return result, nil
}

func (s *eventPatientService) UpdatePatient(id uint, req *dto.UpdatePatientRequest) (*dto.PatientResponse, error) {
	result, err := s.inner.UpdatePatient(id, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicPatientUpdated,
		events.NewPatientUpdatedEvent(
			result.ID, result.PatientCode, result.FullName,
			result.DateOfBirth, result.Gender, result.BloodType,
			result.Phone, result.Email, result.InsuranceProvider, "admin_update",
		))
	return result, nil
}

func (s *eventPatientService) UpdateMyPatientData(userID uint, req *dto.UpdatePatientRequest) (*dto.PatientResponse, error) {
	result, err := s.inner.UpdateMyPatientData(userID, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(kafka.TopicPatientUpdated,
		events.NewPatientUpdatedEvent(
			result.ID, result.PatientCode, result.FullName,
			result.DateOfBirth, result.Gender, result.BloodType,
			result.Phone, result.Email, result.InsuranceProvider, "self_update",
		))
	return result, nil
}

func (s *eventPatientService) SoftDeletePatient(id uint) error {
	// Ambil data sebelum soft delete
	patient, _ := s.inner.GetPatientByID(id)
	if err := s.inner.SoftDeletePatient(id); err != nil {
		return err
	}
	code, name := "", ""
	if patient != nil {
		code = patient.PatientCode
		name = patient.FullName
	}
	s.publisher.PublishAsync(kafka.TopicPatientDeleted,
		events.NewPatientDeletedEvent(id, code, name, "soft_delete"))
	return nil
}

func (s *eventPatientService) HardDeletePatient(id uint) error {
	patient, _ := s.inner.GetPatientByID(id)
	if err := s.inner.HardDeletePatient(id); err != nil {
		return err
	}
	code, name := "", ""
	if patient != nil {
		code = patient.PatientCode
		name = patient.FullName
	}
	s.publisher.PublishAsync(kafka.TopicPatientDeleted,
		events.NewPatientDeletedEvent(id, code, name, "hard_delete"))
	return nil
}

func (s *eventPatientService) RestorePatient(id uint) error {
	if err := s.inner.RestorePatient(id); err != nil {
		return err
	}
	if patient, err := s.inner.GetPatientByID(id); err == nil {
		s.publisher.PublishAsync(kafka.TopicPatientRestored,
			events.NewPatientRestoredEvent(patient.ID, patient.PatientCode, patient.FullName))
	}
	return nil
}
