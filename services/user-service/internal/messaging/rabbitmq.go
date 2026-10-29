package messaging

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}

func NewRabbitMQ(cfg *config.RabbitMQ) (*RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	connection, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQ{Connection: connection, Channel: channel}, nil
}

func (r *RabbitMQ) Close() error {
	return r.Connection.Close()
}
