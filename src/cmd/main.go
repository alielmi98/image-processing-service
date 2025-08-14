package main

import (
	"fmt"
	"log"

	"github.com/alielmi98/image-processing-service/constants"
	"github.com/alielmi98/image-processing-service/docs"
	authRouter "github.com/alielmi98/image-processing-service/internal/auth/api/routers"
	"github.com/alielmi98/image-processing-service/internal/middlewares"
	migration "github.com/alielmi98/image-processing-service/migrations"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/db"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey AuthBearer
// @in header
// @name Authorization
func main() {
	cfg := config.GetConfig()

	err := db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		log.Fatalf("caller:%s  Level:%s  Msg:%s", constants.Postgres, constants.Startup, err.Error())
	}

	// Migrate the database
	migration.Up1()

	InitServer(cfg)

}

func InitServer(cfg *config.Config) {
	r := gin.New()

	r.Use(middlewares.Cors(cfg), middlewares.LimitByRequest())
	RegisterRoutes(r, cfg)
	RegisterSwagger(r, cfg)
	log.Printf("Caller:%s Level:%s Msg:%s", constants.General, constants.Startup, "Started")
	r.Run(fmt.Sprintf(":%s", cfg.Server.InternalPort))

}

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		//Auth
		auth := v1.Group("/auth")
		authRouter.Auth(auth, cfg)
	}

}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "golang web api"
	docs.SwaggerInfo.Description = "This is a sample server for golang web api"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Server.InternalPort)
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
