package models

import (
	"database/sql"
	"time"
)

// ImageStatus represents the status of image processing
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

// Image represents an image record in the database
type Image struct {
	Id           int         `gorm:"primarykey"`
	UserId       int         `gorm:"not null;index"`
	OriginalName string      `gorm:"type:varchar(255);not null"`
	FileName     string      `gorm:"type:varchar(255);not null;unique"`
	FilePath     string      `gorm:"type:text;not null"`
	FileSize     int64       `gorm:"not null"`
	MimeType     string      `gorm:"type:varchar(100);not null"`
	Width        int         `gorm:"not null"`
	Height       int         `gorm:"not null"`
	Status       ImageStatus `gorm:"type:varchar(20);not null;default:'pending'"`

	// Processing metadata
	ProcessingJobs []ProcessingJob `gorm:"foreignKey:ImageId"`

	CreatedAt  time.Time    `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`

	CreatedBy  int            `gorm:"not null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
}

// ProcessingJob represents a processing job for an image
type ProcessingJob struct {
	Id             int                    `gorm:"primarykey"`
	ImageId        int                    `gorm:"not null;index"`
	Image          Image                  `gorm:"foreignKey:ImageId;constraint:OnUpdate:NO ACTION;OnDelete:CASCADE"`
	ProcessingType ProcessingType         `gorm:"type:varchar(50);not null"`
	Parameters     map[string]interface{} `gorm:"type:jsonb"` // JSON parameters for the processing operation
	Status         ImageStatus            `gorm:"type:varchar(20);not null;default:'pending'"`
	ResultPath     sql.NullString         `gorm:"type:text;null"`
	ErrorMessage   sql.NullString         `gorm:"type:text;null"`

	// Processing metrics
	StartedAt   sql.NullTime  `gorm:"type:TIMESTAMP with time zone;null"`
	CompletedAt sql.NullTime  `gorm:"type:TIMESTAMP with time zone;null"`
	Duration    sql.NullInt64 `gorm:"null"` // Duration in milliseconds

	CreatedAt  time.Time    `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`

	CreatedBy  int            `gorm:"not null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
}

// ProcessingResult represents the result of an image processing operation
type ProcessingResult struct {
	Id              int           `gorm:"primarykey"`
	ProcessingJobId int           `gorm:"not null;unique;index"`
	ProcessingJob   ProcessingJob `gorm:"foreignKey:ProcessingJobId;constraint:OnUpdate:NO ACTION;OnDelete:CASCADE"`
	ResultPath      string        `gorm:"type:text;not null"`
	FileSize        int64         `gorm:"not null"`
	Width           int           `gorm:"not null"`
	Height          int           `gorm:"not null"`
	MimeType        string        `gorm:"type:varchar(100);not null"`

	CreatedAt  time.Time    `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`

	CreatedBy  int            `gorm:"not null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
}
