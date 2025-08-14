package repository

import (
	"context"

	model "github.com/alielmi98/image-processing-service/internal/auth/domain/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, u model.User) (model.User, error)
	Update(ctx context.Context, id int, user *model.User) error
	Delete(ctx context.Context, id int) error
	FetchUserInfo(ctx context.Context, username string, password string) (model.User, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
}
