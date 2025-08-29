package contracts

import (
	"time"
)

type ImageStatus string

const (
	ImageStatusPending    ImageStatus = "pending"
	ImageStatusProcessing ImageStatus = "processing"
	ImageStatusCompleted  ImageStatus = "completed"
	ImageStatusFailed     ImageStatus = "failed"
)

// ProcessingType represents different types of image processing operations
type ProcessingType string

const (
	ProcessingTypeResize    ProcessingType = "resize"
	ProcessingTypeCrop      ProcessingType = "crop"
	ProcessingTypeRotate    ProcessingType = "rotate"
	ProcessingTypeFilter    ProcessingType = "filter"
	ProcessingTypeWatermark ProcessingType = "watermark"
	ProcessingTypeCompress  ProcessingType = "compress"
	ProcessingTypeFormat    ProcessingType = "format"
)

type ProcessingMessage struct {
	JobId          int                    `json:"job_id"`
	ImageId        int                    `json:"image_id"`
	UserId         int                    `json:"user_id"`
	ProcessingType ProcessingType         `json:"processing_type"`
	Parameters     map[string]interface{} `json:"parameters"`
	SourcePath     string                 `json:"source_path"`
	Priority       int                    `json:"priority"` // 1-10, higher is more priority
	Timestamp      time.Time              `json:"timestamp"`
	RetryCount     int                    `json:"retry_count"`
	MaxRetries     int                    `json:"max_retries"`
}

// ProcessingResult represents the result of processing
type ProcessingResult struct {
	JobId        int                    `json:"job_id"`
	ImageId      int                    `json:"image_id"`
	UserId       int                    `json:"user_id"`
	Status       ImageStatus            `json:"status"`
	ResultPath   string                 `json:"result_path,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Duration     int64                  `json:"duration"` // Duration in milliseconds
	StartedAt    time.Time              `json:"started_at"`
	ProcessedAt  time.Time              `json:"processed_at"`
}
