package entity

import (
	"time"
	
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
)

// ImageProcessingMessage represents a message sent to RabbitMQ for image processing
type ImageProcessingMessage struct {
	JobId          int                           `json:"job_id"`
	ImageId        int                           `json:"image_id"`
	UserId         int                           `json:"user_id"`
	ProcessingType models.ProcessingType         `json:"processing_type"`
	Parameters     map[string]interface{}        `json:"parameters"`
	SourcePath     string                        `json:"source_path"`
	DestinationDir string                        `json:"destination_dir"`
	Priority       int                           `json:"priority"` // 1-10, higher is more priority
	Timestamp      time.Time                     `json:"timestamp"`
	RetryCount     int                           `json:"retry_count"`
	MaxRetries     int                           `json:"max_retries"`
}

// ImageProcessingResult represents the result message sent back to RabbitMQ
type ImageProcessingResult struct {
	JobId        int                    `json:"job_id"`
	ImageId      int                    `json:"image_id"`
	UserId       int                    `json:"user_id"`
	Status       models.ImageStatus     `json:"status"`
	ResultPath   string                 `json:"result_path,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Duration     int64                  `json:"duration"` // Duration in milliseconds
	ProcessedAt  time.Time              `json:"processed_at"`
}

// ResizeParameters represents parameters for image resize operation
type ResizeParameters struct {
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	MaintainRatio bool   `json:"maintain_ratio"`
	Quality       int    `json:"quality"` // 1-100
	Format        string `json:"format,omitempty"` // jpg, png, webp, etc.
}

// CropParameters represents parameters for image crop operation
type CropParameters struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format,omitempty"`
}

// RotateParameters represents parameters for image rotation
type RotateParameters struct {
	Angle  float64 `json:"angle"` // Rotation angle in degrees
	Format string  `json:"format,omitempty"`
}

// FilterParameters represents parameters for image filters
type FilterParameters struct {
	FilterType string                 `json:"filter_type"` // blur, sharpen, grayscale, sepia, etc.
	Intensity  float64                `json:"intensity"`   // Filter intensity 0.0-1.0
	Options    map[string]interface{} `json:"options,omitempty"`
	Format     string                 `json:"format,omitempty"`
}

// WatermarkParameters represents parameters for watermark operation
type WatermarkParameters struct {
	WatermarkPath string  `json:"watermark_path"`
	Position      string  `json:"position"` // top-left, top-right, bottom-left, bottom-right, center
	Opacity       float64 `json:"opacity"`  // 0.0-1.0
	Scale         float64 `json:"scale"`    // Scale of watermark relative to image
	Format        string  `json:"format,omitempty"`
}

// CompressParameters represents parameters for image compression
type CompressParameters struct {
	Quality int    `json:"quality"` // 1-100
	Format  string `json:"format"`  // jpg, webp, etc.
}

// FormatParameters represents parameters for format conversion
type FormatParameters struct {
	TargetFormat string `json:"target_format"` // jpg, png, webp, gif, etc.
	Quality      int    `json:"quality"`       // 1-100 (for lossy formats)
}
