package dto

import (
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
)

type ProcessingRequest struct {
	ImageId        int
	ProcessingType models.ProcessingType
	Parameters     map[string]interface{}
}

type ProcessingResponse struct {
	JobId int
}
