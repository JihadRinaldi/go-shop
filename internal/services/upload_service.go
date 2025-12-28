package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/JihadRinaldi/go-shop/internal/interfaces"
)

type UploadService struct {
	provider interfaces.UploadProvider
}

func NewUploadService(provider interfaces.UploadProvider) *UploadService {
	return &UploadService{provider: provider}
}

func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))

	if !isValidImageExtension(ext) {
		return "", fmt.Errorf("invalid image extension: %s", ext)
	}

	path := fmt.Sprintf("products/%d/%s", productID, file.Filename)

	return s.provider.UploadFile(file, path)

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
