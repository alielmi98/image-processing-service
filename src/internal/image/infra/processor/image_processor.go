package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alielmi98/image-processing-service/internal/image/entity"
	"github.com/alielmi98/image-processing-service/internal/image/domain/models"
	"github.com/disintegration/imaging"
)

type ImageProcessorImpl struct {
	outputDir string
}

func NewImageProcessor(outputDir string) *ImageProcessorImpl {
	return &ImageProcessorImpl{
		outputDir: outputDir,
	}
}

func (p *ImageProcessorImpl) ProcessImage(ctx context.Context, message *entity.ImageProcessingMessage) (*entity.ImageProcessingResult, error) {
	startTime := time.Now()
	
	// Load the source image
	sourceImg, err := p.loadImage(message.SourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load source image: %w", err)
	}

	var processedImg image.Image
	var outputPath string
	var metadata map[string]interface{}

	// Process based on type
	switch message.ProcessingType {
	case models.ProcessingTypeResize:
		processedImg, outputPath, metadata, err = p.processResize(sourceImg, message)
	case models.ProcessingTypeCrop:
		processedImg, outputPath, metadata, err = p.processCrop(sourceImg, message)
	case models.ProcessingTypeRotate:
		processedImg, outputPath, metadata, err = p.processRotate(sourceImg, message)
	case models.ProcessingTypeFilter:
		processedImg, outputPath, metadata, err = p.processFilter(sourceImg, message)
	case models.ProcessingTypeWatermark:
		processedImg, outputPath, metadata, err = p.processWatermark(sourceImg, message)
	case models.ProcessingTypeCompress:
		processedImg, outputPath, metadata, err = p.processCompress(sourceImg, message)
	case models.ProcessingTypeFormat:
		processedImg, outputPath, metadata, err = p.processFormat(sourceImg, message)
	default:
		return nil, fmt.Errorf("unsupported processing type: %s", message.ProcessingType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to process image: %w", err)
	}

	// Save the processed image
	err = p.saveImage(processedImg, outputPath, p.getOutputFormat(message.Parameters))
	if err != nil {
		return nil, fmt.Errorf("failed to save processed image: %w", err)
	}

	duration := time.Since(startTime).Milliseconds()

	// Create result
	result := &entity.ImageProcessingResult{
		JobId:       message.JobId,
		ImageId:     message.ImageId,
		UserId:      message.UserId,
		Status:      models.ImageStatusCompleted,
		ResultPath:  outputPath,
		Metadata:    metadata,
		Duration:    duration,
		ProcessedAt: time.Now(),
	}

	return result, nil
}

func (p *ImageProcessorImpl) loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (p *ImageProcessorImpl) processResize(img image.Image, message *entity.ImageProcessingMessage) (image.Image, string, map[string]interface{}, error) {
	var params entity.ResizeParameters
	err := p.unmarshalParameters(message.Parameters, &params)
	if err != nil {
		return nil, "", nil, err
	}

	var resizedImg image.Image
	if params.MaintainRatio {
		resizedImg = imaging.Resize(img, params.Width, params.Height, imaging.Lanczos)
	} else {
		resizedImg = imaging.Fit(img, params.Width, params.Height, imaging.Lanczos)
	}

	outputPath := p.generateOutputPath(message, "resized", params.Format)
	
	metadata := map[string]interface{}{
		"original_width":  img.Bounds().Dx(),
		"original_height": img.Bounds().Dy(),
		"new_width":       resizedImg.Bounds().Dx(),
		"new_height":      resizedImg.Bounds().Dy(),
		"maintain_ratio":  params.MaintainRatio,
		"quality":         params.Quality,
	}

	return resizedImg, outputPath, metadata, nil
}

func (p *ImageProcessorImpl) processCrop(img image.Image, message *entity.ImageProcessingMessage) (image.Image, string, map[string]interface{}, error) {
	var params entity.CropParameters
	err := p.unmarshalParameters(message.Parameters, &params)
	if err != nil {
		return nil, "", nil, err
	}

	croppedImg := imaging.Crop(img, image.Rect(params.X, params.Y, params.X+params.Width, params.Y+params.Height))
	outputPath := p.generateOutputPath(message, "cropped", params.Format)
	
	metadata := map[string]interface{}{
		"original_width":  img.Bounds().Dx(),
		"original_height": img.Bounds().Dy(),
		"crop_x":          params.X,
		"crop_y":          params.Y,
		"crop_width":      params.Width,
		"crop_height":     params.Height,
	}

	return croppedImg, outputPath, metadata, nil
}

func (p *ImageProcessorImpl) processRotate(img image.Image, message *entity.ImageProcessingMessage) (image.Image, string, map[string]interface{}, error) {
	var params entity.RotateParameters
	err := p.unmarshalParameters(message.Parameters, &params)
	if err != nil {
		return nil, "", nil, err
	}

	rotatedImg := imaging.Rotate(img, params.Angle, image.Transparent)
	outputPath := p.generateOutputPath(message, "rotated", params.Format)
	
	metadata := map[string]interface{}{
		"rotation_angle": params.Angle,
		"original_width": img.Bounds().Dx(),
		"original_height": img.Bounds().Dy(),
		"new_width":      rotatedImg.Bounds().Dx(),
		"new_height":     rotatedImg.Bounds().Dy(),
	}

	return rotatedImg, outputPath, metadata, nil
}

func (p *ImageProcessorImpl) processFilter(img image.Image, message *entity.ImageProcessingMessage) (image.Image, string, map[string]interface{}, error) {
	var params entity.FilterParameters
	err := p.unmarshalParameters(message.Parameters, &params)
	if err != nil {
		return nil, "", nil, err
	}

	var filteredImg image.Image

	switch strings.ToLower(params.FilterType) {
	case "blur":
		filteredImg = imaging.Blur(img, params.Intensity*10) // Scale intensity for blur
	case "sharpen":
		filteredImg = imaging.Sharpen(img, params.Intensity*2) // Scale intensity for sharpen
	case "grayscale":
		filteredImg = imaging.Grayscale(img)
	case "invert":
		filteredImg = imaging.Invert(img)
	case "brightness":
		filteredImg = imaging.AdjustBrightness(img, params.Intensity*100-50) // -50 to +50
	case "contrast":
		filteredImg = imaging.AdjustContrast(img, params.Intensity*100-50) // -50 to +50
	default:
		return nil, "", nil, fmt.Errorf("unsupported filter type: %s", params.FilterType)
	}

	outputPath := p.generateOutputPath(message, fmt.Sprintf("filtered_%s", params.FilterType), params.Format)
	
	metadata := map[string]interface{}{
		"filter_type": params.FilterType,
		"intensity":   params.Intensity,
		"options":     params.Options,
	}

	return filteredImg, outputPath, metadata, nil
}

func (p *ImageProcessorImpl) processWatermark(img image.Image, message *entity.ImageProcessingMessage) (image.Image, string, map[string]interface{}, error) {
	var params entity.WatermarkParameters
	err := p.unmarshalParameters(message.Parameters, &params)
	if err != nil {
		return nil, "", nil, err
	}

	// Load watermark image
	watermark, err := p.loadImage(params.WatermarkPath)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to load watermark: %w", err)
	}

	// Scale watermark if needed
	if params.Scale != 1.0 {
		watermarkBounds := watermark.Bounds()
		newWidth := int(float64(watermarkBounds.Dx()) * params.Scale)
		newHeight := int(float64(watermarkBounds.Dy()) * params.Scale)
		watermark = imaging.Resize(watermark, newWidth, newHeight, imaging.Lanczos)
	}

	// Calculate position
	imgBounds := img.Bounds()
	watermarkBounds := watermark.Bounds()
	var x, y int

	switch strings.ToLower(params.Position) {
	case "top-left":
		x, y = 0, 0
	case "top-right":
		x, y = imgBounds.Dx()-watermarkBounds.Dx(), 0
	case "bottom-left":
		x, y = 0, imgBounds.Dy()-watermarkBounds.Dy()
	case "bottom-right":
		x, y = imgBounds.Dx()-watermarkBounds.Dx(), imgBounds.Dy()-watermarkBounds.Dy()
	case "center":
		x, y = (imgBounds.Dx()-watermarkBounds.Dx())/2, (imgBounds.Dy()-watermarkBounds.Dy())/2
	default:
		x, y = 0, 0 // Default to top-left
	}

	// Apply watermark
	watermarkedImg := imaging.Overlay(img, watermark, image.Pt(x, y), params.Opacity)
	outputPath := p.generateOutputPath(message, "watermarked", params.Format)
	
	metadata := map[string]interface{}{
		"watermark_path": params.WatermarkPath,
		"position":       params.Position,
		"opacity":        params.Opacity,
		"scale":          params.Scale,
		"watermark_x":    x,
		"watermark_y":    y,
	}

	return watermarkedImg, outputPath, metadata, nil
}

func (p *ImageProcessorImpl) processCompress(img image.Image, message *entity.ImageProcessingMessage) (image.Image, string, map[string]interface{}, error) {
	var params entity.CompressParameters
	err := p.unmarshalParameters(message.Parameters, &params)
	if err != nil {
		return nil, "", nil, err
	}

	outputPath := p.generateOutputPath(message, "compressed", params.Format)
	
	metadata := map[string]interface{}{
		"quality":        params.Quality,
		"target_format":  params.Format,
		"original_width": img.Bounds().Dx(),
		"original_height": img.Bounds().Dy(),
	}

	return img, outputPath, metadata, nil
}

func (p *ImageProcessorImpl) processFormat(img image.Image, message *entity.ImageProcessingMessage) (image.Image, string, map[string]interface{}, error) {
	var params entity.FormatParameters
	err := p.unmarshalParameters(message.Parameters, &params)
	if err != nil {
		return nil, "", nil, err
	}

	outputPath := p.generateOutputPath(message, "converted", params.TargetFormat)
	
	metadata := map[string]interface{}{
		"original_format": filepath.Ext(message.SourcePath),
		"target_format":   params.TargetFormat,
		"quality":         params.Quality,
	}

	return img, outputPath, metadata, nil
}

func (p *ImageProcessorImpl) unmarshalParameters(params map[string]interface{}, target interface{}) error {
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (p *ImageProcessorImpl) generateOutputPath(message *entity.ImageProcessingMessage, suffix, format string) string {
	if format == "" {
		format = "jpg" // Default format
	}
	
	filename := fmt.Sprintf("user_%d_image_%d_job_%d_%s_%d.%s",
		message.UserId,
		message.ImageId,
		message.JobId,
		suffix,
		time.Now().Unix(),
		format,
	)
	
	return filepath.Join(p.outputDir, filename)
}

func (p *ImageProcessorImpl) getOutputFormat(params map[string]interface{}) string {
	if format, ok := params["format"].(string); ok && format != "" {
		return format
	}
	return "jpg" // Default format
}

func (p *ImageProcessorImpl) saveImage(img image.Image, path, format string) error {
	// Ensure output directory exists
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch strings.ToLower(format) {
	case "jpg", "jpeg":
		quality := 85 // Default quality
		return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
	case "png":
		return png.Encode(file, img)
	case "gif":
		return gif.Encode(file, img, nil)
	case "webp":
		// WebP encoding requires a different approach
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 85}) // Fallback to JPEG for now
	default:
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 85})
	}
}
