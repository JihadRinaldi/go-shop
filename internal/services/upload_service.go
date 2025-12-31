package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/JihadRinaldi/go-shop/internal/interfaces"
	"github.com/JihadRinaldi/go-shop/internal/models"
	"github.com/JihadRinaldi/go-shop/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UploadService struct {
	provider   interfaces.UploadProvider
	uploadRepo repositories.UploadRepositoryInterface
}

func NewUploadService(db *gorm.DB, provider interfaces.UploadProvider) *UploadService {
	return &UploadService{
		provider:   provider,
		uploadRepo: repositories.NewUploadRepository(db),
	}
}

func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !isValidImageExtension(ext) {
		return "", fmt.Errorf("invalid image extension: %s", ext)
	}

	newFileName := uuid.New().String()

	path := fmt.Sprintf("products/%d/%s%s", productID, newFileName, ext)

	url, err := s.provider.UploadFile(file, path)
	if err != nil {
		return "", err
	}

	images, _ := s.uploadRepo.GetProductImages(productID)
	isPrimary := len(images) == 0

	image := models.ProductImage{
		ProductID: productID,
		URL:       url,
		AltText:   file.Filename,
		IsPrimary: isPrimary,
	}

	if err := s.uploadRepo.CreateProductImage(&image); err != nil {
		return "", err
	}

	return url, nil
}

func isValidImageExtension(ext string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, v := range validExtensions {
		if ext == v {
			return true
		}
	}
	return false
}
