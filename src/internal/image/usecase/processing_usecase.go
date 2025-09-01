package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/alielmi98/image-processing-service/common"
	"github.com/alielmi98/image-processing-service/constants"
	"github.com/alielmi98/image-processing-service/internal/image/domain/messaging"
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/alielmi98/image-processing-service/internal/image/domain/repository"
	"github.com/alielmi98/image-processing-service/internal/image/usecase/dto"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/contracts"
)

// Ensure ProcessingUsecase implements the ConsumerHandler interface
var _ messaging.ConsumerHandler = (*ProcessingUsecase)(nil)

type ProcessingUsecase struct {
	cfg       *config.Config
	repo      repository.ProcessingRepository
	messaging messaging.MessageSender
}

func NewProcessingUseCase(cfg *config.Config, repo repository.ProcessingRepository, messaging messaging.MessageSender) *ProcessingUsecase {
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
	job, err := uc.repo.CreateProcessingJob(ctx, entity)
	if err != nil {
		return dto.ProcessingResponse{}, err
	}
	processingJob, err := uc.repo.GetProcessingJobByID(ctx, job.Id)
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
	message := &contracts.ProcessingMessage{
		JobId:          job.Id,
		ImageId:        job.ImageId,
		ProcessingType: messaging.ToContractProcessingType(job.ProcessingType),
		Parameters:     job.Parameters,
		UserId:         userId,
		SourcePath:     "../" + job.Image.FilePath + "/" + job.Image.FileName,
		Priority:       1,
		Timestamp:      time.Now(),
		RetryCount:     0,
		MaxRetries:     3,
	}

	// Other fields as necessary

	// Send message to processor via messaging system
	return uc.messaging.SendMessage(ctx, message)
}

func (uc *ProcessingUsecase) GetProcessingJobByID(ctx context.Context, id int) (models.ProcessingJob, error) {
	return uc.repo.GetProcessingJobByID(ctx, id)
}

func (uc *ProcessingUsecase) HandleProcessingResult(ctx context.Context, result *contracts.ProcessingResult) error {
	// Map contract to domain model
	entity := models.ProcessingJob{
		Status:       messaging.FromContractImageStatus(result.Status),
		ResultPath:   sql.NullString{String: result.ResultPath, Valid: result.ResultPath != ""},
		ErrorMessage: sql.NullString{String: result.ErrorMessage, Valid: result.ErrorMessage != ""},
		StartedAt:    sql.NullTime{Time: result.StartedAt, Valid: result.StartedAt != time.Time{}},
		CompletedAt:  sql.NullTime{Time: result.CompletedAt, Valid: result.CompletedAt != time.Time{}},
		Duration:     sql.NullInt64{Int64: result.Duration, Valid: result.Duration != 0},
		ModifiedAt:   sql.NullTime{Time: time.Now().UTC(), Valid: true},
	}
	id := result.JobId
	// Call repository to update image processing job
	_, err := uc.repo.UpdateProcessingJob(ctx, id, entity)
	if err != nil {
		return err
	}
	return nil
}
