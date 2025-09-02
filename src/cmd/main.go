package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alielmi98/image-processing-service/constants"
	"github.com/alielmi98/image-processing-service/di"
	"github.com/alielmi98/image-processing-service/docs"
	authRouter "github.com/alielmi98/image-processing-service/internal/auth/api/routers"
	"github.com/alielmi98/image-processing-service/internal/image/api/handlers"
	imageRouter "github.com/alielmi98/image-processing-service/internal/image/api/routers"
	"github.com/alielmi98/image-processing-service/internal/image/domain/messaging"
	"github.com/alielmi98/image-processing-service/internal/middlewares"
	migration "github.com/alielmi98/image-processing-service/migrations"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/db"
	"github.com/alielmi98/image-processing-service/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()

	err := db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		log.Fatalf("caller:%s  Level:%s  Msg:%s", constants.Postgres, constants.Startup, err.Error())
	}

	//messaging
	broker := di.GetRabbitMQBroker(cfg)
	messageSender := di.GetMessageSender(cfg, broker)
	// usecase
	imageUsecase := di.GetImageUsecase(cfg)
	processorUsecase := di.GetProcessingUsecase(cfg, messageSender)
	imageHandler := handlers.NewImageHandler(cfg, imageUsecase)
	processorHandler := handlers.NewProcessingHandler(cfg, processorUsecase)

	// Migrate the database
	migration.Up1()

	// Start the message consumer in a goroutine
	go func() {
		if err := InitMessageConsumer(cfg, broker, messaging.ConsumerHandler(processorUsecase)); err != nil {
			log.Fatalf("Failed to start message consumer: %v", err)
		}
	}()

	InitServer(cfg, imageHandler, processorHandler)

}

func InitServer(cfg *config.Config, imageHandler *handlers.ImageHandler, processorHandler *handlers.ProcessingHandler) {
	r := gin.New()

	r.Use(middlewares.Cors(cfg))
	RegisterRoutes(r, cfg, imageHandler, processorHandler)
	RegisterSwagger(r, cfg)
	log.Printf("Caller:%s Level:%s Msg:%s", constants.General, constants.Startup, "Started")
	r.Run(fmt.Sprintf(":%s", cfg.Server.InternalPort))

}

func RegisterRoutes(r *gin.Engine, cfg *config.Config, imageHandler *handlers.ImageHandler, processorHandler *handlers.ProcessingHandler) {
	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		//Auth
		auth := v1.Group("/auth")
		authRouter.Auth(auth, cfg)
		//Image
		tokenProvider := di.GetTokenProvider(cfg)
		image := v1.Group("/images")
		image.Use(middlewares.Authentication(cfg, tokenProvider))
		imageRouter.Image(image, cfg, imageHandler)

		//Processing
		processing := v1.Group("/processing")
		processing.Use(middlewares.Authentication(cfg, tokenProvider))
		imageRouter.Processing(processing, cfg, processorHandler)

	}
}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "golang web api"
	docs.SwaggerInfo.Description = "This is a sample server for golang web api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.InternalPort)
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func InitMessageConsumer(cfg *config.Config, broker *rabbitmq.RabbitMQBroker, handler messaging.ConsumerHandler) error {
	consumer := di.GetResultConsumer(cfg, broker, handler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("Starting image processor consumer on queue: %s", "image.result")
	if err := consumer.Start("image.result"); err != nil {
		return fmt.Errorf("error starting consumer: %w", err)
	}

	// Handle graceful shutdown
	go func() {
		<-sigChan
		log.Println("Shutting down message consumer...")
		if err := consumer.Stop(); err != nil {
			log.Printf("Error stopping consumer: %v", err)
		}
	}()

	return nil
}
