package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/alielmi98/image-processing-service/internal/image/domain/messaging"
	"github.com/alielmi98/image-processing-service/pkg/contracts"
	"github.com/alielmi98/image-processing-service/pkg/rabbitmq"
)

// messageConsumer implements MessageConsumer
type messageConsumer struct {
	broker  *rabbitmq.RabbitMQBroker
	handler messaging.ConsumerHandler
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewMessageConsumer creates a new message consumer
func NewMessageConsumer(broker *rabbitmq.RabbitMQBroker, handler messaging.ConsumerHandler) messaging.MessageConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &messageConsumer{
		broker:  broker,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start starts consuming messages from RabbitMQ
func (c *messageConsumer) Start(topic string) error {
	// Connect to RabbitMQ if not already connected
	if !c.broker.IsConnected() {
		if err := c.broker.Connect(); err != nil {
			return err
		}
	}

	// Subscribe to the topic
	err := c.broker.Subscribe(topic, func(ctx context.Context, msg *rabbitmq.Message) error {
		var processingResultMsg contracts.ProcessingResult
		if err := json.Unmarshal(msg.Body, &processingResultMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			return err
		}

		err := c.handler.HandleProcessingResult(ctx, &processingResultMsg)
		if err != nil {
			log.Printf("Error Update ProcessingJob: %v", err)
			return err
		}
		log.Printf("Successfully Update ProcessingJob %d for user %d",
			processingResultMsg.JobId, processingResultMsg.UserId)
		return nil
	})

	if err != nil {
		return err
	}

	// Start consuming messages
	return c.broker.Start(c.ctx)
}

// Stop stops the message consumer
func (c *messageConsumer) Stop() error {
	c.cancel()
	return c.broker.Stop()
}
