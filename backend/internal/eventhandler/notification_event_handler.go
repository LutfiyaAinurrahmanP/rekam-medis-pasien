// NotificationEventHandler mendengarkan event domain tertentu dan
// mengirimkan notifikasi yang relevan (email, SMS, push notification).
// Implementasi saat ini hanya mencetak log; di production hubungkan ke
// layanan email/SMS sungguhan (SendGrid, Twilio, Firebase, dsb.).
package eventhandler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

// NotificationEventHandler mengonsumsi event dan mengirimkan notifikasi.
type NotificationEventHandler struct {
	consumer *kafka.MultiTopicConsumer
}

// notificationTopics adalah daftar topic yang diproses oleh handler ini.
var notificationTopics = []string{
	kafka.TopicUserRegistered,
	kafka.TopicUserLogin,
	kafka.TopicUserCreated,
	kafka.TopicUserDeleted,
	kafka.TopicPatientCreated,
	kafka.TopicPatientUpdated,
	kafka.TopicPatientDeleted,
	kafka.TopicDoctorCreated,
	kafka.TopicDoctorUpdated,
	kafka.TopicDoctorDeleted,
	kafka.TopicDepartmentCreated,
	kafka.TopicDepartmentDeleted,
	kafka.TopicRoomCreated,
	kafka.TopicRoomUpdated,
	kafka.TopicTypeTestCreated,
}

// NewNotificationEventHandler membuat notification handler baru.
func NewNotificationEventHandler(brokers []string) *NotificationEventHandler {
	h := &NotificationEventHandler{}
	h.consumer = kafka.NewMultiTopicConsumer(
		"notification-consumer",
		brokers,
		notificationTopics,
		"sirekam-notification-group",
		h.handle,
	)
	return h
}

// Start menjalankan consumer (blocking). Panggil dalam goroutine.
func (h *NotificationEventHandler) Start(ctx context.Context) {
	h.consumer.Start(ctx)
}

// Close menutup consumer.
func (h *NotificationEventHandler) Close() error {
	return h.consumer.Close()
}

// handle mendispatch event ke notifier yang sesuai berdasarkan topic.
func (h *NotificationEventHandler) handle(ctx context.Context, topic string, key, value []byte) error {
	switch topic {

	// ── User Events ──────────────────────────────────────────────────────────

	case kafka.TopicUserRegistered:
		var e events.UserRegisteredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		h.sendWelcomeEmail(e.Payload.Email, e.Payload.Username)

	case kafka.TopicUserLogin:
		var e events.UserLoginEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🔑 User login: %s (id=%d)", e.Payload.Username, e.Payload.UserID)

	case kafka.TopicUserCreated:
		var e events.UserCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		h.sendAccountCreatedNotification(e.Payload.Email, e.Payload.Username, e.Payload.Role)

	case kafka.TopicUserDeleted:
		var e events.UserDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ⚠️  User account deleted: %s (action=%s)", e.Payload.Username, e.Payload.Action)

	// ── Patient Events ───────────────────────────────────────────────────────

	case kafka.TopicPatientCreated:
		var e events.PatientCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		h.sendPatientRegistrationConfirmation(e.Payload.Email, e.Payload.FullName, e.Payload.PatientCode)

	case kafka.TopicPatientUpdated:
		var e events.PatientUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 📋 Patient data updated: %s (code=%s, action=%s)",
			e.Payload.FullName, e.Payload.PatientCode, e.Payload.Action)

	case kafka.TopicPatientDeleted:
		var e events.PatientDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🗑️  Patient record deleted: %s (action=%s)", e.Payload.FullName, e.Payload.Action)

	// ── Doctor Events ────────────────────────────────────────────────────────

	case kafka.TopicDoctorCreated:
		var e events.DoctorCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		h.sendDoctorOnboardingNotification(e.Payload.Email, e.Payload.FullName, e.Payload.Specialization)

	case kafka.TopicDoctorUpdated:
		var e events.DoctorUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 👨‍⚕️ Doctor profile updated: %s (action=%s)", e.Payload.FullName, e.Payload.Action)

	case kafka.TopicDoctorDeleted:
		var e events.DoctorDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🗑️  Doctor record deleted: %s (action=%s)", e.Payload.FullName, e.Payload.Action)

	// ── Department Events ────────────────────────────────────────────────────

	case kafka.TopicDepartmentCreated:
		var e events.DepartmentCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🏥 New department created: %s (code=%s)", e.Payload.Name, e.Payload.Code)

	case kafka.TopicDepartmentDeleted:
		var e events.DepartmentDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🗑️  Department deleted: %s (action=%s)", e.Payload.Name, e.Payload.Action)

	// ── Room Events ──────────────────────────────────────────────────────────

	case kafka.TopicRoomCreated:
		var e events.RoomCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🛏️  New room created: %s (type=%s, beds=%d)",
			e.Payload.RoomNumber, e.Payload.RoomType, e.Payload.TotalBeds)

	case kafka.TopicRoomUpdated:
		var e events.RoomUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		if e.Payload.Action == "occupy" || e.Payload.Action == "release" {
			log.Printf("[NOTIFICATION] 🛏️  Room %s bed status changed: available=%d (action=%s)",
				e.Payload.RoomNumber, e.Payload.AvailableBeds, e.Payload.Action)
		}

	// ── TypeTest Events ──────────────────────────────────────────────────────

	case kafka.TopicTypeTestCreated:
		var e events.TypeTestCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🧪 New test type added: %s (code=%s, category=%s)",
			e.Payload.Name, e.Payload.Code, e.Payload.Category)
	}

	return nil
}

// ─── Notification Stubs ───────────────────────────────────────────────────────
// Ganti implementasi ini dengan layanan email/SMS asli di production.

func (h *NotificationEventHandler) sendWelcomeEmail(email, username string) {
	log.Printf("[NOTIFICATION] 📧 [EMAIL] Welcome email → %s (%s)", email, username)
	// TODO: Integrate with SendGrid / AWS SES
}

func (h *NotificationEventHandler) sendAccountCreatedNotification(email, username, role string) {
	log.Printf("[NOTIFICATION] 📧 [EMAIL] Account created → %s (%s) role=%s", email, username, role)
	// TODO: Integrate with email service
}

func (h *NotificationEventHandler) sendPatientRegistrationConfirmation(email, fullName, patientCode string) {
	log.Printf("[NOTIFICATION] 📧 [EMAIL] Patient registration confirmation → %s (%s) code=%s",
		email, fullName, patientCode)
	// TODO: Send booklet / patient welcome kit info
}

func (h *NotificationEventHandler) sendDoctorOnboardingNotification(email, fullName, specialization string) {
	log.Printf("[NOTIFICATION] 📧 [EMAIL] Doctor onboarding → %s (%s) spec=%s",
		email, fullName, specialization)
	// TODO: Send onboarding materials
}
