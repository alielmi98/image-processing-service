package routers

import (
	"github.com/alielmi98/image-processing-service/internal/image/api/handlers"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/gin-gonic/gin"
)

func Image(r *gin.RouterGroup, cfg *config.Config) {
	handler := handlers.NewImageHandler(cfg)
	r.POST("/", handler.Create)

}

func Processing(r *gin.RouterGroup, cfg *config.Config) {
	handler := handlers.NewProcessingHandler(cfg)

	r.POST("/", handler.CreateProcessingJob)
}
