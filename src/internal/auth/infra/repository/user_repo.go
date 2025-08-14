package repository

import (
	"context"
	"errors"
	"log"

	"github.com/alielmi98/image-processing-service/constants"
	"github.com/alielmi98/image-processing-service/internal/auth/domain/models"
	model "github.com/alielmi98/image-processing-service/internal/auth/domain/models"
	"github.com/alielmi98/image-processing-service/pkg/db"
	"github.com/alielmi98/image-processing-service/pkg/service_errors"
	"gorm.io/gorm"
)

type PgRepo struct {
	db *gorm.DB
}

func NewUserPgRepo() *PgRepo {
	return &PgRepo{db: db.GetDb()}
}

func (r *PgRepo) Create(ctx context.Context, user *model.User) error {
	tx := r.db.WithContext(ctx).Begin()
	err := tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}
func (r *PgRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.RecordNotFound}
		}
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return nil, err
	}
	return &user, nil
}

func (r *PgRepo) Update(ctx context.Context, id int, user *model.User) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Model(&model.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}
func (r *PgRepo) Delete(ctx context.Context, id int) error {
	tx := r.db.WithContext(ctx).Begin()
	if err := tx.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return err
	}
	tx.Commit()
	return nil
}

func (r *PgRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UsernameOrPasswordInvalid}
		}
		return nil, err
	}
	return &user, nil
}

func (r *PgRepo) ExistsByEmail(email string) (bool, error) {
	var exists bool
	if err := r.db.Model(&model.User{}).
		Select("count(*) > 0").
		Where("email = ?", email).
		Find(&exists).
		Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return false, err
	}
	return exists, nil
}

func (r *PgRepo) ExistsByUsername(username string) (bool, error) {
	var exists bool
	if err := r.db.Model(&model.User{}).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&exists).
		Error; err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Select, err.Error())
		return false, err
	}
	return exists, nil
}

func (r *PgRepo) GetDefaultRole() (roleId int, err error) {

	if err = r.db.Model(&models.Role{}).
		Select("id").
		Where("name = ?", constants.DefaultRoleName).
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}
