package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alielmi98/image-processing-service/internal/processor"
	"github.com/alielmi98/image-processing-service/internal/processor/messaging"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/rabbitmq"
)

func main() {
	// Load configuration
	cfg := config.GetConfig()

	// Build connection URL
	connectionURL := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
		cfg.RabbitMQ.VHost,
	)

	// Convert to rabbitmq.Config
	rbConfig := &rabbitmq.Config{
		URL:                  connectionURL,
		Host:                 cfg.RabbitMQ.Host,
		Port:                 cfg.RabbitMQ.Port,
		Username:             cfg.RabbitMQ.User,
		Password:             cfg.RabbitMQ.Password,
		VHost:                cfg.RabbitMQ.VHost,
		PrefetchCount:        cfg.RabbitMQ.PrefetchCount,
		ReconnectDelay:       cfg.RabbitMQ.ReconnectDelay,
		MaxReconnectAttempts: cfg.RabbitMQ.MaxReconnectAttempts,
	}

	broker := rabbitmq.NewRabbitMQBroker(rbConfig)

	// Connect to RabbitMQ
	if err := broker.Connect(); err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", err)
	}

	// Initialize processor service
	service := processor.NewProcessor()

	// Create message consumer
	consumer := messaging.NewMessageConsumer(broker, service)

	// Create a context that cancels on interrupt signal
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the consumer in a goroutine
	go func() {
		log.Printf("Starting image processor consumer on queue: %s", "image.processing")
		if err := consumer.Start("image.processing"); err != nil {
			log.Fatalf("Error starting consumer: %v", err)
		}
	}()

	// Wait for interrupt signal
	sig := <-sigChan
	log.Printf("Received signal %v. Shutting down...", sig)

	// Stop the consumer
	if err := consumer.Stop(); err != nil {
		log.Printf("Error stopping consumer: %v", err)
	}

	log.Println("Image processor stopped gracefully")
}
