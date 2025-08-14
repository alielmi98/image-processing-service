package router

import (
	"github.com/alielmi98/image-processing-service/internal/auth/api/handler"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/gin-gonic/gin"
)

func Auth(router *gin.RouterGroup, cfg *config.Config) {
	handler := handler.NewAuthHandler(cfg)
	router.POST("/register", handler.RegisterByUsername)
	router.POST("/login", handler.LoginByUsername)
	router.POST("/refresh-token", handler.RefreshToken)

}
