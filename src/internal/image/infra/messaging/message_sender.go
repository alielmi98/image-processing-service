package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/alielmi98/image-processing-service/internal/image/entity"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/rabbitmq"
)

type MessageSender struct {
	config *config.RabbitMQConfig
	broker *rabbitmq.RabbitMQBroker
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMessageSender(config *config.Config) (*MessageSender, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Build connection URL
	connectionURL := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		config.RabbitMQ.User,
		config.RabbitMQ.Password,
		config.RabbitMQ.Host,
		config.RabbitMQ.Port,
		config.RabbitMQ.VHost,
	)

	// Convert to rabbitmq.Config
	rbConfig := &rabbitmq.Config{
		URL:                  connectionURL,
		Host:                 config.RabbitMQ.Host,
		Port:                 config.RabbitMQ.Port,
		Username:             config.RabbitMQ.User,
		Password:             config.RabbitMQ.Password,
		VHost:                config.RabbitMQ.VHost,
		PrefetchCount:        config.RabbitMQ.PrefetchCount,
		ReconnectDelay:       config.RabbitMQ.ReconnectDelay,
		MaxReconnectAttempts: config.RabbitMQ.MaxReconnectAttempts,
	}

	broker := rabbitmq.NewRabbitMQBroker(rbConfig)

	client := &MessageSender{
		config: &config.RabbitMQ,
		broker: broker,
		ctx:    ctx,
		cancel: cancel,
	}

	// Connect to RabbitMQ
	if err := broker.Connect(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	dummyHandler := func(ctx context.Context, msg *rabbitmq.Message) error {
		return nil
	}
	if err := broker.Subscribe("image.processing", dummyHandler); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to subscribe and create queue: %w", err)
	}

	return client, nil
}

func (ms *MessageSender) SendMessage(ctx context.Context, message *entity.ProcessingMessage) error {
	// Marshal message to JSON
	messageBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal processing message: %w", err)
	}

	// Create RabbitMQ message
	rabbitMsg := &rabbitmq.Message{
		ID:         fmt.Sprintf("job_%d", message.JobId),
		Topic:      "image.processing",
		RoutingKey: ms.config.ProcessingRoutingKey,
		Body:       messageBody,
		Headers: map[string]interface{}{
			"content_type": "application/json",
			"job_id":       message.JobId,
			"image_id":     message.ImageId,
			"user_id":      message.UserId,
		},
		Priority:   5, // Default priority
		Timestamp:  message.Timestamp,
		RetryCount: message.RetryCount,
		MaxRetries: message.MaxRetries,
	}

	// Publish message
	err = ms.broker.Publish(ctx, rabbitMsg)
	if err != nil {
		return fmt.Errorf("failed to publish processing message: %w", err)
	}

	log.Printf("Sent processing message for job %d", message.JobId)
	return nil
}

func (ms *MessageSender) Close() error {
	ms.cancel()
	if ms.broker != nil {
		return ms.broker.Close()
	}
	return nil
}
