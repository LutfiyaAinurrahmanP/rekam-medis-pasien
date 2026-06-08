// NotificationEventHandler mendengarkan event domain tertentu dan
// mengirimkan notifikasi yang relevan (email, SMS, push notification).
package eventhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/config"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/events"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/mailer"
)

// NotificationEventHandler mengonsumsi event dan mengirimkan notifikasi.
type NotificationEventHandler struct {
	consumer *kafka.MultiTopicConsumer
	mailer   mailer.Mailer
}

// notificationTopics adalah daftar topic yang diproses oleh handler ini.
var notificationTopics = []string{
	kafka.TopicUserRegistered,
	kafka.TopicUserLogin,
	kafka.TopicUserCreated,
	kafka.TopicUserUpdated,
	kafka.TopicUserDeleted,
	kafka.TopicUserRestored,
	kafka.TopicUserPasswordResetRequested,
	kafka.TopicPatientCreated,
	kafka.TopicPatientUpdated,
	kafka.TopicPatientDeleted,
	kafka.TopicPatientRestored,
	kafka.TopicDoctorSpecializationCreated,
	kafka.TopicDoctorSpecializationUpdated,
	kafka.TopicDoctorSpecializationRestored,
	kafka.TopicDoctorSpecializationDeleted,
	kafka.TopicDoctorCreated,
	kafka.TopicDoctorUpdated,
	kafka.TopicDoctorDeleted,
	kafka.TopicDoctorRestored,
	kafka.TopicDepartmentCreated,
	kafka.TopicDepartmentUpdated,
	kafka.TopicDepartmentDeleted,
	kafka.TopicDepartmentRestored,
	kafka.TopicRoomCreated,
	kafka.TopicRoomUpdated,
	kafka.TopicRoomDeleted,
	kafka.TopicRoomRestored,
	kafka.TopicTypeTestCategoryCreated,
	kafka.TopicTypeTestCategoryUpdated,
	kafka.TopicTypeTestCategoryDeleted,
	kafka.TopicTypeTestCategoryRestored,
	kafka.TopicTypeTestCreated,
	kafka.TopicTypeTestUpdated,
	kafka.TopicTypeTestDeleted,
	kafka.TopicTypeTestRestored,
	kafka.TopicRoomTypeCreated,
	kafka.TopicRoomTypeUpdated,
	kafka.TopicRoomTypeDeleted,
	kafka.TopicRoomTypeRestored,
}

// NewNotificationEventHandler membuat notification handler baru.
func NewNotificationEventHandler(brokers []string, cfg *config.Config) *NotificationEventHandler {
	h := &NotificationEventHandler{}
	h.mailer = mailer.NewSMTPMailer(&cfg.SMTP)
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

	case kafka.TopicUserUpdated:
		var e events.UserUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 👤 User profile updated: %s (action=%s)", e.Payload.Username, e.Payload.Action)

	case kafka.TopicUserDeleted:
		var e events.UserDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ⚠️  User account deleted: %s (action=%s)", e.Payload.Username, e.Payload.Action)

	case kafka.TopicUserRestored:
		var e events.UserRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ♻️  User account restored: %s", e.Payload.Username)

	case kafka.TopicUserPasswordResetRequested:
		var e events.UserPasswordResetRequestedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		return h.sendPasswordResetCodeEmail(e.Payload.Email, e.Payload.Username, e.Payload.ResetCode, e.Payload.ExpiresIn)

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

	case kafka.TopicPatientRestored:
		var e events.PatientRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ♻️  Patient record restored: %s (code=%s)", e.Payload.FullName, e.Payload.PatientCode)

	// ── Doctor Specialization Events ─────────────────────────────────────────

	case kafka.TopicDoctorSpecializationCreated:
		var e events.DoctorSpecializationCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🩺 New doctor specialization created: %s (code=%s)", e.Payload.Name, e.Payload.Code)

	case kafka.TopicDoctorSpecializationUpdated:
		var e events.DoctorSpecializationUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🩺 Doctor specialization updated: %s (code=%s, action=%s)", e.Payload.Name, e.Payload.Code, e.Payload.Action)

	case kafka.TopicDoctorSpecializationRestored:
		var e events.DoctorSpecializationRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🩺 Doctor specialization restored: %s (code=%s)", e.Payload.Name, e.Payload.Code)

	case kafka.TopicDoctorSpecializationDeleted:
		var e events.DoctorSpecializationDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🩺 Doctor specialization deleted: %s (code=%s, action=%s)", e.Payload.Name, e.Payload.Code, e.Payload.Action)

	// ── Doctor Events ────────────────────────────────────────────────────────

	case kafka.TopicDoctorCreated:
		var e events.DoctorCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		h.sendDoctorOnboardingNotification(e.Payload.Email, e.Payload.FullName, e.Payload.SpecializationID)

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
	case kafka.TopicDoctorRestored:
		var e events.DoctorRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ♻️  Doctor record restored: %s", e.Payload.FullName)
	// ── Department Events ────────────────────────────────────────────────────

	case kafka.TopicDepartmentCreated:
		var e events.DepartmentCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🏥 New department created: %s (code=%s)", e.Payload.Name, e.Payload.Code)

	case kafka.TopicDepartmentUpdated:
		var e events.DepartmentUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🏥 Department updated: %s (code=%s)", e.Payload.Name, e.Payload.Code)

	case kafka.TopicDepartmentDeleted:
		var e events.DepartmentDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🗑️  Department deleted: %s (action=%s)", e.Payload.Name, e.Payload.Action)

	case kafka.TopicDepartmentRestored:
		var e events.DepartmentRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ♻️  Department restored: %s", e.Payload.Name)

	// ── Room Events ──────────────────────────────────────────────────────────

	case kafka.TopicRoomCreated:
		var e events.RoomCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}

		typeIDStr := "nil"
		if e.Payload.RoomTypeID != nil {
			typeIDStr = fmt.Sprintf("%d", *e.Payload.RoomTypeID)
		}

		log.Printf("[NOTIFICATION] 🛏️  New room created: %s (type_id=%s, capacity=%d)",
			e.Payload.RoomNumber, typeIDStr, e.Payload.BedCapacity)

	case kafka.TopicRoomUpdated:
		var e events.RoomUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		if e.Payload.Action == "occupy" || e.Payload.Action == "release" {
			log.Printf("[NOTIFICATION] 🛏️  Room %s bed status changed: available=%d (action=%s)",
				e.Payload.RoomNumber, e.Payload.AvailableBeds, e.Payload.Action)
		}

	case kafka.TopicRoomDeleted:
		var e events.RoomDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🗑️  Room deleted: %s (action=%s)", e.Payload.RoomNumber, e.Payload.Action)

	case kafka.TopicRoomRestored:
		var e events.RoomRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ♻️  Room restored: %s", e.Payload.RoomNumber)

	// ── RoomType Events ──────────────────────────────────────────────────────────

	case kafka.TopicRoomTypeCreated:
		var e events.RoomTypeCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🏢 New room type created: %s (code=%s)",
			e.Payload.Name, e.Payload.Code)

	case kafka.TopicRoomTypeUpdated:
		var e events.RoomTypeUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🏢 Room type updated: %s (code=%s, action=%s)",
			e.Payload.Name, e.Payload.Code, e.Payload.Action)

	case kafka.TopicRoomTypeDeleted:
		var e events.RoomTypeDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🗑️  Room type deleted: %s (action=%s)", e.Payload.Name, e.Payload.Action)

	case kafka.TopicRoomTypeRestored:
		var e events.RoomTypeRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ♻️  Room type restored: %s", e.Payload.Name)

	// ── Type Test Category Events ─────────────────────────────────────────
	case kafka.TopicTypeTestCategoryCreated:
		var e events.TypeTestCategoryCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🧪 New type test category created: %s (code=%s)", e.Payload.Name, e.Payload.Code)

	case kafka.TopicTypeTestCategoryUpdated:
		var e events.TypeTestCategoryUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🧪 Type test category updated: %s (code=%s, action=%s)", e.Payload.Name, e.Payload.Code, e.Payload.Action)

	case kafka.TopicTypeTestCategoryRestored:
		var e events.TypeTestCategoryRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ♻️ Type test category restored: %s (code=%s)", e.Payload.Name, e.Payload.Code)

	case kafka.TopicTypeTestCategoryDeleted:
		var e events.TypeTestCategoryDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🗑️ Type test category deleted: %s (code=%s, action=%s)", e.Payload.Name, e.Payload.Code, e.Payload.Action)

	// ── TypeTest Events ──────────────────────────────────────────────────────

	case kafka.TopicTypeTestCreated:
		var e events.TypeTestCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🧪 New test type added: %s (code=%s, category_id=%d)",
			e.Payload.Name, e.Payload.Code, e.Payload.TypeTestCategoryID)

	case kafka.TopicTypeTestUpdated:
		var e events.TypeTestUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🧪 Test type updated: %s (code=%s, action=%s)",
			e.Payload.Name, e.Payload.Code, e.Payload.Action)

	case kafka.TopicTypeTestDeleted:
		var e events.TypeTestDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🗑️ Type test deleted: %s (code=%s, action=%s)",
			e.Payload.Name, e.Payload.Code, e.Payload.Action)

	case kafka.TopicTypeTestRestored:
		var e events.TypeTestRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ♻️ Test type restored: %s (code=%s)", e.Payload.Name, e.Payload.Code)

	// ── MedicineType Events ──────────────────────────────────────────────────────

	case kafka.TopicMedicineTypeCreated:
		var e events.MedicineTypeCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 💊 New medicine type added: %s (code=%s)",
			e.Payload.Name, e.Payload.Code)

	case kafka.TopicMedicineTypeUpdated:
		var e events.MedicineTypeUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 💊 Medicine type updated: %s (code=%s, action=%s)",
			e.Payload.Name, e.Payload.Code, e.Payload.Action)

	case kafka.TopicMedicineTypeDeleted:
		var e events.MedicineTypeDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🗑️ Medicine type deleted: %s (code=%s, action=%s)",
			e.Payload.Name, e.Payload.Code, e.Payload.Action)

	case kafka.TopicMedicineTypeRestored:
		var e events.MedicineTypeRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ♻️ Medicine type restored: %s (code=%s)", e.Payload.Name, e.Payload.Code)

	// ── Medicine Events ──────────────────────────────────────────────────────────

	case kafka.TopicMedicineCreated:
		var e events.MedicineCreatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 💊 New medicine added: %s (generic: %s, brand: %s, type_id=%d)",
			e.Payload.Name, e.Payload.GenericName, e.Payload.BrandName, e.Payload.MedicineTypeID)

	case kafka.TopicMedicineUpdated:
		var e events.MedicineUpdatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 💊 Medicine updated: %s (type_id=%d, action=%s)",
			e.Payload.Name, e.Payload.MedicineTypeID, e.Payload.Action)

	case kafka.TopicMedicineDeleted:
		var e events.MedicineDeletedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 🗑️ Medicine deleted: %s (action=%s)",
			e.Payload.Name, e.Payload.Action)

	case kafka.TopicMedicineRestored:
		var e events.MedicineRestoredEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ♻️ Medicine restored: %s", e.Payload.Name)

	case kafka.TopicMedicineActivated:
		var e events.MedicineActivatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ✅ Medicine activated: %s", e.Payload.Name)

	case kafka.TopicMedicineDeactivated:
		var e events.MedicineDeactivatedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] ❌ Medicine deactivated: %s", e.Payload.Name)

	case kafka.TopicMedicineStockAdded:
		var e events.MedicineStockAddedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 📦 Medicine stock added: %s (+%d, total: %d)", e.Payload.Name, e.Payload.DeltaQuantity, e.Payload.StockQuantity)

	case kafka.TopicMedicineStockReduced:
		var e events.MedicineStockReducedEvent
		if err := json.Unmarshal(value, &e); err != nil {
			return err
		}
		log.Printf("[NOTIFICATION] 📦 Medicine stock reduced: %s (-%d, total: %d)", e.Payload.Name, e.Payload.DeltaQuantity, e.Payload.StockQuantity)
	default:
		log.Printf("[NOTIFICATION] ⚠️  Received unsupported topic: %s", topic)
	}

	return nil
}

// ─── Notification Stubs ───────────────────────────────────────────────────────

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

func (h *NotificationEventHandler) sendDoctorOnboardingNotification(email, fullName string, specializationID uint) {
	log.Printf("[NOTIFICATION] 📧 [EMAIL] Doctor onboarding → %s (%s) specialization_id=%d",
		email, fullName, specializationID)
	// TODO: Send onboarding materials
}

func loadEmailLogoAsset() ([]byte, string, bool) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, "", false
	}

	logoPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "public", "logo", "app-logo.png")
	logoBytes, err := os.ReadFile(logoPath)
	if err != nil {
		log.Printf("[NOTIFICATION] unable to load email logo from %s: %v", logoPath, err)
		return nil, "", false
	}

	return logoBytes, filepath.Base(logoPath), true
}

func (h *NotificationEventHandler) sendPasswordResetCodeEmail(email, username, resetCode string, expiresIn int) error {
	appName := "Medika Health Care"
	subject := fmt.Sprintf("%s - Password Reset Code", appName)
	greeting := username
	if greeting == "" {
		greeting = email
	}

	spaced := ""
	for i, c := range resetCode {
		if i > 0 {
			spaced += "&nbsp;"
		}
		spaced += string(c)
	}

	logoBytes, logoFilename, logoOK := loadEmailLogoAsset()
	logoHTML := `<div class="logo-fallback">M</div>`
	if logoOK {
		logoHTML = `<img src="cid:app-logo" alt="` + appName + `" style="width:52px;height:52px;border-radius:16px;object-fit:contain;display:inline-block;"/>`
	}

	body := fmt.Sprintf(`
<!doctype html>
<html>
<head>
	<meta charset="UTF-8"/>
	<meta name="viewport" content="width=device-width,initial-scale=1.0"/>
	<title>Password Reset</title>
	<style>
		body{margin:0;padding:0;background:#f1faea;font-family:Inter,system-ui,-apple-system,sans-serif}
		.wrap{max-width:600px;margin:0 auto;padding:28px 16px}
		.logo-wrap{text-align:center;padding-bottom:16px}
		.logo-fallback{display:inline-flex;align-items:center;justify-content:center;width:52px;height:52px;border-radius:16px;background:#70b844;color:#fff;font-size:22px;font-weight:800}
		.card{background:#fff;border-radius:24px;overflow:hidden;border:1px solid #b8e09a}
		.header{background:#70b844;padding:28px 32px}
		.header-title{margin:0;color:#fff;font-size:20px;font-weight:800}
		.header-sub{margin:8px 0 0;color:#d6f0c2;font-size:13px;line-height:1.6}
		.body{padding:32px}
		.badge{display:inline-block;background:#70b844;color:#fff;font-size:11px;font-weight:700;padding:4px 12px;border-radius:20px;margin-bottom:16px}
		.body h1{margin:0 0 18px;font-size:24px;font-weight:800;color:#2d5a1b}
		.body p{margin:0 0 24px;font-size:15px;line-height:1.75;color:#374151}
		.code-box{background:#f1faea;border:1px solid #b8e09a;border-radius:16px;padding:24px;text-align:center;margin-bottom:28px}
		.code-label{margin:0 0 10px;font-size:11px;font-weight:700;text-transform:uppercase;letter-spacing:0.18em;color:#70b844}
		.code-value{margin:0;font-size:36px;font-weight:800;letter-spacing:10px;color:#2d5a1b}
		.btn-wrap{text-align:center;margin-bottom:28px}
		.btn{display:inline-block;padding:14px 28px;background:#70b844;color:#fff;text-decoration:none;font-size:15px;font-weight:700;border-radius:12px}
		.divider{height:1px;background:#d6f0c2;margin:24px 0}
		.info-row{display:flex;align-items:flex-start;gap:10px;background:#f1faea;border-radius:10px;padding:12px 16px;margin-bottom:10px;font-size:14px;color:#374151}
		.note{margin:20px 0 0;font-size:13px;line-height:1.8;color:#6b7280;padding:14px 16px;background:#f9fafb;border-left:3px solid #b8e09a;border-radius:0 8px 8px 0}
		.footer{padding:20px 16px;text-align:center;font-size:12px;color:#9ca3af;line-height:1.7}
	</style>
</head>
<body>
<div class="wrap">
	<div class="logo-wrap">
		%s
	</div>
	<div class="card">
		<div class="header">
			<p class="header-title">%s</p>
			<p class="header-sub">Password reset verification</p>
		</div>
		<div class="body">
			<span class="badge">🔒 Security Verification</span>
			<h1>Forgot your password?</h1>
			<p>Hello, <strong style="color:#4a8a28;">%s</strong></p>
			<p>We received a request to reset the password for your account <strong style="color:#4a8a28;">%s</strong>. Use the code below to continue. This code will expire in <strong style="color:#4a8a28;">%d minutes</strong>.</p>
			<div class="code-box">
				<p class="code-label">Your Verification Code</p>
				<p class="code-value">%s</p>
			</div>
			<div class="divider"></div>
			<div class="info-row">
				<span>⚠️</span>
				<span>Never share this code with anyone, including our support team.</span>
			</div>
			<div class="info-row">
				<span>🛡️</span>
				<span>If you did not request this reset, please ignore this email — your password remains secure.</span>
			</div>
			<p class="note">Need help? Open the app and contact our support team if this message looks unfamiliar.</p>
		</div>
	</div>
	<div class="footer">
		<p style="margin:0 0 4px;">%s &nbsp;·&nbsp; Sent automatically</p>
		<p style="margin:0;">© 2025 %s. All rights reserved.</p>
	</div>
</div>
</body>
</html>`,
		logoHTML,
		appName,
		greeting,
		email,
		expiresIn/60,
		spaced,
		appName,
		appName,
	)

	log.Printf("[NOTIFICATION] 📧 [EMAIL] Password reset code → %s (%s)", email, username)
	if h.mailer == nil {
		return nil
	}
	if logoOK {
		if err := h.mailer.SendWithInlineImage(email, subject, body, "app-logo", logoFilename, "image/png", logoBytes); err == nil {
			return nil
		}
	}
	return h.mailer.Send(email, subject, body)
}
