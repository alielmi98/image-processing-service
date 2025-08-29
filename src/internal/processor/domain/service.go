package domain

import "github.com/alielmi98/image-processing-service/pkg/contracts"

// ProcessorService defines the interface for image processing operations
type ProcessorService interface {
	// ProcessImage processes an image based on the provided processing message
	ProcessImage(message contracts.ProcessingMessage) (*contracts.ProcessingResult, error)
}
