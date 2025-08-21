package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/alielmi98/image-processing-service/internal/image/usecase"
	"github.com/gin-gonic/gin"
)

type ProcessingHandler struct {
	processingUseCase usecase.ProcessingUseCase
}

func NewProcessingHandler(processingUseCase usecase.ProcessingUseCase) *ProcessingHandler {
	return &ProcessingHandler{
		processingUseCase: processingUseCase,
	}
}

type ProcessImageRequest struct {
	ImageId        int                    `json:"image_id" binding:"required"`
	ProcessingType string                 `json:"processing_type" binding:"required"`
	Parameters     map[string]interface{} `json:"parameters" binding:"required"`
	SourcePath     string                 `json:"source_path" binding:"required"`
	DestinationDir string                 `json:"destination_dir,omitempty"`
}

type ProcessImageResponse struct {
	Success    bool                   `json:"success"`
	JobId      int                    `json:"job_id,omitempty"`
	Message    string                 `json:"message"`
	ResultPath string                 `json:"result_path,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ProcessImage handles image processing requests
func (h *ProcessingHandler) ProcessImage(c *gin.Context) {
	var req ProcessImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ProcessImageResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context (assuming it's set by auth middleware)
	userIdStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ProcessImageResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	userId, ok := userIdStr.(int)
	if !ok {
		// Try to convert from string if it's stored as string
		if userIdStrVal, ok := userIdStr.(string); ok {
			var err error
			userId, err = strconv.Atoi(userIdStrVal)
			if err != nil {
				c.JSON(http.StatusBadRequest, ProcessImageResponse{
					Success: false,
					Message: "Invalid user ID",
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, ProcessImageResponse{
				Success: false,
				Message: "Invalid user ID format",
			})
			return
		}
	}

	// Set default destination directory if not provided
	if req.DestinationDir == "" {
		req.DestinationDir = "./processed_images"
	}

	// Convert processing type
	processingType := models.ProcessingType(req.ProcessingType)

	// Create processing job
	ctx := context.Background()
	message, err := h.processingUseCase.CreateProcessingJob(
		ctx,
		req.ImageId,
		userId,
		processingType,
		req.Parameters,
		req.SourcePath,
		req.DestinationDir,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ProcessImageResponse{
			Success: false,
			Message: "Failed to create processing job: " + err.Error(),
		})
		return
	}

	// Return success response with job ID
	c.JSON(http.StatusAccepted, ProcessImageResponse{
		Success: true,
		JobId:   message.JobId,
		Message: "Image processing job queued successfully",
	})
}

// ProcessImageAsync handles asynchronous image processing requests
func (h *ProcessingHandler) ProcessImageAsync(c *gin.Context) {
	var req ProcessImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ProcessImageResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Get user ID from context
	userIdStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ProcessImageResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	userId, ok := userIdStr.(int)
	if !ok {
		if userIdStrVal, ok := userIdStr.(string); ok {
			var err error
			userId, err = strconv.Atoi(userIdStrVal)
			if err != nil {
				c.JSON(http.StatusBadRequest, ProcessImageResponse{
					Success: false,
					Message: "Invalid user ID",
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, ProcessImageResponse{
				Success: false,
				Message: "Invalid user ID format",
			})
			return
		}
	}

	// Set default destination directory if not provided
	if req.DestinationDir == "" {
		req.DestinationDir = "./processed_images"
	}

	// Convert processing type
	processingType := models.ProcessingType(req.ProcessingType)

	// Create processing job
	ctx := context.Background()
	message, err := h.processingUseCase.CreateProcessingJob(
		ctx,
		req.ImageId,
		userId,
		processingType,
		req.Parameters,
		req.SourcePath,
		req.DestinationDir,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ProcessImageResponse{
			Success: false,
			Message: "Failed to create processing job: " + err.Error(),
		})
		return
	}

	// Send message to processor for asynchronous processing
	err = h.processingUseCase.SendProcessingMessage(ctx, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ProcessImageResponse{
			Success: false,
			Message: "Failed to send processing message: " + err.Error(),
		})
		return
	}

	// Return success response with job ID
	c.JSON(http.StatusAccepted, ProcessImageResponse{
		Success: true,
		JobId:   message.JobId,
		Message: "Image processing job queued successfully",
	})
}
