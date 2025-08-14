package di

import (
	contractAuth "github.com/alielmi98/image-processing-service/internal/auth/domain/auth"
	contractAuthRepo "github.com/alielmi98/image-processing-service/internal/auth/domain/repository"
	infraAuth "github.com/alielmi98/image-processing-service/internal/auth/infra/auth"
	infraAuthRepo "github.com/alielmi98/image-processing-service/internal/auth/infra/repository"
	"github.com/alielmi98/image-processing-service/pkg/config"
)

// midedlewares
func GetTokenProvider(cfg *config.Config) contractAuth.TokenProvider {
	return infraAuth.NewJwtProvider(cfg)
}

func GetUserRepository(cfg *config.Config) contractAuthRepo.UserRepository {
	return infraAuthRepo.NewUserPgRepo()
}
