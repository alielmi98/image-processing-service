package di

import (
	"fmt"

	contractAuth "github.com/alielmi98/image-processing-service/internal/auth/domain/auth"
	contractAuthRepo "github.com/alielmi98/image-processing-service/internal/auth/domain/repository"
	infraAuth "github.com/alielmi98/image-processing-service/internal/auth/infra/auth"
	infraAuthRepo "github.com/alielmi98/image-processing-service/internal/auth/infra/repository"
	contractMessaging "github.com/alielmi98/image-processing-service/internal/image/domain/messaging"
	contractImageRepo "github.com/alielmi98/image-processing-service/internal/image/domain/repository"
	infraMessaging "github.com/alielmi98/image-processing-service/internal/image/infra/messaging"
	infraImageRepo "github.com/alielmi98/image-processing-service/internal/image/infra/repository"
	"github.com/alielmi98/image-processing-service/internal/image/usecase"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/db"
	"github.com/alielmi98/image-processing-service/pkg/rabbitmq"
)

// midedlewares
func GetTokenProvider(cfg *config.Config) contractAuth.TokenProvider {
	return infraAuth.NewJwtProvider(cfg)
}

func GetUserRepository(cfg *config.Config) contractAuthRepo.UserRepository {
	return infraAuthRepo.NewUserPgRepo()
}

func GetImageRepository(cfg *config.Config) contractImageRepo.ImageRepository {
	var preloads []db.PreloadEntity = []db.PreloadEntity{{Entity: "ProcessingJobs"}}

	return infraImageRepo.NewImagePgRepository(cfg, preloads)
}

func GetProcessingRepository(cfg *config.Config) contractImageRepo.ProcessingRepository {
	var preloads []db.PreloadEntity = []db.PreloadEntity{{Entity: "Image"}}
	return infraImageRepo.NewProcessingRepository(cfg, preloads)
}

func GetMessageSender(cfg *config.Config, broker *rabbitmq.RabbitMQBroker) contractMessaging.MessageSender {
	return infraMessaging.NewMessageSender(cfg, broker)
}

func GetResultConsumer(cfg *config.Config, broker *rabbitmq.RabbitMQBroker, handler contractMessaging.ConsumerHandler) contractMessaging.MessageConsumer {
	return infraMessaging.NewMessageConsumer(broker, handler)
}
func GetRabbitMQBroker(cfg *config.Config) *rabbitmq.RabbitMQBroker {
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
	return rabbitmq.NewRabbitMQBroker(rbConfig)
}

// usecase

func GetImageUsecase(cfg *config.Config) *usecase.ImageUsecase {
	return usecase.NewImageUsecase(cfg, GetImageRepository(cfg))
}

func GetProcessingUsecase(cfg *config.Config, msSender contractMessaging.MessageSender) *usecase.ProcessingUsecase {
	return usecase.NewProcessingUseCase(cfg, GetProcessingRepository(cfg), msSender)
}
