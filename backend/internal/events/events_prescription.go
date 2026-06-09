package events

// ─── Prescription Events ────────────────────────────────────────────────────────

type PrescriptionPayload struct {
	ID               uint   `json:"id"`
	MedicalRecordID  uint   `json:"medical_record_id"`
	DoctorID         uint   `json:"doctor_id"`
	PrescriptionDate string `json:"prescription_date"`
	Notes            string `json:"notes"`
	Status           string `json:"status"`
	Action           string `json:"action,omitempty"`
}

type PrescriptionCreatedEvent struct {
	BaseEvent
	Payload PrescriptionPayload `json:"payload"`
}

type PrescriptionUpdatedEvent struct {
	BaseEvent
	Payload PrescriptionPayload `json:"payload"`
}

type PrescriptionDispensedEvent struct {
	BaseEvent
	Payload PrescriptionPayload `json:"payload"`
}

type PrescriptionCancelledEvent struct {
	BaseEvent
	Payload PrescriptionPayload `json:"payload"`
}

type PrescriptionDeletedEvent struct {
	BaseEvent
	Payload PrescriptionPayload `json:"payload"`
}

type PrescriptionRestoredEvent struct {
	BaseEvent
	Payload PrescriptionPayload `json:"payload"`
}

func NewPrescriptionCreatedEvent(id, mrID, docID uint, date, notes, status string) PrescriptionCreatedEvent {
	return PrescriptionCreatedEvent{
		BaseEvent: newBase("prescription.created"),
		Payload: PrescriptionPayload{
			ID: id, MedicalRecordID: mrID, DoctorID: docID, PrescriptionDate: date, Notes: notes, Status: status,
		},
	}
}

func NewPrescriptionUpdatedEvent(id, mrID, docID uint, date, notes, status, action string) PrescriptionUpdatedEvent {
	return PrescriptionUpdatedEvent{
		BaseEvent: newBase("prescription.updated"),
		Payload: PrescriptionPayload{
			ID: id, MedicalRecordID: mrID, DoctorID: docID, PrescriptionDate: date, Notes: notes, Status: status, Action: action,
		},
	}
}

func NewPrescriptionDispensedEvent(id uint, status string) PrescriptionDispensedEvent {
	return PrescriptionDispensedEvent{
		BaseEvent: newBase("prescription.dispensed"),
		Payload: PrescriptionPayload{
			ID: id, Status: status,
		},
	}
}

func NewPrescriptionCancelledEvent(id uint, status string) PrescriptionCancelledEvent {
	return PrescriptionCancelledEvent{
		BaseEvent: newBase("prescription.cancelled"),
		Payload: PrescriptionPayload{
			ID: id, Status: status,
		},
	}
}

func NewPrescriptionDeletedEvent(id uint, action string) PrescriptionDeletedEvent {
	return PrescriptionDeletedEvent{
		BaseEvent: newBase("prescription.deleted"),
		Payload: PrescriptionPayload{
			ID: id, Action: action,
		},
	}
}

func NewPrescriptionRestoredEvent(id uint) PrescriptionRestoredEvent {
	return PrescriptionRestoredEvent{
		BaseEvent: newBase("prescription.restored"),
		Payload: PrescriptionPayload{
			ID: id,
		},
	}
}

// ─── Prescription Item Events ────────────────────────────────────────────────────────

type PrescriptionItemPayload struct {
	ID             uint   `json:"id"`
	PrescriptionID uint   `json:"prescription_id"`
	MedicineID     uint   `json:"medicine_id"`
	Dosage         string `json:"dosage"`
	Frequency      string `json:"frequency"`
	DurationDays   int    `json:"duration_days"`
	Quantity       int    `json:"quantity"`
	Instructions   string `json:"instructions"`
	Action         string `json:"action,omitempty"`
}

type PrescriptionItemCreatedEvent struct {
	BaseEvent
	Payload PrescriptionItemPayload `json:"payload"`
}

type PrescriptionItemUpdatedEvent struct {
	BaseEvent
	Payload PrescriptionItemPayload `json:"payload"`
}

type PrescriptionItemDeletedEvent struct {
	BaseEvent
	Payload PrescriptionItemPayload `json:"payload"`
}

func NewPrescriptionItemCreatedEvent(id, presID, medID uint, dosage, freq string, duration, qty int, instructions string) PrescriptionItemCreatedEvent {
	return PrescriptionItemCreatedEvent{
		BaseEvent: newBase("prescription_item.created"),
		Payload: PrescriptionItemPayload{
			ID: id, PrescriptionID: presID, MedicineID: medID, Dosage: dosage, Frequency: freq, DurationDays: duration, Quantity: qty, Instructions: instructions,
		},
	}
}

func NewPrescriptionItemUpdatedEvent(id, presID, medID uint, dosage, freq string, duration, qty int, instructions, action string) PrescriptionItemUpdatedEvent {
	return PrescriptionItemUpdatedEvent{
		BaseEvent: newBase("prescription_item.updated"),
		Payload: PrescriptionItemPayload{
			ID: id, PrescriptionID: presID, MedicineID: medID, Dosage: dosage, Frequency: freq, DurationDays: duration, Quantity: qty, Instructions: instructions, Action: action,
		},
	}
}

func NewPrescriptionItemDeletedEvent(id uint, action string) PrescriptionItemDeletedEvent {
	return PrescriptionItemDeletedEvent{
		BaseEvent: newBase("prescription_item.deleted"),
		Payload: PrescriptionItemPayload{
			ID: id, Action: action,
		},
	}
}
