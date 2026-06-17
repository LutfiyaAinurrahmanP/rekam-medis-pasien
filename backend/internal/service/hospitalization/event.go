package hospitalization

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventedHospitalizationService struct {
	inner     HospitalizationService
	publisher kafka.EventPublisher
}

func NewEventedHospitalizationService(inner HospitalizationService, publisher kafka.EventPublisher) HospitalizationService {
	if publisher == nil {
		return inner
	}
	return &eventedHospitalizationService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventedHospitalizationService) List(query *dto.HospitalizationPaginationQuery) (*dto.HospitalizationListResponse, error) {
	return s.inner.List(query)
}

func (s *eventedHospitalizationService) DeletedList(query *dto.HospitalizationPaginationQuery) (*dto.HospitalizationDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}

func (s *eventedHospitalizationService) FindByID(id uint) (*dto.HospitalizationResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventedHospitalizationService) FindByIDUnscoped(id uint) (*dto.HospitalizationResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *eventedHospitalizationService) Create(req *dto.CreateHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	resp, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicHospitalizationCreated,
		events.NewHospitalizationCreatedEvent(
			resp.ID,
			resp.PatientID,
			resp.DoctorID,
			resp.RoomID,
			resp.AdmissionDate,
			resp.AdmissionReason,
			resp.Status,
		),
	)
	return resp, nil
}

func (s *eventedHospitalizationService) Update(id uint, req *dto.UpdateHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	resp, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicHospitalizationUpdated,
		events.NewHospitalizationUpdatedEvent(
			resp.ID,
			resp.PatientID,
			resp.DoctorID,
			resp.RoomID,
			resp.AdmissionDate,
			resp.AdmissionReason,
			resp.Status,
			"update",
		),
	)
	return resp, nil
}

func (s *eventedHospitalizationService) Discharge(id uint, req *dto.DischargeHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	resp, err := s.inner.Discharge(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicHospitalizationDischarged,
		events.NewHospitalizationDischargedEvent(id, req.DischargeSummary),
	)
	return resp, nil
}

func (s *eventedHospitalizationService) Transfer(id uint, req *dto.TransferHospitalizationRequest) (*dto.HospitalizationResponse, error) {
	resp, err := s.inner.Transfer(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicHospitalizationTransferred,
		events.NewHospitalizationTransferredEvent(id, req.Notes),
	)
	return resp, nil
}

func (s *eventedHospitalizationService) Activate(id uint) (*dto.HospitalizationResponse, error) {
	resp, err := s.inner.Activate(id)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicHospitalizationUpdated,
		events.NewHospitalizationUpdatedEvent(
			resp.ID,
			resp.PatientID,
			resp.DoctorID,
			resp.RoomID,
			resp.AdmissionDate,
			resp.AdmissionReason,
			resp.Status,
			"activate",
		),
	)
	return resp, nil
}

func (s *eventedHospitalizationService) Deactivate(id uint) (*dto.HospitalizationResponse, error) {
	resp, err := s.inner.Deactivate(id)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicHospitalizationUpdated,
		events.NewHospitalizationUpdatedEvent(
			resp.ID,
			resp.PatientID,
			resp.DoctorID,
			resp.RoomID,
			resp.AdmissionDate,
			resp.AdmissionReason,
			resp.Status,
			"deactivate",
		),
	)
	return resp, nil
}

func (s *eventedHospitalizationService) SoftDelete(id uint) error {
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicHospitalizationDeleted,
		events.NewHospitalizationDeletedEvent(id, "soft_delete"),
	)
	return nil
}

func (s *eventedHospitalizationService) Restore(id uint) error {
	if err := s.inner.Restore(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicHospitalizationRestored,
		events.NewHospitalizationRestoredEvent(id),
	)
	return nil
}

func (s *eventedHospitalizationService) HardDelete(id uint) error {
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}

	s.publisher.PublishAsync(
		kafka.TopicHospitalizationDeleted,
		events.NewHospitalizationDeletedEvent(id, "hard_delete"),
	)
	return nil
}
