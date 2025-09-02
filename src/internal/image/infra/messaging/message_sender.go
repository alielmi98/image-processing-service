package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/alielmi98/image-processing-service/internal/image/domain/messaging"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/contracts"
	"github.com/alielmi98/image-processing-service/pkg/rabbitmq"
)

type messageSender struct {
	config *config.RabbitMQConfig
	broker *rabbitmq.RabbitMQBroker
	ctx    context.Context
	cancel context.CancelFunc
}

func NewMessageSender(config *config.Config, broker *rabbitmq.RabbitMQBroker) messaging.MessageSender {
	ctx, cancel := context.WithCancel(context.Background())
	return &messageSender{
		config: &config.RabbitMQ,
		broker: broker,
		ctx:    ctx,
		cancel: cancel,
	}

}

func (ms *messageSender) SendMessage(ctx context.Context, message *contracts.ProcessingMessage) error {
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

func (ms *messageSender) Close() error {
	ms.cancel()
	if ms.broker != nil {
		return ms.broker.Close()
	}
	return nil
}
