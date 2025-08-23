package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/alielmi98/image-processing-service/internal/image/entity"
	"github.com/alielmi98/image-processing-service/internal/processor/domain"
	"github.com/alielmi98/image-processing-service/pkg/rabbitmq"
)

// MessageConsumer handles RabbitMQ message consumption for image processing
type MessageConsumer struct {
	broker  *rabbitmq.RabbitMQBroker
	service domain.ProcessorService
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewMessageConsumer creates a new RabbitMQ message consumer
func NewMessageConsumer(broker *rabbitmq.RabbitMQBroker, service domain.ProcessorService) *MessageConsumer {
	ctx, cancel := context.WithCancel(context.Background())

	return &MessageConsumer{
		broker:  broker,
		service: service,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start starts consuming messages from RabbitMQ
func (c *MessageConsumer) Start(topic string) error {
	// Connect to RabbitMQ if not already connected
	if !c.broker.IsConnected() {
		if err := c.broker.Connect(); err != nil {
			return err
		}
	}

	// Subscribe to the topic
	err := c.broker.Subscribe(topic, func(ctx context.Context, msg *rabbitmq.Message) error {
		var processingMsg entity.ProcessingMessage
		if err := json.Unmarshal(msg.Body, &processingMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return err
		}

		// Process the image
		if err := c.service.ProcessImage(processingMsg); err != nil {
			log.Printf("Error processing image: %v", err)
			return err
		}

		log.Printf("Successfully processed image %d for user %d", 
			processingMsg.ImageId, processingMsg.UserId)
		return nil
	})

	if err != nil {
		return err
	}

	// Start consuming messages
	return c.broker.Start(c.ctx)
}

// Stop stops the message consumer
func (c *MessageConsumer) Stop() error {
	c.cancel()
	return c.broker.Stop()
}
