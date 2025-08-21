package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alielmi98/image-processing-service/internal/image/entity"
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/alielmi98/image-processing-service/pkg/rabbitmq"
)

type RabbitMQConfig struct {
	URL                    string
	ProcessingQueue        string
	ResultQueue           string
	ProcessingExchange     string
	ResultExchange        string
	ProcessingRoutingKey   string
	ResultRoutingKey      string
	PrefetchCount         int
	ReconnectDelay        time.Duration
	MaxReconnectAttempts  int
}

type RabbitMQClient struct {
	config     *RabbitMQConfig
	broker     *rabbitmq.RabbitMQBroker
	processor  ImageProcessor
	ctx        context.Context
	cancel     context.CancelFunc
}

type ImageProcessor interface {
	ProcessImage(ctx context.Context, message *entity.ProcessingMessage) (*entity.ProcessingResult, error)
}

func NewRabbitMQClient(config *RabbitMQConfig, processor ImageProcessor) *RabbitMQClient {
	ctx, cancel := context.WithCancel(context.Background())
	
	// Convert to rabbitmq.Config
	rbConfig := &rabbitmq.Config{
		URL:                  config.URL,
		PrefetchCount:        config.PrefetchCount,
		ReconnectDelay:       config.ReconnectDelay,
		MaxReconnectAttempts: config.MaxReconnectAttempts,
	}
	
	broker := rabbitmq.NewRabbitMQBroker(rbConfig)
	
	return &RabbitMQClient{
		config:    config,
		broker:    broker,
		processor: processor,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (r *RabbitMQClient) Connect() error {
	err := r.broker.Connect()
	if err != nil {
		return err
	}
	
	// Only subscribe to processing queue if we have a processor
	if r.processor != nil {
		err = r.broker.Subscribe(r.config.ProcessingExchange, r.handleProcessingMessage)
		if err != nil {
			return fmt.Errorf("failed to subscribe to processing queue: %w", err)
		}
	}
	
	return nil
}

func (r *RabbitMQClient) StartConsumer() error {
	return r.broker.Start(r.ctx)
}

func (r *RabbitMQClient) handleProcessingMessage(ctx context.Context, message *rabbitmq.Message) error {
	var processingMsg entity.ProcessingMessage
	
	err := json.Unmarshal(message.Body, &processingMsg)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return err // Don't requeue malformed messages
	}

	log.Printf("Processing image job %d for user %d", processingMsg.JobId, processingMsg.UserId)

	// Process the image
	result, err := r.processor.ProcessImage(ctx, &processingMsg)
	if err != nil {
		log.Printf("Failed to process image job %d: %v", processingMsg.JobId, err)
		
		// Create error result
		result = &entity.ProcessingResult{
			JobId:        processingMsg.JobId,
			ImageId:      processingMsg.ImageId,
			UserId:       processingMsg.UserId,
			Status:       models.ImageStatusFailed,
			ErrorMessage: err.Error(),
			ProcessedAt:  time.Now(),
		}
	}

	// Send result back to RabbitMQ
	err = r.publishResult(result)
	if err != nil {
		log.Printf("Failed to publish result for job %d: %v", processingMsg.JobId, err)
		return err // Requeue the message
	}

	log.Printf("Successfully processed job %d", processingMsg.JobId)
	return nil
}

func (r *RabbitMQClient) publishResult(result *entity.ProcessingResult) error {
	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	message := &rabbitmq.Message{
		ID:         fmt.Sprintf("result-%d-%d", result.JobId, time.Now().Unix()),
		Topic:      r.config.ResultExchange,
		RoutingKey: r.config.ResultRoutingKey,
		Body:       body,
		Headers:    make(map[string]interface{}),
		Priority:   0,
		Timestamp:  time.Now(),
	}

	return r.broker.Publish(r.ctx, message)
}

func (r *RabbitMQClient) PublishProcessingJob(processingMsg *entity.ProcessingMessage) error {
	body, err := json.Marshal(processingMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	message := &rabbitmq.Message{
		ID:         fmt.Sprintf("job-%d-%d", processingMsg.JobId, time.Now().Unix()),
		Topic:      r.config.ProcessingExchange,
		RoutingKey: r.config.ProcessingRoutingKey,
		Body:       body,
		Headers:    make(map[string]interface{}),
		Priority:   uint8(processingMsg.Priority),
		Timestamp:  time.Now(),
	}

	return r.broker.Publish(r.ctx, message)
}

func (r *RabbitMQClient) Close() error {
	r.cancel()
	return r.broker.Close()
}
