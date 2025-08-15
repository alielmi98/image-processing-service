package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alielmi98/image-processing-service/internal/image/infra/messaging"
	"github.com/alielmi98/image-processing-service/internal/image/infra/processor"
	"github.com/alielmi98/image-processing-service/pkg/config"
)

type ImageProcessingService struct {
	rabbitMQClient *messaging.RabbitMQClient
	processor      *processor.ImageProcessorImpl
	ctx            context.Context
	cancel         context.CancelFunc
}

func NewImageProcessingService(cfg *config.Config, outputDir string) *ImageProcessingService {
	ctx, cancel := context.WithCancel(context.Background())

	// Create image processor
	imageProcessor := processor.NewImageProcessor(outputDir)

	// Convert config to messaging config
	rabbitMQConfig := &messaging.RabbitMQConfig{
		URL:                  fmt.Sprintf("amqp://%s:%s@%s:%s%s", cfg.RabbitMQ.User, cfg.RabbitMQ.Password, cfg.RabbitMQ.Host, cfg.RabbitMQ.Port, cfg.RabbitMQ.VHost),
		ProcessingQueue:      getEnvOrDefault("RABBITMQ_PROCESSING_QUEUE", "image.processing"),
		ResultQueue:          getEnvOrDefault("RABBITMQ_RESULT_QUEUE", "image.results"),
		ProcessingExchange:   getEnvOrDefault("RABBITMQ_PROCESSING_EXCHANGE", "image.processing.exchange"),
		ResultExchange:       getEnvOrDefault("RABBITMQ_RESULT_EXCHANGE", "image.results.exchange"),
		ProcessingRoutingKey: getEnvOrDefault("RABBITMQ_PROCESSING_ROUTING_KEY", cfg.RabbitMQ.ProcessingRoutingKey),
		ResultRoutingKey:     getEnvOrDefault("RABBITMQ_RESULT_ROUTING_KEY", cfg.RabbitMQ.ResultRoutingKey),
		PrefetchCount:        getEnvOrDefaultInt("RABBITMQ_PREFETCH_COUNT", cfg.RabbitMQ.PrefetchCount),
		ReconnectDelay:       time.Duration(getEnvOrDefaultInt("RABBITMQ_RECONNECT_DELAY_SECONDS", int(cfg.RabbitMQ.ReconnectDelay))) * time.Second,
		MaxReconnectAttempts: getEnvOrDefaultInt("RABBITMQ_MAX_RECONNECT_ATTEMPTS", cfg.RabbitMQ.MaxReconnectAttempts),
	}

	// Create RabbitMQ client
	rabbitMQClient := messaging.NewRabbitMQClient(rabbitMQConfig, imageProcessor)

	return &ImageProcessingService{
		rabbitMQClient: rabbitMQClient,
		processor:      imageProcessor,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (s *ImageProcessingService) Start() error {
	log.Println("Starting Image Processing Service...")

	// Connect to RabbitMQ
	err := s.rabbitMQClient.Connect()
	if err != nil {
		return err
	}

	// Start consuming messages
	err = s.rabbitMQClient.StartConsumer()
	if err != nil {
		return err
	}

	log.Println("Image Processing Service started successfully")

	// Wait for interrupt signal to gracefully shutdown
	s.waitForShutdown()

	return nil
}

func (s *ImageProcessingService) Stop() error {
	log.Println("Stopping Image Processing Service...")

	s.cancel()

	err := s.rabbitMQClient.Close()
	if err != nil {
		log.Printf("Error closing RabbitMQ client: %v", err)
	}

	log.Println("Image Processing Service stopped")
	return nil
}

func (s *ImageProcessingService) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down gracefully...", sig)
		s.Stop()
	case <-s.ctx.Done():
		log.Println("Context cancelled, shutting down...")
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		// Simple conversion, in production you might want proper error handling
		if intValue := parseInt(value); intValue > 0 {
			return intValue
		}
	}
	return defaultValue
}

// Simple integer parsing function
func parseInt(s string) int {
	result := 0
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0 // Invalid number
		}
	}
	return result
}
