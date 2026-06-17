package referral

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventReferralService struct {
	next      ReferralService
	publisher kafka.EventPublisher
}

func NewEventReferralService(next ReferralService, publisher kafka.EventPublisher) ReferralService {
	return &eventReferralService{
		next:      next,
		publisher: publisher,
	}
}

func (s *eventReferralService) publishEvent(topic string, event interface{}) {
	if s.publisher == nil {
		return
	}
	s.publisher.PublishAsync(topic, event)
}

func (s *eventReferralService) List(query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	return s.next.List(query)
}

func (s *eventReferralService) DeletedList(query dto.ReferralPaginationQuery) (*dto.ReferralDeletedListResponse, error) {
	return s.next.DeletedList(query)
}

func (s *eventReferralService) FindMyReferrals(patientID uint, status string) (*dto.ReferralMyListResponse, error) {
	return s.next.FindMyReferrals(patientID, status)
}

func (s *eventReferralService) FindByID(id uint) (*dto.ReferralResponse, error) {
	return s.next.FindByID(id)
}

func (s *eventReferralService) FindByIDUnscoped(id uint) (*dto.ReferralResponse, error) {
	return s.next.FindByIDUnscoped(id)
}

func (s *eventReferralService) FindByPatientID(patientID uint, query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	return s.next.FindByPatientID(patientID, query)
}

func (s *eventReferralService) FindByDoctorID(doctorID uint, query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	return s.next.FindByDoctorID(doctorID, query)
}

func (s *eventReferralService) Create(req dto.CreateReferralRequest) (*dto.ReferralResponse, error) {
	resp, err := s.next.Create(req)
	if err == nil && resp != nil {
		event := events.NewReferralCreatedEvent(
			resp.ID,
			resp.PatientID,
			resp.ReferringDoctorID,
			resp.ReferredToDoctorID,
			resp.ReferralNumber,
			resp.Status,
		)
		s.publishEvent(kafka.TopicReferralCreated, event)
	}
	return resp, err
}

func (s *eventReferralService) Update(id uint, req dto.UpdateReferralRequest) (*dto.ReferralResponse, error) {
	resp, err := s.next.Update(id, req)
	if err == nil && resp != nil {
		event := events.NewReferralUpdatedEvent(
			resp.ID,
			resp.PatientID,
			resp.ReferringDoctorID,
			resp.ReferredToDoctorID,
			resp.ReferralNumber,
			resp.Status,
		)
		s.publishEvent(kafka.TopicReferralUpdated, event)
	}
	return resp, err
}

func (s *eventReferralService) Accept(id uint, req dto.AcceptReferralRequest) (*dto.ReferralResponse, error) {
	resp, err := s.next.Accept(id, req)
	if err == nil && resp != nil {
		event := events.NewReferralStatusChangedEvent("referral.accepted", "accepted", resp.ID, resp.PatientID)
		s.publishEvent(kafka.TopicReferralAccepted, event)
	}
	return resp, err
}

func (s *eventReferralService) Reject(id uint, req dto.RejectReferralRequest) (*dto.ReferralResponse, error) {
	resp, err := s.next.Reject(id, req)
	if err == nil && resp != nil {
		event := events.NewReferralStatusChangedEvent("referral.rejected", "rejected", resp.ID, resp.PatientID)
		s.publishEvent(kafka.TopicReferralRejected, event)
	}
	return resp, err
}

func (s *eventReferralService) Complete(id uint, req dto.CompleteReferralRequest) (*dto.ReferralResponse, error) {
	resp, err := s.next.Complete(id, req)
	if err == nil && resp != nil {
		event := events.NewReferralStatusChangedEvent("referral.completed", "completed", resp.ID, resp.PatientID)
		s.publishEvent(kafka.TopicReferralCompleted, event)
	}
	return resp, err
}

func (s *eventReferralService) Cancel(id uint, req dto.CancelReferralRequest) (*dto.ReferralResponse, error) {
	resp, err := s.next.Cancel(id, req)
	if err == nil && resp != nil {
		event := events.NewReferralStatusChangedEvent("referral.cancelled", "cancelled", resp.ID, resp.PatientID)
		s.publishEvent(kafka.TopicReferralCancelled, event)
	}
	return resp, err
}

func (s *eventReferralService) SoftDelete(id uint) error {
	err := s.next.SoftDelete(id)
	if err == nil {
		event := events.NewReferralDeletedEvent(id, "delete")
		s.publishEvent(kafka.TopicReferralDeleted, event)
	}
	return err
}

func (s *eventReferralService) Restore(id uint) error {
	err := s.next.Restore(id)
	if err == nil {
		event := events.NewReferralRestoredEvent(id)
		s.publishEvent(kafka.TopicReferralRestored, event)
	}
	return err
}

func (s *eventReferralService) HardDelete(id uint) error {
	return s.next.HardDelete(id)
}
