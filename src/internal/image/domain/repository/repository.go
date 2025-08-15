package repository

import (
	"context"

	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
)

// ImageRepository defines the contract for image data operations
type ImageRepository interface {
	// Basic CRUD operations
	CreateImage(ctx context.Context, image *models.Image) error
	UpdateImage(ctx context.Context, id int, image *models.Image) error
	DeleteImage(ctx context.Context, id int) error
	GetImageByID(ctx context.Context, id int) (*models.Image, error)
	GetImagesByUserID(ctx context.Context, userID int, offset, limit int) ([]models.Image, error)
	
	// Image-specific operations
	GetImageByFileName(ctx context.Context, fileName string) (*models.Image, error)
	ExistsByFileName(ctx context.Context, fileName string) (bool, error)
	UpdateImageStatus(ctx context.Context, imageID int, status models.ImageStatus) error
	GetImagesByStatus(ctx context.Context, status models.ImageStatus, offset, limit int) ([]models.Image, error)
	CountImagesByUserID(ctx context.Context, userID int) (int64, error)
}

// ProcessingJobRepository defines the contract for processing job operations
type ProcessingJobRepository interface {
	// Basic CRUD operations
	CreateJob(ctx context.Context, job *models.ProcessingJob) error
	UpdateJob(ctx context.Context, id int, job *models.ProcessingJob) error
	DeleteJob(ctx context.Context, id int) error
	GetJobByID(ctx context.Context, id int) (*models.ProcessingJob, error)
	
	// Job-specific operations
	GetJobsByImageID(ctx context.Context, imageID int) ([]models.ProcessingJob, error)
	GetJobsByStatus(ctx context.Context, status models.ImageStatus, offset, limit int) ([]models.ProcessingJob, error)
	GetPendingJobs(ctx context.Context, limit int) ([]models.ProcessingJob, error)
	UpdateJobStatus(ctx context.Context, jobID int, status models.ImageStatus) error
	UpdateJobProgress(ctx context.Context, jobID int, startedAt *int64, completedAt *int64, duration *int64, errorMsg *string) error
	GetJobsByProcessingType(ctx context.Context, processingType models.ProcessingType, offset, limit int) ([]models.ProcessingJob, error)
}

// ProcessingResultRepository defines the contract for processing result operations
type ProcessingResultRepository interface {
	// Basic CRUD operations
	CreateResult(ctx context.Context, result *models.ProcessingResult) error
	UpdateResult(ctx context.Context, id int, result *models.ProcessingResult) error
	DeleteResult(ctx context.Context, id int) error
	GetResultByID(ctx context.Context, id int) (*models.ProcessingResult, error)
	
	// Result-specific operations
	GetResultByJobID(ctx context.Context, jobID int) (*models.ProcessingResult, error)
	GetResultsByImageID(ctx context.Context, imageID int) ([]models.ProcessingResult, error)
	DeleteResultsByJobID(ctx context.Context, jobID int) error
}
