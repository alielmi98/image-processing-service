package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alielmi98/image-processing-service/internal/image/entity"
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/streadway/amqp"
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
	conn       *amqp.Connection
	channel    *amqp.Channel
	processor  ImageProcessor
	ctx        context.Context
	cancel     context.CancelFunc
}

type ImageProcessor interface {
	ProcessImage(ctx context.Context, message *entity.ImageProcessingMessage) (*entity.ImageProcessingResult, error)
}

func NewRabbitMQClient(config *RabbitMQConfig, processor ImageProcessor) *RabbitMQClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &RabbitMQClient{
		config:    config,
		processor: processor,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (r *RabbitMQClient) Connect() error {
	var err error
	r.conn, err = amqp.Dial(r.config.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS for fair dispatching
	err = r.channel.Qos(r.config.PrefetchCount, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Declare exchanges
	err = r.declareExchanges()
	if err != nil {
		return fmt.Errorf("failed to declare exchanges: %w", err)
	}

	// Declare queues
	err = r.declareQueues()
	if err != nil {
		return fmt.Errorf("failed to declare queues: %w", err)
	}

	return nil
}

func (r *RabbitMQClient) declareExchanges() error {
	// Declare processing exchange
	err := r.channel.ExchangeDeclare(
		r.config.ProcessingExchange,
		"direct",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare processing exchange: %w", err)
	}

	// Declare result exchange
	err = r.channel.ExchangeDeclare(
		r.config.ResultExchange,
		"direct",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare result exchange: %w", err)
	}

	return nil
}

func (r *RabbitMQClient) declareQueues() error {
	// Declare processing queue
	_, err := r.channel.QueueDeclare(
		r.config.ProcessingQueue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-message-ttl":             300000, // 5 minutes TTL
			"x-dead-letter-exchange":    r.config.ProcessingExchange + ".dlx",
			"x-dead-letter-routing-key": "failed",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to declare processing queue: %w", err)
	}

	// Bind processing queue to exchange
	err = r.channel.QueueBind(
		r.config.ProcessingQueue,
		r.config.ProcessingRoutingKey,
		r.config.ProcessingExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind processing queue: %w", err)
	}

	// Declare result queue
	_, err = r.channel.QueueDeclare(
		r.config.ResultQueue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare result queue: %w", err)
	}

	// Bind result queue to exchange
	err = r.channel.QueueBind(
		r.config.ResultQueue,
		r.config.ResultRoutingKey,
		r.config.ResultExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind result queue: %w", err)
	}

	return nil
}

func (r *RabbitMQClient) StartConsumer() error {
	msgs, err := r.channel.Consume(
		r.config.ProcessingQueue,
		"image-processor", // consumer tag
		false,            // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Image processing consumer started. Waiting for messages...")

	go func() {
		for {
			select {
			case <-r.ctx.Done():
				log.Println("Consumer context cancelled, stopping...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Message channel closed, attempting to reconnect...")
					r.handleReconnect()
					return
				}
				r.handleMessage(msg)
			}
		}
	}()

	return nil
}

func (r *RabbitMQClient) handleMessage(msg amqp.Delivery) {
	var processingMsg entity.ImageProcessingMessage
	
	err := json.Unmarshal(msg.Body, &processingMsg)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		msg.Nack(false, false) // Don't requeue malformed messages
		return
	}

	log.Printf("Processing image job %d for user %d", processingMsg.JobId, processingMsg.UserId)

	// Process the image
	result, err := r.processor.ProcessImage(r.ctx, &processingMsg)
	if err != nil {
		log.Printf("Failed to process image job %d: %v", processingMsg.JobId, err)
		
		// Create error result
		result = &entity.ImageProcessingResult{
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
		msg.Nack(false, true) // Requeue the message
		return
	}

	// Acknowledge the message
	err = msg.Ack(false)
	if err != nil {
		log.Printf("Failed to acknowledge message: %v", err)
	}

	log.Printf("Successfully processed and acknowledged job %d", processingMsg.JobId)
}

func (r *RabbitMQClient) publishResult(result *entity.ImageProcessingResult) error {
	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	err = r.channel.Publish(
		r.config.ResultExchange,
		r.config.ResultRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			MessageId:    fmt.Sprintf("result-%d-%d", result.JobId, time.Now().Unix()),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish result: %w", err)
	}

	return nil
}

func (r *RabbitMQClient) handleReconnect() {
	for attempt := 1; attempt <= r.config.MaxReconnectAttempts; attempt++ {
		log.Printf("Reconnection attempt %d/%d", attempt, r.config.MaxReconnectAttempts)
		
		time.Sleep(r.config.ReconnectDelay)
		
		err := r.Connect()
		if err != nil {
			log.Printf("Reconnection attempt %d failed: %v", attempt, err)
			continue
		}
		
		err = r.StartConsumer()
		if err != nil {
			log.Printf("Failed to restart consumer on attempt %d: %v", attempt, err)
			continue
		}
		
		log.Printf("Successfully reconnected on attempt %d", attempt)
		return
	}
	
	log.Printf("Failed to reconnect after %d attempts", r.config.MaxReconnectAttempts)
}

func (r *RabbitMQClient) Close() error {
	r.cancel()
	
	if r.channel != nil {
		r.channel.Close()
	}
	
	if r.conn != nil {
		r.conn.Close()
	}
	
	return nil
}

func (r *RabbitMQClient) PublishProcessingJob(message *entity.ImageProcessingMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = r.channel.Publish(
		r.config.ProcessingExchange,
		r.config.ProcessingRoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Priority:     uint8(message.Priority),
			Timestamp:    time.Now(),
			MessageId:    fmt.Sprintf("job-%d-%d", message.JobId, time.Now().Unix()),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
