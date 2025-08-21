package repository

import (
	"context"

	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
)

// ImageRepository defines the contract for image data operations
type ImageRepository interface {
	CreateImage(ctx context.Context, image models.Image) (models.Image, error)
	UpdateImage(ctx context.Context, id int, image map[string]interface{}) (models.Image, error)
	DeleteImage(ctx context.Context, id int) error
	GetImageByID(ctx context.Context, id int) (models.Image, error)
}

// ProcessingRepository defines the contract for processing job data operations
type ProcessingRepository interface {
	CreateProcessingJob(ctx context.Context, job models.ProcessingJob) (models.ProcessingJob, error)
	UpdateProcessingJob(ctx context.Context, id int, job map[string]interface{}) (models.ProcessingJob, error)
	DeleteProcessingJob(ctx context.Context, id int) error
	GetProcessingJobByID(ctx context.Context, id int) (models.ProcessingJob, error)
}
