package usecase

import (
	"context"

	"github.com/alielmi98/image-processing-service/common"
	"github.com/alielmi98/image-processing-service/constants"
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/alielmi98/image-processing-service/internal/image/domain/repository"
	"github.com/alielmi98/image-processing-service/internal/image/usecase/dto"
	"github.com/alielmi98/image-processing-service/pkg/config"
)

type ImageUsecase struct {
	cfg  *config.Config
	repo repository.ImageRepository
}

func NewImageUsecase(cfg *config.Config, repo repository.ImageRepository) *ImageUsecase {
	return &ImageUsecase{
		cfg:  cfg,
		repo: repo,
	}
}

// Create
func (uc *ImageUsecase) CreateImage(ctx context.Context, req dto.CreateImage) (dto.ImageResponse, error) {
	userId := int(ctx.Value(constants.UserIdKey).(float64))
	req.UserID = userId
	// Map DTO to domain model
	entity, _ := common.TypeConverter[models.Image](req)
	// Call repository to save image
	image, err := uc.repo.CreateImage(ctx, entity)
	if err != nil {
		return dto.ImageResponse{}, err
	}

	// Map domain model to response DTO
	response, _ := common.TypeConverter[dto.ImageResponse](image)
	return response, nil
}
