package events

type VitalSignPayload struct {
	ID              uint   `json:"id"`
	MedicalRecordID uint   `json:"medical_record_id"`
	Action          string `json:"action,omitempty"`
}

type VitalSignCreatedEvent struct {
	BaseEvent
	Payload VitalSignPayload `json:"payload"`
}

type VitalSignUpdatedEvent struct {
	BaseEvent
	Payload VitalSignPayload `json:"payload"`
}

type VitalSignDeletedEvent struct {
	BaseEvent
	Payload VitalSignPayload `json:"payload"`
}

type VitalSignRestoredEvent struct {
	BaseEvent
	Payload VitalSignPayload `json:"payload"`
}

func NewVitalSignCreatedEvent(id uint, medicalRecordID uint) VitalSignCreatedEvent {
	return VitalSignCreatedEvent{
		BaseEvent: newBase("vital_sign.created"),
		Payload: VitalSignPayload{
			ID:              id,
			MedicalRecordID: medicalRecordID,
		},
	}
}

func NewVitalSignUpdatedEvent(id uint, medicalRecordID uint) VitalSignUpdatedEvent {
	return VitalSignUpdatedEvent{
		BaseEvent: newBase("vital_sign.updated"),
		Payload: VitalSignPayload{
			ID:              id,
			MedicalRecordID: medicalRecordID,
		},
	}
}

func NewVitalSignDeletedEvent(id uint, medicalRecordID uint, action string) VitalSignDeletedEvent {
	return VitalSignDeletedEvent{
		BaseEvent: newBase("vital_sign.deleted"),
		Payload: VitalSignPayload{
			ID:              id,
			MedicalRecordID: medicalRecordID,
			Action:          action,
		},
	}
}

func NewVitalSignRestoredEvent(id uint, medicalRecordID uint) VitalSignRestoredEvent {
	return VitalSignRestoredEvent{
		BaseEvent: newBase("vital_sign.restored"),
		Payload: VitalSignPayload{
			ID:              id,
			MedicalRecordID: medicalRecordID,
		},
	}
}
