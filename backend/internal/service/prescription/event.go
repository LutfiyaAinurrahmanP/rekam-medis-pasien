package prescription

import (
	"context"
	"encoding/json"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventPrescriptionService struct {
	next     PrescriptionService
	producer kafka.EventPublisher
}

func NewEventPrescriptionService(next PrescriptionService, producer kafka.EventPublisher) PrescriptionService {
	return &eventPrescriptionService{
		next:     next,
		producer: producer,
	}
}

func (s *eventPrescriptionService) publish(topic string, event interface{}) {

	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	if err := s.producer.Publish(context.Background(), topic, data); err != nil {
		log.Printf("Failed to publish event to topic %s: %v", topic, err)
	}
}

func (s *eventPrescriptionService) List(query *dto.PrescriptionPaginationQuery) (*dto.PrescriptionListResponse, error) {
	return s.next.List(query)
}

func (s *eventPrescriptionService) DeletedList(query *dto.PrescriptionPaginationQuery) (*dto.PrescriptionDeletedListResponse, error) {
	return s.next.DeletedList(query)
}

func (s *eventPrescriptionService) FindByID(id uint) (*dto.PrescriptionResponse, error) {
	return s.next.FindByID(id)
}

func (s *eventPrescriptionService) FindByIDUnscoped(id uint) (*dto.PrescriptionResponse, error) {
	return s.next.FindByIDUnscoped(id)
}

func (s *eventPrescriptionService) Create(req *dto.CreatePrescriptionRequest) (*dto.PrescriptionResponse, error) {
	resp, err := s.next.Create(req)
	if err == nil && resp != nil {
		s.publish(kafka.TopicPrescriptionCreated, events.NewPrescriptionCreatedEvent(
			resp.ID, resp.MedicalRecordID, resp.DoctorID, resp.PrescriptionDate, resp.Notes, resp.Status,
		))
	}
	return resp, err
}

func (s *eventPrescriptionService) Update(id uint, req *dto.UpdatePrescriptionRequest) (*dto.PrescriptionResponse, error) {
	resp, err := s.next.Update(id, req)
	if err == nil && resp != nil {
		s.publish(kafka.TopicPrescriptionUpdated, events.NewPrescriptionUpdatedEvent(
			resp.ID, resp.MedicalRecordID, resp.DoctorID, resp.PrescriptionDate, resp.Notes, resp.Status, "update",
		))
	}
	return resp, err
}

func (s *eventPrescriptionService) Dispense(id uint) (*dto.PrescriptionResponse, error) {
	resp, err := s.next.Dispense(id)
	if err == nil && resp != nil {
		s.publish(kafka.TopicPrescriptionDispensed, events.NewPrescriptionDispensedEvent(
			resp.ID, resp.Status,
		))
	}
	return resp, err
}

func (s *eventPrescriptionService) Cancel(id uint) (*dto.PrescriptionResponse, error) {
	resp, err := s.next.Cancel(id)
	if err == nil && resp != nil {
		s.publish(kafka.TopicPrescriptionCancelled, events.NewPrescriptionCancelledEvent(
			resp.ID, resp.Status,
		))
	}
	return resp, err
}

func (s *eventPrescriptionService) SoftDelete(id uint) error {
	err := s.next.SoftDelete(id)
	if err == nil {
		s.publish(kafka.TopicPrescriptionDeleted, events.NewPrescriptionDeletedEvent(
			id, "soft_delete",
		))
	}
	return err
}

func (s *eventPrescriptionService) Restore(id uint) error {
	err := s.next.Restore(id)
	if err == nil {
		s.publish(kafka.TopicPrescriptionRestored, events.NewPrescriptionRestoredEvent(
			id,
		))
	}
	return err
}

func (s *eventPrescriptionService) HardDelete(id uint) error {
	err := s.next.HardDelete(id)
	if err == nil {
		s.publish(kafka.TopicPrescriptionDeleted, events.NewPrescriptionDeletedEvent(
			id, "hard_delete",
		))
	}
	return err
}

func (s *eventPrescriptionService) ListItems(prescriptionID uint) ([]dto.PrescriptionItemResponse, error) {
	return s.next.ListItems(prescriptionID)
}

func (s *eventPrescriptionService) FindItemByID(prescriptionID, itemID uint) (*dto.PrescriptionItemResponse, error) {
	return s.next.FindItemByID(prescriptionID, itemID)
}

func (s *eventPrescriptionService) CreateItem(prescriptionID uint, req *dto.CreatePrescriptionItemRequest) (*dto.PrescriptionItemResponse, error) {
	resp, err := s.next.CreateItem(prescriptionID, req)
	if err == nil && resp != nil {
		s.publish(kafka.TopicPrescriptionItemCreated, events.NewPrescriptionItemCreatedEvent(
			resp.ID, resp.PrescriptionID, resp.MedicineID, resp.Dosage, resp.Frequency, resp.DurationDays, resp.Quantity, resp.Instructions,
		))
	}
	return resp, err
}

func (s *eventPrescriptionService) UpdateItem(prescriptionID, itemID uint, req *dto.UpdatePrescriptionItemRequest) (*dto.PrescriptionItemResponse, error) {
	resp, err := s.next.UpdateItem(prescriptionID, itemID, req)
	if err == nil && resp != nil {
		s.publish(kafka.TopicPrescriptionItemUpdated, events.NewPrescriptionItemUpdatedEvent(
			resp.ID, resp.PrescriptionID, resp.MedicineID, resp.Dosage, resp.Frequency, resp.DurationDays, resp.Quantity, resp.Instructions, "update",
		))
	}
	return resp, err
}

func (s *eventPrescriptionService) DeleteItem(prescriptionID, itemID uint) error {
	err := s.next.DeleteItem(prescriptionID, itemID)
	if err == nil {
		s.publish(kafka.TopicPrescriptionItemDeleted, events.NewPrescriptionItemDeletedEvent(
			itemID, "delete",
		))
	}
	return err
}
