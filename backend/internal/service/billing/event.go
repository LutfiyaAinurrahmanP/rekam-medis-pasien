package billing

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventBillingService struct {
	next      BillingService
	publisher kafka.EventPublisher
}

func NewEventBillingService(next BillingService, publisher kafka.EventPublisher) BillingService {
	return &eventBillingService{
		next:      next,
		publisher: publisher,
	}
}

func (s *eventBillingService) publish(topic string, event interface{}) {
	s.publisher.PublishAsync(topic, event)
}

func (s *eventBillingService) List(query dto.BillingPaginationQuery) (*dto.BillingListResponse, error) {
	return s.next.List(query)
}

func (s *eventBillingService) DeletedList(query dto.BillingPaginationQuery) (*dto.BillingDeletedListResponse, error) {
	return s.next.DeletedList(query)
}

func (s *eventBillingService) FindByID(id uint) (*dto.BillingResponse, error) {
	return s.next.FindByID(id)
}

func (s *eventBillingService) FindByIDUnscoped(id uint) (*dto.BillingResponse, error) {
	return s.next.FindByIDUnscoped(id)
}

func (s *eventBillingService) FindByInvoiceNumber(invoice string) (*dto.BillingResponse, error) {
	return s.next.FindByInvoiceNumber(invoice)
}

func (s *eventBillingService) FindByPatientID(patientID uint, query dto.BillingPaginationQuery) (*dto.BillingListResponse, error) {
	return s.next.FindByPatientID(patientID, query)
}

func (s *eventBillingService) Create(req dto.CreateBillingRequest) (*dto.BillingResponse, error) {
	resp, err := s.next.Create(req)
	if err == nil {
		s.publish(kafka.TopicBillingCreated, events.NewBillingCreatedEvent(
			resp.ID, resp.PatientID, resp.InvoiceNumber, resp.TotalAmount,
		))
	}
	return resp, err
}

func (s *eventBillingService) Update(id uint, req dto.UpdateBillingRequest) (*dto.BillingResponse, error) {
	resp, err := s.next.Update(id, req)
	if err == nil {
		s.publish(kafka.TopicBillingUpdated, events.NewBillingUpdatedEvent(
			resp.ID, resp.PatientID, resp.InvoiceNumber, resp.Status, resp.TotalAmount, resp.PaidAmount,
		))
	}
	return resp, err
}

func (s *eventBillingService) RecordPayment(id uint, req dto.RecordPaymentRequest) (*dto.BillingResponse, error) {
	resp, err := s.next.RecordPayment(id, req)
	if err == nil {
		s.publish(kafka.TopicBillingPaid, events.NewBillingStatusChangedEvent(
			kafka.TopicBillingPaid, resp.Status, resp.ID, resp.PatientID, resp.InvoiceNumber, resp.PaidAmount,
		))
	}
	return resp, err
}

func (s *eventBillingService) Cancel(id uint) (*dto.BillingResponse, error) {
	resp, err := s.next.Cancel(id)
	if err == nil {
		s.publish(kafka.TopicBillingCancelled, events.NewBillingStatusChangedEvent(
			kafka.TopicBillingCancelled, "cancelled", resp.ID, resp.PatientID, resp.InvoiceNumber, resp.PaidAmount,
		))
	}
	return resp, err
}

func (s *eventBillingService) Delete(id uint) error {
	err := s.next.Delete(id)
	if err == nil {
		s.publish(kafka.TopicBillingDeleted, events.NewBillingDeletedEvent(id, "delete"))
	}
	return err
}

func (s *eventBillingService) Restore(id uint) (*dto.BillingResponse, error) {
	resp, err := s.next.Restore(id)
	if err == nil {
		s.publish(kafka.TopicBillingRestored, events.NewBillingRestoredEvent(id))
	}
	return resp, err
}

func (s *eventBillingService) HardDelete(id uint) error {
	err := s.next.HardDelete(id)
	if err == nil {
		s.publish(kafka.TopicBillingDeleted, events.NewBillingDeletedEvent(id, "hard_delete"))
	}
	return err
}

func (s *eventBillingService) ListItems(billingID uint, query dto.BillingItemPaginationQuery) ([]dto.BillingItemResponse, error) {
	return s.next.ListItems(billingID, query)
}

func (s *eventBillingService) FindItemByID(billingID, itemID uint) (*dto.BillingItemResponse, error) {
	return s.next.FindItemByID(billingID, itemID)
}

func (s *eventBillingService) CreateItem(billingID uint, req dto.CreateBillingItemRequest) (*dto.BillingItemResponse, error) {
	resp, err := s.next.CreateItem(billingID, req)
	if err == nil {
		billing, _ := s.next.FindByID(billingID)
		if billing != nil {
			s.publish(kafka.TopicBillingUpdated, events.NewBillingUpdatedEvent(
				billing.ID, billing.PatientID, billing.InvoiceNumber, billing.Status, billing.TotalAmount, billing.PaidAmount,
			))
		}
	}
	return resp, err
}

func (s *eventBillingService) UpdateItem(billingID, itemID uint, req dto.UpdateBillingItemRequest) (*dto.BillingItemResponse, error) {
	resp, err := s.next.UpdateItem(billingID, itemID, req)
	if err == nil {
		billing, _ := s.next.FindByID(billingID)
		if billing != nil {
			s.publish(kafka.TopicBillingUpdated, events.NewBillingUpdatedEvent(
				billing.ID, billing.PatientID, billing.InvoiceNumber, billing.Status, billing.TotalAmount, billing.PaidAmount,
			))
		}
	}
	return resp, err
}

func (s *eventBillingService) DeleteItem(billingID, itemID uint) error {
	err := s.next.DeleteItem(billingID, itemID)
	if err == nil {
		billing, _ := s.next.FindByID(billingID)
		if billing != nil {
			s.publish(kafka.TopicBillingUpdated, events.NewBillingUpdatedEvent(
				billing.ID, billing.PatientID, billing.InvoiceNumber, billing.Status, billing.TotalAmount, billing.PaidAmount,
			))
		}
	}
	return err
}
