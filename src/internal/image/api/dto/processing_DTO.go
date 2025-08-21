package dto

import (
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	usecaseDto "github.com/alielmi98/image-processing-service/internal/image/usecase/dto"
)

type CreateProcessImageRequest struct {
	ImageId        int                    `json:"image_id" binding:"required"`
	ProcessingType models.ProcessingType  `json:"processing_type" binding:"required"`
	Parameters     map[string]interface{} `json:"parameters" binding:"required"`
}

type ProcessImageResponse struct {
	JobId int `json:"job_id,omitempty"`
}

func ToCreateProcessImageRequest(from CreateProcessImageRequest) usecaseDto.ProcessingRequest {
	return usecaseDto.ProcessingRequest{
		ImageId:        from.ImageId,
		ProcessingType: from.ProcessingType,
		Parameters:     from.Parameters,
	}
}

func ToProcessImageResponse(from usecaseDto.ProcessingResponse) ProcessImageResponse {
	return ProcessImageResponse{
		JobId: from.JobId,
	}
}
