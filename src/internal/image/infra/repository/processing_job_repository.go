package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/alielmi98/image-processing-service/constants"
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/alielmi98/image-processing-service/internal/image/domain/repository"
	baseRepo "github.com/alielmi98/image-processing-service/pkg/repository"
	"github.com/alielmi98/image-processing-service/pkg/db"
	"gorm.io/gorm"
)

type ProcessingJobPgRepository struct {
	*baseRepo.BaseRepository[models.ProcessingJob]
	db *gorm.DB
}

func NewProcessingJobPgRepository() repository.ProcessingJobRepository {
	database := db.GetDb()
	return &ProcessingJobPgRepository{
		BaseRepository: baseRepo.NewBaseRepository[models.ProcessingJob](database),
		db:             database,
	}
}

func (r *ProcessingJobPgRepository) CreateJob(ctx context.Context, job *models.ProcessingJob) error {
	return r.Create(ctx, job)
}

func (r *ProcessingJobPgRepository) UpdateJob(ctx context.Context, id int, job *models.ProcessingJob) error {
	return r.Update(ctx, id, job)
}

func (r *ProcessingJobPgRepository) DeleteJob(ctx context.Context, id int) error {
	return r.Delete(ctx, id)
}

func (r *ProcessingJobPgRepository) GetJobByID(ctx context.Context, id int) (*models.ProcessingJob, error) {
	var job models.ProcessingJob
	if err := r.db.WithContext(ctx).
		Preload("Image").
		Where("id = ?", id).
		First(&job).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	return &job, nil
}

func (r *ProcessingJobPgRepository) GetJobsByImageID(ctx context.Context, imageID int) ([]models.ProcessingJob, error) {
	var jobs []models.ProcessingJob
	if err := r.db.WithContext(ctx).
		Where("image_id = ?", imageID).
		Order("created_at DESC").
		Find(&jobs).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	return jobs, nil
}

func (r *ProcessingJobPgRepository) GetJobsByStatus(ctx context.Context, status models.ImageStatus, offset, limit int) ([]models.ProcessingJob, error) {
	var jobs []models.ProcessingJob
	query := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at ASC")
	
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Find(&jobs).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	
	return jobs, nil
}

func (r *ProcessingJobPgRepository) GetPendingJobs(ctx context.Context, limit int) ([]models.ProcessingJob, error) {
	var jobs []models.ProcessingJob
	query := r.db.WithContext(ctx).
		Preload("Image").
		Where("status = ?", models.ImageStatusPending).
		Order("created_at ASC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Find(&jobs).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	
	return jobs, nil
}

func (r *ProcessingJobPgRepository) UpdateJobStatus(ctx context.Context, jobID int, status models.ImageStatus) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Model(&models.ProcessingJob{}).
		Where("id = ?", jobID).
		Update("status", status).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}

func (r *ProcessingJobPgRepository) UpdateJobProgress(ctx context.Context, jobID int, startedAt *int64, completedAt *int64, duration *int64, errorMsg *string) error {
	tx := r.db.WithContext(ctx).Begin()
	
	updates := make(map[string]interface{})
	
	if startedAt != nil {
		updates["started_at"] = sql.NullTime{Valid: true}
	}
	
	if completedAt != nil {
		updates["completed_at"] = sql.NullTime{Valid: true}
	}
	
	if duration != nil {
		updates["duration"] = sql.NullInt64{Int64: *duration, Valid: true}
	}
	
	if errorMsg != nil {
		updates["error_message"] = sql.NullString{String: *errorMsg, Valid: true}
	}
	
	if err := tx.Model(&models.ProcessingJob{}).
		Where("id = ?", jobID).
		Updates(updates).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	
	tx.Commit()
	return nil
}

func (r *ProcessingJobPgRepository) GetJobsByProcessingType(ctx context.Context, processingType models.ProcessingType, offset, limit int) ([]models.ProcessingJob, error) {
	var jobs []models.ProcessingJob
	query := r.db.WithContext(ctx).
		Where("processing_type = ?", processingType).
		Order("created_at DESC")
	
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Find(&jobs).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	
	return jobs, nil
}
