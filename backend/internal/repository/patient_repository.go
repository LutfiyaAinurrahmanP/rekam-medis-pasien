package repository

import (
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

type PatientRepository interface {
	List(query *dto.PatientPaginationQuery) ([]models.Patient, int64, error)
	DeleteList(query *dto.PatientPaginationQuery) ([]models.Patient, int64, error)
	FindById(id uint) (*models.Patient, error)
	FindByCode(code string) (*models.Patient, error)
	Create(patient *models.Patient) error
	Update(patient *models.Patient) error
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
	IsCodeExists(code string, excludeID ...uint) (bool, error)
}

type patientRepository struct {
  db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{
		db: db,
	}
}

func (r patientRepository) List(query *dto.PatientPaginationQuery) ([]models.Patient, int64, error) {
        panic("not implemented") // TODO: Implement
}

func (r patientRepository) DeleteList(query *dto.PatientPaginationQuery) ([]models.Patient, int64, error) {
        panic("not implemented") // TODO: Implement
}

func (r patientRepository) FindById(id uint) (*models.Patient, error) {
        panic("not implemented") // TODO: Implement
}

func (r patientRepository) FindByCode(code string) (*models.Patient, error) {
        panic("not implemented") // TODO: Implement
}

func (r patientRepository) Create(patient *models.Patient) error {   
        panic("not implemented") // TODO: Implement
}

func (r patientRepository) Update(patient *models.Patient) error {   
        panic("not implemented") // TODO: Implement
}

func (r patientRepository) SoftDelete(id uint) error {
        panic("not implemented") // TODO: Implement
}

func (r patientRepository) Restore(id uint) error {
        panic("not implemented") // TODO: Implement
}

func (r patientRepository) HardDelete(id uint) error {
        panic("not implemented") // TODO: Implement
}

func (r patientRepository) IsCodeExists(code string, excludeID ...uint) (bool, error) {
        panic("not implemented") // TODO: Implement
}