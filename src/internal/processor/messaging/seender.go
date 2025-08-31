package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/contracts"
	"github.com/alielmi98/image-processing-service/pkg/rabbitmq"
)

type MessageSender struct {
	broker *rabbitmq.RabbitMQBroker
	ctx    context.Context
	cancel context.CancelFunc
	config *config.RabbitMQConfig
}

func NewMessageSender(config *config.Config, broker *rabbitmq.RabbitMQBroker) *MessageSender {
	ctx, cancel := context.WithCancel(context.Background())
	return  &MessageSender{
		broker: broker,
		ctx:    ctx,
		cancel: cancel,
		config: &config.RabbitMQ,
	}
}

func (ms *MessageSender) SendMessage(ctx context.Context, message *contracts.ProcessingResult) error {
	// Marshal message to JSON
	messageBody, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal processing message: %w", err)
	}

	// Create RabbitMQ message
	rabbitMsg := &rabbitmq.Message{
		ID:         fmt.Sprintf("job_%d", message.JobId),
		Topic:      "image.result",
		RoutingKey: ms.config.ResultRoutingKey,
		Body:       messageBody,
		Headers: map[string]interface{}{
			"content_type": "application/json",
			"job_id":       message.JobId,
			"image_id":     message.ImageId,
			"user_id":      message.UserId,
		},
		Priority:   5, // Default priority
		Timestamp:  time.Now(),
		RetryCount: 0,
		MaxRetries: 3,
	}

	// Publish message
	err = ms.broker.Publish(ctx, rabbitMsg)
	if err != nil {
		return fmt.Errorf("failed to publish result message: %w", err)
	}

	log.Printf("Sent result message for job %d", message.JobId)
	return nil
}

func (ms *MessageSender) Close() error {
	ms.cancel()
	if ms.broker != nil {
		return ms.broker.Close()
	}
	return nil
}
