package vitalsign

import (
	"errors"
	"math"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type VitalSignService interface {
	List(query *dto.VitalSignPaginationQuery) (*dto.VitalSignListResponse, error)
	DeletedList(query *dto.VitalSignPaginationQuery) (*dto.VitalSignDeletedListResponse, error)
	FindByID(id uint) (*dto.VitalSignResponse, error)
	FindByIDUnscoped(id uint) (*dto.VitalSignResponse, error)
	Create(req *dto.CreateVitalSignRequest) (*dto.VitalSignResponse, error)
	Update(id uint, req *dto.UpdateVitalSignRequest) (*dto.VitalSignResponse, error)
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type vitalSignService struct {
	repo              repository.VitalSignRepository
	config            *config.Config
	medicalRecordRepo repository.MedicalRecordRepository
}

func NewVitalSignService(repo repository.VitalSignRepository, config *config.Config, medicalRecordRepo repository.MedicalRecordRepository) VitalSignService {
	return &vitalSignService{
		repo:              repo,
		config:            config,
		medicalRecordRepo: medicalRecordRepo,
	}
}

func (s *vitalSignService) normalizeQuery(query *dto.VitalSignPaginationQuery, defaultSortBy, defaultSortDir string) {
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

func (s *vitalSignService) toResponse(v *models.VitalSign) *dto.VitalSignResponse {
	if v == nil {
		return nil
	}

	resp := &dto.VitalSignResponse{
		ID:                     v.ID,
		MedicalRecordID:        v.MedicalRecordID,
		BloodPressureSystolic:  v.SystolicBP,
		BloodPressureDiastolic: v.DiastolicBP,
		HeartRate:              v.HeartRate,
		Temperature:            v.BodyTemperature,
		RespiratoryRate:        v.RespiratoryRate,
		OxygenSaturation:       v.OxygenSaturation,
		WeightKG:               v.WeightKg,
		HeightCM:               v.HeightCm,
		BMI:                    v.BMI,
		RecordedAt:             parseRecordedAt(v.MeasurementDate, v.MeasurementTime),
		CreatedAt:              v.CreatedAt,
		UpdatedAt:              v.UpdatedAt,
	}

	if v.MedicalRecord.ID != 0 {
		resp.MedicalRecord = &dto.VitalSignMedicalRecord{
			ID:             v.MedicalRecord.ID,
			VisitDate:      v.MedicalRecord.VisitDate,
			ChiefComplaint: v.MedicalRecord.ChiefComplaint,
		}
	}

	return resp
}

func (s *vitalSignService) toDeletedResponse(v *models.VitalSign) *dto.DeletedVitalSignResponse {
	if v == nil {
		return nil
	}
	resp := s.toResponse(v)
	deletedResp := &dto.DeletedVitalSignResponse{
		VitalSignResponse: *resp,
	}
	if v.DeletedAt.Valid {
		deletedResp.DeletedAt = &v.DeletedAt.Time
	}
	return deletedResp
}

func parseRecordedAt(date time.Time, timeStr string) string {
	if timeStr == "" {
		return date.Format(time.RFC3339)
	}
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return date.Format(time.RFC3339)
	}
	fullTime := time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	return fullTime.Format(time.RFC3339)
}

func (s *vitalSignService) List(query *dto.VitalSignPaginationQuery) (*dto.VitalSignListResponse, error) {
	s.normalizeQuery(query, "recorded_at", "desc")

	vitalSigns, total, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.VitalSignResponse, len(vitalSigns))
	for i, v := range vitalSigns {
		responses[i] = *s.toResponse(&v)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.VitalSignListResponse{
		Data: responses,
		Meta: dto.VitalSignPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *vitalSignService) DeletedList(query *dto.VitalSignPaginationQuery) (*dto.VitalSignDeletedListResponse, error) {
	s.normalizeQuery(query, "recorded_at", "desc")

	vitalSigns, total, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DeletedVitalSignResponse, len(vitalSigns))
	for i, v := range vitalSigns {
		responses[i] = *s.toDeletedResponse(&v)
	}

	totalPages := int(math.Ceil(float64(total) / float64(query.PageSize)))
	return &dto.VitalSignDeletedListResponse{
		Data: responses,
		Meta: dto.VitalSignPaginationMeta{
			Page:       query.Page,
			PageSize:   query.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *vitalSignService) FindByID(id uint) (*dto.VitalSignResponse, error) {
	v, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(v), nil
}

func (s *vitalSignService) FindByIDUnscoped(id uint) (*dto.VitalSignResponse, error) {
	v, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(v), nil
}

func (s *vitalSignService) Create(req *dto.CreateVitalSignRequest) (*dto.VitalSignResponse, error) {
	_, err := s.medicalRecordRepo.FindByID(req.MedicalRecordID)
	if err != nil {
		return nil, errors.New("medical record not found")
	}

	existing, _, err := s.repo.List(&dto.VitalSignPaginationQuery{
		MedicalRecordID: &req.MedicalRecordID,
		PageSize:        1,
	})
	if err == nil && existing != nil && len(existing) > 0 {
		return nil, errors.New("vital signs already recorded for this medical record")
	}

	recordedAt, err := time.Parse(time.RFC3339, req.RecordedAt)
	if err != nil {
		return nil, errors.New("invalid recorded_at format, must be ISO8601")
	}

	vitalSign := &models.VitalSign{
		MedicalRecordID:  req.MedicalRecordID,
		MeasurementDate:  recordedAt,
		MeasurementTime:  recordedAt.Format("15:04:05"),
		SystolicBP:       req.BloodPressureSystolic,
		DiastolicBP:      req.BloodPressureDiastolic,
		HeartRate:        req.HeartRate,
		BodyTemperature:  req.Temperature,
		RespiratoryRate:  req.RespiratoryRate,
		OxygenSaturation: req.OxygenSaturation,
		WeightKg:         req.WeightKG,
		HeightCm:         req.HeightCM,
	}

	// Calculate BMI implicitly through BeforeSave hook logic manually or let repo do it.
	if vitalSign.WeightKg != nil && vitalSign.HeightCm != nil && *vitalSign.HeightCm > 0 {
		heightM := float64(*vitalSign.HeightCm) / 100.0
		bmi := *vitalSign.WeightKg / (heightM * heightM)
		vitalSign.BMI = &bmi
	}

	if err := s.repo.Create(vitalSign); err != nil {
		return nil, err
	}

	return s.FindByID(vitalSign.ID)
}

func (s *vitalSignService) Update(id uint, req *dto.UpdateVitalSignRequest) (*dto.VitalSignResponse, error) {
	vitalSign, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.BloodPressureSystolic != nil {
		vitalSign.SystolicBP = req.BloodPressureSystolic
	}
	if req.BloodPressureDiastolic != nil {
		vitalSign.DiastolicBP = req.BloodPressureDiastolic
	}
	if req.HeartRate != nil {
		vitalSign.HeartRate = req.HeartRate
	}
	if req.Temperature != nil {
		vitalSign.BodyTemperature = req.Temperature
	}
	if req.RespiratoryRate != nil {
		vitalSign.RespiratoryRate = req.RespiratoryRate
	}
	if req.OxygenSaturation != nil {
		vitalSign.OxygenSaturation = req.OxygenSaturation
	}
	if req.WeightKG != nil {
		vitalSign.WeightKg = req.WeightKG
	}
	if req.HeightCM != nil {
		vitalSign.HeightCm = req.HeightCM
	}

	if vitalSign.WeightKg != nil && vitalSign.HeightCm != nil && *vitalSign.HeightCm > 0 {
		heightM := float64(*vitalSign.HeightCm) / 100.0
		bmi := *vitalSign.WeightKg / (heightM * heightM)
		vitalSign.BMI = &bmi
	} else {
		vitalSign.BMI = nil
	}

	if err := s.repo.Update(vitalSign); err != nil {
		return nil, err
	}

	return s.FindByID(vitalSign.ID)
}

func (s *vitalSignService) SoftDelete(id uint) error {
	_, err := s.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(id)
}

func (s *vitalSignService) Restore(id uint) error {
	_, err := s.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	return s.repo.Restore(id)
}

func (s *vitalSignService) HardDelete(id uint) error {
	_, err := s.FindByIDUnscoped(id)
	if err != nil {
		return err
	}
	return s.repo.HardDelete(id)
}
