package entity

// ResizeParameters represents parameters for image resize operation
type ResizeParameters struct {
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	MaintainRatio bool   `json:"maintain_ratio"`
	Quality       int    `json:"quality"`          // 1-100
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
