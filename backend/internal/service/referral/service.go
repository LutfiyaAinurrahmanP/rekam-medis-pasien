package referral

import (
	"errors"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/dto"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/repository"
)

type ReferralService interface {
	List(query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error)
	DeletedList(query dto.ReferralPaginationQuery) (*dto.ReferralDeletedListResponse, error)
	FindMyReferrals(patientID uint, status string) (*dto.ReferralMyListResponse, error)
	FindByID(id uint) (*dto.ReferralResponse, error)
	FindByIDUnscoped(id uint) (*dto.ReferralResponse, error)
	FindByPatientID(patientID uint, query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error)
	FindByDoctorID(doctorID uint, query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error)
	Create(req dto.CreateReferralRequest) (*dto.ReferralResponse, error)
	Update(id uint, req dto.UpdateReferralRequest) (*dto.ReferralResponse, error)
	Accept(id uint, req dto.AcceptReferralRequest) (*dto.ReferralResponse, error)
	Reject(id uint, req dto.RejectReferralRequest) (*dto.ReferralResponse, error)
	Complete(id uint, req dto.CompleteReferralRequest) (*dto.ReferralResponse, error)
	Cancel(id uint, req dto.CancelReferralRequest) (*dto.ReferralResponse, error)
	SoftDelete(id uint) error
	Restore(id uint) error
	HardDelete(id uint) error
}

type referralService struct {
	repo repository.ReferralRepository
}

func NewReferralService(repo repository.ReferralRepository) ReferralService {
	return &referralService{repo: repo}
}

func mapModelToResponse(ref models.Referral) dto.ReferralResponse {
	resp := dto.ReferralResponse{
		ID:                     ref.ID,
		ReferralNumber:         ref.ReferralNumber,
		PatientID:              ref.PatientID,
		MedicalRecordID:        ref.MedicalRecordID,
		ReferringDoctorID:      ref.ReferringDoctorID,
		ReferralDate:           ref.ReferralDate,
		ReferralType:           ref.ReferralType,
		ReferredToDepartmentID: ref.ReferredToDepartmentID,
		ReferredToDoctorID:     ref.ReferredToDoctorID,
		Reason:                 ref.Reason,
		Diagnosis:              ref.Diagnosis,
		Priority:               ref.Priority,
		Status:                 ref.Status,
		CreatedAt:              ref.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              ref.UpdatedAt.Format(time.RFC3339),
	}

	if ref.ReferredToFacility != "" {
		resp.ReferredToFacility = &ref.ReferredToFacility
	}
	if ref.AcceptedAt != nil {
		accepted := ref.AcceptedAt.Format(time.RFC3339)
		resp.AcceptedAt = &accepted
	}
	if ref.CompletedAt != nil {
		completed := ref.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &completed
	}
	if ref.RejectionReason != "" {
		resp.RejectionReason = &ref.RejectionReason
	}
	if ref.CancellationReason != "" {
		resp.CancellationReason = &ref.CancellationReason
	}
	if ref.Notes != "" {
		resp.Notes = &ref.Notes
	}

	if ref.Patient != nil {
		resp.Patient = &dto.ReferralPatientInfo{
			ID:                  ref.Patient.ID,
			Name:                ref.Patient.FullName,
			DateOfBirth:         ref.Patient.DateOfBirth,
			MedicalRecordNumber: ref.Patient.PatientCode,
		}
	}
	if ref.MedicalRecord != nil {
		resp.MedicalRecord = &dto.ReferralMedicalRecordInfo{
			ID:             ref.MedicalRecord.ID,
			VisitDate:      ref.MedicalRecord.VisitDate,
			ChiefComplaint: ref.MedicalRecord.ChiefComplaint,
		}
	}
	if ref.ReferringDoctor != nil {
		specName := ""
		if ref.ReferringDoctor.Specialization.ID != 0 {
			specName = ref.ReferringDoctor.Specialization.Name
		}
		resp.ReferringDoctor = &dto.ReferralSimpleDoctor{
			ID:             ref.ReferringDoctor.ID,
			Name:           ref.ReferringDoctor.FullName,
			Specialization: specName,
		}
	}
	if ref.ReferredToDepartment != nil {
		resp.ReferredToDepartment = &dto.ReferralSimpleDepartment{
			ID:   ref.ReferredToDepartment.ID,
			Name: ref.ReferredToDepartment.Name,
		}
	}
	if ref.ReferredToDoctor != nil {
		specName := ""
		if ref.ReferredToDoctor.Specialization.ID != 0 {
			specName = ref.ReferredToDoctor.Specialization.Name
		}
		resp.ReferredToDoctor = &dto.ReferralSimpleDoctor{
			ID:             ref.ReferredToDoctor.ID,
			Name:           ref.ReferredToDoctor.FullName,
			Specialization: specName,
		}
	}

	return resp
}

func mapModelsToResponses(models []models.Referral) []dto.ReferralResponse {
	var resps []dto.ReferralResponse
	for _, m := range models {
		resps = append(resps, mapModelToResponse(m))
	}
	if resps == nil {
		resps = []dto.ReferralResponse{}
	}
	return resps
}

func (s *referralService) List(query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	referrals, meta, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}
	return &dto.ReferralListResponse{
		Data: mapModelsToResponses(referrals),
		Meta: meta,
	}, nil
}

func (s *referralService) DeletedList(query dto.ReferralPaginationQuery) (*dto.ReferralDeletedListResponse, error) {
	referrals, meta, err := s.repo.DeletedList(query)
	if err != nil {
		return nil, err
	}
	return &dto.ReferralDeletedListResponse{
		Data: mapModelsToResponses(referrals),
		Meta: meta,
	}, nil
}

func (s *referralService) FindMyReferrals(patientID uint, status string) (*dto.ReferralMyListResponse, error) {
	referrals, err := s.repo.FindMyReferrals(patientID, status)
	if err != nil {
		return nil, err
	}
	return &dto.ReferralMyListResponse{
		Data: mapModelsToResponses(referrals),
	}, nil
}

func (s *referralService) FindByID(id uint) (*dto.ReferralResponse, error) {
	referral, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if referral == nil {
		return nil, errors.New("referral not found")
	}
	resp := mapModelToResponse(*referral)
	return &resp, nil
}

func (s *referralService) FindByIDUnscoped(id uint) (*dto.ReferralResponse, error) {
	referral, err := s.repo.FindByIDUnscoped(id)
	if err != nil {
		return nil, err
	}
	if referral == nil {
		return nil, errors.New("referral not found")
	}
	resp := mapModelToResponse(*referral)
	return &resp, nil
}

func (s *referralService) FindByPatientID(patientID uint, query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	referrals, meta, err := s.repo.FindByPatientID(patientID, query)
	if err != nil {
		return nil, err
	}
	return &dto.ReferralListResponse{
		Data: mapModelsToResponses(referrals),
		Meta: meta,
	}, nil
}

func (s *referralService) FindByDoctorID(doctorID uint, query dto.ReferralPaginationQuery) (*dto.ReferralListResponse, error) {
	referrals, meta, err := s.repo.FindByDoctorID(doctorID, query)
	if err != nil {
		return nil, err
	}
	return &dto.ReferralListResponse{
		Data: mapModelsToResponses(referrals),
		Meta: meta,
	}, nil
}

func (s *referralService) Create(req dto.CreateReferralRequest) (*dto.ReferralResponse, error) {
	refNum := req.ReferralNumber
	if refNum == "" {
		generated, err := s.repo.GenerateReferralNumber()
		if err != nil {
			return nil, err
		}
		refNum = generated
	}

	referral := models.Referral{
		ReferralNumber:         refNum,
		PatientID:              req.PatientID,
		MedicalRecordID:        req.MedicalRecordID,
		ReferringDoctorID:      req.ReferringDoctorID,
		ReferralDate:           req.ReferralDate,
		ReferralType:           req.ReferralType,
		ReferredToDepartmentID: req.ReferredToDepartmentID,
		ReferredToDoctorID:     req.ReferredToDoctorID,
		ReferredToFacility:     req.ReferredToFacility,
		Reason:                 req.Reason,
		Diagnosis:              req.Diagnosis,
		Priority:               req.Priority,
		Notes:                  req.Notes,
		Status:                 "pending",
	}

	if err := s.repo.Create(&referral); err != nil {
		return nil, err
	}

	return s.FindByID(referral.ID)
}

func (s *referralService) Update(id uint, req dto.UpdateReferralRequest) (*dto.ReferralResponse, error) {
	referral, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if referral == nil {
		return nil, errors.New("referral not found")
	}

	if referral.Status != "pending" {
		return nil, errors.New("cannot update referral that is not in pending status")
	}

	if req.ReferralNumber != nil {
		referral.ReferralNumber = *req.ReferralNumber
	}
	if req.PatientID != nil {
		referral.PatientID = *req.PatientID
	}
	if req.MedicalRecordID != nil {
		referral.MedicalRecordID = *req.MedicalRecordID
	}
	if req.ReferringDoctorID != nil {
		referral.ReferringDoctorID = *req.ReferringDoctorID
	}
	if req.ReferralDate != nil {
		referral.ReferralDate = *req.ReferralDate
	}
	if req.ReferralType != nil {
		referral.ReferralType = *req.ReferralType
	}
	if req.ReferredToDepartmentID != nil {
		referral.ReferredToDepartmentID = req.ReferredToDepartmentID
	}
	if req.ReferredToDoctorID != nil {
		referral.ReferredToDoctorID = req.ReferredToDoctorID
	}
	if req.ReferredToFacility != nil {
		referral.ReferredToFacility = *req.ReferredToFacility
	}
	if req.Reason != nil {
		referral.Reason = *req.Reason
	}
	if req.Diagnosis != nil {
		referral.Diagnosis = *req.Diagnosis
	}
	if req.Priority != nil {
		referral.Priority = *req.Priority
	}
	if req.Notes != nil {
		referral.Notes = *req.Notes
	}

	if err := s.repo.Update(referral); err != nil {
		return nil, err
	}

	return s.FindByID(id)
}

func (s *referralService) Accept(id uint, req dto.AcceptReferralRequest) (*dto.ReferralResponse, error) {
	referral, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if referral == nil {
		return nil, errors.New("referral not found")
	}
	if referral.Status != "pending" {
		return nil, errors.New("only pending referrals can be accepted")
	}

	referral.Status = "accepted"
	now := time.Now()
	referral.AcceptedAt = &now
	if req.Notes != "" {
		referral.Notes = req.Notes
	}

	if err := s.repo.Update(referral); err != nil {
		return nil, err
	}
	return s.FindByID(id)
}

func (s *referralService) Reject(id uint, req dto.RejectReferralRequest) (*dto.ReferralResponse, error) {
	referral, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if referral == nil {
		return nil, errors.New("referral not found")
	}
	if referral.Status != "pending" {
		return nil, errors.New("only pending referrals can be rejected")
	}

	referral.Status = "rejected"
	referral.RejectionReason = req.RejectionReason

	if err := s.repo.Update(referral); err != nil {
		return nil, err
	}
	return s.FindByID(id)
}

func (s *referralService) Complete(id uint, req dto.CompleteReferralRequest) (*dto.ReferralResponse, error) {
	referral, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if referral == nil {
		return nil, errors.New("referral not found")
	}
	if referral.Status != "accepted" {
		return nil, errors.New("only accepted referrals can be completed")
	}

	referral.Status = "completed"
	now := time.Now()
	referral.CompletedAt = &now
	if req.Notes != "" {
		referral.Notes = req.Notes
	}

	if err := s.repo.Update(referral); err != nil {
		return nil, err
	}
	return s.FindByID(id)
}

func (s *referralService) Cancel(id uint, req dto.CancelReferralRequest) (*dto.ReferralResponse, error) {
	referral, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if referral == nil {
		return nil, errors.New("referral not found")
	}
	if referral.Status == "completed" || referral.Status == "cancelled" {
		return nil, errors.New("cannot cancel an already completed or cancelled referral")
	}

	referral.Status = "cancelled"
	if req.CancellationReason != "" {
		referral.CancellationReason = req.CancellationReason
	}

	if err := s.repo.Update(referral); err != nil {
		return nil, err
	}
	return s.FindByID(id)
}

func (s *referralService) SoftDelete(id uint) error {
	return s.repo.SoftDelete(id)
}

func (s *referralService) Restore(id uint) error {
	return s.repo.Restore(id)
}

func (s *referralService) HardDelete(id uint) error {
	return s.repo.HardDelete(id)
}
