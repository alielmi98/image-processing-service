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

type ImagePgRepository struct {
	*baseRepo.BaseRepository[models.Image]
	db *gorm.DB
}

func NewImagePgRepository() repository.ImageRepository {
	database := db.GetDb()
	return &ImagePgRepository{
		BaseRepository: baseRepo.NewBaseRepository[models.Image](database),
		db:             database,
	}
}

func (r *ImagePgRepository) CreateImage(ctx context.Context, image *models.Image) error {
	return r.Create(ctx, image)
}

func (r *ImagePgRepository) UpdateImage(ctx context.Context, id int, image *models.Image) error {
	return r.Update(ctx, id, image)
}

func (r *ImagePgRepository) DeleteImage(ctx context.Context, id int) error {
	return r.Delete(ctx, id)
}

func (r *ImagePgRepository) GetImageByID(ctx context.Context, id int) (*models.Image, error) {
	var image models.Image
	if err := r.db.WithContext(ctx).
		Preload("ProcessingJobs").
		Where("id = ?", id).
		First(&image).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	return &image, nil
}

func (r *ImagePgRepository) GetImagesByUserID(ctx context.Context, userID int, offset, limit int) ([]models.Image, error) {
	var images []models.Image
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC")
	
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Find(&images).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	
	return images, nil
}

func (r *ImagePgRepository) GetImageByFileName(ctx context.Context, fileName string) (*models.Image, error) {
	var image models.Image
	if err := r.db.WithContext(ctx).
		Where("file_name = ?", fileName).
		First(&image).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	return &image, nil
}

func (r *ImagePgRepository) ExistsByFileName(ctx context.Context, fileName string) (bool, error) {
	return r.Exists(ctx, "file_name = ?", fileName)
}

func (r *ImagePgRepository) UpdateImageStatus(ctx context.Context, imageID int, status models.ImageStatus) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Model(&models.Image{}).
		Where("id = ?", imageID).
		Update("status", status).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}

func (r *ImagePgRepository) GetImagesByStatus(ctx context.Context, status models.ImageStatus, offset, limit int) ([]models.Image, error) {
	var images []models.Image
	query := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at ASC")
	
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Find(&images).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	
	return images, nil
}

func (r *ImagePgRepository) CountImagesByUserID(ctx context.Context, userID int) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Image{}).
		Where("user_id = ?", userID).
		Count(&count).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return 0, err
	}
	return count, nil
}
