package dto

type CreateImage struct {
	FileName     string
	OriginalName string
	FilePath     string
	MimeType     string
	UserID       int
	FileSize     int64
	Width        int
	Height       int
}

type UpdateImage struct {
	FileName string
}

type ImageResponse struct {
	Id           int
	FileName     string
	OriginalName string
	FilePath     string
	MimeType     string
	FileSize     int64
	Width        int
	Height       int
}
