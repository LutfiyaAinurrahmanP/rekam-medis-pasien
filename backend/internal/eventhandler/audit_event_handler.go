// Package eventhandler berisi consumer handler untuk setiap jenis event Kafka.
// AuditEventHandler mencatat SEMUA event yang masuk ke structured log —
// ini berfungsi sebagai audit trail untuk keperluan kepatuhan dan debugging.
package eventhandler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/kafka"
)

// ─── Audit Event Handler ─────────────────────────────────────────────────────

// AuditEntry adalah format log audit yang terstruktur.
type AuditEntry struct {
	Timestamp  time.Time       `json:"timestamp"`
	Topic      string          `json:"topic"`
	RawPayload json.RawMessage `json:"payload"`
}

// AuditEventHandler mengonsumsi semua topic dan mencatat event ke log.
// Di production ini bisa diarahkan ke Elasticsearch, ClickHouse, atau DB audit.
type AuditEventHandler struct {
	consumer *kafka.MultiTopicConsumer
}

// NewAuditEventHandler membuat handler audit yang mendengarkan semua topic.
func NewAuditEventHandler(brokers []string, topics []string) *AuditEventHandler {
	h := &AuditEventHandler{}
	h.consumer = kafka.NewMultiTopicConsumer(
		"audit-consumer",
		brokers,
		topics,
		"sirekam-audit-group",
		h.handle,
	)
	return h
}

// Start menjalankan consumer secara blocking. Panggil dalam goroutine.
func (h *AuditEventHandler) Start(ctx context.Context) {
	h.consumer.Start(ctx)
}

// Close menutup consumer.
func (h *AuditEventHandler) Close() error {
	return h.consumer.Close()
}

// handle dipanggil untuk setiap pesan yang masuk.
func (h *AuditEventHandler) handle(ctx context.Context, topic string, key, value []byte) error {
	entry := AuditEntry{
		Timestamp:  time.Now().UTC(),
		Topic:      topic,
		RawPayload: json.RawMessage(value),
	}

	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("[AUDIT] ⚠️  Failed to marshal audit entry for topic %s: %v", topic, err)
		return nil // jangan return error agar tidak stuck
	}

	log.Printf("[AUDIT] 📋 %s", string(data))
	return nil
}
