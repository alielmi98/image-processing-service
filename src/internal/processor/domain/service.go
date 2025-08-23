package domain

import "github.com/alielmi98/image-processing-service/internal/image/entity"

// ProcessorService defines the interface for image processing operations
type ProcessorService interface {
	// ProcessImage processes an image based on the provided processing message
	ProcessImage(message entity.ProcessingMessage) error
	// GetProcessingStatus returns the status of a processing job
	GetProcessingStatus(jobID int) (*entity.ProcessingResult, error)
}
