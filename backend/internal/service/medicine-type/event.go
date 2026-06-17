package medicinetype

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventMedicineTypeService struct {
	inner     MedicineTypeService
	publisher kafka.EventPublisher
}

func NewEventMedicineTypeService(inner MedicineTypeService, publisher kafka.EventPublisher) MedicineTypeService {
	if publisher == nil {
		return inner
	}
	return &eventMedicineTypeService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventMedicineTypeService) List(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
	return s.inner.List(query)
}
func (s *eventMedicineTypeService) ActiveList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
	return s.inner.ActiveList(query)
}
func (s *eventMedicineTypeService) InactiveList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeListResponse, error) {
	return s.inner.InactiveList(query)
}
func (s *eventMedicineTypeService) DeletedList(query *dto.MedicineTypePaginationQuery) (*dto.MedicineTypeDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}
func (s *eventMedicineTypeService) FindByID(id uint) (*dto.MedicineTypeResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventMedicineTypeService) Create(req *dto.CreateMedicineTypeRequest) (*dto.MedicineTypeResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicineTypeCreated,
		events.NewMedicineTypeCreatedEvent(result.ID, result.Name, result.Code, result.Description, result.IsActive),
	)
	return result, nil
}
func (s *eventMedicineTypeService) Update(id uint, req *dto.UpdateMedicineTypeRequest) (*dto.MedicineTypeResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}
	s.publisher.PublishAsync(
		kafka.TopicMedicineTypeUpdated,
		events.NewMedicineTypeUpdatedEvent(result.ID, result.Name, result.Code, result.Description, result.IsActive, "update"),
	)
	return result, nil
}
func (s *eventMedicineTypeService) SoftDelete(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}
	name := ""
	code := ""
	if ds != nil {
		name = ds.Name
		code = ds.Code
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicineTypeDeleted,
		events.NewMedicineTypeDeletedEvent(id, name, code, "soft_delete"),
	)
	return nil
}
func (s *eventMedicineTypeService) Restore(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.Restore(id); err != nil {
		return err
	}
	name := ""
	code := ""
	if ds != nil {
		name = ds.Name
		code = ds.Code
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicineTypeRestored,
		events.NewMedicineTypeRestoredEvent(id, name, code),
	)
	return nil
}
func (s *eventMedicineTypeService) HardDelete(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	name := ""
	code := ""
	if ds != nil {
		name = ds.Name
		code = ds.Code
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicineTypeDeleted,
		events.NewMedicineTypeDeletedEvent(id, name, code, "hard_delete"),
	)
	return nil
}

func (s *eventMedicineTypeService) Activate(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.Activate(id); err != nil {
		return err
	}
	name := ""
	code := ""
	if ds != nil {
		name = ds.Name
		code = ds.Code
	}
	s.publisher.PublishAsync(
		kafka.TopicMedicineTypeUpdated,
		events.NewMedicineTypeUpdatedEvent(id, name, code, ds.Description, true, "activate"),
	)
	return nil
}

func (s *eventMedicineTypeService) Deactivate(id uint) error {
	ds, _ := s.inner.FindByID(id)
	if err := s.inner.Deactivate(id); err != nil {
		return err
	}
	name := ""
	code := ""
	if ds != nil {
		name = ds.Name
		code = ds.Code
	}
	s.publisher.PublishAsync(
		kafka.TopicMedicineTypeUpdated,
		events.NewMedicineTypeUpdatedEvent(id, name, code, ds.Description, false, "deactivate"),
	)
	return nil
}
