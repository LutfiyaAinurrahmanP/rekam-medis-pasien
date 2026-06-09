// Package events mendefinisikan semua tipe event yang diproduksi oleh aplikasi
// dan di-publish ke Apache Kafka. Setiap event memiliki metadata standar
// (EventID, EventType, Timestamp) ditambah payload yang spesifik terhadap domain.
package events

import (
	"time"

	"github.com/google/uuid"
)

// ─── Base Event ───────────────────────────────────────────────────────────────

// BaseEvent adalah metadata standar yang ada di setiap event.
type BaseEvent struct {
	EventID   string    `json:"event_id"`   // UUID unik per event
	EventType string    `json:"event_type"` // contoh: "patient.created"
	Timestamp time.Time `json:"timestamp"`  // waktu event terjadi (UTC)
	Version   string    `json:"version"`    // versi schema event
}

func newBase(eventType string) BaseEvent {
	return BaseEvent{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
	}
}

// ─── User Events ─────────────────────────────────────────────────────────────

type UserPayload struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	Action    string `json:"action,omitempty"`
}

type UserRegisteredEvent struct {
	BaseEvent
	Payload UserPayload `json:"payload"`
}

type UserLoginEvent struct {
	BaseEvent
	Payload struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"payload"`
}

type UserCreatedEvent struct {
	BaseEvent
	Payload UserPayload `json:"payload"`
}

type UserUpdatedEvent struct {
	BaseEvent
	Payload UserPayload `json:"payload"`
}

type UserDeletedEvent struct {
	BaseEvent
	Payload UserPayload `json:"payload"`
}

type UserRestoredEvent struct {
	BaseEvent
	Payload UserPayload `json:"payload"`
}

type UserPasswordResetRequestedEvent struct {
	BaseEvent
	Payload struct {
		UserID    uint   `json:"user_id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		ResetCode string `json:"reset_code"`
		ExpiresIn int    `json:"expires_in"`
	} `json:"payload"`
}

// Constructors ────────────────────────────────────────────────────────────────

func NewUserRegisteredEvent(id uint, username, email, phone, role string) UserRegisteredEvent {
	return UserRegisteredEvent{
		BaseEvent: newBase("user.registered"),
		Payload: UserPayload{
			ID: id, Username: username, Email: email,
			Phone: phone, Role: role, IsActive: true,
		},
	}
}

func NewUserLoginEvent(userID uint, username, email, role string) UserLoginEvent {
	e := UserLoginEvent{BaseEvent: newBase("user.login")}
	e.Payload.UserID = userID
	e.Payload.Username = username
	e.Payload.Email = email
	e.Payload.Role = role
	return e
}

func NewUserCreatedEvent(id uint, username, email, phone, role string, isActive bool) UserCreatedEvent {
	return UserCreatedEvent{
		BaseEvent: newBase("user.created"),
		Payload: UserPayload{
			ID: id, Username: username, Email: email,
			Phone: phone, Role: role, IsActive: isActive,
		},
	}
}

func NewUserUpdatedEvent(id uint, username, email, phone, role string, isActive bool, action string) UserUpdatedEvent {
	return UserUpdatedEvent{
		BaseEvent: newBase("user.updated"),
		Payload: UserPayload{
			ID: id, Username: username, Email: email,
			Phone: phone, Role: role, IsActive: isActive,
			Action: action,
		},
	}
}

func NewUserDeletedEvent(userID uint, username, action string) UserDeletedEvent {
	return UserDeletedEvent{
		BaseEvent: newBase("user.deleted"),
		Payload: UserPayload{
			ID:       userID,
			Username: username,
			Action:   action,
		},
	}
}

func NewUserRestoredEvent(userID uint, username string) UserRestoredEvent {
	return UserRestoredEvent{
		BaseEvent: newBase("user.restored"),
		Payload: UserPayload{
			ID:       userID,
			Username: username,
		},
	}
}

func NewUserPasswordResetRequestedEvent(userID uint, username, email, resetCode string, expiresIn int) UserPasswordResetRequestedEvent {
	e := UserPasswordResetRequestedEvent{BaseEvent: newBase("user.password_reset_requested")}
	e.Payload.UserID = userID
	e.Payload.Username = username
	e.Payload.Email = email
	e.Payload.ResetCode = resetCode
	e.Payload.ExpiresIn = expiresIn
	return e
}

// ─── Patient Events ───────────────────────────────────────────────────────────

type PatientPayload struct {
	ID                    uint   `json:"id"`
	UserID                *uint  `json:"user_id,omitempty"`
	PatientCode           string `json:"patient_code"`
	FullName              string `json:"full_name"`
	DateOfBirth           string `json:"date_of_birth"`
	Gender                string `json:"gender"`
	BloodType             string `json:"blood_type"`
	Phone                 string `json:"phone"`
	Email                 string `json:"email"`
	Address               string `json:"address"`
	EmergencyContactName  string `json:"emergency_contact_name"`
	EmergencyContactPhone string `json:"emergency_contact_phone"`
	InsuranceNumber       string `json:"insurance_number"`
	InsuranceProvider     string `json:"insurance_provider"`
	Allergies             string `json:"allergies"`
	CreatedAt             string `json:"created_at,omitempty"`
	UpdatedAt             string `json:"updated_at,omitempty"`
	DeletedAt             string `json:"deleted_at,omitempty"`
	Action                string `json:"action,omitempty"`
}

type PatientCreatedEvent struct {
	BaseEvent
	Payload PatientPayload `json:"payload"`
}

type PatientUpdatedEvent struct {
	BaseEvent
	Payload PatientPayload `json:"payload"`
}

type PatientDeletedEvent struct {
	BaseEvent
	Payload PatientPayload `json:"payload"`
}

type PatientRestoredEvent struct {
	BaseEvent
	Payload PatientPayload `json:"payload"`
}

// Constructors ────────────────────────────────────────────────────────────────

func NewPatientCreatedEvent(id uint, code, fullName, dob, gender, blood, phone, email, insurance string) PatientCreatedEvent {
	return PatientCreatedEvent{
		BaseEvent: newBase("patient.created"),
		Payload: PatientPayload{
			ID: id, PatientCode: code, FullName: fullName, DateOfBirth: dob,
			Gender: gender, BloodType: blood, Phone: phone, Email: email,
			InsuranceProvider: insurance,
		},
	}
}

func NewPatientUpdatedEvent(id uint, code, fullName, dob, gender, blood, phone, email, insurance, action string) PatientUpdatedEvent {
	return PatientUpdatedEvent{
		BaseEvent: newBase("patient.updated"),
		Payload: PatientPayload{
			ID: id, PatientCode: code, FullName: fullName, DateOfBirth: dob,
			Gender: gender, BloodType: blood, Phone: phone, Email: email,
			InsuranceProvider: insurance, Action: action,
		},
	}
}

func NewPatientDeletedEvent(id uint, code, fullName, action string) PatientDeletedEvent {
	return PatientDeletedEvent{
		BaseEvent: newBase("patient.deleted"),
		Payload: PatientPayload{
			ID: id, PatientCode: code, FullName: fullName, Action: action,
		},
	}
}

func NewPatientRestoredEvent(id uint, code, fullName string) PatientRestoredEvent {
	return PatientRestoredEvent{
		BaseEvent: newBase("patient.restored"),
		Payload: PatientPayload{
			ID: id, PatientCode: code, FullName: fullName,
		},
	}
}

// ─── Doctor Specialization Events ─────────────────────────────────────────────

type DoctorSpecializationPayload struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	DeletedAt   string `json:"deleted_at,omitempty"`
	Action      string `json:"action,omitempty"`
}

type DoctorSpecializationCreatedEvent struct {
	BaseEvent
	Payload DoctorSpecializationPayload `json:"payload"`
}

type DoctorSpecializationUpdatedEvent struct {
	BaseEvent
	Payload DoctorSpecializationPayload `json:"payload"`
}

type DoctorSpecializationDeletedEvent struct {
	BaseEvent
	Payload DoctorSpecializationPayload `json:"payload"`
}

type DoctorSpecializationRestoredEvent struct {
	BaseEvent
	Payload DoctorSpecializationPayload `json:"payload"`
}

// Constructors ─────────────────────────────────────────────────────────────────

func NewDoctorSpecializationCreatedEvent(id uint, name, code, desc string, isActive bool) DoctorSpecializationCreatedEvent {
	return DoctorSpecializationCreatedEvent{
		BaseEvent: newBase("doctorspecialization.created"),
		Payload: DoctorSpecializationPayload{
			ID:          id,
			Name:        name,
			Code:        code,
			Description: desc,
			IsActive:    isActive,
		},
	}
}

func NewDoctorSpecializationUpdatedEvent(id uint, name, code, desc string, isActive bool, action string) DoctorSpecializationUpdatedEvent {
	return DoctorSpecializationUpdatedEvent{
		BaseEvent: newBase("doctor_specialization.updated"),
		Payload: DoctorSpecializationPayload{
			ID:          id,
			Name:        name,
			Code:        code,
			Description: desc,
			IsActive:    isActive,
			Action:      action,
		},
	}
}

func NewDoctorSpecializationDeletedEvent(id uint, name, code, action string) DoctorSpecializationDeletedEvent {
	return DoctorSpecializationDeletedEvent{
		BaseEvent: newBase("doctor_specialization.deleted"),
		Payload: DoctorSpecializationPayload{
			ID:     id,
			Name:   name,
			Code:   code,
			Action: action,
		},
	}
}

func NewDoctorSpecializationRestoredEvent(id uint, name, code string) DoctorSpecializationRestoredEvent {
	return DoctorSpecializationRestoredEvent{
		BaseEvent: newBase("doctor_specialization.restored"),
		Payload: DoctorSpecializationPayload{
			ID:   id,
			Name: name,
			Code: code,
		},
	}
}

// ─── Doctor Events ────────────────────────────────────────────────────────────

type DoctorPayload struct {
	ID               uint   `json:"id"`
	UserID           *uint  `json:"user_id,omitempty"`
	EmployeeID       string `json:"employee_id"`
	FullName         string `json:"full_name"`
	SpecializationID uint   `json:"specialization_id"`
	LicenseNumber    string `json:"license_number"`
	Phone            string `json:"phone"`
	Email            string `json:"email"`
	DepartmentID     *uint  `json:"department_id,omitempty"`
	IsActive         bool   `json:"is_active"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
	DeletedAt        string `json:"deleted_at,omitempty"`
	Action           string `json:"action,omitempty"`
}

type DoctorCreatedEvent struct {
	BaseEvent
	Payload DoctorPayload `json:"payload"`
}

type DoctorUpdatedEvent struct {
	BaseEvent
	Payload DoctorPayload `json:"payload"`
}

type DoctorDeletedEvent struct {
	BaseEvent
	Payload DoctorPayload `json:"payload"`
}

type DoctorRestoredEvent struct {
	BaseEvent
	Payload DoctorPayload `json:"payload"`
}

// Constructors ────────────────────────────────────────────────────────────────

func NewDoctorCreatedEvent(id uint, fullName string, specializationID uint, phone, email string, deptID *uint, isActive bool) DoctorCreatedEvent {
	return DoctorCreatedEvent{
		BaseEvent: newBase("doctor.created"),
		Payload: DoctorPayload{
			ID: id, FullName: fullName, SpecializationID: specializationID,
			Phone: phone, Email: email, DepartmentID: deptID, IsActive: isActive,
		},
	}
}

func NewDoctorUpdatedEvent(id uint, fullName string, specializationID uint, phone, email string, deptID *uint, isActive bool, action string) DoctorUpdatedEvent {
	return DoctorUpdatedEvent{
		BaseEvent: newBase("doctor.updated"),
		Payload: DoctorPayload{
			ID: id, FullName: fullName, SpecializationID: specializationID,
			Phone: phone, Email: email, DepartmentID: deptID, IsActive: isActive,
			Action: action,
		},
	}
}

func NewDoctorDeletedEvent(id uint, fullName, action string) DoctorDeletedEvent {
	return DoctorDeletedEvent{
		BaseEvent: newBase("doctor.deleted"),
		Payload: DoctorPayload{
			ID:       id,
			FullName: fullName,
			Action:   action,
		},
	}
}

func NewDoctorRestoredEvent(id uint, fullName string) DoctorRestoredEvent {
	return DoctorRestoredEvent{
		BaseEvent: newBase("doctor.restored"),
		Payload: DoctorPayload{
			ID:       id,
			FullName: fullName,
		},
	}
}

// ─── Department Events ────────────────────────────────────────────────────────

type DepartmentPayload struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	Description   string `json:"description"`
	FloorLocation string `json:"floor_location"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	DeletedAt     string `json:"deleted_at,omitempty"`
	Action        string `json:"action,omitempty"`
}

type DepartmentCreatedEvent struct {
	BaseEvent
	Payload DepartmentPayload `json:"payload"`
}

type DepartmentUpdatedEvent struct {
	BaseEvent
	Payload DepartmentPayload `json:"payload"`
}

type DepartmentDeletedEvent struct {
	BaseEvent
	Payload DepartmentPayload `json:"payload"`
}

type DepartmentRestoredEvent struct {
	BaseEvent
	Payload DepartmentPayload `json:"payload"`
}

// Constructors ────────────────────────────────────────────────────────────────

func NewDepartmentCreatedEvent(id uint, name, description, code string) DepartmentCreatedEvent {
	return DepartmentCreatedEvent{
		BaseEvent: newBase("department.created"),
		Payload:   DepartmentPayload{ID: id, Name: name, Description: description, Code: code},
	}
}

func NewDepartmentUpdatedEvent(id uint, name, description, code string) DepartmentUpdatedEvent {
	return DepartmentUpdatedEvent{
		BaseEvent: newBase("department.updated"),
		Payload:   DepartmentPayload{ID: id, Name: name, Description: description, Code: code},
	}
}

func NewDepartmentDeletedEvent(id uint, name, action string) DepartmentDeletedEvent {
	return DepartmentDeletedEvent{
		BaseEvent: newBase("department.deleted"),
		Payload: DepartmentPayload{
			ID:     id,
			Name:   name,
			Action: action,
		},
	}
}

func NewDepartmentRestoredEvent(id uint, name string) DepartmentRestoredEvent {
	return DepartmentRestoredEvent{
		BaseEvent: newBase("department.restored"),
		Payload: DepartmentPayload{
			ID:   id,
			Name: name,
		},
	}
}

// ─── Room Type Events ──────────────────────────────────────────────────────────────

type RoomTypePayload struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	DeletedAt   string `json:"deleted_at,omitempty"`
	Action      string `json:"action,omitempty"`
}

type RoomTypeCreatedEvent struct {
	BaseEvent
	Payload RoomTypePayload `json:"payload"`
}

type RoomTypeUpdatedEvent struct {
	BaseEvent
	Payload RoomTypePayload `json:"payload"`
}

type RoomTypeDeletedEvent struct {
	BaseEvent
	Payload RoomTypePayload `json:"payload"`
}

type RoomTypeRestoredEvent struct {
	BaseEvent
	Payload RoomTypePayload `json:"payload"`
}

// Constructors ────────────────────────────────────────────────────────────────

func NewRoomTypeCreatedEvent(id uint, name, code, description string, isActive bool) RoomTypeCreatedEvent {
	return RoomTypeCreatedEvent{
		BaseEvent: newBase("room_type.created"),
		Payload: RoomTypePayload{
			ID:          id,
			Name:        name,
			Code:        code,
			Description: description,
			IsActive:    isActive,
		},
	}
}

func NewRoomTypeUpdatedEvent(id uint, name, code, description string, isActive bool, action string) RoomTypeUpdatedEvent {
	return RoomTypeUpdatedEvent{
		BaseEvent: newBase("room_type.updated"),
		Payload: RoomTypePayload{
			ID:          id,
			Name:        name,
			Code:        code,
			Description: description,
			IsActive:    isActive,
			Action:      action,
		},
	}
}

func NewRoomTypeDeletedEvent(id uint, name, action string) RoomTypeDeletedEvent {
	return RoomTypeDeletedEvent{
		BaseEvent: newBase("room_type.deleted"),
		Payload: RoomTypePayload{
			ID:     id,
			Name:   name,
			Action: action,
		},
	}
}

func NewRoomTypeRestoredEvent(id uint, name string) RoomTypeRestoredEvent {
	return RoomTypeRestoredEvent{
		BaseEvent: newBase("room_type.restored"),
		Payload: RoomTypePayload{
			ID:   id,
			Name: name,
		},
	}
}

// ─── Room Events ──────────────────────────────────────────────────────────────

type RoomPayload struct {
	ID            uint    `json:"id"`
	RoomNumber    string  `json:"room_number"`
	RoomTypeID    *uint   `json:"room_type_id,omitempty"`
	DepartmentID  *uint   `json:"department_id,omitempty"`
	BedCapacity   int     `json:"bed_capacity"`
	AvailableBeds int     `json:"available_beds"`
	PricePerDay   float64 `json:"price_per_day"`
	IsActive      bool    `json:"is_active"`
	CreatedAt     string  `json:"created_at,omitempty"`
	UpdatedAt     string  `json:"updated_at,omitempty"`
	DeletedAt     string  `json:"deleted_at,omitempty"`
	Action        string  `json:"action,omitempty"`
}

type RoomCreatedEvent struct {
	BaseEvent
	Payload RoomPayload `json:"payload"`
}

type RoomUpdatedEvent struct {
	BaseEvent
	Payload RoomPayload `json:"payload"`
}

type RoomDeletedEvent struct {
	BaseEvent
	Payload RoomPayload `json:"payload"`
}

type RoomRestoredEvent struct {
	BaseEvent
	Payload RoomPayload `json:"payload"`
}

// Constructors ────────────────────────────────────────────────────────────────

func NewRoomCreatedEvent(id uint, roomNumber string, roomTypeID *uint, deptID *uint, total, available int, price float64, isActive bool) RoomCreatedEvent {
	return RoomCreatedEvent{
		BaseEvent: newBase("room.created"),
		Payload: RoomPayload{
			ID: id, RoomNumber: roomNumber, RoomTypeID: roomTypeID, DepartmentID: deptID,
			BedCapacity: total, AvailableBeds: available, PricePerDay: price, IsActive: isActive,
		},
	}
}

func NewRoomUpdatedEvent(id uint, roomNumber string, roomTypeID *uint, deptID *uint, total, available int, price float64, isActive bool, action string) RoomUpdatedEvent {
	return RoomUpdatedEvent{
		BaseEvent: newBase("room.updated"),
		Payload: RoomPayload{
			ID: id, RoomNumber: roomNumber, RoomTypeID: roomTypeID, DepartmentID: deptID,
			BedCapacity: total, AvailableBeds: available, PricePerDay: price, IsActive: isActive,
			Action: action,
		},
	}
}

func NewRoomDeletedEvent(id uint, roomNumber, action string) RoomDeletedEvent {
	return RoomDeletedEvent{
		BaseEvent: newBase("room.deleted"),
		Payload: RoomPayload{
			ID:         id,
			RoomNumber: roomNumber,
			Action:     action,
		},
	}
}

func NewRoomRestoredEvent(id uint, roomNumber string) RoomRestoredEvent {
	return RoomRestoredEvent{
		BaseEvent: newBase("room.restored"),
		Payload: RoomPayload{
			ID:         id,
			RoomNumber: roomNumber,
		},
	}
}

// ─── Type Test Category Events ─────────────────────────────────────────────

type TypeTestCategoryPayload struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	DeletedAt   string `json:"deleted_at,omitempty"`
	Action      string `json:"action,omitempty"`
}

type TypeTestCategoryCreatedEvent struct {
	BaseEvent
	Payload TypeTestCategoryPayload `json:"payload"`
}

type TypeTestCategoryUpdatedEvent struct {
	BaseEvent
	Payload TypeTestCategoryPayload `json:"payload"`
}

type TypeTestCategoryDeletedEvent struct {
	BaseEvent
	Payload TypeTestCategoryPayload `json:"payload"`
}

type TypeTestCategoryRestoredEvent struct {
	BaseEvent
	Payload TypeTestCategoryPayload `json:"payload"`
}

// Constructors ─────────────────────────────────────────────────────────────────

func NewTypeTestCategoryCreatedEvent(id uint, name, code, desc string, isActive bool) TypeTestCategoryCreatedEvent {
	return TypeTestCategoryCreatedEvent{
		BaseEvent: newBase("type_test_category.created"),
		Payload: TypeTestCategoryPayload{
			ID:          id,
			Name:        name,
			Code:        code,
			Description: desc,
			IsActive:    isActive,
		},
	}
}

func NewTypeTestCategoryUpdatedEvent(id uint, name, code, desc string, isActive bool, action string) TypeTestCategoryUpdatedEvent {
	return TypeTestCategoryUpdatedEvent{
		BaseEvent: newBase("type_test_category.updated"),
		Payload: TypeTestCategoryPayload{
			ID:          id,
			Name:        name,
			Code:        code,
			Description: desc,
			IsActive:    isActive,
			Action:      action,
		},
	}
}

func NewTypeTestCategoryDeletedEvent(id uint, name, code, action string) TypeTestCategoryDeletedEvent {
	return TypeTestCategoryDeletedEvent{
		BaseEvent: newBase("type_test_category.deleted"),
		Payload: TypeTestCategoryPayload{
			ID:     id,
			Name:   name,
			Code:   code,
			Action: action,
		},
	}
}

func NewTypeTestCategoryRestoredEvent(id uint, name, code string) TypeTestCategoryRestoredEvent {
	return TypeTestCategoryRestoredEvent{
		BaseEvent: newBase("type_test_category.restored"),
		Payload: TypeTestCategoryPayload{
			ID:   id,
			Name: name,
			Code: code,
		},
	}
}

// ─── TypeTest Events ──────────────────────────────────────────────────────────

type TypeTestPayload struct {
	ID                 uint    `json:"id"`
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	TypeTestCategoryID uint    `json:"type_test_category_id"`
	Description        string  `json:"description"`
	Price              float64 `json:"price"`
	IsActive           bool    `json:"is_active"`
	CreatedAt          string  `json:"created_at,omitempty"`
	UpdatedAt          string  `json:"updated_at,omitempty"`
	DeletedAt          string  `json:"deleted_at,omitempty"`
	Action             string  `json:"action,omitempty"`
}

type TypeTestCreatedEvent struct {
	BaseEvent
	Payload TypeTestPayload `json:"payload"`
}

type TypeTestUpdatedEvent struct {
	BaseEvent
	Payload TypeTestPayload `json:"payload"`
}

type TypeTestDeletedEvent struct {
	BaseEvent
	Payload TypeTestPayload `json:"payload"`
}

type TypeTestRestoredEvent struct {
	BaseEvent
	Payload TypeTestPayload `json:"payload"`
}

// Constructors ────────────────────────────────────────────────────────────────

func NewTypeTestCreatedEvent(id uint, code, name string, typeTestCategoryID uint, description string, price float64, isActive bool) TypeTestCreatedEvent {
	return TypeTestCreatedEvent{
		BaseEvent: newBase("type_test.created"),
		Payload:   TypeTestPayload{ID: id, Code: code, Name: name, TypeTestCategoryID: typeTestCategoryID, Description: description, Price: price, IsActive: isActive},
	}
}

func NewTypeTestUpdatedEvent(id uint, code, name string, typeTestCategoryID uint, description string, price float64, isActive bool, action string) TypeTestUpdatedEvent {
	return TypeTestUpdatedEvent{
		BaseEvent: newBase("type_test.updated"),
		Payload: TypeTestPayload{
			ID: id, Code: code, Name: name, TypeTestCategoryID: typeTestCategoryID, Description: description, Price: price, IsActive: isActive,
			Action: action,
		},
	}
}

func NewTypeTestDeletedEvent(id uint, code, name, action string) TypeTestDeletedEvent {
	return TypeTestDeletedEvent{
		BaseEvent: newBase("type_test.deleted"),
		Payload: TypeTestPayload{
			ID:     id,
			Code:   code,
			Name:   name,
			Action: action,
		},
	}
}

func NewTypeTestRestoredEvent(id uint, code, name string) TypeTestRestoredEvent {
	return TypeTestRestoredEvent{
		BaseEvent: newBase("type_test.restored"),
		Payload: TypeTestPayload{
			ID:   id,
			Code: code,
			Name: name,
		},
	}
}

// ─── Medicine Types Events ─────────────────────────────────────────────

type MedicineTypePayload struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	DeletedAt   string `json:"deleted_at,omitempty"`
	Action      string `json:"action,omitempty"`
}

type MedicineTypeCreatedEvent struct {
	BaseEvent
	Payload MedicineTypePayload `json:"payload"`
}

type MedicineTypeUpdatedEvent struct {
	BaseEvent
	Payload MedicineTypePayload `json:"payload"`
}

type MedicineTypeDeletedEvent struct {
	BaseEvent
	Payload MedicineTypePayload `json:"payload"`
}

type MedicineTypeRestoredEvent struct {
	BaseEvent
	Payload MedicineTypePayload `json:"payload"`
}

// Constructors ─────────────────────────────────────────────────────────────────

func NewMedicineTypeCreatedEvent(id uint, name, code, desc string, isActive bool) MedicineTypeCreatedEvent {
	return MedicineTypeCreatedEvent{
		BaseEvent: newBase("medicine_type.created"),
		Payload: MedicineTypePayload{
			ID:          id,
			Name:        name,
			Code:        code,
			Description: desc,
			IsActive:    isActive,
		},
	}
}

func NewMedicineTypeUpdatedEvent(id uint, name, code, desc string, isActive bool, action string) MedicineTypeUpdatedEvent {
	return MedicineTypeUpdatedEvent{
		BaseEvent: newBase("medicine_type.updated"),
		Payload: MedicineTypePayload{
			ID:          id,
			Name:        name,
			Code:        code,
			Description: desc,
			IsActive:    isActive,
			Action:      action,
		},
	}
}

func NewMedicineTypeDeletedEvent(id uint, name, code, action string) MedicineTypeDeletedEvent {
	return MedicineTypeDeletedEvent{
		BaseEvent: newBase("medicine_type.deleted"),
		Payload: MedicineTypePayload{
			ID:     id,
			Name:   name,
			Code:   code,
			Action: action,
		},
	}
}

func NewMedicineTypeRestoredEvent(id uint, name, code string) MedicineTypeRestoredEvent {
	return MedicineTypeRestoredEvent{
		BaseEvent: newBase("medicine_type.restored"),
		Payload: MedicineTypePayload{
			ID:   id,
			Name: name,
			Code: code,
		},
	}
}

// ─── Medicine Events ────────────────────────────────────────────────────────────

type MedicinePayload struct {
	ID             uint    `json:"id"`
	Name           string  `json:"name,omitempty"`
	GenericName    string  `json:"generic_name,omitempty"`
	BrandName      string  `json:"brand_name,omitempty"`
	MedicineTypeID uint    `json:"medicine_type_id,omitempty"`
	Strength       string  `json:"strength,omitempty"`
	Manufacturer   string  `json:"manufacturer,omitempty"`
	Unit           string  `json:"unit,omitempty"`
	StockQuantity  int     `json:"stock_quantity,omitempty"`
	Price          float64 `json:"price,omitempty"`
	IsActive       bool    `json:"is_active,omitempty"`
	CreatedAt      string  `json:"created_at,omitempty"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
	DeletedAt      string  `json:"deleted_at,omitempty"`
	Action         string  `json:"action,omitempty"`
	DeltaQuantity  int     `json:"delta_quantity,omitempty"` // For add/reduce stock
}

type MedicineCreatedEvent struct {
	BaseEvent
	Payload MedicinePayload `json:"payload"`
}

type MedicineUpdatedEvent struct {
	BaseEvent
	Payload MedicinePayload `json:"payload"`
}

type MedicineDeletedEvent struct {
	BaseEvent
	Payload MedicinePayload `json:"payload"`
}

type MedicineRestoredEvent struct {
	BaseEvent
	Payload MedicinePayload `json:"payload"`
}

type MedicineActivatedEvent struct {
	BaseEvent
	Payload MedicinePayload `json:"payload"`
}

type MedicineDeactivatedEvent struct {
	BaseEvent
	Payload MedicinePayload `json:"payload"`
}

type MedicineStockAddedEvent struct {
	BaseEvent
	Payload MedicinePayload `json:"payload"`
}

type MedicineStockReducedEvent struct {
	BaseEvent
	Payload MedicinePayload `json:"payload"`
}

// Constructors ─────────────────────────────────────────────────────────────────

func NewMedicineCreatedEvent(id uint, name, genericName, brandName string, medicineTypeID uint, strength, manufacturer, unit string, stock int, price float64, isActive bool) MedicineCreatedEvent {
	return MedicineCreatedEvent{
		BaseEvent: newBase("medicine.created"),
		Payload: MedicinePayload{
			ID: id, Name: name, GenericName: genericName, BrandName: brandName,
			MedicineTypeID: medicineTypeID, Strength: strength, Manufacturer: manufacturer,
			Unit: unit, StockQuantity: stock, Price: price, IsActive: isActive,
		},
	}
}

func NewMedicineUpdatedEvent(id uint, name, genericName, brandName string, medicineTypeID uint, strength, manufacturer, unit string, stock int, price float64, isActive bool, action string) MedicineUpdatedEvent {
	return MedicineUpdatedEvent{
		BaseEvent: newBase("medicine.updated"),
		Payload: MedicinePayload{
			ID: id, Name: name, GenericName: genericName, BrandName: brandName,
			MedicineTypeID: medicineTypeID, Strength: strength, Manufacturer: manufacturer,
			Unit: unit, StockQuantity: stock, Price: price, IsActive: isActive, Action: action,
		},
	}
}

func NewMedicineDeletedEvent(id uint, name, action string) MedicineDeletedEvent {
	return MedicineDeletedEvent{
		BaseEvent: newBase("medicine.deleted"),
		Payload: MedicinePayload{
			ID:     id,
			Name:   name,
			Action: action,
		},
	}
}

func NewMedicineRestoredEvent(id uint, name string) MedicineRestoredEvent {
	return MedicineRestoredEvent{
		BaseEvent: newBase("medicine.restored"),
		Payload: MedicinePayload{
			ID:   id,
			Name: name,
		},
	}
}

func NewMedicineActivatedEvent(id uint, name string) MedicineActivatedEvent {
	return MedicineActivatedEvent{
		BaseEvent: newBase("medicine.activated"),
		Payload: MedicinePayload{
			ID:   id,
			Name: name,
		},
	}
}

func NewMedicineDeactivatedEvent(id uint, name string) MedicineDeactivatedEvent {
	return MedicineDeactivatedEvent{
		BaseEvent: newBase("medicine.deactivated"),
		Payload: MedicinePayload{
			ID:   id,
			Name: name,
		},
	}
}

func NewMedicineStockAddedEvent(id uint, name string, newStock, addedQuantity int) MedicineStockAddedEvent {
	return MedicineStockAddedEvent{
		BaseEvent: newBase("medicine.stock_added"),
		Payload: MedicinePayload{
			ID:            id,
			Name:          name,
			StockQuantity: newStock,
			DeltaQuantity: addedQuantity,
		},
	}
}

func NewMedicineStockReducedEvent(id uint, name string, newStock, reducedQuantity int) MedicineStockReducedEvent {
	return MedicineStockReducedEvent{
		BaseEvent: newBase("medicine.stock_reduced"),
		Payload: MedicinePayload{
			ID:            id,
			Name:          name,
			StockQuantity: newStock,
			DeltaQuantity: reducedQuantity,
		},
	}
}

// ─── Appointment Events ────────────────────────────────────────────────────────

type AppointmentPayload struct {
	ID              uint   `json:"id"`
	PatientID       uint   `json:"patient_id,omitempty"`
	DoctorID        uint   `json:"doctor_id,omitempty"`
	AppointmentDate string `json:"appointment_date,omitempty"`
	AppointmentTime string `json:"appointment_time,omitempty"`
	DurationMinutes int    `json:"duration_minutes,omitempty"`
	Status          string `json:"status,omitempty"`
	Reason          string `json:"reason,omitempty"`
	Notes           string `json:"notes,omitempty"`
	Action          string `json:"action,omitempty"`
}

type AppointmentCreatedEvent struct {
	BaseEvent
	Payload AppointmentPayload `json:"payload"`
}

type AppointmentUpdatedEvent struct {
	BaseEvent
	Payload AppointmentPayload `json:"payload"`
}

type AppointmentDeletedEvent struct {
	BaseEvent
	Payload AppointmentPayload `json:"payload"`
}

type AppointmentRestoredEvent struct {
	BaseEvent
	Payload AppointmentPayload `json:"payload"`
}

type AppointmentConfirmedEvent struct {
	BaseEvent
	Payload AppointmentPayload `json:"payload"`
}

type AppointmentStartedEvent struct {
	BaseEvent
	Payload AppointmentPayload `json:"payload"`
}

type AppointmentCompletedEvent struct {
	BaseEvent
	Payload AppointmentPayload `json:"payload"`
}

type AppointmentCancelledEvent struct {
	BaseEvent
	Payload AppointmentPayload `json:"payload"`
}

type AppointmentRescheduledEvent struct {
	BaseEvent
	Payload AppointmentPayload `json:"payload"`
}

type AppointmentNoShowEvent struct {
	BaseEvent
	Payload AppointmentPayload `json:"payload"`
}

// Constructors

func NewAppointmentCreatedEvent(id, patientID, doctorID uint, date, timeStr string, duration int, status string) AppointmentCreatedEvent {
	return AppointmentCreatedEvent{
		BaseEvent: newBase("appointment.created"),
		Payload: AppointmentPayload{
			ID:              id,
			PatientID:       patientID,
			DoctorID:        doctorID,
			AppointmentDate: date,
			AppointmentTime: timeStr,
			DurationMinutes: duration,
			Status:          status,
		},
	}
}

func NewAppointmentUpdatedEvent(id, patientID, doctorID uint, date, timeStr string, duration int, status, action string) AppointmentUpdatedEvent {
	return AppointmentUpdatedEvent{
		BaseEvent: newBase("appointment.updated"),
		Payload: AppointmentPayload{
			ID:              id,
			PatientID:       patientID,
			DoctorID:        doctorID,
			AppointmentDate: date,
			AppointmentTime: timeStr,
			DurationMinutes: duration,
			Status:          status,
			Action:          action,
		},
	}
}

func NewAppointmentDeletedEvent(id uint, action string) AppointmentDeletedEvent {
	return AppointmentDeletedEvent{
		BaseEvent: newBase("appointment.deleted"),
		Payload: AppointmentPayload{
			ID:     id,
			Action: action,
		},
	}
}

func NewAppointmentRestoredEvent(id uint) AppointmentRestoredEvent {
	return AppointmentRestoredEvent{
		BaseEvent: newBase("appointment.restored"),
		Payload: AppointmentPayload{
			ID: id,
		},
	}
}

func NewAppointmentConfirmedEvent(id uint) AppointmentConfirmedEvent {
	return AppointmentConfirmedEvent{
		BaseEvent: newBase("appointment.confirmed"),
		Payload: AppointmentPayload{
			ID: id,
		},
	}
}

func NewAppointmentStartedEvent(id uint) AppointmentStartedEvent {
	return AppointmentStartedEvent{
		BaseEvent: newBase("appointment.started"),
		Payload: AppointmentPayload{
			ID: id,
		},
	}
}

func NewAppointmentCompletedEvent(id uint) AppointmentCompletedEvent {
	return AppointmentCompletedEvent{
		BaseEvent: newBase("appointment.completed"),
		Payload: AppointmentPayload{
			ID: id,
		},
	}
}

func NewAppointmentCancelledEvent(id uint, reason string) AppointmentCancelledEvent {
	return AppointmentCancelledEvent{
		BaseEvent: newBase("appointment.cancelled"),
		Payload: AppointmentPayload{
			ID:     id,
			Reason: reason,
		},
	}
}

func NewAppointmentRescheduledEvent(id uint, date, timeStr string) AppointmentRescheduledEvent {
	return AppointmentRescheduledEvent{
		BaseEvent: newBase("appointment.rescheduled"),
		Payload: AppointmentPayload{
			ID:              id,
			AppointmentDate: date,
			AppointmentTime: timeStr,
		},
	}
}

func NewAppointmentNoShowEvent(id uint) AppointmentNoShowEvent {
	return AppointmentNoShowEvent{
		BaseEvent: newBase("appointment.no_show"),
		Payload: AppointmentPayload{
			ID: id,
		},
	}
}

// ─── Medical Record Events ──────────────────────────────────────────────────

type MedicalRecordPayload struct {
	ID             uint   `json:"id"`
	PatientID      uint   `json:"patient_id,omitempty"`
	DoctorID       uint   `json:"doctor_id,omitempty"`
	VisitDate      string `json:"visit_date,omitempty"`
	ChiefComplaint string `json:"chief_complaint,omitempty"`
	Diagnosis      string `json:"diagnosis,omitempty"`
	Status         string `json:"status,omitempty"`
	Action         string `json:"action,omitempty"`
}

type MedicalRecordCreatedEvent struct {
	BaseEvent
	Payload MedicalRecordPayload `json:"payload"`
}

type MedicalRecordUpdatedEvent struct {
	BaseEvent
	Payload MedicalRecordPayload `json:"payload"`
}

type MedicalRecordDeletedEvent struct {
	BaseEvent
	Payload MedicalRecordPayload `json:"payload"`
}

type MedicalRecordRestoredEvent struct {
	BaseEvent
	Payload MedicalRecordPayload `json:"payload"`
}

type MedicalRecordFinalizedEvent struct {
	BaseEvent
	Payload MedicalRecordPayload `json:"payload"`
}

// Constructors

func NewMedicalRecordCreatedEvent(id, patientID, doctorID uint, visitDate, chiefComplaint, diagnosis, status string) MedicalRecordCreatedEvent {
	return MedicalRecordCreatedEvent{
		BaseEvent: newBase("medical_record.created"),
		Payload: MedicalRecordPayload{
			ID:             id,
			PatientID:      patientID,
			DoctorID:       doctorID,
			VisitDate:      visitDate,
			ChiefComplaint: chiefComplaint,
			Diagnosis:      diagnosis,
			Status:         status,
		},
	}
}

func NewMedicalRecordUpdatedEvent(id, patientID, doctorID uint, visitDate, chiefComplaint, diagnosis, status, action string) MedicalRecordUpdatedEvent {
	return MedicalRecordUpdatedEvent{
		BaseEvent: newBase("medical_record.updated"),
		Payload: MedicalRecordPayload{
			ID:             id,
			PatientID:      patientID,
			DoctorID:       doctorID,
			VisitDate:      visitDate,
			ChiefComplaint: chiefComplaint,
			Diagnosis:      diagnosis,
			Status:         status,
			Action:         action,
		},
	}
}

func NewMedicalRecordDeletedEvent(id uint, action string) MedicalRecordDeletedEvent {
	return MedicalRecordDeletedEvent{
		BaseEvent: newBase("medical_record.deleted"),
		Payload: MedicalRecordPayload{
			ID:     id,
			Action: action,
		},
	}
}

func NewMedicalRecordRestoredEvent(id uint) MedicalRecordRestoredEvent {
	return MedicalRecordRestoredEvent{
		BaseEvent: newBase("medical_record.restored"),
		Payload: MedicalRecordPayload{
			ID: id,
		},
	}
}

func NewMedicalRecordFinalizedEvent(id uint) MedicalRecordFinalizedEvent {
	return MedicalRecordFinalizedEvent{
		BaseEvent: newBase("medical_record.finalized"),
		Payload: MedicalRecordPayload{
			ID: id,
		},
	}
}

// ─── Hospitalization Events ────────────────────────────────────────────────────────

type HospitalizationPayload struct {
	ID                 uint   `json:"id"`
	PatientID          uint   `json:"patient_id"`
	DoctorID           uint   `json:"doctor_id"`
	RoomID             uint   `json:"room_id"`
	AdmissionDate      string `json:"admission_date"`
	ReasonForAdmission string `json:"reason_for_admission"`
	Status             string `json:"status"`
	Action             string `json:"action,omitempty"`
}

type HospitalizationCreatedEvent struct {
	BaseEvent
	Payload HospitalizationPayload `json:"payload"`
}

type HospitalizationUpdatedEvent struct {
	BaseEvent
	Payload HospitalizationPayload `json:"payload"`
}

type HospitalizationDeletedEvent struct {
	BaseEvent
	Payload struct {
		ID     uint   `json:"id"`
		Action string `json:"action"` // "soft_delete" atau "hard_delete"
	} `json:"payload"`
}

type HospitalizationRestoredEvent struct {
	BaseEvent
	Payload struct {
		ID uint `json:"id"`
	} `json:"payload"`
}

type HospitalizationDischargedEvent struct {
	BaseEvent
	Payload struct {
		ID               uint   `json:"id"`
		DischargeSummary string `json:"discharge_summary"`
	} `json:"payload"`
}

type HospitalizationTransferredEvent struct {
	BaseEvent
	Payload struct {
		ID    uint   `json:"id"`
		Notes string `json:"notes"`
	} `json:"payload"`
}

func NewHospitalizationCreatedEvent(id, patientID, doctorID, roomID uint, admissionDate, reasonForAdmission, status string) HospitalizationCreatedEvent {
	return HospitalizationCreatedEvent{
		BaseEvent: newBase("hospitalization.created"),
		Payload: HospitalizationPayload{
			ID:                 id,
			PatientID:          patientID,
			DoctorID:           doctorID,
			RoomID:             roomID,
			AdmissionDate:      admissionDate,
			ReasonForAdmission: reasonForAdmission,
			Status:             status,
		},
	}
}

func NewHospitalizationUpdatedEvent(id, patientID, doctorID, roomID uint, admissionDate, reasonForAdmission, status, action string) HospitalizationUpdatedEvent {
	return HospitalizationUpdatedEvent{
		BaseEvent: newBase("hospitalization.updated"),
		Payload: HospitalizationPayload{
			ID:                 id,
			PatientID:          patientID,
			DoctorID:           doctorID,
			RoomID:             roomID,
			AdmissionDate:      admissionDate,
			ReasonForAdmission: reasonForAdmission,
			Status:             status,
			Action:             action,
		},
	}
}

func NewHospitalizationDeletedEvent(id uint, action string) HospitalizationDeletedEvent {
	return HospitalizationDeletedEvent{
		BaseEvent: newBase("hospitalization.deleted"),
		Payload: struct {
			ID     uint   `json:"id"`
			Action string `json:"action"`
		}{
			ID:     id,
			Action: action,
		},
	}
}

func NewHospitalizationRestoredEvent(id uint) HospitalizationRestoredEvent {
	return HospitalizationRestoredEvent{
		BaseEvent: newBase("hospitalization.restored"),
		Payload: struct {
			ID uint `json:"id"`
		}{
			ID: id,
		},
	}
}

func NewHospitalizationDischargedEvent(id uint, dischargeSummary string) HospitalizationDischargedEvent {
	return HospitalizationDischargedEvent{
		BaseEvent: newBase("hospitalization.discharged"),
		Payload: struct {
			ID               uint   `json:"id"`
			DischargeSummary string `json:"discharge_summary"`
		}{
			ID:               id,
			DischargeSummary: dischargeSummary,
		},
	}
}

func NewHospitalizationTransferredEvent(id uint, notes string) HospitalizationTransferredEvent {
	return HospitalizationTransferredEvent{
		BaseEvent: newBase("hospitalization.transferred"),
		Payload: struct {
			ID    uint   `json:"id"`
			Notes string `json:"notes"`
		}{
			ID:    id,
			Notes: notes,
		},
	}
}

// ─── LabTest Events ─────────────────────────────────────────────────────────

type LabTestCreatedEvent struct {
	BaseEvent
	Payload struct {
		ID                uint   `json:"id"`
		MedicalRecordID   uint   `json:"medical_record_id"`
		TestTypeID        uint   `json:"test_type_id"`
		OrderedByDoctorID uint   `json:"ordered_by_doctor_id"`
		OrderDate         string `json:"order_date"`
		Status            string `json:"status"`
	} `json:"payload"`
}

func NewLabTestCreatedEvent(id, medicalRecordID, testTypeID, orderedByDoctorID uint, orderDate, status string) LabTestCreatedEvent {
	return LabTestCreatedEvent{
		BaseEvent: newBase("lab_test.created"),
		Payload: struct {
			ID                uint   `json:"id"`
			MedicalRecordID   uint   `json:"medical_record_id"`
			TestTypeID        uint   `json:"test_type_id"`
			OrderedByDoctorID uint   `json:"ordered_by_doctor_id"`
			OrderDate         string `json:"order_date"`
			Status            string `json:"status"`
		}{
			ID:                id,
			MedicalRecordID:   medicalRecordID,
			TestTypeID:        testTypeID,
			OrderedByDoctorID: orderedByDoctorID,
			OrderDate:         orderDate,
			Status:            status,
		},
	}
}

type LabTestUpdatedEvent struct {
	BaseEvent
	Payload struct {
		ID     uint   `json:"id"`
		Action string `json:"action"`
	} `json:"payload"`
}

func NewLabTestUpdatedEvent(id uint, action string) LabTestUpdatedEvent {
	return LabTestUpdatedEvent{
		BaseEvent: newBase("lab_test.updated"),
		Payload: struct {
			ID     uint   `json:"id"`
			Action string `json:"action"`
		}{
			ID:     id,
			Action: action,
		},
	}
}

type LabTestSampleCollectedEvent struct {
	BaseEvent
	Payload struct {
		ID                   uint   `json:"id"`
		SampleCollectionDate string `json:"sample_collection_date"`
	} `json:"payload"`
}

func NewLabTestSampleCollectedEvent(id uint, sampleCollectionDate string) LabTestSampleCollectedEvent {
	return LabTestSampleCollectedEvent{
		BaseEvent: newBase("lab_test.sample_collected"),
		Payload: struct {
			ID                   uint   `json:"id"`
			SampleCollectionDate string `json:"sample_collection_date"`
		}{
			ID:                   id,
			SampleCollectionDate: sampleCollectionDate,
		},
	}
}

type LabTestStartedEvent struct {
	BaseEvent
	Payload struct {
		ID            uint   `json:"id"`
		TestStartDate string `json:"test_start_date"`
	} `json:"payload"`
}

func NewLabTestStartedEvent(id uint, testStartDate string) LabTestStartedEvent {
	return LabTestStartedEvent{
		BaseEvent: newBase("lab_test.started"),
		Payload: struct {
			ID            uint   `json:"id"`
			TestStartDate string `json:"test_start_date"`
		}{
			ID:            id,
			TestStartDate: testStartDate,
		},
	}
}

type LabTestCompletedEvent struct {
	BaseEvent
	Payload struct {
		ID         uint   `json:"id"`
		ResultDate string `json:"result_date"`
	} `json:"payload"`
}

func NewLabTestCompletedEvent(id uint, resultDate string) LabTestCompletedEvent {
	return LabTestCompletedEvent{
		BaseEvent: newBase("lab_test.completed"),
		Payload: struct {
			ID         uint   `json:"id"`
			ResultDate string `json:"result_date"`
		}{
			ID:         id,
			ResultDate: resultDate,
		},
	}
}

type LabTestCancelledEvent struct {
	BaseEvent
	Payload struct {
		ID uint `json:"id"`
	} `json:"payload"`
}

func NewLabTestCancelledEvent(id uint) LabTestCancelledEvent {
	return LabTestCancelledEvent{
		BaseEvent: newBase("lab_test.cancelled"),
		Payload: struct {
			ID uint `json:"id"`
		}{
			ID: id,
		},
	}
}

type LabTestDeletedEvent struct {
	BaseEvent
	Payload struct {
		ID     uint   `json:"id"`
		Action string `json:"action"` // "soft_delete" | "hard_delete"
	} `json:"payload"`
}

func NewLabTestDeletedEvent(id uint, action string) LabTestDeletedEvent {
	return LabTestDeletedEvent{
		BaseEvent: newBase("lab_test.deleted"),
		Payload: struct {
			ID     uint   `json:"id"`
			Action string `json:"action"`
		}{
			ID:     id,
			Action: action,
		},
	}
}

type LabTestRestoredEvent struct {
	BaseEvent
	Payload struct {
		ID uint `json:"id"`
	} `json:"payload"`
}

func NewLabTestRestoredEvent(id uint) LabTestRestoredEvent {
	return LabTestRestoredEvent{
		BaseEvent: newBase("lab_test.restored"),
		Payload: struct {
			ID uint `json:"id"`
		}{
			ID: id,
		},
	}
}
