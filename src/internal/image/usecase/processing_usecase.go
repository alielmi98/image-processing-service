package usecase

import (
	"context"
	"time"

	"github.com/alielmi98/image-processing-service/common"
	"github.com/alielmi98/image-processing-service/constants"
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/alielmi98/image-processing-service/internal/image/domain/repository"
	"github.com/alielmi98/image-processing-service/internal/image/entity"
	"github.com/alielmi98/image-processing-service/internal/image/infra/messaging"
	"github.com/alielmi98/image-processing-service/internal/image/usecase/dto"
	"github.com/alielmi98/image-processing-service/pkg/config"
)

type ProcessingUsecase struct {
	cfg       *config.Config
	repo      repository.ProcessingRepository
	messaging *messaging.MessageSender
}

func NewProcessingUseCase(cfg *config.Config, repo repository.ProcessingRepository, messaging *messaging.MessageSender) *ProcessingUsecase {
	return &ProcessingUsecase{
		cfg:       cfg,
		repo:      repo,
		messaging: messaging,
	}
}

func (uc *ProcessingUsecase) CreateProcessingJob(ctx context.Context, req dto.ProcessingRequest) (dto.ProcessingResponse, error) {
	// Map DTO to domain model
	entity, _ := common.TypeConverter[models.ProcessingJob](req)
	// Call repository to save image
	processingJob, err := uc.repo.CreateProcessingJob(ctx, entity)
	if err != nil {
		return dto.ProcessingResponse{}, err
	}
	err = uc.SendProcessingMessage(ctx, &processingJob)
	if err != nil {
		return dto.ProcessingResponse{}, err
	}
	// Map domain model to response DTO
	response := dto.ProcessingResponse{
		JobId: processingJob.Id,
	}
	return response, nil
}

func (uc *ProcessingUsecase) SendProcessingMessage(ctx context.Context, job *models.ProcessingJob) error {
	userId := int(ctx.Value(constants.UserIdKey).(float64))
	message := &entity.ProcessingMessage{
		JobId:          job.Id,
		ImageId:        job.ImageId,
		ProcessingType: job.ProcessingType,
		Parameters:     job.Parameters,
		UserId:         userId,
		SourcePath:     "/uploads",
		DestinationDir: "/uploads/processed",
		Priority:       1,
		Timestamp:      time.Now(),
		RetryCount:     0,
		MaxRetries:     3,
	}

	// Other fields as necessary

	// Send message to processor via messaging system
	return uc.messaging.SendMessage(ctx, message)
}

func (uc *ProcessingUsecase) HandleProcessingResult(ctx context.Context, result *entity.ProcessingResult) error {
	return nil
}
