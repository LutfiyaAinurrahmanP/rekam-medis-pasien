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

// Medicine Events
type MedicinePayload struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	GenericName   string  `json:"generic_name"`
	BrandName     string  `json:"brand_name"`
	Type          string  `json:"type"`
	Strength      string  `json:"strength"`
	Manufacturer  string  `json:"manufacturer"`
	Unit          string  `json:"unit"`
	StockQuantity int     `json:"stock_quantity"`
	Price         float64 `json:"price"`
	IsActive      bool    `json:"is_active"`
	IsLowStock    bool    `json:"is_low_stock"`
	IsOutOfStock  bool    `json:"is_out_of_stock"`
}
