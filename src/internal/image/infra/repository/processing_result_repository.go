package repository

import (
	"context"
	"log"

	"github.com/alielmi98/image-processing-service/constants"
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/alielmi98/image-processing-service/internal/image/domain/repository"
	baseRepo "github.com/alielmi98/image-processing-service/pkg/repository"
	"github.com/alielmi98/image-processing-service/pkg/db"
	"gorm.io/gorm"
)

type ProcessingResultPgRepository struct {
	*baseRepo.BaseRepository[models.ProcessingResult]
	db *gorm.DB
}

func NewProcessingResultPgRepository() repository.ProcessingResultRepository {
	database := db.GetDb()
	return &ProcessingResultPgRepository{
		BaseRepository: baseRepo.NewBaseRepository[models.ProcessingResult](database),
		db:             database,
	}
}

func (r *ProcessingResultPgRepository) CreateResult(ctx context.Context, result *models.ProcessingResult) error {
	return r.Create(ctx, result)
}

func (r *ProcessingResultPgRepository) UpdateResult(ctx context.Context, id int, result *models.ProcessingResult) error {
	return r.Update(ctx, id, result)
}

func (r *ProcessingResultPgRepository) DeleteResult(ctx context.Context, id int) error {
	return r.Delete(ctx, id)
}

func (r *ProcessingResultPgRepository) GetResultByID(ctx context.Context, id int) (*models.ProcessingResult, error) {
	var result models.ProcessingResult
	if err := r.db.WithContext(ctx).
		Preload("ProcessingJob").
		Preload("ProcessingJob.Image").
		Where("id = ?", id).
		First(&result).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	return &result, nil
}

func (r *ProcessingResultPgRepository) GetResultByJobID(ctx context.Context, jobID int) (*models.ProcessingResult, error) {
	var result models.ProcessingResult
	if err := r.db.WithContext(ctx).
		Preload("ProcessingJob").
		Where("processing_job_id = ?", jobID).
		First(&result).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	return &result, nil
}

func (r *ProcessingResultPgRepository) GetResultsByImageID(ctx context.Context, imageID int) ([]models.ProcessingResult, error) {
	var results []models.ProcessingResult
	if err := r.db.WithContext(ctx).
		Preload("ProcessingJob").
		Joins("JOIN processing_jobs ON processing_results.processing_job_id = processing_jobs.id").
		Where("processing_jobs.image_id = ?", imageID).
		Order("processing_results.created_at DESC").
		Find(&results).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	return results, nil
}

func (r *ProcessingResultPgRepository) DeleteResultsByJobID(ctx context.Context, jobID int) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Where("processing_job_id = ?", jobID).Delete(&models.ProcessingResult{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}
