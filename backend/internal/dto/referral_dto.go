package dto

// ─── Shared Referral Response Types ───────────────────────────────────────

type ReferralSimpleDoctor struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization,omitempty"`
}

type ReferralSimpleDepartment struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ReferralPatientInfo struct {
	ID                  uint   `json:"id"`
	Name                string `json:"name"`
	DateOfBirth         string `json:"date_of_birth,omitempty"`
	MedicalRecordNumber string `json:"medical_record_number,omitempty"`
}

type ReferralMedicalRecordInfo struct {
	ID             uint   `json:"id"`
	VisitDate      string `json:"visit_date"`
	ChiefComplaint string `json:"chief_complaint,omitempty"`
}

// ─── Responses ────────────────────────────────────────────────────────────

type ReferralResponse struct {
	ID                     uint                      `json:"id"`
	ReferralNumber         string                    `json:"referral_number"`
	PatientID              uint                      `json:"patient_id"`
	MedicalRecordID        uint                      `json:"medical_record_id"`
	ReferringDoctorID      uint                      `json:"referring_doctor_id"`
	ReferralDate           string                    `json:"referral_date"`
	ReferralType           string                    `json:"referral_type"`
	ReferredToDepartmentID *uint                     `json:"referred_to_department_id"`
	ReferredToDoctorID     *uint                     `json:"referred_to_doctor_id"`
	ReferredToFacility     *string                   `json:"referred_to_facility"`
	Reason                 string                    `json:"reason"`
	Diagnosis              string                    `json:"diagnosis"`
	Priority               string                    `json:"priority"`
	Status                 string                    `json:"status"`
	AcceptedAt             *string                   `json:"accepted_at"`
	CompletedAt            *string                   `json:"completed_at"`
	RejectionReason        *string                   `json:"rejection_reason"`
	CancellationReason     *string                   `json:"cancellation_reason"`
	Notes                  *string                   `json:"notes"`
	CreatedAt              string                    `json:"created_at"`
	UpdatedAt              string                    `json:"updated_at"`
	Patient                *ReferralPatientInfo      `json:"patient,omitempty"`
	MedicalRecord          *ReferralMedicalRecordInfo`json:"medical_record,omitempty"`
	ReferringDoctor        *ReferralSimpleDoctor     `json:"referring_doctor,omitempty"`
	ReferredToDepartment   *ReferralSimpleDepartment `json:"referred_to_department,omitempty"`
	ReferredToDoctor       *ReferralSimpleDoctor     `json:"referred_to_doctor,omitempty"`
}

type ReferralListResponse struct {
	Data []ReferralResponse     `json:"data"`
	Meta ReferralPaginationMeta `json:"meta"`
}

type ReferralDeletedListResponse struct {
	Data []ReferralResponse     `json:"data"`
	Meta ReferralPaginationMeta `json:"meta"`
}

type ReferralMyListResponse struct {
	Data []ReferralResponse `json:"data"`
}

// ─── Requests ─────────────────────────────────────────────────────────────

type ReferralPaginationQuery struct {
	Page              int    `form:"page,default=1" binding:"omitempty,min=1"`
	PageSize          int    `form:"page_size,default=10" binding:"omitempty,min=1,max=100"`
	SortBy            string `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at updated_at id referral_date priority status"`
	SortDir           string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
	Search            string `form:"search" binding:"omitempty"`
	Status            string `form:"status" binding:"omitempty,oneof=pending accepted rejected completed cancelled"`
	Priority          string `form:"priority" binding:"omitempty,oneof=routine urgent emergency"`
	ReferralType      string `form:"referral_type" binding:"omitempty,oneof=internal external"`
	ReferringDoctorID *uint  `form:"referring_doctor_id" binding:"omitempty"`
	PatientID         *uint  `form:"-"` // set by handler logic if needed
	DoctorID          *uint  `form:"-"` // set by handler logic if needed
}

type ReferralPaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type CreateReferralRequest struct {
	ReferralNumber         string `json:"referral_number" binding:"omitempty,max=50"`
	PatientID              uint   `json:"patient_id" binding:"required"`
	MedicalRecordID        uint   `json:"medical_record_id" binding:"required"`
	ReferringDoctorID      uint   `json:"referring_doctor_id" binding:"required"`
	ReferralDate           string `json:"referral_date" binding:"required,datetime=2006-01-02"`
	ReferralType           string `json:"referral_type" binding:"required,oneof=internal external"`
	ReferredToDepartmentID *uint  `json:"referred_to_department_id" binding:"omitempty"`
	ReferredToDoctorID     *uint  `json:"referred_to_doctor_id" binding:"omitempty"`
	ReferredToFacility     string `json:"referred_to_facility" binding:"omitempty,max=255"`
	Reason                 string `json:"reason" binding:"required"`
	Diagnosis              string `json:"diagnosis" binding:"omitempty"`
	Priority               string `json:"priority" binding:"required,oneof=routine urgent emergency"`
	Notes                  string `json:"notes" binding:"omitempty"`
}

type UpdateReferralRequest struct {
	ReferralNumber         *string `json:"referral_number" binding:"omitempty,max=50"`
	PatientID              *uint   `json:"patient_id" binding:"omitempty"`
	MedicalRecordID        *uint   `json:"medical_record_id" binding:"omitempty"`
	ReferringDoctorID      *uint   `json:"referring_doctor_id" binding:"omitempty"`
	ReferralDate           *string `json:"referral_date" binding:"omitempty,datetime=2006-01-02"`
	ReferralType           *string `json:"referral_type" binding:"omitempty,oneof=internal external"`
	ReferredToDepartmentID *uint   `json:"referred_to_department_id" binding:"omitempty"`
	ReferredToDoctorID     *uint   `json:"referred_to_doctor_id" binding:"omitempty"`
	ReferredToFacility     *string `json:"referred_to_facility" binding:"omitempty,max=255"`
	Reason                 *string `json:"reason" binding:"omitempty"`
	Diagnosis              *string `json:"diagnosis" binding:"omitempty"`
	Priority               *string `json:"priority" binding:"omitempty,oneof=routine urgent emergency"`
	Notes                  *string `json:"notes" binding:"omitempty"`
}

type AcceptReferralRequest struct {
	Notes string `json:"notes" binding:"omitempty"`
}

type RejectReferralRequest struct {
	RejectionReason string `json:"rejection_reason" binding:"required"`
}

type CompleteReferralRequest struct {
	Notes string `json:"notes" binding:"omitempty"`
}

type CancelReferralRequest struct {
	CancellationReason string `json:"cancellation_reason" binding:"omitempty"`
}
