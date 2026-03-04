// Package kafka menyediakan consumer base yang memudahkan pembuatan
// Kafka consumer group. Setiap handler domain akan di-wrap dalam Consumer.
package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// MessageHandler adalah fungsi callback yang dipanggil saat ada pesan masuk.
// key adalah message key (bisa kosong), value adalah raw JSON bytes.
type MessageHandler func(ctx context.Context, topic string, key, value []byte) error

// Consumer membungkus kafka.Reader dan menjalankan loop konsumsi.
type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
	name    string
}

// ConsumerConfig konfigurasi untuk membuat Consumer.
type ConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
	// MinBytes / MaxBytes opsional, default Kafka digunakan jika 0.
	MinBytes int
	MaxBytes int
}

// NewConsumer membuat Consumer baru.
//
// name    — nama deskriptif untuk log (contoh: "audit-consumer")
// cfg     — konfigurasi broker, topic, group
// handler — fungsi yang dipanggil untuk setiap pesan
func NewConsumer(name string, cfg ConsumerConfig, handler MessageHandler) *Consumer {
	minBytes := cfg.MinBytes
	if minBytes == 0 {
		minBytes = 1
	}
	maxBytes := cfg.MaxBytes
	if maxBytes == 0 {
		maxBytes = 10e6 // 10 MB
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		GroupID:     cfg.GroupID,
		MinBytes:    minBytes,
		MaxBytes:    maxBytes,
		StartOffset: kafka.LastOffset,
		Logger:      kafka.LoggerFunc(func(msg string, args ...interface{}) {}), // silent info
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("[KAFKA-CONSUMER:%s ERROR] "+msg, append([]interface{}{name}, args...)...)
		}),
	})

	log.Printf("✅ Kafka consumer [%s] initialized — topic: %s, group: %s, brokers: %v",
		name, cfg.Topic, cfg.GroupID, cfg.Brokers)

	return &Consumer{
		reader:  r,
		handler: handler,
		name:    name,
	}
}

// Start memulai loop konsumsi secara blocking. Jalankan dalam goroutine tersendiri.
// Loop berhenti saat ctx dibatalkan.
func (c *Consumer) Start(ctx context.Context) {
	log.Printf("[KAFKA-CONSUMER:%s] 🚀 Starting...", c.name)
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("[KAFKA-CONSUMER:%s] 🛑 Context cancelled, stopping.", c.name)
				return
			}
			log.Printf("[KAFKA-CONSUMER:%s] ⚠️  FetchMessage error: %v", c.name, err)
			continue
		}

		if err := c.handler(ctx, msg.Topic, msg.Key, msg.Value); err != nil {
			log.Printf("[KAFKA-CONSUMER:%s] ⚠️  Handler error (topic=%s): %v",
				c.name, msg.Topic, err)
			// Tetap commit agar tidak stuck pada pesan yang selalu gagal.
			// Di production, pertimbangkan dead-letter topic.
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("[KAFKA-CONSUMER:%s] ⚠️  CommitMessages error: %v", c.name, err)
		}
	}
}

// Close menutup reader Kafka.
func (c *Consumer) Close() error {
	log.Printf("[KAFKA-CONSUMER:%s] 🛑 Closing...", c.name)
	return c.reader.Close()
}

// ─── Multi-Topic Fan-Out Consumer ────────────────────────────────────────────

// MultiTopicConsumer mengonsumsi beberapa topic sekaligus dalam satu consumer group.
type MultiTopicConsumer struct {
	reader  *kafka.Reader
	handler MessageHandler
	name    string
}

// NewMultiTopicConsumer membuat consumer yang mengkonsumsi banyak topic sekaligus.
func NewMultiTopicConsumer(name string, brokers []string, topics []string, groupID string, handler MessageHandler) *MultiTopicConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupTopics: topics,
		GroupID:     groupID,
		MinBytes:    1,
		MaxBytes:    10e6,
		StartOffset: kafka.LastOffset,
		Logger:      kafka.LoggerFunc(func(msg string, args ...interface{}) {}),
		ErrorLogger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			log.Printf("[KAFKA-CONSUMER:%s ERROR] "+msg, append([]interface{}{name}, args...)...)
		}),
	})

	log.Printf("✅ Kafka multi-topic consumer [%s] initialized — topics: %v, group: %s, brokers: %v",
		name, topics, groupID, brokers)

	return &MultiTopicConsumer{reader: r, handler: handler, name: name}
}

// Start memulai loop konsumsi (blocking). Jalankan dalam goroutine.
func (c *MultiTopicConsumer) Start(ctx context.Context) {
	log.Printf("[KAFKA-CONSUMER:%s] 🚀 Starting multi-topic consumer...", c.name)
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Printf("[KAFKA-CONSUMER:%s] 🛑 Context cancelled, stopping.", c.name)
				return
			}
			log.Printf("[KAFKA-CONSUMER:%s] ⚠️  FetchMessage error: %v", c.name, err)
			continue
		}

		if err := c.handler(ctx, msg.Topic, msg.Key, msg.Value); err != nil {
			log.Printf("[KAFKA-CONSUMER:%s] ⚠️  Handler error (topic=%s): %v",
				c.name, msg.Topic, err)
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("[KAFKA-CONSUMER:%s] ⚠️  CommitMessages error: %v", c.name, err)
		}
	}
}

// Close menutup reader.
func (c *MultiTopicConsumer) Close() error {
	log.Printf("[KAFKA-CONSUMER:%s] 🛑 Closing...", c.name)
	return c.reader.Close()
}

// ─── Utility ─────────────────────────────────────────────────────────────────

// UnmarshalPayload helper untuk unmarshal JSON message value ke struct tujuan.
func UnmarshalPayload(value []byte, target any) error {
	return json.Unmarshal(value, target)
}
