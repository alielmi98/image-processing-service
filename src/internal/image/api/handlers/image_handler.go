package handlers

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/alielmi98/image-processing-service/di"
	"github.com/alielmi98/image-processing-service/internal/image/api/dto"
	"github.com/alielmi98/image-processing-service/internal/image/usecase"
	"github.com/alielmi98/image-processing-service/pkg/config"
	"github.com/alielmi98/image-processing-service/pkg/helper"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageHandler struct {
	usecase *usecase.ImageUsecase
}

func NewImageHandler(cfg *config.Config) *ImageHandler {
	return &ImageHandler{
		usecase: usecase.NewImageUsecase(cfg, di.GetImageRepository(cfg)),
	}
}

// CreateImage godoc
// @Summary Create an image
// @Description Create an image
// @Tags Images
// @Accept multipart/form-data
// @produces json
// @Param file formData file true "Image file to upload"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.ImageResponse} "Image response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/images/ [post]
// @Security AuthBearer
func (h *ImageHandler) Create(c *gin.Context) {
	upload := dto.UploadImageRequest{}
	err := c.ShouldBind(&upload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	req := dto.CreateImageRequest{}
	req.MimeType = upload.Image.Header.Get("Content-Type")
	req.FilePath = "uploads"
	req.FileName, req.OriginalName, err = saveUploadedFile(upload.Image, req.FilePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	req.FileSize = upload.Image.Size
	req.Width, req.Height, err = extractImageMetadata(upload.Image)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}

	res, err := h.usecase.CreateImage(c, dto.ToCreateImage(req))
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, helper.Success))

}

func saveUploadedFile(file *multipart.FileHeader, directory string) (fileName, originalName string, err error) {
	allowedExtensions := map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
	}
	// test.txt -> 95239855629856.txt
	randFileName := uuid.New()
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return "", "", err
	}
	fileName = file.Filename
	fileNameArr := strings.Split(fileName, ".")
	originalName = fileNameArr[0]
	fileExt := fileNameArr[len(fileNameArr)-1]
	fileName = fmt.Sprintf("%s.%s", randFileName, fileExt)
	dst := fmt.Sprintf("%s/%s", directory, fileName)

	if !allowedExtensions[fileExt] {
		return "", "", fmt.Errorf("unsupported file format: %s", fileExt)
	}

	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", "", err
	}
	return fileName, originalName, nil
}

func extractImageMetadata(fileHeader *multipart.FileHeader) (width, height int, err error) {
	file, err := fileHeader.Open()
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	file.Seek(0, 0) // Reset pointer
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return config.Width, config.Height, nil
}
