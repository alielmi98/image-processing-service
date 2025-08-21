package rabbitmq

import (
	"context"
	"time"
)

// Message represents a generic message structure
type Message struct {
	ID          string                 `json:"id"`
	Topic       string                 `json:"topic"`
	RoutingKey  string                 `json:"routing_key"`
	Body        []byte                 `json:"body"`
	Headers     map[string]interface{} `json:"headers"`
	Priority    uint8                  `json:"priority"`
	Timestamp   time.Time              `json:"timestamp"`
	RetryCount  int                    `json:"retry_count"`
	MaxRetries  int                    `json:"max_retries"`
}

// MessageHandler defines the function signature for message handlers
type MessageHandler func(ctx context.Context, message *Message) error

// Config represents RabbitMQ configuration
type Config struct {
	URL                  string        `mapstructure:"url"`
	Host                 string        `mapstructure:"host"`
	Port                 string        `mapstructure:"port"`
	Username             string        `mapstructure:"username"`
	Password             string        `mapstructure:"password"`
	VHost                string        `mapstructure:"vhost"`
	PrefetchCount        int           `mapstructure:"prefetch_count"`
	ReconnectDelay       time.Duration `mapstructure:"reconnect_delay"`
	MaxReconnectAttempts int           `mapstructure:"max_reconnect_attempts"`
}

// Publisher defines the interface for publishing messages
type Publisher interface {
	Publish(ctx context.Context, message *Message) error
	PublishBatch(ctx context.Context, messages []*Message) error
	Close() error
}

// Consumer defines the interface for consuming messages
type Consumer interface {
	Subscribe(topic string, handler MessageHandler) error
	Unsubscribe(topic string) error
	Start(ctx context.Context) error
	Stop() error
}

// RabbitMQ defines the main RabbitMQ interface
type RabbitMQ interface {
	Publisher
	Consumer
	Connect() error
	IsConnected() bool
	Health() error
}
