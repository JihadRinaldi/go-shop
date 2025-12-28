package providers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalUploadProvider struct {
	basePath string
}

func NewLocalUploadProvider(basePath string) *LocalUploadProvider {
	return &LocalUploadProvider{basePath: basePath}
}

func (p *LocalUploadProvider) UploadFile(fileHeader *multipart.FileHeader, path string) (string, error) {
	fullPath := filepath.Join(p.basePath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
		return "", err
	}

	// Open the uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Read from source and write to destination
	if _, err := dst.ReadFrom(src); err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/%s", path), nil
}

func (p *LocalUploadProvider) DeleteFile(path string) error {
	fullPath := filepath.Join(p.basePath, path)
	return os.Remove(fullPath)
}
