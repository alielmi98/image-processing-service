package handlers

import (
	"net/http"
	"strconv"

	"github.com/alielmi98/image-processing-service/internal/image/api/dto"
	"github.com/alielmi98/image-processing-service/internal/image/usecase"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/helper"
	"github.com/gin-gonic/gin"
)

type ProcessingHandler struct {
	usecase *usecase.ProcessingUsecase
}

func NewProcessingHandler(cfg *config.Config, usecase *usecase.ProcessingUsecase) *ProcessingHandler {
	return &ProcessingHandler{
		usecase: usecase,
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

	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(response, true, 0))
}

// GetProcessingJobByID godoc
// @Summary Get a processing job by ID
// @Description Get a processing job by ID
// @Tags Processing
// @Accept json
// @produces json
// @param id path int true "Processing job ID"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.ProcessImageResponse} "Processing response"
// @Failure 404 {object} helper.BaseHttpResponse "Not found"
// @Router /v1/processing/{id} [get]
// @Security AuthBearer
func (h *ProcessingHandler) GetProcessingJobByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponse(nil, false, helper.ValidationError))
		return
	}
	job, err := h.usecase.GetProcessingJobByID(c, id)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	c.JSON(http.StatusOK, helper.GenerateBaseResponse(job, true, 0))
}
