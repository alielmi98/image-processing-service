package messaging

import (
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/alielmi98/image-processing-service/pkg/contracts"
)

// ---- ProcessingType ----
func ToContractProcessingType(t models.ProcessingType) contracts.ProcessingType {
	switch t {
	case models.ProcessingTypeResize:
		return contracts.ProcessingTypeResize
	case models.ProcessingTypeCrop:
		return contracts.ProcessingTypeCrop
	case models.ProcessingTypeRotate:
		return contracts.ProcessingTypeRotate
	case models.ProcessingTypeFilter:
		return contracts.ProcessingTypeFilter
	case models.ProcessingTypeWatermark:
		return contracts.ProcessingTypeWatermark
	case models.ProcessingTypeCompress:
		return contracts.ProcessingTypeCompress
	case models.ProcessingTypeFormat:
		return contracts.ProcessingTypeFormat
	default:
		return contracts.ProcessingType("") // unknown
	}
}

func FromContractProcessingType(t contracts.ProcessingType) models.ProcessingType {
	switch t {
	case contracts.ProcessingTypeResize:
		return models.ProcessingTypeResize
	case contracts.ProcessingTypeCrop:
		return models.ProcessingTypeCrop
	case contracts.ProcessingTypeRotate:
		return models.ProcessingTypeRotate
	case contracts.ProcessingTypeFilter:
		return models.ProcessingTypeFilter
	case contracts.ProcessingTypeWatermark:
		return models.ProcessingTypeWatermark
	case contracts.ProcessingTypeCompress:
		return models.ProcessingTypeCompress
	case contracts.ProcessingTypeFormat:
		return models.ProcessingTypeFormat
	default:
		return models.ProcessingType("") // unknown
	}
}

// ---- ImageStatus ----
func ToContractImageStatus(s models.ImageStatus) contracts.ImageStatus {
	switch s {
	case models.ImageStatusPending:
		return contracts.ImageStatusPending
	case models.ImageStatusProcessing:
		return contracts.ImageStatusProcessing
	case models.ImageStatusCompleted:
		return contracts.ImageStatusCompleted
	case models.ImageStatusFailed:
		return contracts.ImageStatusFailed
	default:
		return contracts.ImageStatus("") // unknown
	}
}

func FromContractImageStatus(s contracts.ImageStatus) models.ImageStatus {
	switch s {
	case contracts.ImageStatusPending:
		return models.ImageStatusPending
	case contracts.ImageStatusProcessing:
		return models.ImageStatusProcessing
	case contracts.ImageStatusCompleted:
		return models.ImageStatusCompleted
	case contracts.ImageStatusFailed:
		return models.ImageStatusFailed
	default:
		return models.ImageStatus("") // unknown
	}
}
