package handlers

import (
	"net/http"

	"github.com/alielmi98/image-processing-service/di"
	"github.com/alielmi98/image-processing-service/internal/image/api/dto"
	"github.com/alielmi98/image-processing-service/internal/image/usecase"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/helper"
	"github.com/gin-gonic/gin"
)

type ProcessingHandler struct {
	usecase *usecase.ProcessingUsecase
}

func NewProcessingHandler(cfg *config.Config) *ProcessingHandler {
	return &ProcessingHandler{
		usecase: usecase.NewProcessingUseCase(cfg, di.GetProcessingRepository(cfg), di.GetMessageSender(cfg)),
	}
}

// CreateProcessingJob godoc
// @Summary Create an image processing job
// @Description Create an image processing job
// @Tags Processing
// @Accept json
// @produces json
// @param request body dto.CreateProcessImageRequest true "Processing request"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.ProcessImageResponse} "Processing response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/processing [post]
// @Security AuthBearer
func (h *ProcessingHandler) CreateProcessingJob(c *gin.Context) {
	var request dto.CreateProcessImageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, helper.BaseHttpResponse{Error: err.Error()})
		return
	}

	response, err := h.usecase.CreateProcessingJob(c, dto.ToCreateProcessImageRequest(request))
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.BaseHttpResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, helper.BaseHttpResponse{Result: response})
}
