package dto

import (
	"time"

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

type ProcessingUpdate struct {
	Status       models.ImageStatus
	ResultPath   string
	ErrorMessage string
	StartedAt    time.Time
	CompletedAt  time.Time
	Duration     int64
}
