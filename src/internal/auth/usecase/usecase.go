package usecase

import (
	"context"
	"log"

	"github.com/alielmi98/image-processing-service/constants"
	"github.com/alielmi98/image-processing-service/internal/auth/api/dto"
	"github.com/alielmi98/image-processing-service/internal/auth/domain/auth"
	model "github.com/alielmi98/image-processing-service/internal/auth/domain/models"
	repository "github.com/alielmi98/image-processing-service/internal/auth/domain/repository"
	"github.com/alielmi98/image-processing-service/internal/auth/entity"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/service_errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	cfg   *config.Config
	repo  repository.UserRepository
	token auth.TokenProvider
}

func NewUserUsecase(cfg *config.Config, repository repository.UserRepository, token auth.TokenProvider) *UserUsecase {
	return &UserUsecase{
		cfg:   cfg,
		repo:  repository,
		token: token,
	}
}

// Register by username
func (s *UserUsecase) RegisterByUsername(ctx context.Context, req *dto.RegisterUserByUsernameRequest) error {
	u := model.User{
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
	// Check if username already exists
	if existing, _ := s.repo.ExistsByUsername(req.Username); existing {
		return &service_errors.ServiceError{EndUserMessage: service_errors.UsernameExists}
	}
	// Check if email already exists
	if existing, _ := s.repo.ExistsByEmail(req.Email); existing {
		return &service_errors.ServiceError{EndUserMessage: service_errors.EmailExists}
	}
	// Hash password
	bp := []byte(req.Password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Caller:%s Level:%s Msg:%s", constants.General, constants.HashPassword, err.Error())
		return err
	}
	req.Password = string(hp)
	u.Password = req.Password

	_, err = s.repo.CreateUser(ctx, u)
	if err != nil {
		return err
	}
	return nil

}

func (s *UserUsecase) LoginByUsername(ctx context.Context, req *dto.LoginByUsernameRequest) (*dto.TokenDetail, error) {
	user, err := s.repo.FetchUserInfo(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UsernameOrPasswordInvalid}
	}

	token, err := s.generateToken(&user)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *UserUsecase) RefreshToken(refreshToken string) (*dto.TokenDetail, error) {
	tokenDetail, err := s.token.RefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return tokenDetail, nil
}

func (s *UserUsecase) generateToken(user *model.User) (*dto.TokenDetail, error) {
	tokenDto := entity.TokenPayload{UserId: user.Id, FirstName: user.FirstName, LastName: user.LastName,
		Email: user.Email, MobileNumber: user.MobileNumber}

	if len(*user.UserRoles) > 0 {
		for _, ur := range *user.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := s.token.GenerateToken(&tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil
}
