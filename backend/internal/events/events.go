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
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
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
	Payload struct {
		UserPayload
		Action string `json:"action"` // "profile_update", "password_change", "activate", "deactivate"
	} `json:"payload"`
}

type UserDeletedEvent struct {
	BaseEvent
	Payload struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Action   string `json:"action"` // "soft_delete", "hard_delete"
	} `json:"payload"`
}

type UserRestoredEvent struct {
	BaseEvent
	Payload struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
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
	e := UserUpdatedEvent{BaseEvent: newBase("user.updated")}
	e.Payload.UserPayload = UserPayload{
		ID: id, Username: username, Email: email,
		Phone: phone, Role: role, IsActive: isActive,
	}
	e.Payload.Action = action
	return e
}

func NewUserDeletedEvent(userID uint, username, action string) UserDeletedEvent {
	e := UserDeletedEvent{BaseEvent: newBase("user.deleted")}
	e.Payload.UserID = userID
	e.Payload.Username = username
	e.Payload.Action = action
	return e
}

func NewUserRestoredEvent(userID uint, username string) UserRestoredEvent {
	e := UserRestoredEvent{BaseEvent: newBase("user.restored")}
	e.Payload.UserID = userID
	e.Payload.Username = username
	return e
}

// ─── Patient Events ───────────────────────────────────────────────────────────

type PatientPayload struct {
	ID                uint   `json:"id"`
	PatientCode       string `json:"patient_code"`
	FullName          string `json:"full_name"`
	DateOfBirth       string `json:"date_of_birth"`
	Gender            string `json:"gender"`
	BloodType         string `json:"blood_type"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	InsuranceProvider string `json:"insurance_provider"`
}

type PatientCreatedEvent struct {
	BaseEvent
	Payload PatientPayload `json:"payload"`
}

type PatientUpdatedEvent struct {
	BaseEvent
	Payload struct {
		PatientPayload
		Action string `json:"action"` // "admin_update", "self_update"
	} `json:"payload"`
}

type PatientDeletedEvent struct {
	BaseEvent
	Payload struct {
		PatientID   uint   `json:"patient_id"`
		PatientCode string `json:"patient_code"`
		FullName    string `json:"full_name"`
		Action      string `json:"action"` // "soft_delete", "hard_delete"
	} `json:"payload"`
}

type PatientRestoredEvent struct {
	BaseEvent
	Payload struct {
		PatientID   uint   `json:"patient_id"`
		PatientCode string `json:"patient_code"`
		FullName    string `json:"full_name"`
	} `json:"payload"`
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
	e := PatientUpdatedEvent{BaseEvent: newBase("patient.updated")}
	e.Payload.PatientPayload = PatientPayload{
		ID: id, PatientCode: code, FullName: fullName, DateOfBirth: dob,
		Gender: gender, BloodType: blood, Phone: phone, Email: email,
		InsuranceProvider: insurance,
	}
	e.Payload.Action = action
	return e
}

func NewPatientDeletedEvent(id uint, code, fullName, action string) PatientDeletedEvent {
	e := PatientDeletedEvent{BaseEvent: newBase("patient.deleted")}
	e.Payload.PatientID = id
	e.Payload.PatientCode = code
	e.Payload.FullName = fullName
	e.Payload.Action = action
	return e
}

func NewPatientRestoredEvent(id uint, code, fullName string) PatientRestoredEvent {
	e := PatientRestoredEvent{BaseEvent: newBase("patient.restored")}
	e.Payload.PatientID = id
	e.Payload.PatientCode = code
	e.Payload.FullName = fullName
	return e
}

// ─── Doctor Events ────────────────────────────────────────────────────────────

type DoctorPayload struct {
	ID             uint    `json:"id"`
	FullName       string  `json:"full_name"`
	Specialization string  `json:"specialization"`
	Phone          string  `json:"phone"`
	Email          string  `json:"email"`
	DepartmentID   *uint   `json:"department_id,omitempty"`
	IsActive       bool    `json:"is_active"`
}

type DoctorCreatedEvent struct {
	BaseEvent
	Payload DoctorPayload `json:"payload"`
}

type DoctorUpdatedEvent struct {
	BaseEvent
	Payload struct {
		DoctorPayload
		Action string `json:"action"`
	} `json:"payload"`
}

type DoctorDeletedEvent struct {
	BaseEvent
	Payload struct {
		DoctorID uint   `json:"doctor_id"`
		FullName string `json:"full_name"`
		Action   string `json:"action"` // "soft_delete", "hard_delete"
	} `json:"payload"`
}

type DoctorRestoredEvent struct {
	BaseEvent
	Payload struct {
		DoctorID uint   `json:"doctor_id"`
		FullName string `json:"full_name"`
	} `json:"payload"`
}

// Constructors ────────────────────────────────────────────────────────────────

func NewDoctorCreatedEvent(id uint, fullName, specialization, phone, email string, deptID *uint, isActive bool) DoctorCreatedEvent {
	return DoctorCreatedEvent{
		BaseEvent: newBase("doctor.created"),
		Payload: DoctorPayload{
			ID: id, FullName: fullName, Specialization: specialization,
			Phone: phone, Email: email, DepartmentID: deptID, IsActive: isActive,
		},
	}
}

func NewDoctorUpdatedEvent(id uint, fullName, specialization, phone, email string, deptID *uint, isActive bool, action string) DoctorUpdatedEvent {
	e := DoctorUpdatedEvent{BaseEvent: newBase("doctor.updated")}
	e.Payload.DoctorPayload = DoctorPayload{
		ID: id, FullName: fullName, Specialization: specialization,
		Phone: phone, Email: email, DepartmentID: deptID, IsActive: isActive,
	}
	e.Payload.Action = action
	return e
}

func NewDoctorDeletedEvent(id uint, fullName, action string) DoctorDeletedEvent {
	e := DoctorDeletedEvent{BaseEvent: newBase("doctor.deleted")}
	e.Payload.DoctorID = id
	e.Payload.FullName = fullName
	e.Payload.Action = action
	return e
}

func NewDoctorRestoredEvent(id uint, fullName string) DoctorRestoredEvent {
	e := DoctorRestoredEvent{BaseEvent: newBase("doctor.restored")}
	e.Payload.DoctorID = id
	e.Payload.FullName = fullName
	return e
}

// ─── Department Events ────────────────────────────────────────────────────────

type DepartmentPayload struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
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
	Payload struct {
		DepartmentID uint   `json:"department_id"`
		Name         string `json:"name"`
		Action       string `json:"action"` // "soft_delete", "hard_delete"
	} `json:"payload"`
}

type DepartmentRestoredEvent struct {
	BaseEvent
	Payload struct {
		DepartmentID uint   `json:"department_id"`
		Name         string `json:"name"`
	} `json:"payload"`
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
	e := DepartmentDeletedEvent{BaseEvent: newBase("department.deleted")}
	e.Payload.DepartmentID = id
	e.Payload.Name = name
	e.Payload.Action = action
	return e
}

func NewDepartmentRestoredEvent(id uint, name string) DepartmentRestoredEvent {
	e := DepartmentRestoredEvent{BaseEvent: newBase("department.restored")}
	e.Payload.DepartmentID = id
	e.Payload.Name = name
	return e
}

// ─── Room Events ──────────────────────────────────────────────────────────────

type RoomPayload struct {
	ID            uint    `json:"id"`
	RoomNumber    string  `json:"room_number"`
	RoomType      string  `json:"room_type"`
	DepartmentID  *uint   `json:"department_id,omitempty"`
	TotalBeds     int     `json:"total_beds"`
	AvailableBeds int     `json:"available_beds"`
	PricePerDay   float64 `json:"price_per_day"`
	IsActive      bool    `json:"is_active"`
}

type RoomCreatedEvent struct {
	BaseEvent
	Payload RoomPayload `json:"payload"`
}

type RoomUpdatedEvent struct {
	BaseEvent
	Payload struct {
		RoomPayload
		Action string `json:"action"`
	} `json:"payload"`
}

type RoomDeletedEvent struct {
	BaseEvent
	Payload struct {
		RoomID     uint   `json:"room_id"`
		RoomNumber string `json:"room_number"`
		Action     string `json:"action"` // "soft_delete", "hard_delete"
	} `json:"payload"`
}

type RoomRestoredEvent struct {
	BaseEvent
	Payload struct {
		RoomID     uint   `json:"room_id"`
		RoomNumber string `json:"room_number"`
	} `json:"payload"`
}

// Constructors ────────────────────────────────────────────────────────────────

func NewRoomCreatedEvent(id uint, roomNumber, roomType string, deptID *uint, total, available int, price float64, isActive bool) RoomCreatedEvent {
	return RoomCreatedEvent{
		BaseEvent: newBase("room.created"),
		Payload: RoomPayload{
			ID: id, RoomNumber: roomNumber, RoomType: roomType, DepartmentID: deptID,
			TotalBeds: total, AvailableBeds: available, PricePerDay: price, IsActive: isActive,
		},
	}
}

func NewRoomUpdatedEvent(id uint, roomNumber, roomType string, deptID *uint, total, available int, price float64, isActive bool, action string) RoomUpdatedEvent {
	e := RoomUpdatedEvent{BaseEvent: newBase("room.updated")}
	e.Payload.RoomPayload = RoomPayload{
		ID: id, RoomNumber: roomNumber, RoomType: roomType, DepartmentID: deptID,
		TotalBeds: total, AvailableBeds: available, PricePerDay: price, IsActive: isActive,
	}
	e.Payload.Action = action
	return e
}

func NewRoomDeletedEvent(id uint, roomNumber, action string) RoomDeletedEvent {
	e := RoomDeletedEvent{BaseEvent: newBase("room.deleted")}
	e.Payload.RoomID = id
	e.Payload.RoomNumber = roomNumber
	e.Payload.Action = action
	return e
}

func NewRoomRestoredEvent(id uint, roomNumber string) RoomRestoredEvent {
	e := RoomRestoredEvent{BaseEvent: newBase("room.restored")}
	e.Payload.RoomID = id
	e.Payload.RoomNumber = roomNumber
	return e
}

// ─── TypeTest Events ──────────────────────────────────────────────────────────

type TypeTestPayload struct {
	ID          uint    `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	IsActive    bool    `json:"is_active"`
}

type TypeTestCreatedEvent struct {
	BaseEvent
	Payload TypeTestPayload `json:"payload"`
}

type TypeTestUpdatedEvent struct {
	BaseEvent
	Payload struct {
		TypeTestPayload
		Action string `json:"action"`
	} `json:"payload"`
}

type TypeTestDeletedEvent struct {
	BaseEvent
	Payload struct {
		TypeTestID uint   `json:"typetest_id"`
		Code       string `json:"code"`
		Name       string `json:"name"`
		Action     string `json:"action"` // "soft_delete", "hard_delete"
	} `json:"payload"`
}

type TypeTestRestoredEvent struct {
	BaseEvent
	Payload struct {
		TypeTestID uint   `json:"typetest_id"`
		Code       string `json:"code"`
		Name       string `json:"name"`
	} `json:"payload"`
}

// Constructors ────────────────────────────────────────────────────────────────

func NewTypeTestCreatedEvent(id uint, code, name, category string, price float64, isActive bool) TypeTestCreatedEvent {
	return TypeTestCreatedEvent{
		BaseEvent: newBase("typetest.created"),
		Payload:   TypeTestPayload{ID: id, Code: code, Name: name, Category: category, Price: price, IsActive: isActive},
	}
}

func NewTypeTestUpdatedEvent(id uint, code, name, category string, price float64, isActive bool, action string) TypeTestUpdatedEvent {
	e := TypeTestUpdatedEvent{BaseEvent: newBase("typetest.updated")}
	e.Payload.TypeTestPayload = TypeTestPayload{
		ID: id, Code: code, Name: name, Category: category, Price: price, IsActive: isActive,
	}
	e.Payload.Action = action
	return e
}

func NewTypeTestDeletedEvent(id uint, code, name, action string) TypeTestDeletedEvent {
	e := TypeTestDeletedEvent{BaseEvent: newBase("typetest.deleted")}
	e.Payload.TypeTestID = id
	e.Payload.Code = code
	e.Payload.Name = name
	e.Payload.Action = action
	return e
}

func NewTypeTestRestoredEvent(id uint, code, name string) TypeTestRestoredEvent {
	e := TypeTestRestoredEvent{BaseEvent: newBase("typetest.restored")}
	e.Payload.TypeTestID = id
	e.Payload.Code = code
	e.Payload.Name = name
	return e
}
