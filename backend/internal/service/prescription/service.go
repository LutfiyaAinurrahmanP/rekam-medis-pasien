package prescription

import (
	"errors"
	"math"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type PrescriptionService interface {
	List(query *dto.PrescriptionPaginationQuery) (*dto.PrescriptionListResponse, error)
	DeletedList(query *dto.PrescriptionPaginationQuery) (*dto.PrescriptionDeletedListResponse, error)
	FindByID(id uint) (*dto.PrescriptionResponse, error)
	FindByIDUnscoped(id uint) (*dto.PrescriptionResponse, error)
	Create(req *dto.CreatePrescriptionRequest) (*dto.PrescriptionResponse, error)
	Update(id uint, req *dto.UpdatePrescriptionRequest) (*dto.PrescriptionResponse, error)
	Dispense(id uint) (*dto.PrescriptionResponse, error)
	Cancel(id uint) (*dto.PrescriptionResponse, error)
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error

	ListItems(prescriptionID uint) ([]dto.PrescriptionItemResponse, error)
	FindItemByID(prescriptionID, itemID uint) (*dto.PrescriptionItemResponse, error)
	CreateItem(prescriptionID uint, req *dto.CreatePrescriptionItemRequest) (*dto.PrescriptionItemResponse, error)
	UpdateItem(prescriptionID, itemID uint, req *dto.UpdatePrescriptionItemRequest) (*dto.PrescriptionItemResponse, error)
	DeleteItem(prescriptionID, itemID uint) error
}

type prescriptionService struct {
	repo   repository.PrescriptionRepository
	config *config.Config
}

func NewPrescriptionService(repo repository.PrescriptionRepository, config *config.Config) PrescriptionService {
	return &prescriptionService{
		repo:   repo,
		config: config,
	}
}

func (s *prescriptionService) normalizeQuery(query *dto.PrescriptionPaginationQuery, defaultSortBy, defaultSortDir string) {
	if query.Page < 1 {
		query.Page = 1
	}

	if query.PageSize < 1 {
		query.PageSize = s.config.Pagination.DefaultPageSize
	}

	if query.PageSize > s.config.Pagination.MaxPageSize {
		query.PageSize = s.config.Pagination.MaxPageSize
	}

	if query.SortBy == "" {
		query.SortBy = defaultSortBy
	}

	if query.SortDir == "" {
		query.SortDir = defaultSortDir
	}
}

func (s *prescriptionService) toResponse(m *models.Prescription) *dto.PrescriptionResponse {
	if m == nil {
		return nil
	}

	var mrResp *dto.PrescriptionMedicalRecordResponse
	if m.MedicalRecord != nil {
		mrResp = &dto.PrescriptionMedicalRecordResponse{
			ID:        m.MedicalRecord.ID,
			VisitDate: m.MedicalRecord.VisitDate,
		}
	}

	var drResp *dto.PrescriptionDoctorResponse
	if m.Doctor != nil {
		specName := "Unspecified"
		if m.Doctor.Specialization.Name != "" {
			specName = m.Doctor.Specialization.Name
		}
		drResp = &dto.PrescriptionDoctorResponse{
			ID:             m.Doctor.ID,
			Name:           m.Doctor.FullName,
			Specialization: specName,
		}
	}

	items := make([]dto.PrescriptionItemResponse, 0, len(m.Items))
	for _, it := range m.Items {
		items = append(items, *s.toItemResponse(&it))
	}

	return &dto.PrescriptionResponse{
		ID:               m.ID,
		MedicalRecordID:  m.MedicalRecordID,
		MedicalRecord:    mrResp,
		DoctorID:         m.DoctorID,
		Doctor:           drResp,
		PrescriptionDate: m.PrescriptionDate,
		Notes:            m.Notes,
		Status:           m.Status,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		Items:            items,
	}
}

func (s *prescriptionService) toDeletedResponse(m *models.Prescription) *dto.DeletedPrescriptionResponse {
	resp := s.toResponse(m)
	deletedResp := &dto.DeletedPrescriptionResponse{
		PrescriptionResponse: *resp,
	}

	if m.DeletedAt.Valid {
		deletedResp.DeletedAt = &m.DeletedAt.Time
	}

	return deletedResp
}

func (s *prescriptionService) toItemResponse(m *models.PrescriptionItem) *dto.PrescriptionItemResponse {
	if m == nil {
		return nil
	}
	var medResp *dto.PrescriptionMedicineResponse
	if m.Medicine != nil {
		medResp = &dto.PrescriptionMedicineResponse{
			ID:   m.Medicine.ID,
			Name: m.Medicine.Name,
			Unit: m.Medicine.Unit,
		}
	}
	return &dto.PrescriptionItemResponse{
		ID:             m.ID,
		PrescriptionID: m.PrescriptionID,
		MedicineID:     m.MedicineID,
		Dosage:         m.Dosage,
		Frequency:      m.Frequency,
		DurationDays:   m.DurationDays,
		Quantity:       m.Quantity,
		Instructions:   m.Instructions,
		Medicine:       medResp,
	}
}

func (s *prescriptionService) List(query *dto.PrescriptionPaginationQuery) (*dto.PrescriptionListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")

	prescriptions, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PrescriptionResponse, len(prescriptions))
	for i, r := range prescriptions {
		responses[i] = *s.toResponse(&r)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.PrescriptionListResponse{
		Data: responses,
		Meta: dto.PrescriptionPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *prescriptionService) DeletedList(query *dto.PrescriptionPaginationQuery) (*dto.PrescriptionDeletedListResponse, error) {
	s.normalizeQuery(query, "created_at", "desc")

	prescriptions, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DeletedPrescriptionResponse, len(prescriptions))
	for i, r := range prescriptions {
		responses[i] = *s.toDeletedResponse(&r)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.PrescriptionDeletedListResponse{
		Data: responses,
		Meta: dto.PrescriptionPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *prescriptionService) FindByID(id uint) (*dto.PrescriptionResponse, error) {
	prescription, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(prescription), nil
}

func (s *prescriptionService) FindByIDUnscoped(id uint) (*dto.PrescriptionResponse, error) {
	prescription, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(prescription), nil
}

func (s *prescriptionService) Create(req *dto.CreatePrescriptionRequest) (*dto.PrescriptionResponse, error) {
	status := "pending"
	if req.Status != "" {
		status = req.Status
	}

	prescription := &models.Prescription{
		MedicalRecordID:  req.MedicalRecordID,
		DoctorID:         req.DoctorID,
		PrescriptionDate: req.PrescriptionDate,
		Notes:            req.Notes,
		Status:           status,
	}

	if err := s.repo.Create(prescription); err != nil {
		return nil, err
	}

	created, _ := s.repo.FindByID(prescription.ID)
	if created == nil {
		created = prescription
	}

	return s.toResponse(created), nil
}

func (s *prescriptionService) Update(id uint, req *dto.UpdatePrescriptionRequest) (*dto.PrescriptionResponse, error) {
	prescription, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Notes != nil {
		prescription.Notes = *req.Notes
	}

	if err := s.repo.Update(prescription); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(id)
	if updated == nil {
		updated = prescription
	}

	return s.toResponse(updated), nil
}

func (s *prescriptionService) Dispense(id uint) (*dto.PrescriptionResponse, error) {
	prescription, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if prescription.Status == "cancelled" {
		return nil, errors.New("cannot dispense a cancelled prescription")
	}

	prescription.Status = "dispensed"

	if err := s.repo.Update(prescription); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(id)
	return s.toResponse(updated), nil
}

func (s *prescriptionService) Cancel(id uint) (*dto.PrescriptionResponse, error) {
	prescription, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if prescription.Status == "dispensed" {
		return nil, errors.New("cannot cancel a dispensed prescription")
	}

	prescription.Status = "cancelled"

	if err := s.repo.Update(prescription); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindByID(id)
	return s.toResponse(updated), nil
}

func (s *prescriptionService) SoftDelete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}

func (s *prescriptionService) Restore(id uint) error {
	_, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	return s.repo.Restore(id)
}

func (s *prescriptionService) HardDelete(id uint) error {
	_, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	return s.repo.HardDelete(id)
}

func (s *prescriptionService) ListItems(prescriptionID uint) ([]dto.PrescriptionItemResponse, error) {
	p, err := s.repo.FindByID(prescriptionID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.PrescriptionItemResponse, len(p.Items))
	for i, it := range p.Items {
		items[i] = *s.toItemResponse(&it)
	}
	return items, nil
}

func (s *prescriptionService) FindItemByID(prescriptionID, itemID uint) (*dto.PrescriptionItemResponse, error) {
	item, err := s.repo.FindItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if item.PrescriptionID != prescriptionID {
		return nil, errors.New("prescription item not found in this prescription")
	}
	return s.toItemResponse(item), nil
}

func (s *prescriptionService) CreateItem(prescriptionID uint, req *dto.CreatePrescriptionItemRequest) (*dto.PrescriptionItemResponse, error) {
	_, err := s.repo.FindByID(prescriptionID)
	if err != nil {
		return nil, err
	}

	item := &models.PrescriptionItem{
		PrescriptionID: prescriptionID,
		MedicineID:     req.MedicineID,
		Dosage:         req.Dosage,
		Frequency:      req.Frequency,
		DurationDays:   req.DurationDays,
		Quantity:       req.Quantity,
		Instructions:   req.Instructions,
	}

	if err := s.repo.CreateItem(item); err != nil {
		return nil, err
	}

	created, _ := s.repo.FindItemByID(item.ID)
	if created == nil {
		created = item
	}

	return s.toItemResponse(created), nil
}

func (s *prescriptionService) UpdateItem(prescriptionID, itemID uint, req *dto.UpdatePrescriptionItemRequest) (*dto.PrescriptionItemResponse, error) {
	item, err := s.repo.FindItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if item.PrescriptionID != prescriptionID {
		return nil, errors.New("prescription item not found in this prescription")
	}

	if req.Dosage != nil {
		item.Dosage = *req.Dosage
	}
	if req.Frequency != nil {
		item.Frequency = *req.Frequency
	}
	if req.DurationDays != nil {
		item.DurationDays = *req.DurationDays
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.Instructions != nil {
		item.Instructions = *req.Instructions
	}

	if err := s.repo.UpdateItem(item); err != nil {
		return nil, err
	}

	updated, _ := s.repo.FindItemByID(itemID)
	if updated == nil {
		updated = item
	}

	return s.toItemResponse(updated), nil
}

func (s *prescriptionService) DeleteItem(prescriptionID, itemID uint) error {
	item, err := s.repo.FindItemByID(itemID)
	if err != nil {
		return err
	}
	if item.PrescriptionID != prescriptionID {
		return errors.New("prescription item not found in this prescription")
	}
	return s.repo.DeleteItem(itemID)
}
