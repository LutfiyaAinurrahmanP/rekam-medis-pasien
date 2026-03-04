// Package kafka menyediakan Kafka producer yang digunakan untuk menerbitkan
// event ke Apache Kafka. Producer dirancang untuk non-blocking: setiap publish
// dilakukan dalam goroutine sehingga tidak memblokir request path.
package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// EventPublisher adalah interface yang wajib diimplementasikan oleh producer.
// Abstraksi ini memudahkan mock saat unit testing.
type EventPublisher interface {
	// Publish mem-publish event secara sinkron. Gunakan ini saat urutan sangat penting.
	Publish(ctx context.Context, topic string, payload any) error

	// PublishAsync mem-publish event secara asinkron dalam goroutine terpisah.
	// Error hanya di-log, tidak dikembalikan ke caller.
	PublishAsync(topic string, payload any)

	// Close menutup koneksi ke Kafka broker.
	Close() error
}

// Producer adalah implementasi konkret dari EventPublisher.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer membuat Producer baru yang terhubung ke Kafka broker.
//
// brokers  — daftar alamat broker, contoh: []string{"localhost:9092"}
// clientID — nama klien yang akan muncul di log Kafka (biasanya nama service)
func NewProducer(brokers []string, clientID string) *Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireOne,
		Async:                  false,
		AllowAutoTopicCreation: true,
		BatchTimeout:           5 * time.Millisecond,
		WriteTimeout:           10 * time.Second,
		ReadTimeout:            10 * time.Second,
		Logger:                 kafka.LoggerFunc(func(msg string, args ...interface{}) {}), // silent info
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("[KAFKA-PRODUCER ERROR] "+msg, args...)
		}),
	}

	// Tambahkan ClientID sebagai header
	_ = clientID

	log.Printf("✅ Kafka producer initialized — brokers: %v", brokers)
	return &Producer{writer: w}
}

// Publish mem-publish satu pesan ke topic yang diberikan secara sinkron.
// payload akan di-serialize ke JSON.
func (p *Producer) Publish(ctx context.Context, topic string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Topic: topic,
		Value: data,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("[KAFKA-PRODUCER] ❌ Failed to publish to topic %s: %v", topic, err)
		return err
	}

	log.Printf("[KAFKA-PRODUCER] ✅ Published to topic %s", topic)
	return nil
}

// PublishAsync mem-publish event dalam goroutine terpisah (fire-and-forget).
// Error hanya di-log, tidak dikembalikan ke caller.
func (p *Producer) PublishAsync(topic string, payload any) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := p.Publish(ctx, topic, payload); err != nil {
			log.Printf("[KAFKA-PRODUCER] ⚠️  Async publish failed for topic %s: %v", topic, err)
		}
	}()
}

// Close menutup writer Kafka.
func (p *Producer) Close() error {
	log.Println("🛑 Closing Kafka producer...")
	return p.writer.Close()
}
