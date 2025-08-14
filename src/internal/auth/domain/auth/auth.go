package auth

import (
	"github.com/alielmi98/image-processing-service/internal/auth/api/dto"
	"github.com/alielmi98/image-processing-service/internal/auth/entity"
	"github.com/golang-jwt/jwt"
)

type TokenProvider interface {
	GenerateToken(token *entity.TokenPayload) (*dto.TokenDetail, error)
	VerifyToken(token string) (*jwt.Token, error)
	GetClaims(token string) (map[string]interface{}, error)
	RefreshToken(refreshToken string) (*dto.TokenDetail, error)
}
