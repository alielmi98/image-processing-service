package repository

import (
	"context"
	"log"

	"github.com/alielmi98/image-processing-service/constants"
	"gorm.io/gorm"
)

// BaseRepository provides common CRUD operations for all entities
type BaseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository instance
func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// Create creates a new entity
func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Create(entity).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}

// Update updates an existing entity by ID
func (r *BaseRepository[T]) Update(ctx context.Context, id int, entity *T) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Model(new(T)).Where("id = ?", id).Updates(entity).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}

// Delete soft deletes an entity by ID
func (r *BaseRepository[T]) Delete(ctx context.Context, id int) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Where("id = ?", id).Delete(new(T)).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}

// GetByID retrieves an entity by ID
func (r *BaseRepository[T]) GetByID(ctx context.Context, id int) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	return &entity, nil
}

// GetAll retrieves all entities with pagination
func (r *BaseRepository[T]) GetAll(ctx context.Context, offset, limit int) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx)

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&entities).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}

	return entities, nil
}

// Count returns the total count of entities
func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(new(T)).Count(&count).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return 0, err
	}
	return count, nil
}

// Exists checks if an entity exists by a given condition
func (r *BaseRepository[T]) Exists(ctx context.Context, condition string, args ...interface{}) (bool, error) {
	var exists bool
	if err := r.db.WithContext(ctx).Model(new(T)).
		Select("count(*) > 0").
		Where(condition, args...).
		Find(&exists).Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return false, err
	}
	return exists, nil
}
