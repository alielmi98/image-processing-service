package di

import (
	contractAuth "github.com/alielmi98/image-processing-service/internal/auth/domain/auth"
	infraAuth "github.com/alielmi98/image-processing-service/internal/auth/infra/auth"
	"github.com/alielmi98/image-processing-service/pkg/config"
)

// midedlewares
func GetTokenProvider(cfg *config.Config) contractAuth.TokenProvider {
	return infraAuth.NewJwtProvider(cfg)
}
