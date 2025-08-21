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

type ImagePgRepository struct {
	*baseRepo.BaseRepository[models.Image]
	db *gorm.DB
}

func NewImagePgRepository(cfg *config.Config, preloads []db.PreloadEntity) repository.ImageRepository {
	database := db.GetDb()
	return &ImagePgRepository{
		BaseRepository: baseRepo.NewBaseRepository[models.Image](cfg, database, preloads),
		db:             database,
	}
}

func (r *ImagePgRepository) CreateImage(ctx context.Context, image models.Image) (models.Image, error) {
	return r.Create(ctx, image)
}

func (r *ImagePgRepository) UpdateImage(ctx context.Context, id int, image map[string]interface{}) (models.Image, error) {
	return r.Update(ctx, id, image)
}

func (r *ImagePgRepository) DeleteImage(ctx context.Context, id int) error {
	return r.Delete(ctx, id)
}

func (r *ImagePgRepository) GetImageByID(ctx context.Context, id int) (models.Image, error) {
	return r.GetById(ctx, id)
}
