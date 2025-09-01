package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alielmi98/image-processing-service/constants"
	"gorm.io/gorm"
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

type JSONMap map[string]interface{}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONMap: %v", value)
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// ProcessingJob represents a processing job for an image
type ProcessingJob struct {
	Id             int            `gorm:"primarykey"`
	ImageId        int            `gorm:"not null;index"`
	Image          Image          `gorm:"foreignKey:ImageId;constraint:OnUpdate:NO ACTION;OnDelete:CASCADE"`
	ProcessingType ProcessingType `gorm:"type:varchar(50);not null"`
	Parameters     JSONMap        `gorm:"type:jsonb"` // JSON parameters for the processing operation
	Status         ImageStatus    `gorm:"type:varchar(20);not null;default:'pending'"`
	ResultPath     sql.NullString `gorm:"type:text;null"`
	ErrorMessage   sql.NullString `gorm:"type:text;null"`

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

func (m *ProcessingJob) BeforeCreate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = -1
	if value != nil {
		userId = int(value.(float64))
	}
	m.CreatedAt = time.Now().UTC()
	m.CreatedBy = userId
	return
}

func (m *Image) BeforeCreate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = -1
	if value != nil {
		userId = int(value.(float64))
	}
	m.CreatedAt = time.Now().UTC()
	m.CreatedBy = userId
	return
}
