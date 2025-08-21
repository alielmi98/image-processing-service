package repository

import (
	"context"

	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/alielmi98/image-processing-service/internal/image/domain/repository"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/db"
	baseRepo "github.com/alielmi98/image-processing-service/pkg/repository"
	"gorm.io/gorm"
)

type ProcessingRepository struct {
	*baseRepo.BaseRepository[models.ProcessingJob]
	db *gorm.DB
}

func NewProcessingRepository(cfg *config.Config, preloads []db.PreloadEntity) repository.ProcessingRepository {
	database := db.GetDb()
	return &ProcessingRepository{
		BaseRepository: baseRepo.NewBaseRepository[models.ProcessingJob](cfg, database, preloads),
		db:             database,
	}
}

func (r *ProcessingRepository) CreateProcessingJob(ctx context.Context, job models.ProcessingJob) (models.ProcessingJob, error) {
	return r.Create(ctx, job)
}

func (r *ProcessingRepository) UpdateProcessingJob(ctx context.Context, id int, job map[string]interface{}) (models.ProcessingJob, error) {
	return r.Update(ctx, id, job)
}

func (r *ProcessingRepository) DeleteProcessingJob(ctx context.Context, id int) error {
	return r.Delete(ctx, id)
}

func (r *ProcessingRepository) GetProcessingJobByID(ctx context.Context, id int) (models.ProcessingJob, error) {
	return r.GetById(ctx, id)
}
