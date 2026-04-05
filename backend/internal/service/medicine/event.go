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