package dto

import (
	"mime/multipart"

	"github.com/alielmi98/image-processing-service/internal/image/usecase/dto"
)

type UploadImageRequest struct {
	Image *multipart.FileHeader `json:"file" form:"file" binding:"required" swaggerignore:"true"`
}

type CreateImageRequest struct {
	FileName     string `json:"file-name"`
	OriginalName string `json:"original-name"`
	FilePath     string `json:"file-path"`
	MimeType     string `json:"mime-type"`
	FileSize     int64  `json:"file-size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

type UpdateImageRequest struct {
	FileName string `json:"file-name"`
}

type ImageResponse struct {
	Id           int    `json:"id"`
	FileName     string `json:"file-name"`
	OriginalName string `json:"original-name"`
	FilePath     string `json:"file-path"`
	MimeType     string `json:"mime-type"`
	FileSize     int64  `json:"file-size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

func ToImageResponse(from dto.ImageResponse) ImageResponse {
	return ImageResponse{
		Id:           from.Id,
		FileName:     from.FileName,
		OriginalName: from.OriginalName,
		FilePath:     from.FilePath,
		MimeType:     from.MimeType,
		FileSize:     from.FileSize,
		Width:        from.Width,
		Height:       from.Height,
	}
}

func ToCreateImage(from CreateImageRequest) dto.CreateImage {
	return dto.CreateImage{
		FileName:     from.FileName,
		OriginalName: from.OriginalName,
		FilePath:     from.FilePath,
		MimeType:     from.MimeType,
		FileSize:     from.FileSize,
		Width:        from.Width,
		Height:       from.Height,
	}
}

func ToUpdateImage(from UpdateImageRequest) dto.UpdateImage {
	return dto.UpdateImage{
		FileName: from.FileName,
	}
}
