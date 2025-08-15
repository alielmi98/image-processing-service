package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type RabbitMQBroker struct {
	config    *Config
	conn      *amqp.Connection
	channel   *amqp.Channel
	handlers  map[string]MessageHandler
	mu        sync.RWMutex
	connected bool
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewRabbitMQBroker creates a new RabbitMQ broker instance
func NewRabbitMQBroker(config *Config) *RabbitMQBroker {
	ctx, cancel := context.WithCancel(context.Background())
	return &RabbitMQBroker{
		config:   config,
		handlers: make(map[string]MessageHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Connect establishes connection to RabbitMQ
func (r *RabbitMQBroker) Connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.connected {
		return nil
	}

	var connectionURL string
	if r.config.URL != "" {
		connectionURL = r.config.URL
	} else {
		connectionURL = fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
			r.config.Username,
			r.config.Password,
			r.config.Host,
			r.config.Port,
			r.config.VHost,
		)
	}

	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS if prefetch count is configured
	if r.config.PrefetchCount > 0 {
		err = channel.Qos(r.config.PrefetchCount, 0, false)
		if err != nil {
			channel.Close()
			conn.Close()
			return fmt.Errorf("failed to set QoS: %w", err)
		}
	}

	r.conn = conn
	r.channel = channel
	r.connected = true

	// Start connection monitoring
	go r.monitorConnection()

	log.Println("Successfully connected to RabbitMQ")
	return nil
}

// IsConnected returns the connection status
func (r *RabbitMQBroker) IsConnected() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.connected && r.conn != nil && !r.conn.IsClosed()
}

// Health checks the broker health
func (r *RabbitMQBroker) Health() error {
	if !r.IsConnected() {
		return fmt.Errorf("RabbitMQ connection is not healthy")
	}
	return nil
}

// Publish publishes a message to RabbitMQ
func (r *RabbitMQBroker) Publish(ctx context.Context, message *Message) error {
	if !r.IsConnected() {
		if err := r.Connect(); err != nil {
			return fmt.Errorf("failed to connect before publishing: %w", err)
		}
	}

	r.mu.RLock()
	channel := r.channel
	r.mu.RUnlock()

	if channel == nil {
		return fmt.Errorf("channel is not available")
	}

	// Declare exchange if needed
	err := channel.ExchangeDeclare(
		message.Topic, // exchange name
		"topic",       // exchange type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Convert headers
	headers := make(amqp.Table)
	for k, v := range message.Headers {
		headers[k] = v
	}

	// Publish message
	err = channel.Publish(
		message.Topic,      // exchange
		message.RoutingKey, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			MessageId:   message.ID,
			ContentType: "application/json",
			Body:        message.Body,
			Headers:     headers,
			Priority:    message.Priority,
			Timestamp:   message.Timestamp,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// PublishBatch publishes multiple messages
func (r *RabbitMQBroker) PublishBatch(ctx context.Context, messages []*Message) error {
	for _, msg := range messages {
		if err := r.Publish(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

// Subscribe subscribes to a topic with a handler
func (r *RabbitMQBroker) Subscribe(topic string, handler MessageHandler) error {
	r.mu.Lock()
	r.handlers[topic] = handler
	r.mu.Unlock()

	if !r.IsConnected() {
		if err := r.Connect(); err != nil {
			return fmt.Errorf("failed to connect before subscribing: %w", err)
		}
	}

	r.mu.RLock()
	channel := r.channel
	r.mu.RUnlock()

	if channel == nil {
		return fmt.Errorf("channel is not available")
	}

	// Declare exchange
	err := channel.ExchangeDeclare(
		topic,   // exchange name
		"topic", // exchange type
		true,    // durable
		false,   // auto-deleted
		false,   // internal
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue
	queue, err := channel.QueueDeclare(
		fmt.Sprintf("%s_queue", topic), // queue name
		true,                           // durable
		false,                          // delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil,                            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	err = channel.QueueBind(
		queue.Name, // queue name
		"#",        // routing key (wildcard for topic exchange)
		topic,      // exchange
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

// Unsubscribe removes a subscription
func (r *RabbitMQBroker) Unsubscribe(topic string) error {
	r.mu.Lock()
	delete(r.handlers, topic)
	r.mu.Unlock()
	return nil
}

// Start starts consuming messages
func (r *RabbitMQBroker) Start(ctx context.Context) error {
	if !r.IsConnected() {
		if err := r.Connect(); err != nil {
			return fmt.Errorf("failed to connect before starting: %w", err)
		}
	}

	r.mu.RLock()
	handlers := make(map[string]MessageHandler)
	for k, v := range r.handlers {
		handlers[k] = v
	}
	channel := r.channel
	r.mu.RUnlock()

	if channel == nil {
		return fmt.Errorf("channel is not available")
	}

	// Start consuming for each subscribed topic
	for topic, handler := range handlers {
		go r.consume(ctx, topic, handler)
	}

	return nil
}

// consume handles message consumption for a specific topic
func (r *RabbitMQBroker) consume(ctx context.Context, topic string, handler MessageHandler) {
	r.mu.RLock()
	channel := r.channel
	r.mu.RUnlock()

	if channel == nil {
		log.Printf("Channel not available for topic %s", topic)
		return
	}

	queueName := fmt.Sprintf("%s_queue", topic)
	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Printf("Failed to register consumer for topic %s: %v", topic, err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.ctx.Done():
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Printf("Message channel closed for topic %s", topic)
				return
			}

			// Convert AMQP message to our message format
			headers := make(map[string]interface{})
			for k, v := range msg.Headers {
				headers[k] = v
			}

			rabbitMsg := &Message{
				ID:         msg.MessageId,
				Topic:      topic,
				RoutingKey: msg.RoutingKey,
				Body:       msg.Body,
				Headers:    headers,
				Priority:   msg.Priority,
				Timestamp:  msg.Timestamp,
			}

			// Handle message
			if err := handler(ctx, rabbitMsg); err != nil {
				log.Printf("Error handling message: %v", err)
				msg.Nack(false, true) // Requeue on error
			} else {
				msg.Ack(false)
			}
		}
	}
}

// Stop stops the broker
func (r *RabbitMQBroker) Stop() error {
	r.cancel()

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel != nil {
		r.channel.Close()
		r.channel = nil
	}

	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}

	r.connected = false
	log.Println("RabbitMQ broker stopped")
	return nil
}

// Close closes the broker connection
func (r *RabbitMQBroker) Close() error {
	return r.Stop()
}

// monitorConnection monitors the connection and reconnects if needed
func (r *RabbitMQBroker) monitorConnection() {
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-time.After(30 * time.Second):
			if !r.IsConnected() {
				log.Println("Connection lost, attempting to reconnect...")
				if err := r.reconnect(); err != nil {
					log.Printf("Reconnection failed: %v", err)
				}
			}
		}
	}
}

// reconnect attempts to reconnect to RabbitMQ
func (r *RabbitMQBroker) reconnect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Close existing connections
	if r.channel != nil {
		r.channel.Close()
		r.channel = nil
	}
	if r.conn != nil {
		r.conn.Close()
		r.conn = nil
	}
	r.connected = false

	// Wait before reconnecting
	time.Sleep(r.config.ReconnectDelay)

	// Attempt to reconnect
	return r.Connect()
}
