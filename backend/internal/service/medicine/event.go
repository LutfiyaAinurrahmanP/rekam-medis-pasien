package medicine

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

type eventMedicineService struct {
  inner MedicineService
  publisher kafka.EventPublisher
}

func NewMedicineEventService(inner MedicineService, publisher kafka.EventPublisher) MedicineService {
	if publisher == nil {
		return inner
	}
	return &eventMedicineService{
		inner: inner,
		publisher: publisher,
	}
}

func (s *eventMedicineService) List(query *dto.MedicinePaginationQuery) (*dto.MedicineListResponse, error) {
	return s.inner.List(query)
}

func (s *eventMedicineService) DeletedList(query *dto.MedicinePaginationQuery) (*dto.MedicineDeletedListResponse, error) {
	return s.inner.DeletedList(query)
}

func (s *eventMedicineService) ListByAvailable(query *dto.MedicinePaginationQuery) (*dto.MedicineAvailableResponse, error) {
	return s.inner.ListByAvailable(query)
}

func (s *eventMedicineService) ListByLowStock(query *dto.MedicinePaginationQuery) (*dto.MedicineLowStockResponse, error) {
	return s.inner.ListByLowStock(query)
}

func (s *eventMedicineService) ListByOutStock(query *dto.MedicinePaginationQuery) (*dto.MedicineOutOfStockResponse, error) {
	return s.inner.ListByOutStock(query)
}

func (s *eventMedicineService) ListByInactive(query *dto.MedicinePaginationQuery) (*dto.MedicineInactiveResponse, error) {
	return  s.inner.ListByInactive(query)
}

func (s *eventMedicineService) FindByID(id uint) (*dto.MedicineResponse, error) {
	return s.inner.FindByID(id)
} 

func (s *eventMedicineService) FindByName(name string) (*dto.MedicineResponse, error) {
	return s.inner.FindByName(name)
}

func (s *eventMedicineService) ListByType(query *dto.MedicinePaginationQuery) (*dto.MedicineByTypeResponse, error) {
	return s.inner.ListByType(query)
}