package main

import (
	"log"
	"os"

	"github.com/alielmi98/image-processing-service/internal/image/infra/service"
	"github.com/alielmi98/image-processing-service/pkg/config"
)

func main() {
	log.Println("Image Processing Service starting...")

	// Load configuration
	cfg := config.GetConfig()
	if cfg == nil {
		log.Fatalf("Failed to load configuration")
	}

	// Get output directory from environment or use default
	outputDir := os.Getenv("IMAGE_OUTPUT_DIR")
	if outputDir == "" {
		outputDir = "./processed_images"
	}

	// Create output directory if it doesn't exist
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create and start the image processing service
	imageService := service.NewImageProcessingService(cfg, outputDir)

	err = imageService.Start()
	if err != nil {
		log.Fatalf("Failed to start image processing service: %v", err)
	}
}
