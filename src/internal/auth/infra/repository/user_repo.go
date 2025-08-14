package repository

import (
	"context"
	"log"

	"github.com/alielmi98/image-processing-service/constants"
	model "github.com/alielmi98/image-processing-service/internal/auth/domain/models"
	"github.com/alielmi98/image-processing-service/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const userFilterExp string = "username = ?"
const countFilterExp string = "count(*) > 0"

type PgRepo struct {
	db *gorm.DB
}

func NewUserPgRepo() *PgRepo {
	return &PgRepo{db: db.GetDb()}
}

func (r *PgRepo) CreateUser(ctx context.Context, u model.User) (model.User, error) {

	roleId, err := r.GetDefaultRole(ctx)
	if err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.DefaultRoleNotFound, err.Error())

		return u, err
	}
	tx := r.db.WithContext(ctx).Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())

		return u, err
	}
	err = tx.Create(&model.UserRole{RoleId: roleId, UserId: u.Id}).Error
	if err != nil {
		tx.Rollback()
		log.Printf("Caller:%s Level:%s Msg:%s", constants.Postgres, constants.Rollback, err.Error())
		return u, err
	}
	tx.Commit()
	return u, nil
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

func (r *PgRepo) GetDefaultRole(ctx context.Context) (roleId int, err error) {

	if err = r.db.WithContext(ctx).Model(&model.Role{}).
		Select("id").
		Where("name = ?", constants.DefaultRoleName).
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}

func (r *PgRepo) FetchUserInfo(ctx context.Context, username string, password string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where(userFilterExp, username).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&user).Error

	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}
