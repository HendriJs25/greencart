package kafka

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/config"

	kafkago "github.com/segmentio/kafka-go"
)

type producer struct {
	writer *kafkago.Writer
}

type Producer interface {
	Publish(ctx context.Context, topic string, key string, payload []byte) error
	Close() error
}

func NewProducer(cfg *config.Kafka) Producer {
	writer := &kafkago.Writer{
		Addr:                   kafkago.TCP(cfg.Brokers...),
		Balancer:               &kafkago.Hash{},
		RequiredAcks:           kafkago.RequireAll,
		Async:                  false,
		MaxAttempts:            3,
		BatchTimeout:           10 * time.Millisecond,
		AllowAutoTopicCreation: false,
	}
	return &producer{
		writer: writer,
	}
}

func (p *producer) Publish(ctx context.Context, topic string, key string, payload []byte) error {
	publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.writer.WriteMessages(publishCtx, kafkago.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	})

	if err != nil {
		return fmt.Errorf("failed to publish message to topic %s: %w", topic, err)
	}
	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
