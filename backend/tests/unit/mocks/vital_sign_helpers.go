package mocks

import (
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func NewTestVitalSignWithData(id uint, medicalRecordID uint, isDeleted bool) *models.VitalSign {
	now := time.Now()
	sbp := 120
	dbp := 80
	hr := 75
	temp := 36.5
	rr := 16
	os := 98.0
	weight := 65.0
	height := 170
	bmi := 22.49

	v := &models.VitalSign{
		ID:               id,
		MedicalRecordID:  medicalRecordID,
		MeasurementDate:  now,
		MeasurementTime:  now.Format("15:04:05"),
		SystolicBP:       &sbp,
		DiastolicBP:      &dbp,
		HeartRate:        &hr,
		BodyTemperature:  &temp,
		RespiratoryRate:  &rr,
		OxygenSaturation: &os,
		WeightKg:         &weight,
		HeightCm:         &height,
		BMI:              &bmi,
		Notes:            "Normal",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if isDeleted {
		v.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	}

	return v
}

func NewTestVitalSignList(count int) []models.VitalSign {
	var list []models.VitalSign
	for i := 1; i <= count; i++ {
		list = append(list, *NewTestVitalSignWithData(uint(i), 1, false))
	}
	return list
}

func NewTestVitalSignResponse(v *models.VitalSign) *dto.VitalSignResponse {
	return &dto.VitalSignResponse{
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
		RecordedAt:             v.MeasurementDate.Format(time.RFC3339),
		CreatedAt:              v.CreatedAt,
		UpdatedAt:              v.UpdatedAt,
	}
}

func NewVitalSignPaginationQuery(page, pageSize int) *dto.VitalSignPaginationQuery {
	return &dto.VitalSignPaginationQuery{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "recorded_at",
		SortDir:  "desc",
	}
}

func NewCreateVitalSignRequest(medicalRecordID uint) *dto.CreateVitalSignRequest {
	sbp := 120
	dbp := 80
	hr := 75
	temp := 36.5
	rr := 16
	os := 98.0
	weight := 65.0
	height := 170

	return &dto.CreateVitalSignRequest{
		MedicalRecordID:        medicalRecordID,
		RecordedAt:             time.Now().Format(time.RFC3339),
		BloodPressureSystolic:  &sbp,
		BloodPressureDiastolic: &dbp,
		HeartRate:              &hr,
		Temperature:            &temp,
		RespiratoryRate:        &rr,
		OxygenSaturation:       &os,
		WeightKG:               &weight,
		HeightCM:               &height,
	}
}

func NewUpdateVitalSignRequest() *dto.UpdateVitalSignRequest {
	sbp := 125
	return &dto.UpdateVitalSignRequest{
		BloodPressureSystolic: &sbp,
	}
}
