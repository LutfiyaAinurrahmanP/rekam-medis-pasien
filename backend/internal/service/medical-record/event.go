package medicalrecord

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventedMedicalRecordService struct {
	inner     MedicalRecordService
	publisher kafka.EventPublisher
}

func NewEventedMedicalRecordService(inner MedicalRecordService, publisher kafka.EventPublisher) MedicalRecordService {
	if publisher == nil {
		return inner
	}
	return &eventedMedicalRecordService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventedMedicalRecordService) List(query *dto.MedicalRecordPaginationQuery) (*dto.MedicalRecordListResponse, error) {
	return s.inner.List(query)
}

func (s *eventedMedicalRecordService) DeletedList(query *dto.MedicalRecordPaginationQuery) (*dto.MedicalRecordDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}

func (s *eventedMedicalRecordService) FindByID(id uint) (*dto.MedicalRecordResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventedMedicalRecordService) FindByIDUnscoped(id uint) (*dto.MedicalRecordResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *eventedMedicalRecordService) Create(req *dto.CreateMedicalRecordRequest) (*dto.MedicalRecordResponse, error) {
	resp, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicalRecordCreated,
		events.NewMedicalRecordCreatedEvent(
			resp.ID,
			resp.PatientID,
			resp.DoctorID,
			resp.VisitDate,
			resp.ChiefComplaint,
			resp.Diagnosis,
			resp.Status,
		),
	)
	return resp, nil
}

func (s *eventedMedicalRecordService) Update(id uint, req *dto.UpdateMedicalRecordRequest) (*dto.MedicalRecordResponse, error) {
	resp, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicalRecordUpdated,
		events.NewMedicalRecordUpdatedEvent(
			resp.ID,
			resp.PatientID,
			resp.DoctorID,
			resp.VisitDate,
			resp.ChiefComplaint,
			resp.Diagnosis,
			resp.Status,
			"update",
		),
	)
	return resp, nil
}

func (s *eventedMedicalRecordService) Finalize(id uint) error {
	if err := s.inner.Finalize(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicalRecordFinalized,
		events.NewMedicalRecordFinalizedEvent(id),
	)
	return nil
}

func (s *eventedMedicalRecordService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicalRecordDeleted,
		events.NewMedicalRecordDeletedEvent(id, "soft_delete"),
	)
	return nil
}

func (s *eventedMedicalRecordService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicalRecordRestored,
		events.NewMedicalRecordRestoredEvent(id),
	)
	return nil
}

func (s *eventedMedicalRecordService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicalRecordDeleted,
		events.NewMedicalRecordDeletedEvent(id, "hard_delete"),
	)
	return nil
}
