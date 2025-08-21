package routers

import (
	"github.com/alielmi98/image-processing-service/di"
	"github.com/alielmi98/image-processing-service/internal/image/api/handlers"
	"github.com/alielmi98/image-processing-service/internal/middlewares"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/gin-gonic/gin"
)

func Image(r *gin.RouterGroup, cfg *config.Config) {
	handler := handlers.NewImageHandler(cfg)
	tokenProvider := di.GetTokenProvider(cfg)
	r.POST("/", middlewares.Authentication(cfg, tokenProvider), handler.Create)
}
