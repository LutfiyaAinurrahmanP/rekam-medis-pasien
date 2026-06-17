package medicine

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventMedicineService struct {
	inner     MedicineService
	publisher kafka.EventPublisher
}

func NewEventMedicineService(inner MedicineService, publisher kafka.EventPublisher) MedicineService {
	if publisher == nil {
		return inner
	}
	return &eventMedicineService{
		inner:     inner,
		publisher: publisher,
	}
}

func (s *eventMedicineService) List(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	return s.inner.List(query)
}

func (s *eventMedicineService) DeletedList(query *dto.MedicinePaginationQuery) (*dto.MedicineDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}

func (s *eventMedicineService) AvailableList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	return s.inner.AvailableList(query)
}

func (s *eventMedicineService) LowStockList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	return s.inner.LowStockList(query)
}

func (s *eventMedicineService) OutStockList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	return s.inner.OutStockList(query)
}

func (s *eventMedicineService) ActiveList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	return s.inner.ActiveList(query)
}

func (s *eventMedicineService) InactiveList(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	return s.inner.InactiveList(query)
}

func (s *eventMedicineService) FindByID(id uint) (*dto.MedicineResponse, error) {
	return s.inner.FindByID(id)
}

func (s *eventMedicineService) FindByIDUnscoped(id uint) (*dto.DeletedMedicineResponse, error) {
	return s.inner.FindByIDUnscoped(id)
}

func (s *eventMedicineService) Create(req *dto.CreateMedicineRequest) (*dto.MedicineResponse, error) {
	result, err := s.inner.Create(req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicineCreated,
		events.NewMedicineCreatedEvent(
			result.ID, result.Name, result.GenericName, result.BrandName, result.MedicineTypeID,
			result.Strength, result.Manufacturer, result.Unit, result.StockQuantity, result.Price, result.IsActive,
		),
	)
	return result, nil
}

func (s *eventMedicineService) Update(id uint, req *dto.UpdateMedicineRequest) (*dto.MedicineResponse, error) {
	result, err := s.inner.Update(id, req)
	if err != nil {
		return nil, err
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicineUpdated,
		events.NewMedicineUpdatedEvent(
			result.ID, result.Name, result.GenericName, result.BrandName, result.MedicineTypeID,
			result.Strength, result.Manufacturer, result.Unit, result.StockQuantity, result.Price, result.IsActive, "update",
		),
	)
	return result, nil
}

func (s *eventMedicineService) AddStock(id uint, req *dto.AddStockRequest) error {
	m, _ := s.inner.FindByID(id)
	if err := s.inner.AddStock(id, req); err != nil {
		return err
	}

	if m != nil {
		// New stock is old stock + added quantity
		newStock := m.StockQuantity + req.Quantity
		s.publisher.PublishAsync(
			kafka.TopicMedicineStockAdded,
			events.NewMedicineStockAddedEvent(id, m.Name, newStock, req.Quantity),
		)
	}
	return nil
}

func (s *eventMedicineService) ReduceStock(id uint, req *dto.ReduceStockRequest) error {
	m, _ := s.inner.FindByID(id)
	if err := s.inner.ReduceStock(id, req); err != nil {
		return err
	}

	if m != nil {
		// New stock is old stock - reduced quantity
		newStock := m.StockQuantity - req.Quantity
		s.publisher.PublishAsync(
			kafka.TopicMedicineStockReduced,
			events.NewMedicineStockReducedEvent(id, m.Name, newStock, req.Quantity),
		)
	}
	return nil
}

func (s *eventMedicineService) Activate(id uint) error {
	m, _ := s.inner.FindByID(id)
	if err := s.inner.Activate(id); err != nil {
		return err
	}

	if m != nil {
		s.publisher.PublishAsync(
			kafka.TopicMedicineActivated,
			events.NewMedicineActivatedEvent(id, m.Name),
		)
	}
	return nil
}

func (s *eventMedicineService) Deactivate(id uint, req *dto.DeactivateMedicineRequest) error {
	m, _ := s.inner.FindByID(id)
	if err := s.inner.Deactivate(id, req); err != nil {
		return err
	}

	if m != nil {
		s.publisher.PublishAsync(
			kafka.TopicMedicineDeactivated,
			events.NewMedicineDeactivatedEvent(id, m.Name),
		)
	}
	return nil
}

func (s *eventMedicineService) SoftDelete(id uint) error {
	m, _ := s.inner.FindByID(id)
	if err := s.inner.SoftDelete(id); err != nil {
		return err
	}

	name := ""
	if m != nil {
		name = m.Name
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicineDeleted,
		events.NewMedicineDeletedEvent(id, name, "soft_delete"),
	)
	return nil
}

func (s *eventMedicineService) Restore(id uint) error {
	m, _ := s.inner.FindByIDUnscoped(id)
	if err := s.inner.Restore(id); err != nil {
		return err
	}

	name := ""
	if m != nil {
		name = m.Name
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicineRestored,
		events.NewMedicineRestoredEvent(id, name),
	)
	return nil
}

func (s *eventMedicineService) HardDelete(id uint) error {
	m, _ := s.inner.FindByIDUnscoped(id)
	if err := s.inner.HardDelete(id); err != nil {
		return err
	}
	name := ""
	if m != nil {
		name = m.Name
	}

	s.publisher.PublishAsync(
		kafka.TopicMedicineDeleted,
		events.NewMedicineDeletedEvent(id, name, "hard_delete"),
	)
	return nil
}
