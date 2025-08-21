package di

import (
	"log"

	contractAuth "github.com/alielmi98/image-processing-service/internal/auth/domain/auth"
	contractAuthRepo "github.com/alielmi98/image-processing-service/internal/auth/domain/repository"
	infraAuth "github.com/alielmi98/image-processing-service/internal/auth/infra/auth"
	infraAuthRepo "github.com/alielmi98/image-processing-service/internal/auth/infra/repository"
	contractImageRepo "github.com/alielmi98/image-processing-service/internal/image/domain/repository"
	"github.com/alielmi98/image-processing-service/internal/image/infra/messaging"
	infraImageRepo "github.com/alielmi98/image-processing-service/internal/image/infra/repository"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/db"
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
	var preloads []db.PreloadEntity = []db.PreloadEntity{}
	return infraImageRepo.NewProcessingRepository(cfg, preloads)
}

func GetMessageSender(cfg *config.Config) *messaging.MessageSender {
	messageSender, err := messaging.NewMessageSender(cfg)
	if err != nil {
		log.Fatalf("failed to create message sender: %v", err)
	}
	return messageSender
}
